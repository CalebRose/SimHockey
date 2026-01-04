package engine

import (
	"fmt"
	"math"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/structs"
)

func HandleBaseEvents(gs *GameState) {
	// Check puck state first
	switch gs.PuckState {
	case PuckStateContested:
		// Puck battle is ongoing, resolve it
		handlePuckBattle(gs, gs.ContestedPlayers)
		return
	case PuckStateCovered:
		// Goalie has covered puck, handle faceoff setup
		gs.ClearPuckBattleState()
		return
	case PuckStateLoose:
		// Loose puck, try to establish possession
		// This should have been resolved by previous battle, continue normal flow
		gs.PuckState = PuckStateClear
	}

	pc := gs.PuckCarrier
	if pc == nil || pc.ID == 0 {
		// No clear puck carrier, this shouldn't happen but handle gracefully
		return
	}

	homePossession := pc.TeamID == uint16(gs.HomeTeamID)
	switch gs.PuckLocation {
	case HomeGoal:
		handleGoalZoneEvents(gs, homePossession, false)
		return
	case HomeZone, AwayZone:
		handleZoneEvents(gs, homePossession, gs.PuckLocation == HomeZone)
		return
	case NeutralZone:
		handleNeutralZoneEvents(gs)
		return
	case AwayGoal:
		handleGoalZoneEvents(gs, homePossession, true)
		return
	}
}

func handleGoalZoneEvents(gs *GameState, homePossession bool, isAwayGoal bool) {
	if (homePossession && isAwayGoal) || (!homePossession && !isAwayGoal) {
		handleOffensiveGoalZoneEvents(gs)
		return
	} else {
		handleDefensiveGoalZoneEvents(gs)
		return
	}
}

func handleZoneEvents(gs *GameState, homePossession bool, isHomeZone bool) {
	if (homePossession && isHomeZone) || (!homePossession && !isHomeZone) {
		handleDefensiveZoneEvents(gs)
		return
	} else {
		handleOffensiveZoneEvents(gs)
		return
	}
}

func handleOffensiveGoalZoneEvents(gs *GameState) {
	// Calculate base event weights
	eventWeights := CalculateEventWeights(gs)

	// Apply system modifiers
	ApplySystemModifiersToEventWeights(gs, &eventWeights)

	slapshot := eventWeights.ShotWeight
	pass := eventWeights.PassWeight
	passBack := eventWeights.PassBackWeight
	stickCheck := eventWeights.StickCheckWeight
	bodyCheck := eventWeights.BodyCheckWeight
	penalty := 1
	totalSkill := slapshot + stickCheck + bodyCheck + pass + passBack + penalty
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := float64(stickCheckCutoff) + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	passBackCheckCutoff := passCheckCutoff + float64(passBack)
	shotCutoff := passBackCheckCutoff + float64(slapshot)
	penaltyCutoff := shotCutoff + 0.1
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
		return
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
		return
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs, false, false)
		return
	} else if dr <= passBackCheckCutoff {
		handlePassCheck(gs, false, true)
		return
	} else if dr <= shotCutoff {
		HandleShot(gs, true)
		return
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
		return
	}
}

func handleOffensiveZoneEvents(gs *GameState) {
	// Calculate base event weights
	eventWeights := CalculateEventWeights(gs)

	// Apply system modifiers
	ApplySystemModifiersToEventWeights(gs, &eventWeights)

	wristshot := eventWeights.ShotWeight
	agility := eventWeights.AgilityWeight
	pass := eventWeights.PassWeight
	longPass := eventWeights.LongPassWeight
	stickCheck := eventWeights.StickCheckWeight
	bodyCheck := eventWeights.BodyCheckWeight
	penalty := 1
	totalSkill := wristshot + stickCheck + bodyCheck + pass + penalty + agility
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	longPassCheckCutoff := passCheckCutoff + float64(longPass)
	agilityCutoff := longPassCheckCutoff + float64(agility)
	wristshotCutoff := agilityCutoff + float64(wristshot)
	penaltyCutoff := wristshotCutoff + 0.1
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
		return
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
		return
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs, false, false)
		return
	} else if dr <= longPassCheckCutoff {
		handlePassCheck(gs, true, false)
		return
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
		return
	} else if dr <= wristshotCutoff {
		HandleShot(gs, false)
		return
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
		return
	}
}

func handleDefensiveGoalZoneEvents(gs *GameState) {
	// Calculate base event weights
	eventWeights := CalculateEventWeights(gs)

	// Apply system modifiers
	ApplySystemModifiersToEventWeights(gs, &eventWeights)

	pc := gs.PuckCarrier
	isHome := pc.TeamID == uint16(gs.HomeTeamID)
	penalty := 1
	agility := eventWeights.AgilityWeight
	pass := eventWeights.PassWeight
	longPass := eventWeights.LongPassWeight
	stickCheck := eventWeights.StickCheckWeight
	bodyCheck := eventWeights.BodyCheckWeight
	faceOffCheck := 0
	if pc.Position == Goalie {
		agility = 0
		stickCheck = 0
		bodyCheck = 0
		faceOffCheck = 20
	}
	totalSkill := stickCheck + bodyCheck + pass + penalty + agility + faceOffCheck
	faceoffCutoff := float64(faceOffCheck)
	stickCheckCutoff := faceoffCutoff + float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	longPassCheckCutoff := passCheckCutoff + float64(longPass)
	agilityCutoff := passCheckCutoff + float64(agility)
	penaltyCutoff := agilityCutoff + 0.1
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= faceoffCutoff {
		gs.SetFaceoffOnCenterIce(true)
		newZone := HomeZone
		if !isHome {
			newZone = AwayZone
		}
		_, zoneEnum := getZoneID(newZone, gs.HomeTeamID, gs.AwayTeamID)
		gs.SetNewZone(newZone)
		RecordPlay(gs, FaceoffID, GoalieHoldID, zoneEnum, 0, 0, 0, 0, 0, false, pc.ID, 0, 0, 0, pc.ID, false)
		HandleFaceoff(gs)
	} else if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
		return
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
		return
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs, false, false)
		return
	} else if dr <= longPassCheckCutoff {
		handlePassCheck(gs, true, false)
		return
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
		return
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
		return
	}
}

func handleDefensiveZoneEvents(gs *GameState) {
	// Calculate base event weights
	eventWeights := CalculateEventWeights(gs)

	// Apply system modifiers
	ApplySystemModifiersToEventWeights(gs, &eventWeights)

	penalty := 1
	agility := eventWeights.AgilityWeight
	pass := eventWeights.PassWeight
	passBack := eventWeights.PassBackWeight
	stickCheck := eventWeights.StickCheckWeight
	bodyCheck := eventWeights.BodyCheckWeight
	totalSkill := stickCheck + bodyCheck + pass + penalty + agility
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	passBackCheckCutoff := passCheckCutoff + float64(passBack)
	agilityCutoff := passCheckCutoff + float64(agility)
	penaltyCutoff := agilityCutoff + 0.1
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
		return
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
		return
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs, false, false)
		return
	} else if dr <= passBackCheckCutoff {
		handlePassCheck(gs, false, true)
		return
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
		return
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
		return
	}
}

