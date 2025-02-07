package structs

import "gorm.io/gorm"

type ProCapsheet struct {
	gorm.Model
	TeamID   uint
	Y1Bonus  float32
	Y1Salary float32
	Y1CapHit float32
	Y2Bonus  float32
	Y2Salary float32
	Y2CapHit float32
	Y3Bonus  float32
	Y3Salary float32
	Y3CapHit float32
	Y4Bonus  float32
	Y4Salary float32
	Y4CapHit float32
	Y5Bonus  float32
	Y5Salary float32
	Y5CapHit float32
}

func (nc *ProCapsheet) AssignCapsheet(id uint) {
	nc.ID = id
	nc.TeamID = id
}

func (nc *ProCapsheet) ResetCapsheet() {
	nc.Y1Bonus = 0
	nc.Y1Salary = 0
	nc.Y2Bonus = 0
	nc.Y2Salary = 0
	nc.Y3Bonus = 0
	nc.Y3Salary = 0
	nc.Y4Bonus = 0
	nc.Y4Salary = 0
	nc.Y5Bonus = 0
	nc.Y5Salary = 0
}

func (nc *ProCapsheet) AddContractToCapsheet(contract ProContract) {
	nc.Y1Bonus += contract.Y1Bonus
	nc.Y1Salary += contract.Y1BaseSalary
	nc.Y2Bonus += contract.Y2Bonus
	nc.Y2Salary += contract.Y2BaseSalary
	nc.Y3Bonus += contract.Y3Bonus
	nc.Y3Salary += contract.Y3BaseSalary
	nc.Y4Bonus += contract.Y4Bonus
	nc.Y4Salary += contract.Y4BaseSalary
	nc.Y5Bonus += contract.Y5Bonus
	nc.Y5Salary += contract.Y5BaseSalary
}

func (nc *ProCapsheet) SubtractFromCapsheet(contract ProContract) {
	nc.Y1CapHit += contract.Y1Bonus
	nc.Y1Bonus -= contract.Y1Bonus
	nc.Y2Bonus -= contract.Y2Bonus
	nc.Y3Bonus -= contract.Y3Bonus
	nc.Y4Bonus -= contract.Y4Bonus
	nc.Y5Bonus -= contract.Y5Bonus
	nc.Y1Salary -= contract.Y1BaseSalary
	nc.Y2Salary -= contract.Y2BaseSalary
	nc.Y3Salary -= contract.Y3BaseSalary
	nc.Y4Salary -= contract.Y4BaseSalary
	nc.Y5Salary -= contract.Y5BaseSalary
}

func (nc *ProCapsheet) CutPlayerFromCapsheet(contract ProContract) {
	nc.Y1CapHit += contract.Y1Bonus + contract.Y2Bonus + contract.Y3Bonus + contract.Y4Bonus + contract.Y5Bonus
	nc.Y1Bonus -= contract.Y1Bonus
	nc.Y2Bonus -= contract.Y2Bonus
	nc.Y3Bonus -= contract.Y3Bonus
	nc.Y4Bonus -= contract.Y4Bonus
	nc.Y5Bonus -= contract.Y5Bonus
	nc.Y1Salary -= contract.Y1BaseSalary
	nc.Y2Salary -= contract.Y2BaseSalary
	nc.Y3Salary -= contract.Y3BaseSalary
	nc.Y4Salary -= contract.Y4BaseSalary
	nc.Y5Salary -= contract.Y5BaseSalary
}

func (nc *ProCapsheet) SubtractFromCapsheetViaTrade(contract ProContract) {
	nc.Y1CapHit += contract.Y1Bonus + contract.Y2Bonus + contract.Y3Bonus + contract.Y4Bonus + contract.Y5Bonus
	nc.Y1Bonus -= contract.Y1Bonus
	nc.Y2Bonus -= contract.Y2Bonus
	nc.Y3Bonus -= contract.Y3Bonus
	nc.Y4Bonus -= contract.Y4Bonus
	nc.Y5Bonus -= contract.Y5Bonus
	nc.Y2Salary -= contract.Y2BaseSalary
	nc.Y3Salary -= contract.Y3BaseSalary
	nc.Y4Salary -= contract.Y4BaseSalary
	nc.Y5Salary -= contract.Y5BaseSalary
}

