package managers

import (
	"strconv"
	"sync"

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
		if !game.GameComplete || game.ID == 2449 || game.ID == 2451 || game.ID == 2452 {
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
			Model:          p.Model,
			BasePlayer:     p.BasePlayer,
			BaseInjuryData: p.BaseInjuryData,
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
			Model:          p.Model,
			BasePlayer:     p.BasePlayer,
			BaseInjuryData: p.BaseInjuryData,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}