func handleNeutralZoneEvents(gs *GameState) {
	// Calculate base event weights
	eventWeights := CalculateEventWeights(gs)

	// Apply system modifiers
	ApplySystemModifiersToEventWeights(gs, &eventWeights)

	agility := eventWeights.AgilityWeight
	pass := eventWeights.PassWeight
	stickCheck := eventWeights.StickCheckWeight
	bodyCheck := eventWeights.BodyCheckWeight
	penalty := 1
	totalSkill := stickCheck + bodyCheck + pass + penalty + agility
	stickCheckCutoff := float64(stickCheck)
	bodyCheckCutoff := stickCheckCutoff + float64(bodyCheck)
	passCheckCutoff := bodyCheckCutoff + float64(pass)
	agilityCutoff := passCheckCutoff + float64(agility)
	penaltyCutoff := agilityCutoff + 0.1
	dr := util.GenerateFloatFromRange(1, float64(totalSkill))
	if dr <= stickCheckCutoff {
		handleDefenseCheck(gs, true)
		return
	} else if dr <= bodyCheckCutoff {
		handleDefenseCheck(gs, false)
		return
	} else if dr <= passCheckCutoff {
		handlePassCheck(gs, false, false)
		return
	} else if dr <= agilityCutoff {
		handleAgilityCheck(gs)
		return
	} else if dr <= penaltyCutoff {
		handlePenalties(gs)
		return
	}
}

