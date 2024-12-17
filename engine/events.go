package engine

import (
	"math"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/structs"
)

func HandleBaseEvents(gs *GameState) {
	pb := gs.PuckCarrier
	homePossession := pb.TeamID == uint16(gs.HomeTeamID)
	switch gs.PuckLocation {
	case HomeGoal:
		handleGoalZoneEvents(gs, homePossession, false)
	case HomeZone, AwayZone:
		handleZoneEvents(gs, homePossession, gs.PuckLocation == HomeZone)
	case NeutralZone:
		handleNeutralZoneEvents(gs)
	case AwayGoal:
		handleGoalZoneEvents(gs, homePossession, true)
	}
}

func handleGoalZoneEvents(gs *GameState, homePossession bool, isAwayGoal bool) {
	if (homePossession && isAwayGoal) || (!homePossession && !isAwayGoal) {
		handleOffensiveGoalZoneEvents(gs)
	} else {
		handleDefensiveGoalZoneEvents(gs)
	}
}

func handleZoneEvents(gs *GameState, homePossession bool, isHomeZone bool) {
	if (homePossession && isHomeZone) || (!homePossession && !isHomeZone) {
		handleDefensiveZoneEvents(gs)
	} else {
		handleOffensiveZoneEvents(gs)
	}
}

func handleOffensiveGoalZoneEvents(gs *GameState) {
	// Passing, close shots, defensive checks
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	attackStrategy := gs.GetLineStrategy(isHome, 1)
	defendStrategy := gs.GetLineStrategy(!isHome, 2)
	slapshot := 0
	pass := 0
	stickCheck := 0
	bodyCheck := 0
	penalty := 1
	slapshot = int(attackStrategy.AGZShot) + int(pb.AGZShot) + int(gs.Momentum)
	pass = int(attackStrategy.AGZPass) + int(pb.AGZPass)
	stickCheck = int(defendStrategy.DGZStickCheck)
	bodyCheck = int(defendStrategy.DGZBodyCheck)
	totalSkill := slapshot + stickCheck + bodyCheck + pass + penalty
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := float64(stickCheckCutoff) + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	shotCutoff := passCheckCutoff + float64(slapshot)
	penaltyCutoff := shotCutoff + 1
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs)
	} else if dr <= shotCutoff {
		HandleShot(gs, true)
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
	}
}

func handleOffensiveZoneEvents(gs *GameState) {
	// Movement, passing, defense, and long range shots
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	attackStrategy := gs.GetLineStrategy(isHome, 1)
	defendStrategy := gs.GetLineStrategy(!isHome, 2)
	penalty := 1
	slapshot := int(attackStrategy.AZSlapshot) + int(pb.AZSlapshot) + int(gs.Momentum)
	wristshot := int(attackStrategy.AZWristshot) + int(pb.AZWristshot) + int(gs.Momentum)
	agility := int(attackStrategy.AZAgility) + int(pb.AZAgility)
	pass := int(attackStrategy.AZPass) + int(pb.AZPass)
	stickCheck := int(defendStrategy.DZStickCheck)
	bodyCheck := int(defendStrategy.DZBodyCheck)
	totalSkill := slapshot + wristshot + stickCheck + bodyCheck + pass + penalty + int(agility)
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	agilityCutoff := passCheckCutoff + float64(agility)
	slapshotCutoff := agilityCutoff + float64(slapshot)
	wristshotCutoff := slapshotCutoff + float64(wristshot)
	penaltyCutoff := wristshotCutoff + 1.0
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs)
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
	} else if dr <= slapshotCutoff {
		HandleShot(gs, true)
	} else if dr <= wristshotCutoff {
		HandleShot(gs, false)
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
	}
}

