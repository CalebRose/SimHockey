package engine

import (
	"github.com/CalebRose/SimHockey/structs"
)

type GameState struct {
	GameID                uint
	WeekID                uint
	HomeTeamID            uint
	HomeTeam              string
	HomeTeamCoach         string
	HomeStrategy          GamePlaybook
	HomeTeamStats         TeamStatDTO
	AwayTeamID            uint
	AwayTeam              string
	AwayTeamCoach         string
	AwayStrategy          GamePlaybook
	AwayTeamStats         TeamStatDTO
	FaceoffOnCenterIce    bool
	Period                uint8
	TimeOnClock           uint16
	SecondsConsumed       uint16
	IsOvertime            bool
	IsOvertimeShootout    bool
	HomeTeamScore         uint8
	AwayTeamScore         uint8
	HomeTeamShootoutScore uint8
	AwayTeamShootoutScore uint8
	HomeTeamWin           bool
	AwayTeamWin           bool
	TieGame               bool
	GameComplete          bool
	IsPowerPlay           bool
	PowerPlayTeamID       uint
	PowerPlayIDCount      uint
	ActivePowerPlays      []PowerPlayState
	Zones                 []string // Home Goal, Home Zone, Neutral Zone, Away Zone, Away Goal
	PuckLocation          string
	PuckCarrier           *GamePlayer
	AssistingPlayer       *GamePlayer
	Momentum              float64
	PossessingTeam        uint // Use the ID of the team
	IsCollegeGame         bool
	IsPlayoffGame         bool
	Collector             structs.PbPCollector
}

type PowerPlayState struct {
	PowerPlayID     uint
	PowerPlayTeamID uint
	PenaltyInfo     Penalty
	PowerPlayTime   uint16
}

// Todo -- Change powerplay structure to a list of Power Play Objects. Iterate over list and decrement over time.
// If a score is made, check list for who the powerplay team is and if they scored. Then return penalized player
// and remove powerplay item from list

func (gs *GameState) RecordPlay(play structs.PbP) {
	gs.Collector.AppendPlay(play)
}

func (gs *GameState) SetTime(isNewPeriod, isOvertime bool) {
	if isNewPeriod && !isOvertime {
		handleNewPeriod(gs)
	} else {
		reduceTimeOnClock(gs)
		handleLineups(gs)
	}
	gs.SecondsConsumed = 1
}

func handleNewPeriod(gs *GameState) {
	// Increment period
	gs.Period++

	// Determine if overtime or shootout
	if gs.Period > 3 && gs.HomeTeamScore == gs.AwayTeamScore {
		gs.IsOvertime = true
	} else if gs.Period > 3 && gs.HomeTeamScore != gs.AwayTeamScore {
		gs.GameComplete = true
	}
	if gs.Period > 4 && gs.HomeTeamScore == gs.AwayTeamScore && !gs.IsPlayoffGame {
		gs.IsOvertimeShootout = true
	}

	// Set time on clock
	if gs.IsOvertime && !gs.IsPlayoffGame {
		gs.TimeOnClock = OvertimePeriodTime
	} else {
		gs.TimeOnClock = RegularPeriodTime
	}
}

func reduceTimeOnClock(gs *GameState) {
	// Calculate remaining time
	gs.TimeOnClock -= gs.SecondsConsumed
	if len(gs.ActivePowerPlays) > 0 {
		for _, pp := range gs.ActivePowerPlays {
			if pp.PowerPlayTime > 0 {
				pp.PowerPlayTime -= gs.SecondsConsumed
				if pp.PowerPlayTime > MaxTimeOnClock {
					pp.PowerPlayTime = 0
					gs.TurnOffPowerPlay(pp)
				}
			}
		}
	}

	// Ensure time doesn't overflow
	if gs.TimeOnClock > MaxTimeOnClock {
		gs.TimeOnClock = 0
	}
}

func handleLineups(gs *GameState) {
	gs.HomeStrategy.HandleLineups(int(gs.SecondsConsumed), gs.IsPowerPlay)
	gs.AwayStrategy.HandleLineups(int(gs.SecondsConsumed), gs.IsPowerPlay)
}

func (gs *GameState) SetSecondsConsumed(seconds uint16) {
	gs.SecondsConsumed += seconds
}

