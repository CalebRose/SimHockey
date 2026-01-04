package managers

import (
	"strconv"
	"sync"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCHLPlayByPlayStreamData(streamType string) []structs.StreamResponse {
	ts := GetTimestamp()
	teamMap := GetCollegeTeamMap()
	rosterMap := GetAllCollegePlayersMapByTeam()
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()
	games := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)

	streams := []structs.StreamResponse{}

	for _, game := range games {
		if !game.GameComplete || (ts.Week == 18 && game.ID < 2633) {
			continue
		}
		homeTeam := teamMap[uint(game.HomeTeamID)]
		awayTeam := teamMap[uint(game.AwayTeamID)]

		if streamType == "1" {
			if !homeTeam.IsUserCoached && !awayTeam.IsUserCoached {
				continue
			}
			mod := game.ID % 2
			if mod == 0 {
				continue
			}
		}
		if streamType == "2" {
			if !homeTeam.IsUserCoached && !awayTeam.IsUserCoached {
				continue
			}
			mod := game.ID % 2
			if mod == 1 {
				continue
			}
		}

		if streamType == "3" {
			if homeTeam.IsUserCoached || awayTeam.IsUserCoached {
				continue
			}
		}

		gameID := strconv.Itoa(int(game.ID))
		var wg sync.WaitGroup
		var (
			playByPlays []structs.CollegePlayByPlay
			homePlayers []structs.CollegePlayer
			awayPlayers []structs.CollegePlayer
			homeStats   []structs.CollegePlayerGameStats
			awayStats   []structs.CollegePlayerGameStats
		)
		homePlayers = rosterMap[game.HomeTeamID]
		awayPlayers = rosterMap[game.AwayTeamID]
		wg.Add(2)

		go func() {
			defer wg.Done()
			stats := repository.FindCollegePlayerStatsRecordByGame(gameID)
			for _, s := range stats {
				if s.TeamID == game.HomeTeamID {
					homeStats = append(homeStats, s)
				} else {
					awayStats = append(awayStats, s)
				}
			}
		}()

		go func() {
			defer wg.Done()
			playByPlays = GetCHLPlayByPlaysByGameID(gameID)
		}()

		wg.Wait()

		totalList := []structs.CollegePlayer{}
		totalList = append(totalList, homePlayers...)
		totalList = append(totalList, awayPlayers...)

		participantMap := MakeCollegePlayerMap(totalList)
		playbyPlayResponse := GenerateCHLPlayByPlayResponse(playByPlays, teamMap, participantMap, true, game.HomeTeamID, game.AwayTeamID)

		stream := structs.StreamResponse{
			GameID:            game.ID,
			HomeTeamID:        uint(game.HomeTeamID),
			HomeTeam:          game.HomeTeam,
			HomeTeamCoach:     homeTeam.Coach,
			HomeTeamRank:      game.HomeTeamRank,
			HomeLabel:         homeTeam.TeamName + " " + homeTeam.Mascot,
			HomeTeamDiscordID: homeTeam.DiscordID,
			AwayTeamID:        uint(game.AwayTeamID),
			AwayTeam:          game.AwayTeam,
			AwayTeamCoach:     awayTeam.Coach,
			AwayTeamRank:      game.AwayTeamRank,
			AwayTeamDiscordID: awayTeam.DiscordID,
			AwayLabel:         awayTeam.TeamName + " " + awayTeam.Mascot,
			Streams:           playbyPlayResponse,
			City:              game.City,
			State:             game.State,
			Country:           game.Country,
			ArenaID:           game.ArenaID,
			Arena:             game.Arena,
			Attendance:        uint(game.AttendanceCount),
		}

		streams = append(streams, stream)
	}

	return streams
}