func handleDefensiveGoalZoneEvents(gs *GameState) {
	// Movement, pass, defensive checks
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	attackStrategy := gs.GetLineStrategy(isHome, 1)
	defendStrategy := gs.GetLineStrategy(!isHome, 2)
	penalty := 1
	agility := int(attackStrategy.DGZAgility) + int(pb.DGZAgility)
	pass := int(attackStrategy.DGZPass) + int(pb.DGZPass)
	stickCheck := int(defendStrategy.AGZStickCheck)
	bodyCheck := int(defendStrategy.AGZBodyCheck)
	faceOffCheck := 0
	if pb.Position == Goalie {
		agility = 0
		stickCheck = 0
		bodyCheck = 0
		faceOffCheck = 20
	}
	totalSkill := stickCheck + bodyCheck + pass + penalty + int(agility) + faceOffCheck
	faceoffCutoff := float64(faceOffCheck)
	stickCheckCutoff := faceoffCutoff + float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	agilityCutoff := passCheckCutoff + float64(agility)
	penaltyCutoff := agilityCutoff + 1.0
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= faceoffCutoff {
		gs.SetFaceoffOnCenterIce(true)
		newZone := HomeZone
		if !isHome {
			newZone = AwayZone
		}
		gs.SetNewZone(newZone)
		HandleFaceoff(gs)
	} else if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs)
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
	}
}

func handleDefensiveZoneEvents(gs *GameState) {
	// Movement, passing, defensive checks made by opposing offensive forwards
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	attackStrategy := gs.GetLineStrategy(isHome, 1)
	defendStrategy := gs.GetLineStrategy(!isHome, 2)
	penalty := 1
	agility := int(attackStrategy.DZAgility) + int(pb.DZAgility)
	pass := int(attackStrategy.DZPass) + int(pb.DZPass)
	stickCheck := int(defendStrategy.AZStickCheck)
	bodyCheck := int(defendStrategy.AZBodyCheck)
	totalSkill := stickCheck + bodyCheck + pass + penalty + int(agility)
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	agilityCutoff := passCheckCutoff + float64(agility)
	penaltyCutoff := agilityCutoff + 1.0
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs)
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
	}
}

func handleNeutralZoneEvents(gs *GameState) {
	// Movement, passing, defensive checks
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	attackStrategy := gs.GetLineStrategy(isHome, 1)
	defendStrategy := gs.GetLineStrategy(!isHome, 2)
	penalty := 1
	agility := int(attackStrategy.NAgility) + int(pb.NAgility)
	pass := int(attackStrategy.NPass) + int(pb.NPass)
	stickCheck := int(defendStrategy.NStickCheck)
	bodyCheck := int(defendStrategy.NBodyCheck)
	totalSkill := stickCheck + bodyCheck + pass + penalty + int(agility)
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	agilityCutoff := passCheckCutoff + float64(agility)
	penaltyCutoff := agilityCutoff + 1.0
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs)
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
	}
}

// handleDefenseCheck -- For preventing the defense from getting the puck
func handleDefenseCheck(gs *GameState, isStickCheck bool) {
	pb := gs.PuckCarrier
	// Select player on defense
	defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	defender := selectDefendingPlayer(gs, defendingTeamID)

	chance := CalculatePenaltyChance()
	if chance {
		shouldReturn := handlePenalty(gs, !isStickCheck, defender)
		if shouldReturn {
			return
		}
	}
	diceRoll := util.GenerateIntFromRange(1, 20)

	// Critical checks
	if diceRoll == CritFail {
		// Defender gets puck
		gs.SetPuckBearer(defender)
		return
	} else if diceRoll == CritSuccess {
		// Defense DOES NOT get puck
		return
	} else {
		// Determine if physical check or non-physical check
		puckHandling := DiffReq + pb.HandlingMod
		if isStickCheck {
			puckHandling -= defender.StickCheckMod
		} else {
			puckHandling -= defender.BodyCheckMod
		}

		puckHandling = math.Max(puckHandling, 1.0)

		if float64(diceRoll) < puckHandling {
			// fmt.Println(pb.Team + " " + pb.FirstName + " " + pb.LastName + " holds onto the puck!")
		} else {
			// Logger(defender.FirstName + " GETS THE PUCK FOR " + defender.Team + "!")
			gs.SetPuckBearer(defender)
		}
	}
	secondsConsumed := util.GenerateIntFromRange(2, 4)
	gs.SetSecondsConsumed(uint16(secondsConsumed))
}

