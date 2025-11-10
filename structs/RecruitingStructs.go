package structs

import (
	"math"
	"sort"

	util "github.com/CalebRose/SimHockey/_util"
	"gorm.io/gorm"
)

// RecruitingTeamProfile - The profile for a team for recruiting
type RecruitingTeamProfile struct {
	gorm.Model
	TeamID                uint
	Team                  string
	State                 string
	Country               string
	ScholarshipsAvailable uint8
	WeeklyPoints          float32
	WeeklyScoutingPoints  uint8
	SpentPoints           float32
	TotalCommitments      uint8
	RecruitClassSize      uint8
	PortalReputation      uint8 // A value between 1-100 signifying the coach's reputation and behavior in the transfer portal.
	ESPNScore             float32
	RivalsScore           float32
	Rank247Score          float32
	CompositeScore        float32
	ThreeStars            uint8
	FourStars             uint8
	FiveStars             uint8
	RecruitingClassRank   uint8
	CaughtCheating        bool
	IsAI                  bool
	IsUserTeam            bool
	AIBehavior            string
	AIQuality             string
	WeeksMissed           uint8
	BattlesWon            uint8
	BattlesLost           uint8
	AIMinThreshold        uint8
	AIMaxThreshold        uint8
	AIStarMin             uint8
	AIStarMax             uint8
	Recruiter             string
	OffensiveScheme       string
	DefensiveScheme       string
	Y1Rank                uint16
	Y2Rank                uint16
	Y3Rank                uint16
	Y4Rank                uint16
	Y5Rank                uint16
	Recruits              []RecruitPlayerProfile `gorm:"foreignKey:ProfileID"`
}

func (r *RecruitingTeamProfile) ResetSpentPoints() {
	if r.SpentPoints == 0 && r.TotalCommitments < r.RecruitClassSize {
		r.WeeksMissed += 1
	} else {
		r.WeeksMissed = 0
	}
	if r.TotalCommitments == r.RecruitClassSize {
		r.WeeksMissed = 0
	}
	r.SpentPoints = 0
}

func (r *RecruitingTeamProfile) ResetScoutingPoints(week int) {
	if week == 0 {
		r.WeeklyScoutingPoints = 30
	} else {
		r.WeeklyScoutingPoints = 10
	}
}

func (r *RecruitingTeamProfile) SubtractScoutingPoints() {
	if r.WeeklyScoutingPoints > 0 {
		r.WeeklyScoutingPoints--
	}
}

func (r *RecruitingTeamProfile) SubtractScholarshipsAvailable() {
	if r.ScholarshipsAvailable > 0 {
		r.ScholarshipsAvailable--
	}
}

func (r *RecruitingTeamProfile) ReallocateScholarship() {
	if r.ScholarshipsAvailable < 20 {
		r.ScholarshipsAvailable++
	}
}

func (r *RecruitingTeamProfile) ResetScholarshipCount() {
	r.ScholarshipsAvailable = 20
}

func (r *RecruitingTeamProfile) AdjustPortalReputation(points int8) {
	adj := int8(r.PortalReputation) + points
	if adj < 0 {
		adj = 1
	}
	r.PortalReputation = uint8(adj)
}

func (r *RecruitingTeamProfile) AllocateSpentPoints(points float32) {
	r.SpentPoints = points
}

func (r *RecruitingTeamProfile) AIAllocateSpentPoints(points float32) {
	r.SpentPoints += points
}

func (r *RecruitingTeamProfile) ResetWeeklyPoints(points float32, isOffseason bool) {
	r.WeeklyPoints = points
	if isOffseason {
		r.WeeklyScoutingPoints = 30
	} else {
		r.WeeklyScoutingPoints = 10
	}
}

func (r *RecruitingTeamProfile) AddRecruitsToProfile(croots []RecruitPlayerProfile) {
	r.Recruits = croots
}

func (r *RecruitingTeamProfile) AssignRivalsRank(score float32) {
	r.RivalsScore = score
}

