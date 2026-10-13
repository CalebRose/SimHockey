package managers

import (
	"context"
	"fmt"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
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

	sendingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	receivingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
	if !isCHLUserTeam(receivingTeam) {
		request.Accepted()
	} else {
		ctx := context.Background()
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{receivingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         receivingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("You have received a new CHL game request from %s.", sendingTeam.TeamName),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "new", strconv.Itoa(int(request.ID))),
		})
	}
	repository.CreateCHLGameRequest(request, db)
}

// CreateCHLInvitationalRequest
func CreateCHLInvitationalRequest(request structs.HockeyInvitationalRequest) {
	db := dbprovider.GetInstance().GetDB()
	repository.CreateInvitationalRequest(request, db)
}

// CancelCHLInvitationalRequest
func CancelCHLInvitationalRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()
	requests := repository.FindInvitationalRequestRecords(repository.SchedulerQuery{ID: requestID})
	for _, request := range requests {
		repository.DeleteInvitationalRequest(request, db)
	}
}

// ProcessHockeyInvitationalRequest -- If there are any existing college games that the user has requested for, they should be removed that week
func ProcessHockeyInvitationalRequest(requestID string) {
	ts := GetTimestamp()
	db := dbprovider.GetInstance().GetDB()
	requests := repository.FindInvitationalRequestRecords(repository.SchedulerQuery{ID: requestID})
	for _, request := range requests {
		request.Approved()
		repository.SaveInvitationalRequest(request, db)
		teamID := strconv.Itoa(int(request.TeamID))
		sendingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.TeamID)))
		invitational := repository.FindInvitationalRecords(repository.SchedulerQuery{ID: strconv.Itoa(int(request.InvitationalID))})
		if len(invitational) == 0 {
			continue
		}
		week := invitational[0].Week
		weekID := util.GetWeekID(ts.SeasonID, uint(week))

		gamesToRemove := []structs.CollegeGame{}
		collegeGames := repository.FindCollegeGames(repository.GamesClauses{WeekID: strconv.Itoa(int(weekID)), TeamID: teamID})
		for _, game := range collegeGames {
			if !game.IsConference && game.WeekID == weekID && (game.HomeTeamID == request.TeamID || game.AwayTeamID == request.TeamID) {
				gamesToRemove = append(gamesToRemove, game)
			}
		}
		if isCHLUserTeam(sendingTeam) {
			receivingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.TeamID)))
			ctx := context.Background()
			uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
			_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
				League:         "chl",
				Domain:         fbsvc.DomainCHL,
				TeamID:         sendingTeam.ID,
				RecipientUIDs:  uids,
				Message:        fmt.Sprintf("%s has accepted your CHL Invitational request.", receivingTeam.TeamName),
				SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "accept", requestID),
			})
		}
	}
}

// VetoHockeyInvitationalRequest marks the request as vetoed and saves it to the database.
func VetoHockeyInvitationalRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()
	requests := repository.FindInvitationalRequestRecords(repository.SchedulerQuery{ID: requestID})
	for _, request := range requests {
		repository.DeleteInvitationalRequest(request, db)
		sendingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.TeamID)))
		if isCHLUserTeam(sendingTeam) {
			receivingTeam := GetCollegeTeamByTeamID(strconv.Itoa(int(request.TeamID)))
			ctx := context.Background()
			uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
			_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
				League:         "chl",
				Domain:         fbsvc.DomainCHL,
				TeamID:         sendingTeam.ID,
				RecipientUIDs:  uids,
				Message:        fmt.Sprintf("%s has vetoed your CHL Invitational request.", receivingTeam.TeamName),
				SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "accept", requestID),
			})
		}
	}
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
			Message:        fmt.Sprintf("%s has accepted your CHL game request for Week %d.", receivingTeam.TeamName, request.Week),
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

	ctx := context.Background()
	if isCHLUserTeam(homeTeam) {
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{homeTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         homeTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("Your game request against %s has been approved.", awayTeam.TeamName),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "new", strconv.Itoa(int(request.ID))),
		})

	}

	if isCHLUserTeam(awayTeam) {
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{awayTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "chl",
			Domain:         fbsvc.DomainCHL,
			TeamID:         awayTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("Your game request against %s has been approved.", homeTeam.TeamName),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "chl", "new", strconv.Itoa(int(request.ID))),
		})
	}

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