func handleAgilityCheck(gs *GameState) {
	// Logger("Agility Check in " + gs.PuckLocation + ".")
	// Get Current Zone
	nextZone := getNextZone(gs)

	pb := gs.PuckCarrier
	agilityMod := pb.AgilityMod
	momentumMod := gs.Momentum
	critCheck := util.GenerateIntFromRange(1, 20)
	secondsConsumed := util.GenerateIntFromRange(1, 4)
	defenseCheck := true
	if critCheck == CritFail {
		secondsConsumed += 3
	} else if critCheck == CritSuccess || float64(critCheck) > DiffReq+agilityMod+momentumMod {
		defenseCheck = false
	}
	if critCheck == CritSuccess {
		gs.TriggerBreakaway()
	}

	if defenseCheck {
		defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
		defender := selectDefendingPlayer(gs, defendingTeamID)
		diceRoll := util.GenerateIntFromRange(1, 20)
		puckHandling := DiffReq + pb.HandlingMod
		coinFlip := util.CoinFlip()
		if coinFlip == Heads {
			puckHandling -= defender.BodyCheckMod
		} else {
			puckHandling -= defender.StickCheckMod
		}

		chance := CalculatePenaltyChance()
		if chance {
			// Logger("Defensive Stick Check in " + gs.PuckLocation + ".")
			// Logger("Defensive Body Check in " + gs.PuckLocation + ".")
			shouldReturn := handlePenalty(gs, coinFlip == Heads, defender)
			if shouldReturn {
				return
			}
		}

		puckHandling = math.Max(puckHandling, 1.0)

		if float64(diceRoll) < puckHandling {
			// fmt.Println(pb.Team + " " + pb.FirstName + " " + pb.LastName + " holds onto the puck!")
		} else {
			RecordPlay(gs, AgilityCheckID, DefenseStopAgilityID, 0, 0, 0, 0, pb.ID, 0, 0, defender.ID, 0, false)
			gs.SetPuckBearer(defender)
			// Logger(defender.FirstName + " GETS THE PUCK FOR " + defender.Team + "!")
			return
		}
	}

	// Logger(pb.FirstName + " " + pb.LastName + " moves up to " + nextZone + "!")
	// Move up zone
	gs.SetNewZone(nextZone)
	// Seconds consumed
	gs.SetSecondsConsumed(uint16(secondsConsumed))
}

func handlePenalty(gs *GameState, isBodyCheck bool, defender GamePlayer) bool {
	zoneID := 0
	switch gs.PuckLocation {
	case HomeGoal, AwayGoal:
		zoneID = 2
	case HomeZone, AwayZone:
		zoneID = 1
	}
	var penalty Penalty
	if isBodyCheck {
		penalty = SelectPenalty(defender, uint(zoneID), BodyCheck)
	} else {
		penalty = SelectPenalty(defender, uint(zoneID), StickCheck)
	}

	if penalty.PenaltyID > 0 {
		ApplyPenalty(gs, penalty, defender)
		return true
	}
	return false
}

