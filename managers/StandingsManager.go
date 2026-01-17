package managers

import (
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func GetCollegeStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(repository.StandingsQuery{
		ConferenceID: conferenceID,
		SeasonID:     seasonID,
	})
}

// HasSeriesBeenClinched returns true if either team has reached the required wins to clinch the series.
func HasSeriesBeenClinched(series structs.BaseSeries) bool {
	requiredWins := (series.BestOfCount / 2) + 1
	return series.HomeTeamWins >= requiredWins || series.AwayTeamWins >= requiredWins
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

func UpdateStandings(ts structs.Timestamp, gameDay string) structs.Timestamp {
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
				if !HasSeriesBeenClinched(series.BaseSeries) && series.GameCount <= series.BestOfCount {
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
						arenaID := 0
						if game.HomeTeamWin {
							homeTeam := collegeTeamMap[HomeID]
							winningTeamID = int(game.HomeTeamID)
							winningTeam = game.HomeTeam
							winningTeamRank = int(game.HomeTeamRank)
							winningCoach = game.HomeTeamCoach
							arenaID = int(homeTeam.ArenaID)
							arena = homeTeam.Arena
							city = homeTeam.City
							state = homeTeam.State
						} else {
							winningTeamID = int(game.AwayTeamID)
							winningTeam = game.AwayTeam
							winningTeamRank = int(game.AwayTeamRank)
							winningCoach = game.AwayTeamCoach
							awayTeam := collegeTeamMap[AwayID]
							arenaID = int(awayTeam.ArenaID)
							arena = awayTeam.Arena
							city = awayTeam.City
							state = awayTeam.State
						}

						nextGame := GetCollegeGameByID(nextGameID)

						nextGame.AddTeam(series.NextGameHOA == "H", uint(winningTeamID), uint(arenaID), uint(winningTeamRank),
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
				arenaID := 0
				if game.HomeTeamWin {
					homeTeam := collegeTeamMap[HomeID]
					winningTeamID = int(game.HomeTeamID)
					winningTeam = game.HomeTeam
					winningTeamRank = int(game.HomeTeamRank)
					winningCoach = game.HomeTeamCoach
					arena = homeTeam.Arena
					city = homeTeam.City
					state = homeTeam.State
					arenaID = int(homeTeam.ArenaID)
				} else {
					winningTeamID = int(game.AwayTeamID)
					winningTeam = game.AwayTeam
					winningTeamRank = int(game.AwayTeamRank)
					winningCoach = game.AwayTeamCoach
					awayTeam := collegeTeamMap[AwayID]
					arena = awayTeam.Arena
					city = awayTeam.City
					state = awayTeam.State
					arenaID = int(awayTeam.ArenaID)
				}

				nextGame := GetCollegeGameByID(nextGameID)

				nextGame.AddTeam(game.NextGameHOA == "H", uint(winningTeamID), uint(arenaID), uint(winningTeamRank),
					winningTeam, winningCoach, arena, city, state)

				repository.SaveCollegeGameRecord(nextGame, db)
			}

			if game.IsNationalChampionship {
				ts.EndTheCollegeSeason()
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

				// Change second half of condition to reflect either best of 7 or best of 5
				if !HasSeriesBeenClinched(series.BaseSeries) && series.GameCount <= series.BestOfCount {
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
					} else {
						switch series.GameCount {
						case 1, 3:
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
						case 2:
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
					} else if series.IsTheFinals && series.NextGameID == 0 {

						ts.EndTheProfessionalSeason()
					}
				}
				repository.SavePlayoffSeriesRecord(series, db)
			}
		}
	}
	return ts
}

