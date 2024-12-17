package engine

import "github.com/CalebRose/SimHockey/structs"

type GameState struct {
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
	PowerPlayTime         uint8
	IsPowerPlay           bool
	PowerPlayTeamID       uint
	Zones                 []string // Home Goal, Home Zone, Neutral Zone, Away Zone, Away Goal
	PuckLocation          string
	PuckCarrier           GamePlayer
	AssistingPlayer       GamePlayer
	Momentum              float64
	PossessingTeam        uint // Use the ID of the team
	IsCollegeGame         bool
	IsPlayoffGame         bool
	Collector             structs.PbPCollector
}

func (gs *GameState) RecordPlay(play structs.PlayByPlay) {
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

	// Ensure time doesn't overflow
	if gs.TimeOnClock > MaxTimeOnClock {
		gs.TimeOnClock = 0
	}
}

func handleLineups(gs *GameState) {
	gs.HomeStrategy.HandleLineups(int(gs.SecondsConsumed))
	gs.AwayStrategy.HandleLineups(int(gs.SecondsConsumed))
}

func (gs *GameState) SetSecondsConsumed(seconds uint16) {
	gs.SecondsConsumed += seconds
}

func (gs *GameState) SetPowerPlay(duration, teamID int) {
	gs.Momentum = 0
	gs.PowerPlayTime = uint8(duration)
	gs.IsPowerPlay = true
	gs.PowerPlayTeamID = uint(teamID)
}

func (gs *GameState) SetNewZone(zone string) {
	gs.PuckLocation = zone
}

func (gs *GameState) SetFaceoffOnCenterIce(check bool) {
	gs.FaceoffOnCenterIce = check
}

func (gs *GameState) SetPuckBearer(player GamePlayer) {
	if gs.PuckCarrier.TeamID != player.TeamID {
		gs.ResetMomentum()
		gs.AssistingPlayer = GamePlayer{}
	} else {
		gs.Momentum += 0.125
		gs.AssistingPlayer = gs.PuckCarrier
	}
	gs.PuckCarrier = player
	gs.PossessingTeam = uint(player.TeamID)
}

func (gs *GameState) TriggerBreakaway() {
	gs.Momentum += 0.275
}

func (gs *GameState) ResetMomentum() {
	gs.Momentum = 0
}

func (gs *GameState) GetCenter(isHome bool) GamePlayer {
	var currentLineup []GamePlayer
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
		gs.HomeTeamStats.AddShot(true, !gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID), gs.Period > 3)
		gs.AwayTeamStats.AddShotAgainst(true)
	} else {
		gs.AwayTeamScore += 1
		gs.AwayTeamStats.AddShot(true, !gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID), gs.Period > 3)
		gs.HomeTeamStats.AddShotAgainst(true)
	}
	gs.PuckCarrier.Stats.AddShot(true, !gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID), gs.Period > 3)
	if gs.AssistingPlayer.ID > 0 {
		gs.AssistingPlayer.Stats.AddAssist(!gs.IsPowerPlay, gs.IsPowerPlay, gs.IsPowerPlay && gs.PowerPlayTeamID != uint(gs.PuckCarrier.TeamID))
	}
	gs.HomeStrategy.HandlePlusMinus(isHome, gs.PuckCarrier.ID, gs.AssistingPlayer.ID)
	gs.AwayStrategy.HandlePlusMinus(isHome, gs.PuckCarrier.ID, gs.AssistingPlayer.ID)
	gs.SetFaceoffOnCenterIce(true)
	gs.SetNewZone(NeutralZone)
	gs.SetPuckBearer(GamePlayer{})
}