func (nc *ProCapsheet) NegotiateSalaryDifference(SalaryDifference float32, CapHit float32) {
	nc.Y1Salary -= SalaryDifference
	nc.Y1CapHit += CapHit
}

func (nc *ProCapsheet) AddContractViaTrade(contract ProContract, differenceValue float32) {
	// nc.Y1Bonus += contract.Y1Bonus
	nc.Y1Salary += differenceValue
	nc.Y2Bonus += contract.Y2Bonus
	nc.Y2Salary += contract.Y2BaseSalary
	nc.Y3Bonus += contract.Y3Bonus
	nc.Y3Salary += contract.Y3BaseSalary
	nc.Y4Bonus += contract.Y4Bonus
	nc.Y4Salary += contract.Y4BaseSalary
	nc.Y5Bonus += contract.Y5Bonus
	nc.Y5Salary += contract.Y5BaseSalary
}

func (nc *ProCapsheet) ProgressCapsheet() {
	nc.Y1Salary = nc.Y2Salary
	nc.Y1Bonus = nc.Y2Bonus
	nc.Y1CapHit = nc.Y2CapHit
	nc.Y2Salary = nc.Y3Salary
	nc.Y2Bonus = nc.Y3Bonus
	nc.Y2CapHit = nc.Y3CapHit
	nc.Y3Salary = nc.Y4Salary
	nc.Y3Bonus = nc.Y4Bonus
	nc.Y3CapHit = nc.Y4CapHit
	nc.Y4Salary = nc.Y5Salary
	nc.Y4Bonus = nc.Y5Bonus
	nc.Y4CapHit = nc.Y5CapHit
	nc.Y5Salary = 0
	nc.Y5Bonus = 0
	nc.Y5CapHit = 0
}

type ProContract struct {
	gorm.Model
	PlayerID        int
	TeamID          uint
	Team            string
	OriginalTeamID  uint
	OriginalTeam    string
	ContractLength  int
	Y1BaseSalary    float32
	Y1Bonus         float32
	Y2BaseSalary    float32
	Y2Bonus         float32
	Y3BaseSalary    float32
	Y3Bonus         float32
	Y4BaseSalary    float32
	Y4Bonus         float32
	Y5BaseSalary    float32
	Y5Bonus         float32
	BonusPercentage float32
	ContractType    string // Pro Bowl, Starter, Veteran, New ?
	ContractValue   float32
	SigningValue    float32
	IsActive        bool
	IsComplete      bool
	IsExtended      bool
	HasProgressed   bool
	PlayerRetired   bool
	TagType         uint8
	IsTagged        bool
	IsCut           bool
}

func (c *ProContract) DeactivateContract() {
	c.IsActive = false
}

func (c *ProContract) CutContract() {
	c.IsActive = false
	c.IsCut = true
}

func (c *ProContract) ReassignTeam(TeamID uint, Team string) {
	c.TeamID = TeamID
	c.Team = Team
}

func (c *ProContract) TradePlayer(TeamID uint, Team string, percentage float32) {
	c.TeamID = TeamID
	c.Team = Team
	c.Y1BaseSalary = c.Y1BaseSalary * percentage
	c.Y1Bonus = 0
	c.Y2Bonus = 0
	c.Y3Bonus = 0
	c.Y4Bonus = 0
	c.Y5Bonus = 0
}

