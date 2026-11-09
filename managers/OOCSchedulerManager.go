package managers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func GeneratePermanentInvitationals(ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()
	seasonID := ts.SeasonID
	collegeTeams := GetAllCollegeTeams()
	collegeTeamMap := MakeCollegeTeamMap(collegeTeams)
	arenas := repository.FindAllArenas(repository.ArenaQuery{})
	arenaMap := MakeArenaMap(arenas)
	baseID := repository.FindLatestCHLGameID() + 1
	collegeGamesBatch := []structs.CollegeGame{}
	collegeStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: strconv.Itoa(int(seasonID))})
	standingsMap := MakeCollegeStandingsMap(collegeStandings)

	hockeyInvitationals := repository.FindInvitationalRecords(repository.SchedulerQuery{})

	mk := func(id uint, invId uint8, title string, week int, winningGameID, losingGameID uint, hoa string, arena structs.Arena) structs.CollegeGame {
		return structs.CollegeGame{
			Model:          gorm.Model{ID: id},
			IsInvitational: true,
			InvitationalID: invId,
			LosingGameID:   losingGameID,
			LosingGameHOA:  hoa,
			BaseGame: structs.BaseGame{
				GameTitle:     title,
				SeasonID:      seasonID,
				WeekID:        util.GetWeekID(seasonID, uint(week)),
				Week:          week,
				IsNeutralSite: true,
				LeagueID:      1,
				ArenaID:       arena.ID,
				Arena:         arena.Name,
				City:          arena.City,
				State:         arena.State,
				Country:       arena.Country,
				NextGameID:    winningGameID,
				NextGameHOA:   hoa,
			},
		}
	}

	for _, invitational := range hockeyInvitationals {
		if invitational.ID != 1 && invitational.ID != 3 {
			continue
		}
		participatingArena := arenaMap[uint(invitational.ArenaID)]

		participatingTeams := []structs.CollegeTeam{}
		switch invitational.ID {
		case 1:
			bostonU := collegeTeamMap[11]
			bostonCollege := collegeTeamMap[10]
			harvard := collegeTeamMap[23]
			northeastern := collegeTeamMap[43]
			participatingTeams = append(participatingTeams, bostonU, bostonCollege, harvard, northeastern)
		case 3:
			UConn := collegeTeamMap[18]
			quinnipiac := collegeTeamMap[51]
			sacredHeart := collegeTeamMap[55]
			yale := collegeTeamMap[64]
			participatingTeams = append(participatingTeams, UConn, quinnipiac, sacredHeart, yale)
		}

		ShuffleTeams(participatingTeams)

		// Game 1
		game1 := mk(baseID, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 1", ts.Season, invitational.Name), int(invitational.Week), baseID+2, baseID+3, "h", participatingArena)
		game1.HomeTeam = participatingTeams[0].Abbreviation
		game1.HomeTeamID = participatingTeams[0].ID
		game1.HomeTeamCoach = participatingTeams[0].Coach
		game1.HomeTeamRank = standingsMap[participatingTeams[0].ID].Rank
		game1.AwayTeam = participatingTeams[1].Abbreviation
		game1.AwayTeamID = participatingTeams[1].ID
		game1.AwayTeamCoach = participatingTeams[1].Coach
		game1.AwayTeamRank = standingsMap[participatingTeams[1].ID].Rank
		game1.GameDay = "B"
		// Game 2
		game2 := mk(baseID+1, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 2", ts.Season, invitational.Name), int(invitational.Week), baseID+2, baseID+3, "a", participatingArena)
		game2.HomeTeam = participatingTeams[2].Abbreviation
		game2.HomeTeamID = participatingTeams[2].ID
		game2.HomeTeamCoach = participatingTeams[2].Coach
		game2.HomeTeamRank = standingsMap[participatingTeams[2].ID].Rank
		game2.AwayTeam = participatingTeams[3].Abbreviation
		game2.AwayTeamID = participatingTeams[3].ID
		game2.AwayTeamCoach = participatingTeams[3].Coach
		game2.AwayTeamRank = standingsMap[participatingTeams[3].ID].Rank
		game2.GameDay = "B"
		// Finals Game
		finals := mk(baseID+2, uint8(invitational.ID), fmt.Sprintf("%d %s - Finals", ts.Season, invitational.Name), int(invitational.Week), 0, 0, "", participatingArena)
		finals.GameDay = "C"
		// Consolidation Game
		consolation := mk(baseID+3, uint8(invitational.ID), fmt.Sprintf("%d %s - Consolation", ts.Season, invitational.Name), int(invitational.Week), 0, 0, "", participatingArena)
		consolation.GameDay = "C"

		collegeGamesBatch = append(collegeGamesBatch, game1, game2, finals, consolation)

		// Iterate BaseID
		baseID = baseID + 4
	}

	repository.CreateCHLGamesRecordsBatch(db, collegeGamesBatch, 20)
}

func GenerateOOCSchedule() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	collegeTeams := GetAllCollegeTeams()

	maxAttempts := 5000
	var matchesToUpload []structs.CollegeGame
	var err error

	// Try multiple times to generate a complete schedule
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		matchesToUpload, err = attemptGenerateOOCSchedule(ts, seasonID, collegeTeams)
		if err == nil {
			fmt.Printf("Successfully generated OOC schedule on attempt %d with %d matches.\n", attempt, len(matchesToUpload))
			break
		}
		if attempt < maxAttempts {
			fmt.Printf("Attempt %d failed: %v. Retrying...\n", attempt, err)
		} else {
			fmt.Printf("Failed to generate complete schedule after %d attempts. Last error: %v\n", maxAttempts, err)
			fmt.Printf("Uploading partial schedule with %d matches.\n", len(matchesToUpload))
		}
	}

	repository.CreateCHLGamesRecordsBatch(db, matchesToUpload, 200)
}