func (gs *GameState) SetPowerPlay(duration, teamID int, penalty Penalty) {
	gs.ResetMomentum()
	id := gs.PowerPlayIDCount
	gs.PowerPlayIDCount++
	pp := PowerPlayState{
		PowerPlayID:     id,
		PowerPlayTime:   uint16(duration),
		PowerPlayTeamID: uint(teamID),
		PenaltyInfo:     penalty,
	}
	gs.ActivePowerPlays = append(gs.ActivePowerPlays, pp)
	gs.IsPowerPlay = true
	if gs.PowerPlayTeamID == 0 {
		gs.PowerPlayTeamID = uint(teamID)
	} else if gs.PowerPlayTeamID > 0 {
		// Even power play
		gs.PowerPlayTeamID = 0
	}
}

func (gs *GameState) TurnOffPowerPlay(powerPlay PowerPlayState) {
	homePowerPlay := gs.HomeTeamID == powerPlay.PowerPlayTeamID
	if homePowerPlay {
		// Yeah I might need to just setup the new structure
		gs.AwayStrategy.ReturnPlayerFromPowerPlay(powerPlay.PenaltyInfo.PenalizedPlayerID, powerPlay.PenaltyInfo.PenalizedPlayerPosition)
	} else {
		gs.HomeStrategy.ReturnPlayerFromPowerPlay(powerPlay.PenaltyInfo.PenalizedPlayerID, powerPlay.PenaltyInfo.PenalizedPlayerPosition)
	}
	gs.FilterActivePowerPlayList(powerPlay.PowerPlayID)
	gs.IsPowerPlay = len(gs.ActivePowerPlays) > 0
}

func (gs *GameState) FilterActivePowerPlayList(powerPlayID uint) {
	activeList := []PowerPlayState{}
	for _, pp := range gs.ActivePowerPlays {
		if pp.PowerPlayID != powerPlayID {
			activeList = append(activeList, pp)
		}
	}
	gs.ActivePowerPlays = activeList
}

func (gs *GameState) SetNewZone(zone string) {
	gs.PuckLocation = zone
}

func (gs *GameState) SetFaceoffOnCenterIce(check bool) {
	gs.FaceoffOnCenterIce = check
}

func (gs *GameState) SetPuckBearer(player *GamePlayer) {
	if gs.PuckCarrier != nil {
		if gs.PuckCarrier.ID > 0 && gs.PuckCarrier.TeamID != player.TeamID {
			gs.ResetMomentum()
			gs.AssistingPlayer = &GamePlayer{}
		} else {
			gs.Momentum += 0.125
			gs.AssistingPlayer = gs.PuckCarrier
		}
		gs.PuckCarrier = player
		gs.PossessingTeam = uint(player.TeamID)
	} else {
		gs.PuckCarrier = &GamePlayer{}
	}
}

func (gs *GameState) TriggerBreakaway() {
	gs.Momentum += 0.275
}

func (gs *GameState) ResetMomentum() {
	gs.Momentum = 0
}

func (gs *GameState) GetCenter(isHome bool) *GamePlayer {
	var currentLineup []*GamePlayer
	if isHome {
		s := gs.HomeStrategy
		idx := s.CurrentForwards
		currentLineup = s.Forwards[idx].Players
	} else {
		s := gs.AwayStrategy
		idx := s.CurrentForwards
		currentLineup = s.Forwards[idx].Players
	}

	player := GetGamePlayerByPosition(currentLineup, "C")
	if player.ID == 0 || player.IsOut {
		player = GetGamePlayerByPosition(currentLineup, "F")
	}
	if player.ID == 0 || player.IsOut {
		player = GetGamePlayerByPosition(currentLineup, "D")
	}
	return player
}

func (gs *GameState) IncrementShootoutScore(isHome bool) {
	if isHome {
		gs.HomeTeamShootoutScore += 1
	} else {
		gs.AwayTeamShootoutScore += 1
	}
}

func (gs *GameState) CalculateWinner() {
	gs.GameComplete = true
	if (gs.HomeTeamScore > gs.AwayTeamScore) || (gs.HomeTeamShootoutScore > gs.AwayTeamShootoutScore) {
		gs.HomeTeamWin = true
	} else if (gs.AwayTeamScore > gs.HomeTeamScore) || gs.AwayTeamShootoutScore > gs.HomeTeamShootoutScore {
		gs.AwayTeamWin = true
	} else {
		gs.TieGame = true
	}
}

