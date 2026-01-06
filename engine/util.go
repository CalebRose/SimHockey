package engine

import (
	"math"
	"sort"

	util "github.com/CalebRose/SimHockey/_util"
)

/*
BaseCheck : The base number of which the dice roll must be met.
MomentumMod : The number of passes & agility checks passed. Rebounds after shot. The lower the number, the better the chance of a shot
GoalieMod : Depending on if it is a slapshot or wrist shot, goalie needs to utilize either agility or strength
*/
func CalculateShot(shotModifier, momentumModifier, goaltendingModifier, goalieModifier, baseCheck float64) bool {
	max := 20.0
	min := 1.0

	// Calculate the effective shot score with a secondary logarithmic scaling
	combinedModifiers := shotModifier + momentumModifier
	adjustedShotScore := calculateModifier(combinedModifiers, ScaleFactor)
	combinedGoalieModifiers := goaltendingModifier + goalieModifier
	adjustedGoalieScore := calculateModifier(combinedGoalieModifiers, ScaleFactor)

	num := util.GenerateFloatFromRange(min, max) + adjustedShotScore
	return num > (baseCheck + adjustedGoalieScore)
}

func CalculateAccuracy(accuracyModifier float64, isCloseShot bool) bool {
	max := 20.0
	min := 1.0
	req := SlightlyDiffReq
	if isCloseShot {
		req = BaseReq
	}
	adjAccuracy := calculateModifier(accuracyModifier, ScaleFactor)
	accuracyCheck := util.GenerateFloatFromRange(min, max)
	return accuracyCheck < req+adjAccuracy
}

func CalculateShotBlock(shotblockModifier float64) bool {
	adjShotBlock := calculateModifier(shotblockModifier, ScaleFactor)
	dr := util.GenerateFloatFromRange(1, 20)
	return dr+adjShotBlock >= ToughReq
}

func CalculateFaceoff(homeFaceoffMod, awayFaceoffMod float64) bool {
	adjHomeFaceOff := calculateModifier(homeFaceoffMod, ScaleFactor)
	adjAwayFaceOff := calculateModifier(awayFaceoffMod, ScaleFactor)
	homeFaceoffVal := 10.0 + adjHomeFaceOff
	awayFaceoffVal := 10.0 + adjAwayFaceOff
	totalFaceoffVal := homeFaceoffVal + awayFaceoffVal
	faceoffCheck := util.GenerateFloatFromRange(0, totalFaceoffVal)
	return faceoffCheck <= homeFaceoffVal
}

func CalculateSafePass(passModifier, stickCheckModifier float64, longPass bool) bool {
	adjustedPassMod := calculateModifier(passModifier, ScaleFactor)
	adjustedStickCheckMod := calculateModifier(stickCheckModifier, ScaleFactor)
	req := VeryEasyReq
	if longPass {
		req = EasyReq
	}
	return util.DiceRoll(adjustedPassMod-adjustedStickCheckMod, req)
}

func GetPuckLocationAfterMiss(possessingTeam, homeTeam uint) string {
	isHome := possessingTeam == homeTeam
	list := []string{HomeGoal, HomeZone}
	if !isHome {
		list = []string{AwayGoal, AwayZone}
	}
	return util.PickFromStringList(list)
}

func RetrievePuckAfterFaceoffCheck(players []*GamePlayer, CurrentZone string, FaceoffWinID uint, homeFaceoffWin bool) uint {
	faceoffWeights, totalWeight := getPlayerWeights(homeFaceoffWin, players, FaceoffWinID, Faceoff)
	return selectPlayerIDByWeights(totalWeight, faceoffWeights)
}

func PassPuckToPlayer(gs *GameState, players []*GamePlayer, CurrentZone string) uint {
	faceoffWeights, totalWeight := getPlayerWeightsWithSystems(gs, false, players, 0, CurrentZone, Pass)
	return selectPlayerIDByWeights(totalWeight, faceoffWeights)
}

func reboundCheck(gs *GameState, players []*GamePlayer, CurrentZone string) uint {
	// Apply temporary home ice advantage during weight calculation
	for _, p := range players {
		if uint(p.TeamID) == gs.HomeTeamID {
			p.AgilityMod += 0.025 // Temporary home ice advantage
		}
	}

	reboundWeights, totalWeight := getPlayerWeightsWithSystems(gs, false, players, 0, CurrentZone, Rebound)

	// Restore original values
	for _, p := range players {
		if uint(p.TeamID) == gs.HomeTeamID {
			p.AgilityMod -= 0.025 // Restore original value
		}
	}

	// Select
	return selectPlayerIDByWeights(totalWeight, reboundWeights)
}

func selectPlayerIDByWeights(totalWeight float64, playerWeights []PlayerWeight) uint {
	selectedWeight := util.GenerateFloatFromRange(0, totalWeight)
	currWeight := 0.0
	lastID := uint(0)

	for _, rw := range playerWeights {
		currWeight += rw.Weight
		if currWeight >= selectedWeight {
			return rw.PlayerID
		}
		lastID = rw.PlayerID
	}
	return lastID
}