func handlePassCheck(gs *GameState) {
	// Logger("PASS EVENT")
	pb := gs.PuckCarrier

	// Roll to see if puck is intercepted by defense
	defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	defender := selectDefendingPlayer(gs, defendingTeamID)

	safePass := CalculateSafePass(pb.PassMod, defender.StickCheckMod)

	if !safePass {
		RecordPlay(gs, PassCheckID, InterceptedPassID, 0, 0, 0, 0, pb.ID, 0, 0, defender.ID, 0, false)
		gs.SetPuckBearer(defender)
		// Logger(defender.FirstName + " INTERCEPTS THE PASS FOR " + defender.Team + "!")
		return
	}

	// Get available player on own team
	playerList := getFullPlayerListByTeamID(uint(pb.TeamID), gs)
	filteredList := getAvailablePlayers(pb.ID, playerList)

	receivingPlayer := PassPuckToPlayer(filteredList, gs.PuckLocation, gs.HomeTeamID, gs.AwayTeamID)
	retrievingPlayer, _ := findPlayerByID(playerList, receivingPlayer)
	HandleMissingPlayer(*retrievingPlayer, "PASSING PUCK")
	secondsConsumed := util.GenerateIntFromRange(1, 3)
	gs.SetSecondsConsumed(uint16(secondsConsumed))
	RecordPlay(gs, PassCheckID, InterceptedPassID, 0, 0, 0, 0, pb.ID, receivingPlayer, 0, defender.ID, 0, false)
	gs.SetPuckBearer(*retrievingPlayer)
	// Logger(pb.FirstName + " PASSES THE PUCK TO " + retrievingPlayer.FirstName + " " + retrievingPlayer.LastName + "!")
}

func handlePenalties(gs *GameState) {
	// Determine penalty type
	// Minor, major, misconduct, game misconduct, match
	// Logger("PENALTY WOULD BE HERE!")
	pb := gs.PuckCarrier
	defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	zoneID := 0
	switch gs.PuckLocation {
	case HomeGoal, AwayGoal:
		zoneID = 2
	case HomeZone, AwayZone:
		zoneID = 1
	}

	player := selectDefendingPlayer(gs, defendingTeamID)
	secondPlayer := GamePlayer{}
	penaltyTypeID := GeneralPenaltyID
	penaltyType := General
	diceRoll := util.DiceRoll(0, 20)
	if diceRoll {
		penaltyTypeID = FightPenaltyID
		penaltyType = Fight
		// If a fight occurs, then two players should probably get placed in the penalty box.
		secondPlayer = selectDefendingPlayer(gs, uint(pb.TeamID))
	}

	penalty := SelectPenalty(player, uint(zoneID), penaltyType)
	if penalty.PenaltyID == 0 {
		return
	}
	RecordPlay(gs, PenaltyCheckID, penaltyTypeID, 0, 0, 0, uint8(penalty.PenaltyID), pb.ID, 0, 0, player.ID, secondPlayer.ID, false)

	// Apply Penalty to Player
	ApplyPenalty(gs, penalty, player)
	ApplyPenalty(gs, penalty, secondPlayer)
}

func HandleFaceoff(gs *GameState) {
	// Get Centers from current lineups
	homeCenter := gs.GetCenter(true)
	awayCenter := gs.GetCenter(false)
	HandleMissingPlayer(homeCenter, "HandleFaceoff Home Center")
	HandleMissingPlayer(awayCenter, "HandleFaceoff Away Center")
	homeFaceoffWin := CalculateFaceoff(homeCenter.FaceoffMod, awayCenter.FaceoffMod)
	faceOffWinID := homeCenter.TeamID
	// Away wins faceoff
	if !homeFaceoffWin {
		faceOffWinID = awayCenter.TeamID
		awayCenter.Stats.AddFaceoff(true)
		homeCenter.Stats.AddFaceoff(false)
		gs.AwayTeamStats.AddFaceoff(true)
		gs.HomeTeamStats.AddFaceoff(false)
	} else {
		homeCenter.Stats.AddFaceoff(true)
		awayCenter.Stats.AddFaceoff(false)
		gs.HomeTeamStats.AddFaceoff(true)
		gs.AwayTeamStats.AddFaceoff(false)
	}
	if gs.FaceoffOnCenterIce {
		gs.SetFaceoffOnCenterIce(false)
	}
	// Select player who gets puck after faceoff
	HandleFaceoffRetrieval(gs, homeFaceoffWin, uint(faceOffWinID))
}