func (gs *GameState) IncrementScore(isHome bool) {
	if isHome {
		gs.HomeTeamScore += 1
		gs.HomeTeamStats.AddShot(true, !gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID), gs.Period > 3, gs.Period)
		gs.AwayTeamStats.AddShotAgainst(true)
	} else {
		gs.AwayTeamScore += 1
		gs.AwayTeamStats.AddShot(true, !gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID), gs.Period > 3, gs.Period)
		gs.HomeTeamStats.AddShotAgainst(true)
	}
	gs.PuckCarrier.AddShot(true, !gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID), gs.Period > 3)
	if gs.AssistingPlayer.ID > 0 {
		gs.AssistingPlayer.Stats.AddAssist(!gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID))
	}
	gs.HomeStrategy.HandlePlusMinus(isHome, gs.PuckCarrier.ID, gs.AssistingPlayer.ID)
	gs.AwayStrategy.HandlePlusMinus(!isHome, gs.PuckCarrier.ID, gs.AssistingPlayer.ID)
	gs.SetFaceoffOnCenterIce(true)
	gs.SetNewZone(NeutralZone)
	gs.SetPuckBearer(&GamePlayer{})
	if gs.IsPowerPlay {
		for _, pp := range gs.ActivePowerPlays {
			if ((isHome && pp.PowerPlayTeamID == gs.HomeTeamID) ||
				(!isHome && pp.PowerPlayTeamID == gs.AwayTeamID)) &&
				pp.PenaltyInfo.PenaltyType == 1 {
				gs.TurnOffPowerPlay(pp)
				break
			}
		}
	}
}

func (gs *GameState) AddShots(isHome bool) {
	if isHome {
		gs.HomeTeamStats.AddShot(false, false, false, false, false, 0)
		gs.AwayTeamStats.AddShotAgainst(false)
	} else {
		gs.AwayTeamStats.AddShot(false, false, false, false, false, 0)
		gs.HomeTeamStats.AddShotAgainst(false)
	}
}

func (gs *GameState) GetPlaybook(isHome bool) GamePlaybook {
	if isHome {
		return gs.HomeStrategy
	}
	return gs.AwayStrategy
}

func (gs *GameState) GetLineStrategy(isHome bool, lineUpType int) LineStrategy {
	pb := GamePlaybook{}
	if isHome {
		pb = gs.HomeStrategy
	} else {
		pb = gs.AwayStrategy
	}
	if lineUpType == 1 {
		return pb.Forwards[pb.CurrentForwards]
	}
	if lineUpType == 2 {
		return pb.Defenders[pb.CurrentDefenders]
	}
	return pb.Goalies[pb.CurrentGoalie]
}

func GetGamePlayerByPosition(currentLineup []*GamePlayer, pos string) *GamePlayer {
	for _, p := range currentLineup {
		if p.Position != pos || p.IsOut {
			continue
		}
		return p
	}
	return &GamePlayer{}
}

func (gs *GameState) GetLineupType(isHome bool) int {
	if isHome {
		return gs.HomeStrategy.GetLineupType()
	}
	return gs.AwayStrategy.GetLineupType()
}

func (gs *GameState) RemovePlayerFromLine(isHome bool, playerID uint) {
	if isHome {
		gs.HomeStrategy.filterOutPlayer(playerID)
	} else {
		gs.AwayStrategy.filterOutPlayer(playerID)
	}
}

type GamePlaybook struct {
	Forwards                    []LineStrategy
	Defenders                   []LineStrategy
	CurrentForwards             int
	CurrentDefenders            int
	Goalies                     []LineStrategy
	CurrentGoalie               int
	MinForwardStaminaThreshold  int
	MinDefenderStaminaThreshold int
	BenchPlayers                []*GamePlayer
	CenterOut                   bool
	Forward1Out                 bool
	Forward2Out                 bool
	Defender1Out                bool
	Defender2Out                bool
	ShootoutLineUp              structs.ShootoutPlayerIDs
	RosterMap                   map[uint]*GamePlayer
}

func (gp *GamePlaybook) ActivatePowerPlayer(playerID uint, position string) {
	gp.HandlePositionToggle(playerID, position, true)
}

func (gp *GamePlaybook) ReturnPlayerFromPowerPlay(powerPlayID uint, powerPlayPostion string) {
	// Iterate through forwards and defenders to find the penalized player and return them in play
	if powerPlayPostion == Forward {
		gp.Forwards[gp.CurrentForwards].ReturnPlayerFromPowerPlay(powerPlayID)
	} else if powerPlayPostion == Defender {
		gp.Defenders[gp.CurrentDefenders].ReturnPlayerFromPowerPlay(powerPlayID)
	}
	gp.HandlePositionToggle(powerPlayID, powerPlayPostion, false)
}