func UpdateCollegeRankings() {
	db := dbprovider.GetInstance().GetDB()
	timestamp := GetTimestamp()
	seasonID := strconv.Itoa(int(timestamp.SeasonID - 1))

	// Get all teams and games for the season
	teams := repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})
	standings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: seasonID})
	games := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, IsPreseason: false})

	// Create maps for easier access
	teamMap := MakeCollegeTeamMap(teams)
	// standingsMap := MakeCollegeStandingsMap(standings)

	// Calculate opponent maps for each team
	teamOpponents := calculateTeamOpponents(games)

	// Step 1: Calculate basic win percentages
	teamWinPercentages := calculateWinPercentages(standings)

	// Step 2: Calculate RPI for each team
	teamRPIs := calculateRPI(teamWinPercentages, teamOpponents, standings)

	// Step 3: Calculate SOS (Strength of Schedule)
	teamSOS := calculateSOS(teamRPIs, teamOpponents)

	// Step 4: Calculate SOR (Strength of Record)
	teamSOR := calculateSOR(teamWinPercentages, teamSOS, standings)

	// Step 5: Calculate Quality Win metrics
	calculateQualityWins(teamRPIs, teamOpponents, games, &standings)

	// Step 6: Calculate Conference Strength Adjustments
	conferenceStrengths := calculateConferenceStrengths(teamRPIs, teams)

	// Step 7: Update standings with calculated values
	for i := range standings {
		teamID := standings[i].TeamID

		// Store the actual RPI value and scaled SOS/SOR
		standings[i].RPI = teamRPIs[teamID]
		standings[i].SOS = teamSOS[teamID]
		standings[i].SOR = teamSOR[teamID]

		if conferenceID := teamMap[teamID].ConferenceID; conferenceID > 0 {
			standings[i].ConferenceStrengthAdj = conferenceStrengths[uint(conferenceID)]
		}

		// Save updated standings
		repository.SaveCollegeStandingsRecord(standings[i], db)
	}

	// Step 8: Calculate and assign final rankings
	assignRPIRanks(teamRPIs, &standings, db)
	assignPairwiseRanks(teamRPIs, teamOpponents, games, &standings, db)
}

// Helper function to calculate team opponents from games
func calculateTeamOpponents(games []structs.CollegeGame) map[uint][]uint {
	opponents := make(map[uint][]uint)

	for _, game := range games {
		if !game.GameComplete {
			continue
		}

		homeID := game.HomeTeamID
		awayID := game.AwayTeamID

		opponents[homeID] = append(opponents[homeID], awayID)
		opponents[awayID] = append(opponents[awayID], homeID)
	}

	return opponents
}

// Calculate win percentages for all teams
func calculateWinPercentages(standings []structs.CollegeStandings) map[uint]float32 {
	winPercentages := make(map[uint]float32)

	for _, standing := range standings {
		totalGames := standing.TotalWins + standing.TotalLosses + standing.TotalOTLosses
		if totalGames > 0 {
			// In hockey, OT losses typically count as 0.5 wins for RPI
			adjustedWins := float32(standing.TotalWins) + (float32(standing.TotalOTLosses) * 0.5)
			winPercentages[standing.TeamID] = adjustedWins / float32(totalGames)
		} else {
			winPercentages[standing.TeamID] = 0.0
		}
	}

	return winPercentages
}

// Calculate RPI: (0.25 × Win%) + (0.50 × Opponents' Win%) + (0.25 × Opponents' Opponents' Win%)
func calculateRPI(winPercentages map[uint]float32, teamOpponents map[uint][]uint, standings []structs.CollegeStandings) map[uint]float32 {
	teamRPIs := make(map[uint]float32)

	for _, standing := range standings {
		teamID := standing.TeamID
		teamWinPct := winPercentages[teamID]

		// Calculate opponents' win percentage
		var opponentsWinPct float32
		opponents := teamOpponents[teamID]
		if len(opponents) > 0 {
			for _, oppID := range opponents {
				opponentsWinPct += winPercentages[oppID]
			}
			opponentsWinPct /= float32(len(opponents))
		}

		// Calculate opponents' opponents' win percentage
		var opponentsOpponentsWinPct float32
		totalOppOpps := 0
		for _, oppID := range opponents {
			oppOpponents := teamOpponents[oppID]
			for _, oppOppID := range oppOpponents {
				if oppOppID != teamID { // Don't count self
					opponentsOpponentsWinPct += winPercentages[oppOppID]
					totalOppOpps++
				}
			}
		}
		if totalOppOpps > 0 {
			opponentsOpponentsWinPct /= float32(totalOppOpps)
		}

		// Calculate final RPI
		rpi := (0.25 * teamWinPct) + (0.50 * opponentsWinPct) + (0.25 * opponentsOpponentsWinPct)
		teamRPIs[teamID] = rpi
	}

	return teamRPIs
}

