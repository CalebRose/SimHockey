package engine

import util "github.com/CalebRose/SimHockey/_util"

func CalculatePenaltyChance() bool {
	chance := 0.1
	return util.GenerateFloatFromRange(1, 100) <= chance
}

func ApplyPenalty(gs *GameState, penalty Penalty, player *GamePlayer) {
	switch penalty.Severity {
	case MinorPenalty:
		HandlePenaltyLogic(gs, penalty, 120, player) // 2 minutes
	case MajorPenalty:
		HandlePenaltyLogic(gs, penalty, 300, player) // 5 minutes
	case GameMisconduct:
		RemovePlayerFromGame(gs, player)
	case MatchPenalty:
		RemovePlayerFromGame(gs, player)
	}

	// Reset game with faceoff
	gs.SetFaceoffOnCenterIce(true)
}

func HandlePenaltyLogic(gs *GameState, penalty Penalty, duration int, player *GamePlayer) {
	// Put Player into Penalty Box
	player.GoToPenaltyBox(penalty.PenaltyType > 2)

	// Set Gamestate duration and power play team
	powerPlayTeamID := gs.HomeTeamID
	if player.TeamID == uint16(gs.HomeTeamID) {
		powerPlayTeamID = gs.AwayTeamID
		gs.HomeStrategy.ActivatePowerPlayer(player.ID, player.Position)
	} else {
		gs.AwayStrategy.ActivatePowerPlayer(player.ID, player.Position)
	}

	gs.SetPowerPlay(duration, int(powerPlayTeamID), penalty)
}

func HandleFight(gs *GameState) {

}

func RemovePlayerFromGame(gs *GameState, player *GamePlayer) {
	isHome := player.TeamID == uint16(gs.HomeTeamID)
	hasSubstitutablePlayers := true
	benchPlayers := 0
	// Check for benchable players
	if isHome {
		benchPlayers = len(gs.HomeStrategy.BenchPlayers)
	} else {
		benchPlayers = len(gs.AwayStrategy.BenchPlayers)
	}
	hasSubstitutablePlayers = benchPlayers > 0

	if hasSubstitutablePlayers {
		gs.RemovePlayerFromLine(isHome, player.ID)
	}
}

func SelectPenalty(player *GamePlayer, typeID uint, context string) Penalty {
	selectablePenalties := GetPenaltiesByID(player, typeID, context)
	totalWeight := 0.0
	for _, p := range selectablePenalties {
		totalWeight += p.Weight
	}
	curr := 0.0
	chance := util.GenerateFloatFromRange(0, totalWeight)
	for _, p := range selectablePenalties {
		if chance <= curr {
			p.ApplyPlayerInfo(player.ID, player.Position)
			return p
		}
		curr += p.Weight
	}
	// If no penalty for somne reason is selected, return an empty one
	return Penalty{}
}

func GetPenaltiesByID(player *GamePlayer, id uint, context string) []Penalty {
	penaltyList := GetAllPenalties()
	if id == 2 {
		return penaltyList
	}
	filteredList := []Penalty{}
	for _, p := range penaltyList {
		validPenalty := player.Aggression >= p.AggressionReq && player.Discipline <= p.DisciplineReq
		if context == General && validPenalty {
			filteredList = append(filteredList, p)
		} else if p.PenaltyType <= id && (p.Context == context) && validPenalty {
			filteredList = append(filteredList, p)
		}
	}

	return filteredList
}

