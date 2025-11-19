package structs

import util "github.com/CalebRose/SimHockey/_util"

// SystemModifiers defines how each system affects different game situations
type SystemModifiers struct {
	// Zone-based modifiers (-10 to +10)
	AttackingGoalZone ZoneModifiers
	AttackingZone     ZoneModifiers
	NeutralZone       ZoneModifiers
	DefendingZone     ZoneModifiers
	DefendingGoalZone ZoneModifiers

	// Special situation modifiers
	PowerPlayModifier   int8
	PenaltyKillModifier int8

	// Archetype preferences (higher values = better fit)
	ArchetypeWeights map[string]int8
}

type ZoneModifiers struct {
	ShotBonus       int8 // Bonus to shot attempts
	PassBonus       int8 // Bonus to passing
	AgilityBonus    int8 // Bonus to agility/movement
	StickCheckBonus int8 // Bonus to stick checking
	BodyCheckBonus  int8 // Bonus to body checking
	TurnoverChance  int8 // Modifier to turnover probability
}

// OffensiveSystemType represents the offensive systems
type OffensiveSystemType uint8

const (
	Offensive122Forecheck OffensiveSystemType = iota + 1
	Offensive212Forecheck
	Offensive113Forecheck
	OffensiveCycleGame
	OffensiveQuickTransition
	OffensiveUmbrella
	OffensiveEastWestMotion
	OffensiveCrashNet
)

// DefensiveSystemType represents the defensive systems
type DefensiveSystemType uint8

const (
	DefensiveBalanced DefensiveSystemType = iota + 1
	DefensiveManToMan
	DefensiveZone
	DefensiveNeutralTrap
	DefensiveLeftWingLock
	DefensiveAggressiveForecheck
	DefensiveCollapsing
	DefensiveBox
)