// Calculate SOS: Average RPI of all opponents
func calculateSOS(teamRPIs map[uint]float32, teamOpponents map[uint][]uint) map[uint]float32 {
	teamSOS := make(map[uint]float32)

	for teamID, opponents := range teamOpponents {
		var totalOpponentRPI float32

		if len(opponents) > 0 {
			for _, oppID := range opponents {
				totalOpponentRPI += teamRPIs[oppID]
			}
			teamSOS[teamID] = totalOpponentRPI / float32(len(opponents))
		} else {
			teamSOS[teamID] = 0.0
		}
	}

	return teamSOS
}

// Calculate SOR: Expected wins given schedule difficulty
func calculateSOR(winPercentages map[uint]float32, teamSOS map[uint]float32, standings []structs.CollegeStandings) map[uint]float32 {
	teamSOR := make(map[uint]float32)

	for _, standing := range standings {
		teamID := standing.TeamID
		actualWinPct := winPercentages[teamID]
		scheduleStrength := teamSOS[teamID]

		// Simple SOR calculation: actual performance vs expected performance based on schedule
		// Higher SOS means harder schedule, so good record vs hard schedule = higher SOR
		if scheduleStrength > 0 {
			teamSOR[teamID] = actualWinPct / scheduleStrength
		} else {
			teamSOR[teamID] = actualWinPct
		}
	}

	return teamSOR
}

// Calculate quality wins and bad losses
func calculateQualityWins(teamRPIs map[uint]float32, teamOpponents map[uint][]uint, games []structs.CollegeGame, standings *[]structs.CollegeStandings) {
	// Create sorted list of teams by RPI for tier determination
	type teamRPI struct {
		teamID uint
		rpi    float32
	}

	var sortedRPIs []teamRPI
	for teamID, rpi := range teamRPIs {
		sortedRPIs = append(sortedRPIs, teamRPI{teamID, rpi})
	}

	// Sort by RPI descending
	for i := 0; i < len(sortedRPIs)-1; i++ {
		for j := i + 1; j < len(sortedRPIs); j++ {
			if sortedRPIs[i].rpi < sortedRPIs[j].rpi {
				sortedRPIs[i], sortedRPIs[j] = sortedRPIs[j], sortedRPIs[i]
			}
		}
	}

	// Determine tier thresholds (top 15 = Tier 1, 16-30 = Tier 2, 45+ = bad loss threshold)
	tier1Threshold := 15
	tier2Threshold := 30
	badLossThreshold := 45

	tier1Teams := make(map[uint]bool)
	tier2Teams := make(map[uint]bool)
	badLossTeams := make(map[uint]bool)

	for i, team := range sortedRPIs {
		if i < tier1Threshold {
			tier1Teams[team.teamID] = true
		} else if i < tier2Threshold {
			tier2Teams[team.teamID] = true
		} else if i >= badLossThreshold {
			badLossTeams[team.teamID] = true
		}
	}

	// Create head-to-head record map for enhanced quality win analysis
	headToHeadMap := make(map[uint]map[uint]struct{ wins, losses int })
	for _, game := range games {
		if !game.GameComplete {
			continue
		}

		homeID := game.HomeTeamID
		awayID := game.AwayTeamID

		if headToHeadMap[homeID] == nil {
			headToHeadMap[homeID] = make(map[uint]struct{ wins, losses int })
		}
		if headToHeadMap[awayID] == nil {
			headToHeadMap[awayID] = make(map[uint]struct{ wins, losses int })
		}

		if game.HomeTeamWin {
			rec := headToHeadMap[homeID][awayID]
			rec.wins++
			headToHeadMap[homeID][awayID] = rec

			rec = headToHeadMap[awayID][homeID]
			rec.losses++
			headToHeadMap[awayID][homeID] = rec
		} else {
			rec := headToHeadMap[awayID][homeID]
			rec.wins++
			headToHeadMap[awayID][homeID] = rec

			rec = headToHeadMap[homeID][awayID]
			rec.losses++
			headToHeadMap[homeID][awayID] = rec
		}
	}

	// Count quality wins and bad losses for each team
	for i := range *standings {
		standing := &(*standings)[i]
		teamID := standing.TeamID

		// Reset quality metrics before recalculating
		standing.Tier1Wins = 0
		standing.Tier2Wins = 0
		standing.BadLosses = 0

		for _, game := range games {
			if !game.GameComplete {
				continue
			}

			var isWin bool
			var opponentID uint

			if game.HomeTeamID == teamID {
				isWin = game.HomeTeamWin
				opponentID = game.AwayTeamID
			} else if game.AwayTeamID == teamID {
				isWin = game.AwayTeamWin
				opponentID = game.HomeTeamID
			} else {
				continue
			}

			if isWin {
				if tier1Teams[opponentID] {
					// Count each individual win against tier 1 teams
					standing.Tier1Wins++
				} else if tier2Teams[opponentID] {
					standing.Tier2Wins++
				}
			} else {
				if badLossTeams[opponentID] {
					// Count each individual bad loss
					standing.BadLosses++
				}
			}
		}
	}
}