func (gp *GamePlaybook) HandlePositionToggle(playerID uint, position string, isOut bool) {
	enum := gp.getPlayerPositionEnum(playerID, position)
	if enum == 1 {
		gp.CenterOut = isOut
	} else if enum == 2 {
		gp.Forward1Out = isOut
	} else if enum == 3 {
		gp.Forward2Out = isOut
	} else if enum == 4 {
		gp.Defender1Out = isOut
	} else if enum == 5 {
		gp.Defender2Out = isOut
	}
}

func (gp *GamePlaybook) HandlePlusMinus(isScore bool, puckBearerID, assistingPlayerID uint) {
	forwards := gp.Forwards[gp.CurrentForwards]
	for _, p := range forwards.Players {
		if p.ID == puckBearerID || p.ID == assistingPlayerID {
			continue
		}
		p.Stats.AdjustPlusMinus(isScore)
	}
	defenders := gp.Defenders[gp.CurrentDefenders]
	for _, p := range defenders.Players {
		if p.ID == puckBearerID || p.ID == assistingPlayerID {
			continue
		}
		p.Stats.AdjustPlusMinus(isScore)
	}
}

func (gp *GamePlaybook) GetLineupType() int {
	currentForwards := gp.Forwards[gp.CurrentForwards]
	currentDefenders := gp.Defenders[gp.CurrentDefenders]
	if len(currentForwards.Players) < 3 {
		return 1
	}
	if len(currentDefenders.Players) < 2 {
		return 2
	}
	return 0
}

func (gp *GamePlaybook) filterOutPlayer(playerID uint) {
	isForward := false

	// Handle Forward Replacement
	forwardIdx := gp.CurrentForwards
	if playerIDInLineup(playerID, gp.Forwards[forwardIdx].Players) {
		isForward = true
		gp.Forwards[forwardIdx].Players = gp.handleLineReplacement(
			gp.Forwards[forwardIdx].Players, playerID, 3, 1)
	}

	// If the player was a forward, we don't need to process defenders
	if isForward {
		return
	}

	// Handle Defender Replacement
	defenderIdx := gp.CurrentDefenders
	if playerIDInLineup(playerID, gp.Defenders[defenderIdx].Players) {
		gp.Defenders[defenderIdx].Players = gp.handleLineReplacement(
			gp.Defenders[defenderIdx].Players, playerID, 2, 2)
	}
}

func (gp *GamePlaybook) getPlayerPositionEnum(playerID uint, position string) uint {
	// Return an enum (1 == C, 2 == F1, 3 == F2, 4 == D1, 5==D2) that will indicate which position to shutoff during a power play.
	if position == Center {
		return 1
	}
	if position == Forward {
		currentForwards := gp.Forwards[gp.CurrentForwards]
		for idx, p := range currentForwards.Players {
			if p.ID == playerID {
				return 1 + uint(idx)
			}
		}
	} else if position == Defender {
		currentDefenders := gp.Defenders[gp.CurrentDefenders]
		for idx, p := range currentDefenders.Players {
			if p.ID == playerID {
				return 4 + uint(idx)
			}
		}
	}
	return 0
}

func (gp *GamePlaybook) handleLineReplacement(players []*GamePlayer, playerID uint, requiredCount, lineType uint) []*GamePlayer {
	filteredPlayers, queue := filterOutPlayerFromLineup(players, playerID, lineType)

	for len(filteredPlayers) < int(requiredCount) {
		var replacement *GamePlayer

		// Check if a substitution ID exists for the current queue player
		if queue.SubstitutionID > 0 {
			nextLine := gp.getNextLine(lineType)
			replacement = queue.Player
			queue = GetPlayerFromLine(queue.SubstitutionID, nextLine)
		} else {
			// Use a player from the bench
			replacement, gp.BenchPlayers = popPlayerFromBench(gp.BenchPlayers)
		}

		// Add replacement player to the lineup
		filteredPlayers = append(filteredPlayers, replacement)
	}

	return filteredPlayers
}

func (gp *GamePlaybook) InitializeStamina() {
	for idx := range gp.Forwards {
		gp.Forwards[idx].InitializeBoostedStamina(false)
	}

	for idx := range gp.Defenders {
		gp.Defenders[idx].InitializeBoostedStamina(false)
	}

	for idx := range gp.Goalies {
		gp.Goalies[idx].InitializeBoostedStamina(true)
	}
}