func attemptGenerateOOCSchedule(ts structs.Timestamp, seasonID string, collegeTeams []structs.CollegeTeam) ([]structs.CollegeGame, error) {
	maxNumberOfGamesPerTeam := 30
	oocWeeks := []uint{17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1} // Schedule backwards: fill constrained weeks (7-15) first
	timeSlots := []string{"A", "B"}
	sID := ts.SeasonID
	baseID := repository.FindLatestCHLGameID() + 1
	// Initialize tracking maps
	collegeTeamMap := MakeCollegeTeamMap(collegeTeams)
	remainingGamesLeftByTeam := make(map[uint]int)
	timeSlotMap := make(map[uint]map[uint]map[string]bool)
	homeGamesByTeam := make(map[uint]int)
	awayGamesByTeam := make(map[uint]int)
	opponentsFacedMap := make(map[uint]map[uint]bool) // teamID -> opponentTeamID -> bool
	// Count the number of games per week and time slot to ensure we are filling them correctly
	weekTimeSlotCounts := make(map[uint]map[string]int)
	numberOfGamesExpectedPerSlot := len(collegeTeams) / 2 // Each game has 2 teams
	gamesScheduledPerTeam := make(map[uint]int)
	// Process existing matches to update counters
	allActiveMatches := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID})
	arenas := repository.FindAllArenas(repository.ArenaQuery{})
	arenaMap := MakeArenaMap(arenas)
	gamesToUpload := []structs.CollegeGame{}
	mk := func(id uint, invId uint8, title string, week int, isNeutralSite bool, arena structs.Arena) structs.CollegeGame {
		return structs.CollegeGame{
			Model:          gorm.Model{ID: id},
			IsInvitational: true,
			InvitationalID: invId,
			LosingGameID:   0,
			LosingGameHOA:  "",
			BaseGame: structs.BaseGame{
				GameTitle:     title,
				SeasonID:      sID,
				WeekID:        util.GetWeekID(sID, uint(week)),
				Week:          week,
				IsNeutralSite: isNeutralSite,
				LeagueID:      1,
				ArenaID:       arena.ID,
				Arena:         arena.Name,
				City:          arena.City,
				State:         arena.State,
				Country:       arena.Country,
				NextGameID:    0,
				NextGameHOA:   "",
			},
		}
	}

	// Create a map to track user vs AI teams
	isUserTeam := make(map[uint]bool)
	for _, team := range collegeTeams {
		remainingGamesLeftByTeam[team.ID] = maxNumberOfGamesPerTeam
		timeSlotMap[team.ID] = make(map[uint]map[string]bool)
		homeGamesByTeam[team.ID] = 0
		awayGamesByTeam[team.ID] = 0
		opponentsFacedMap[team.ID] = make(map[uint]bool)
		// Assume teams with a coach are user teams
		isUserTeam[team.ID] = team.IsUserCoached
	}

	// Check all active games to ensure no team is scheduled twice in the same timeslot
	// teamSlotSeen: teamID -> "week:timeSlot" -> matchID
	teamSlotSeen := make(map[uint]map[string]uint)
	for _, match := range allActiveMatches {
		if match.AwayTeamID == 0 && match.HomeTeamID == 0 {
			continue
		}
		if weekTimeSlotCounts[uint(match.Week)] == nil {
			weekTimeSlotCounts[uint(match.Week)] = make(map[string]int)
		}
		weekTimeSlotCounts[uint(match.Week)][match.GameDay]++
		slotKey := fmt.Sprintf("%d:%s", match.Week, match.GameDay)
		for _, teamID := range []uint{match.HomeTeamID, match.AwayTeamID} {
			if teamSlotSeen[teamID] == nil {
				teamSlotSeen[teamID] = make(map[string]uint)
			}
			if firstMatchID, exists := teamSlotSeen[teamID][slotKey]; exists {
				fmt.Printf("DOUBLE-SCHEDULED: TeamID %d appears in Week %d TimeSlot %s in both MatchID %d and MatchID %d\n",
					teamID, match.Week, match.GameDay, firstMatchID, match.ID)
			} else {
				teamSlotSeen[teamID][slotKey] = match.ID
			}
		}
	}

	for _, match := range allActiveMatches {
		if match.AwayTeamID == 0 && match.HomeTeamID == 0 {
			continue
		}
		// Initialize nested maps if needed, guarding against team IDs not present in collegeTeams
		if timeSlotMap[match.HomeTeamID] == nil {
			timeSlotMap[match.HomeTeamID] = make(map[uint]map[string]bool)
		}
		if timeSlotMap[match.HomeTeamID][uint(match.Week)] == nil {
			timeSlotMap[match.HomeTeamID][uint(match.Week)] = make(map[string]bool)
		}
		if timeSlotMap[match.AwayTeamID] == nil {
			timeSlotMap[match.AwayTeamID] = make(map[uint]map[string]bool)
		}
		if timeSlotMap[match.AwayTeamID][uint(match.Week)] == nil {
			timeSlotMap[match.AwayTeamID][uint(match.Week)] = make(map[string]bool)
		}

		gamesScheduledPerTeam[match.HomeTeamID]++
		gamesScheduledPerTeam[match.AwayTeamID]++

		// Mark time slots as occupied
		timeSlotMap[match.HomeTeamID][uint(match.Week)][match.GameDay] = true
		timeSlotMap[match.AwayTeamID][uint(match.Week)][match.GameDay] = true
		if opponentsFacedMap[match.HomeTeamID] == nil {
			opponentsFacedMap[match.HomeTeamID] = make(map[uint]bool)
		}
		if opponentsFacedMap[match.AwayTeamID] == nil {
			opponentsFacedMap[match.AwayTeamID] = make(map[uint]bool)
		}
		opponentsFacedMap[match.HomeTeamID][match.AwayTeamID] = true
		opponentsFacedMap[match.AwayTeamID][match.HomeTeamID] = true

		// Update remaining games and home/away counts
		remainingGamesLeftByTeam[match.HomeTeamID]--
		remainingGamesLeftByTeam[match.AwayTeamID]--
		homeGamesByTeam[match.HomeTeamID]++
		awayGamesByTeam[match.AwayTeamID]++
	}

	/* Handle Invitational Requests */
	invitationalRequests := repository.FindInvitationalRequestRecords(repository.SchedulerQuery{SeasonID: seasonID})
	hockeyInvitationals := repository.FindInvitationalRecords(repository.SchedulerQuery{})
	invitationalMap := make(map[uint][]structs.HockeyInvitationalRequest)
	for _, request := range invitationalRequests {
		invitationalMap[request.InvitationalID] = append(invitationalMap[request.InvitationalID], request)
	}

	for _, invitational := range hockeyInvitationals {
		// Bean Pot & Lake Placid Invitationals are handled separately in the offseason, just skip over
		if invitational.ID == 1 || invitational.ID == 3 {
			continue
		}
		requests := invitationalMap[invitational.ID]
		// No Need to Add More Teams to Invitationals
		if len(requests) == int(invitational.TotalTeams) {
			continue
		}

		participatingTeams := []structs.CollegeTeam{}
		for _, request := range requests {
			participatingTeams = append(participatingTeams, collegeTeamMap[request.TeamID])
		}

		ShuffleTeams(participatingTeams)
		participatingArena := arenaMap[uint(invitational.ArenaID)]

		// Game 1
		game1 := mk(baseID, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 1", ts.Season, invitational.Name), int(invitational.Week), false, participatingArena)
		game1.HomeTeam = participatingTeams[0].Abbreviation
		game1.HomeTeamID = participatingTeams[0].ID
		game1.HomeTeamCoach = participatingTeams[0].Coach
		game1.AwayTeam = participatingTeams[1].Abbreviation
		game1.AwayTeamID = participatingTeams[1].ID
		game1.AwayTeamCoach = participatingTeams[1].Coach
		game1.GameDay = "A"
		// Game 2
		game2 := mk(baseID+1, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 2", ts.Season, invitational.Name), int(invitational.Week), true, participatingArena)
		game2.HomeTeam = participatingTeams[2].Abbreviation
		game2.HomeTeamID = participatingTeams[2].ID
		game2.HomeTeamCoach = participatingTeams[2].Coach
		game2.AwayTeam = participatingTeams[3].Abbreviation
		game2.AwayTeamID = participatingTeams[3].ID
		game2.AwayTeamCoach = participatingTeams[3].Coach
		game2.GameDay = "A"

		// Game 3
		game3 := mk(baseID+2, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 3", ts.Season, invitational.Name), int(invitational.Week), true, participatingArena)
		game3.HomeTeam = participatingTeams[3].Abbreviation
		game3.HomeTeamID = participatingTeams[3].ID
		game3.HomeTeamCoach = participatingTeams[3].Coach
		game3.AwayTeam = participatingTeams[0].Abbreviation
		game3.AwayTeamID = participatingTeams[0].ID
		game3.AwayTeamCoach = participatingTeams[0].Coach
		game3.GameDay = "B"
		// Game 4
		game4 := mk(baseID+3, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 4", ts.Season, invitational.Name), int(invitational.Week), true, participatingArena)
		game4.HomeTeam = participatingTeams[1].Abbreviation
		game4.HomeTeamID = participatingTeams[1].ID
		game4.HomeTeamCoach = participatingTeams[1].Coach
		game4.AwayTeam = participatingTeams[2].Abbreviation
		game4.AwayTeamID = participatingTeams[2].ID
		game4.AwayTeamCoach = participatingTeams[2].Coach
		game4.GameDay = "B"
		// Game 5
		game5 := mk(baseID+4, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 5", ts.Season, invitational.Name), int(invitational.Week), true, participatingArena)
		game5.HomeTeam = participatingTeams[0].Abbreviation
		game5.HomeTeamID = participatingTeams[0].ID
		game5.HomeTeamCoach = participatingTeams[0].Coach
		game5.AwayTeam = participatingTeams[2].Abbreviation
		game5.AwayTeamID = participatingTeams[2].ID
		game5.AwayTeamCoach = participatingTeams[2].Coach
		game5.GameDay = "B"
		// Game 6
		game6 := mk(baseID+5, uint8(invitational.ID), fmt.Sprintf("%d %s - Game 6", ts.Season, invitational.Name), int(invitational.Week), true, participatingArena)
		game6.HomeTeam = participatingTeams[1].Abbreviation
		game6.HomeTeamID = participatingTeams[1].ID
		game6.HomeTeamCoach = participatingTeams[1].Coach
		game6.AwayTeam = participatingTeams[3].Abbreviation
		game6.AwayTeamID = participatingTeams[3].ID
		game6.AwayTeamCoach = participatingTeams[3].Coach
		game6.GameDay = "B"

		gamesToUpload = append(gamesToUpload, game1, game2, game3, game4, game5, game6)
	}

	// Iterate through each week and time slot
	for _, week := range oocWeeks {
		for _, timeSlot := range timeSlots {
			// Build list of available teams for this week/timeslot
			var availableTeams []structs.CollegeTeam
			for _, team := range collegeTeams {
				// Check if team needs games and has the time slot available
				if remainingGamesLeftByTeam[team.ID] > 0 {
					if timeSlotMap[team.ID][week] == nil || !timeSlotMap[team.ID][week][timeSlot] {
						availableTeams = append(availableTeams, team)
					}
				}
			}

			// Shuffle available teams for randomness
			availableTeams = ShuffleTeams(availableTeams)

			// Track which teams have been matched in this slot
			matchedInSlot := make(map[uint]bool)

			// Try to pair teams
			for i := 0; i < len(availableTeams); i++ {
				team := availableTeams[i]

				// Skip if already matched in this slot
				if matchedInSlot[team.ID] {
					continue
				}

				// Find an opponent
				opponentIndex := -1
				for j := i + 1; j < len(availableTeams); j++ {
					opponent := availableTeams[j]

					// Check if valid opponent: different conference, not yet matched, has games remaining, and haven't faced each other
					if !matchedInSlot[opponent.ID] &&
						opponent.ConferenceID != team.ConferenceID &&
						remainingGamesLeftByTeam[opponent.ID] > 0 &&
						!opponentsFacedMap[team.ID][opponent.ID] {
						opponentIndex = j
						break
					}
				}

				// If no opponent found, this attempt fails
				if opponentIndex == -1 {
					continue // Try next team
				}

				opponent := availableTeams[opponentIndex]

				// Decide which team is home based on home/away balance
				var homeTeam, awayTeam structs.CollegeTeam

				// Prefer to balance user teams first
				if isUserTeam[team.ID] && !isUserTeam[opponent.ID] {
					// User team: prioritize balance
					if homeGamesByTeam[team.ID] <= awayGamesByTeam[team.ID] {
						homeTeam, awayTeam = team, opponent
					} else {
						homeTeam, awayTeam = opponent, team
					}
				} else if !isUserTeam[team.ID] && isUserTeam[opponent.ID] {
					// Opponent is user team: prioritize their balance
					if homeGamesByTeam[opponent.ID] <= awayGamesByTeam[opponent.ID] {
						homeTeam, awayTeam = opponent, team
					} else {
						homeTeam, awayTeam = team, opponent
					}
				} else {
					// Both user or both AI: balance based on current counts
					teamHomeDeficit := homeGamesByTeam[team.ID] - awayGamesByTeam[team.ID]
					oppHomeDeficit := homeGamesByTeam[opponent.ID] - awayGamesByTeam[opponent.ID]

					if teamHomeDeficit < oppHomeDeficit {
						homeTeam, awayTeam = team, opponent
					} else if teamHomeDeficit > oppHomeDeficit {
						homeTeam, awayTeam = opponent, team
					} else {
						// Equal deficit, randomize
						if rand.Intn(2) == 0 {
							homeTeam, awayTeam = team, opponent
						} else {
							homeTeam, awayTeam = opponent, team
						}
					}
				}

				// Create match
				match := structs.CollegeGame{
					BaseGame: structs.BaseGame{
						SeasonID:      ts.SeasonID,
						Week:          int(week),
						WeekID:        ts.WeekID + week,
						GameDay:       timeSlot,
						HomeTeamID:    homeTeam.ID,
						HomeTeam:      homeTeam.Abbreviation,
						HomeTeamCoach: homeTeam.Coach,
						ArenaID:       uint(homeTeam.ArenaID),
						Arena:         homeTeam.Arena,
						City:          homeTeam.City,
						State:         homeTeam.State,
						AwayTeamID:    awayTeam.ID,
						AwayTeam:      awayTeam.Abbreviation,
						AwayTeamCoach: awayTeam.Coach,
						IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID && homeTeam.ConferenceID != 7,
					},
				}

				gamesToUpload = append(gamesToUpload, match)

				// Update tracking data
				matchedInSlot[team.ID] = true
				matchedInSlot[opponent.ID] = true

				if timeSlotMap[homeTeam.ID][week] == nil {
					timeSlotMap[homeTeam.ID][week] = make(map[string]bool)
				}
				if timeSlotMap[awayTeam.ID][week] == nil {
					timeSlotMap[awayTeam.ID][week] = make(map[string]bool)
				}

				timeSlotMap[homeTeam.ID][week][timeSlot] = true
				timeSlotMap[awayTeam.ID][week][timeSlot] = true

				remainingGamesLeftByTeam[homeTeam.ID]--
				remainingGamesLeftByTeam[awayTeam.ID]--

				homeGamesByTeam[homeTeam.ID]++
				awayGamesByTeam[awayTeam.ID]++

				// Mark that these teams have faced each other
				opponentsFacedMap[homeTeam.ID][awayTeam.ID] = true
				opponentsFacedMap[awayTeam.ID][homeTeam.ID] = true
			}

			// Log teams that were available this slot but could not be paired
			for _, team := range availableTeams {
				if !matchedInSlot[team.ID] && remainingGamesLeftByTeam[team.ID] > 0 {
					fmt.Printf("Week %d, Slot %s: Could not schedule TeamID %d (%s)\n", week, timeSlot, team.ID, team.Abbreviation)
				}
			}
		}
	}

	// Cleanup pass: schedule any remaining games with relaxed constraints.
	// The only requirement is that the two teams are from different conferences;
	// the opponentsFacedMap check is intentionally dropped to allow rematches.
	fmt.Println("Starting cleanup scheduling pass for remaining unscheduled games...")
	for _, week := range oocWeeks {
		for _, timeSlot := range timeSlots {
			var cleanupAvailable []structs.CollegeTeam
			for _, team := range collegeTeams {
				if remainingGamesLeftByTeam[team.ID] > 0 {
					if timeSlotMap[team.ID][week] == nil || !timeSlotMap[team.ID][week][timeSlot] {
						cleanupAvailable = append(cleanupAvailable, team)
					}
				}
			}

			if len(cleanupAvailable) < 2 {
				continue
			}

			cleanupAvailable = ShuffleTeams(cleanupAvailable)
			cleanupMatchedInSlot := make(map[uint]bool)

			for i := 0; i < len(cleanupAvailable); i++ {
				team := cleanupAvailable[i]
				if cleanupMatchedInSlot[team.ID] {
					continue
				}

				opponentIndex := -1
				for j := i + 1; j < len(cleanupAvailable); j++ {
					opponent := cleanupAvailable[j]
					// Relaxed: only require different conference and no previous matchup
					if !cleanupMatchedInSlot[opponent.ID] &&
						opponent.ConferenceID != team.ConferenceID &&
						remainingGamesLeftByTeam[opponent.ID] > 0 &&
						!opponentsFacedMap[team.ID][opponent.ID] {
						opponentIndex = j
						break
					}
				}

				if opponentIndex == -1 {
					continue
				}

				opponent := cleanupAvailable[opponentIndex]

				var homeTeam, awayTeam structs.CollegeTeam
				if isUserTeam[team.ID] && !isUserTeam[opponent.ID] {
					if homeGamesByTeam[team.ID] <= awayGamesByTeam[team.ID] {
						homeTeam, awayTeam = team, opponent
					} else {
						homeTeam, awayTeam = opponent, team
					}
				} else if !isUserTeam[team.ID] && isUserTeam[opponent.ID] {
					if homeGamesByTeam[opponent.ID] <= awayGamesByTeam[opponent.ID] {
						homeTeam, awayTeam = opponent, team
					} else {
						homeTeam, awayTeam = team, opponent
					}
				} else {
					teamHomeDeficit := homeGamesByTeam[team.ID] - awayGamesByTeam[team.ID]
					oppHomeDeficit := homeGamesByTeam[opponent.ID] - awayGamesByTeam[opponent.ID]
					if teamHomeDeficit < oppHomeDeficit {
						homeTeam, awayTeam = team, opponent
					} else if teamHomeDeficit > oppHomeDeficit {
						homeTeam, awayTeam = opponent, team
					} else {
						if rand.Intn(2) == 0 {
							homeTeam, awayTeam = team, opponent
						} else {
							homeTeam, awayTeam = opponent, team
						}
					}
				}

				match := structs.CollegeGame{

					BaseGame: structs.BaseGame{
						SeasonID:      ts.SeasonID,
						Week:          int(week),
						WeekID:        ts.WeekID + week,
						GameDay:       timeSlot,
						HomeTeamID:    homeTeam.ID,
						HomeTeam:      homeTeam.Abbreviation,
						HomeTeamCoach: homeTeam.Coach,
						Arena:         homeTeam.Arena,
						City:          homeTeam.City,
						State:         homeTeam.State,
						AwayTeamID:    awayTeam.ID,
						AwayTeam:      awayTeam.Abbreviation,
						AwayTeamCoach: awayTeam.Coach,
						IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID && homeTeam.ConferenceID != 7,
					},
				}

				gamesToUpload = append(gamesToUpload, match)

				cleanupMatchedInSlot[team.ID] = true
				cleanupMatchedInSlot[opponent.ID] = true

				if timeSlotMap[homeTeam.ID][week] == nil {
					timeSlotMap[homeTeam.ID][week] = make(map[string]bool)
				}
				if timeSlotMap[awayTeam.ID][week] == nil {
					timeSlotMap[awayTeam.ID][week] = make(map[string]bool)
				}

				timeSlotMap[homeTeam.ID][week][timeSlot] = true
				timeSlotMap[awayTeam.ID][week][timeSlot] = true

				remainingGamesLeftByTeam[homeTeam.ID]--
				remainingGamesLeftByTeam[awayTeam.ID]--

				homeGamesByTeam[homeTeam.ID]++
				awayGamesByTeam[awayTeam.ID]++

				opponentsFacedMap[homeTeam.ID][awayTeam.ID] = true
				opponentsFacedMap[awayTeam.ID][homeTeam.ID] = true

				fmt.Printf("Cleanup: Scheduled Week %d, Slot %s: TeamID %d (%s) vs TeamID %d (%s)\n",
					week, timeSlot, homeTeam.ID, homeTeam.Abbreviation, awayTeam.ID, awayTeam.Abbreviation)
			}
		}
	}

	// Validate the schedule - check if all teams that need OOC games got enough
	expectedGamesPerTeam := len(oocWeeks) * len(timeSlots) // 4 weeks * 2 slots = 8 games expected per team

	for _, match := range gamesToUpload {
		gamesScheduledPerTeam[match.HomeTeamID]++
		gamesScheduledPerTeam[match.AwayTeamID]++
	}

	// Count teams with insufficient games
	insufficientTeams := 0
	for _, team := range collegeTeams {
		if gamesScheduledPerTeam[team.ID] < expectedGamesPerTeam {
			insufficientTeams++
		}
	}

	for _, match := range gamesToUpload {
		if weekTimeSlotCounts[uint(match.Week)] == nil {
			weekTimeSlotCounts[uint(match.Week)] = make(map[string]int)
		}
		weekTimeSlotCounts[uint(match.Week)][match.GameDay]++
	}

	// Print out the distribution of games per week and time slot
	for week, slotCounts := range weekTimeSlotCounts {
		for slot, count := range slotCounts {
			fmt.Printf("Week %d, Time Slot %s: %d games scheduled (expected ~%d)\n", week, slot, count, numberOfGamesExpectedPerSlot)
		}
	}

	// If too many teams don't have enough games, consider this attempt failed
	// No, we need a perfect schedule generated.
	tolerance := 0
	if insufficientTeams > tolerance {
		return gamesToUpload, fmt.Errorf("incomplete schedule: %d teams have insufficient games", insufficientTeams)
	}

	return gamesToUpload, nil
}

// ShuffleTeams randomizes the order of teams in a slice
func ShuffleTeams(teams []structs.CollegeTeam) []structs.CollegeTeam {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]structs.CollegeTeam, len(teams))
	copy(shuffled, teams)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}