// Helper function to find common opponents between two teams
func getCommonOpponents(teamA, teamB uint, teamOpponents map[uint][]uint) []uint {
	opponentsA := make(map[uint]bool)
	for _, opp := range teamOpponents[teamA] {
		opponentsA[opp] = true
	}

	var commonOpponents []uint
	for _, opp := range teamOpponents[teamB] {
		if opp != teamA && opp != teamB && opponentsA[opp] {
			commonOpponents = append(commonOpponents, opp)
		}
	}

	return commonOpponents
}

// Calculate conference strengths
func calculateConferenceStrengths(teamRPIs map[uint]float32, teams []structs.CollegeTeam) map[uint]float32 {
	conferenceRPIs := make(map[uint][]float32)
	conferenceStrengths := make(map[uint]float32)

	// Group teams by conference
	for _, team := range teams {
		confID := uint(team.ConferenceID)
		if rpi, exists := teamRPIs[team.ID]; exists {
			conferenceRPIs[confID] = append(conferenceRPIs[confID], rpi)
		}
	}

	// Calculate average RPI for each conference
	for confID, rpis := range conferenceRPIs {
		var total float32
		for _, rpi := range rpis {
			total += rpi
		}
		if len(rpis) > 0 {
			conferenceStrengths[confID] = total / float32(len(rpis))
		}
	}

	return conferenceStrengths
}

// Assign RPI rankings
func assignRPIRanks(teamRPIs map[uint]float32, standings *[]structs.CollegeStandings, db *gorm.DB) {
	type teamRPI struct {
		teamID uint
		rpi    float32
	}

	var sortedRPIs []teamRPI
	for teamID, rpi := range teamRPIs {
		sortedRPIs = append(sortedRPIs, teamRPI{teamID, rpi})
	}

	// Sort by RPI descending
	for i := 0; i < len(sortedRPIs)-1; i++ {
		for j := i + 1; j < len(sortedRPIs); j++ {
			if sortedRPIs[i].rpi < sortedRPIs[j].rpi {
				sortedRPIs[i], sortedRPIs[j] = sortedRPIs[j], sortedRPIs[i]
			}
		}
	}

	// Assign RPI ranks
	for i := range *standings {
		standing := &(*standings)[i]
		for rank, teamRPI := range sortedRPIs {
			if teamRPI.teamID == standing.TeamID {
				standing.RPIRank = uint8(rank + 1)
				repository.SaveCollegeStandingsRecord(*standing, db)
				break
			}
		}
	}
}