func (gp *GamePlaybook) HandleLineups(secondsConsumed int, isPowerPlay bool) {
	for idx := range gp.Forwards {
		if idx == gp.CurrentForwards {
			gp.Forwards[idx].DecrementStamina()
			gp.Forwards[idx].HandleTimeOnIce(secondsConsumed)
		} else {
			gp.Forwards[idx].RecoverStaminaOffRink()
		}
	}

	for idx := range gp.Defenders {
		if idx == gp.CurrentDefenders {
			gp.Defenders[idx].DecrementStamina()
			gp.Defenders[idx].HandleTimeOnIce(secondsConsumed)
		} else {
			gp.Defenders[idx].RecoverStaminaOffRink()
		}
	}

	for idx := range gp.Goalies {
		if idx == gp.CurrentGoalie {
			gp.Goalies[idx].DecrementStamina()
			gp.Goalies[idx].HandleTimeOnIce(secondsConsumed)
		} else {
			gp.Goalies[idx].RecoverStaminaOffRink()
		}
	}

	gp.CheckAndRotateLineup()
}

func (gp *GamePlaybook) CheckAndRotateLineup() {
	if gp.Forwards[gp.CurrentForwards].CurrentStamina < gp.MinForwardStaminaThreshold {
		gp.CurrentForwards = (gp.CurrentForwards + 1) % len(gp.Forwards)
	}
	if gp.Defenders[gp.CurrentDefenders].CurrentStamina < gp.MinDefenderStaminaThreshold {
		gp.CurrentDefenders = (gp.CurrentDefenders + 1) % len(gp.Defenders)
	}
	// Goalies do not swap unless overtime
	if float64(gp.Goalies[gp.CurrentGoalie].CurrentStamina) < float64(gp.Goalies[gp.CurrentGoalie].TotalStamina)*(0.1) {
		gp.CurrentGoalie = (gp.CurrentGoalie + 1) & len(gp.Goalies)
	}
}

func (gp *GamePlaybook) getNextLine(lineType uint) []*GamePlayer {
	if lineType == 1 {
		for i := gp.CurrentForwards + 1; i < len(gp.Forwards); i++ {
			return gp.Forwards[i].Players
		}
	} else if lineType == 2 {
		for i := gp.CurrentDefenders + 1; i < len(gp.Defenders); i++ {
			return gp.Defenders[i].Players
		}
	}
	return []*GamePlayer{}
}

type LineStrategy struct {
	structs.Allocations
	Players        []*GamePlayer
	CenterID       uint
	Forward1ID     uint
	Forward2ID     uint
	Defender1ID    uint
	Defender2ID    uint
	TotalStamina   int
	CurrentStamina int
	Threshold      int
}

func (ls *LineStrategy) ReturnPlayerFromPowerPlay(playerID uint) {
	for _, p := range ls.Players {
		if p.ID == playerID {
			p.ReturnToPlay()
		}
	}
}

func (ls *LineStrategy) SetNewLineup(players []*GamePlayer) {
	ls.Players = players
}

func (ls *LineStrategy) InitializeBoostedStamina(isGoalie bool) {
	staminaBoostFactor := 1.75 // Adjust as needed
	if isGoalie {
		staminaBoostFactor = 75
	}
	ls.TotalStamina = 0
	mode := 0.0
	if isGoalie {
		ls.CurrentStamina = int(float64(ls.Players[0].CurrentStamina) * staminaBoostFactor)
		return
	}
	for i := range ls.Players {
		initialStamina := int(float64(ls.Players[i].CurrentStamina) * staminaBoostFactor)
		mode += float64(initialStamina)
	}
	avg := mode / float64(len(ls.Players))
	ls.CurrentStamina = int(avg)
	ls.TotalStamina = int(avg)
}

func (ls *LineStrategy) RecoverStaminaOffRink() {
	if ls.CurrentStamina < ls.TotalStamina {
		ls.CurrentStamina += 1
	}
}

func (ls *LineStrategy) DecrementStamina() {
	// Update CurrentStamina for lineup after decrement
	ls.CurrentStamina -= 1
}

func (ls *LineStrategy) HandleTimeOnIce(secondsConsumed int) {
	for _, p := range ls.Players {
		p.Stats.AddTimeOnIce(secondsConsumed, p.IsOut)
	}
}

// GamePlayer -- for
type GamePlayer struct {
	ID uint
	structs.BasePlayer
	CurrentStamina    int
	OneTimerMod       float64
	AgilityMod        float64
	StrengthMod       float64
	LongShotPowerMod  float64
	LongShotAccMod    float64
	CloseShotPowerMod float64
	CloseShotAccMod   float64
	FaceoffMod        float64
	HandlingMod       float64
	PassMod           float64
	StickCheckMod     float64
	BodyCheckMod      float64
	GoalkeepingMod    float64
	GoalieVisionMod   float64
	ShotblockingMod   float64
	IsOut             bool
	FoulOut           bool
	ForcedOut         bool
	SubstituteID      uint
	Stats             PlayerStatsDTO
}