func (c *ProContract) ProgressContract() {
	c.Y1BaseSalary = c.Y2BaseSalary
	c.Y1Bonus = c.Y2Bonus
	c.Y2BaseSalary = c.Y3BaseSalary
	c.Y2Bonus = c.Y3Bonus
	c.Y3BaseSalary = c.Y4BaseSalary
	c.Y3Bonus = c.Y4Bonus
	c.Y4BaseSalary = c.Y5BaseSalary
	c.Y4Bonus = c.Y5Bonus
	c.Y5BaseSalary = 0
	c.Y5Bonus = 0
	c.ContractLength -= 1
	c.CalculateContract()
	c.HasProgressed = true

	if c.Y1BaseSalary == 0 && c.Y1Bonus == 0 {
		c.IsComplete = true
		c.DeactivateContract()
	}
}

func (c *ProContract) CalculateContract() {
	// Calculate Value
	y1SalaryVal := c.Y1BaseSalary * 0.8
	y1BonusVal := c.Y1Bonus * 1
	y2SalaryVal := c.Y2BaseSalary * 0.4
	y2BonusVal := c.Y2Bonus * 0.9
	y3SalaryVal := c.Y3BaseSalary * 0.2
	y3BonusVal := c.Y3Bonus * 0.8
	y4SalaryVal := c.Y4BaseSalary * 0.1
	y4BonusVal := c.Y4Bonus * 0.7
	y5SalaryVal := c.Y5BaseSalary * 0.05
	y5BonusVal := c.Y5Bonus * 0.6
	c.ContractValue = y1SalaryVal + y1BonusVal + y2SalaryVal + y2BonusVal + y3SalaryVal + y3BonusVal + y4SalaryVal + y4BonusVal + y5SalaryVal + y5BonusVal
}

func (c *ProContract) MapExtension(e ExtensionOffer) {
	c.ContractLength = e.ContractLength
	c.Y1BaseSalary = e.Y1BaseSalary
	c.Y1Bonus = e.Y1Bonus
	c.Y2BaseSalary = e.Y2BaseSalary
	c.Y2Bonus = e.Y2Bonus
	c.Y3BaseSalary = e.Y3BaseSalary
	c.Y3Bonus = e.Y3Bonus
	c.Y4BaseSalary = e.Y4BaseSalary
	c.Y4Bonus = e.Y4Bonus
	c.Y5BaseSalary = e.Y5BaseSalary
	c.Y5Bonus = e.Y5Bonus
	c.BonusPercentage = e.BonusPercentage
	c.CalculateContract()
	c.SigningValue = c.ContractValue
	c.IsActive = true
	c.IsComplete = false
	c.IsExtended = true
}

func (c *ProContract) ToggleRetirement() {
	c.PlayerRetired = true
}

func (c *ProContract) TagContract(tagType uint8, salary, bonus float32) {
	if c.ContractLength == 1 {
		c.IsTagged = true
		c.ContractLength += 1
		c.Y2BaseSalary = salary
		c.Y2Bonus = bonus
	}
}

type TagDTO struct {
	PlayerID uint
	TagType  uint8
	Position string
}

type FreeAgencyOfferDTO struct {
	ID             uint
	PlayerID       uint
	TeamID         uint
	ContractLength int
	Y1BaseSalary   float32
	Y1Bonus        float32
	Y2BaseSalary   float32
	Y2Bonus        float32
	Y3BaseSalary   float32
	Y3Bonus        float32
	Y4BaseSalary   float32
	Y4Bonus        float32
	Y5BaseSalary   float32
	Y5Bonus        float32
}

type FreeAgencyOffer struct {
	gorm.Model
	PlayerID        uint
	TeamID          uint
	ContractLength  int
	Y1BaseSalary    float32
	Y1Bonus         float32
	Y2BaseSalary    float32
	Y2Bonus         float32
	Y3BaseSalary    float32
	Y3Bonus         float32
	Y4BaseSalary    float32
	Y4Bonus         float32
	Y5BaseSalary    float32
	Y5Bonus         float32
	TotalBonus      float32
	TotalSalary     float32
	ContractValue   float32
	BonusPercentage float32
	IsActive        bool
}