func (r *RecruitingTeamProfile) Assign247Rank(score float32) {
	r.Rank247Score = score
}

func (r *RecruitingTeamProfile) AssignESPNRank(score float32) {
	r.ESPNScore = score
}

func (r *RecruitingTeamProfile) AssignCompositeRank(score float32) {
	r.CompositeScore = score
}

func (r *RecruitingTeamProfile) AssignHistoricRank(rank int) {
	r.Y1Rank = uint16(rank)
	r.Y2Rank = r.Y1Rank
	r.Y3Rank = r.Y2Rank
	r.Y4Rank = r.Y3Rank
	r.Y5Rank = r.Y4Rank
}

func (r *RecruitingTeamProfile) UpdateTotalSignedRecruits(num uint8) {
	r.TotalCommitments = num
}

func (r *RecruitingTeamProfile) IncreaseCommitCount() {
	r.TotalCommitments++
}

func (r *RecruitingTeamProfile) ApplyCaughtCheating() {
	r.CaughtCheating = true
}

func (r *RecruitingTeamProfile) ActivateAI() {
	r.IsAI = true
	r.IsUserTeam = false
}

func (r *RecruitingTeamProfile) DeactivateAI() {
	r.IsAI = false
	r.IsUserTeam = true
}

func (r *RecruitingTeamProfile) ToggleAIBehavior() {
	r.IsAI = !r.IsAI
}

func (r *RecruitingTeamProfile) UpdateAIBehavior(isAi bool, starMax, starMin, min, max uint8) {
	r.IsAI = isAi
	r.AIStarMax = starMax
	r.AIStarMin = starMin
	r.AIMinThreshold = min
	r.AIMaxThreshold = max
}

func (r *RecruitingTeamProfile) SetRecruitingClassSize(val uint8) {
	if val > 10 {
		r.RecruitClassSize = 10
	} else {
		r.RecruitClassSize = val
	}
}

func (r *RecruitingTeamProfile) IncrementClassSize() {
	if r.RecruitClassSize < 25 {
		r.RecruitClassSize += 1
	}
}

func (r *RecruitingTeamProfile) AddBattleWon() {
	r.BattlesWon += 1
}

func (r *RecruitingTeamProfile) AddBattleLost() {
	r.BattlesLost += 1
}

func (r *RecruitingTeamProfile) ResetStarCount() {
	r.ThreeStars = 0
	r.FourStars = 0
	r.FiveStars = 0
}

func (r *RecruitingTeamProfile) AddStarPlayer(stars uint8) {
	switch stars {
	case 3:
		r.ThreeStars += 1
	case 4:
		r.FourStars += 1
	case 5:
		r.FiveStars += 1
	}
}

func (r *RecruitingTeamProfile) AssignRecruiter(name string) {
	r.Recruiter = name
}

// RecruitPlayerProfile - Individual points profile for a Team's Recruiting Portfolio
type RecruitPlayerProfile struct {
	gorm.Model
	SeasonID             uint
	RecruitID            uint
	ProfileID            uint
	TotalPoints          float32
	CurrentWeeksPoints   float32
	PreviousWeekPoints   float32
	Modifier             float32
	IsHomeState          bool
	IsPipelineState      bool
	SpendingCount        uint8
	Scholarship          bool
	ScholarshipRevoked   bool
	RemovedFromBoard     bool
	IsSigned             bool
	IsLocked             bool
	CaughtCheating       bool
	TeamReachedMax       bool
	Agility              bool
	Faceoffs             bool
	LongShotAccuracy     bool
	LongShotPower        bool
	CloseShotAccuracy    bool
	CloseShotPower       bool
	OneTimer             bool
	Passing              bool
	PuckHandling         bool
	Strength             bool
	BodyChecking         bool
	StickChecking        bool
	ShotBlocking         bool
	Goalkeeping          bool
	GoalieVision         bool
	GoalieReboundControl bool
	Recruit              Recruit `gorm:"foreignKey:RecruitID"`
	// RecruitPoints             []RecruitPointAllocation `gorm:"foreignKey:RecruitProfileID"`
}

