package managers

import (
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(repository.StandingsQuery{
		ConferenceID: conferenceID,
		SeasonID:     seasonID,
	})
}

func GetAllCollegeStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(repository.StandingsQuery{
		SeasonID: seasonID,
	})
}

func GetAllProfessionalStandingsBySeasonID(seasonID string) []structs.ProfessionalStandings {
	return repository.FindAllProfessionalStandings(repository.StandingsQuery{
		SeasonID: seasonID,
	})
}

func GetCollegeStandingsMap(seasonID string) map[uint]structs.CollegeStandings {
	standings := repository.FindAllCollegeStandings(repository.StandingsQuery{
		SeasonID: seasonID,
	})
	return MakeCollegeStandingsMap(standings)
}

func GetProStandingsMap(seasonID string) map[uint]structs.ProfessionalStandings {
	standings := repository.FindAllProfessionalStandings(repository.StandingsQuery{
		SeasonID: seasonID,
	})
	return MakeProfessionalStandingsMap(standings)
}

func GetStandingsHistoryByTeamID(id string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(repository.StandingsQuery{
		TeamID: id,
	})
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
			homeStandings.UpdateSeasonStatus(game)
			awayStandings.UpdateSeasonStatus(game)

			repository.SaveCollegeStandingsRecord(homeStandings, db)
			repository.SaveCollegeStandingsRecord(awayStandings, db)

			if game.SeriesID > 0 {
				seriesID := strconv.Itoa(int(game.SeriesID))

				series := repository.FindCollegeSeriesRecord(seriesID)

				winningID := 0
				if game.HomeTeamWin {
					winningID = int(game.HomeTeamID)
				} else {
					winningID = int(game.AwayTeamID)
				}
				series.UpdateWinCount(winningID)

				if series.GameCount <= series.BestOfCount && (series.HomeTeamWins < 2 && series.AwayTeamWins < 2) {
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
					arenaID := 0
					state := ""
					country := ""
					if series.BestOfCount == 5 || series.BestOfCount == 7 {
						switch series.GameCount {
						case 1, 2, 5, 7:
							homeTeam := collegeTeamMap[series.HomeTeamID]
							homeTeamID = int(series.HomeTeamID)
							nextHomeTeam = series.HomeTeam
							nextHomeTeamCoach = series.HomeTeamCoach
							nextHomeRank = int(series.HomeTeamRank)
							city = homeTeam.City
							arena = homeTeam.Arena
							arenaID = int(homeTeam.ArenaID)
							state = homeTeam.State
							country = homeTeam.Country
							awayTeamID = int(series.AwayTeamID)
							nextAwayTeam = series.AwayTeam
							nextAwayTeamCoach = series.AwayTeamCoach
							nextAwayRank = int(series.AwayTeamRank)
						case 3, 4, 6:
							awayTeam := collegeTeamMap[series.AwayTeamID]
							homeTeamID = int(series.AwayTeamID)
							nextHomeTeam = series.AwayTeam
							nextHomeTeamCoach = series.AwayTeamCoach
							nextHomeRank = int(series.AwayTeamRank)
							city = awayTeam.City
							arena = awayTeam.Arena
							arenaID = int(awayTeam.ArenaID)
							state = awayTeam.State
							country = awayTeam.Country
							awayTeamID = int(series.HomeTeamID)
							nextAwayTeam = series.HomeTeam
							nextAwayTeamCoach = series.HomeTeamCoach
							nextAwayRank = int(series.HomeTeamRank)
						}
					} else {
						switch series.GameCount {
						case 1, 3:
							homeTeam := collegeTeamMap[series.HomeTeamID]
							homeTeamID = int(series.HomeTeamID)
							nextHomeTeam = series.HomeTeam
							nextHomeTeamCoach = series.HomeTeamCoach
							nextHomeRank = int(series.HomeTeamRank)
							city = homeTeam.City
							arena = homeTeam.Arena
							arenaID = int(homeTeam.ArenaID)
							state = homeTeam.State
							country = homeTeam.Country
							awayTeamID = int(series.AwayTeamID)
							nextAwayTeam = series.AwayTeam
							nextAwayTeamCoach = series.AwayTeamCoach
							nextAwayRank = int(series.AwayTeamRank)
						case 2:
							awayTeam := collegeTeamMap[series.AwayTeamID]
							homeTeamID = int(series.AwayTeamID)
							nextHomeTeam = series.AwayTeam
							nextHomeTeamCoach = series.AwayTeamCoach
							nextHomeRank = int(series.AwayTeamRank)
							city = awayTeam.City
							arena = awayTeam.Arena
							arenaID = int(awayTeam.ArenaID)
							state = awayTeam.State
							country = awayTeam.Country
							awayTeamID = int(series.HomeTeamID)
							nextAwayTeam = series.HomeTeam
							nextAwayTeamCoach = series.HomeTeamCoach
							nextAwayRank = int(series.HomeTeamRank)
						default:
							// Should not happen
						}

					}
					weekID := ts.WeekID
					week := ts.Week
					matchOfWeek := "A"
					switch game.GameDay {
					case "A":
						matchOfWeek = "B"
					case "B":
						matchOfWeek = "C"
					case "C":
						matchOfWeek = "D"
					case "D":
						// Move game to next week
						weekID += 1
						week += 1
					}
					matchTitle := series.SeriesName + ": " + nextHomeTeam + " vs. " + nextAwayTeam
					nextGame := structs.CollegeGame{
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
							ArenaID:       uint(arenaID),
							State:         state,
							Country:       country,
							IsPlayoffGame: false,
							SeriesID:      series.ID,
						},
						IsConferenceTournament: true,
					}
					repository.CreateCHLGamesRecordsBatch(db, []structs.CollegeGame{nextGame}, 1)
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
						if series.HomeTeamSeriesWin {
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
					} else if series.NextGameID > 0 {

						nextGameID := strconv.Itoa(int(series.NextGameID))
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
				}
				repository.SaveCollegeSeriesRecord(series, db)
			} else if game.NextGameID > 0 {

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
					arenaID := 0
					state := ""
					country := ""
					switch series.GameCount {
					case 1, 2, 5, 7:
						homeTeam := proTeamMap[series.HomeTeamID]
						homeTeamID = int(series.HomeTeamID)
						nextHomeTeam = series.HomeTeam
						nextHomeTeamCoach = series.HomeTeamCoach
						nextHomeRank = int(series.HomeTeamRank)
						city = homeTeam.City
						arena = homeTeam.Arena
						arenaID = int(homeTeam.ArenaID)
						state = homeTeam.State
						country = homeTeam.Country
						awayTeamID = int(series.AwayTeamID)
						nextAwayTeam = series.AwayTeam
						nextAwayTeamCoach = series.AwayTeamCoach
						nextAwayRank = int(series.AwayTeamRank)
					case 3, 4, 6:
						awayTeam := proTeamMap[series.AwayTeamID]
						homeTeamID = int(series.AwayTeamID)
						nextHomeTeam = series.AwayTeam
						nextHomeTeamCoach = series.AwayTeamCoach
						nextHomeRank = int(series.AwayTeamRank)
						city = awayTeam.City
						arena = awayTeam.Arena
						arenaID = int(awayTeam.ArenaID)
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
					switch game.GameDay {
					case "A":
						matchOfWeek = "B"
					case "B":
						matchOfWeek = "C"
					case "C":
						matchOfWeek = "D"
					case "D":
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
							ArenaID:       uint(arenaID),
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