func HandleFaceoffRetrieval(gs *GameState, homeFaceoffWin bool, faceoffWinID uint) {
	puckLocation := NeutralZone
	playerList := []GamePlayer{}
	// Get Available Players in Home Forward Line
	hgs := gs.HomeStrategy
	ags := gs.AwayStrategy
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(true).ID, hgs.Forwards[hgs.CurrentForwards].Players)...)
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(true).ID, hgs.Defenders[hgs.CurrentDefenders].Players)...)
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(false).ID, ags.Forwards[ags.CurrentForwards].Players)...)
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(false).ID, ags.Defenders[ags.CurrentDefenders].Players)...)

	faceoffRetrievalCheck := RetrievePuckAfterFaceoffCheck(playerList, puckLocation, gs.HomeTeamID, gs.AwayTeamID, faceoffWinID, homeFaceoffWin)
	retrievingPlayer, _ := findPlayerByID(playerList, faceoffRetrievalCheck)
	HandleMissingPlayer(*retrievingPlayer, "REBOUNDING AFTER FACEOFF")
	gs.SetPuckBearer(*retrievingPlayer)
	outcomeID := HomeFaceoffWinID
	if !homeFaceoffWin {
		outcomeID = AwayFaceoffWinID
	}
	RecordPlay(gs, FaceoffID, outcomeID, 0, 0, 0, 0, retrievingPlayer.ID, 0, 0, 0, 0, false)
	// Logger(retrievingPlayer.Team + " gets the puck from the faceoff with " + retrievingPlayer.FirstName + " " + retrievingPlayer.LastName + " in possession!")
}

func HandleReboundAfterShot(gs *GameState, eventID uint8, outcomeID uint8, puckCarrierID uint, goalieID uint) {
	puckLocation := GetPuckLocationAfterMiss(1, 1)
	reboundPlayerList := []GamePlayer{}
	hgs := gs.HomeStrategy
	ags := gs.AwayStrategy
	if puckLocation == HomeGoal {
		reboundPlayerList = append(reboundPlayerList, hgs.Forwards[hgs.CurrentForwards].Players...)
		reboundPlayerList = append(reboundPlayerList, ags.Defenders[ags.CurrentDefenders].Players...)
	} else if puckLocation == HomeZone || puckLocation == AwayZone {
		reboundPlayerList = append(reboundPlayerList, hgs.Forwards[hgs.CurrentForwards].Players...)
		reboundPlayerList = append(reboundPlayerList, ags.Defenders[ags.CurrentDefenders].Players...)
		reboundPlayerList = append(reboundPlayerList, ags.Forwards[ags.CurrentForwards].Players...)
		reboundPlayerList = append(reboundPlayerList, hgs.Defenders[hgs.CurrentDefenders].Players...)
	} else if puckLocation == AwayGoal {
		reboundPlayerList = append(reboundPlayerList, ags.Forwards[ags.CurrentForwards].Players...)
		reboundPlayerList = append(reboundPlayerList, hgs.Defenders[hgs.CurrentDefenders].Players...)
	}
	reboundCheck := reboundCheck(reboundPlayerList, puckLocation, gs.HomeTeamID, gs.AwayTeamID)
	reboundingPlayer, _ := findPlayerByID(reboundPlayerList, reboundCheck)
	HandleMissingPlayer(*reboundingPlayer, "REBOUNDING AFTER SHOT")
	// Record Play after inaccurate shot
	RecordPlay(gs, eventID, outcomeID, 0, 0, 0, 0, puckCarrierID, reboundingPlayer.ID, 0, 0, goalieID, false)
	gs.SetPuckBearer(*reboundingPlayer)
	// Logger(reboundingPlayer.FirstName + " " + reboundingPlayer.LastName + " gets the puck for " + reboundingPlayer.Team + " on the rebound!")
}

