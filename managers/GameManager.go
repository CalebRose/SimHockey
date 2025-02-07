package managers

import (
	"fmt"
	"math/rand"
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
	collegeTeamMap := GetCollegeTeamMap()
	proTeamMap := GetProTeamMap()
	collegePlayerMap := GetCollegePlayersMap()
	proPlayersMap := GetProPlayersMap()
	for _, r := range results {
		// Iterate through all lines, players, accumulate stats to upload
		WriteBoxScoreFile(r, "test_results/test_twelve/box_score/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv")

		// Iterate through Play By Plays and record them to a CSV
		pbps := r.Collector.PlayByPlays
		if r.IsCollegeGame {
			WritePlayByPlayCSVFile(pbps, "test_results/test_twelve/play_by_play/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv", collegePlayerMap, collegeTeamMap)
		} else {
			WriteProPlayByPlayCSVFile(pbps, "test_results/test_twelve/play_by_play/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv", proPlayersMap, proTeamMap)
		}
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
	collegeTeamMap := GetCollegeTeamMap()
	collegeTeamRosterMap := GetAllCollegePlayersMapByTeam()
	collegeLineupMap := GetCollegeLineupsMap()
	collegeShootoutLineupMap := GetCollegeShootoutLineups()
	// weekID := strconv.Itoa(int(ts.WeekID))
	// seasonID := strconv.Itoa(int(ts.SeasonID))
	// gameDay := ts.GetGameDay()
	arenaMap := GetArenaMap()
	// collegeGames := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay)
	collegeGames := GetCollegeGamesForTesting(collegeTeamMap)
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
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	var proGamesWg sync.WaitGroup
	professionalTeamMap := GetProTeamMap()
	proGames := GetProGamesForTesting(professionalTeamMap)
	proTeamRosterMap := GetAllProPlayersMapByTeam()
	proLineupMap := GetProLineupsMap()
	proShootoutLineupMap := GetProShootoutLineups()
	proGamesWg.Add(len(proGames))
	proSem := make(chan struct{}, 20)
	for _, g := range proGames {
		proSem <- struct{}{}
		localC := g
		go func(g structs.ProfessionalGame) {
			defer func() { <-proSem }()
			defer proGamesWg.Done()
			if g.GameComplete {
				return
			}
			mutex.Lock()
			htr := proTeamRosterMap[g.HomeTeamID]
			atr := proTeamRosterMap[g.AwayTeamID]
			htl := proLineupMap[g.HomeTeamID]
			atl := proLineupMap[g.AwayTeamID]
			htsl := proShootoutLineupMap[g.HomeTeamID]
			atsl := proShootoutLineupMap[g.AwayTeamID]
			hp := getProfessionalPlaybookDTO(htl, htr, htsl)
			ap := getProfessionalPlaybookDTO(atl, atr, atsl)
			capacity := 0

			arena := arenaMap[g.ArenaID]
			if arena.ID == 0 {
				capacity = 6000
			} else {
				capacity = int(arena.Capacity)
			}
			mutex.Unlock()

			match := structs.GameDTO{
				GameInfo:      g.BaseGame,
				HomeStrategy:  hp,
				AwayStrategy:  ap,
				IsCollegeGame: false,
				Attendance:    uint32(capacity),
			}

			mutex.Lock()
			gameDTOList = append(gameDTOList, match)
			mutex.Unlock()
		}(localC)
	}
	collegeGamesWg.Wait()
	proGamesWg.Wait()
	for i := 0; i < cap(proSem); i++ {
		proSem <- struct{}{}
	}
	return gameDTOList
}

func GetCollegeGamesForTesting(teamMap map[uint]structs.CollegeTeam) []structs.CollegeGame {
	numbers := make([]uint, 66)
	for i := uint(1); i <= 66; i++ {
		numbers[i-1] = i
	}

	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	pairs, _ := createProGamePairings(numbers)
	games := []structs.CollegeGame{}

	for _, pair := range pairs {
		game := generateCollegeGame(pair[0], pair[1], teamMap)
		games = append(games, game)
	}
	return games
}

func GetProGamesForTesting(teamMap map[uint]structs.ProfessionalTeam) []structs.ProfessionalGame {
	numbers := make([]uint, 32)
	for i := uint(1); i <= 32; i++ {
		numbers[i-1] = i
	}

	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	pairs, _ := createProGamePairings(numbers)
	games := []structs.ProfessionalGame{}

	for _, pair := range pairs {
		game := generateProfessionalGame(pair[0], pair[1], teamMap)
		games = append(games, game)
	}
	return games
}

func createProGamePairings(numbers []uint) ([][2]uint, error) {
	// Check if the number of elements is even
	if len(numbers)%2 != 0 {
		return nil, fmt.Errorf("the list must have an even number of elements")
	}

	// Create pairings
	var pairings [][2]uint
	for i := 0; i < len(numbers); i += 2 {
		pair := [2]uint{numbers[i], numbers[i+1]}
		pairings = append(pairings, pair)
	}

	return pairings, nil
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

func getProfessionalPlaybookDTO(lineups []structs.ProfessionalLineup, roster []structs.ProfessionalPlayer, shootoutLineup structs.ProfessionalShootoutLineup) structs.PlayBookDTO {
	forwards, defenders, goalies := getProfessionalForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:           forwards,
		Defenders:          defenders,
		Goalies:            goalies,
		ProfessionalRoster: roster,
		ShootoutLineup:     shootoutLineup.ShootoutPlayerIDs,
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

func getProfessionalForwardDefenderGoalieLineups(lineups []structs.ProfessionalLineup) ([]structs.BaseLineup, []structs.BaseLineup, []structs.BaseLineup) {
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

func generateCollegeGame(hid, aid uint, teamMap map[uint]structs.CollegeTeam) structs.CollegeGame {
	return structs.CollegeGame{
		BaseGame: structs.BaseGame{
			WeekID:     1,
			Week:       1,
			GameDay:    "A",
			SeasonID:   1,
			HomeTeamID: hid,
			HomeTeam:   teamMap[hid].TeamName,
			AwayTeamID: aid,
			AwayTeam:   teamMap[aid].TeamName,
			ArenaID:    uint(teamMap[hid].ArenaID),
		},
	}
}

func generateProfessionalGame(hid, aid uint, teamMap map[uint]structs.ProfessionalTeam) structs.ProfessionalGame {
	return structs.ProfessionalGame{
		BaseGame: structs.BaseGame{
			WeekID:     1,
			Week:       1,
			GameDay:    "A",
			SeasonID:   1,
			HomeTeamID: hid,
			HomeTeam:   teamMap[hid].Abbreviation,
			AwayTeamID: aid,
			AwayTeam:   teamMap[aid].Abbreviation,
			ArenaID:    uint(teamMap[hid].ArenaID),
		},
	}
}