func (rp *RecruitPlayerProfile) AllocateCurrentWeekPoints(points float32) {
	rp.CurrentWeeksPoints = points
}

func (rp *RecruitPlayerProfile) AddCurrentWeekPointsToTotal(CurrentPoints float32) {
	// If user spends points on a recruit
	if CurrentPoints > 0 {
		rp.TotalPoints += CurrentPoints
		if rp.SpendingCount < 5 && CurrentPoints >= 1 {
			rp.SpendingCount++
			// In the event that someone tries to exploit the consistency system with a value between 0.00001 and 0.99999
		} else if CurrentPoints > 0 && CurrentPoints < 1 {
			rp.SpendingCount = 0
		}
	} else {
		rp.TotalPoints = 0
		rp.CaughtCheating = true
		rp.SpendingCount = 0
	}
	rp.PreviousWeekPoints = rp.CurrentWeeksPoints
	rp.CurrentWeeksPoints = 0
}

func (rp *RecruitPlayerProfile) ToggleRemoveFromBoard() {
	rp.RemovedFromBoard = !rp.RemovedFromBoard
	rp.CurrentWeeksPoints = 0
}

func (rp *RecruitPlayerProfile) ToggleScholarship(rewardScholarship bool, revokeScholarship bool) {
	rp.Scholarship = rewardScholarship
	rp.ScholarshipRevoked = revokeScholarship
}

func (rp *RecruitPlayerProfile) SignPlayer() {
	if rp.Scholarship {
		rp.IsSigned = true
		rp.IsLocked = true
	}
}

func (rp *RecruitPlayerProfile) LockPlayer() {
	rp.IsLocked = true
}

func (rp *RecruitPlayerProfile) ResetSpendingCount() {
	rp.SpendingCount = 0
}

func (rp *RecruitPlayerProfile) ResetTotalPoints() {
	rp.TotalPoints = 0
	rp.TeamReachedMax = true
}

func (rp *RecruitPlayerProfile) ApplyModifier(mod float32) {
	rp.Modifier = mod
}

func (rp *RecruitPlayerProfile) ApplyScoutingAttribute(attr string) {
	if attr == util.Faceoffs {
		rp.Faceoffs = true
	}
	if attr == util.Agility {
		rp.Agility = true
	}
	if attr == util.LongShotAccuracy {
		rp.LongShotAccuracy = true
	}
	if attr == util.LongShotPower {
		rp.LongShotPower = true
	}
	if attr == util.CloseShotAccuracy {
		rp.CloseShotAccuracy = true
	}
	if attr == util.CloseShotPower {
		rp.CloseShotPower = true
	}
	if attr == util.Strength {
		rp.Strength = true
	}
	if attr == util.Passing {
		rp.Passing = true
	}
	if attr == util.PuckHandling {
		rp.PuckHandling = true
	}
	if attr == util.BodyChecking {
		rp.BodyChecking = true
	}
	if attr == util.StickChecking {
		rp.StickChecking = true
	}
	if attr == util.Goalkeeping {
		rp.Goalkeeping = true
	}
	if attr == util.GoalieVision {
		rp.GoalieVision = true
	}
	if attr == util.ShotBlocking {
		rp.ShotBlocking = true
	}
}