func HandleShot(gs *GameState, isCloseShot bool) {
	if isCloseShot {
		// Logger("CLOSE SHOT EVENT")
	} else {
		// Logger("LONG SHOT EVENT")
	}
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	goalie := GamePlayer{}
	goalieStrategy := gs.GetLineStrategy(!isHome, 3)
	goalie = goalieStrategy.Players[0]
	accuracy := 0
	power := 0
	goalieMod := 0.0
	goalKeepingMod := 0.0
	eventTypeID := SlapshotCheckID
	if isCloseShot {
		accuracy = int(pb.SlapshotAccMod)
		power = int(pb.SlapshotMod)
		goalieMod = goalie.StrengthMod
		goalKeepingMod = goalie.GoalkeepingMod
	} else {
		eventTypeID = WristshotCheckID
		accuracy = int(pb.WristShotAccMod)
		power = int(pb.WristShotMod)
		goalieMod = goalie.AgilityMod
		goalKeepingMod = float64(goalie.GoalieVision)

		// Detect shotblocking
		defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
		defender := selectBlockingPlayer(gs, defendingTeamID)
		shotBlocked := CalculateShotBlock(defender.ShotblockingMod)
		if shotBlocked {
			secondsConsumed := util.GenerateIntFromRange(2, 5)
			gs.SetSecondsConsumed(uint16(secondsConsumed))
			RecordPlay(gs, eventTypeID, ShotBlockedID, 0, 0, 0, 0, pb.ID, 0, 0, defender.ID, goalie.ID, false)
			gs.SetPuckBearer(defender)
			return
		}
	}
	secondsConsumed := util.GenerateIntFromRange(2, 5)
	gs.SetSecondsConsumed(uint16(secondsConsumed))
	accuracyCheck := CalculateAccuracy(float64(accuracy))
	if !accuracyCheck {
		// Might need to record who gets the rebound...
		HandleReboundAfterShot(gs, eventTypeID, InAccurateShotID, pb.ID, goalie.ID)
		pb.Stats.AddShot(false, false, false, false, false)
		return
	}

	baseCheck := 17.25
	// If Overtime, open more opportunities for shooting
	if gs.IsOvertime {
		baseCheck -= 0.75
	}

	shotAttempt := CalculateShot(float64(power),
		float64(gs.PuckCarrier.OneTimerMod)*gs.Momentum,
		goalKeepingMod,
		goalieMod, baseCheck)

	gs.ResetMomentum()

	if shotAttempt {
		// fmt.Println("SCORE BY " + pb.FirstName + " " + pb.LastName + " for " + pb.Team + "!")
		RecordPlay(gs, eventTypeID, ShotOnGoalID, 0, 0, 0, 0, pb.ID, 0, gs.AssistingPlayer.ID, 0, goalie.ID, false)
		gs.IncrementScore(isHome)
		goalie.Stats.AddShotAgainst(true)
	} else {
		gs.AddShots(isHome)
		goalie.Stats.AddShotAgainst(false)
		pb.Stats.AddShot(false, false, false, false, false)
		RecordPlay(gs, eventTypeID, GoalieSaveID, 0, 0, 0, 0, pb.ID, 0, gs.AssistingPlayer.ID, 0, goalie.ID, false)
		gs.SetPuckBearer(goalie)
		// fmt.Println("SAVE BY " + goalie.Team + " " + goalie.FirstName + " " + goalie.LastName + "!")
	}
}

func HandleShootoutAttempt(gs *GameState) {
	pb := gs.PuckCarrier
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	goalie := GamePlayer{}
	goalieStrategy := gs.GetLineStrategy(!isHome, 3)
	goalie = goalieStrategy.Players[0]
	accuracy := int(pb.SlapshotAccMod)
	power := int(pb.SlapshotMod)
	goalieMod := goalie.StrengthMod
	goalKeepingMod := goalie.GoalkeepingMod

	accuracyCheck := CalculateAccuracy(float64(accuracy))
	if !accuracyCheck {
		// If fails the accuracy check, add stat, then return
		return
	}

	baseCheck := 12.0
	shotAttempt := CalculateShot(float64(power),
		float64(gs.PuckCarrier.OneTimerMod)*ShootoutMomenumModifier,
		goalKeepingMod,
		goalieMod, baseCheck)

	if shotAttempt {
		// fmt.Println("SCORE BY " + pb.FirstName + " " + pb.LastName + " for " + pb.Team + "!")
		gs.IncrementShootoutScore(isHome)
	}
}

