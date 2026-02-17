package engine

import (
	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/structs"
)

// GetSystemModifiersForZone returns the system modifiers for a specific zone and possession state
func GetSystemModifiersForZone(gs *GameState, isHomePossession bool, zone string) (structs.SystemModifiers, structs.SystemModifiers) {
	// Get gameplans from playbooks
	homePlaybook := gs.GetPlaybook(true)
	awayPlaybook := gs.GetPlaybook(false)

	// Extract system information from gameplans
	homeGameplan := homePlaybook.Gameplan
	awayGameplan := awayPlaybook.Gameplan

	homeOffensiveSystem := homeGameplan.OffensiveSystem
	homeDefensiveSystem := homeGameplan.DefensiveSystem
	awayOffensiveSystem := awayGameplan.OffensiveSystem
	awayDefensiveSystem := awayGameplan.DefensiveSystem

	homeOffensiveIntensity := homeGameplan.OffensiveIntensity
	homeDefensiveIntensity := homeGameplan.DefensiveIntensity
	awayOffensiveIntensity := awayGameplan.OffensiveIntensity
	awayDefensiveIntensity := awayGameplan.DefensiveIntensity

	// Get modifiers based on possession
	var homeModifiers, awayModifiers structs.SystemModifiers

	if isHomePossession {
		// Home team is on offense, away team is on defense
		homeModifiers = structs.GetOffensiveSystemModifiers(structs.OffensiveSystemType(homeOffensiveSystem), homeOffensiveIntensity)
		awayModifiers = structs.GetDefensiveSystemModifiers(structs.DefensiveSystemType(awayDefensiveSystem), awayDefensiveIntensity)
	} else {
		// Away team is on offense, home team is on defense
		homeModifiers = structs.GetDefensiveSystemModifiers(structs.DefensiveSystemType(homeDefensiveSystem), homeDefensiveIntensity)
		awayModifiers = structs.GetOffensiveSystemModifiers(structs.OffensiveSystemType(awayOffensiveSystem), awayOffensiveIntensity)
	}

	return homeModifiers, awayModifiers
}

// GetZoneModifiersForEvent returns the specific zone modifiers for the current event
func GetZoneModifiersForEvent(modifiers structs.SystemModifiers, zone string, isHome bool) structs.ZoneModifiers {
	switch zone {
	case HomeGoal:
		if isHome {
			return modifiers.DefendingGoalZone // Home team defending their goal
		}
		return modifiers.AttackingGoalZone // Away team attacking home goal

	case HomeZone:
		if isHome {
			return modifiers.DefendingZone // Home team defending their zone
		}
		return modifiers.AttackingZone // Away team attacking home zone

	case NeutralZone:
		return modifiers.NeutralZone

	case AwayZone:
		if !isHome {
			return modifiers.DefendingZone // Away team defending their zone
		}
		return modifiers.AttackingZone // Home team attacking away zone

	case AwayGoal:
		if !isHome {
			return modifiers.DefendingGoalZone // Away team defending their goal
		}
		return modifiers.AttackingGoalZone // Home team attacking away goal
	}

	return structs.ZoneModifiers{} // No modifiers
}

// ApplySystemModifiersToEventWeights adjusts event weights based on systems
func ApplySystemModifiersToEventWeights(gs *GameState, eventWeights *EventWeights) {
	pc := gs.PuckCarrier
	isHomePossession := pc.TeamID == uint16(gs.HomeTeamID)

	homeModifiers, awayModifiers := GetSystemModifiersForZone(gs, isHomePossession, gs.PuckLocation)

	// Apply modifiers to the possessing team
	var modifiers structs.SystemModifiers
	var zoneModifiers structs.ZoneModifiers

	if isHomePossession {
		modifiers = homeModifiers
		zoneModifiers = GetZoneModifiersForEvent(modifiers, gs.PuckLocation, true)
	} else {
		modifiers = awayModifiers
		zoneModifiers = GetZoneModifiersForEvent(modifiers, gs.PuckLocation, false)
	}

	// Apply zone-specific bonuses to event weights
	eventWeights.ShotWeight += int(zoneModifiers.ShotBonus)
	eventWeights.PassWeight += int(zoneModifiers.PassBonus)
	eventWeights.AgilityWeight += int(zoneModifiers.AgilityBonus)

	// Apply defensive modifiers to defending team
	var defenseModifiers structs.SystemModifiers
	var defenseZoneModifiers structs.ZoneModifiers

	if isHomePossession {
		defenseModifiers = awayModifiers
		defenseZoneModifiers = GetZoneModifiersForEvent(defenseModifiers, gs.PuckLocation, false)
	} else {
		defenseModifiers = homeModifiers
		defenseZoneModifiers = GetZoneModifiersForEvent(defenseModifiers, gs.PuckLocation, true)
	}

	eventWeights.StickCheckWeight += int(defenseZoneModifiers.StickCheckBonus)
	eventWeights.BodyCheckWeight += int(defenseZoneModifiers.BodyCheckBonus)
}