type Croot struct {
	ID                uint
	TeamID            uint
	College           string
	FirstName         string
	LastName          string
	Position          string
	Archetype         string
	Height            uint8
	Weight            uint16
	Stars             uint8
	PotentialGrade    string
	Personality       string
	RecruitingBias    string
	AcademicBias      string
	WorkEthic         string
	HighSchool        string
	City              string
	State             string
	Country           string
	AffinityOne       string
	AffinityTwo       string
	RecruitingStatus  string
	RecruitModifier   float32
	IsCustomCroot     bool
	CustomCrootFor    string
	IsSigned          bool
	OverallGrade      string
	Agility           uint8 // How fast a player can go in a zone without a defense check
	Faceoffs          uint8 // Ability to win faceoffs
	LongShotAccuracy  uint8 // Accuracy on non-close shots
	LongShotPower     uint8 // Power on non-close shots. High power means less shotblocking
	CloseShotAccuracy uint8 // Accuracy on close shots. Great on pass plays
	CloseShotPower    uint8 // Power on Close shots
	OneTimer          uint8 // Shots bassed on passing. Essentially a modifier that gets greater with each pass in a zone
	Passing           uint8 // Passing ability
	PuckHandling      uint8 // Ability to handle the puck when going between zones.
	Strength          uint8 // General modifier on all physical attributes. Also used in fights
	BodyChecking      uint8 // Physical defense check.
	StickChecking     uint8 // Non-phyisical defense check
	ShotBlocking      uint8 // Ability for defensemen to block a shot being made
	Goalkeeping       uint8 // Goalkeepers' ability to block a shot
	GoalieVision      uint8 // Goalkeepers' vision
	TotalRank         float32
	InjuryRating      string
	Stamina           string
	BaseRecruitingGrades
	PlayerPreferences
	LeadingTeams []LeadingTeams
}

type BaseRecruitingGrades struct {
	// Each attribute has a chance to grow at a different rate. These are all small modifiers
	AgilityGrade           string // Ability to switch between zones
	FaceoffsGrade          string // Ability to win faceoffs
	CloseShotAccuracyGrade string // Accuracy on close shots
	CloseShotPowerGrade    string // Power on close shots. High power means less shotblocking
	LongShotAccuracyGrade  string // Accuracy on far shots. Great on pass plays
	LongShotPowerGrade     string // Accuracy on far shots
	PassingGrade           string // Power on close shots. Great on pass plays
	PuckHandlingGrade      string // Ability to handle the puck when going between zones.
	StrengthGrade          string // General modifier on all physical attributes. Also used in fights
	BodyCheckingGrade      string // Physical defense check.
	StickCheckingGrade     string // Non-phyisical defense check
	ShotBlockingGrade      string // Ability for defensemen to block a shot being made
	GoalkeepingGrade       string // Goalkeepers' ability to block a shot
	GoalieVisionGrade      string // Goalkeepers' ability to block a shot
	GoalieReboundGrade     string // Goalkeepers' ability to block a shot
}

func (g *BaseRecruitingGrades) MapLetterGrades(p BasePotentials) {
	g.AgilityGrade = util.GetPotentialGrade(int(p.AgilityPotential))
	g.FaceoffsGrade = util.GetPotentialGrade(int(p.FaceoffsPotential))
	g.CloseShotAccuracyGrade = util.GetPotentialGrade(int(p.CloseShotAccuracyPotential))
	g.CloseShotPowerGrade = util.GetPotentialGrade(int(p.CloseShotPowerPotential))
	g.LongShotAccuracyGrade = util.GetPotentialGrade(int(p.LongShotAccuracyPotential))
	g.LongShotPowerGrade = util.GetPotentialGrade(int(p.LongShotPowerPotential))
	g.PassingGrade = util.GetPotentialGrade(int(p.PassingPotential))
	g.PuckHandlingGrade = util.GetPotentialGrade(int(p.PuckHandlingPotential))
	g.StrengthGrade = util.GetPotentialGrade(int(p.StrengthPotential))
	g.BodyCheckingGrade = util.GetPotentialGrade(int(p.BodyCheckingPotential))
	g.StickCheckingGrade = util.GetPotentialGrade(int(p.StickCheckingPotential))
	g.ShotBlockingGrade = util.GetPotentialGrade(int(p.ShotBlockingPotential))
	g.GoalkeepingGrade = util.GetPotentialGrade(int(p.GoalkeepingPotential))
	g.GoalieVisionGrade = util.GetPotentialGrade(int(p.GoalieVisionPotential))
}

