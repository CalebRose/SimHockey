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
	// numStr := util.ConvertFloatToString(num)
	// goalTendString := util.ConvertFloatToString(baseCheck + goaltending + goalieMod)
	// fmt.Println("Num: " + numStr + " Goaltending: " + goalTendString)
	return num > (baseCheck + adjustedGoalieScore)
}

func CalculateAccuracy(accuracyModifier float64, isCloseShot bool) bool {
	max := 20.0
	min := 1.0
	req := DiffReq
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
	return dr+adjShotBlock >= DiffReq
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

func CalculateSafePass(passModifier, stickCheckModifier float64) bool {
	adjustedPassMod := calculateModifier(passModifier, ScaleFactor)
	adjustedStickCheckMod := calculateModifier(stickCheckModifier, ScaleFactor)
	return util.DiceRoll(adjustedPassMod-adjustedStickCheckMod, EasyReq)
}

func GetPuckLocationAfterMiss(possessingTeam, homeTeam uint) string {
	isHome := possessingTeam == homeTeam
	list := []string{HomeGoal, HomeZone}
	if !isHome {
		list = []string{AwayGoal, AwayZone}
	}
	return util.PickFromStringList(list)
}

func RetrievePuckAfterFaceoffCheck(players []GamePlayer, CurrentZone string, HomeTeamID, AwayTeamID, FaceoffWinID uint, homeFaceoffWin bool) uint {
	zoneID, _ := getZoneID(CurrentZone, HomeTeamID, AwayTeamID)
	faceoffWeights, totalWeight := getPlayerWeights(true, homeFaceoffWin, players, zoneID, FaceoffWinID, CurrentZone, Faceoff)
	return selectPlayerIDByWeights(totalWeight, faceoffWeights)
}

func PassPuckToPlayer(players []GamePlayer, CurrentZone string, HomeTeamID, AwayTeamID uint) uint {
	zoneID, _ := getZoneID(CurrentZone, HomeTeamID, AwayTeamID)
	faceoffWeights, totalWeight := getPlayerWeights(false, false, players, zoneID, 0, CurrentZone, Pass)
	return selectPlayerIDByWeights(totalWeight, faceoffWeights)
}

func reboundCheck(players []GamePlayer, CurrentZone string, HomeTeamID, AwayTeamID uint) uint {
	zoneID, _ := getZoneID(CurrentZone, HomeTeamID, AwayTeamID)
	reboundWeights, totalWeight := getPlayerWeights(false, false, players, zoneID, 0, CurrentZone, Rebound)

	// Select
	return selectPlayerIDByWeights(totalWeight, reboundWeights)
}

func selectPlayerIDByWeights(totalWeight float64, playerWeights []PlayerWeight) uint {
	selectedWeight := util.GenerateFloatFromRange(0, totalWeight)
	currWeight := 0.0
	lastID := uint(0)

	for _, rw := range playerWeights {
		currWeight += rw.Weight
		if currWeight <= selectedWeight {
			return rw.PlayerID
		}
		lastID = rw.PlayerID
	}
	return lastID
}

func getPlayerWeights(isFaceoff, homeTeamFaceoffWin bool, players []GamePlayer, zoneID, faceoffWinID uint, CurrentZone, event string) ([]PlayerWeight, float64) {
	playerWeights := []PlayerWeight{}
	totalWeight := 0.0
	for _, p := range players {
		if p.IsOut {
			continue
		}
		mod := getAttributeModifier(event, p)
		weight := 0.0
		if isFaceoff {
			weight = getFaceoffWeight(uint(p.TeamID), faceoffWinID, mod, CurrentZone, homeTeamFaceoffWin)
		} else {
			weight = getPlayerWeight(uint(p.TeamID), zoneID, mod, p.Position, CurrentZone)
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

func getAttributeModifier(event string, p GamePlayer) float64 {
	mod := 0.0
	if event == Rebound || event == Faceoff {
		mod = p.AgilityMod
	} else if event == Defense {
		mod = p.StrengthMod
	} else if event == Pass {
		mod = p.PassMod
	} else if event == ShotBlock {
		mod = p.StrengthMod
	}
	return mod
}

func getZoneID(currentZone string, homeTeamID, awayTeamID uint) (uint, uint8) {
	var zoneID uint = 0
	var zoneIDEnum uint8 = NeutralZoneID
	if currentZone == HomeGoal {
		zoneID = homeTeamID
		zoneIDEnum = HomeGoalZoneID
	} else if currentZone == HomeZone {
		zoneID = homeTeamID
		zoneIDEnum = HomeZoneID
	} else if currentZone == AwayZone {
		zoneID = awayTeamID
		zoneIDEnum = AwayZoneID
	} else if currentZone == AwayGoal {
		zoneID = awayTeamID
		zoneIDEnum = AwayGoalZoneID
	}
	return zoneID, zoneIDEnum
}

func getFaceoffWeight(playerTeamID, faceoffWinID uint, mod float64, CurrentZone string, homeTeamFaceoffWin bool) float64 {
	weight := mod
	defendingPlayer := faceoffWinID == playerTeamID && homeTeamFaceoffWin
	if defendingPlayer {
		weight += 1.2
	}
	if (CurrentZone == HomeGoal || CurrentZone == HomeZone) && defendingPlayer {
		weight++
	} else if (CurrentZone == AwayGoal || CurrentZone == AwayZone) && defendingPlayer {
		weight++
	}
	return weight
}

func getPlayerWeight(playerTeamID, ZoneID uint, mod float64, Position, CurrentZone string) float64 {
	weight := mod
	defendingPlayer := playerTeamID == ZoneID
	if defendingPlayer {
		weight += 0.5
	}
	if (CurrentZone == HomeGoal || CurrentZone == HomeZone) && Position == Defender && defendingPlayer {
		weight += 0.5
	} else if (CurrentZone == AwayGoal || CurrentZone == AwayZone) && Position == Defender && defendingPlayer {
		weight += 0.5
	}
	return weight
}

func findPlayerByID(people []GamePlayer, id uint) (*GamePlayer, bool) {
	for _, person := range people {
		if person.ID == id {
			return &person, true // Return the found person and true
		}
	}
	return nil, false // Return nil if not found
}

func selectDefendingPlayer(gs *GameState, defendingTeamID uint) GamePlayer {
	playerList := getFullPlayerListByTeamID(defendingTeamID, gs)
	playerMap := getGameplayerMap(playerList)
	zoneID, _ := getZoneID(gs.PuckLocation, gs.HomeTeamID, gs.AwayTeamID)

	playerWeights, totalWeight := getPlayerWeights(false, false, playerList, zoneID, 0, gs.PuckLocation, Defense)
	playerID := selectPlayerIDByWeights(totalWeight, playerWeights)
	return playerMap[playerID]
}

func selectBlockingPlayer(gs *GameState, defendingTeamID uint) GamePlayer {
	playerList := getFullPlayerListByTeamID(defendingTeamID, gs)
	playerMap := getGameplayerMap(playerList)
	zoneID, _ := getZoneID(gs.PuckLocation, gs.HomeTeamID, gs.AwayTeamID)

	playerWeights, totalWeight := getPlayerWeights(false, false, playerList, zoneID, 0, gs.PuckLocation, ShotBlock)
	playerID := selectPlayerIDByWeights(totalWeight, playerWeights)
	return playerMap[playerID]
}

func getFullPlayerListByTeamID(teamID uint, gs *GameState) []GamePlayer {
	playerList := []GamePlayer{}
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

	playerList = append(playerList, currentDefenders...)

	return playerList
}

// getAvailablePlayers -- For passes and faceoffs
func getAvailablePlayers(possessingPlayerID uint, line []GamePlayer) []GamePlayer {
	list := []GamePlayer{}
	for _, p := range line {
		if p.ID != possessingPlayerID && !p.IsOut {
			list = append(list, p)
		}
	}
	return list
}

func getGameplayerMap(list []GamePlayer) map[uint]GamePlayer {
	playerMap := make(map[uint]GamePlayer)

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

// Define the calculateModifier function with logarithmic scaling
func calculateModifier(attribute float64, scaleFactor float64) float64 {
	return scaleFactor * math.Log(float64(attribute)+1)
}

func calculateAttributeModifier(attribute float64, scaleFactor float64) float64 {
	// return scaleFactor * math.Log(float64(attribute)+1)
	return attribute / 10
}
