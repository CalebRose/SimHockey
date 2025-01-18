package managers

import (
	"fmt"
	"sync"

	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func RunGames() {
	// Get GameDTOs
	gameDTOs := PrepareGames()
	// RUN THE GAMES!
	results := engine.RunGames(gameDTOs)
	teamMap := GetCollegeTeamMap()
	collegePlayerMap := GetCollegePlayersMap()
	for _, r := range results {
		// Iterate through all lines, players, accumulate stats to upload

		// Iterate through Play By Plays and record them to a CSV
		pbps := r.Collector.PlayByPlays
		WritePlayByPlayCSVFile(pbps, "test_results/test_five/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv", collegePlayerMap, teamMap)
	}
}

func PrepareGames() []structs.GameDTO {
	fmt.Println("Loading Games...")

	// ts := GetTimestamp()
	// Wait Groups
	var collegeGamesWg sync.WaitGroup

	// Mutex Lock
	var mutex sync.Mutex

	// College Only
	collegeTeamRosterMap := GetAllCollegePlayersMapByTeam()
	collegeLineupMap := GetCollegeLineupsMap()
	collegeShootoutLineupMap := GetCollegeShootoutLineups()
	// weekID := strconv.Itoa(int(ts.WeekID))
	// seasonID := strconv.Itoa(int(ts.SeasonID))
	// gameDay := ts.GetGameDay()
	arenaMap := GetArenaMap()
	// collegeGames := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay)
	collegeGames := GetCollegeGamesForTesting()

	collegeGamesWg.Add(len(collegeGames))
	gameDTOList := make([]structs.GameDTO, 0, len(collegeGames))
	sem := make(chan struct{}, 20)

	for _, c := range collegeGames {
		sem <- struct{}{}
		localC := c
		go func(c structs.CollegeGame) {
			defer func() { <-sem }()
			defer collegeGamesWg.Done()
			if c.GameComplete {
				return
			}
			mutex.Lock()
			htr := collegeTeamRosterMap[c.HomeTeamID]
			atr := collegeTeamRosterMap[c.AwayTeamID]
			htl := collegeLineupMap[c.HomeTeamID]
			atl := collegeLineupMap[c.AwayTeamID]
			htsl := collegeShootoutLineupMap[c.HomeTeamID]
			atsl := collegeShootoutLineupMap[c.AwayTeamID]
			hp := getCollegePlaybookDTO(htl, htr, htsl)
			ap := getCollegePlaybookDTO(atl, atr, atsl)
			capacity := 0

			arena := arenaMap[c.ArenaID]
			if arena.ID == 0 {
				capacity = 6000
			} else {
				capacity = int(arena.Capacity)
			}
			mutex.Unlock()

			match := structs.GameDTO{
				GameInfo:      c.BaseGame,
				HomeStrategy:  hp,
				AwayStrategy:  ap,
				IsCollegeGame: true,
				Attendance:    uint32(capacity),
			}

			mutex.Lock()
			gameDTOList = append(gameDTOList, match)
			mutex.Unlock()
		}(localC)
	}
	collegeGamesWg.Wait()
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	return gameDTOList
}

func GetCollegeGamesForTesting() []structs.CollegeGame {
	games := []structs.CollegeGame{
		// generateCollegeGame(63, 38, "Wisconsin", "Minnesota-Duluth"),
		// generateCollegeGame(22, 45, "Ferris State", "Notre Dame"),
		// generateCollegeGame(38, 4, "Minnesota-Duluth", "American International"),
		// generateCollegeGame(48, 14, "Penn State", "Canisius"),
		// generateCollegeGame(1, 42, "Air Force", "North Dakota"),
		// generateCollegeGame(28, 15, "Maine", "Clarkson"),
		// generateCollegeGame(11, 53, "Boston U", "Robert Morris"),
		// generateCollegeGame(56, 17, "St. Cloud State", "Colorado College"),
		// generateCollegeGame(32, 30, "Merrimack", "UMass-Lowell"),
		// generateCollegeGame(20, 60, "Dartmouth", "Union"),
		// generateCollegeGame(27, 55, "Long Island", "Sacred Heart"),
		// generateCollegeGame(50, 31, "Providence", "Mercyhurst"),
		// generateCollegeGame(29, 41, "UMass", "Niagara"),
		// generateCollegeGame(49, 35, "Princeton", "Michigan State"),
		// generateCollegeGame(6, 66, "Army", "Binghamton"),
		// generateCollegeGame(7, 64, "Augustana", "Yale"),
		// generateCollegeGame(13, 62, "Brown", "WMU"),
		// generateCollegeGame(58, 16, "St. Thomas", "Colgate"),
		// generateCollegeGame(23, 47, "Harvard", "Omaha"),
		generateCollegeGame(18, 37, "UConn", "Minnesota"),
		// generateCollegeGame(26, 9, "Lindenwood", "Bentley"),
		// generateCollegeGame(8, 63, "Bemidji State", "Wisconsin"),
		generateCollegeGame(10, 61, "Boston College", "Vermont"),
		// generateCollegeGame(65, 34, "Tennessee State", "Michigan"),
		// generateCollegeGame(54, 21, "Rochester", "Denver"),
		// generateCollegeGame(36, 52, "Michigan Tech", "RPI"),
		// generateCollegeGame(40, 46, "New Hampshire", "Ohio State"),
		// generateCollegeGame(39, 43, "Minnesota State", "Northeastern"),
		// generateCollegeGame(3, 51, "Anchorage", "Quinnipiac"),
		// generateCollegeGame(33, 12, "Miami (OH)", "Bowling Green"),
		// generateCollegeGame(5, 44, "Arizona State", "Northern Michigan"),
		// generateCollegeGame(25, 59, "Lake Superior State", "Stonehill"),
		// generateCollegeGame(24, 19, "Holy Cross", "Cornell"),
	}
	return games
}

func GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay string) []structs.CollegeGame {
	return repository.FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay)
}

func GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay string) []structs.ProfessionalGame {
	return repository.FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay)
}

func GetCollegeGamesByTeamIDAndSeasonID(teamID, seasonID string) []structs.CollegeGame {
	return repository.FindCollegeGames(seasonID, teamID)
}

func GetProfessionalGamesByTeamIDAndSeasonID(teamID, seasonID string) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(seasonID, teamID)
}

func GetCollegeGamesBySeasonID(seasonID string) []structs.CollegeGame {
	return repository.FindCollegeGames(seasonID, "")
}

func GetProfessionalGamesBySeasonID(seasonID string) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(seasonID, "")
}

func GetArenaMap() map[uint]structs.Arena {
	arenas := repository.FindAllArenas()
	return MakeArenaMap(arenas)
}

func getCollegePlaybookDTO(lineups []structs.CollegeLineup, roster []structs.CollegePlayer, shootoutLineup structs.CollegeShootoutLineup) structs.PlayBookDTO {
	forwards, defenders, goalies := getCollegeForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:       forwards,
		Defenders:      defenders,
		Goalies:        goalies,
		CollegeRoster:  roster,
		ShootoutLineup: shootoutLineup.ShootoutPlayerIDs,
	}
}

// func getBaseRoster(roster []structs.CollegePlayer) []structs.BasePlayer {
// 	basePlayers := []structs.BasePlayer{}

// 	for _, p := range roster {
// 		basePlayers = append(basePlayers, p.BasePlayer)
// 	}
// 	return basePlayers
// }

func getCollegeForwardDefenderGoalieLineups(lineups []structs.CollegeLineup) ([]structs.BaseLineup, []structs.BaseLineup, []structs.BaseLineup) {
	forwards := []structs.BaseLineup{}
	defenders := []structs.BaseLineup{}
	goalies := []structs.BaseLineup{}
	for _, l := range lineups {
		if l.LineType == 1 {
			forwards = append(forwards, l.BaseLineup)
		} else if l.LineType == 2 {
			defenders = append(defenders, l.BaseLineup)
		} else {
			goalies = append(goalies, l.BaseLineup)
		}
	}
	return forwards, defenders, goalies
}

func generateCollegeGame(hid, aid uint, ht, at string) structs.CollegeGame {
	return structs.CollegeGame{
		BaseGame: structs.BaseGame{
			WeekID:     1,
			Week:       1,
			GameDay:    "A",
			SeasonID:   1,
			HomeTeamID: hid,
			HomeTeam:   ht,
			AwayTeamID: aid,
			AwayTeam:   at,
			ArenaID:    1,
		},
	}
}