func GetPHLPlayByPlayStreamData(streamType string) []structs.StreamResponse {
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()
	games := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)
	teamMap := GetProTeamMap()
	rosterMap := GetAllProPlayersMapByTeam()
	streams := []structs.StreamResponse{}

	for _, game := range games {
		if !game.GameComplete {
			continue
		}

		homeTeam := teamMap[uint(game.HomeTeamID)]
		awayTeam := teamMap[uint(game.AwayTeamID)]
		if streamType == "1" {
			if homeTeam.Owner == "" && awayTeam.Owner == "" {
				continue
			}
			mod := game.ID % 2
			if mod == 1 {
				continue
			}
		}
		if streamType == "2" {
			if homeTeam.Owner == "" && awayTeam.Owner == "" {
				continue
			}
			mod := game.ID % 2
			if mod == 0 {
				continue
			}
		}

		if streamType == "3" {
			if homeTeam.Owner != "" || awayTeam.Owner != "" {
				continue
			}

		}

		gameID := strconv.Itoa(int(game.ID))
		var wg sync.WaitGroup
		var (
			playByPlays []structs.ProPlayByPlay
			homeStats   []structs.ProfessionalPlayerGameStats
			awayStats   []structs.ProfessionalPlayerGameStats
		)
		homePlayers := rosterMap[game.HomeTeamID]
		awayPlayers := rosterMap[game.AwayTeamID]

		wg.Add(2)

		go func() {
			defer wg.Done()
			stats := repository.FindProPlayerStatsRecordByGame(gameID)
			for _, s := range stats {
				if s.TeamID == game.HomeTeamID {
					homeStats = append(homeStats, s)
				} else {
					awayStats = append(awayStats, s)
				}
			}
		}()

		go func() {
			defer wg.Done()
			playByPlays = GetPHLPlayByPlaysByGameID(gameID)
		}()

		wg.Wait()

		totalList := []structs.ProfessionalPlayer{}
		totalList = append(totalList, homePlayers...)
		totalList = append(totalList, awayPlayers...)

		participantMap := MakeProfessionalPlayerMap(totalList)
		playbyPlayResponse := GeneratePHLPlayByPlayResponse(playByPlays, teamMap, participantMap, true, game.HomeTeamID, game.AwayTeamID)

		stream := structs.StreamResponse{
			GameID:            game.ID,
			HomeTeamID:        uint(game.HomeTeamID),
			HomeTeam:          game.HomeTeam,
			HomeTeamCoach:     game.HomeTeamCoach,
			HomeTeamDiscordID: homeTeam.DiscordID,
			HomeLabel:         game.HomeTeam,
			AwayTeamID:        uint(game.AwayTeamID),
			AwayTeam:          game.AwayTeam,
			AwayTeamCoach:     game.AwayTeamCoach,
			AwayLabel:         game.AwayTeam,
			AwayTeamDiscordID: awayTeam.DiscordID,
			Streams:           playbyPlayResponse,
			City:              game.City,
			State:             game.State,
		}

		streams = append(streams, stream)
	}

	return streams
}

func GetAllCollegePlayersWithGameStatsByTeamID(GameID string, stats []structs.CollegePlayerGameStats) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()
	ids := []string{}
	statMap := make(map[uint]structs.CollegePlayerGameStats)
	for _, s := range stats {
		playerID := strconv.Itoa(int(s.PlayerID))
		ids = append(ids, playerID)
		statMap[uint(s.PlayerID)] = s
	}

	var collegePlayers []structs.CollegePlayer
	var matchRows []structs.CollegePlayer

	db.Where("id in (?)", ids).Find(&collegePlayers)

	for _, p := range collegePlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.TimeOnIce == 0 {
			continue
		}

		matchRows = append(matchRows, p)
	}

	historicPlayers := []structs.HistoricCollegePlayer{}
	db.Where("id in (?)", ids).Find(&historicPlayers)

	for _, p := range historicPlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.TimeOnIce == 0 {
			continue
		}

		row := structs.CollegePlayer{
			Model:      p.Model,
			BasePlayer: p.BasePlayer,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}

func GetAllPHLPlayersWithGameStatsByTeamID(GameID string, stats []structs.ProfessionalPlayerGameStats) []structs.ProfessionalPlayer {
	db := dbprovider.GetInstance().GetDB()

	ids := []string{}
	statMap := make(map[uint]structs.ProfessionalPlayerGameStats)
	for _, s := range stats {
		playerID := strconv.Itoa(int(s.PlayerID))
		ids = append(ids, playerID)
		statMap[uint(s.PlayerID)] = s
	}

	var proPlayers []structs.ProfessionalPlayer
	var matchRows []structs.ProfessionalPlayer

	db.Where("id in (?)", ids).Find(&proPlayers)

	for _, p := range proPlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.TimeOnIce == 0 {
			continue
		}

		matchRows = append(matchRows, p)
	}

	historicPlayers := []structs.RetiredPlayer{}
	db.Where("id in (?)", ids).Find(&historicPlayers)

	for _, p := range historicPlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.TimeOnIce == 0 {
			continue
		}
		row := structs.ProfessionalPlayer{
			Model:      p.Model,
			BasePlayer: p.BasePlayer,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}