// GetSystemPlayerWeight calculates player selection weight based on system preferences
// Excludes goalies to preserve user choice in goalie selection
func GetSystemPlayerWeight(player *GamePlayer, modifiers structs.SystemModifiers, baseWeight float64) float64 {
	// Skip goalies - let users choose their preferred goalie style
	if player.Position == util.Goalie {
		return baseWeight
	}

	systemWeight := baseWeight

	// Apply archetype bonuses from the system for skaters only
	if archetypeBonus, exists := modifiers.ArchetypeWeights[player.Archetype]; exists {
		systemWeight += float64(archetypeBonus) * 0.1 // Scale the bonus appropriately

		// Apply position-specific validation and adjustments
		switch player.Position {
		case util.Center, util.Forward:
			// Valid forward archetypes get proper bonuses
			switch player.Archetype {
			case util.Offensive, util.Defensive:
				// Invalid archetype for forwards
				systemWeight -= 0.5 // Heavy penalty
			case util.Grinder:
				// Defensive forwards excel in certain systems
				systemWeight += 0.02
			}

		case util.Defender:
			// Valid defenseman archetypes get proper bonuses
			switch player.Archetype {
			case util.Grinder, util.Playmaker, util.Power, util.Sniper:
				// Invalid archetype for defensemen
				systemWeight -= 0.5 // Heavy penalty
			case util.Offensive:
				// Offensive defensemen get bonus
				systemWeight += 0.05
			case util.Defensive:
				// Defensive defensemen get bonus
				systemWeight += 0.05
			}
		}
	}

	return systemWeight
}

// ApplySystemAttributeModifiers applies system compatibility bonuses to player attributes
// Excludes goalies to preserve user choice in goalie selection
func ApplySystemAttributeModifiers(player *GamePlayer, offSystem structs.OffensiveSystemType, defSystem structs.DefensiveSystemType) {
	// Skip goalies - let users choose their preferred goalie style
	if player.Position == util.Goalie {
		return
	}

	modifier := structs.GetSystemAttributeModifier(offSystem, defSystem, player.Archetype, player.Position)

	if modifier != 1.0 {
		// Apply modifier to skater attributes only
		player.AgilityMod *= modifier
		player.StrengthMod *= modifier
		player.LongShotAccMod *= modifier
		player.LongShotPowerMod *= modifier
		player.CloseShotAccMod *= modifier
		player.CloseShotPowerMod *= modifier
		player.FaceoffMod *= modifier
		player.HandlingMod *= modifier
		player.PassMod *= modifier
		player.StickCheckMod *= modifier
		player.BodyCheckMod *= modifier
		player.ShotblockingMod *= modifier
	}
}

// EventWeights represents the calculated weights for different events in a zone
type EventWeights struct {
	ShotWeight       int
	PassWeight       int
	PassBackWeight   int
	LongPassWeight   int
	AgilityWeight    int
	StickCheckWeight int
	BodyCheckWeight  int
	TotalWeight      int
}

