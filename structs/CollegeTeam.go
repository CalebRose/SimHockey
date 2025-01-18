package structs

import "gorm.io/gorm"

type BaseTeam struct {
	TeamName         string
	Mascot           string
	Abbreviation     string
	ConferenceID     uint8
	Conference       string
	Coach            string
	City             string
	State            string
	Country          string
	ArenaID          uint16
	Arena            string
	ArenaCapacity    uint16
	RecordAttendance uint16
	FirstPlayed      uint16
	DiscordID        string
	ColorOne         string
	ColorTwo         string
	ColorThree       string
	OverallGrade     string
	OffenseGrade     string
	DefenseGrade     string
}

type ProfileAttributes struct {
	ProgramPrestige      uint8
	ProfessionalPrestige uint8
	Traditions           uint8
	Facilities           uint8
	Atmosphere           uint8
	Academics            uint8
	ConferencePrestige   uint8
	CoachRating          uint8
	SeasonMomentum       uint8
}

type CollegeTeam struct {
	gorm.Model
	BaseTeam
	IsUserCoached     bool
	IsClub            bool
	IsActive          bool
	PlayersProgressed bool
	ProfileAttributes
}

type ProfessionalTeam struct {
	gorm.Model
	BaseTeam
}

type ProfessionalTeamFranchise struct {
	TeamID                    uint
	HomeMarket                uint8   // 1 == Small, 2 == Medium, 3 == Large
	TicketPrice               float32 // Price for a generic ticket
	TicketBoxPrice            float32 // Price for a ticket in a luxury box
	TicketDemand              uint8   // How likely a fan will buy a ticket to go to a game. Is impacted by Market Type, Ticket price, and W/L record
	BoxDemand                 uint8
	Food1ID                   uint    // Customizable food item
	Food1Price                float32 // Cost of said food item
	Food2ID                   uint    // See 1
	Food2Price                float32
	Food3ID                   uint // Se 1
	Food3Price                float32
	Drink1ID                  uint // Likely water
	Drink1Price               float32
	Drink2ID                  uint // Likely soda on tap
	Drink2Price               float32
	Drink3ID                  uint // Beer
	Drink3Price               float32
	Drink4ID                  uint // Hard liquor
	Drink4Price               float32
	TeamShirtPrice            float32 // Price of generic team shirt
	TeamHatPrice              float32 // Price of generice baseball cap
	TeamBeaniePrice           float32 // Price of team beanie
	TeamJerseyPrice           float32 // Generic Team Jersey price
	FacilitiesMaintenanceCost float32 // Energy, cleanliness, sewage, zamboni, 3% of each game revenue
	EquipmentCost             float32 // Team equipment cost Likely 5% of each game revenue
	OperationsCost            float32 // Can be set between 1% and 10%. This is what you pay your facility employees.
	BathroomPrice             float32 // Give users the option to price customers
	FanHappinessRating        uint8   // General happiness of fanbase. Wins and losses impact this, along with merchandise/ticket prices
	AtmosphereRating          uint8   // General
	FoodRating                uint8   // General rating of the food selected & prices
	EmployeeHappinessRating   uint8   // How happy your employees are with you. Reflected by operactions cost
}

type Arena struct {
	gorm.Model
	Name             string
	TeamID           uint
	City             string
	State            string
	Country          string
	Capacity         uint16
	RecordAttendance uint16
}