// Assign Pairwise rankings using true NCAA Pairwise methodology
func assignPairwiseRanks(teamRPIs map[uint]float32, teamOpponents map[uint][]uint, games []structs.CollegeGame, standings *[]structs.CollegeStandings, db *gorm.DB) {
	// Create head-to-head records
	headToHeadMap := make(map[uint]map[uint]struct{ wins, losses int })
	for _, game := range games {
		if !game.GameComplete {
			continue
		}

		homeID := game.HomeTeamID
		awayID := game.AwayTeamID

		if headToHeadMap[homeID] == nil {
			headToHeadMap[homeID] = make(map[uint]struct{ wins, losses int })
		}
		if headToHeadMap[awayID] == nil {
			headToHeadMap[awayID] = make(map[uint]struct{ wins, losses int })
		}

		if game.HomeTeamWin {
			rec := headToHeadMap[homeID][awayID]
			rec.wins++
			headToHeadMap[homeID][awayID] = rec

			rec = headToHeadMap[awayID][homeID]
			rec.losses++
			headToHeadMap[awayID][homeID] = rec
		} else {
			rec := headToHeadMap[awayID][homeID]
			rec.wins++
			headToHeadMap[awayID][homeID] = rec

			rec = headToHeadMap[homeID][awayID]
			rec.losses++
			headToHeadMap[homeID][awayID] = rec
		}
	}

	// Calculate common opponent records for each team pair
	commonOpponentWins := make(map[uint]map[uint]int)
	for _, standing := range *standings {
		teamA := standing.TeamID
		commonOpponentWins[teamA] = make(map[uint]int)

		for _, otherStanding := range *standings {
			teamB := otherStanding.TeamID
			if teamA == teamB {
				continue
			}

			// Find common opponents
			commonOpps := getCommonOpponents(teamA, teamB, teamOpponents)
			if len(commonOpps) < 3 {
				continue // Need at least 3 common opponents for meaningful comparison
			}

			// Count wins against common opponents
			teamAWins := 0
			teamBWins := 0

			for _, commonOpp := range commonOpps {
				if h2h := headToHeadMap[teamA][commonOpp]; h2h.wins > h2h.losses {
					teamAWins++
				}
				if h2h := headToHeadMap[teamB][commonOpp]; h2h.wins > h2h.losses {
					teamBWins++
				}
			}

			commonOpponentWins[teamA][teamB] = teamAWins - teamBWins
		}
	}

	// Perform pairwise comparisons
	type pairwiseResult struct {
		teamID uint
		wins   int
	}

	var pairwiseResults []pairwiseResult
	for _, standingA := range *standings {
		teamA := standingA.TeamID
		wins := 0

		for _, standingB := range *standings {
			teamB := standingB.TeamID
			if teamA == teamB {
				continue
			}

			// Pairwise comparison criteria (in order of priority):
			// 1. Head-to-head record (if teams have played)
			if h2h := headToHeadMap[teamA][teamB]; h2h.wins+h2h.losses > 0 {
				if h2h.wins > h2h.losses {
					wins++
				}
				continue
			}

			// 2. Record against common opponents (if sufficient common opponents)
			if commonWinDiff, exists := commonOpponentWins[teamA][teamB]; exists {
				if commonWinDiff > 0 {
					wins++
				}
				continue
			}

			// 3. RPI comparison
			if teamRPIs[teamA] > teamRPIs[teamB] {
				wins++
			}
		}

		pairwiseResults = append(pairwiseResults, pairwiseResult{teamA, wins})
	}

	// Sort by pairwise wins (descending)
	for i := 0; i < len(pairwiseResults)-1; i++ {
		for j := i + 1; j < len(pairwiseResults); j++ {
			if pairwiseResults[i].wins < pairwiseResults[j].wins {
				pairwiseResults[i], pairwiseResults[j] = pairwiseResults[j], pairwiseResults[i]
			}
		}
	}

	// Assign Pairwise ranks based on comparison wins
	for i := range *standings {
		standing := &(*standings)[i]
		for rank, result := range pairwiseResults {
			if result.teamID == standing.TeamID {
				standing.PairwiseRank = uint8(rank + 1)
				repository.SaveCollegeStandingsRecord(*standing, db)
				break
			}
		}
	}
}