// handleDefenseCheck -- For preventing the defense from getting the puck
func handleDefenseCheck(gs *GameState, isStickCheck bool) {
	pc := gs.PuckCarrier
	// Select player on defense
	defendingTeamID := getDefendingTeamID(uint(pc.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	defender := selectDefendingPlayer(gs, defendingTeamID)
	if defender == nil {
		fmt.Println("ERROR: Could not find defending player for defense check")
		return
	}
	eventID := DexDefenseCheckID
	if !isStickCheck {
		eventID = PhysDefenseCheckID
	}
	secondsConsumed := util.GenerateIntFromRange(2, 4)
	gs.SetSecondsConsumed(uint16(secondsConsumed))
	chance := CalculatePenaltyChance()
	if chance {
		shouldReturn := handlePenalty(gs, !isStickCheck, defender, eventID, pc.ID)
		if shouldReturn {
			return
		}
	}
	diceRoll := util.GenerateIntFromRange(1, 20)
	outcomeID := CarrierKeepsPuckID

	// Critical checks
	switch diceRoll {
	case CritFail:
		// Defender gets puck
		outcomeID = DefenseTakesPuckID
		RecordPlay(gs, eventID, outcomeID, 0, 0, 0, 0, 0, 0, false, pc.ID, 0, 0, defender.ID, 0, false)
		defender.AddDefensiveHit(!isStickCheck)
		gs.SetPuckBearer(defender, false)
		return
	case CritSuccess:
		// Defense DOES NOT get puck
		RecordPlay(gs, eventID, outcomeID, 0, 0, 0, 0, 0, 0, false, pc.ID, 0, 0, defender.ID, 0, false)
		return
	default:
		// Determine if physical check or non-physical check
		puckHandling := DiffReq + pc.HandlingMod
		eventType := StickCheckEvent

		if isStickCheck {
			puckHandling -= defender.StickCheckMod
		} else {
			eventType = BodyCheckEvent
			puckHandling -= defender.BodyCheckMod
		}

		// Check puck carrier (receiver) for injury - higher risk
		pcInjuryChance := CalculateInjuryRisk(eventType, int(pc.InjuryRating), util.PCDefenderIntensity) // 20% higher risk for receiver
		pcInjuryOccurs := IsPlayerInjured(pcInjuryChance)
		if pcInjuryOccurs {
			HandleInjuryEvent(gs, eventType, pc)
		}

		// Check defender (deliverer) for injury - lower risk
		defenderInjuryChance := CalculateInjuryRisk(eventType, int(defender.InjuryRating), util.DefenderIntensity) // 40% lower risk for deliverer
		defenderInjuryOccurs := IsPlayerInjured(defenderInjuryChance)
		if defenderInjuryOccurs {
			HandleInjuryEvent(gs, eventType, defender)
		}

		if pc.IsInjured || defender.IsInjured {
			switch gs.PuckLocation {
			case HomeGoal:
				gs.SetNewZone(HomeZone)
			case AwayGoal:
				gs.SetNewZone(AwayZone)
			}
			HandleFaceoff(gs)
			return
		}

		puckHandling = math.Max(puckHandling, 1.0)

		// Check if this should trigger a puck battle
		attackWeight := pc.HandlingMod + pc.StrengthMod
		defenseWeight := defender.StickCheckMod + defender.BodyCheckMod
		if !isStickCheck {
			defenseWeight = defender.BodyCheckMod + defender.StrengthMod
		}

		shouldBattle := shouldTriggerPuckBattle(gs, attackWeight, defenseWeight)

		if shouldBattle && diceRoll > 5 && diceRoll < 15 {
			// Trigger puck battle instead of clean turnover
			triggerPuckBattle(gs, pc, defender)
			return
		}

		// Normal resolution
		if float64(diceRoll) >= puckHandling {
			defender.AddDefensiveHit(!isStickCheck)
			outcomeID = DefenseTakesPuckID
			gs.SetPuckBearer(defender, false)
		}
		RecordPlay(gs, eventID, outcomeID, 0, 0, 0, 0, 0, 0, false, pc.ID, 0, 0, defender.ID, 0, false)
	}
}

func handleAgilityCheck(gs *GameState) {
	// Get Current Zone
	nextZone := getNextZone(gs)

	pc := gs.PuckCarrier
	agilityMod := pc.AgilityMod
	momentumMod := gs.Momentum
	critCheck := util.GenerateIntFromRange(1, 20)
	secondsConsumed := util.GenerateIntFromRange(1, 4)
	defenseCheck := true
	isBreakaway := false
	if critCheck == CritFail {
		secondsConsumed += 3
	} else if critCheck == CritSuccess || float64(critCheck) > EasyReq+agilityMod+momentumMod {
		defenseCheck = false
	}
	if critCheck == CritSuccess {
		isBreakaway = true
		gs.TriggerBreakaway()
	}
	gs.SetSecondsConsumed(uint16(secondsConsumed))
	eventId := AgilityCheckID
	_, nextZoneEnum := getZoneID(nextZone, gs.HomeTeamID, gs.AwayTeamID)
	if defenseCheck {
		defendingTeamID := getDefendingTeamID(uint(pc.TeamID), gs.HomeTeamID, gs.AwayTeamID)
		defender := selectDefendingPlayer(gs, defendingTeamID)
		if defender == nil {
			fmt.Println("ERROR: Could not find defending player for agility check")
			// Continue with the agility check without defense interference
		} else {
			diceRoll := util.GenerateIntFromRange(1, 20)
			puckHandling := ToughReq + 1 + pc.HandlingMod
			// if !gs.IsCollegeGame {
			// 	puckHandling = ToughReq + pb.HandlingMod
			// }
			coinFlip := util.CoinFlip()
			eventType := StickCheckEvent
			if coinFlip == Heads {
				puckHandling -= defender.BodyCheckMod
				eventType = BodyCheckEvent
			} else {
				puckHandling -= defender.StickCheckMod
			}

			chance := CalculatePenaltyChance()
			if chance {
				shouldReturn := handlePenalty(gs, coinFlip == Heads, defender, eventId, pc.ID)
				if shouldReturn {
					return
				}
			}

			pcInjuryChance := CalculateInjuryRisk(eventType, int(pc.InjuryRating), util.PCDefenderIntensity) // 20% higher risk for receiver
			if IsPlayerInjured(pcInjuryChance) {
				HandleInjuryEvent(gs, eventType, pc)
			}

			// Check defender (deliverer) for injury - lower risk
			defenderInjuryChance := CalculateInjuryRisk(eventType, int(defender.InjuryRating), util.DefenderIntensity) // 40% lower risk for deliverer
			if IsPlayerInjured(defenderInjuryChance) {
				HandleInjuryEvent(gs, eventType, defender)
			}

			if pc.IsInjured || defender.IsInjured {
				switch gs.PuckLocation {
				case HomeGoal:
					gs.SetNewZone(HomeZone)
				case AwayGoal:
					gs.SetNewZone(AwayZone)
				}
				HandleFaceoff(gs)
				return
			}

			puckHandling = math.Max(puckHandling, 1.0)

			if float64(diceRoll) < puckHandling {
				RecordPlay(gs, eventId, OffenseMovesUpID, nextZoneEnum, 0, 0, 0, 0, 0, false, pc.ID, 0, 0, defender.ID, 0, isBreakaway)
				gs.SetNewZone(nextZone)
				return
			} else {
				RecordPlay(gs, eventId, DefenseStopAgilityID, 0, 0, 0, 0, 0, 0, false, pc.ID, 0, 0, defender.ID, 0, false)
				defender.AddDefensiveHit(coinFlip == Heads)
				gs.SetPuckBearer(defender, false)
				// Logger(defender.FirstName + " GETS THE PUCK FOR " + defender.Team + "!")
				return
			}
		}

	}

	// Logger(pb.FirstName + " " + pb.LastName + " moves up to " + nextZone + "!")
	// Move up zone
	RecordPlay(gs, eventId, OffenseMovesUpID, uint8(nextZoneEnum), 0, 0, 0, 0, 0, false, pc.ID, 0, 0, 0, 0, isBreakaway)
	gs.SetNewZone(nextZone)
}

func handlePenalty(gs *GameState, isBodyCheck bool, defender *GamePlayer, eventID uint8, pcId uint) bool {
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
		sevId := GetSeverityID(penalty.Severity)
		RecordPlay(gs, eventID, uint8(penalty.PenaltyID), 0, 0, 0, 0, uint8(penalty.PenaltyID), sevId, penalty.IsFight, pcId, 0, 0, defender.ID, 0, false)
		ApplyPenalty(gs, penalty, defender)
		return true
	}
	return false
}

func handlePassCheck(gs *GameState, longPass, backPass bool) {
	pb := gs.PuckCarrier

	// Roll to see if puck is intercepted by defense
	defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	defender := selectDefendingPlayer(gs, defendingTeamID)

	// If cannot find a defending player
	if defender == nil {
		fmt.Println("ERROR: Could not find defending player")
		// Fall back to no interception attempt
		playerList := getFullPlayerListByTeamID(uint(pb.TeamID), gs)
		filteredList := getAvailablePlayers(pb.ID, playerList)
		receivingPlayer := PassPuckToPlayer(gs, filteredList, gs.PuckLocation)
		if receivingPlayer == 0 {
			// No available player, keep puck
			secondsConsumed := util.GenerateIntFromRange(1, 3)
			gs.SetSecondsConsumed(uint16(secondsConsumed))
			return
		}
		retrievingPlayer, _ := findPlayerByID(playerList, receivingPlayer)
		if retrievingPlayer == nil {
			// Still can't find player, keep puck
			secondsConsumed := util.GenerateIntFromRange(1, 3)
			gs.SetSecondsConsumed(uint16(secondsConsumed))
			return
		}
		gs.SetPuckBearer(retrievingPlayer, longPass || backPass)
		secondsConsumed := util.GenerateIntFromRange(1, 3)
		gs.SetSecondsConsumed(uint16(secondsConsumed))
		return
	}
	passID := PassCheckID
	if longPass {
		passID = LongPassCheckID
	} else if backPass {
		passID = PassBackCheckID
	}
	safePass := CalculateSafePass(pb.PassMod, defender.StickCheckMod, longPass || backPass)

	if !safePass {

		defenderInjuryChance := CalculateInjuryRisk(MissedPassInterception, int(defender.InjuryRating), util.DefenderIntensity) // 40% lower risk for deliverer
		if IsPlayerInjured(defenderInjuryChance) {
			HandleInjuryEvent(gs, MissedPassInterception, defender)
		}

		if defender.IsInjured {
			switch gs.PuckLocation {
			case HomeGoal:
				gs.SetNewZone(HomeZone)
			case AwayGoal:
				gs.SetNewZone(AwayZone)
			}
			HandleFaceoff(gs)
			return
		}

		gs.SetPuckBearer(defender, false)
		RecordPlay(gs, passID, InterceptedPassID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, 0, defender.ID, 0, false)
		// Logger(defender.FirstName + " INTERCEPTS THE PASS FOR " + defender.Team + "!")
		return
	}

	// Get available player on own team
	playerList := getFullPlayerListByTeamID(uint(pb.TeamID), gs)
	filteredList := getAvailablePlayers(pb.ID, playerList)

	receivingPlayer := PassPuckToPlayer(gs, filteredList, gs.PuckLocation)
	if receivingPlayer == 0 {
		fmt.Println("Cannot find open player")
		// No available player to pass to, hold onto puck
		secondsConsumed := util.GenerateIntFromRange(1, 3)
		gs.SetSecondsConsumed(uint16(secondsConsumed))
		RecordPlay(gs, passID, NoOneOpenID, 0, 0, 0, 0, 0, 0, false, pb.ID, receivingPlayer, 0, defender.ID, 0, false)
		return
	}
	retrievingPlayer, _ := findPlayerByID(playerList, receivingPlayer)
	if retrievingPlayer == nil {
		fmt.Printf("ERROR: Could not find receiving player with ID %d\n", receivingPlayer)
		// Fall back to keeping current puck carrier
		secondsConsumed := util.GenerateIntFromRange(1, 3)
		gs.SetSecondsConsumed(uint16(secondsConsumed))
		RecordPlay(gs, passID, NoOneOpenID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, 0, defender.ID, 0, false)
		return
	}
	HandleMissingPlayer(*retrievingPlayer, "PASSING PUCK")
	enum := uint8(0)
	zone := ""
	if longPass {
		zone = getNextZone(gs)
		_, nextZoneEnum := getZoneID(zone, gs.HomeTeamID, gs.AwayTeamID)
		enum = nextZoneEnum
	} else if backPass {
		zone = getPreviousZone(gs)
		_, nextZoneEnum := getZoneID(zone, gs.HomeTeamID, gs.AwayTeamID)
		enum = nextZoneEnum
	}
	// If a long pass or back pass, move to next zone
	if len(zone) > 0 {
		gs.SetNewZone(zone)
	}
	secondsConsumed := util.GenerateIntFromRange(1, 3)
	gs.SetSecondsConsumed(uint16(secondsConsumed))
	RecordPlay(gs, passID, ReceivedPassID, enum, 0, 0, 0, 0, 0, false, pb.ID, receivingPlayer, 0, defender.ID, 0, false)
	gs.SetPuckBearer(retrievingPlayer, longPass || backPass)
}

func handlePenalties(gs *GameState) {
	// Determine penalty type
	// Minor, major, misconduct, game misconduct, match
	pc := gs.PuckCarrier
	defendingTeamID := getDefendingTeamID(uint(pc.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	zoneID := 0
	switch gs.PuckLocation {
	case HomeGoal, AwayGoal:
		zoneID = 2
	case HomeZone, AwayZone:
		zoneID = 1
	}

	player := selectDefendingPlayer(gs, defendingTeamID)
	if player == nil {
		fmt.Println("ERROR: Could not find player for penalty")
		return
	}
	secondPlayer := &GamePlayer{}
	penaltyTypeID := GeneralPenaltyID
	penaltyType := General
	diceRoll := util.DiceRoll(0, 20)
	if diceRoll {
		penaltyTypeID = FightPenaltyID
		penaltyType = Fight
		// If a fight occurs, then two players should probably get placed in the penalty box.
		secondPlayer = selectDefendingPlayer(gs, uint(pc.TeamID))
		if secondPlayer == nil {
			// If we can't find a second player, treat it as a regular penalty
			secondPlayer = &GamePlayer{}
			penaltyTypeID = GeneralPenaltyID
			penaltyType = General
		}
	}

	penalty := SelectPenalty(player, uint(zoneID), penaltyType)
	if penalty.PenaltyID == 0 {
		return
	}
	sevId := GetSeverityID(penalty.Severity)
	RecordPlay(gs, PenaltyCheckID, penaltyTypeID, 0, 0, 0, 0, uint8(penalty.PenaltyID), sevId, penalty.IsFight, pc.ID, 0, 0, player.ID, secondPlayer.ID, false)

	// Roll for injury if major
	if penalty.Severity == MajorPenalty || penalty.Severity == MatchPenalty {
		intensity := util.SevereInjuryIntensity
		if penalty.Severity == MatchPenalty {
			intensity = util.CriticalInjuryIntensity
		}
		injuryChance := CalculateInjuryRisk(PenaltyEvent, int(pc.InjuryRating), intensity)
		if IsPlayerInjured(injuryChance) {
			HandleInjuryEvent(gs, PenaltyEvent, pc)
		}
	}

	// Apply Penalty to Player
	ApplyPenalty(gs, penalty, player)
	if player.ID != secondPlayer.ID {
		ApplyPenalty(gs, penalty, secondPlayer)
	}
}

func HandleFaceoff(gs *GameState) {
	// Get Centers from current lineups
	homeCenter := gs.GetCenter(true)
	awayCenter := gs.GetCenter(false)
	if homeCenter == nil || awayCenter == nil {
		fmt.Printf("ERROR: Missing centers for faceoff - Home: %v, Away: %v\n", homeCenter != nil, awayCenter != nil)
		return
	}
	HandleMissingPlayer(*homeCenter, "HandleFaceoff Home Center")
	HandleMissingPlayer(*awayCenter, "HandleFaceoff Away Center")
	homeFaceoffWin := CalculateFaceoff(homeCenter.FaceoffMod, awayCenter.FaceoffMod)
	faceOffWinID := homeCenter.TeamID
	// Away wins faceoff
	if !homeFaceoffWin {
		faceOffWinID = awayCenter.TeamID
		awayCenter.AddFaceoff(true)
		homeCenter.AddFaceoff(false)
		gs.AwayTeamStats.AddFaceoff(true)
		gs.HomeTeamStats.AddFaceoff(false)
	} else {
		homeCenter.AddFaceoff(true)
		awayCenter.AddFaceoff(false)
		gs.HomeTeamStats.AddFaceoff(true)
		gs.AwayTeamStats.AddFaceoff(false)
	}
	if gs.FaceoffOnCenterIce {
		gs.SetFaceoffOnCenterIce(false)
	}
	// Select player who gets puck after faceoff
	HandleFaceoffRetrieval(gs, homeFaceoffWin, uint(faceOffWinID), homeCenter.ID, awayCenter.ID)
}

func HandleFaceoffRetrieval(gs *GameState, homeFaceoffWin bool, faceoffWinID, homeCenterID, awayCenterID uint) {
	puckLocation := NeutralZone
	playerList := []*GamePlayer{}
	// Get Available Players in Home Forward Line
	hgs := gs.HomeStrategy
	ags := gs.AwayStrategy
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(true).ID, hgs.Forwards[hgs.CurrentForwards].Players)...)
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(true).ID, hgs.Defenders[hgs.CurrentDefenders].Players)...)
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(false).ID, ags.Forwards[ags.CurrentForwards].Players)...)
	playerList = append(playerList, getAvailablePlayers(gs.GetCenter(false).ID, ags.Defenders[ags.CurrentDefenders].Players)...)

	faceoffRetrievalCheck := RetrievePuckAfterFaceoffCheck(playerList, puckLocation, faceoffWinID, homeFaceoffWin)
	retrievingPlayer, _ := findPlayerByID(playerList, faceoffRetrievalCheck)
	if retrievingPlayer == nil {
		fmt.Printf("ERROR: Could not find retrieving player after faceoff with ID %d\n", faceoffRetrievalCheck)
		return
	}
	HandleMissingPlayer(*retrievingPlayer, "REBOUNDING AFTER FACEOFF")
	gs.SetPuckBearer(retrievingPlayer, false)
	outcomeID := HomeFaceoffWinID
	if !homeFaceoffWin {
		outcomeID = AwayFaceoffWinID
	}
	RecordPlay(gs, FaceoffID, outcomeID, 0, 0, 0, 0, 0, 0, false, homeCenterID, awayCenterID, retrievingPlayer.ID, 0, 0, false)
	// Logger(retrievingPlayer.Team + " gets the puck from the faceoff with " + retrievingPlayer.FirstName + " " + retrievingPlayer.LastName + " in possession!")
}

