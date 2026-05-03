package managers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/structs"
)

// StartLiveScoreboardSession spins up active games and routes JSON payloads to the frontend
func StartLiveScoreboardSession(ctx context.Context, leagueType string, gameLimit int, outChannel chan<- string) {
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()

	var activeGames []structs.GameDTO

	if leagueType == "CHL" {
		games := GetCollegeGamesBySeasonID("", false)
		collegeStandingsMap := GetCollegeStandingsMap(seasonID)
		activeGames = PrepareGames(games, nil, collegeStandingsMap, nil)
	} else {
		games := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)
		proStandingsMap := GetProStandingsMap(seasonID)
		activeGames = PrepareGames(nil, games, nil, proStandingsMap)
	}

	// Filter down to only games that are NOT complete and haven't been revealed
	filteredGames := []structs.GameDTO{}
	for _, g := range activeGames {
		if true { // Force all games to run for testing
			filteredGames = append(filteredGames, g)
		}
	}

	// Limit games based on config (4 or 8)
	if len(filteredGames) > gameLimit {
		filteredGames = filteredGames[:gameLimit]
	}

	// Channel to receive raw plays from the engine
	engineChannel := make(chan structs.PbP)

	// Spin up a goroutine for each filtered game
	for _, game := range filteredGames {
		go engine.RunLiveHockeyGame(game, engineChannel)
	}

	// Load Maps needed for translating raw IDs to readable text
	collegePlayerMap := GetCollegePlayersMap()
	collegeTeamMap := GetCollegeTeamMap()
	proPlayerMap := GetProPlayersMap()
	proTeamMap := GetProTeamMap()

	for {
		select {
		case <-ctx.Done():
			// The user closed the browser, stop processing
			return
		case play := <-engineChannel:
			// 1. Translate the raw IDs to readable strings
			eventString := GetEventName(play.EventID)
			outcomeString := GetOutcomeName(play.Outcome)

			var text string

			// 2. Generate the readable play-by-play text
			if leagueType == "CHL" {
				possessingTeam := collegeTeamMap[uint(play.TeamID)]
				text = generateCollegeResultsString(play, eventString, outcomeString, collegePlayerMap, possessingTeam)
			} else {
				possessingTeam := proTeamMap[uint(play.TeamID)]
				text = generateProResultsString(play, eventString, outcomeString, proPlayerMap, possessingTeam)
			}

			// 3. Package it into a UI-friendly object
			payloadObj := map[string]interface{}{
				"GameID":            play.GameID,
				"Period":            play.Period,
				"TimeOnClock":       play.TimeOnClock,
				"HomeScore":         play.HomeTeamScore,
				"AwayScore":         play.AwayTeamScore,
				"HomeShootoutScore": play.HomeTeamShootoutScore,
				"AwayShootoutScore": play.AwayTeamShootoutScore,
				"Zone":              play.ZoneID,
				"PlayText":          text,
			}

			payload, _ := json.Marshal(payloadObj)
			outChannel <- string(payload)
		}
	}
}

// GetEventName translates integer Event IDs to string constants
func GetEventName(eventID uint8) string {
	switch eventID {
	case 1:
		return Faceoff
	case 2:
		return PhysDefenseCheck
	case 3:
		return DexDefenseCheck
	case 4:
		return PassCheck
	case 5:
		return AgilityCheck
	case 6:
		return WristshotCheck
	case 7:
		return SlapshotCheck
	case 8:
		return PenaltyCheck
	case 34:
		return EnteringShootout
	case 35, 36:
		return Shootout
	case 37:
		return PuckBattle
	case 40:
		return PuckScramble
	case 41:
		return PuckCovered
	case 42:
		return LongPassCheck
	case 43:
		return PassBackCheck
	case 44:
		return "Injury Check"
	default:
		return ""
	}
}

// GetOutcomeName translates integer Outcome IDs to string constants
func GetOutcomeName(outcomeID uint8) string {
	switch outcomeID {
	case 14:
		return DefenseTakesPuck
	case 15:
		return CarrierKeepsPuck
	case 16:
		return DefenseStopAgility
	case 17:
		return OffenseMovesUp
	case 18:
		return GeneralPenalty
	case 20:
		return FightPenalty
	case 21:
		return InterceptedPass
	case 22:
		return ReceivedPass
	case 23:
		return HomeFaceoffWin
	case 24:
		return AwayFaceoffWin
	case 25:
		return InAccurateShot
	case 26:
		return ShotBlocked
	case 27:
		return GoalieSave
	case 28:
		return GoalieReboundOutcome
	case 29:
		return ShotOnGoal
	case 30:
		return "Goalie Hold"
	case 32:
		return ReceivedLongPass
	case 33:
		return ReceivedBackPass
	case 38:
		return PuckBattleWin
	case 39:
		return PuckBattleLose
	default:
		return ""
	}
}
