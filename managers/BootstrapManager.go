package managers

import (
	"log"
	"strconv"
	"sync"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

type BootstrapDataNews struct {
	CollegeNews []structs.NewsLog
	ProNews     []structs.NewsLog
}

type BootstrapDataStats struct {
	// PostSeasonAwards       AwardsList
	HistoricCollegePlayers []structs.HistoricCollegePlayer
	RetiredPlayers         []structs.RetiredPlayer
}

func GetAllTeamsData() structs.BootstrapData {
	var wg sync.WaitGroup

	var (
		allCollegeTeams []structs.CollegeTeam
		allProTeams     []structs.ProfessionalTeam
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		allCollegeTeams = GetAllCollegeTeams()
	}()
	go func() {
		defer wg.Done()
		allProTeams = GetAllProfessionalTeams()
	}()

	wg.Wait()

	return structs.BootstrapData{
		AllCollegeTeams: allCollegeTeams,
		AllProTeams:     allProTeams,
	}
}

func GetBootstrapData(collegeID, proID string) structs.BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		collegeTeam            structs.CollegeTeam
		collegeStandings       []structs.CollegeStandings
		collegePlayerMap       map[uint][]structs.CollegePlayer
		teamProfileMap         map[uint]*structs.RecruitingTeamProfile
		recruitProfiles        []structs.RecruitPlayerProfile
		portalPlayers          []structs.CollegePlayer
		transferPortalProfiles []structs.TransferPortalProfile
		injuredCollegePlayers  []structs.CollegePlayer
		chlGoals               []structs.CollegePlayer
		chlAssists             []structs.CollegePlayer
		chlSaves               []structs.CollegePlayer
		collegeNews            []structs.NewsLog
		collegeNotifications   []structs.Notification
		collegeGames           []structs.CollegeGame
		recruits               []structs.Croot
		collegeGameplanMap     map[uint]structs.CollegeGameplan
		collegeLineups         []structs.CollegeLineup
		collegeShootoutLineup  structs.CollegeShootoutLineup
		faceDataMap            map[uint]structs.FaceDataResponse
		collegePromises        []structs.CollegePromise
	)

	// Professional Data
	var (
		proTeam             structs.ProfessionalTeam
		proStandings        []structs.ProfessionalStandings
		proRosterMap        map[uint][]structs.ProfessionalPlayer
		capsheetMap         map[uint]structs.ProCapsheet
		injuredProPlayers   []structs.ProfessionalPlayer
		affiliatePlayers    []structs.ProfessionalPlayer
		phlGoals            []structs.ProfessionalPlayer
		phlAssists          []structs.ProfessionalPlayer
		phlSaves            []structs.ProfessionalPlayer
		freeAgentOffers     []structs.FreeAgencyOffer
		waiverWireOffers    []structs.WaiverOffer
		proNews             []structs.NewsLog
		proNotifications    []structs.Notification
		proGames            []structs.ProfessionalGame
		proGameplanMap      map[uint]structs.ProGameplan
		proLineups          []structs.ProfessionalLineup
		proShootoutLineup   structs.ProfessionalShootoutLineup
		contractMap         map[uint]structs.ProContract
		extensionMap        map[uint]structs.ExtensionOffer
		tradeProposalMap    map[uint][]structs.TradeProposal
		tradePreferencesMap map[uint]structs.TradePreferences
		draftPicks          map[uint][]structs.DraftPick
	)
	ts := GetTimestamp()
	_, collegeGameType := ts.GetCurrentGameType(true)
	_, proGameType := ts.GetCurrentGameType(false)
	seasonID := strconv.Itoa(int(ts.SeasonID))

	// Start concurrent queries

	wg.Add(1)
	go func() {
		defer wg.Done()
		collegePlayers := GetAllCollegePlayers()
		chlStats := GetCollegePlayerSeasonStatsBySeason(seasonID, collegeGameType)

		var teamID uint = 0
		if len(collegeID) > 0 && collegeID != "0" {
			parsed, _ := strconv.Atoi(collegeID)
			teamID = uint(parsed)
		}

		mu.Lock()
		collegePlayerMap = MakeCollegePlayerMapByTeamID(collegePlayers)
		collegePlayerIndvMap := MakeCollegePlayerMap(collegePlayers)
		chlGoals = GetCollegeOrderedListByStatType("GOALS", teamID, chlStats, collegePlayerIndvMap)
		chlAssists = GetCollegeOrderedListByStatType("ASSISTS", teamID, chlStats, collegePlayerIndvMap)
		chlSaves = GetCollegeOrderedListByStatType("SAVES", teamID, chlStats, collegePlayerIndvMap)
		injuredCollegePlayers = MakeCollegeInjuryList(collegePlayers)
		portalPlayers = MakeCollegePortalList(collegePlayers)
		mu.Unlock()
	}()

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(4)
		go func() {
			defer wg.Done()
			mu.Lock()
			collegeTeam = GetCollegeTeamByTeamID(collegeID)
			collegeTeam.UpdateLatestInstance()
			repository.SaveCollegeTeamRecord(dbprovider.GetInstance().GetDB(), collegeTeam)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			teamProfiles := repository.FindTeamRecruitingProfiles(false)
			teamProfileMap = MakeTeamProfileMap(teamProfiles)
		}()
		go func() {
			defer wg.Done()
			recruits = GetAllCrootRecords()
		}()
		go func() {
			defer wg.Done()
			transferPortalProfiles = repository.FindTransferPortalProfileRecords(repository.TransferPortalQuery{RemovedFromBoard: "N"})
		}()
		wg.Wait()
		wg.Add(3)
		go func() {
			defer wg.Done()
			collegeGames = GetCollegeGamesBySeasonID("", ts.IsPreseason)
		}()

		go func() {
			defer wg.Done()
			collegeNotifications = GetNotificationByTeamIDAndLeague("CHL", collegeID)
		}()
		go func() {
			defer wg.Done()
			collegeStandings = GetAllCollegeStandingsBySeasonID("")
		}()
		wg.Wait()
		wg.Add(4)
		go func() {
			defer wg.Done()
			collegeLineups = GetCollegeLineupsByTeamID(collegeID)
		}()
		go func() {
			defer wg.Done()
			recruitProfiles = repository.FindRecruitPlayerProfileRecords(collegeID, "", false, false, true)
		}()
		go func() {
			defer wg.Done()
			collegeShootoutLineup = GetCollegeShootoutLineupByTeamID(collegeID)
		}()
		go func() {
			defer wg.Done()
			collegeGameplans := repository.FindCollegeGameplanRecords()
			collegeGameplanMap = MakeCollegeGameplanMap(collegeGameplans)
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			collegePromises = repository.FindCollegePromiseRecords(repository.TransferPortalQuery{IsActive: "Y"})
		}()

		wg.Wait()

	}

	// Pros
	if len(proID) > 0 && proID != "0" {
		wg.Add(4)
		go func() {
			defer wg.Done()
			mu.Lock()
			proTeam = GetProTeamByTeamID(proID)
			proTeam.UpdateLatestInstance()
			repository.SaveProTeamRecord(dbprovider.GetInstance().GetDB(), proTeam)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			proPlayers := GetAllProPlayers()
			phlStats := GetProPlayerSeasonStatsBySeason(seasonID, proGameType)
			mu.Lock()
			proRosterMap = MakeProfessionalPlayerMapByTeamID(proPlayers)
			proPlayerMap := MakeProfessionalPlayerMap(proPlayers)
			phlGoals = GetProOrderedListByStatType("GOALS", proTeam.ID, phlStats, proPlayerMap)
			phlAssists = GetProOrderedListByStatType("ASSISTS", proTeam.ID, phlStats, proPlayerMap)
			phlSaves = GetProOrderedListByStatType("SAVES", proTeam.ID, phlStats, proPlayerMap)
			affiliatePlayers = MakeProAffiliateList(proPlayers)
			injuredProPlayers = MakeProInjuryList(proPlayers)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			proGames = GetProfessionalGamesBySeasonID("", ts.IsPreseason)
		}()
		go func() {
			defer wg.Done()
			proNotifications = GetNotificationByTeamIDAndLeague("PHL", proID)
		}()
		wg.Wait()
		wg.Add(4)
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
			allDraftPicks := GetAllDraftPicksBySeasonID("")
			draftPicks = MakeDraftPickMapByTeamID(allDraftPicks)
		}()
		go func() {
			defer wg.Done()
			proGameplans := repository.FindProfessionalGameplanRecords()
			proGameplanMap = MakeProGameplanMap(proGameplans)
		}()
		wg.Wait()
		wg.Add(8)
		go func() {
			defer wg.Done()
			proLineups = GetProLineupsByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			proShootoutLineup = GetProShootoutLineupByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			contractMap = GetContractMap()
		}()

		go func() {
			defer wg.Done()
			extensionMap = GetExtensionMap()
		}()
		go func() {
			defer wg.Done()
			freeAgentOffers = repository.FindAllFreeAgentOffers("", "", "", true)
		}()
		go func() {
			defer wg.Done()
			waiverWireOffers = repository.FindAllWaiverWireOffers("", "", "", true)
		}()
		go func() {
			defer wg.Done()
			tradeProposalMap = GetTradeProposalsMap()
		}()
		go func() {
			defer wg.Done()
			tradePreferencesMap = GetTradePreferencesMap()
		}()
		wg.Wait()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		faceDataMap = GetAllFaces()
	}()

	wg.Wait()

	return structs.BootstrapData{
		CollegeTeam:               collegeTeam,
		CollegeStandings:          collegeStandings,
		CollegeRosterMap:          collegePlayerMap,
		Recruits:                  recruits,
		RecruitProfiles:           recruitProfiles,
		TeamProfileMap:            teamProfileMap,
		PortalPlayers:             portalPlayers,
		TransferPortalProfiles:    transferPortalProfiles,
		CollegePromises:           collegePromises,
		CollegeInjuryReport:       injuredCollegePlayers,
		CollegeNews:               collegeNews,
		CollegeNotifications:      collegeNotifications,
		CHLGameplanMap:            collegeGameplanMap,
		CollegeTeamLineups:        collegeLineups,
		CollegeTeamShootoutLineup: collegeShootoutLineup,
		AllCollegeGames:           collegeGames,
		ProTeam:                   proTeam,
		ProStandings:              proStandings,
		ProRosterMap:              proRosterMap,
		CapsheetMap:               capsheetMap,
		WaiverWireOffers:          waiverWireOffers,
		FreeAgentOffers:           freeAgentOffers,
		AffiliatePlayers:          affiliatePlayers,
		ProInjuryReport:           injuredProPlayers,
		ProNews:                   proNews,
		ProNotifications:          proNotifications,
		AllProGames:               proGames,
		PHLGameplanMap:            proGameplanMap,
		ProTeamLineups:            proLineups,
		ProTeamShootoutLineup:     proShootoutLineup,
		FaceData:                  faceDataMap,
		ContractMap:               contractMap,
		ExtensionMap:              extensionMap,
		TopCHLGoals:               chlGoals,
		TopCHLAssists:             chlAssists,
		TopCHLSaves:               chlSaves,
		TopPHLGoals:               phlGoals,
		TopPHLAssists:             phlAssists,
		TopPHLSaves:               phlSaves,
		ProTradeProposalMap:       tradeProposalMap,
		ProTradePreferenceMap:     tradePreferencesMap,
		DraftPicks:                draftPicks,
	}
}