func generateCollegeResultsString(play structs.PbP, event, outcome string, playerMap map[uint]structs.CollegePlayer, possessingTeam structs.CollegeTeam) string {
	puckCarrier := playerMap[play.PuckCarrierID]
	puckCarrierLabel := getPlayerLabel(puckCarrier.BasePlayer)
	receivingPlayer := playerMap[play.PassedPlayerID]
	receivingPlayerLabel := getPlayerLabel(receivingPlayer.BasePlayer)
	assistingPlayer := playerMap[play.AssistingPlayerID]
	assistingPlayerLabel := getPlayerLabel(assistingPlayer.BasePlayer)
	defendingPlayer := playerMap[play.DefenderID]
	defendingPlayerLabel := getPlayerLabel(defendingPlayer.BasePlayer)
	goalie := playerMap[play.GoalieID]
	goalieLabel := getPlayerLabel(goalie.BasePlayer)
	statement := ""
	nextZoneLabel := getZoneLabel(play.NextZoneID)
	teamLabel := possessingTeam.TeamName
	// First Segment
	switch event {
	case Faceoff:
		switch outcome {
		case "Home Faceoff Win":
			statement = puckCarrierLabel + " wins the faceoff! "
		case util.GoalieHold:
			statement = puckCarrierLabel + " holds onto the puck, and it's going to a faceoff."
		default:
			statement = receivingPlayerLabel + " wins the faceoff! "
		}
		// Mention receiving player
		if outcome != util.GoalieHold {
			statement += assistingPlayerLabel + " receives the puck on the faceoff."
		}
	case PhysDefenseCheck:
		switch outcome {
		case DefenseTakesPuck:
			statement = defendingPlayerLabel + " bodies " + puckCarrierLabel + " right into the boards and snatches the puck away!"
		case CarrierKeepsPuck:
			statement = defendingPlayerLabel + " attempts to body right into " + puckCarrierLabel + ", but " + puckCarrierLabel + " maneuvers effortlessly within the zone!"
		}
	case DexDefenseCheck:
		switch outcome {
		case DefenseTakesPuck:
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		case CarrierKeepsPuck:
			statement = defendingPlayerLabel + " attempts to swipe the puck from " + puckCarrierLabel + ", but his stick is batted away!"
		}
	case PassCheck:
		switch outcome {
		case InterceptedPass:
			statement = defendingPlayerLabel + " intercepts the pass right from " + puckCarrierLabel + "!"
		case ReceivedPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the pass!"
		case ReceivedLongPass, ReceivedBackPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " in the " + nextZoneLabel + " and makes the pass!"
		}
	case LongPassCheck:
		switch outcome {
		case InterceptedPass:
			statement = defendingPlayerLabel + " intercepts the long pass right from " + puckCarrierLabel + "!"
		case ReceivedPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the long pass in the " + nextZoneLabel + "!"
		case ReceivedLongPass, ReceivedBackPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " in the " + nextZoneLabel + " and makes the pass!"
		}
	case PassBackCheck:
		switch outcome {
		case InterceptedPass:
			statement = defendingPlayerLabel + " intercepts the back pass right from " + puckCarrierLabel + "!"
		case ReceivedPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the pass back to the " + nextZoneLabel + "!"
		case ReceivedLongPass, ReceivedBackPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " in the " + nextZoneLabel + " and makes the pass!"
		}
	case AgilityCheck:
		switch outcome {
		case DefenseStopAgility:
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		case OffenseMovesUp:
			statement = puckCarrierLabel + " moves the puck up to the " + nextZoneLabel + "."
		}
	case WristshotCheck:
		statement = puckCarrierLabel + " attempts a long shot on goal..."
		switch outcome {
		case ShotBlocked:
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		case GoalieSave:
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		case InAccurateShot:
			statement += " and he misses the goal! It's a loose puck! Picked up by " + receivingPlayerLabel + "!"
		case ShotOnGoal:
			statement += " and he scores! That's a point for " + teamLabel + "!"
		}
	case SlapshotCheck:
		statement = puckCarrierLabel + " attempts a slapshot on goal..."
		switch outcome {
		case ShotBlocked:
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		case GoalieSave:
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		case InAccurateShot:
			if !play.IsShootout {
				statement += " and he misses the goal! It's a loose puck! Picked up by " + receivingPlayerLabel + "!"
			} else {
				statement += " and he misses the net!"
			}
		case ShotOnGoal:
			statement += " and he scores! That's a point for " + teamLabel + "!"
		case PenaltyCheck:
			penalty := getPenaltyByID(uint(play.PenaltyID))
			severity := getSeverityByID(play.Severity)
			penaltyMinutes := "two"
			if play.Severity > 1 {
				penaltyMinutes = "five"
			}
			statement += " and a penalty is called! " + defendingPlayerLabel + " has been called for a " + severity + " " + penalty + " on " + puckCarrierLabel + ". This will lead into a faceoff. Power play for " + penaltyMinutes + " minutes."
		}
	case PenaltyCheck:
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		penaltyMinutes := "two"
		if play.Severity > 1 {
			penaltyMinutes = "five"
		}
		if play.IsFight {
			statement = "There's a fight on center ice! " + defendingPlayerLabel + " and " + goalieLabel + " are right at with the fisticuffs. Refs are breaking up the fight. Both players are out for " + penaltyMinutes + " minutes. Resetting play with a faceoff. "
		} else {
			statement = "Penalty called! " + defendingPlayerLabel + " has been called for " + severity + " " + penalty + " on " + puckCarrierLabel + ". Power play for " + penaltyMinutes + " minutes."
		}
	case EnteringShootout:
		statement = "END OF OVERTIME, STARTING SHOOTOUT"
	case Shootout:
		statement = puckCarrierLabel + " faces " + goalieLabel + " in the shootout..."
		switch outcome {
		case GoalieSave:
			statement += " and the shot is SAVED by " + goalieLabel + "! The next player is up!"
		case ShotOnGoal:
			statement += " and he scores! That's a point for " + teamLabel + "!"
		case InAccurateShot:
			statement += " and he misses the net! What an inaccurate shot!"
		}

	case PuckBattle:
		statement = "Puck battle between " + receivingPlayerLabel + " and " + defendingPlayerLabel + "!"
		switch outcome {
		case PuckBattleWin:
			statement += " " + puckCarrierLabel + " retains the puck!"
		case PuckBattleLose:
			statement += " " + defendingPlayerLabel + " comes out with the puck!"
		}

	case PuckScramble:
		statement = "The puck is loose! It's a scramble for the puck!"

	case util.InjuryCheck:
		injury := ""
		injuryType := ""
		injuryDuration := ""
		if event == util.InjuryCheck {
			injury = util.GetInjuryNameByID(play.InjuryID)
			injuryType = util.GetInjurySeverityByID(play.InjuryType)
			injuryDuration = strconv.Itoa(int(play.InjuryDuration)) + " games"
		}
		statement = "Injury Update: " + puckCarrierLabel + " has sustained a " + injuryType + " injury (" + injury + ") and will be out for approximately " + injuryDuration + "."
	}

	return statement
}

