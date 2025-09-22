package structs

type BootstrapData struct {
	CollegeTeam               CollegeTeam
	AllCollegeTeams           []CollegeTeam
	CollegeStandings          []CollegeStandings
	CollegeRosterMap          map[uint][]CollegePlayer
	Recruits                  []Croot
	RecruitProfiles           []RecruitPlayerProfile
	TeamProfileMap            map[uint]*RecruitingTeamProfile
	PortalPlayers             []CollegePlayer
	TransferPortalProfiles    []TransferPortalProfile
	CollegeInjuryReport       []CollegePlayer
	CollegeNews               []NewsLog
	CollegeNotifications      []Notification
	AllCollegeGames           []CollegeGame
	CHLGameplan               CollegeGameplan
	CollegeTeamLineups        []CollegeLineup
	CollegeTeamShootoutLineup CollegeShootoutLineup
	TopCHLGoals               []CollegePlayer
	TopCHLAssists             []CollegePlayer
	TopCHLSaves               []CollegePlayer
	CollegePoll               CollegePollSubmission
	OfficialPolls             []CollegePollOfficial
	// Player Profiles by Team?
	// Portal profiles?
	ProTeam               ProfessionalTeam
	AllProTeams           []ProfessionalTeam
	ProStandings          []ProfessionalStandings
	ProRosterMap          map[uint][]ProfessionalPlayer
	AffiliatePlayers      []ProfessionalPlayer
	TopPHLGoals           []ProfessionalPlayer
	TopPHLAssists         []ProfessionalPlayer
	TopPHLSaves           []ProfessionalPlayer
	FreeAgentOffers       []FreeAgencyOffer
	WaiverWireOffers      []WaiverOffer
	CapsheetMap           map[uint]ProCapsheet
	ProInjuryReport       []ProfessionalPlayer
	ProNews               []NewsLog
	ProNotifications      []Notification
	AllProGames           []ProfessionalGame
	PHLGameplan           ProGameplan
	ProTeamLineups        []ProfessionalLineup
	ProTeamShootoutLineup ProfessionalShootoutLineup
	FaceData              map[uint]FaceDataResponse
	ContractMap           map[uint]ProContract
	ExtensionMap          map[uint]ExtensionOffer
	ProTradeProposalMap   map[uint][]TradeProposal
	ProTradePreferenceMap map[uint]TradePreferences
	DraftPicks            []DraftPick
}
