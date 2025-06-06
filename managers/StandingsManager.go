package managers

import (
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(seasonID, conferenceID, "")
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

func GetProStandingsMap(seasonID string) map[uint]structs.ProfessionalStandings {
	standings := repository.FindAllProfessionalStandings(seasonID, "", "")
	return MakeProfessionalStandingsMap(standings)
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

	historicGames := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, TeamID: TeamID, IsPreseason: false})
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
	historicGames := repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, TeamID: TeamID, IsPreseason: false})
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

func UpdateStandings(ts structs.Timestamp, gameDay string) {
	db := dbprovider.GetInstance().GetDB()

	if !ts.IsOffSeason {
		weekID := strconv.Itoa(int(ts.WeekID))
		seasonID := strconv.Itoa(int(ts.SeasonID))
		collegeStandingsMap := GetCollegeStandingsMap(seasonID)
		collegeTeamMap := GetCollegeTeamMap()
		proTeamMap := GetProTeamMap()
		proStandingsMap := GetProStandingsMap(seasonID)

		collegeGames := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)

		for _, game := range collegeGames {
			if !game.GameComplete || game.IsPreseason {
				continue
			}

			HomeID := game.HomeTeamID
			AwayID := game.AwayTeamID
			homeStandings := collegeStandingsMap[HomeID]
			awayStandings := collegeStandingsMap[AwayID]

			homeStandings.UpdateStandings(game.BaseGame)
			awayStandings.UpdateStandings(game.BaseGame)

			repository.SaveCollegeStandingsRecord(homeStandings, db)
			repository.SaveCollegeStandingsRecord(awayStandings, db)

			if game.NextGameID > 0 {

				nextGameID := strconv.Itoa(int(game.NextGameID))
				winningTeamID := 0
				winningTeam := ""
				winningCoach := ""
				winningTeamRank := 0
				arena := ""
				city := ""
				state := ""
				if game.HomeTeamWin {
					homeTeam := collegeTeamMap[HomeID]
					winningTeamID = int(game.HomeTeamID)
					winningTeam = game.HomeTeam
					winningTeamRank = int(game.HomeTeamRank)
					winningCoach = game.HomeTeamCoach
					arena = homeTeam.Arena
					city = homeTeam.City
					state = homeTeam.State
				} else {
					winningTeamID = int(game.AwayTeamID)
					winningTeam = game.AwayTeam
					winningTeamRank = int(game.AwayTeamRank)
					winningCoach = game.AwayTeamCoach
					awayTeam := collegeTeamMap[AwayID]
					arena = awayTeam.Arena
					city = awayTeam.City
					state = awayTeam.State
				}

				nextGame := GetCollegeGameByID(nextGameID)

				nextGame.AddTeam(game.NextGameHOA == "H", uint(winningTeamID), uint(winningTeamRank),
					winningTeam, winningCoach, arena, city, state)

				repository.SaveCollegeGameRecord(nextGame, db)
			}

			if game.IsNationalChampionship {
				ts.EndTheCollegeSeason()
				repository.SaveTimestamp(ts, db)
			}
		}

		proGames := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)

		for _, game := range proGames {
			if !game.GameComplete || game.IsPreseason {
				continue
			}

			HomeID := game.HomeTeamID
			AwayID := game.AwayTeamID
			if !game.IsPlayoffGame {
				homeStandings := proStandingsMap[HomeID]
				awayStandings := proStandingsMap[AwayID]

				homeStandings.UpdateStandings(game.BaseGame)
				awayStandings.UpdateStandings(game.BaseGame)

				repository.SaveProfessionalStandingsRecord(homeStandings, db)
				repository.SaveProfessionalStandingsRecord(awayStandings, db)
			}

			if game.IsPlayoffGame && game.SeriesID > 0 {
				seriesID := strconv.Itoa(int(game.SeriesID))

				series := GetPlayoffSeriesBySeriesID(seriesID)

				winningID := 0
				if game.HomeTeamWin {
					winningID = int(game.HomeTeamID)
				} else {
					winningID = int(game.AwayTeamID)
				}
				series.UpdateWinCount(winningID)

				if series.GameCount <= 7 && (series.HomeTeamWins < 4 && series.AwayTeamWins < 4) {
					homeTeamID := 0
					nextHomeTeam := ""
					nextHomeTeamCoach := ""
					nextHomeRank := 0
					awayTeamID := 0
					nextAwayTeam := ""
					nextAwayTeamCoach := ""
					nextAwayRank := 0
					city := ""
					arena := ""
					state := ""
					country := ""
					if series.GameCount == 1 || series.GameCount == 2 || series.GameCount == 5 || series.GameCount == 7 {
						homeTeam := proTeamMap[series.HomeTeamID]
						homeTeamID = int(series.HomeTeamID)
						nextHomeTeam = series.HomeTeam
						nextHomeTeamCoach = series.HomeTeamCoach
						nextHomeRank = int(series.HomeTeamRank)
						city = homeTeam.City
						arena = homeTeam.Arena
						state = homeTeam.State
						country = homeTeam.Country
						awayTeamID = int(series.AwayTeamID)
						nextAwayTeam = series.AwayTeam
						nextAwayTeamCoach = series.AwayTeamCoach
						nextAwayRank = int(series.AwayTeamRank)
					} else if series.GameCount == 3 || series.GameCount == 4 || series.GameCount == 6 {
						awayTeam := proTeamMap[series.AwayTeamID]
						homeTeamID = int(series.AwayTeamID)
						nextHomeTeam = series.AwayTeam
						nextHomeTeamCoach = series.AwayTeamCoach
						nextHomeRank = int(series.AwayTeamRank)
						city = awayTeam.City
						arena = awayTeam.Arena
						state = awayTeam.State
						country = awayTeam.Country
						awayTeamID = int(series.HomeTeamID)
						nextAwayTeam = series.HomeTeam
						nextAwayTeamCoach = series.HomeTeamCoach
						nextAwayRank = int(series.HomeTeamRank)
					}
					weekID := ts.WeekID
					week := ts.Week
					matchOfWeek := "A"
					if game.GameDay == "A" {
						matchOfWeek = "B"
					} else if game.GameDay == "B" {
						matchOfWeek = "C"
					} else if game.GameDay == "C" {
						matchOfWeek = "D"
					} else if game.GameDay == "D" {
						// Move game to next week
						weekID += 1
						week += 1
					}
					matchTitle := series.SeriesName + ": " + nextHomeTeam + " vs. " + nextAwayTeam
					nextGame := structs.ProfessionalGame{
						BaseGame: structs.BaseGame{
							WeekID:        weekID,
							Week:          int(week),
							SeasonID:      ts.SeasonID,
							GameDay:       matchOfWeek,
							GameTitle:     matchTitle,
							HomeTeamID:    uint(homeTeamID),
							HomeTeam:      nextHomeTeam,
							HomeTeamCoach: nextHomeTeamCoach,
							HomeTeamRank:  uint(nextHomeRank),
							AwayTeamID:    uint(awayTeamID),
							AwayTeam:      nextAwayTeam,
							AwayTeamCoach: nextAwayTeamCoach,
							AwayTeamRank:  uint(nextAwayRank),
							City:          city,
							Arena:         arena,
							State:         state,
							Country:       country,
							IsPlayoffGame: true,
						},
						SeriesID:        series.ID,
						IsInternational: series.IsInternational,
					}
					repository.CreatePHLGamesRecordsBatch(db, []structs.ProfessionalGame{nextGame}, 1)
				} else {
					if !series.IsTheFinals && series.NextSeriesID > 0 {
						// Promote Team to Next Series
						nextSeriesID := strconv.Itoa(int(series.NextSeriesID))
						nextSeriesHoa := series.NextSeriesHOA
						nextSeries := GetPlayoffSeriesBySeriesID(nextSeriesID)
						var teamID uint = 0
						teamLabel := ""
						teamCoach := ""
						teamRank := 0
						if series.HomeTeamWin {
							teamID = series.HomeTeamID
							teamLabel = series.HomeTeam
							teamCoach = series.HomeTeamCoach
							teamRank = int(series.HomeTeamRank)
						} else {
							teamID = series.AwayTeamID
							teamLabel = series.AwayTeam
							teamCoach = series.AwayTeamCoach
							teamRank = int(series.AwayTeamRank)
						}
						nextSeries.AddTeam(nextSeriesHoa == "H", teamID, uint(teamRank), teamLabel, teamCoach)
						repository.SavePlayoffSeriesRecord(nextSeries, db)
					}
				}
				repository.SavePlayoffSeriesRecord(series, db)
			}
		}
	}
}
