package engine

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimHockey/structs"
)

func RunGames(games []structs.GameDTO) {
	fmt.Println("Running Games...")
	for _, g := range games {
		RunTheGame(g)
	}
}

func RunTheGame(game structs.GameDTO) {
	gs := generateGameState(game)
	gs.HomeStrategy.InitializeStamina()
	gs.AwayStrategy.InitializeStamina()

	// While the game is at 3 periods
	totalPeriods := 3
	for gs.Period <= uint8(totalPeriods) && !gs.GameComplete {
		playPeriod(&gs)
	}

	if gs.HomeTeamScore == gs.AwayTeamScore {
		playOvertime(&gs)
	}

	if gs.HomeTeamScore == gs.AwayTeamScore && gs.IsOvertime {
		gs.IsOvertimeShootout = true
		HandleOvertimeShootout(&gs)
		gs.CalculateWinner()
	}

	logGameResults(gs)
}

// Play a single period
func playPeriod(gs *GameState) {
	gs.SetTime(true, false)
	for gs.TimeOnClock > 0 && !gs.GameComplete {
		if gs.FaceoffOnCenterIce {
			HandleFaceoff(gs)
		}
		HandleBaseEvents(gs)
		gs.SetTime(false, false)
	}
}

// Handle the overtime phase
func playOvertime(gs *GameState) {
	gs.SetTime(true, true)
	for gs.TimeOnClock > 0 {
		if gs.FaceoffOnCenterIce {
			HandleFaceoff(gs)
		}
		HandleBaseEvents(gs)
		gs.SetTime(false, false)

		// If a team scores, the game ends immediately
		if gs.HomeTeamScore != gs.AwayTeamScore {
			gs.CalculateWinner()
			break
		}
	}
}

// Log game results
func logGameResults(gs GameState) {
	// Handle Team Stats
	overTimeString := ""
	if gs.IsOvertime {
		overTimeCount := gs.Period - 3
		overTimeString = " (" + strconv.Itoa(int(overTimeCount)) + "OT)"
		if gs.IsOvertimeShootout {
			overTimeString = " (Shootout) | Home Shootout: " + strconv.Itoa(int(gs.HomeTeamShootoutScore)) + " | Away Team Shootout: " + strconv.Itoa(int(gs.AwayTeamShootoutScore))
		}
	}
	scoreStr := fmt.Sprintf(
		"%s: %d | %s: %d%s",
		gs.HomeTeam, gs.HomeTeamScore, gs.AwayTeam, gs.AwayTeamScore, overTimeString,
	)
	Logger(scoreStr)
}

func generateGameState(game structs.GameDTO) GameState {
	gameInfo := game.GameInfo
	gs := GameState{
		HomeTeamID:         gameInfo.HomeTeamID,
		HomeTeam:           gameInfo.HomeTeam,
		HomeStrategy:       loadGamePlaybook(game.IsCollegeGame, game.HomeStrategy),
		AwayTeamID:         gameInfo.AwayTeamID,
		AwayTeam:           gameInfo.AwayTeam,
		AwayStrategy:       loadGamePlaybook(game.IsCollegeGame, game.AwayStrategy),
		Period:             0,
		Momentum:           0,
		FaceoffOnCenterIce: true,
		PuckLocation:       NeutralZone,
		TimeOnClock:        1200, // Number of seconds per period
		IsPlayoffGame:      gameInfo.IsPlayoffGame,
		IsCollegeGame:      game.IsCollegeGame,
		IsOvertime:         false,
		IsOvertimeShootout: false,
		GameComplete:       false,
		TieGame:            false,
	}

	return gs
}

func loadGamePlaybook(isCollegeGame bool, pb structs.PlayBookDTO) GamePlaybook {
	gameRoster := LoadGameRoster(isCollegeGame, pb.CollegeRoster, pb.ProfessionalRoster)
	forwardLines, defenderLines, goalieLines, activeIDs := LoadAllLineStrategies(pb, gameRoster)
	benchPlayers := LoadBenchPlayers(activeIDs, gameRoster)
	return GamePlaybook{
		Forwards:         forwardLines,
		Defenders:        defenderLines,
		Goalies:          goalieLines,
		CurrentForwards:  0,
		CurrentDefenders: 0,
		CurrentGoalie:    0,
		BenchPlayers:     benchPlayers,
	}
}