func getPlayerWeights(isFaceoff bool, players []*GamePlayer, faceoffWinID uint, event string) ([]PlayerWeight, float64) {
	playerWeights := []PlayerWeight{}
	totalWeight := 0.0
	for _, p := range players {
		if p.IsOut {
			continue
		}
		mod := getAttributeModifier(event, p)
		weight := 0.0
		if isFaceoff {
			weight = getFaceoffWeight(uint(p.TeamID), faceoffWinID, mod)
		} else {
			weight = getPlayerWeight(mod)
		}

		rw := PlayerWeight{
			PlayerID: p.ID,
			Weight:   weight,
		}
		playerWeights = append(playerWeights, rw)
		totalWeight += weight
	}

	// Sort weights
	sort.Slice(playerWeights, func(i, j int) bool {
		return playerWeights[i].Weight > playerWeights[j].Weight
	})

	return playerWeights, totalWeight
}

// getPlayerWeightsWithSystems - Enhanced version that considers offensive/defensive systems
func getPlayerWeightsWithSystems(gs *GameState, isFaceoff bool, players []*GamePlayer, faceoffWinID uint, CurrentZone, event string) ([]PlayerWeight, float64) {
	playerWeights := []PlayerWeight{}
	totalWeight := 0.0

	// Get system modifiers for current game state
	pc := gs.PuckCarrier
	isHomePossession := pc.TeamID == uint16(gs.HomeTeamID)
	homeModifiers, awayModifiers := GetSystemModifiersForZone(gs, isHomePossession, CurrentZone)

	for _, p := range players {
		if p.IsOut {
			continue
		}

		mod := getAttributeModifier(event, p)
		weight := 0.0

		if isFaceoff {
			weight = getFaceoffWeight(uint(p.TeamID), faceoffWinID, mod)
		} else {
			weight = getPlayerWeight(mod)
		}

		// Apply system modifiers based on player's team
		isHomePlayer := uint(p.TeamID) == gs.HomeTeamID
		if isHomePlayer {
			weight = GetSystemPlayerWeight(p, homeModifiers, weight)
		} else {
			weight = GetSystemPlayerWeight(p, awayModifiers, weight)
		}

		rw := PlayerWeight{
			PlayerID: p.ID,
			Weight:   weight,
		}
		playerWeights = append(playerWeights, rw)
		totalWeight += weight
	}

	// Sort weights
	sort.Slice(playerWeights, func(i, j int) bool {
		return playerWeights[i].Weight > playerWeights[j].Weight
	})

	return playerWeights, totalWeight
}

func getAttributeModifier(event string, p *GamePlayer) float64 {
	mod := 0.0
	switch event {
	case Rebound, Faceoff:
		mod = p.AgilityMod
	case Defense:
		mod = p.StrengthMod
	case Pass:
		mod = p.PassMod
	case ShotBlock:
		mod = p.StrengthMod
	}
	return mod
}

func getZoneID(currentZone string, homeTeamID, awayTeamID uint) (uint, uint8) {
	var zoneID uint = 0
	var zoneIDEnum uint8 = NeutralZoneID
	switch currentZone {
	case HomeGoal:
		zoneID = homeTeamID
		zoneIDEnum = HomeGoalZoneID
	case HomeZone:
		zoneID = homeTeamID
		zoneIDEnum = HomeZoneID
	case AwayZone:
		zoneID = awayTeamID
		zoneIDEnum = AwayZoneID
	case AwayGoal:
		zoneID = awayTeamID
		zoneIDEnum = AwayGoalZoneID
	}
	return zoneID, zoneIDEnum
}

func getFaceoffWeight(playerTeamID, faceoffWinID uint, mod float64) float64 {
	weight := mod
	// Players on the faceoff-winning team get a significant advantage in puck retrieval
	isOnWinningTeam := faceoffWinID == playerTeamID
	if isOnWinningTeam {
		weight += 2.5 // Substantial bonus for faceoff-winning team
	}
	return weight
}

func getPlayerWeight(mod float64) float64 {
	weight := mod
	return weight
}

func findPlayerByID(people []*GamePlayer, id uint) (*GamePlayer, bool) {
	player := GamePlayer{}
	for _, person := range people {
		if person.ID == id {
			return person, true // Return the found person and true
		}
	}
	return &player, false // Return nil if not found
}

func selectDefendingPlayer(gs *GameState, defendingTeamID uint) *GamePlayer {
	playerList := getFullPlayerListByTeamID(defendingTeamID, gs)
	playerMap := getGameplayerMap(playerList)
	playerWeights, totalWeight := getPlayerWeightsWithSystems(gs, false, playerList, 0, gs.PuckLocation, Defense)
	playerID := selectPlayerIDByWeights(totalWeight, playerWeights)
	player := playerMap[playerID]
	return player
}