func (gs *GameState) AddShots(isHome bool) {
	if isHome {
		gs.HomeTeamStats.AddShot(false, false, false, false, false)
		gs.AwayTeamStats.AddShotAgainst(false)
	} else {
		gs.AwayTeamStats.AddShot(false, false, false, false, false)
		gs.HomeTeamStats.AddShotAgainst(false)
	}
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

func GetGamePlayerByPosition(currentLineup []GamePlayer, pos string) GamePlayer {
	for _, p := range currentLineup {
		if p.Position != "C" || p.IsOut {
			continue
		}
		return p
	}
	return GamePlayer{}
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
	BenchPlayers                []GamePlayer
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

func (gp *GamePlaybook) handleLineReplacement(players []GamePlayer, playerID uint, requiredCount, lineType uint) []GamePlayer {
	filteredPlayers, queue := filterOutPlayerFromLineup(players, playerID, lineType)

	for len(filteredPlayers) < int(requiredCount) {
		var replacement GamePlayer

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

func (gp *GamePlaybook) HandleLineups(secondsConsumed int) {
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

func (gp *GamePlaybook) getNextLine(lineType uint) []GamePlayer {
	if lineType == 1 {
		for i := gp.CurrentForwards + 1; i < len(gp.Forwards); i++ {
			return gp.Forwards[i].Players
		}
	} else if lineType == 2 {
		for i := gp.CurrentDefenders + 1; i < len(gp.Defenders); i++ {
			return gp.Defenders[i].Players
		}
	}
	return []GamePlayer{}
}

type LineStrategy struct {
	structs.Allocations
	Players        []GamePlayer
	TotalStamina   int
	CurrentStamina int
	Threshold      int
}

func (ls *LineStrategy) SetNewLineup(players []GamePlayer) {
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
	CurrentStamina  int
	OneTimerMod     float64
	AgilityMod      float64
	StrengthMod     float64
	WristShotMod    float64
	WristShotAccMod float64
	SlapshotMod     float64
	SlapshotAccMod  float64
	FaceoffMod      float64
	HandlingMod     float64
	PassMod         float64
	StickCheckMod   float64
	BodyCheckMod    float64
	GoalkeepingMod  float64
	GoalieVisionMod float64
	ShotblockingMod float64
	IsOut           bool
	FoulOut         bool
	SubstituteID    uint
	Stats           PlayerStatsDTO
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
	g.WristShotAccMod = calculateAttributeModifier(float64(g.WristShotAccuracy), ModifierFactor)
	g.WristShotMod = calculateAttributeModifier(float64(g.WristShotPower), ModifierFactor)
	g.SlapshotMod = calculateAttributeModifier(float64(g.SlapshotPower), ModifierFactor)
	g.SlapshotAccMod = calculateAttributeModifier(float64(g.SlapshotAccuracy), ModifierFactor)
	g.FaceoffMod = calculateAttributeModifier(float64(g.Faceoffs), ModifierFactor)
	g.HandlingMod = calculateAttributeModifier(float64(g.PuckHandling), ModifierFactor)
	g.PassMod = calculateAttributeModifier(float64(g.Passing), ModifierFactor)
	g.StickCheckMod = calculateAttributeModifier(float64(g.StickChecking), ModifierFactor)
	g.BodyCheckMod = calculateAttributeModifier(float64(g.BodyChecking), ModifierFactor)
	g.GoalkeepingMod = calculateAttributeModifier(float64(g.Goalkeeping), ModifierFactor)
	g.GoalieVisionMod = calculateAttributeModifier(float64(g.GoalieVision), ModifierFactor)
	g.ShotblockingMod = calculateAttributeModifier(float64(g.ShotBlocking), ModifierFactor)
}

func (g *GamePlayer) GoToPenaltyBox() {
	g.IsOut = true
}

func (g *GamePlayer) ReturnToPlay() {
	g.IsOut = false
}

// Util Structs
// PlayerWeight -- For event checks
type PlayerWeight struct {
	PlayerID uint
	Weight   float64
}

type Penalty struct {
	PenaltyID     uint
	PenaltyName   string
	PenaltyType   uint // 0 == Can occur anywhere, 1 == Defending Zones, 2== Goal Defending Zones
	Severity      string
	Weight        float64
	IsFight       bool
	AggressionReq uint8
	DisciplineReq uint8
	Context       string
}

// Line Management Functions and Structs
type RemovalQueue struct {
	PlayerID       uint
	Player         GamePlayer
	SubstitutionID uint
	LineType       uint
	Idx            uint
}

func filterOutPlayerFromLineup(players []GamePlayer, playerID uint, lineType uint) ([]GamePlayer, RemovalQueue) {
	filtered := []GamePlayer{}
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

func GetPlayerFromLine(playerID uint, players []GamePlayer) RemovalQueue {
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

func popPlayerFromBench(bench []GamePlayer) (GamePlayer, []GamePlayer) {
	if len(bench) == 0 {
		return GamePlayer{}, bench
	}
	return bench[0], bench[1:]
}

func playerIDInLineup(playerID uint, players []GamePlayer) bool {
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

func (t *TeamStatDTO) AddShot(isScore, isEvenStrength, isPowerPlay, isShorthanded, isOvertime bool) {
	t.Shots++
	if isScore {
		t.AddGoal(isEvenStrength, isPowerPlay, isShorthanded, isOvertime)
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
