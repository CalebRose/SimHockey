package engine

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimHockey/structs"
)

func RunGames(games []structs.GameDTO) []GameState {
	results := []GameState{}
	fmt.Println("Running Games...")
	for _, g := range games {
		fmt.Println("Starting Game... " + g.GameInfo.HomeTeam + " vs. " + g.GameInfo.AwayTeam)
		res := RunTheGame(g)
		results = append(results, res)
	}
	return results
}

func RunTheGame(game structs.GameDTO) GameState {
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
	if gs.HomeTeamScore != gs.AwayTeamScore && gs.IsOvertime {
		gs.CalculateWinner()
	}

	if gs.HomeTeamScore == gs.AwayTeamScore && gs.IsOvertime {
		gs.IsOvertimeShootout = true
		HandleOvertimeShootout(&gs)
		gs.CalculateWinner()
	}

	logGameResults(gs)
	return gs
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
	gs.SetNewZone(NeutralZone)
	gs.SetFaceoffOnCenterIce(true)
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
	hra := float64(game.Attendance) / float64(game.Capacity)
	// If no arena generated, no home rink advantage. Granted, this edge case should never happen.
	if game.Capacity == 0 {
		hra = 1.0
	}
	gs := GameState{
		GameID:            game.GameID,
		WeekID:            game.GameInfo.WeekID,
		Attendance:        game.Attendance,
		HomeRinkAdvantage: hra,
		HomeTeamID:        gameInfo.HomeTeamID,
		HomeTeam:          gameInfo.HomeTeam,
		HomeStrategy:      loadGamePlaybook(game.IsCollegeGame, true, game.HomeStrategy, game.GameInfo.SeasonID, game.GameInfo.GameDay, hra),
		HomeTeamStats: TeamStatDTO{
			TeamID:  game.GameInfo.HomeTeamID,
			Team:    game.GameInfo.HomeTeam,
			GameDay: game.GameInfo.GameDay,
		},
		AwayTeamID:   gameInfo.AwayTeamID,
		AwayTeam:     gameInfo.AwayTeam,
		AwayStrategy: loadGamePlaybook(game.IsCollegeGame, false, game.AwayStrategy, game.GameInfo.SeasonID, game.GameInfo.GameDay, hra),
		AwayTeamStats: TeamStatDTO{
			TeamID:  game.GameInfo.AwayTeamID,
			Team:    game.GameInfo.AwayTeam,
			GameDay: game.GameInfo.GameDay,
		},
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
		PuckCarrier:        &GamePlayer{},
		AssistingPlayer:    &GamePlayer{},
	}
	gs.EnableStartedGameStat()

	return gs
}

func loadGamePlaybook(isCollegeGame, isHome bool, pb structs.PlayBookDTO, seasonID uint, gameDay string, hra float64) GamePlaybook {
	gameRoster := LoadGameRoster(isCollegeGame, pb.CollegeRoster, pb.ProfessionalRoster, seasonID, gameDay, isHome, hra)
	rosterMap := getGameRosterMap(gameRoster)
	forwardLines, defenderLines, goalieLines, activeIDs := LoadAllLineStrategies(pb, gameRoster)
	benchPlayers := LoadBenchPlayers(activeIDs, gameRoster)
	return GamePlaybook{
		Forwards:           forwardLines,
		Defenders:          defenderLines,
		Goalies:            goalieLines,
		CurrentForwards:    0,
		CurrentDefenders:   0,
		CurrentGoalie:      0,
		BenchPlayers:       benchPlayers,
		ShootoutLineUp:     pb.ShootoutLineup,
		RosterMap:          rosterMap,
		ForwardShiftTimer:  0,
		DefenderShiftTimer: 0,
		ForwardShiftLimit:  ForwardShiftLimit,
		DefenderShiftLimit: DefenderShiftLimit,
	}
}