// getSystemAwareReboundPlayers - Build player list based on offensive systems and zone
func getSystemAwareReboundPlayers(gs *GameState, reboundZone string) []*GamePlayer {
	reboundPlayerList := []*GamePlayer{}
	hgs := gs.HomeStrategy
	ags := gs.AwayStrategy

	pc := gs.PuckCarrier
	isHomePossession := pc.TeamID == uint16(gs.HomeTeamID)

	// Get current offensive system for the team with possession
	var offensiveSystem uint8
	var intensity uint8
	if isHomePossession {
		offensiveSystem = hgs.Gameplan.OffensiveSystem
		intensity = hgs.Gameplan.OffensiveIntensity
	} else {
		offensiveSystem = ags.Gameplan.OffensiveSystem
		intensity = ags.Gameplan.OffensiveIntensity
	}

	switch reboundZone {
	case HomeGoal:
		// Goal area - depends on offensive system philosophy
		if isHomePossession {
			// Home team shooting at away goal - they're attacking
			reboundPlayerList = append(reboundPlayerList, getOffensiveReboundPlayers(hgs, offensiveSystem, intensity, true)...)
			reboundPlayerList = append(reboundPlayerList, getDefensiveReboundPlayers(ags, false)...)
		} else {
			// Away team shooting at home goal - home team defending
			reboundPlayerList = append(reboundPlayerList, getDefensiveReboundPlayers(hgs, true)...)
			reboundPlayerList = append(reboundPlayerList, getOffensiveReboundPlayers(ags, offensiveSystem, intensity, false)...)
		}

	case AwayGoal:
		// Goal area - depends on offensive system philosophy
		if isHomePossession {
			// Home team defending their goal
			reboundPlayerList = append(reboundPlayerList, getDefensiveReboundPlayers(hgs, true)...)
			reboundPlayerList = append(reboundPlayerList, getOffensiveReboundPlayers(ags, offensiveSystem, intensity, false)...)
		} else {
			// Away team shooting at home goal - they're attacking
			reboundPlayerList = append(reboundPlayerList, getOffensiveReboundPlayers(ags, offensiveSystem, intensity, false)...)
			reboundPlayerList = append(reboundPlayerList, getDefensiveReboundPlayers(hgs, true)...)
		}

	case HomeZone, AwayZone:
		// Zone play - all players in current lines involved
		reboundPlayerList = append(reboundPlayerList, hgs.Forwards[hgs.CurrentForwards].Players...)
		reboundPlayerList = append(reboundPlayerList, ags.Defenders[ags.CurrentDefenders].Players...)
		reboundPlayerList = append(reboundPlayerList, ags.Forwards[ags.CurrentForwards].Players...)
		reboundPlayerList = append(reboundPlayerList, hgs.Defenders[hgs.CurrentDefenders].Players...)
	}

	return reboundPlayerList
}

