package engine

import (
	"fmt"
	"time"

	"github.com/CalebRose/SimHockey/structs"
)

// RunLiveHockeyGame simulates the game in memory and broadcasts events in "real-time"
// CRITICAL: This function must NEVER return the GameState to be saved to the DB.
func RunLiveHockeyGame(game structs.GameDTO, broadcastChannel chan<- structs.PbP) {
	gs := generateGameState(game)
	gs.HomeStrategy.ScrubInjuredPlayersFromLineups()
	gs.AwayStrategy.ScrubInjuredPlayersFromLineups()
	gs.HomeStrategy.InitializeStamina()
	gs.AwayStrategy.InitializeStamina()

	totalPeriods := 3
	for gs.Period <= uint8(totalPeriods) && !gs.GameComplete {
		playLivePeriod(&gs, broadcastChannel)
	}

	if gs.HomeTeamScore == gs.AwayTeamScore {
		playLiveOvertime(&gs, broadcastChannel)
	}
	if gs.HomeTeamScore != gs.AwayTeamScore && gs.IsOvertime {
		gs.CalculateWinner()
	}

	if gs.HomeTeamScore == gs.AwayTeamScore && gs.IsOvertime {
		gs.IsOvertimeShootout = true
		HandleLiveOvertimeShootout(&gs, broadcastChannel)
		gs.CalculateWinner()
	}

	fmt.Printf("Live Game Complete: %s %d - %s %d\n", gs.HomeTeam, gs.HomeTeamScore, gs.AwayTeam, gs.AwayTeamScore)
}

func playLivePeriod(gs *GameState, broadcastChannel chan<- structs.PbP) {
	gs.SetTime(true, false)
	for gs.TimeOnClock > 0 && !gs.GameComplete {
		// Track how many plays we had before the event
		playsCountBefore := len(gs.Collector.PlayByPlays)

		if gs.FaceoffOnCenterIce {
			HandleFaceoff(gs)
		}
		HandleBaseEvents(gs)
		gs.SetTime(false, false)

		// If a new play was generated, broadcast it
		playsCountAfter := len(gs.Collector.PlayByPlays)
		if playsCountAfter > playsCountBefore {
			newPlays := gs.Collector.PlayByPlays[playsCountBefore:playsCountAfter]

			for _, play := range newPlays {
				broadcastChannel <- play
				// Pause to simulate real-time. Adjust this to make the game faster/slower.
				time.Sleep(3 * time.Second)
			}
		}
	}
	gs.SetNewZone(NeutralZone)
	gs.SetFaceoffOnCenterIce(true)
}

func playLiveOvertime(gs *GameState, broadcastChannel chan<- structs.PbP) {
	gs.SetTime(true, true)
	for gs.TimeOnClock > 0 {
		playsCountBefore := len(gs.Collector.PlayByPlays)

		if gs.FaceoffOnCenterIce {
			HandleFaceoff(gs)
		}
		HandleBaseEvents(gs)
		gs.SetTime(false, false)

		playsCountAfter := len(gs.Collector.PlayByPlays)
		if playsCountAfter > playsCountBefore {
			newPlays := gs.Collector.PlayByPlays[playsCountBefore:playsCountAfter]
			for _, play := range newPlays {
				broadcastChannel <- play
				time.Sleep(3 * time.Second)
			}
		}

		if gs.HomeTeamScore != gs.AwayTeamScore {
			gs.CalculateWinner()
			break
		}
	}
}

func HandleLiveOvertimeShootout(gs *GameState, broadcastChannel chan<- structs.PbP) {
	isRepeat := false
	shootoutQueue := formShootoutQueue(gs.HomeStrategy, gs.AwayStrategy)

	// FIXED: Using "EnteringShootout" from your constants.go instead of EnteringShootoutID
	RecordPlay(gs, EnteringShootout, 0, 0, 0, 0, 0, 0, 0, false, 0, 0, 0, 0, 0, false)
	broadcastChannel <- gs.Collector.PlayByPlays[len(gs.Collector.PlayByPlays)-1]
	time.Sleep(2 * time.Second)

	for gs.HomeTeamShootoutScore == gs.AwayTeamShootoutScore {
		for idx, player := range shootoutQueue {
			if (idx > 5 && gs.HomeTeamShootoutScore != gs.AwayTeamShootoutScore && !isRepeat) || (isRepeat && gs.HomeTeamShootoutScore != gs.AwayTeamShootoutScore) {
				break
			}
			playsCountBefore := len(gs.Collector.PlayByPlays)

			gs.SetPuckBearer(player, false)
			HandleShootoutAttempt(gs)

			playsCountAfter := len(gs.Collector.PlayByPlays)
			if playsCountAfter > playsCountBefore {
				newPlays := gs.Collector.PlayByPlays[playsCountBefore:playsCountAfter]
				for _, play := range newPlays {
					broadcastChannel <- play
					time.Sleep(4 * time.Second) // Shootout shots get a bit more dramatic pause
				}
			}
		}
	}
}
