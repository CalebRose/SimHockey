package structs

type GameDTO struct {
	GameID        uint
	GameInfo      BaseGame
	HomeStrategy  PlayBookDTO
	AwayStrategy  PlayBookDTO
	IsCollegeGame bool
	Attendance    uint32
	Capacity      uint32
}

type PlayBookDTO struct {
	Forwards           []BaseLineup
	Defenders          []BaseLineup
	Goalies            []BaseLineup
	CollegeRoster      []CollegePlayer
	ProfessionalRoster []ProfessionalPlayer
	ShootoutLineup     ShootoutPlayerIDs
	Gameplan           BaseGameplan
}

type GameResultDTO struct {
}