// getOffensiveReboundPlayers - Get players positioned for offensive rebounds based on system
func getOffensiveReboundPlayers(strategy GamePlaybook, system, intensity uint8, isHome bool) []*GamePlayer {
	players := []*GamePlayer{}

	// Base players - always include current forward line
	forwards := strategy.Forwards[strategy.CurrentForwards].Players
	defensemen := strategy.Defenders[strategy.CurrentDefenders].Players

	switch system {
	case 8: // Crash the Net - maximum net presence
		// All forwards crash, some defensemen join
		players = append(players, forwards...)
		if intensity >= 6 { // High intensity - defensemen also crash
			players = append(players, defensemen...)
		} else {
			// Only add one defenseman for moderate crashing
			if len(defensemen) > 0 {
				players = append(players, defensemen[0])
			}
		}

	case 4: // Cycle Game - patient possession, spread positioning
		// Forwards cycle, defensemen stay back unless high intensity
		players = append(players, forwards...)
		if intensity >= 7 { // Very high intensity cycling
			players = append(players, defensemen...)
		}

	case 6: // Umbrella - structured positioning with D-man quarterback
		// All forwards, plus the quarterback defenseman
		players = append(players, forwards...)
		if len(defensemen) > 0 {
			players = append(players, defensemen[0]) // Quarterback D
		}

	case 1, 2, 3: // Forechecking systems - forwards press, D support varies
		players = append(players, forwards...)
		if system == 3 && intensity >= 6 { // 1-1-3 with high intensity - more forwards up
			// Already have all forwards
		}
		if system == 2 && intensity >= 7 { // 2-1-2 with very high intensity
			if len(defensemen) > 0 {
				players = append(players, defensemen[0]) // One D joins attack
			}
		}

	case 5: // Quick Transition - speed over positioning
		players = append(players, forwards...)
		// Offensive defensemen more likely to be up ice
		players = append(players, defensemen...)

	case 7: // East-West Motion - all forwards, selective D involvement
		players = append(players, forwards...)
		if intensity >= 6 {
			if len(defensemen) > 0 {
				players = append(players, defensemen[0])
			}
		}

	default: // Balanced/Unknown system
		players = append(players, forwards...)
		if intensity >= 6 {
			if len(defensemen) > 0 {
				players = append(players, defensemen[0])
			}
		}
	}

	return players
}

// getDefensiveReboundPlayers - Get players positioned for defensive rebounds
func getDefensiveReboundPlayers(strategy GamePlaybook, isHome bool) []*GamePlayer {
	players := []*GamePlayer{}

	// Defensive positioning - typically defensemen plus some forwards
	defensemen := strategy.Defenders[strategy.CurrentDefenders].Players
	forwards := strategy.Forwards[strategy.CurrentForwards].Players

	// Always include defensemen for defensive rebounds
	players = append(players, defensemen...)

	// Add defensive-minded forwards (typically checking line players)
	// For simplicity, add the current forward line but they'll be weighted lower
	players = append(players, forwards...)

	return players
}

func HandleReboundAfterShot(gs *GameState, eventID uint8, outcomeID uint8, puckCarrierID uint, goalieID uint) {
	puckLocation := GetPuckLocationAfterMiss(1, 1)
	reboundPlayerList := getSystemAwareReboundPlayers(gs, puckLocation)

	// Check if this should be a puck battle scenario
	// Rebounds are prime candidates for battles, especially in goal areas
	shouldBeBattle := false

	switch gs.PuckLocation {
	case HomeGoal, AwayGoal:
		// High chance of battle in goal areas
		shouldBeBattle = util.GenerateIntFromRange(1, 100) <= 70 // 70% chance
	case HomeZone, AwayZone:
		// Medium chance in regular zones
		shouldBeBattle = util.GenerateIntFromRange(1, 100) <= 40 // 40% chance
	}

	if shouldBeBattle && len(reboundPlayerList) >= 2 {
		// Multiple players converge on loose puck - trigger battle
		gs.PuckState = PuckStateLoose
		handlePuckBattle(gs, reboundPlayerList)
		// Record the rebound event but note it became contested
		puckCarrierIDForRecord := uint(0)
		if gs.PuckCarrier != nil {
			puckCarrierIDForRecord = gs.PuckCarrier.ID
		}
		RecordPlay(gs, eventID, PuckScrambleID, 0, 0, 0, 0, 0, 0, false, puckCarrierID, puckCarrierIDForRecord, 0, 0, goalieID, false)
		return
	}

	// Normal rebound resolution
	reboundCheck := reboundCheck(gs, reboundPlayerList, puckLocation)
	reboundingPlayer, _ := findPlayerByID(reboundPlayerList, reboundCheck)
	if reboundingPlayer == nil {
		fmt.Printf("ERROR: Could not find rebounding player with ID %d\n", reboundCheck)
		return
	}
	HandleMissingPlayer(*reboundingPlayer, "REBOUNDING AFTER SHOT")
	// Record Play after inaccurate shot
	RecordPlay(gs, eventID, outcomeID, 0, 0, 0, 0, 0, 0, false, puckCarrierID, reboundingPlayer.ID, 0, 0, goalieID, false)
	gs.SetPuckBearer(reboundingPlayer, false)
	// Logger(reboundingPlayer.FirstName + " " + reboundingPlayer.LastName + " gets the puck for " + reboundingPlayer.Team + " on the rebound!")
}

func HandleShot(gs *GameState, isCloseShot bool) {
	pb := gs.PuckCarrier
	if pb == nil {
		fmt.Println("ERROR: No puck carrier for shot attempt")
		return
	}
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	goalie := &GamePlayer{}
	goalieStrategy := gs.GetLineStrategy(!isHome, 3)
	goalie = goalieStrategy.Players[0]
	accuracy := 0
	power := 0
	goalieMod := 0.0
	goalKeepingMod := 0.0
	eventTypeID := SlapshotCheckID
	if isCloseShot {
		accuracy = int(pb.CloseShotAccMod)
		power = int(pb.CloseShotPowerMod)
		goalieMod = goalie.StrengthMod
		goalKeepingMod = goalie.GoalkeepingMod
	} else {
		eventTypeID = WristshotCheckID
		accuracy = int(pb.LongShotAccMod)
		power = int(pb.LongShotPowerMod)
		goalieMod = goalie.AgilityMod
		goalKeepingMod = float64(goalie.GoalieVision)
	}

	// Seconds Consumed
	secondsConsumed := util.GenerateIntFromRange(2, 5)
	gs.SetSecondsConsumed(uint16(secondsConsumed))

	// Detect shotblocking
	defendingTeamID := getDefendingTeamID(uint(pb.TeamID), gs.HomeTeamID, gs.AwayTeamID)
	defender := selectBlockingPlayer(gs, defendingTeamID)
	if defender == nil {
		fmt.Println("ERROR: Could not find blocking player")
		// Continue with shot attempt without blocking
	} else {
		shotBlocked := CalculateShotBlock(defender.ShotblockingMod)
		if shotBlocked {
			// Check defender (deliverer) for injury - lower risk
			defenderInjuryChance := CalculateInjuryRisk(MissedShotBlocked, int(defender.InjuryRating), util.SevereInjuryIntensity) // 40% lower risk for deliverer
			if IsPlayerInjured(defenderInjuryChance) {
				HandleInjuryEvent(gs, MissedShotBlocked, defender)
			}

			if defender.IsInjured {
				switch gs.PuckLocation {
				case HomeGoal:
					gs.SetNewZone(HomeZone)
				case AwayGoal:
					gs.SetNewZone(AwayZone)
				}
				HandleFaceoff(gs)
				return
			}

			pb.AddShot(false, false, false, false, gs.Period > 3)
			defender.AddShotBlocked()
			RecordPlay(gs, eventTypeID, ShotBlockedID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, 0, defender.ID, goalie.ID, false)
			gs.SetPuckBearer(defender, false)
			return
		}
	}
	accuracyCheck := CalculateAccuracy(float64(accuracy), isCloseShot)
	if !accuracyCheck {
		// Might need to record who gets the rebound...
		HandleReboundAfterShot(gs, eventTypeID, InAccurateShotID, pb.ID, goalie.ID)
		pb.AddShot(false, false, false, false, false)
		return
	}

	baseCheck := CollegeBaseShot
	if !gs.IsCollegeGame {
		baseCheck = ProBaseShot
	}
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
		RecordPlay(gs, eventTypeID, ShotOnGoalID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, gs.AssistingPlayer.ID, 0, goalie.ID, false)
		gs.IncrementScore(isHome)
		goalie.AddShotAgainst(true)
	} else {
		gs.AddShots(isHome)
		goalie.AddShotAgainst(false)
		pb.AddShot(false, false, false, false, false)

		goalieInjuryChance := CalculateInjuryRisk(PuckContact, int(goalie.InjuryRating), util.BaseInjuryIntensity) // 40% lower risk for deliverer
		if IsPlayerInjured(goalieInjuryChance) {
			HandleInjuryEvent(gs, PuckContact, goalie)
		}

		if defender.IsInjured {
			switch gs.PuckLocation {
			case HomeGoal:
				gs.SetNewZone(HomeZone)
			case AwayGoal:
				gs.SetNewZone(AwayZone)
			}
			HandleFaceoff(gs)
			return
		}

		switch gs.PuckLocation {
		case HomeZone:
			gs.SetNewZone(HomeGoal)
		case AwayZone:
			gs.SetNewZone(AwayGoal)
		}
		RecordPlay(gs, eventTypeID, GoalieSaveID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, gs.AssistingPlayer.ID, 0, goalie.ID, false)
		gs.SetPuckBearer(goalie, false)
	}
}