// GetCollegeRankingsWithMetrics returns all college standings with ranking metrics for a season
func GetCollegeRankingsWithMetrics(seasonID string) []structs.CollegeStandings {
	standings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: seasonID})

	// Sort by Pairwise Rank (primary ranking system)
	for i := 0; i < len(standings)-1; i++ {
		for j := i + 1; j < len(standings); j++ {
			if standings[i].PairwiseRank > standings[j].PairwiseRank {
				standings[i], standings[j] = standings[j], standings[i]
			}
		}
	}

	return standings
}

// GetTop25CollegeRankings returns the top 25 teams by Pairwise ranking
func GetTop25CollegeRankings(seasonID string) []structs.CollegeStandings {
	standings := GetCollegeRankingsWithMetrics(seasonID)

	if len(standings) > 25 {
		return standings[:25]
	}

	return standings
}

func CreatePreseasonRanking() {
	db := dbprovider.GetInstance().GetDB()
	timestamp := GetTimestamp()
	currentSeasonID := strconv.Itoa(int(timestamp.SeasonID))
	previousSeasonID := strconv.Itoa(int(timestamp.SeasonID - 1))

	// Get all teams for current season
	teams := repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})

	// Get previous season data for baseline
	previousStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: previousSeasonID})

	// Get current season standings (should be empty/reset for preseason)
	currentStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: currentSeasonID})

	// currentStandingsMap := MakeCollegeStandingsMap(currentStandings)

	previousStandingsMap := MakeCollegeStandingsMap(previousStandings)

	// TODO: Implement these data sources as needed:
	// 1. Get roster data with player ratings
	players := repository.FindAllCollegePlayers(repository.PlayerQuery{})

	collegeRosterMap := MakeCollegePlayerMapByTeamID(players)

	preseasonScores := calculatePreseasonScores(teams, previousStandingsMap, collegeRosterMap)

	assignPreseasonRanks(preseasonScores, &currentStandings, db)
}

// calculatePreseasonScores computes preseason ranking scores based on multiple factors
func calculatePreseasonScores(teams []structs.CollegeTeam, previousStandings map[uint]structs.CollegeStandings, collegeRosterMap map[uint][]structs.CollegePlayer) map[uint]float32 {
	scores := make(map[uint]float32)

	for _, team := range teams {
		var score float32 = 0.0

		// Factor 1: Previous Season Performance (40% weight)
		if prevStanding, exists := previousStandings[team.ID]; exists {
			// Base score on previous win percentage and RPI
			winPct := prevStanding.GetWinPercentage()
			rpiScore := prevStanding.RPI

			prevSeasonScore := (winPct * 0.6) + (rpiScore * 0.4)
			score += prevSeasonScore * 0.40
		} else {
			// New program or missing data - use neutral score
			score += 0.5 * 0.40
		}

		// Factor 2: Program Prestige (20% weight)
		prestigeScore := float32(team.ProgramPrestige) / 10.0
		score += prestigeScore * 0.20

		// Factor 3: Roster Talent (30% weight)
		roster := collegeRosterMap[team.ID]
		rosterScore := calculateRosterTalent(roster)
		score += rosterScore * 0.30

		// Factor 4: Coaching Stability (10% weight)
		// CoachRating is a number between 1 and 10
		coachingScore := float32(team.CoachRating) / 10.0
		score += coachingScore * 0.10

		scores[team.ID] = score
	}

	return scores
}

