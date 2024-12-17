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
	engine.RunGames(gameDTOs)
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
			hp := getCollegePlaybookDTO(htl, htr)
			ap := getCollegePlaybookDTO(atl, atr)
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
		generateCollegeGame(22, 38, "Ferris State", "Minnesota-Duluth"),
		generateCollegeGame(17, 2, "Colorado College", "Fairbanks"),
		generateCollegeGame(57, 11, "St. Lawrence", "Boston University"),
		generateCollegeGame(4, 50, "American International", "Providence"),
		generateCollegeGame(41, 1, "Niagara", "Air Force"),
		generateCollegeGame(18, 28, "UConn", "Maine"),
		generateCollegeGame(15, 26, "Clarkson", "Lindenwood"),
		generateCollegeGame(30, 31, "UMass-Lowell", "Marcyhurst"),
		generateCollegeGame(35, 42, "Michigan State", "North Dakota"),
		generateCollegeGame(47, 45, "Omaha", "Notre Dame"),
		generateCollegeGame(48, 62, "Penn State", "WMU"),
		generateCollegeGame(13, 56, "Brown", "St. Cloud State"),
		generateCollegeGame(23, 20, "Harvard", "Dartmouth"),
		generateCollegeGame(34, 32, "Michigan", "Merrimack"),
		generateCollegeGame(63, 53, "Wisconsin", "Robert Morris"),
		generateCollegeGame(55, 64, "Sacred Heart", "Yale"),
		generateCollegeGame(9, 66, "Bentley", "Binghamton"),
		generateCollegeGame(40, 29, "New Hampshire", "UMass"),
		generateCollegeGame(49, 54, "Princeton", "Rochester"),
		generateCollegeGame(58, 65, "St. Thomas", "Tennessee State"),
		generateCollegeGame(6, 3, "Army", "Anchorage"),
		generateCollegeGame(7, 51, "Augustana", "Quinnipiac"),
		generateCollegeGame(37, 5, "Minnesota", "Arizona State"),
		generateCollegeGame(16, 33, "Colgate", "Miami (OH)"),
		generateCollegeGame(61, 12, "Vermont", "Bowling Green State"),
		generateCollegeGame(7, 59, "Bemidji State", "Stonehill"),
		generateCollegeGame(44, 10, "Northern Michigan", "Boston College"),
		generateCollegeGame(21, 52, "Denver", "RPI"),
		generateCollegeGame(46, 24, "Ohio State", "Holy Cross"),
		generateCollegeGame(39, 25, "Minnesota State", "Lake Superior State"),
		generateCollegeGame(43, 19, "Northeastern", "Cornell"),
		generateCollegeGame(60, 36, "Union", "Michigan Tech"),
		generateCollegeGame(14, 27, "Canisius", "Long Island"),
	}
	return games
}

func GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay string) []structs.CollegeGame {
	return repository.FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay)
}

func GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay string) []structs.ProfessionalGame {
	return repository.FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay)
}

func GetArenaMap() map[uint]structs.Arena {
	arenas := repository.FindAllArenas()
	return MakeArenaMap(arenas)
}

func getCollegePlaybookDTO(lineups []structs.CollegeLineup, roster []structs.CollegePlayer) structs.PlayBookDTO {
	forwards, defenders, goalies := getCollegeForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:      forwards,
		Defenders:     defenders,
		Goalies:       goalies,
		CollegeRoster: roster,
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