func HandleShootoutAttempt(gs *GameState) {
	pb := gs.PuckCarrier
	if pb == nil {
		fmt.Println("ERROR: No puck carrier for shootout attempt")
		return
	}
	isHome := pb.TeamID == uint16(gs.HomeTeamID)
	goalie := &GamePlayer{}
	goalieStrategy := gs.GetLineStrategy(!isHome, 3)
	goalie = goalieStrategy.Players[0]
	slapPower := pb.CloseShotAccuracy + pb.CloseShotPower
	wsPower := pb.LongShotPower + pb.LongShotAccuracy
	accuracy := int(pb.CloseShotAccMod)
	power := int(pb.CloseShotPowerMod)
	isSlapshot := true
	eventID := CSShootoutID
	if wsPower > slapPower {
		eventID = WSShootoutID
		isSlapshot = false
		accuracy = int(pb.LongShotAccMod)
		power = int(pb.LongShotPowerMod)
	}

	accuracyCheck := CalculateAccuracy(float64(accuracy), isSlapshot)
	if !accuracyCheck {
		RecordPlay(gs, eventID, InAccurateShotID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, 0, 0, goalie.ID, false)
		return
	}

	goalieMod := goalie.StrengthMod
	goalKeepingMod := goalie.GoalkeepingMod
	if !isSlapshot {
		goalieMod = goalie.AgilityMod
		goalKeepingMod = goalie.GoalieVisionMod
	}

	baseCheck := 12.0
	shotAttempt := CalculateShot(float64(power),
		float64(gs.PuckCarrier.OneTimerMod)*ShootoutMomenumModifier,
		goalKeepingMod,
		goalieMod, baseCheck)

	if shotAttempt {
		gs.IncrementShootoutScore(isHome)
		RecordPlay(gs, eventID, ShotOnGoalID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, 0, 0, goalie.ID, false)

	} else {
		RecordPlay(gs, eventID, GoalieSaveID, 0, 0, 0, 0, 0, 0, false, pb.ID, 0, 0, 0, goalie.ID, false)
	}
}

func HandleOvertimeShootout(gs *GameState) {
	// Shootouts should be close shots
	isRepeat := false
	// Organize shootout queue. Home, then away, repeat (H,A,H,A,H,A,H,A...)
	shootoutQueue := formShootoutQueue(gs.HomeStrategy, gs.AwayStrategy)
	// Run For Loop for Shootout. Once six players from each team have made an attempt, compare the shootout scores
	// If score is still the same, keep running through loop

	RecordPlay(gs, EnteringShootout, 0, 0, 0, 0, 0, 0, 0, false, 0, 0, 0, 0, 0, false)

	for gs.HomeTeamShootoutScore == gs.AwayTeamShootoutScore {
		for idx, player := range shootoutQueue {
			// If we go beyond the first six shots and the score is still the same, continue
			if (idx > 5 && gs.HomeTeamShootoutScore != gs.AwayTeamShootoutScore && !isRepeat) || (isRepeat && gs.HomeTeamShootoutScore != gs.AwayTeamShootoutScore) {
				break
			}
			gs.SetPuckBearer(player, false)
			HandleShootoutAttempt(gs)
		}
	}
}

func formShootoutQueue(homeGP, awayGP GamePlaybook) []*GamePlayer {
	queue := []*GamePlayer{}

	homeLineup := getShootoutLineup(homeGP)
	awayLineup := getShootoutLineup(awayGP)

	// Loop through the same 6 shooters
	for i := 0; i < len(homeLineup); i++ {
		homePlayer := homeLineup[i]
		if homePlayer == nil {
			fmt.Printf("ERROR: Missing home player in shootout lineup at index %d\n", i)
			continue
		}
		HandleMissingPlayer(*homePlayer, "SHOOTOUT QUEUE, INDEX "+strconv.Itoa(i)+" TeamID:"+strconv.Itoa(int(homeGP.ShootoutLineUp.TeamID)))
		awayPlayer := awayLineup[i]
		if awayPlayer == nil {
			fmt.Printf("ERROR: Missing away player in shootout lineup at index %d\n", i)
			continue
		}
		HandleMissingPlayer(*awayPlayer, "SHOOTOUT QUEUE, INDEX "+strconv.Itoa(i)+" TeamID:"+strconv.Itoa(int(awayGP.ShootoutLineUp.TeamID)))
		queue = append(queue, homePlayer)
		queue = append(queue, awayPlayer)
	}

	return queue
}

