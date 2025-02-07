package managers

import (
	"strconv"

	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(seasonID, conferenceID, "")
}

func GetProfessionalStandingsBySeasonID(seasonID string) []structs.ProfessionalStandings {
	return repository.FindAllProfessionalStandings(seasonID, "", "")
}

func GetAllCollegeStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(seasonID, "", "")
}

func GetAllProfessionalStandingsBySeasonID(seasonID string) []structs.ProfessionalStandings {
	return repository.FindAllProfessionalStandings(seasonID, "", "")
}

func GetCollegeStandingsMap(seasonID string) map[uint]structs.CollegeStandings {
	standings := repository.FindAllCollegeStandings(seasonID, "", "")
	return MakeCollegeStandingsMap(standings)
}

// GetHistoricalRecordsByTeamID
func GetHistoricalRecordsByTeamID(TeamID string) structs.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	seasonID := strconv.Itoa(int(timestamp.SeasonID))

	historicGames := repository.FindCollegeGames(seasonID, TeamID)
	var conferenceChampionships []string
	var divisionTitles []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	playoffWins := 0
	playoffLosses := 0

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.SeasonID && game.WeekID == timestamp.WeekID) {
			continue
		}
		winningSeason := game.SeasonID + 2024
		winningSeasonStr := strconv.Itoa(int(winningSeason))
		isAway := strconv.Itoa(int(game.AwayTeamID)) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonWins++
			}

			if game.IsPlayoffGame {
				playoffWins++
			}

			if game.IsNationalChampionship {
				nationalChampionships = append(nationalChampionships, winningSeasonStr)
			}
		} else {
			overallLosses++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonLosses++
			}

			if game.IsPlayoffGame {
				playoffLosses++
			}
		}
	}

	response := structs.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		PostSeasonWins:          playoffWins,
		PostSeasonLosses:        playoffLosses,
		ConferenceChampionships: conferenceChampionships,
		DivisionTitles:          divisionTitles,
		NationalChampionships:   nationalChampionships,
	}

	return response
}

func GetHistoricalProRecordsByTeamID(TeamID string) structs.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)
	seasonID := strconv.Itoa(int(timestamp.SeasonID))
	historicGames := repository.FindProfessionalGames(seasonID, TeamID)
	var conferenceChampionships []string
	var divisionTitles []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	playoffWins := 0
	playoffLosses := 0

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.SeasonID && game.WeekID == timestamp.WeekID) {
			continue
		}
		gameSeason := game.SeasonID + 2020
		isAway := strconv.Itoa(int(game.AwayTeamID)) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonWins++
			}

			if game.IsPlayoffGame {
				playoffWins++
			}

			if game.IsStanleyCup {
				nationalChampionships = append(nationalChampionships, strconv.Itoa(int(gameSeason)))
			}
		} else {
			overallLosses++

			if game.SeasonID == timestamp.SeasonID {
				currentSeasonLosses++
			}

			if game.IsPlayoffGame {
				playoffLosses++
			}
		}
	}

	response := structs.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		PostSeasonWins:          playoffWins,
		PostSeasonLosses:        playoffLosses,
		ConferenceChampionships: conferenceChampionships,
		DivisionTitles:          divisionTitles,
		NationalChampionships:   nationalChampionships,
	}

	return response
}