// calculateRosterTalent evaluates the overall talent level of a team's roster
func calculateRosterTalent(roster []structs.CollegePlayer) float32 {
	if len(roster) == 0 {
		return 0.3 // Below average for teams with no roster data
	}

	var totalTalent float32
	var playerCount int
	var starterBonus float32

	// Position weights for roster balance evaluation
	positionCounts := make(map[string]int)

	for _, player := range roster {
		// Use player's Overall rating as primary talent metric
		// Assuming Overall is 1-50, normalize to 0-1
		playerTalent := float32(player.Overall) / 50.0

		// Weight by player year (experience matters)
		yearMultiplier := float32(1.0)
		switch player.Year {
		case 1: // Freshman
			yearMultiplier = 0.8
		case 2: // Sophomore
			yearMultiplier = 0.9
		case 3: // Junior
			yearMultiplier = 1.0
		case 4: // Senior
			yearMultiplier = 1.1
		case 5: // Graduate/5th year
			yearMultiplier = 1.15
		}

		// Apply experience multiplier
		adjustedTalent := playerTalent * yearMultiplier

		// Bonus for star players (45+ overall)
		if player.Overall >= 45 {
			starterBonus += 0.05 // Each elite player adds 5% bonus
		}

		totalTalent += adjustedTalent
		playerCount++

		// Track position balance
		positionCounts[player.Position]++
	}

	// Calculate average roster talent
	avgTalent := totalTalent / float32(playerCount)

	// Roster depth bonus/penalty
	depthModifier := float32(1.0)
	if playerCount < 18 { // Thin roster penalty
		depthModifier = 0.85
	} else if playerCount > 25 { // Good depth bonus
		depthModifier = 1.1
	}

	// Position balance bonus (ensure reasonable distribution)
	balanceModifier := calculatePositionBalance(positionCounts)

	// Final roster score with all modifiers
	finalScore := (avgTalent + starterBonus) * depthModifier * balanceModifier

	// Ensure score stays within reasonable bounds (0.0 to 1.0)
	if finalScore > 1.0 {
		finalScore = 1.0
	}
	if finalScore < 0.0 {
		finalScore = 0.0
	}

	return finalScore
}

// calculatePositionBalance evaluates roster balance across positions
func calculatePositionBalance(positionCounts map[string]int) float32 {
	// Expected minimum players by position for hockey
	expectedMins := map[string]int{
		"C": 2, // Centers
		"F": 2, // Forwards
		"D": 3, // Defensemen
		"G": 2, // Goalies
	}

	balanceScore := float32(1.0)

	for position, expectedMin := range expectedMins {
		actual := positionCounts[position]
		if actual < expectedMin {
			// Penalty for being short at a position
			shortfall := float32(expectedMin - actual)
			balanceScore -= shortfall * 0.05 // 5% penalty per missing player
		}
	}

	// Ensure balance score doesn't go below 0.7 (max 30% penalty)
	if balanceScore < 0.7 {
		balanceScore = 0.7
	}

	return balanceScore
}

// assignPreseasonRanks assigns preseason rankings based on calculated scores
func assignPreseasonRanks(preseasonScores map[uint]float32, standings *[]structs.CollegeStandings, db *gorm.DB) {
	type teamScore struct {
		teamID uint
		score  float32
	}

	var sortedScores []teamScore
	for teamID, score := range preseasonScores {
		sortedScores = append(sortedScores, teamScore{teamID, score})
	}

	// Sort by score descending
	for i := 0; i < len(sortedScores)-1; i++ {
		for j := i + 1; j < len(sortedScores); j++ {
			if sortedScores[i].score < sortedScores[j].score {
				sortedScores[i], sortedScores[j] = sortedScores[j], sortedScores[i]
			}
		}
	}

	// Assign preseason ranks
	for i := range *standings {
		standing := &(*standings)[i]
		for rank, teamScore := range sortedScores {
			if teamScore.teamID == standing.TeamID {
				// Set preseason ranking - you might want a separate PreseasonRank field
				standing.PreseasonRank = uint8(rank + 1)

				// Initialize RPI with preseason score for display
				standing.RPI = teamScore.score

				repository.SaveCollegeStandingsRecord(*standing, db)
				break
			}
		}
	}
}