func GetNewsBootstrap(collegeID, proID string) BootstrapDataNews {
	var wg sync.WaitGroup

	var (
		collegeNews []structs.NewsLog
		proNews     []structs.NewsLog
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllCHLNewsLogs()
			log.Println("Fetched College News Logs, count:", len(collegeNews))
		}()
		log.Println("Initiated all College data queries.")
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			proNews = GetAllPHLNewsLogs()
		}()

	}

	wg.Wait()

	return BootstrapDataNews{
		CollegeNews: collegeNews,
		ProNews:     proNews,
	}
}

func GetStatsBootstrap(collegeID, proID string) BootstrapDataStats {
	var wg sync.WaitGroup

	var (
		historicCollegePlayers []structs.HistoricCollegePlayer
		retiredPlayers         []structs.RetiredPlayer
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		historicCollegePlayers = GetAllHistoricCollegePlayers()
	}()

	if len(proID) > 0 && proID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Fetching Retired Players...")
			retiredPlayers = GetAllRetiredPlayers()
			log.Println("Fetched Retired Players, count:", len(retiredPlayers))
		}()
	}

	wg.Wait()

	return BootstrapDataStats{
		HistoricCollegePlayers: historicCollegePlayers,
		RetiredPlayers:         retiredPlayers,
	}
}