func (g *GamePlayer) MapFromCollegePlayer(c structs.CollegePlayer) {
	g.BasePlayer = c.BasePlayer
	g.ID = c.ID
	g.IsOut = c.IsInjured || c.IsRedshirting
	g.CalculateModifiers()
}

func (g *GamePlayer) CalculateModifiers() {
	g.OneTimerMod = calculateAttributeModifier(float64(g.OneTimer), ModifierFactor)
	g.AgilityMod = calculateAttributeModifier(float64(g.Agility), ModifierFactor)
	g.StrengthMod = calculateAttributeModifier(float64(g.Strength), ModifierFactor)
	g.LongShotAccMod = calculateAttributeModifier(float64(g.LongShotAccuracy), ModifierFactor)
	g.LongShotPowerMod = calculateAttributeModifier(float64(g.LongShotPower), ModifierFactor)
	g.CloseShotPowerMod = calculateAttributeModifier(float64(g.CloseShotPower), ModifierFactor)
	g.CloseShotAccMod = calculateAttributeModifier(float64(g.CloseShotAccuracy), ModifierFactor)
	g.FaceoffMod = calculateAttributeModifier(float64(g.Faceoffs), ModifierFactor)
	g.HandlingMod = calculateAttributeModifier(float64(g.PuckHandling), ModifierFactor)
	g.PassMod = calculateAttributeModifier(float64(g.Passing), ModifierFactor)
	g.StickCheckMod = calculateAttributeModifier(float64(g.StickChecking), ModifierFactor)
	g.BodyCheckMod = calculateAttributeModifier(float64(g.BodyChecking), ModifierFactor)
	g.GoalkeepingMod = calculateAttributeModifier(float64(g.Goalkeeping), ModifierFactor)
	g.GoalieVisionMod = calculateAttributeModifier(float64(g.GoalieVision), ModifierFactor)
	g.ShotblockingMod = calculateAttributeModifier(float64(g.ShotBlocking), ModifierFactor)
}

func (g *GamePlayer) GoToPenaltyBox(outOfGame bool) {
	g.IsOut = true
	if outOfGame {
		g.ForcedOut = true
	}
}

func (g *GamePlayer) ReturnToPlay() {
	if !g.ForcedOut {
		g.IsOut = false
	}
}

func (p *GamePlayer) AdjustPlusMinus(isScore bool) {
	p.Stats.AdjustPlusMinus(isScore)
}

func (p *GamePlayer) AddShot(isScore, isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool) {
	p.Stats.AddShot(isScore, isEvenStrength, isPowerPlay, isShorthanded, isOvertime)
}

func (p *GamePlayer) AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool) {
	p.Stats.AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime)
}

func (p *GamePlayer) AddGoalAgainst() {
	p.Stats.AddGoalAgainst()
}

func (p *GamePlayer) AddAssist(isEvenStrength, isPowerPlay, isShorthanded bool) {
	p.Stats.AddAssist(isEvenStrength, isPowerPlay, isShorthanded)
}

func (p *GamePlayer) AddTimeOnIce(seconds int, isPenalty bool) {
	p.Stats.AddTimeOnIce(seconds, isPenalty)
}

func (p *GamePlayer) AddPenaltyTime(seconds int) {
	p.Stats.AddPenaltyTime(seconds)
}

func (p *GamePlayer) AddShotAgainst(isScore bool) {
	p.Stats.AddShotAgainst(isScore)
}

func (p *GamePlayer) AddShotBlocked() {
	p.Stats.AddShotBlocked()
}

func (p *GamePlayer) AddDefensiveHit(isBodyCheck bool) {
	p.Stats.AddDefensiveHit(isBodyCheck)
}

func (p *GamePlayer) AddFaceoff(isWin bool) {
	p.Stats.AddFaceoff(isWin)
}

func (p *GamePlayer) AddShutout() {
	p.Stats.AddShutout()
}

func (p *GamePlayer) AddGoalieStat(scoreType int, isOvertimeLoss bool) {
	p.Stats.AddGoalieStat(scoreType, isOvertimeLoss)
}

// Util Structs
// PlayerWeight -- For event checks
type PlayerWeight struct {
	PlayerID uint
	Weight   float64
}

type Penalty struct {
	PenaltyID               uint
	PenaltyName             string
	PenaltyType             uint // 0 == Can occur anywhere, 1 == Defending Zones, 2== Goal Defending Zones
	Severity                string
	Weight                  float64
	IsFight                 bool
	AggressionReq           uint8
	DisciplineReq           uint8
	Context                 string
	PenalizedPlayerID       uint
	PenalizedPlayerPosition string
}

