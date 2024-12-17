package engine

import "github.com/CalebRose/SimHockey/structs"

func LoadForwardLineStrategyTEST(isHome bool) []LineStrategy {
	return []LineStrategy{
		{
			Players: LoadTeamOffenseTest(isHome)[0],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
		{
			Players: LoadTeamOffenseTest(isHome)[1],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
		{
			Players: LoadTeamOffenseTest(isHome)[2],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
		{
			Players: LoadTeamOffenseTest(isHome)[3],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
	}
}

func LoadDefenderLineStrategyTEST(isHome bool) []LineStrategy {
	return []LineStrategy{
		{
			Players: LoadTeamDefenseTest(isHome)[0],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
		{
			Players: LoadTeamDefenseTest(isHome)[1],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
		{
			Players: LoadTeamDefenseTest(isHome)[2],
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
	}
}

func LoadGoalieStrategyTEST(isHome bool) []LineStrategy {
	goalies := LoadGoalieTest(isHome)
	return []LineStrategy{
		{
			Players: []GamePlayer{goalies[0]},
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
		{
			Players: []GamePlayer{goalies[1]},
			Allocations: structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			},
		},
	}
}

func LoadTeamOffenseTest(isHome bool) [][]GamePlayer {
	teamID := uint(1)
	team := "HOME"
	baseAttr := uint(30)
	if !isHome {
		teamID = 2
		team = "AWAY"
		baseAttr = uint(30)
	}
	if isHome {
		return [][]GamePlayer{
			{
				LoadPlayerTEST("Buddy", "Cheque", "C", team, 1, teamID, baseAttr),
				LoadPlayerTEST("EJ", "Faust", "F", team, 2, teamID, baseAttr),
				LoadPlayerTEST("Dmitri", "Petrovich", "F", team, 3, teamID, baseAttr),
			},
			{
				LoadPlayerTEST("Marky", "DuBois", "C", team, 4, teamID, baseAttr),
				LoadPlayerTEST("Reese", "Worthington", "F", team, 5, teamID, baseAttr),
				LoadPlayerTEST("Jocinda", "Smith", "F", team, 6, teamID, baseAttr),
			},
			{
				LoadPlayerTEST("Blake", "Wheeler", "C", team, 7, teamID, baseAttr),
				LoadPlayerTEST("Toner", "Wild", "F", team, 8, teamID, baseAttr),
				LoadPlayerTEST("Maximum", "Justice", "F", team, 9, teamID, baseAttr),
			},
			{
				LoadPlayerTEST("Cookie", "Dawg", "C", team, 10, teamID, baseAttr),
				LoadPlayerTEST("Shawn", "O'Conner", "F", team, 11, teamID, baseAttr),
				LoadPlayerTEST("Seamus", "O'Conner", "F", team, 12, teamID, baseAttr),
			},
		}
	}
	return [][]GamePlayer{
		{
			LoadPlayerTEST("David", "Ross", "C", team, 13, teamID, baseAttr),
			LoadPlayerTEST("Bread", "Man", "F", team, 14, teamID, baseAttr),
			LoadPlayerTEST("Rocket", "Can", "F", team, 15, teamID, baseAttr),
		},
		{
			LoadPlayerTEST("Worst", "Admin", "C", team, 16, teamID, baseAttr),
			LoadPlayerTEST("Jake", "Kirby", "F", team, 17, teamID, baseAttr),
			LoadPlayerTEST("Jacob", "Toth", "F", team, 18, teamID, baseAttr),
		},
		{
			LoadPlayerTEST("Dan", "Molenaar", "C", team, 19, teamID, baseAttr),
			LoadPlayerTEST("Bryce", "Langley", "F", team, 20, teamID, baseAttr),
			LoadPlayerTEST("Connor", "Ryan", "F", team, 21, teamID, baseAttr),
		},
		{
			LoadPlayerTEST("Alex", "Moran", "C", team, 22, teamID, baseAttr),
			LoadPlayerTEST("Tall", "Person", "F", team, 23, teamID, baseAttr),
			LoadPlayerTEST("Average", "Person", "F", team, 24, teamID, baseAttr),
		},
	}
}

func LoadTeamDefenseTest(isHome bool) [][]GamePlayer {
	teamID := uint(1)
	team := "HOME"
	baseAttr := uint(30)
	if !isHome {
		teamID = 2
		team = "AWAY"
		baseAttr = uint(30)
	}
	if isHome {
		return [][]GamePlayer{
			{
				LoadPlayerTEST("Jacob", "Trouba", "D", team, 25, teamID, baseAttr),
				LoadPlayerTEST("Erik", "Johnson", "D", team, 26, teamID, baseAttr),
			},
			{
				LoadPlayerTEST("Patrick", "Johnson", "D", team, 27, teamID, baseAttr),
				LoadPlayerTEST("Sponge", "Bob", "D", team, 28, teamID, baseAttr),
			},
			{
				LoadPlayerTEST("Toasted", "Bread", "D", team, 29, teamID, baseAttr),
				LoadPlayerTEST("Whole Wheat", "Bread", "D", team, 30, teamID, baseAttr),
			},
		}
	}
	return [][]GamePlayer{
		{
			LoadPlayerTEST("Tee", "Sweezy", "D", team, 31, teamID, baseAttr),
			LoadPlayerTEST("Doug", "Dimmadome", "D", team, 32, teamID, baseAttr),
		},
		{
			LoadPlayerTEST("Toucan", "Soda", "D", team, 33, teamID, baseAttr),
			LoadPlayerTEST("Kyle", "Greene", "D", team, 34, teamID, baseAttr),
		},
		{
			LoadPlayerTEST("Tin", "Missile", "D", team, 35, teamID, baseAttr),
			LoadPlayerTEST("Missile", "Jar", "D", team, 36, teamID, baseAttr),
		},
	}
}

func LoadGoalieTest(isHome bool) []GamePlayer {
	teamID := uint(1)
	team := "HOME"
	baseAttr := uint(30)
	if !isHome {
		teamID = 2
		team = "AWAY"
		baseAttr = uint(30)
	}
	if isHome {
		return []GamePlayer{
			LoadPlayerTEST("Patrick", "Roy", "G", team, 37, teamID, baseAttr),
			LoadPlayerTEST("Martin", "Brodeur", "G", team, 38, teamID, baseAttr),
		}
	}
	return []GamePlayer{
		LoadPlayerTEST("Kenny", "Kawaguchi", "G", team, 39, teamID, baseAttr),
		LoadPlayerTEST("Ajon", "Rodney", "G", team, 40, teamID, baseAttr)}
}

func LoadPlayerTEST(firstname, lastname, position, team string, ID, teamID, baseAttr uint) GamePlayer {
	primaryAttr := baseAttr
	secondaryAttr := baseAttr - 4
	forwardAttr := 0
	defAttr := 0
	goalieAttr := 0
	physAttr := primaryAttr
	if position == "F" || position == "C" {
		forwardAttr = int(primaryAttr)
		defAttr = int(secondaryAttr)
		goalieAttr = int(secondaryAttr) - 1
	} else if position == "D" {
		defAttr = int(primaryAttr)
		forwardAttr = int(secondaryAttr)
		goalieAttr = int(secondaryAttr) - 1
	} else {
		goalieAttr = int(primaryAttr)
		defAttr = int(secondaryAttr)
		forwardAttr = int(secondaryAttr) - 1
	}
	baseStats := structs.BasePlayer{
		FirstName:            firstname,
		LastName:             lastname,
		Position:             position,
		TeamID:               uint16(teamID),
		Team:                 team,
		WristShotAccuracy:    uint8(forwardAttr),
		WristShotPower:       uint8(forwardAttr),
		SlapshotAccuracy:     uint8(forwardAttr),
		SlapshotPower:        uint8(forwardAttr),
		OneTimer:             uint8(forwardAttr),
		Strength:             uint8(physAttr),
		Faceoffs:             uint8(forwardAttr),
		Agility:              uint8(physAttr),
		Goalkeeping:          uint8(goalieAttr),
		GoalieVision:         uint8(goalieAttr),
		GoalieReboundControl: uint8(goalieAttr),
		ShotBlocking:         uint8(goalieAttr),
		Passing:              uint8(physAttr),
		PuckHandling:         uint8(defAttr),
		BodyChecking:         uint8(defAttr),
		StickChecking:        uint8(defAttr),
		Discipline:           uint8(baseAttr),
		Aggression:           uint8(baseAttr),
		Stamina:              50,
		InjuryRating:         uint8(baseAttr),
	}

	gamePlayer := GamePlayer{
		ID:             ID,
		BasePlayer:     baseStats,
		CurrentStamina: 50,
	}
	gamePlayer.CalculateModifiers()

	return gamePlayer
}