type LeadingTeams struct {
	TeamID         uint
	TeamName       string
	TeamAbbr       string
	Odds           float32
	HasScholarship bool
}

// Sorting Funcs
type ByLeadingPoints []LeadingTeams

func (rp ByLeadingPoints) Len() int      { return len(rp) }
func (rp ByLeadingPoints) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByLeadingPoints) Less(i, j int) bool {
	return rp[i].Odds > rp[j].Odds
}

func (c *Croot) Map(r Recruit) {
	c.ID = r.ID
	c.TeamID = uint(r.TeamID)
	c.FirstName = r.FirstName
	c.LastName = r.LastName
	c.Position = r.Position
	c.Archetype = r.Archetype
	c.Height = r.Height
	c.Weight = r.Weight
	c.Stars = r.Stars
	c.Personality = r.Personality
	c.HighSchool = r.HighSchool
	c.City = r.City
	c.State = r.State
	c.Country = r.Country
	c.College = r.College
	c.OverallGrade = util.GetLetterGrade(int(r.Overall), 1)
	c.IsSigned = r.IsSigned
	c.RecruitingStatus = r.RecruitingStatus
	c.RecruitModifier = r.RecruitingModifier
	c.IsCustomCroot = r.IsCustomCroot
	c.CustomCrootFor = r.CustomCrootFor
	c.BaseRecruitingGrades.MapLetterGrades(r.BasePotentials)
	c.PlayerPreferences = r.PlayerPreferences
	c.Agility = r.Agility
	c.Faceoffs = r.Faceoffs
	c.LongShotAccuracy = r.LongShotAccuracy
	c.LongShotPower = r.LongShotPower
	c.CloseShotAccuracy = r.CloseShotAccuracy
	c.CloseShotPower = r.CloseShotPower
	c.Passing = r.Passing
	c.PuckHandling = r.PuckHandling
	c.Strength = r.Strength
	c.BodyChecking = r.BodyChecking
	c.StickChecking = r.StickChecking
	c.ShotBlocking = r.ShotBlocking
	c.Goalkeeping = r.Goalkeeping
	c.GoalieVision = r.GoalieVision
	c.Stamina = util.GetPotentialGrade(int(r.Stamina))
	c.InjuryRating = util.GetPotentialGrade(int(r.InjuryRating))

	mod := r.TopRankModifier
	if mod == 0 {
		mod = 1
	}
	c.TotalRank = (r.RivalsRank + r.ESPNRank + r.Rank247) / r.TopRankModifier
	if math.IsNaN(float64(c.TotalRank)) {
		c.TotalRank = 0
	}

	var totalPoints float32 = 0
	var runningThreshold float32 = 0

	sortedProfiles := r.RecruitPlayerProfiles

	sort.Sort(ByPoints(sortedProfiles))

	for _, recruitProfile := range sortedProfiles {
		if recruitProfile.TeamReachedMax {
			continue
		}
		if runningThreshold == 0 {
			runningThreshold = float32(recruitProfile.TotalPoints) * 0.66
		}

		if recruitProfile.TotalPoints >= runningThreshold {
			totalPoints += float32(recruitProfile.TotalPoints)
		}

	}

	for i := 0; i < len(sortedProfiles); i++ {
		if sortedProfiles[i].TeamReachedMax || sortedProfiles[i].RemovedFromBoard {
			continue
		}
		var odds float32 = 0

		if sortedProfiles[i].TotalPoints >= runningThreshold && runningThreshold > 0 {
			odds = float32(sortedProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamID:         uint(r.RecruitPlayerProfiles[i].ProfileID),
			Odds:           odds,
			HasScholarship: r.RecruitPlayerProfiles[i].Scholarship,
		}
		c.LeadingTeams = append(c.LeadingTeams, leadingTeam)
	}
	sort.Sort(ByLeadingPoints(c.LeadingTeams))
}

type ByCrootRank []Croot

func (c ByCrootRank) Len() int      { return len(c) }
func (c ByCrootRank) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByCrootRank) Less(i, j int) bool {
	return c[i].TotalRank > c[j].TotalRank || c[i].Stars > c[j].Stars
}