func (f *FreeAgencyOffer) CalculateOffer(offer FreeAgencyOfferDTO) {
	f.PlayerID = offer.PlayerID
	f.TeamID = offer.TeamID
	f.ContractLength = offer.ContractLength
	f.Y1BaseSalary = offer.Y1BaseSalary
	f.Y1Bonus = offer.Y1Bonus
	f.Y2BaseSalary = offer.Y2BaseSalary
	f.Y2Bonus = offer.Y2Bonus
	f.Y3BaseSalary = offer.Y3BaseSalary
	f.Y3Bonus = offer.Y3Bonus
	f.Y4BaseSalary = offer.Y4BaseSalary
	f.Y4Bonus = offer.Y4Bonus
	f.Y5BaseSalary = offer.Y5BaseSalary
	f.Y5Bonus = offer.Y5Bonus
	f.IsActive = true

	// Calculate Value
	y1SalaryVal := f.Y1BaseSalary * 0.8
	y1BonusVal := f.Y1Bonus * 1
	y2SalaryVal := f.Y2BaseSalary * 0.4
	y2BonusVal := f.Y2Bonus * 0.9
	y3SalaryVal := f.Y3BaseSalary * 0.2
	y3BonusVal := f.Y3Bonus * 0.8
	y4SalaryVal := f.Y4BaseSalary * 0.1
	y4BonusVal := f.Y4Bonus * 0.7
	y5SalaryVal := f.Y5BaseSalary * 0.05
	y5BonusVal := f.Y5Bonus * 0.6
	f.ContractValue = y1SalaryVal + y1BonusVal + y2SalaryVal + y2BonusVal + y3SalaryVal + y3BonusVal + y4SalaryVal + y4BonusVal + y5SalaryVal + y5BonusVal
	f.TotalBonus = f.Y1Bonus + f.Y2Bonus + f.Y3Bonus + f.Y4Bonus + f.Y5Bonus
	f.TotalSalary = f.Y1BaseSalary + f.Y2BaseSalary + f.Y3BaseSalary + f.Y4BaseSalary + f.Y5BaseSalary
	total := f.TotalBonus + f.TotalSalary
	f.BonusPercentage = f.TotalBonus / (total)
}

func (f *FreeAgencyOffer) CancelOffer() {
	f.IsActive = false
}

func (f *FreeAgencyOffer) AssignID(id uint) {
	f.ID = id
}

// Sorting Funcs
type ByContractValue []FreeAgencyOffer

func (fo ByContractValue) Len() int      { return len(fo) }
func (fo ByContractValue) Swap(i, j int) { fo[i], fo[j] = fo[j], fo[i] }
func (fo ByContractValue) Less(i, j int) bool {
	return fo[i].ContractValue > fo[j].ContractValue
}

// Table for storing Extensions for contracted players
type ExtensionOffer struct {
	gorm.Model
	PlayerID        uint
	TeamID          uint
	SeasonID        uint
	ContractLength  int
	Y1BaseSalary    float32
	Y1Bonus         float32
	Y2BaseSalary    float32
	Y2Bonus         float32
	Y3BaseSalary    float32
	Y3Bonus         float32
	Y4BaseSalary    float32
	Y4Bonus         float32
	Y5BaseSalary    float32
	Y5Bonus         float32
	TotalBonus      float32
	TotalSalary     float32
	ContractValue   float32
	BonusPercentage float32
	Rejections      int
	IsAccepted      bool
	IsActive        bool
	IsRejected      bool
}

func (f *ExtensionOffer) AssignID(id uint) {
	f.ID = id
}