func selectBlockingPlayer(gs *GameState, defendingTeamID uint) *GamePlayer {
	playerList := getFullPlayerListByTeamID(defendingTeamID, gs)
	playerMap := getGameplayerMap(playerList)
	playerWeights, totalWeight := getPlayerWeightsWithSystems(gs, false, playerList, 0, gs.PuckLocation, ShotBlock)
	playerID := selectPlayerIDByWeights(totalWeight, playerWeights)
	player := playerMap[playerID]
	return player
}

func getFullPlayerListByTeamID(teamID uint, gs *GameState) []*GamePlayer {
	playerList := []*GamePlayer{}
	isHome := teamID == gs.HomeTeamID
	forwardStrategy := gs.GetLineStrategy(isHome, 1)
	defenderStrategy := gs.GetLineStrategy(isHome, 2)

	playbook := gs.GetPlaybook(isHome)
	currentForwards := forwardStrategy.Players
	currentDefenders := defenderStrategy.Players

	// Filter out player depending on if there's a power play
	for idx, p := range currentForwards {
		if (idx == 0 && playbook.CenterOut) || (idx == 1 && playbook.Forward1Out) || (idx == 2 && playbook.Forward2Out) {
			continue
		}
		playerList = append(playerList, p)
	}

	for idx, p := range currentDefenders {
		if (idx == 0 && playbook.Defender1Out) || (idx == 1 && playbook.Defender2Out) {
			continue
		}
		playerList = append(playerList, p)
	}

	return playerList
}

// getAvailablePlayers -- For passes and faceoffs
func getAvailablePlayers(possessingPlayerID uint, line []*GamePlayer) []*GamePlayer {
	list := []*GamePlayer{}
	for _, p := range line {
		if p.ID != possessingPlayerID && !p.IsOut {
			list = append(list, p)
		}
	}
	return list
}

func getGameplayerMap(list []*GamePlayer) map[uint]*GamePlayer {
	playerMap := make(map[uint]*GamePlayer)

	for _, p := range list {
		playerMap[p.ID] = p
	}

	return playerMap
}

func getDefendingTeamID(teamID, ht, at uint) uint {
	if teamID == ht {
		return at
	}
	return ht
}

func isHomeTeam(teamID, ht uint) bool {
	return teamID == ht
}

func getNextZone(gs *GameState) string {
	currentZone := gs.PuckLocation
	nextZone := NeutralZone
	pb := gs.PuckCarrier
	isHT := isHomeTeam(uint(pb.TeamID), gs.HomeTeamID)
	if isHT && currentZone == NeutralZone {
		nextZone = AwayZone
	} else if isHT && currentZone == AwayZone {
		nextZone = AwayGoal
	} else if isHT && currentZone == HomeZone {
		nextZone = NeutralZone
	} else if isHT && currentZone == HomeGoal {
		nextZone = HomeZone
	} else if !isHT && currentZone == AwayGoal {
		nextZone = AwayZone
	} else if !isHT && currentZone == AwayZone {
		nextZone = NeutralZone
	} else if !isHT && currentZone == NeutralZone {
		nextZone = HomeZone
	} else if !isHT && currentZone == HomeZone {
		nextZone = HomeGoal
	}
	return nextZone
}

func getPreviousZone(gs *GameState) string {
	currentZone := gs.PuckLocation
	prevZone := NeutralZone
	pb := gs.PuckCarrier
	isHT := isHomeTeam(uint(pb.TeamID), gs.HomeTeamID)
	if isHT && currentZone == AwayGoal {
		prevZone = AwayZone
	} else if isHT && currentZone == AwayZone {
		prevZone = NeutralZone
	} else if isHT && currentZone == NeutralZone {
		prevZone = HomeZone
	} else if isHT && currentZone == HomeZone {
		prevZone = HomeGoal
	} else if !isHT && currentZone == HomeGoal {
		prevZone = HomeZone
	} else if !isHT && currentZone == HomeZone {
		prevZone = NeutralZone
	} else if !isHT && currentZone == NeutralZone {
		prevZone = AwayZone
	} else if !isHT && currentZone == AwayZone {
		prevZone = AwayGoal
	}
	return prevZone
}

// Define the calculateModifier function with logarithmic scaling
func calculateModifier(attribute, scaleFactor float64) float64 {
	return scaleFactor * math.Log(attribute+1)
}

func calculateAttributeModifier(attribute, scaleFactor float64, isHome bool, hra float64) float64 {
	hraMod := 1.0
	penalty := 0.0
	if !isHome {
		penalty = 0.03 + (0.12 * hra)
		hraMod = hraMod - penalty
	}
	// return scaleFactor * math.Log(float64(attribute)+1)
	return (attribute / 10) * hraMod
}