func generateProResultsString(play structs.PbP, event, outcome string, playerMap map[uint]structs.ProfessionalPlayer, possessingTeam structs.ProfessionalTeam) string {
	puckCarrier := playerMap[play.PuckCarrierID]
	puckCarrierLabel := getPlayerLabel(puckCarrier.BasePlayer)
	receivingPlayer := playerMap[play.PassedPlayerID]
	receivingPlayerLabel := getPlayerLabel(receivingPlayer.BasePlayer)
	assistingPlayer := playerMap[play.AssistingPlayerID]
	assistingPlayerLabel := getPlayerLabel(assistingPlayer.BasePlayer)
	defendingPlayer := playerMap[play.DefenderID]
	defendingPlayerLabel := getPlayerLabel(defendingPlayer.BasePlayer)
	goalie := playerMap[play.GoalieID]
	goalieLabel := getPlayerLabel(goalie.BasePlayer)
	statement := ""
	nextZoneLabel := getZoneLabel(play.NextZoneID)
	teamLabel := possessingTeam.TeamName
	// First Segment
	switch event {
	case Faceoff:
		switch outcome {
		case "Home Faceoff Win":
			statement = puckCarrierLabel + " wins the faceoff! "
		case util.GoalieHold:
			statement = puckCarrierLabel + " holds onto the puck, and it's going to a faceoff."
		default:
			statement = receivingPlayerLabel + " wins the faceoff! "
		}
		// Mention receiving player
		if outcome != util.GoalieHold {
			statement += assistingPlayerLabel + " receives the puck on the faceoff."
		}
	case PhysDefenseCheck:
		switch outcome {
		case DefenseTakesPuck:
			statement = defendingPlayerLabel + " bodies " + puckCarrierLabel + " right into the boards and snatches the puck away!"
		case CarrierKeepsPuck:
			statement = defendingPlayerLabel + " attempts to body right into " + puckCarrierLabel + ", but " + puckCarrierLabel + " maneuvers effortlessly within the zone!"
		}
	case DexDefenseCheck:
		switch outcome {
		case DefenseTakesPuck:
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		case CarrierKeepsPuck:
			statement = defendingPlayerLabel + " attempts to swipe the puck from " + puckCarrierLabel + ", but his stick is batted away!"
		}
	case PassCheck:
		switch outcome {
		case InterceptedPass:
			statement = defendingPlayerLabel + " intercepts the pass right from " + puckCarrierLabel + "!"
		case ReceivedPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the pass!"
		case ReceivedLongPass, ReceivedBackPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " in the " + nextZoneLabel + " and makes the pass!"
		}
	case LongPassCheck:
		switch outcome {
		case InterceptedPass:
			statement = defendingPlayerLabel + " intercepts the long pass right from " + puckCarrierLabel + "!"
		case ReceivedPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the long pass in the " + nextZoneLabel + "!"
		case ReceivedLongPass, ReceivedBackPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " in the " + nextZoneLabel + " and makes the pass!"
		}
	case PassBackCheck:
		switch outcome {
		case InterceptedPass:
			statement = defendingPlayerLabel + " intercepts the back pass right from " + puckCarrierLabel + "!"
		case ReceivedPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the pass back to the " + nextZoneLabel + "!"
		case ReceivedLongPass, ReceivedBackPass:
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " in the " + nextZoneLabel + " and makes the pass!"
		}
	case AgilityCheck:
		switch outcome {
		case DefenseStopAgility:
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		case OffenseMovesUp:
			statement = puckCarrierLabel + " moves the puck up to the " + nextZoneLabel + "."
		}
	case WristshotCheck:
		statement = puckCarrierLabel + " attempts a long shot on goal..."
		switch outcome {
		case ShotBlocked:
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		case GoalieSave:
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		case InAccurateShot:
			statement += " and he misses the goal! It's a loose puck! Picked up by " + receivingPlayerLabel + "!"
		case ShotOnGoal:
			statement += " and he scores! That's a point for " + teamLabel + "!"
		}
	case SlapshotCheck:
		statement = puckCarrierLabel + " attempts a slapshot on goal..."
		switch outcome {
		case ShotBlocked:
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		case GoalieSave:
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		case InAccurateShot:
			if !play.IsShootout {
				statement += " and he misses the goal! It's a loose puck! Picked up by " + receivingPlayerLabel + "!"
			} else {
				statement += " and he misses the net!"
			}
		case ShotOnGoal:
			statement += " and he scores! That's a point for " + teamLabel + "!"
		case PenaltyCheck:
			penalty := getPenaltyByID(uint(play.PenaltyID))
			severity := getSeverityByID(play.Severity)
			penaltyMinutes := "two"
			if play.Severity > 1 {
				penaltyMinutes = "five"
			}
			statement += " and a penalty is called! " + defendingPlayerLabel + " has been called for a " + severity + " " + penalty + " on " + puckCarrierLabel + ". This will lead into a faceoff. Power play for " + penaltyMinutes + " minutes."
		}
	case PenaltyCheck:
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		penaltyMinutes := "two"
		if play.Severity > 1 {
			penaltyMinutes = "five"
		}
		if play.IsFight {
			statement = "There's a fight on center ice! " + defendingPlayerLabel + " and " + goalieLabel + " are right at with the fisticuffs. Refs are breaking up the fight. Both players are out for " + penaltyMinutes + " minutes. Resetting play with a faceoff. "
		} else {
			statement = "Penalty called! " + defendingPlayerLabel + " has been called for " + severity + " " + penalty + " on " + puckCarrierLabel + ". Power play for " + penaltyMinutes + " minutes."
		}
	case EnteringShootout:
		statement = "END OF OVERTIME, STARTING SHOOTOUT"
	case Shootout:
		statement = puckCarrierLabel + " faces " + goalieLabel + " in the shootout..."
		switch outcome {
		case GoalieSave:
			statement += " and the shot is SAVED by " + goalieLabel + "! The next player is up!"
		case ShotOnGoal:
			statement += " and he scores! That's a point for " + teamLabel + "!"
		case InAccurateShot:
			statement += " and he misses the net! What an inaccurate shot!"
		}

	case PuckBattle:
		statement = "Puck battle between " + receivingPlayerLabel + " and " + defendingPlayerLabel + "!"
		switch outcome {
		case PuckBattleWin:
			statement += " " + puckCarrierLabel + " retains the puck!"
		case PuckBattleLose:
			statement += " " + defendingPlayerLabel + " comes out with the puck!"
		}

	case PuckScramble:
		statement = "The puck is loose! It's a scramble for the puck!"

	case util.InjuryCheck:
		injury := ""
		injuryType := ""
		injuryDuration := ""
		if event == util.InjuryCheck {
			injury = util.GetInjuryNameByID(play.InjuryID)
			injuryType = util.GetInjurySeverityByID(play.InjuryType)
			injuryDuration = strconv.Itoa(int(play.InjuryDuration)) + " games"
		}
		statement = "Injury Update: " + puckCarrierLabel + " has sustained a " + injuryType + " injury (" + injury + ") and will be out for approximately " + injuryDuration + "."

	}

	return statement
}

