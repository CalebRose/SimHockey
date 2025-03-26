package structs

type BootstrapData struct {
	CollegeTeam               CollegeTeam
	AllCollegeTeams           []CollegeTeam
	CollegeStandings          []CollegeStandings
	CollegeRosterMap          map[uint][]CollegePlayer
	Recruits                  []Croot
	RecruitProfiles           []RecruitPlayerProfile
	TeamProfileMap            map[uint]RecruitingTeamProfile
	PortalPlayers             []CollegePlayer
	CollegeInjuryReport       []CollegePlayer
	CollegeNews               []NewsLog
	CollegeNotifications      []Notification
	AllCollegeGames           []CollegeGame
	CollegeTeamLineups        []CollegeLineup
	CollegeTeamShootoutLineup CollegeShootoutLineup
	// Player Profiles by Team?
	// Portal profiles?
	ProTeam               ProfessionalTeam
	AllProTeams           []ProfessionalTeam
	ProStandings          []ProfessionalStandings
	ProRosterMap          map[uint][]ProfessionalPlayer
	CapsheetMap           map[uint]ProCapsheet
	FreeAgency            FreeAgencyResponse
	ProInjuryReport       []ProfessionalPlayer
	ProNews               []NewsLog
	ProNotifications      []Notification
	AllProGames           []ProfessionalGame
	ProTeamLineups        []ProfessionalLineup
	ProTeamShootoutLineup ProfessionalShootoutLineup
	FaceData              map[uint]FaceDataResponse
}