func (f *ExtensionOffer) CalculateOffer(offer FreeAgencyOfferDTO) {
	f.PlayerID = offer.PlayerID
	f.TeamID = offer.TeamID
	f.ContractLength = offer.ContractLength
	f.Y1BaseSalary = offer.Y1BaseSalary
	f.Y1Bonus = offer.Y1Bonus
	f.Y2BaseSalary = offer.Y2BaseSalary
	f.Y2Bonus = offer.Y2Bonus
	f.Y3BaseSalary = offer.Y3BaseSalary
	f.Y3Bonus = offer.Y3Bonus
	f.Y4BaseSalary = offer.Y4BaseSalary
	f.Y4Bonus = offer.Y4Bonus
	f.Y5BaseSalary = offer.Y5BaseSalary
	f.Y5Bonus = offer.Y5Bonus
	f.IsActive = true

	// Calculate Value
	y1SalaryVal := f.Y1BaseSalary * 0.8
	y1BonusVal := f.Y1Bonus * 1
	y2SalaryVal := f.Y2BaseSalary * 0.4
	y2BonusVal := f.Y2Bonus * 0.9
	y3SalaryVal := f.Y3BaseSalary * 0.2
	y3BonusVal := f.Y3Bonus * 0.8
	y4SalaryVal := f.Y4BaseSalary * 0.1
	y4BonusVal := f.Y4Bonus * 0.7
	y5SalaryVal := f.Y5BaseSalary * 0.05
	y5BonusVal := f.Y5Bonus * 0.6
	f.ContractValue = y1SalaryVal + y1BonusVal + y2SalaryVal + y2BonusVal + y3SalaryVal + y3BonusVal + y4SalaryVal + y4BonusVal + y5SalaryVal + y5BonusVal
	f.TotalBonus = f.Y1Bonus + f.Y2Bonus + f.Y3Bonus + f.Y4Bonus + f.Y5Bonus
	f.TotalSalary = f.Y1BaseSalary + f.Y2BaseSalary + f.Y3BaseSalary + f.Y4BaseSalary + f.Y5BaseSalary
	total := f.TotalBonus + f.TotalSalary
	f.BonusPercentage = f.TotalBonus / (total)
}

func (f *ExtensionOffer) AcceptOffer() {
	f.IsAccepted = true
	f.CancelOffer()
}

func (f *ExtensionOffer) DeclineOffer(week int) {
	f.Rejections += 1
	if f.Rejections > 2 || week >= 23 {
		f.IsRejected = true
	}
}

func (f *ExtensionOffer) CancelOffer() {
	f.IsActive = false
}

type WaiverOfferDTO struct {
	ID          uint
	PlayerID    uint
	TeamID      uint
	Team        string
	WaiverOrder uint
	IsActive    bool
}

type WaiverOffer struct {
	ID          uint
	PlayerID    uint
	TeamID      uint
	Team        string
	WaiverOrder uint
	IsActive    bool
}

func (wo *WaiverOffer) AssignID(id uint) {
	wo.ID = id
}

func (wo *WaiverOffer) AssignNewWaiverOrder(val uint) {
	wo.WaiverOrder = val
}

func (wo *WaiverOffer) Map(offer WaiverOfferDTO) {
	wo.TeamID = offer.TeamID
	wo.Team = offer.Team
	wo.PlayerID = offer.PlayerID
	wo.WaiverOrder = offer.WaiverOrder
	wo.IsActive = true
}

func (wo *WaiverOffer) DeactivateWaiverOffer() {
	wo.IsActive = false
}

type FreeAgencyResponse struct {
	FreeAgents    []ProfessionalPlayer
	WaiverPlayers []ProfessionalPlayer
	PracticeSquad []ProfessionalPlayer
	TeamOffers    []FreeAgencyOffer
	RosterCount   uint
}

type FreeAgentResponse struct {
	ID uint
	BasePlayer
	SeasonStats ProfessionalPlayerSeasonStats
	Offers      []FreeAgencyOffer
}

type WaiverWirePlayerResponse struct {
	ID uint
	BasePlayer
	SeasonStats  ProfessionalPlayerSeasonStats
	WaiverOffers []WaiverOffer
	Contract     ProContract
}