func GetLineUpBootstrapData(collegeID, proID string) structs.BootstrapData {
	var wg sync.WaitGroup

	var (
		collegeLineupsMap         map[uint][]structs.CollegeLineup
		proLineupsMap             map[uint][]structs.ProfessionalLineup
		collegeShootoutLineupsMap map[uint]structs.CollegeShootoutLineup
		proShootoutLineupsMap     map[uint]structs.ProfessionalShootoutLineup
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(2)
		go func() {
			defer wg.Done()
			collegeLineups := repository.FindAllCollegeLineups()
			collegeLineupsMap = MakeCollegeLineupMap(collegeLineups)
		}()

		go func() {
			defer wg.Done()
			collegeShootoutLineups := repository.FindAllCollegeShootoutLineups()
			collegeShootoutLineupsMap = MakeCollegeShootoutLineupMap(collegeShootoutLineups)
		}()
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(2)
		go func() {
			defer wg.Done()
			proShootoutLineups := repository.FindAllProShootoutLineups()
			proShootoutLineupsMap = MakeProfessionalShootoutLineupMap(proShootoutLineups)
		}()

		go func() {
			defer wg.Done()
			proLineups := repository.FindAllProLineups()
			proLineupsMap = MakeProfessionalLineupMap(proLineups)
		}()
	}
	wg.Wait()
	return structs.BootstrapData{
		CollegeLineupsMap:         collegeLineupsMap,
		ProLineupsMap:             proLineupsMap,
		CollegeShootoutLineupsMap: collegeShootoutLineupsMap,
		ProShootoutLineupsMap:     proShootoutLineupsMap,
	}
}