func getShootoutLineup(gp GamePlaybook) []*GamePlayer {
	lineupIDs := gp.ShootoutLineUp
	forwards := gp.Forwards
	defenders := gp.Defenders
	allPlayers := []*GamePlayer{}
	queue := []*GamePlayer{}
	s1 := &GamePlayer{}
	s2 := &GamePlayer{}
	s3 := &GamePlayer{}
	s4 := &GamePlayer{}
	s5 := &GamePlayer{}
	s6 := &GamePlayer{}
	for _, f := range forwards {
		for _, player := range f.Players {
			if player.IsOut {
				continue
			}
			allPlayers = append(allPlayers, player)
		}
	}
	for _, d := range defenders {
		for _, player := range d.Players {
			if player.IsOut {
				continue
			}
			allPlayers = append(allPlayers, player)
		}
	}

	selectedMap := make(map[uint]bool)

	for _, p := range allPlayers {
		switch p.ID {
		case lineupIDs.Shooter1ID:
			s1 = p
			selectedMap[p.ID] = true
		case lineupIDs.Shooter2ID:
			s2 = p
			selectedMap[p.ID] = true
		case lineupIDs.Shooter3ID:
			s3 = p
			selectedMap[p.ID] = true
		case lineupIDs.Shooter4ID:
			s4 = p
			selectedMap[p.ID] = true
		case lineupIDs.Shooter5ID:
			s5 = p
			selectedMap[p.ID] = true
		case lineupIDs.Shooter6ID:
			s6 = p
			selectedMap[p.ID] = true
		}
	}
	// If any of the positions are STILL empty, go through the loop of all players and fix
	missingPlayer := s1.ID == 0 || s2.ID == 0 || s3.ID == 0 || s4.ID == 0 || s5.ID == 0 || s6.ID == 0
	if missingPlayer {
		// Get Bench
		for _, player := range gp.BenchPlayers {
			if player.IsOut || player.Position == Goalie {
				continue
			}
			allPlayers = append(allPlayers, player)
		}
	}
	for missingPlayer {
		for _, p := range allPlayers {
			if selectedMap[p.ID] {
				continue
			}
			if s1.ID == 0 {
				selectedMap[p.ID] = true
				s1 = p
			} else if s2.ID == 0 {
				selectedMap[p.ID] = true
				s2 = p
			} else if s3.ID == 0 {
				selectedMap[p.ID] = true
				s3 = p
			} else if s4.ID == 0 {
				selectedMap[p.ID] = true
				s4 = p
			} else if s5.ID == 0 {
				selectedMap[p.ID] = true
				s5 = p
			} else if s6.ID == 0 {
				selectedMap[p.ID] = true
				s6 = p
			}
		}
		missingPlayer = s1.ID == 0 || s2.ID == 0 || s3.ID == 0 || s4.ID == 0 || s5.ID == 0 || s6.ID == 0
	}

	queue = append(queue, s1, s2, s3, s4, s5, s6)

	return queue
}

func RecordPlay(gs *GameState, eventID, outcomeID, nextZoneID, injuryID, injuryType, injuryDuration, penaltyID, severity uint8, isFight bool, pcID, ppID, apID, dpID, gpID uint, isBreakaway bool) {
	_, zoneID := getZoneID(gs.PuckLocation, gs.HomeTeamID, gs.AwayTeamID)
	play := structs.PbP{
		GameID:                gs.GameID,
		Period:                gs.Period,
		TimeOnClock:           gs.TimeOnClock,
		SecondsConsumed:       uint8(gs.SecondsConsumed),
		EventID:               eventID,
		ZoneID:                uint8(zoneID),
		NextZoneID:            nextZoneID,
		Outcome:               outcomeID,
		HomeTeamScore:         gs.HomeTeamScore,
		AwayTeamScore:         gs.AwayTeamScore,
		HomeTeamShootoutScore: gs.HomeTeamShootoutScore,
		AwayTeamShootoutScore: gs.AwayTeamShootoutScore,
		TeamID:                uint8(gs.PossessingTeam),
		PuckCarrierID:         pcID,
		PassedPlayerID:        ppID,
		AssistingPlayerID:     apID,
		DefenderID:            dpID,
		GoalieID:              gpID,
		InjuryID:              injuryID,
		InjuryType:            injuryType,
		InjuryDuration:        injuryDuration,
		PenaltyID:             penaltyID,
		IsFight:               false,
		Severity:              severity,
		IsBreakaway:           isBreakaway,
		IsShootout:            gs.IsOvertimeShootout,
		HOS:                   gs.HomeStrategy.Gameplan.OffensiveSystem,
		AOS:                   gs.AwayStrategy.Gameplan.OffensiveSystem,
		HDS:                   gs.HomeStrategy.Gameplan.DefensiveSystem,
		ADS:                   gs.AwayStrategy.Gameplan.DefensiveSystem,
	}

	gs.RecordPlay(play)
}

// ============================================================================
// PUCK BATTLE SYSTEM
// ============================================================================

// handlePuckBattle manages contested puck scenarios between multiple players
func handlePuckBattle(gs *GameState, contestedPlayers []*GamePlayer) {
	if len(contestedPlayers) < 2 {
		// Not enough players for a battle, default to first player
		if len(contestedPlayers) == 1 {
			gs.SetPuckBearer(contestedPlayers[0], false)
		}
		return
	}

	// Set game state to contested
	gs.PuckState = PuckStateContested
	gs.ContestedPlayers = contestedPlayers

	// Calculate weights for each player based on system modifiers and attributes
	playerWeights := []PlayerWeight{}

	for _, player := range contestedPlayers {
		weight := calculatePuckBattleWeight(gs, player)
		playerWeights = append(playerWeights, PlayerWeight{
			PlayerID: player.ID,
			Weight:   weight,
		})
	}

	// Select winner based on weights
	totalWeight := 0.0
	for _, pw := range playerWeights {
		totalWeight += pw.Weight
	}
	winnerID := selectPlayerIDByWeights(totalWeight, playerWeights)

	// Find the winning player
	var winner *GamePlayer
	for _, player := range contestedPlayers {
		if player.ID == winnerID {
			winner = player
			break
		}
	}
	if winner == nil {
		winner = contestedPlayers[0] // Fallback
	}

	// Determine battle outcome
	secondsConsumed := util.GenerateIntFromRange(2, 5) // Battles take time
	gs.SetSecondsConsumed(uint16(secondsConsumed))

	// Separate players by team to track participants from both sides
	var homeTeamPlayers, awayTeamPlayers []*GamePlayer
	for _, player := range contestedPlayers {
		if player.TeamID == uint16(gs.HomeTeamID) {
			homeTeamPlayers = append(homeTeamPlayers, player)
		} else {
			awayTeamPlayers = append(awayTeamPlayers, player)
		}
	}

	// Find opposing team player and additional participant to track
	var opposingPlayer *GamePlayer
	var additionalPlayer *GamePlayer

	// Get a player from the opposing team (not the winner)
	if winner.TeamID == uint16(gs.HomeTeamID) {
		// Winner is home team, get away team player
		for _, player := range awayTeamPlayers {
			opposingPlayer = player
			break
		}
		// Get additional home team participant if available
		for _, player := range homeTeamPlayers {
			if player.ID != winner.ID {
				additionalPlayer = player
				break
			}
		}
	} else {
		// Winner is away team, get home team player
		for _, player := range homeTeamPlayers {
			opposingPlayer = player
			break
		}
		// Get additional away team participant if available
		for _, player := range awayTeamPlayers {
			if player.ID != winner.ID {
				additionalPlayer = player
				break
			}
		}
	}

	// Record the puck battle event with participant tracking
	eventID := PuckBattleID
	outcomeID := PuckBattleWinID
	if gs.PuckCarrier != nil && gs.PuckCarrier.TeamID != winner.TeamID {
		outcomeID = PuckBattleLoseID
	}

	// Track participants: winner as pcID, opposing team as dpID, additional as ppID
	opposingPlayerID := uint(0)
	if opposingPlayer != nil {
		opposingPlayerID = opposingPlayer.ID
	}
	additionalPlayerID := uint(0)
	if additionalPlayer != nil {
		additionalPlayerID = additionalPlayer.ID
	}

	RecordPlay(gs, eventID, outcomeID, 0, 0, 0, 0, 0, 0, false, winner.ID, additionalPlayerID, 0, opposingPlayerID, 0, false)

	// Winner takes possession
	gs.SetPuckBearer(winner, false)
}