// Sorting Funcs
type ByPoints []RecruitPlayerProfile

func (rp ByPoints) Len() int      { return len(rp) }
func (rp ByPoints) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByPoints) Less(i, j int) bool {
	return rp[i].TotalPoints > rp[j].TotalPoints
}

type CreateRecruitProfileDto struct {
	PlayerID      int
	SeasonID      int
	RecruitID     int
	ProfileID     int
	Team          CollegeTeam
	PlayerRecruit Croot
	Recruiter     string
}

type UpdateRecruitProfileDto struct {
	RecruitPointsID   int
	RecruitID         int
	ProfileID         int
	Team              string
	WeekID            int
	AllocationID      int
	SpentPoints       int
	RewardScholarship bool
	RevokeScholarship bool
}

type CrootProfile struct {
	ID                   uint
	SeasonID             uint
	RecruitID            uint
	ProfileID            uint
	TotalPoints          float32
	CurrentWeeksPoints   float32
	Modifier             float32
	SpendingCount        uint8
	Scholarship          bool
	ScholarshipRevoked   bool
	IsCloseToHome        bool
	IsPipeline           bool
	TeamAbbreviation     string
	RemovedFromBoard     bool
	IsSigned             bool
	IsLocked             bool
	CaughtCheating       bool
	Agility              bool
	Faceoffs             bool
	LongShotAccuracy     bool
	LongShotPower        bool
	CloseShotAccuracy    bool
	CloseShotPower       bool
	OneTimer             bool
	Passing              bool
	PuckHandling         bool
	Strength             bool
	BodyChecking         bool
	StickChecking        bool
	ShotBlocking         bool
	Goalkeeping          bool
	GoalieVision         bool
	GoalieReboundControl bool
	Recruit              Croot
}

func (cp *CrootProfile) Map(rp RecruitPlayerProfile, c Croot) {
	cp.ID = rp.ID
	cp.SeasonID = rp.SeasonID
	cp.RecruitID = rp.RecruitID
	cp.ProfileID = rp.ProfileID
	cp.TotalPoints = rp.TotalPoints
	cp.CurrentWeeksPoints = rp.CurrentWeeksPoints
	cp.SpendingCount = rp.SpendingCount
	cp.Modifier = rp.Modifier
	cp.Scholarship = rp.Scholarship
	cp.ScholarshipRevoked = rp.ScholarshipRevoked
	cp.IsCloseToHome = rp.IsHomeState
	cp.IsPipeline = rp.IsPipelineState
	cp.RemovedFromBoard = rp.RemovedFromBoard
	cp.IsSigned = rp.IsSigned
	cp.IsLocked = rp.IsLocked
	cp.CaughtCheating = rp.CaughtCheating
	cp.Agility = rp.Agility
	cp.Faceoffs = rp.Faceoffs
	cp.LongShotAccuracy = rp.LongShotAccuracy
	cp.LongShotPower = rp.LongShotPower
	cp.CloseShotAccuracy = rp.CloseShotAccuracy
	cp.CloseShotPower = rp.CloseShotPower
	cp.OneTimer = rp.OneTimer
	cp.Passing = rp.Passing
	cp.PuckHandling = rp.PuckHandling
	cp.Strength = rp.Strength
	cp.BodyChecking = rp.BodyChecking
	cp.StickChecking = rp.StickChecking
	cp.ShotBlocking = rp.ShotBlocking
	cp.Goalkeeping = rp.Goalkeeping
	cp.GoalieVision = rp.GoalieVision
	cp.GoalieReboundControl = rp.GoalieReboundControl
	cp.Recruit = c
}

