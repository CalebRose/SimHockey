package managers

import (
	"strconv"
	"sync"

	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetBootstrapData(collegeID, proID string) structs.BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		allCollegeTeams       []structs.CollegeTeam
		collegeTeam           structs.CollegeTeam
		collegeStandings      []structs.CollegeStandings
		collegePlayerMap      map[uint][]structs.CollegePlayer
		teamProfileMap        map[uint]*structs.RecruitingTeamProfile
		recruitProfiles       []structs.RecruitPlayerProfile
		portalPlayers         []structs.CollegePlayer
		injuredCollegePlayers []structs.CollegePlayer
		collegeNews           []structs.NewsLog
		collegeNotifications  []structs.Notification
		collegeGames          []structs.CollegeGame
		recruits              []structs.Croot
		collegeLineups        []structs.CollegeLineup
		collegeShootoutLineup structs.CollegeShootoutLineup
		faceDataMap           map[uint]structs.FaceDataResponse
	)

	// Professional Data
	var (
		proTeam           structs.ProfessionalTeam
		allProTeams       []structs.ProfessionalTeam
		proStandings      []structs.ProfessionalStandings
		proRosterMap      map[uint][]structs.ProfessionalPlayer
		capsheetMap       map[uint]structs.ProCapsheet
		freeAgency        structs.FreeAgencyResponse
		injuredProPlayers []structs.ProfessionalPlayer
		proNews           []structs.NewsLog
		proNotifications  []structs.Notification
		proGames          []structs.ProfessionalGame
		proLineups        []structs.ProfessionalLineup
		proShootoutLineup structs.ProfessionalShootoutLineup
	)

	freeAgencyCh := make(chan structs.FreeAgencyResponse, 1)

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	// Start concurrent queries
	wg.Add(2)
	go func() {
		defer wg.Done()
		allCollegeTeams = GetAllCollegeTeams()
	}()
	go func() {
		defer wg.Done()
		allProTeams = GetAllProfessionalTeams()
	}()

	if len(collegeID) > 0 {
		wg.Add(6)
		go func() {
			defer wg.Done()
			collegeTeam = GetCollegeTeamByTeamID(collegeID)
		}()
		go func() {
			defer wg.Done()
			teamProfiles := repository.FindTeamRecruitingProfiles(false)
			teamProfileMap = MakeTeamProfileMap(teamProfiles)
		}()
		go func() {
			defer wg.Done()
			collegePlayers := GetAllCollegePlayers()
			mu.Lock()
			collegePlayerMap = MakeCollegePlayerMapByTeamID(collegePlayers)
			injuredCollegePlayers = MakeCollegeInjuryList(collegePlayers)
			portalPlayers = MakeCollegePortalList(collegePlayers)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			recruits = GetAllCrootRecords()
		}()
		go func() {
			defer wg.Done()
			collegeGames = GetCollegeGamesBySeasonID(seasonID)
		}()
		go func() {
			defer wg.Done()
			collegeNews = GetAllCHLNewsLogs()
		}()
		wg.Wait()
		wg.Add(5)
		go func() {
			defer wg.Done()
			collegeNotifications = GetNotificationByTeamIDAndLeague("CHL", collegeID)
		}()
		go func() {
			defer wg.Done()
			collegeStandings = GetAllCollegeStandingsBySeasonID(seasonID)
		}()
		go func() {
			defer wg.Done()
			collegeLineups = GetCollegeLineupsByTeamID(collegeID)
		}()
		go func() {
			defer wg.Done()
			recruitProfiles = repository.FindRecruitPlayerProfileRecords(collegeID, "", false, false)
		}()
		go func() {
			defer wg.Done()
			collegeShootoutLineup = GetCollegeShootoutLineupByTeamID(collegeID)
		}()
		wg.Wait()
	}

	// Pros
	if len(proID) > 0 {
		wg.Add(9)
		go func() {
			defer wg.Done()
			proTeam = GetProTeamByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			proPlayers := GetAllProPlayers()
			mu.Lock()
			proRosterMap = MakeProfessionalPlayerMapByTeamID(proPlayers)
			injuredProPlayers = MakeProInjuryList(proPlayers)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			proGames = GetProfessionalGamesBySeasonID("")
		}()
		go func() {
			defer wg.Done()
			proNews = GetAllPHLNewsLogs()
		}()
		go func() {
			defer wg.Done()
			proNotifications = GetNotificationByTeamIDAndLeague("PHL", proID)
		}()
		go func() {
			defer wg.Done()
			capsheetMap = GetProCapsheetMap()
		}()
		go func() {
			defer wg.Done()
			proStandings = GetAllProfessionalStandingsBySeasonID("")
		}()
		go func() {
			defer wg.Done()
			proLineups = GetProLineupsByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			proShootoutLineup = GetProShootoutLineupByTeamID(proID)
		}()
		go GetAllAvailableProPlayers(proID, freeAgencyCh)
		wg.Wait()
	}

	// Wait for all goroutines to finish

	// Receive Free Agency data from channel
	if len(proID) > 0 {
		freeAgency = <-freeAgencyCh
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		faceDataMap = GetAllFaces()
	}()

	wg.Wait()

	return structs.BootstrapData{
		AllCollegeTeams:           allCollegeTeams,
		AllProTeams:               allProTeams,
		CollegeTeam:               collegeTeam,
		CollegeStandings:          collegeStandings,
		CollegeRosterMap:          collegePlayerMap,
		Recruits:                  recruits,
		RecruitProfiles:           recruitProfiles,
		TeamProfileMap:            teamProfileMap,
		PortalPlayers:             portalPlayers,
		CollegeInjuryReport:       injuredCollegePlayers,
		CollegeNews:               collegeNews,
		CollegeNotifications:      collegeNotifications,
		CollegeTeamLineups:        collegeLineups,
		CollegeTeamShootoutLineup: collegeShootoutLineup,
		AllCollegeGames:           collegeGames,
		ProTeam:                   proTeam,
		ProStandings:              proStandings,
		ProRosterMap:              proRosterMap,
		CapsheetMap:               capsheetMap,
		FreeAgency:                freeAgency,
		ProInjuryReport:           injuredProPlayers,
		ProNews:                   proNews,
		ProNotifications:          proNotifications,
		AllProGames:               proGames,
		ProTeamLineups:            proLineups,
		ProTeamShootoutLineup:     proShootoutLineup,
		FaceData:                  faceDataMap,
	}
}
