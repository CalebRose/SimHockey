package managers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	fbsvc "github.com/CalebRose/SimHockey/firebase"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

// ─────────────────────────────────────────────
// CHL Game Request
// ─────────────────────────────────────────────

// CreateCHLGameRequest saves a new CHLGameRequest record to the database.
func CreateCHLGameRequest(request structs.CHLGameRequest) {
	db := dbprovider.GetInstance().GetDB()
	repository.CreateCHLGameRequest(request, db)
}

// AcceptCHLGameRequest marks the request as accepted and notifies the sending
// team's coach if they are a user-managed team.
func AcceptCHLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCHLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	request.Accepted()
	repository.SaveCHLGameRequest(request, db)

	sendingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	if isCHLUserTeam(sendingTeam) {
		receivingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has accepted your game request for Week %d.", receivingTeam.TeamName, request.Week),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "accept", requestID),
		})
	}
}

// RejectCHLGameRequest deletes the request and notifies the sending team's coach
// if they are a user-managed team.
func RejectCHLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCHLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))

	repository.DeleteCHLGameRequest(request, db)

	if isCHLUserTeam(sendingTeam) {
		receivingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has rejected your game request for Week %d.", receivingTeam.TeamName, request.Week),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "reject", requestID),
		})
	}
}

// ProcessCHLGameRequest creates a CollegeGame record from the accepted game
// request and marks the request as approved.
func ProcessCHLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCHLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	homeTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.HomeTeamID)))
	awayTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.AwayTeamID)))

	arenaMap := GetArenaMap()
	arena := arenaMap[request.ArenaID]

	game := structs.CollegeGame{
		BaseGame: structs.BaseGame{
			SeasonID:      request.SeasonID,
			WeekID:        request.WeekID,
			Week:          int(request.Week),
			HomeTeamID:    request.HomeTeamID,
			HomeTeam:      homeTeam.TeamName,
			HomeTeamCoach: homeTeam.Coach,
			AwayTeamID:    request.AwayTeamID,
			AwayTeam:      awayTeam.TeamName,
			AwayTeamCoach: awayTeam.Coach,
			ArenaID:       arena.ID,
			Arena:         arena.Name,
			City:          arena.City,
			State:         arena.State,
			Country:       arena.Country,
			GameDay:       request.Timeslot,
			IsNeutralSite: request.IsNeutralSite,
			IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID,
		},
	}

	repository.CreateCHLGamesRecordsBatch(db, []structs.CollegeGame{game}, 1)

	request.Approved()
	repository.SaveCHLGameRequest(request, db)
}

// VetoCHLGameRequest deletes the request and notifies both teams' coaches
// if either is a user-managed team.
func VetoCHLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCHLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	receivingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))

	repository.DeleteCHLGameRequest(request, db)

	ctx := context.Background()
	msg := fmt.Sprintf("The game request between %s and %s for Week %d has been vetoed.", sendingTeam.TeamName, receivingTeam.TeamName, request.Week)
	vetoKey := fbsvc.BuildSourceEventKey("gamerequest", "chl", "veto", requestID)

	if isCHLUserTeam(sendingTeam) {
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":sending",
		})
	}
	if isCHLUserTeam(receivingTeam) {
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{receivingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         receivingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":receiving",
		})
	}
}

// isCHLUserTeam returns true if the given CollegeTeam is managed by a human coach.
func isCHLUserTeam(team structs.CollegeTeam) bool {
	return team.IsUserCoached && team.Coach != "" && team.Coach != "AI"
}