// Sorting Funcs
type ByCrootProfileTotal []CrootProfile

func (rp ByCrootProfileTotal) Len() int      { return len(rp) }
func (rp ByCrootProfileTotal) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByCrootProfileTotal) Less(i, j int) bool {
	return rp[i].TotalPoints > rp[j].TotalPoints
}

type SimTeamBoardResponse struct {
	ID                        uint
	TeamID                    uint
	Team                      string
	State                     string
	ScholarshipsAvailable     int
	WeeklyPoints              float32
	SpentPoints               float32
	TotalCommitments          uint
	RecruitClassSize          uint
	BaseEfficiencyScore       float32
	RecruitingEfficiencyScore float32
	PreviousOverallWinPer     float32
	PreviousConferenceWinPer  float32
	CurrentOverallWinPer      float32
	CurrentConferenceWinPer   float32
	ESPNScore                 float32
	RivalsScore               float32
	Rank247Score              float32
	CompositeScore            float32
	IsAI                      bool
	IsUserTeam                bool
	BattlesWon                uint
	BattlesLost               uint
	AIMinThreshold            uint
	AIMaxThreshold            uint
	AIStarMin                 uint
	AIStarMax                 uint
	AIAutoOfferscholarships   bool
	OffensiveScheme           string
	DefensiveScheme           string
	Recruiter                 string
	RecruitingClassRank       int
	Recruits                  []CrootProfile
}

func (stbr *SimTeamBoardResponse) Map(rtp RecruitingTeamProfile, c []CrootProfile) {
	stbr.ID = rtp.ID
	stbr.TeamID = rtp.TeamID
	stbr.Team = rtp.Team
	stbr.IsAI = rtp.IsAI
	stbr.State = rtp.State
	stbr.ScholarshipsAvailable = int(rtp.ScholarshipsAvailable)
	stbr.WeeklyPoints = rtp.WeeklyPoints
	stbr.SpentPoints = rtp.SpentPoints
	stbr.TotalCommitments = uint(rtp.TotalCommitments)
	stbr.ESPNScore = rtp.ESPNScore
	stbr.RivalsScore = rtp.RivalsScore
	stbr.Rank247Score = rtp.Rank247Score
	stbr.CompositeScore = rtp.CompositeScore
	stbr.RecruitingClassRank = int(rtp.RecruitingClassRank)
	stbr.Recruits = c
	stbr.RecruitClassSize = uint(rtp.RecruitClassSize)
	stbr.IsUserTeam = rtp.IsUserTeam
	stbr.BattlesWon = uint(rtp.BattlesWon)
	stbr.BattlesLost = uint(rtp.BattlesLost)
	stbr.AIMinThreshold = uint(rtp.AIMinThreshold)
	stbr.AIMaxThreshold = uint(rtp.AIMaxThreshold)
	stbr.AIStarMin = uint(rtp.AIStarMin)
	stbr.AIStarMax = uint(rtp.AIStarMax)
	stbr.Recruiter = rtp.Recruiter
}

// UpdateRecruitingBoardDTO - Data Transfer Object from UI to API
type UpdateRecruitingBoardDTO struct {
	Profile  RecruitingTeamProfile
	Recruits []RecruitPlayerProfile
	TeamID   int
}

type RecruitPointAllocation struct {
	gorm.Model
	RecruitID          uint
	TeamProfileID      uint
	RecruitProfileID   uint
	WeekID             uint
	Points             float32
	ModAffectedPoints  float32
	IsHomeStateApplied bool
	IsPipelineApplied  bool
	CaughtCheating     bool
}

type RecruitingOdds struct {
	Odds          int
	IsCloseToHome bool
	IsPipeline    bool
}

type ScoutAttributeDTO struct {
	ProfileID uint
	RecruitID uint
	Attribute string
}

type SchemeFits struct {
	GoodFits []string
	BadFits  []string
}