func HandleOvertimeShootout(gs *GameState) {
	// Shootouts should be close shots
	isRepeat := false
	// Organize shootout queue. Home, then away, repeat (H,A,H,A,H,A,H,A...)
	shootoutQueue := formShootoutQueue(gs.HomeStrategy, gs.AwayStrategy)
	// Run For Loop for Shootout. Once six players from each team have made an attempt, compare the shootout scores
	// If score is still the same, keep running through loop

	for gs.HomeTeamShootoutScore == gs.AwayTeamShootoutScore {
		for idx, player := range shootoutQueue {
			// If we go beyond the first six shots and the score is still the same, continue
			if (idx > 5 && gs.HomeTeamShootoutScore != gs.AwayTeamShootoutScore && !isRepeat) || (isRepeat && gs.HomeTeamShootoutScore != gs.AwayTeamShootoutScore) {
				break
			}
			gs.SetPuckBearer(player)
			HandleShootoutAttempt(gs)
		}
	}
}

func formShootoutQueue(homeGP, awayGP GamePlaybook) []GamePlayer {
	queue := []GamePlayer{}

	homeForwards := homeGP.Forwards
	awayForwards := awayGP.Forwards
	homeDefenders := homeGP.Defenders
	awayDefenders := awayGP.Defenders
	homeBench := homeGP.BenchPlayers
	awayBench := awayGP.BenchPlayers

	// Loop through forward lines
	// Each team has the same number of lines AND players
	for i := 0; i < len(homeForwards); i++ {
		hFLine := homeForwards[i]
		aFLine := awayForwards[i]
		for j := 0; j < len(hFLine.Players); j++ {
			homePlayer := hFLine.Players[j]
			awayPlayer := aFLine.Players[j]
			queue = append(queue, homePlayer)
			queue = append(queue, awayPlayer)
		}
	}

	// Add defending players
	for i := 0; i < len(homeDefenders); i++ {
		hFLine := homeDefenders[i]
		aFLine := awayDefenders[i]
		for j := 0; j < len(hFLine.Players); j++ {
			homePlayer := hFLine.Players[j]
			awayPlayer := aFLine.Players[j]
			queue = append(queue, homePlayer)
			queue = append(queue, awayPlayer)
		}
	}

	// Add Bench. Because the bench size may vary,
	// only loop based on the smallest size of bench
	for i := 0; i < len(homeBench) && i < len(awayBench); i++ {
		homePlayer := homeBench[i]
		awayPlayer := awayBench[i]
		queue = append(queue, homePlayer)
		queue = append(queue, awayPlayer)
	}

	return queue
}

func RecordPlay(gs *GameState, eventID, outcomeID, injuryID, injuryType, injuryDuration, penaltyID uint8, pcID, ppID, apID, dpID, gpID uint, isBreakaway bool) {
	pc := gs.PuckCarrier
	zoneID := getZoneID(gs.PuckLocation, gs.HomeTeamID, gs.AwayTeamID)
	play := structs.PlayByPlay{
		Period:            gs.Period,
		TimeOnClock:       gs.TimeOnClock,
		SecondsConsumed:   uint8(gs.SecondsConsumed),
		EventID:           eventID,
		ZoneID:            uint8(zoneID),
		Outcome:           outcomeID,
		HomeTeamScore:     gs.HomeTeamScore,
		AwayTeamScore:     gs.AwayTeamScore,
		TeamID:            uint8(pc.TeamID),
		PuckCarrierID:     pcID,
		PassedPlayerID:    ppID,
		AssistingPlayerID: apID,
		DefenderID:        dpID,
		GoalieID:          gpID,
		InjuryID:          injuryID,
		InjuryType:        injuryType,
		InjuryDuration:    injuryDuration,
		PenaltyID:         penaltyID,
		IsBreakaway:       isBreakaway,
	}

	gs.RecordPlay(play)
}