func getPlayerLabel(player structs.BasePlayer) string {
	if len(player.FirstName) == 0 {
		return ""
	}
	return player.Team + " " + player.Position + " " + player.FirstName + " " + player.LastName
}

func getZoneLabel(zoneID uint8) string {
	if zoneID == 0 {
		return ""
	}
	if zoneID == HomeGoalZoneID {
		return HomeGoal
	}
	if zoneID == HomeZoneID {
		return HomeZone
	}
	if zoneID == NeutralZoneID {
		return NeutralZone
	}
	if zoneID == AwayZoneID {
		return AwayZone
	}
	if zoneID == AwayGoalZoneID {
		return AwayGoal
	}
	return ""
}

func getPenaltyByID(penaltyType uint) string {
	var penaltyMap = map[uint]string{
		1:  "Aggressor Penalty",
		2:  "Attempt to Injure",
		3:  "Biting",
		4:  "Boarding",
		5:  "Boarding",
		6:  "Stabbing",
		7:  "Charging",
		8:  "Charging",
		9:  "Checking from Behind",
		10: "Checking from Behind",
		11: "Clipping",
		12: "Clipping",
		13: "Cross Checking",
		14: "Cross Checking",
		15: "Delay of Game",
		16: "Diving",
		17: "Elbowing",
		18: "Elbowing",
		19: "Eye-Gouging",
		20: "Fighting",
		21: "Goaltender Interference",
		22: "Headbutting",
		23: "High-sticking",
		24: "High-sticking",
		25: "Holding",
		26: "Hooking",
		27: "Hooking",
		28: "Kicking",
		29: "Kicking",
		30: "Kneeing",
		31: "Kneeing",
		32: "Roughing",
		33: "Roughing",
		34: "Slashing",
		35: "Slashing",
		36: "Slew footing",
		37: "Slew footing",
		38: "Throwing the stick",
		39: "Too many men on the ice",
		40: "Tripping",
		41: "Tripping",
		42: "Unsportsmanlike conduct",
	}
	return penaltyMap[penaltyType]
}

func getSeverityByID(sevId uint8) string {
	var severityMap = map[uint8]string{
		1: "Minor Penalty",
		2: "Major Penalty",
		3: "Game Misconduct",
		4: "Match Penalty",
	}
	return severityMap[sevId]
}