func GetAllPenalties() []Penalty {
	return []Penalty{
		GetPenalty(1, 0, "Aggressor Penalty", Fight, MatchPenalty, 0.001, false, 80, 40),
		GetPenalty(2, 0, "Attempt to Injure", BodyCheck, MatchPenalty, 0.001, false, 90, 30),
		GetPenalty(3, 0, "Biting", Fight, MajorPenalty, 0.1, false, 85, 20),
		GetPenalty(4, 0, "Boarding", BodyCheck, MinorPenalty, 1, false, 70, 50),
		GetPenalty(5, 0, "Boarding", BodyCheck, MajorPenalty, 0.1, false, 75, 40),
		GetPenalty(6, 0, "Stabbing", General, GameMisconduct, 0.01, false, 95, 10),
		GetPenalty(7, 0, "Charging", BodyCheck, MinorPenalty, 1, false, 75, 40),
		GetPenalty(8, 0, "Charging", BodyCheck, MajorPenalty, 0.1, false, 85, 30),
		GetPenalty(9, 0, "Checking from Behind", BodyCheck, MinorPenalty, 1, false, 80, 50),
		GetPenalty(10, 0, "Checking from Behind", BodyCheck, MajorPenalty, 0.1, false, 85, 40),
		GetPenalty(11, 0, "Clipping", BodyCheck, MinorPenalty, 1, false, 60, 60),
		GetPenalty(12, 0, "Clipping", BodyCheck, MajorPenalty, 0.1, false, 65, 50),
		GetPenalty(13, 0, "Cross Checking", StickCheck, MinorPenalty, 1, false, 75, 50),
		GetPenalty(14, 0, "Cross Checking", StickCheck, MajorPenalty, 0.1, false, 80, 40),
		GetPenalty(15, 0, "Delay of Game", General, MinorPenalty, 1, false, 50, 80),
		GetPenalty(16, 0, "Diving", General, MinorPenalty, 1, false, 30, 70),
		GetPenalty(17, 0, "Elbowing", BodyCheck, MinorPenalty, 1, false, 65, 60),
		GetPenalty(18, 0, "Elbowing", BodyCheck, MajorPenalty, 0.1, false, 70, 50),
		GetPenalty(19, 0, "Eye-Gouging", Fight, MajorPenalty, 0.1, false, 90, 20),
		GetPenalty(20, 0, "Fighting", Fight, MajorPenalty, 0.5, true, 80, 30),
		GetPenalty(21, 2, "Goaltender Interference", General, MinorPenalty, 1.5, false, 55, 70),
		GetPenalty(22, 0, "Headbutting", Fight, MatchPenalty, 0.001, false, 95, 20),
		GetPenalty(23, 0, "High-sticking", StickCheck, MinorPenalty, 1, false, 60, 70),
		GetPenalty(24, 0, "High-sticking", StickCheck, MajorPenalty, 0.1, false, 65, 60),
		GetPenalty(25, 0, "Holding", BodyCheck, MinorPenalty, 1, false, 50, 80),
		GetPenalty(26, 0, "Hooking", StickCheck, MinorPenalty, 1, false, 50, 80),
		GetPenalty(27, 0, "Hooking", StickCheck, MajorPenalty, 0.1, false, 55, 70),
		GetPenalty(28, 0, "Kicking", BodyCheck, MinorPenalty, 1, false, 60, 60),
		GetPenalty(29, 0, "Kicking", BodyCheck, MajorPenalty, 0.1, false, 75, 50),
		GetPenalty(30, 0, "Kneeing", BodyCheck, MinorPenalty, 1, false, 55, 70),
		GetPenalty(31, 0, "Kneeing", BodyCheck, MajorPenalty, 0.1, false, 65, 60),
		GetPenalty(32, 0, "Roughing", BodyCheck, MinorPenalty, 1, false, 70, 50),
		GetPenalty(33, 0, "Roughing", BodyCheck, MajorPenalty, 0.1, false, 85, 50),
		GetPenalty(34, 0, "Slashing", StickCheck, MinorPenalty, 1, false, 60, 60),
		GetPenalty(35, 0, "Slashing", StickCheck, MajorPenalty, 0.1, false, 75, 50),
		GetPenalty(36, 0, "Slew footing", BodyCheck, MinorPenalty, 1, false, 60, 50),
		GetPenalty(37, 0, "Slew footing", BodyCheck, MajorPenalty, 0.1, false, 75, 50),
		GetPenalty(38, 1, "Throwing the stick", StickCheck, MinorPenalty, 1, false, 50, 70),
		GetPenalty(39, 0, "Too many men on the ice", General, MinorPenalty, 1, false, 10, 90),
		GetPenalty(40, 0, "Tripping", StickCheck, MinorPenalty, 1, false, 50, 80),
		GetPenalty(41, 0, "Tripping", StickCheck, MajorPenalty, 0.1, false, 55, 70),
		GetPenalty(42, 0, "Unsportsmanlike conduct", General, MinorPenalty, 1, false, 85, 50),
	}
}

func GetPenalty(id uint, penaltyType uint, name, context, sev string, weight float64, isFight bool, aggReq, disReq uint8) Penalty {
	return Penalty{
		PenaltyID:     id,
		PenaltyName:   name,
		PenaltyType:   penaltyType,
		Severity:      sev,
		Weight:        weight,
		IsFight:       isFight,
		AggressionReq: aggReq,
		DisciplineReq: disReq,
		Context:       context,
	}
}

func GetSeverityID(severity string) uint8 {
	if severity == MinorPenalty {
		return 1
	}
	if severity == MajorPenalty {
		return 2
	}
	if severity == GameMisconduct {
		return 3
	}
	// Match Penalty
	return 4
}
