package structs

import (
	"database/sql"
	"sort"

	util "github.com/CalebRose/SimHockey/_util"
	"gorm.io/gorm"
)

type CollegePromise struct {
	gorm.Model
	TeamID          uint
	CollegePlayerID uint
	PromiseType     string // Snaps (at least minimum), Wins (varies), Bowl Game (Medium), Conf Championship (High), Playoffs (Very High), National Championship (very High), Gameplan Fit (medium), Adjust Gameplan (Low), Play Game In State (Low)
	PromiseWeight   string // The impact the promise will have on their decision. Low, Medium, High
	Benchmark       int    // The value that must be met. For wins & minutes
	BenchmarkStr    string // Needed value for benchmarks that are a string
	PromiseMade     bool   // The player has agreed to the premise of the promise
	IsFullfilled    bool   // If the promise was accomplished
	IsActive        bool   // Whether the promise is active
}

func (p *CollegePromise) Reactivate(promtype, weight string, benchmark int) {
	p.IsActive = true
	p.PromiseType = promtype
	p.PromiseWeight = weight
	p.Benchmark = benchmark
}

func (p *CollegePromise) UpdatePromise(promtype, weight string, benchmark int) {
	p.PromiseType = promtype
	p.PromiseWeight = weight
	p.Benchmark = benchmark
}

func (p *CollegePromise) Deactivate() {
	p.IsActive = false
}

func (p *CollegePromise) MakePromise() {
	p.PromiseMade = true
}

func (p *CollegePromise) FulfillPromise() {
	p.IsFullfilled = true
}

type TransferPortalBoardDto struct {
	Profiles []TransferPortalProfile
}

// Player Profile For the Transfer Portal?
type TransferPortalProfile struct {
	gorm.Model
	SeasonID              uint
	CollegePlayerID       uint
	ProfileID             uint
	PromiseID             sql.NullInt64
	TeamAbbreviation      string
	TotalPoints           float64
	CurrentWeeksPoints    int
	PreviouslySpentPoints int
	SpendingCount         int
	RemovedFromBoard      bool
	RolledOnPromise       bool
	LockProfile           bool
	IsSigned              bool
	Agility               bool
	Faceoffs              bool
	LongShotAccuracy      bool
	LongShotPower         bool
	CloseShotAccuracy     bool
	CloseShotPower        bool
	OneTimer              bool
	Passing               bool
	PuckHandling          bool
	Strength              bool
	BodyChecking          bool
	StickChecking         bool
	ShotBlocking          bool
	Goalkeeping           bool
	GoalieVision          bool
	GoalieReboundControl  bool
	Recruiter             string
}

func (p *TransferPortalProfile) Reactivate() {
	p.RemovedFromBoard = false
}

func (p *TransferPortalProfile) RemovePromise() {
	p.PromiseID = sql.NullInt64{
		Int64: 0,
		Valid: false,
	}
}

func (p *TransferPortalProfile) SignPlayer() {
	p.IsSigned = true
	p.LockProfile = true
	p.CurrentWeeksPoints = 0
}

func (p *TransferPortalProfile) Lock() {
	p.LockProfile = true
	p.CurrentWeeksPoints = 0
}

func (p *TransferPortalProfile) Deactivate() {
	p.RemovedFromBoard = true
	p.CurrentWeeksPoints = 0
}

func (p *TransferPortalProfile) AllocatePoints(points int) {
	p.CurrentWeeksPoints = points
}

func (p *TransferPortalProfile) AddPointsToTotal(multiplier float64) {
	sum := (float64(p.CurrentWeeksPoints) * multiplier)
	p.TotalPoints += sum
	if p.CurrentWeeksPoints == 0 {
		p.SpendingCount = 0
	} else {
		p.SpendingCount += 1
	}
	p.PreviouslySpentPoints = p.CurrentWeeksPoints
}

func (p *TransferPortalProfile) AssignPromise(id uint) {
	p.PromiseID = sql.NullInt64{Valid: true, Int64: int64(id)}
}
func (p *TransferPortalProfile) ToggleRolledOnPromise() {
	p.RolledOnPromise = true
}

