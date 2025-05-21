package managers

import (
	"sort"
	"strconv"
	"sync"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCHLTeamForAvailableTeamsPage(teamID string) structs.TeamRecordResponse {
	historicalDataResponse := GetHistoricalRecordsByTeamID(teamID)

	// Get top 3 players on roster
	roster := repository.FindCollegePlayersByTeamID(teamID)
	sort.Slice(roster, func(i, j int) bool {
		return roster[i].Overall > roster[j].Overall
	})

	topPlayers := []structs.TopPlayer{}

	for i := range roster {
		if i > 4 {
			break
		}
		tp := structs.TopPlayer{}
		tp.MapCollegePlayer(roster[i])
		topPlayers = append(topPlayers, tp)
	}

	historicalDataResponse.AddTopPlayers(topPlayers)

	return historicalDataResponse
}

func GetProTeamForAvailableTeamsPage(teamID string) structs.TeamRecordResponse {
	historicalDataResponse := GetHistoricalProRecordsByTeamID(teamID)

	// Get top 3 players on roster
	roster := repository.FindAllProPlayers(repository.PlayerQuery{TeamID: teamID})
	sort.Slice(roster, func(i, j int) bool {
		return roster[i].Overall > roster[j].Overall
	})

	topPlayers := []structs.TopPlayer{}

	for i := range roster {
		if i > 4 {
			break
		}
		tp := structs.TopPlayer{}
		tp.MapProPlayer(roster[i])
		topPlayers = append(topPlayers, tp)
	}

	historicalDataResponse.AddTopPlayers(topPlayers)

	return historicalDataResponse
}

func GetAllHCKRequests() structs.TeamRequestsResponse {

	var wg sync.WaitGroup
	wg.Add(3)

	var (
		collegeRequests []structs.CollegeTeamRequest
		proRequests     []structs.ProTeamRequest
		acceptedTrades  []structs.TradeProposal
	)

	go func() {
		defer wg.Done()
		collegeRequests = repository.FindAllCHLTeamRequests(true)
	}()

	go func() {
		defer wg.Done()
		proRequests = repository.FindAllPHLTeamRequests(true)
	}()

	go func() {
		defer wg.Done()
		acceptedTrades = repository.FindAllTradeProposalsRecords(repository.TradeClauses{IsAccepted: true, PreloadTradeOptions: true})
	}()

	wg.Wait()

	return structs.TeamRequestsResponse{
		CollegeRequests: collegeRequests,
		ProRequest:      proRequests,
		AcceptedTrades:  acceptedTrades,
	}
}

func CreateCHLTeamRequest(request structs.CollegeTeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	existingRequest := repository.FindCHLRequestRecord(repository.RequestQuery{
		TeamID:   strconv.Itoa(int(request.TeamID)),
		Username: request.Username,
	})

	if existingRequest.ID == 0 {
		repository.CreateCHLTeamRequest(db, request)
		return
	}
	existingRequest.Reactivate()
	repository.SaveCHLTeamRequest(db, existingRequest)
}

func CreatePHLTeamRequest(request structs.ProTeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	existingRequest := repository.FindCHLRequestRecord(repository.RequestQuery{
		TeamID:   strconv.Itoa(int(request.TeamID)),
		Username: request.Username,
		Role:     request.Role,
	})

	if existingRequest.ID == 0 {
		repository.CreatePHLTeamRequest(db, request)
		return
	}
	existingRequest.Reactivate()
	repository.SavePHLTeamRequest(db, structs.ProTeamRequest(existingRequest))
}

func ApproveCHLTeamRequest(request structs.CollegeTeamRequest) structs.CollegeTeamRequest {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	req := repository.FindCHLRequestRecord(repository.RequestQuery{
		ID: strconv.Itoa(int(request.ID)),
	})

	req.ApproveTeamRequest()

	repository.SaveCHLTeamRequest(db, req)

	// Additional changes

	// Team Table
	teamID := strconv.Itoa(int(request.TeamID))
	team := repository.FindCollegeTeamRecord(teamID)
	team.AssignToUser(request.Username)

	repository.SaveCollegeTeamRecord(db, team)

	games := repository.FindCollegeGames(repository.GamesClauses{SeasonID: strconv.Itoa(int(ts.SeasonID)), TeamID: teamID, IsPreseason: false})

	for _, g := range games {
		if g.Week >= int(ts.Week) {
			g.UpdateCoach(req.TeamID, req.Username)
			repository.SaveCollegeGameRecord(g, db)
		}
	}

	recruitingProfile := repository.FindTeamRecruitingProfile(teamID, false, false)

	recruitingProfile.AssignRecruiter(req.Username)
	recruitingProfile.DeactivateAI()

	repository.SaveTeamProfileRecord(db, recruitingProfile)

	CreateNewsLog("CHL", "Breaking News! The "+team.TeamName+" "+team.Mascot+" have hired "+req.Username+" as their new coach for the "+strconv.Itoa(int(ts.Season))+" season!", "CoachJob", 0, ts, true)

	return request
}

func RejectCollegeTeamRequest(request structs.CollegeTeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	req := repository.FindCHLRequestRecord(repository.RequestQuery{
		ID: strconv.Itoa(int(request.ID)),
	})

	req.RejectTeamRequest()

	repository.SaveCHLTeamRequest(db, req)
}

func ApprovePHLTeamRequest(request structs.ProTeamRequest) structs.ProTeamRequest {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	req := repository.FindPHLRequestRecord(repository.RequestQuery{
		ID: strconv.Itoa(int(request.ID)),
	})

	req.ApproveTeamRequest()

	repository.SavePHLTeamRequest(db, req)

	// Team Table
	teamID := strconv.Itoa(int(request.TeamID))
	team := repository.FindProTeamRecord(teamID)
	team.AssignUser(request.Username, request.Role)
	repository.SaveProTeamRecord(db, team)

	CreateNewsLog("CHL", "Breaking News! The "+team.TeamName+" "+team.Mascot+" have hired "+req.Username+" as their new "+request.Role+" for the "+strconv.Itoa(int(ts.Season))+" season!", "CoachJob", 0, ts, true)

	return request
}

func RejectPHLTeamRequest(request structs.ProTeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	req := repository.FindPHLRequestRecord(repository.RequestQuery{
		ID: strconv.Itoa(int(request.ID)),
	})

	req.RejectTeamRequest()

	repository.SavePHLTeamRequest(db, req)
}

func RemoveUserFromCollegeTeam(teamID string) {
	db := dbprovider.GetInstance().GetDB()

	team := repository.FindCollegeTeamRecord(teamID)
	username := team.Coach
	team.RemoveUser()
	team.AssignDiscordID("")

	repository.SaveCollegeTeamRecord(db, team)

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	seasonalGames := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, TeamID: teamID, IsPreseason: false})

	for _, game := range seasonalGames {
		if game.Week >= int(ts.Week) {
			game.UpdateCoach(team.ID, "AI")
			repository.SaveCollegeGameRecord(game, db)
		}

	}

	recruitingProfile := repository.FindTeamRecruitingProfile(teamID, false, false)

	recruitingProfile.ActivateAI()

	repository.SaveTeamProfileRecord(db, recruitingProfile)

	CreateNewsLog("CHL", username+" has decided to step down as the head coach of the "+team.TeamName+" "+team.Mascot+"!", "CoachJob", 0, ts, true)
}

func RemoveUserFromProTeam(request structs.ProTeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	teamID := strconv.Itoa(int(request.TeamID))

	team := repository.FindProTeamRecord(teamID)

	message := ""

	if team.Owner == request.Username {
		message = request.Username + " has decided to step down as Owner of the " + team.TeamName + " " + team.Mascot + "!"
	}

	if team.GM == request.Username {
		message = request.Username + " has decided to step down as Manager of the " + team.TeamName + " " + team.Mascot + "!"
	}

	if team.Coach == request.Username {
		message = request.Username + " has decided to step down as Head Coach of the " + team.TeamName + " " + team.Mascot + "!"
	}

	if team.Scout == request.Username {
		message = request.Username + " has decided to step down as an Assistant of the " + team.TeamName + " " + team.Mascot + "!"
	}

	team.RemoveUser(request.Role)

	repository.SaveProTeamRecord(db, team)

	ts := GetTimestamp()

	CreateNewsLog("PHL", message, "CoachJob", 0, ts, true)
}