func (p *Penalty) ApplyPlayerInfo(id uint, position string) {
	p.PenalizedPlayerID = id
	p.PenalizedPlayerPosition = position
}

// Line Management Functions and Structs
type RemovalQueue struct {
	PlayerID       uint
	Player         *GamePlayer
	SubstitutionID uint
	LineType       uint
	Idx            uint
}

func filterOutPlayerFromLineup(players []*GamePlayer, playerID uint, lineType uint) ([]*GamePlayer, RemovalQueue) {
	filtered := []*GamePlayer{}
	queue := RemovalQueue{}

	for _, p := range players {
		if p.ID == playerID {
			queue = RemovalQueue{
				PlayerID:       p.ID,
				Player:         p,
				SubstitutionID: p.SubstituteID,
				LineType:       lineType,
			}
			continue
		}
		filtered = append(filtered, p)
	}

	return filtered, queue
}

func GetPlayerFromLine(playerID uint, players []*GamePlayer) RemovalQueue {
	for _, p := range players {
		if p.ID == playerID {
			return RemovalQueue{
				PlayerID:       p.ID,
				Player:         p,
				SubstitutionID: p.SubstituteID,
			}
		}
	}
	return RemovalQueue{}
}

func popPlayerFromBench(bench []*GamePlayer) (*GamePlayer, []*GamePlayer) {
	if len(bench) == 0 {
		return &GamePlayer{}, bench
	}
	return bench[0], bench[1:]
}

func playerIDInLineup(playerID uint, players []*GamePlayer) bool {
	for _, p := range players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}

type TeamStatDTO struct {
	TeamID               uint
	Team                 string
	GoalsFor             uint16
	GoalsAgainst         uint16
	Assists              uint16
	Points               uint16
	Period1Score         uint8
	Period2Score         uint8
	Period3Score         uint8
	OTScore              uint8
	PlusMinus            int8
	PenaltyMinutes       uint16
	EvenStrengthGoals    uint8
	EvenStrengthPoints   uint8
	PowerPlayGoals       uint8
	PowerPlayPoints      uint8
	ShorthandedGoals     uint8
	ShorthandedPoints    uint8
	OvertimeGoals        uint8
	Shots                uint16
	ShootingPercentage   float32
	FaceOffWinPercentage float32
	FaceOffsWon          uint
	FaceOffs             uint
	ShotsAgainst         uint16
	Saves                uint16
	SavePercentage       float32
	Shutouts             uint16
}

func (t *TeamStatDTO) AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool) {
	t.GoalsFor++
	t.Points++
	t.PlusMinus++
	if isEvenStrength {
		t.EvenStrengthGoals++
		t.EvenStrengthPoints++
	}
	if isPowerPlay {
		t.PowerPlayGoals++
		t.PowerPlayPoints++
	}
	if isShorthanded {
		t.ShorthandedGoals++
		t.ShorthandedPoints++
	}
	if isOvertime {
		t.OvertimeGoals++
	}
}

func (t *TeamStatDTO) AddGoalAgainst() {
	t.GoalsAgainst++
	t.PlusMinus--
}

func (t *TeamStatDTO) AddAssist(isEvenStrength, isPowerPlay, isShorthanded bool) {
	t.Assists++
	t.Points++
	if isEvenStrength {
		t.EvenStrengthPoints++
	}
	if isPowerPlay {
		t.PowerPlayPoints++
	}
	if isShorthanded {
		t.ShorthandedPoints++
	}
}

func (t *TeamStatDTO) AddPenaltyTime(seconds int) {
	t.PenaltyMinutes += uint16(seconds)
}

func (t *TeamStatDTO) AddShot(isScore, isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool, period uint8) {
	t.Shots++
	if isScore {
		t.AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime)
		if period == 1 {
			t.Period1Score += 1
		} else if period == 2 {
			t.Period2Score += 1
		} else if period == 3 {
			t.Period3Score += 1
		} else if period > 3 {
			t.OTScore += 1
		}
	}

	t.ShootingPercentage = float32(t.GoalsFor) / float32(t.Shots)
}

func (t *TeamStatDTO) AddShotAgainst(isScore bool) {
	t.ShotsAgainst++
	if isScore {
		t.AddGoalAgainst()
	} else {
		t.Saves++
	}

	t.SavePercentage = float32(t.Saves) / float32(t.ShotsAgainst)
}