// CalculateEventWeights computes the base event weights for the current situation
func CalculateEventWeights(gs *GameState) EventWeights {
	pc := gs.PuckCarrier
	isHome := pc.TeamID == uint16(gs.HomeTeamID)

	attackStrategy := gs.GetLineStrategy(isHome, 1)
	defendStrategy := gs.GetLineStrategy(!isHome, 2)

	weights := EventWeights{}

	// Base weights from your existing system
	switch gs.PuckLocation {
	case HomeGoal, AwayGoal:
		if (gs.PuckLocation == HomeGoal && !isHome) || (gs.PuckLocation == AwayGoal && isHome) {
			// Offensive goal zone
			weights.ShotWeight = int(attackStrategy.AGZShot) + int(pc.AGZShot) + int(gs.Momentum)
			if weights.ShotWeight < 1 {
				weights.ShotWeight = 1 // Ensure minimum weight for shot attempts in goal zone
			}
			weights.PassWeight = int(attackStrategy.AGZPass) + int(pc.AGZPass)
			if weights.PassWeight < 1 {
				weights.PassWeight = 1 // Ensure minimum weight for passes in goal zone
			}
			weights.PassBackWeight = int(attackStrategy.AGZPassBack) + int(pc.AGZPassBack)
			if weights.PassBackWeight < 1 {
				weights.PassBackWeight = 1 // Ensure minimum weight for pass backs in goal zone
			}
			weights.StickCheckWeight = int(defendStrategy.DGZStickCheck)
			if weights.StickCheckWeight < 1 {
				weights.StickCheckWeight = 1 // Ensure minimum weight for stick checks in goal zone
			}
			weights.BodyCheckWeight = int(defendStrategy.DGZBodyCheck)
			if weights.BodyCheckWeight < 1 {
				weights.BodyCheckWeight = 1 // Ensure minimum weight for body checks in goal zone
			}
		} else {
			// Defensive goal zone
			weights.AgilityWeight = int(attackStrategy.DGZAgility) + int(pc.DGZAgility)
			if weights.AgilityWeight < 1 {
				weights.AgilityWeight = 1 // Ensure minimum weight for agility in defensive goal zone
			}
			weights.PassWeight = int(attackStrategy.DGZPass) + int(pc.DGZPass)
			if weights.PassWeight < 1 {
				weights.PassWeight = 1 // Ensure minimum weight for passes in defensive goal zone
			}
			weights.LongPassWeight = int(attackStrategy.DGZLongPass) + int(pc.DGZLongPass)
			if weights.LongPassWeight < 1 {
				weights.LongPassWeight = 1 // Ensure minimum weight for long passes in defensive goal zone
			}
			weights.StickCheckWeight = int(defendStrategy.AGZStickCheck)
			if weights.StickCheckWeight < 1 {
				weights.StickCheckWeight = 1 // Ensure minimum weight for stick checks in defensive goal zone
			}
			weights.BodyCheckWeight = int(defendStrategy.AGZBodyCheck)
			if weights.BodyCheckWeight < 1 {
				weights.BodyCheckWeight = 1 // Ensure minimum weight for body checks in defensive goal zone
			}
		}

	case HomeZone, AwayZone:
		if (gs.PuckLocation == HomeZone && !isHome) || (gs.PuckLocation == AwayZone && isHome) {
			// Offensive zone
			weights.ShotWeight = int(attackStrategy.AZShot) + int(pc.AZShot) + int(gs.Momentum)
			if weights.ShotWeight < 1 {
				weights.ShotWeight = 1 // Ensure minimum weight for shot attempts in offensive zone
			}
			weights.AgilityWeight = int(attackStrategy.AZAgility) + int(pc.AZAgility)
			if weights.AgilityWeight < 1 {
				weights.AgilityWeight = 1 // Ensure minimum weight for agility in offensive zone
			}
			weights.PassWeight = int(attackStrategy.AZPass) + int(pc.AZPass)
			if weights.PassWeight < 1 {
				weights.PassWeight = 1 // Ensure minimum weight for passes in offensive zone
			}
			weights.LongPassWeight = int(attackStrategy.AZLongPass) + int(pc.AZLongPass)
			if weights.LongPassWeight < 1 {
				weights.LongPassWeight = 1 // Ensure minimum weight for long passes in offensive zone
			}
			weights.StickCheckWeight = int(defendStrategy.DZStickCheck)
			if weights.StickCheckWeight < 1 {
				weights.StickCheckWeight = 1 // Ensure minimum weight for stick checks in offensive zone
			}
			weights.BodyCheckWeight = int(defendStrategy.DZBodyCheck)
			if weights.BodyCheckWeight < 1 {
				weights.BodyCheckWeight = 1 // Ensure minimum weight for body checks in offensive zone
			}
		} else {
			// Defensive zone
			weights.AgilityWeight = int(attackStrategy.DZAgility) + int(pc.DZAgility)
			if weights.AgilityWeight < 1 {
				weights.AgilityWeight = 1 // Ensure minimum weight for agility in defensive zone
			}
			weights.PassWeight = int(attackStrategy.DZPass) + int(pc.DZPass)
			if weights.PassWeight < 1 {
				weights.PassWeight = 1 // Ensure minimum weight for passes in defensive zone
			}
			weights.StickCheckWeight = int(defendStrategy.AZStickCheck)
			if weights.StickCheckWeight < 1 {
				weights.StickCheckWeight = 1 // Ensure minimum weight for stick checks in defensive zone
			}
			weights.BodyCheckWeight = int(defendStrategy.AZBodyCheck)
			if weights.BodyCheckWeight < 1 {
				weights.BodyCheckWeight = 1 // Ensure minimum weight for body checks in defensive zone
			}
		}

	case NeutralZone:
		weights.AgilityWeight = int(attackStrategy.NAgility) + int(pc.NAgility)
		if weights.AgilityWeight < 1 {
			weights.AgilityWeight = 1 // Ensure minimum weight for agility in neutral zone
		}
		weights.PassWeight = int(attackStrategy.NPass) + int(pc.NPass)
		if weights.PassWeight < 1 {
			weights.PassWeight = 1 // Ensure minimum weight for passes in neutral zone
		}
		weights.StickCheckWeight = int(defendStrategy.NStickCheck)
		if weights.StickCheckWeight < 1 {
			weights.StickCheckWeight = 1 // Ensure minimum weight for stick checks in neutral zone
		}
		weights.BodyCheckWeight = int(defendStrategy.NBodyCheck)
		if weights.BodyCheckWeight < 1 {
			weights.BodyCheckWeight = 1 // Ensure minimum weight for body checks in neutral zone
		}
	}

	weights.TotalWeight = weights.ShotWeight + weights.PassWeight + weights.AgilityWeight +
		weights.StickCheckWeight + weights.BodyCheckWeight + 1 // +1 for penalty

	return weights
}