func (rp *TransferPortalProfile) ApplyScoutingAttribute(attr string) {
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

type TransferPlayerResponse struct {
	FirstName           string
	LastName            string
	Archetype           string
	Position            string
	PositionTwo         string
	ArchetypeTwo        string
	Age                 int
	Year                int
	State               string
	Country             string
	Stars               int
	Height              int
	Weight              int
	PotentialGrade      string
	Overall             int
	Stamina             int
	Injury              int
	FootballIQ          int
	Speed               int
	Carrying            int
	Agility             int
	Catching            int
	RouteRunning        int
	ZoneCoverage        int
	ManCoverage         int
	Strength            int
	Tackle              int
	PassBlock           int
	RunBlock            int
	PassRush            int
	RunDefense          int
	ThrowPower          int
	ThrowAccuracy       int
	KickAccuracy        int
	KickPower           int
	PuntAccuracy        int
	PuntPower           int
	OverallGrade        string
	Personality         string
	RecruitingBias      string
	RecruitingBiasValue string
	WorkEthic           string
	AcademicBias        string
	PlayerID            uint
	TeamID              uint
	TeamAbbr            string
	IsRedshirting       bool
	IsRedshirt          bool
	PreviousTeamID      uint
	PreviousTeam        string
	TransferStatus      int    // 1 == Intends, 2 == Is Transferring
	TransferLikeliness  string // Low, Medium, High
	LegacyID            uint   // Either a legacy school or a legacy coach
	SeasonStats         CollegePlayerSeasonStats
	LeadingTeams        []LeadingTeams
}

func (c *TransferPlayerResponse) Map(r CollegePlayer, ovr string) {
	c.PlayerID = uint(r.ID)
	c.TeamID = uint(r.TeamID)
	c.FirstName = r.FirstName
	c.LastName = r.LastName
	c.Position = r.Position
	c.Archetype = r.Archetype
	c.State = r.State
	c.Year = r.Year
	c.IsRedshirt = r.IsRedshirt
	c.IsRedshirting = r.IsRedshirting

	var totalPoints float32 = 0
	var runningThreshold float32 = 0

	sortedProfiles := r.Profiles

	sort.Slice(sortedProfiles, func(i, j int) bool {
		return sortedProfiles[i].TotalPoints > sortedProfiles[j].TotalPoints
	})
	for _, profile := range sortedProfiles {
		if profile.RemovedFromBoard {
			continue
		}
		if runningThreshold == 0 {
			runningThreshold = float32(profile.TotalPoints) * 0.66
		}

		if float32(profile.TotalPoints) >= runningThreshold {
			totalPoints += float32(profile.TotalPoints)
		}

	}

	for i := 0; i < len(sortedProfiles); i++ {
		if sortedProfiles[i].RemovedFromBoard {
			continue
		}
		var odds float32 = 0

		if float32(sortedProfiles[i].TotalPoints) >= runningThreshold && runningThreshold > 0 {
			odds = float32(sortedProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamID:   r.Profiles[i].ProfileID,
			TeamAbbr: r.Profiles[i].TeamAbbreviation,
			Odds:     odds,
		}
		c.LeadingTeams = append(c.LeadingTeams, leadingTeam)
	}
	sort.Sort(ByLeadingPoints(c.LeadingTeams))
}

// Player Profile For the Transfer Portal?
type TransferPortalProfileResponse struct {
	ID                    uint
	SeasonID              uint
	CollegePlayerID       uint
	ProfileID             uint
	PromiseID             uint
	TeamAbbreviation      string
	TotalPoints           float64
	CurrentWeeksPoints    int
	PreviouslySpentPoints int
	SpendingCount         int
	RemovedFromBoard      bool
	RolledOnPromise       bool
	LockProfile           bool
	IsSigned              bool
	Recruiter             string
	CollegePlayer         TransferPlayerResponse `gorm:"foreignKey:CollegePlayerID"`
	Promise               CollegePromise         `gorm:"foreignKey:PromiseID"`
}

type TransferPortalResponse struct {
	Team         RecruitingTeamProfile
	TeamBoard    []TransferPortalProfileResponse
	TeamPromises []CollegePromise         // List of all promises
	Players      []TransferPlayerResponse // List of all Transfer Portal Players
	TeamList     []CollegeTeam
}

// UpdateTransferPortalBoard - Data Transfer Object from UI to API
type UpdateTransferPortalBoard struct {
	Profile SimTeamBoardResponse
	Players []TransferPortalProfileResponse
	TeamID  int
}