// GetOffensiveSystemModifiers returns the modifiers for a given offensive system
func GetOffensiveSystemModifiers(system OffensiveSystemType, intensity uint8) SystemModifiers {
	baseIntensity := float64(intensity) / 5.0 // Normalize to 0-2 multiplier

	switch system {
	case Offensive122Forecheck:
		return SystemModifiers{
			AttackingZone: ZoneModifiers{
				ShotBonus:       int8(2 * baseIntensity),
				PassBonus:       int8(3 * baseIntensity),
				StickCheckBonus: int8(4 * baseIntensity),
				BodyCheckBonus:  int8(2 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				PassBonus:       int8(2 * baseIntensity),
				AgilityBonus:    int8(3 * baseIntensity),
				StickCheckBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Grinder:   int8(3 * baseIntensity),
				util.TwoWay:    int8(4 * baseIntensity),
				util.Playmaker: int8(2 * baseIntensity),
				util.Defensive: int8(3 * baseIntensity),  // Defensive D good for balanced forecheck
				util.Power:     int8(-2 * baseIntensity), // Anti-fit: Power forwards too slow for quick forechecking
			},
		}

	case Offensive212Forecheck:
		return SystemModifiers{
			AttackingZone: ZoneModifiers{
				ShotBonus:       int8(4 * baseIntensity),
				StickCheckBonus: int8(3 * baseIntensity),
				TurnoverChance:  int8(5 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				AgilityBonus:    int8(4 * baseIntensity),
				StickCheckBonus: int8(4 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Grinder:   int8(4 * baseIntensity),
				util.Enforcer:  int8(5 * baseIntensity),
				util.Playmaker: int8(3 * baseIntensity),
				util.TwoWay:    int8(2 * baseIntensity),
				util.Sniper:    int8(-2 * baseIntensity), // Anti-fit: poor defense for aggressive system
				util.Defensive: int8(2 * baseIntensity),  // Defensive D help with aggressive forechecking
			},
		}

	case OffensiveCycleGame:
		return SystemModifiers{
			AttackingGoalZone: ZoneModifiers{
				ShotBonus: int8(5 * baseIntensity),
				PassBonus: int8(4 * baseIntensity),
			},
			AttackingZone: ZoneModifiers{
				PassBonus: int8(5 * baseIntensity),
				ShotBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Playmaker: int8(5 * baseIntensity),
				util.Power:     int8(4 * baseIntensity),
				util.Sniper:    int8(4 * baseIntensity),
				util.Grinder:   int8(-3 * baseIntensity),
				util.Enforcer:  int8(-2 * baseIntensity),
			},
		}

	case OffensiveQuickTransition:
		return SystemModifiers{
			NeutralZone: ZoneModifiers{
				PassBonus:    int8(5 * baseIntensity),
				AgilityBonus: int8(5 * baseIntensity),
			},
			DefendingZone: ZoneModifiers{
				PassBonus:    int8(4 * baseIntensity),
				AgilityBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Offensive: int8(5 * baseIntensity), // Offensive Defenseman
				util.Sniper:    int8(4 * baseIntensity),
				util.TwoWay:    int8(3 * baseIntensity),
				util.Enforcer:  int8(-3 * baseIntensity), // Anti-fit: enforcers are too slow for quick transition
				util.Power:     int8(-6 * baseIntensity), // Anti-fit: enforcers are too slow for quick transition
			},
		}

	case OffensiveCrashNet:
		return SystemModifiers{
			AttackingGoalZone: ZoneModifiers{
				ShotBonus:      int8(6 * baseIntensity),
				BodyCheckBonus: int8(4 * baseIntensity),
			},
			AttackingZone: ZoneModifiers{
				ShotBonus:      int8(4 * baseIntensity),
				BodyCheckBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Power:     int8(6 * baseIntensity), // Power forwards perfect for crashing net - close shot power + strength
				util.Enforcer:  int8(4 * baseIntensity),
				util.Grinder:   int8(3 * baseIntensity),
				util.Playmaker: int8(-3 * baseIntensity), // Anti-fit: finesse doesn't suit crashing
				util.Sniper:    int8(-3 * baseIntensity), // Anti-fit: worse strength, poor for net-front play
			},
		}

	case Offensive113Forecheck:
		return SystemModifiers{
			AttackingZone: ZoneModifiers{
				ShotBonus:      int8(6 * baseIntensity), // Heavy shot focus with 3 forwards
				PassBonus:      int8(3 * baseIntensity),
				TurnoverChance: int8(3 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				AgilityBonus: int8(-2 * baseIntensity), // Less agility with only 1 forechecker
			},
			ArchetypeWeights: map[string]int8{
				util.Sniper:    int8(5 * baseIntensity),  // Perfect for 3-forward attack
				util.Playmaker: int8(4 * baseIntensity),  // Helps set up 3 forwards
				util.Power:     int8(4 * baseIntensity),  // Good for net-front in 3-forward setup
				util.Grinder:   int8(-3 * baseIntensity), // Anti-fit: Defensive forwards poor in offensive system
				util.Defensive: int8(-4 * baseIntensity), // Anti-fit: Too defensive for aggressive offense
			},
		}

	case OffensiveUmbrella:
		return SystemModifiers{
			AttackingGoalZone: ZoneModifiers{
				ShotBonus: int8(4 * baseIntensity), // Good shots from umbrella setup
				PassBonus: int8(5 * baseIntensity), // Excellent passing in umbrella
			},
			AttackingZone: ZoneModifiers{
				PassBonus: int8(6 * baseIntensity), // Core of umbrella is passing
			},
			ArchetypeWeights: map[string]int8{
				util.Playmaker: int8(6 * baseIntensity), // Perfect for umbrella passing
				util.Offensive: int8(5 * baseIntensity), // Offensive D crucial for umbrella
				util.Sniper:    int8(4 * baseIntensity), // Good for umbrella shooting
				util.TwoWay:    int8(3 * baseIntensity),
				util.Enforcer:  int8(-4 * baseIntensity), // Anti-fit: Too aggressive for finesse system
				util.Grinder:   int8(-3 * baseIntensity), // Anti-fit: Not skilled enough for umbrella
			},
		}

	case OffensiveEastWestMotion:
		return SystemModifiers{
			AttackingGoalZone: ZoneModifiers{
				PassBonus:    int8(5 * baseIntensity), // Lots of lateral passing
				AgilityBonus: int8(4 * baseIntensity), // Need agility for east-west movement
			},
			AttackingZone: ZoneModifiers{
				PassBonus:    int8(6 * baseIntensity), // Heavy passing system
				AgilityBonus: int8(5 * baseIntensity), // Constant movement
			},
			ArchetypeWeights: map[string]int8{
				util.Playmaker: int8(6 * baseIntensity),  // Perfect for east-west passing
				util.Sniper:    int8(4 * baseIntensity),  // Good agility and passing
				util.TwoWay:    int8(3 * baseIntensity),  // Versatile for motion system
				util.Offensive: int8(3 * baseIntensity),  // Offensive D help with puck movement
				util.Power:     int8(-3 * baseIntensity), // Anti-fit: Too slow for east-west motion
				util.Enforcer:  int8(-4 * baseIntensity), // Anti-fit: Too aggressive/slow for finesse
			},
		}

	default:
		return SystemModifiers{} // No modifiers for unknown systems
	}
}

// GetDefensiveSystemModifiers returns the modifiers for a given defensive system
func GetDefensiveSystemModifiers(system DefensiveSystemType, intensity uint8) SystemModifiers {
	baseIntensity := float64(intensity) / 5.0 // Normalize to 0-2 multiplier

	switch system {
	case DefensiveManToMan:
		return SystemModifiers{
			DefendingZone: ZoneModifiers{
				StickCheckBonus: int8(3 * baseIntensity),
				BodyCheckBonus:  int8(5 * baseIntensity),
			},
			DefendingGoalZone: ZoneModifiers{
				StickCheckBonus: int8(4 * baseIntensity),
				BodyCheckBonus:  int8(6 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Defensive: int8(4 * baseIntensity), // Defensive Defenseman
				util.Grinder:   int8(3 * baseIntensity),
			},
		}

	case DefensiveNeutralTrap:
		return SystemModifiers{
			NeutralZone: ZoneModifiers{
				StickCheckBonus: int8(5 * baseIntensity),
				AgilityBonus:    int8(4 * baseIntensity),
				TurnoverChance:  int8(4 * baseIntensity),
			},
			AttackingZone: ZoneModifiers{
				StickCheckBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Grinder:   int8(5 * baseIntensity), // Grinders are "defensive forwards" - perfect for trap
				util.Defensive: int8(4 * baseIntensity),
				util.TwoWay:    int8(3 * baseIntensity),
				util.Offensive: int8(-3 * baseIntensity), // Anti-fit: offensive defensemen don't suit trap
				util.Sniper:    int8(-2 * baseIntensity), // Anti-fit: worst defense, poor for trap system
			},
		}

	case DefensiveLeftWingLock:
		return SystemModifiers{
			DefendingZone: ZoneModifiers{
				StickCheckBonus: int8(4 * baseIntensity),
				BodyCheckBonus:  int8(3 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				StickCheckBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.TwoWay:    int8(5 * baseIntensity), // Especially for LW
				util.Defensive: int8(4 * baseIntensity),
				util.Grinder:   int8(3 * baseIntensity),  // Good defensive forwards
				util.Offensive: int8(-3 * baseIntensity), // Anti-fit: Offensive D don't suit structured defense
				util.Power:     int8(-2 * baseIntensity), // Anti-fit: Power forwards poor at disciplined defense
			},
		}

	case DefensiveAggressiveForecheck:
		return SystemModifiers{
			AttackingZone: ZoneModifiers{
				BodyCheckBonus:  int8(6 * baseIntensity),
				StickCheckBonus: int8(4 * baseIntensity),
				TurnoverChance:  int8(5 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				BodyCheckBonus: int8(5 * baseIntensity),
				AgilityBonus:   int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Enforcer:  int8(6 * baseIntensity),  // Increased from 5 - perfect for maximum physicality
				util.Grinder:   int8(5 * baseIntensity),  // Increased from 4 - great for relentless pressure
				util.TwoWay:    int8(3 * baseIntensity),  // Added - can handle both aggression and responsibility
				util.Playmaker: int8(-4 * baseIntensity), // Increased penalty from -3 - really conflicts with physical style
				util.Sniper:    int8(-4 * baseIntensity), // Keep same - terrible for aggressive forecheck
			},
		}

	case DefensiveCollapsing:
		return SystemModifiers{
			DefendingGoalZone: ZoneModifiers{
				BodyCheckBonus:  int8(4 * baseIntensity),
				StickCheckBonus: int8(3 * baseIntensity),
			},
			DefendingZone: ZoneModifiers{
				BodyCheckBonus:  int8(3 * baseIntensity),
				StickCheckBonus: int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Defensive: int8(5 * baseIntensity),
				util.TwoWay:    int8(3 * baseIntensity),  // Two-Way players help with structure
				util.Enforcer:  int8(2 * baseIntensity),  // Physical presence near goal
				util.Offensive: int8(-4 * baseIntensity), // Anti-fit: Offensive D poor at collapsing
				util.Sniper:    int8(-3 * baseIntensity), // Anti-fit: Poor defense, bad for goal-area defense
			},
		}

	case DefensiveBox:
		return SystemModifiers{
			DefendingGoalZone: ZoneModifiers{
				StickCheckBonus: int8(4 * baseIntensity),
				BodyCheckBonus:  int8(3 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.Defensive: int8(4 * baseIntensity),
				util.Grinder:   int8(3 * baseIntensity),
				util.TwoWay:    int8(3 * baseIntensity),  // Structured defensive play
				util.Enforcer:  int8(2 * baseIntensity),  // Physical presence
				util.Offensive: int8(-3 * baseIntensity), // Anti-fit: Poor at structured defense
				util.Playmaker: int8(-2 * baseIntensity), // Anti-fit: Too finesse for box structure
			},
		}

	case DefensiveBalanced:
		return SystemModifiers{
			// Balanced system - slight bonuses across all zones but no extremes
			DefendingZone: ZoneModifiers{
				StickCheckBonus: int8(2 * baseIntensity),
				BodyCheckBonus:  int8(2 * baseIntensity),
			},
			DefendingGoalZone: ZoneModifiers{
				StickCheckBonus: int8(2 * baseIntensity),
				BodyCheckBonus:  int8(2 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				PassBonus: int8(1 * baseIntensity),
			},
			ArchetypeWeights: map[string]int8{
				util.TwoWay:    int8(4 * baseIntensity), // Perfect for balanced system
				util.Defensive: int8(2 * baseIntensity),
				util.Grinder:   int8(2 * baseIntensity),
				util.Offensive: int8(1 * baseIntensity),
				util.Enforcer:  int8(-1 * baseIntensity), // Anti-fit: Too one-dimensional for balanced
			},
		}

	case DefensiveZone:
		return SystemModifiers{
			DefendingZone: ZoneModifiers{
				StickCheckBonus: int8(3 * baseIntensity),
				PassBonus:       int8(4 * baseIntensity), // Zone coverage requires good passing
			},
			DefendingGoalZone: ZoneModifiers{
				StickCheckBonus: int8(4 * baseIntensity),
				BodyCheckBonus:  int8(3 * baseIntensity),
			},
			NeutralZone: ZoneModifiers{
				PassBonus: int8(3 * baseIntensity), // Zone transitions need passing
			},
			ArchetypeWeights: map[string]int8{
				util.Defensive: int8(4 * baseIntensity), // Good for structured zone play
				util.TwoWay:    int8(3 * baseIntensity), // Versatile for zone coverage
				util.Grinder:   int8(2 * baseIntensity),
				util.Playmaker: int8(2 * baseIntensity),  // Forwards help with zone passing
				util.Enforcer:  int8(-2 * baseIntensity), // Anti-fit: Too aggressive for structured zones
				util.Power:     int8(-2 * baseIntensity), // Anti-fit: Too slow for zone coverage
			},
		}

	default:
		return SystemModifiers{} // No modifiers for unknown systems
	}
}

// GetSystemCompatibility checks how well player archetypes fit with systems
// Note: Goalie archetypes are intentionally excluded to preserve user choice in goalie selection
func GetSystemCompatibility(offSystem OffensiveSystemType, defSystem DefensiveSystemType, archetype string, position string) int8 {
	// Skip goalies entirely
	if position == "G" {
		return 0
	}

	offMods := GetOffensiveSystemModifiers(offSystem, 5) // Use medium intensity for compatibility check
	defMods := GetDefensiveSystemModifiers(defSystem, 5)

	compatibility := int8(0)

	if offWeight, exists := offMods.ArchetypeWeights[archetype]; exists {
		compatibility += offWeight
	}

	if defWeight, exists := defMods.ArchetypeWeights[archetype]; exists {
		compatibility += defWeight
	}

	// Apply position-specific logic for certain archetypes
	compatibility = applyPositionModifiers(compatibility, archetype, position, offSystem, defSystem)

	return compatibility
}

// applyPositionModifiers adjusts compatibility based on position-archetype combinations
func applyPositionModifiers(baseCompatibility int8, archetype string, position string, offSystem OffensiveSystemType, defSystem DefensiveSystemType) int8 {
	compatibility := baseCompatibility

	// Validate archetype-position combinations
	switch position {
	case util.Center, util.Forward:
		// Valid forward archetypes: Enforcer, Grinder, Playmaker, Power, Sniper, Two-Way
		switch archetype {
		case util.Offensive, util.Defensive:
			// These archetypes don't exist for forwards - this should never happen
			compatibility -= 10 // Heavy penalty for invalid combination
		case util.Grinder:
			// Grinder forwards excel in defensive systems
			if defSystem == DefensiveNeutralTrap || defSystem == DefensiveAggressiveForecheck {
				compatibility += 1
			}
		case util.TwoWay:
			// Two-Way forwards get bonus in versatility systems, centers especially
			if position == util.Center && (offSystem == Offensive122Forecheck || defSystem == DefensiveLeftWingLock) {
				compatibility += 1
			}
		}

	case util.Defender:
		// Valid defenseman archetypes: Enforcer, Offensive, Defensive, Two-Way
		switch archetype {
		case util.Grinder, util.Playmaker, util.Power, util.Sniper:
			// These archetypes don't exist for defensemen - this should never happen
			compatibility -= 10 // Heavy penalty for invalid combination
		case util.Offensive:
			// Offensive defensemen excel in transition systems
			if offSystem == OffensiveQuickTransition {
				compatibility += 2 // Strong bonus for perfect fit
			}
		case util.Defensive:
			// Defensive defensemen excel in defensive systems
			if defSystem == DefensiveManToMan || defSystem == DefensiveCollapsing {
				compatibility += 2 // Strong bonus for perfect fit
			}
		}
	}

	return compatibility
}

// GetSystemAttributeModifier calculates attribute modifier based on system compatibility
// Returns a multiplier that should be applied to base attributes
func GetSystemAttributeModifier(offSystem OffensiveSystemType, defSystem DefensiveSystemType, archetype string, position string) float64 {
	compatibility := GetSystemCompatibility(offSystem, defSystem, archetype, position)

	// Convert compatibility score to attribute modifier
	switch {
	case compatibility >= 8:
		return 1.10 // +10% to attributes (excellent fit)
	case compatibility >= 5:
		return 1.05 // +5% to attributes (good fit)
	case compatibility >= 3:
		return 1.03 // +3% to attributes (decent fit)
	case compatibility <= -8:
		return 0.90 // -10% to attributes (terrible fit)
	case compatibility <= -5:
		return 0.95 // -5% to attributes (poor fit)
	case compatibility <= -3:
		return 0.97 // -3% to attributes (slight mismatch)
	default:
		return 1.00 // No modifier (neutral fit)
	}
}

// GetSystemNames returns human-readable names for systems
func GetOffensiveSystemName(system OffensiveSystemType) string {
	switch system {
	case Offensive122Forecheck:
		return "1-2-2 Forecheck"
	case Offensive212Forecheck:
		return "2-1-2 Forecheck"
	case Offensive113Forecheck:
		return "1-1-3 Forecheck"
	case OffensiveCycleGame:
		return "Cycle Game"
	case OffensiveQuickTransition:
		return "Quick Transition"
	case OffensiveUmbrella:
		return "Umbrella (1-3-1)"
	case OffensiveEastWestMotion:
		return "East-West Motion"
	case OffensiveCrashNet:
		return "Crash the Net"
	default:
		return "Unknown System"
	}
}

func GetDefensiveSystemName(system DefensiveSystemType) string {
	switch system {
	case DefensiveBalanced:
		return "Balanced"
	case DefensiveManToMan:
		return "Man-to-Man"
	case DefensiveZone:
		return "Zone Defense"
	case DefensiveNeutralTrap:
		return "Neutral Zone Trap"
	case DefensiveLeftWingLock:
		return "Left-Wing Lock"
	case DefensiveAggressiveForecheck:
		return "Aggressive Forecheck"
	case DefensiveCollapsing:
		return "Collapsing Defense"
	case DefensiveBox:
		return "Box Defense"
	default:
		return "Unknown System"
	}
}
