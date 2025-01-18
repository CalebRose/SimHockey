package structs

type GameDTO struct {
	GameInfo      BaseGame
	HomeStrategy  PlayBookDTO
	AwayStrategy  PlayBookDTO
	IsCollegeGame bool
	Attendance    uint32
}

type PlayBookDTO struct {
	Forwards           []BaseLineup
	Defenders          []BaseLineup
	Goalies            []BaseLineup
	CollegeRoster      []CollegePlayer
	ProfessionalRoster []ProfessionalPlayer
	ShootoutLineup     ShootoutPlayerIDs
}

type GameResultDTO struct {
}