func (t *TeamStatDTO) AddFaceoff(isWin bool) {
	t.FaceOffs++
	if isWin {
		t.FaceOffsWon++
	}
	t.FaceOffWinPercentage = float32(t.FaceOffsWon) / float32(t.FaceOffs)
}

func (t *TeamStatDTO) AddShutout() {
	t.Shutouts++
}

type PlayerStatsDTO struct {
	TeamID               uint
	SeasonID             uint
	PlayerID             uint
	Goals                uint16
	Assists              uint16
	Points               uint16
	PlusMinus            int8
	PenaltyMinutes       uint16
	EvenStrengthGoals    uint8
	EvenStrengthPoints   uint8
	PowerPlayGoals       uint8
	PowerPlayPoints      uint8
	ShorthandedGoals     uint8
	ShorthandedPoints    uint8
	OvertimeGoals        uint8
	GameWinningGoals     uint8
	Shots                uint16
	ShootingPercentage   float32
	TimeOnIce            uint
	FaceOffWinPercentage float32
	FaceOffsWon          uint
	FaceOffs             uint
	GoalieWins           uint8
	GoalieLosses         uint8
	GoalieTies           uint8
	OvertimeLosses       uint8
	ShotsAgainst         uint16
	Saves                uint16
	GoalsAgainst         uint16
	SavePercentage       float32
	Shutouts             uint16
	ShotsBlocked         uint16
	BodyChecks           uint16
	StickChecks          uint16
}

func (p *PlayerStatsDTO) AdjustPlusMinus(isScore bool) {
	if isScore {
		p.PlusMinus++
	} else {
		p.PlusMinus--
	}
}

func (p *PlayerStatsDTO) AddShot(isScore, isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool) {
	p.Shots++
	if isScore {
		p.AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime)
	}
	p.ShootingPercentage = float32(p.Goals) / float32(p.Shots)
}

func (p *PlayerStatsDTO) AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool) {
	p.Goals++
	p.Points++
	p.PlusMinus++
	if isEvenStrength {
		p.EvenStrengthGoals++
		p.EvenStrengthPoints++
	}
	if isPowerPlay {
		p.PowerPlayGoals++
		p.PowerPlayPoints++
	}
	if isShorthanded {
		p.ShorthandedGoals++
		p.ShorthandedPoints++
	}
	if isOvertime {
		p.OvertimeGoals++
	}
}

func (p *PlayerStatsDTO) AddGoalAgainst() {
	p.GoalsAgainst++
	p.PlusMinus--
}

func (p *PlayerStatsDTO) AddAssist(isEvenStrength, isPowerPlay, isShorthanded bool) {
	p.Assists++
	p.Points++
	if isEvenStrength {
		p.EvenStrengthPoints++
	}
	if isPowerPlay {
		p.PowerPlayPoints++
	}
	if isShorthanded {
		p.ShorthandedPoints++
	}
}

func (p *PlayerStatsDTO) AddTimeOnIce(seconds int, isPenalty bool) {
	if isPenalty {
		p.AddPenaltyTime(seconds)
	} else {
		p.TimeOnIce += uint(seconds)
	}
}

func (p *PlayerStatsDTO) AddPenaltyTime(seconds int) {
	p.PenaltyMinutes += uint16(seconds)
}

func (p *PlayerStatsDTO) AddShotAgainst(isScore bool) {
	p.ShotsAgainst++
	if isScore {
		p.AddGoalAgainst()
	} else {
		p.Saves++
	}

	p.SavePercentage = float32(p.Saves) / float32(p.ShotsAgainst)
}

func (p *PlayerStatsDTO) AddShotBlocked() {
	p.ShotsBlocked += 1
}

func (p *PlayerStatsDTO) AddDefensiveHit(isBodyCheck bool) {
	if isBodyCheck {
		p.BodyChecks += 1
	} else {
		p.StickChecks += 1
	}
}

func (p *PlayerStatsDTO) AddFaceoff(isWin bool) {
	p.FaceOffs++
	if isWin {
		p.FaceOffsWon++
	}
	p.FaceOffWinPercentage = float32(p.FaceOffsWon) / float32(p.FaceOffs)
}

func (p *PlayerStatsDTO) AddShutout() {
	p.Shutouts++
}

func (p *PlayerStatsDTO) AddGoalieStat(scoreType int, isOvertimeLoss bool) {
	if scoreType == 1 {
		p.GoalieWins++
	}
	if scoreType == 2 {
		p.GoalieLosses++
		if isOvertimeLoss {
			p.OvertimeLosses++
		}
	}
	if scoreType == 3 {
		p.GoalieTies++
	}
}