func GetScheduleBootstrap(collegeID, username string) structs.BootstrapData {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	var wg sync.WaitGroup

	var (
		hockeyInvitationals        []structs.HockeyInvitational
		hockeyInvitationalRequests []structs.HockeyInvitationalRequest
		collegeGameRequests        []structs.CHLGameRequest
		poll                       structs.CollegePollSubmission
		officialPolls              []structs.CollegePollOfficial
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(5)

		go func() {
			defer wg.Done()
			hockeyInvitationals = repository.FindInvitationalRecords(repository.SchedulerQuery{
				SeasonID: seasonID,
			})
		}()

		go func() {
			defer wg.Done()
			hockeyInvitationalRequests = repository.FindInvitationalRequestRecords(repository.SchedulerQuery{
				SeasonID: seasonID,
			})
		}()
		go func() {
			defer wg.Done()
			collegeGameRequests = repository.FindCHLGameRequestRecords(repository.SchedulerQuery{
				SeasonID: seasonID,
			})
		}()

		go func() {
			defer wg.Done()
			poll = GetPollSubmissionByUsernameWeekAndSeason(username)
		}()
		go func() {
			defer wg.Done()
			officialPolls = GetOfficialPollBySeasonID("")
		}()
		wg.Wait()
	}

	return structs.BootstrapData{
		HockeyInvitationals:        hockeyInvitationals,
		HockeyInvitationalRequests: hockeyInvitationalRequests,
		CollegeGameRequests:        collegeGameRequests,
		CollegePoll:                poll,
		OfficialPolls:              officialPolls,
	}
}