// calculatePuckBattleWeight determines a player's effectiveness in puck battles
func calculatePuckBattleWeight(gs *GameState, player *GamePlayer) float64 {
	baseWeight := 1.0

	// Core attributes for puck battles
	strength := player.StrengthMod
	agility := player.AgilityMod
	handling := player.HandlingMod

	// Base calculation: strength for physical battles, agility for positioning, handling for puck control
	weight := baseWeight + (strength * 0.4) + (agility * 0.3) + (handling * 0.3)

	// Apply system modifiers based on team and zone
	isHome := player.TeamID == uint16(gs.HomeTeamID)
	homePossession := gs.PuckCarrier != nil && gs.PuckCarrier.TeamID == uint16(gs.HomeTeamID)

	// Get system modifiers
	homeModifiers, awayModifiers := GetSystemModifiersForZone(gs, homePossession, gs.PuckLocation)

	var systemModifiers structs.SystemModifiers
	if isHome {
		systemModifiers = homeModifiers
	} else {
		systemModifiers = awayModifiers
	}

	// Apply system-specific puck battle weight
	systemWeight := GetSystemPlayerWeight(player, systemModifiers, weight)

	// Zone-specific modifiers for puck battles
	zoneModifiers := GetZoneModifiersForEvent(systemModifiers, gs.PuckLocation, isHome)

	// Use offensive checking bonuses for puck battles (now they have a purpose!)
	battleBonus := 0.0
	if player.TeamID == gs.PuckCarrier.TeamID {
		// Offensive player in puck battle - use any offensive system checking bonuses
		battleBonus += float64(zoneModifiers.StickCheckBonus) * 0.1
		battleBonus += float64(zoneModifiers.BodyCheckBonus) * 0.1
	} else {
		// Defensive player in puck battle - use defensive checking bonuses
		battleBonus += float64(zoneModifiers.StickCheckBonus) * 0.15
		battleBonus += float64(zoneModifiers.BodyCheckBonus) * 0.15
	}

	return math.Max(systemWeight+battleBonus, 0.1) // Minimum weight
}

// triggerPuckBattle initiates a contested puck scenario
func triggerPuckBattle(gs *GameState, offensivePlayer, defensivePlayer *GamePlayer) {
	// Gather nearby players for the battle
	contestedPlayers := []*GamePlayer{offensivePlayer, defensivePlayer}

	// Add additional players based on zone and situation
	additionalPlayers := getNearbyPlayersForBattle(gs, gs.PuckLocation)
	contestedPlayers = append(contestedPlayers, additionalPlayers...)

	// Remove duplicates
	contestedPlayers = removeDuplicatePlayers(contestedPlayers)

	// Handle the battle
	handlePuckBattle(gs, contestedPlayers)
}

// getNearbyPlayersForBattle gets additional players who might join the battle
func getNearbyPlayersForBattle(gs *GameState, zone string) []*GamePlayer {
	nearbyPlayers := []*GamePlayer{}

	// In goal zones, more players get involved
	switch zone {
	case HomeGoal, AwayGoal:
		// More chaotic battles near the net - up to 2 additional players
		if util.GenerateIntFromRange(1, 100) <= 60 { // 60% chance
			// Add a forward from each team
			homeForward := getRandomPlayerByPosition(gs, true, []string{Center, Forward})
			if homeForward != nil {
				nearbyPlayers = append(nearbyPlayers, homeForward)
			}

			awayForward := getRandomPlayerByPosition(gs, false, []string{Center, Forward})
			if awayForward != nil {
				nearbyPlayers = append(nearbyPlayers, awayForward)
			}
		}
	case HomeZone, AwayZone:
		// Medium chance of additional player in regular zones
		if util.GenerateIntFromRange(1, 100) <= 40 { // 40% chance
			// Add one additional player
			isHome := util.GenerateIntFromRange(1, 2) == 1
			player := getRandomPlayerByPosition(gs, isHome, []string{Center, Forward, Defender})
			if player != nil {
				nearbyPlayers = append(nearbyPlayers, player)
			}
		}
	case NeutralZone:
		// Lower chance in neutral zone
		if util.GenerateIntFromRange(1, 100) <= 25 { // 25% chance
			isHome := util.GenerateIntFromRange(1, 2) == 1
			player := getRandomPlayerByPosition(gs, isHome, []string{Center, Forward})
			if player != nil {
				nearbyPlayers = append(nearbyPlayers, player)
			}
		}
	}

	return nearbyPlayers
}

// getRandomPlayerByPosition selects a random player from current lineup by position
func getRandomPlayerByPosition(gs *GameState, isHome bool, positions []string) *GamePlayer {
	var lineup []*GamePlayer

	if isHome {
		// Combine forwards and defenders for home team
		hgs := &gs.HomeStrategy
		lineup = append(lineup, hgs.Forwards[hgs.CurrentForwards].Players...)
		lineup = append(lineup, hgs.Defenders[hgs.CurrentDefenders].Players...)
	} else {
		// Combine forwards and defenders for away team
		ags := &gs.AwayStrategy
		lineup = append(lineup, ags.Forwards[ags.CurrentForwards].Players...)
		lineup = append(lineup, ags.Defenders[ags.CurrentDefenders].Players...)
	}

	// Filter by positions
	validPlayers := []*GamePlayer{}
	for _, player := range lineup {
		for _, pos := range positions {
			if player.Position == pos && !player.IsOut {
				validPlayers = append(validPlayers, player)
				break
			}
		}
	}

	if len(validPlayers) == 0 {
		return nil
	}

	// Return random valid player
	idx := util.GenerateIntFromRange(0, len(validPlayers)-1)
	return validPlayers[idx]
}

// removeDuplicatePlayers removes duplicate players from the list
func removeDuplicatePlayers(players []*GamePlayer) []*GamePlayer {
	seen := make(map[uint]bool)
	unique := []*GamePlayer{}

	for _, player := range players {
		if !seen[player.ID] && player != nil {
			seen[player.ID] = true
			unique = append(unique, player)
		}
	}

	return unique
}

// modifyDefenseCheckForPuckBattles updates defense check to sometimes trigger battles
func shouldTriggerPuckBattle(gs *GameState, attackWeight, defenseWeight float64) bool {
	// Close battles (within 20% of each other) can become contested
	diff := math.Abs(attackWeight - defenseWeight)
	maxWeight := math.Max(attackWeight, defenseWeight)

	if maxWeight == 0 {
		return false
	}

	ratio := diff / maxWeight

	// If weights are close (less than 30% difference), chance of puck battle
	if ratio <= 0.3 {
		battleChance := util.GenerateIntFromRange(1, 100)

		// Zone-dependent battle chances
		switch gs.PuckLocation {
		case HomeGoal, AwayGoal:
			return battleChance <= 45 // 45% chance in goal zones
		case HomeZone, AwayZone:
			return battleChance <= 30 // 30% chance in regular zones
		case NeutralZone:
			return battleChance <= 20 // 20% chance in neutral zone
		}
	}

	return false
}

// handleCoveredPuck manages scenarios where goalie covers the puck
func handleCoveredPuck(gs *GameState, goalie *GamePlayer) {
	gs.PuckState = PuckStateCovered
	gs.ContestedPlayers = []*GamePlayer{}

	// Covered puck leads to faceoff
	secondsConsumed := util.GenerateIntFromRange(3, 6)
	gs.SetSecondsConsumed(uint16(secondsConsumed))

	// Record the covered puck event
	RecordPlay(gs, PenaltyCheckID, 0, 0, 0, 0, 0, 0, 0, false, 0, 0, 0, 0, goalie.ID, false)

	// Set up faceoff
	gs.SetFaceoffOnCenterIce(false)        // Faceoff in zone where puck was covered
	gs.SetPuckBearer(&GamePlayer{}, false) // Clear puck carrier for faceoff
}

// updateGameStateWithPuckBattle manages game state during puck battles
func (gs *GameState) UpdatePuckBattleState(contested []*GamePlayer) {
	gs.PuckState = PuckStateContested
	gs.ContestedPlayers = contested
}

// clearPuckBattleState resets game state after puck battle resolution
func (gs *GameState) ClearPuckBattleState() {
	gs.PuckState = PuckStateClear
	gs.ContestedPlayers = []*GamePlayer{}
}
