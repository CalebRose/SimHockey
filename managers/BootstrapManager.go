package managers

import (
	"strconv"
	"sync"

	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

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
		collegeTeam           structs.CollegeTeam
		collegeStandings      []structs.CollegeStandings
		collegePlayerMap      map[uint][]structs.CollegePlayer
		teamProfileMap        map[uint]*structs.RecruitingTeamProfile
		recruitProfiles       []structs.RecruitPlayerProfile
		portalPlayers         []structs.CollegePlayer
		injuredCollegePlayers []structs.CollegePlayer
		chlGoals              []structs.CollegePlayer
		chlAssists            []structs.CollegePlayer
		chlSaves              []structs.CollegePlayer
		collegeNews           []structs.NewsLog
		collegeNotifications  []structs.Notification
		collegeGames          []structs.CollegeGame
		recruits              []structs.Croot
		collegeGameplan       structs.CollegeGameplan
		collegeLineups        []structs.CollegeLineup
		collegeShootoutLineup structs.CollegeShootoutLineup
		faceDataMap           map[uint]structs.FaceDataResponse
		poll                  structs.CollegePollSubmission
		officialPolls         []structs.CollegePollOfficial
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
		proGameplan         structs.ProGameplan
		proLineups          []structs.ProfessionalLineup
		proShootoutLineup   structs.ProfessionalShootoutLineup
		contractMap         map[uint]structs.ProContract
		extensionMap        map[uint]structs.ExtensionOffer
		tradeProposalMap    map[uint][]structs.TradeProposal
		tradePreferencesMap map[uint]structs.TradePreferences
		draftPicks          []structs.DraftPick
	)
	ts := GetTimestamp()
	_, collegeGameType := ts.GetCurrentGameType(true)
	_, proGameType := ts.GetCurrentGameType(false)
	seasonID := strconv.Itoa(int(ts.SeasonID))

	// Start concurrent queries

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(4)
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
			chlStats := GetCollegePlayerSeasonStatsBySeason(seasonID, collegeGameType)
			mu.Lock()
			collegePlayerMap = MakeCollegePlayerMapByTeamID(collegePlayers)
			collegePlayerIndvMap := MakeCollegePlayerMap(collegePlayers)
			chlGoals = GetCollegeOrderedListByStatType("GOALS", collegeTeam.ID, chlStats, collegePlayerIndvMap)
			chlAssists = GetCollegeOrderedListByStatType("ASSISTS", collegeTeam.ID, chlStats, collegePlayerIndvMap)
			chlSaves = GetCollegeOrderedListByStatType("SAVES", collegeTeam.ID, chlStats, collegePlayerIndvMap)
			injuredCollegePlayers = MakeCollegeInjuryList(collegePlayers)
			portalPlayers = MakeCollegePortalList(collegePlayers)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			recruits = GetAllCrootRecords()
		}()
		wg.Wait()
		wg.Add(5)
		go func() {
			defer wg.Done()
			collegeGames = GetCollegeGamesBySeasonID(seasonID, ts.IsPreseason)
		}()
		go func() {
			defer wg.Done()
			collegeNews = GetAllCHLNewsLogs()
		}()

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
			officialPolls = GetOfficialPollBySeasonID(seasonID)
		}()
		wg.Wait()
		wg.Add(5)
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
			collegeGameplan = repository.FindCollegeGameplanRecord(collegeID)
		}()
		go func() {
			defer wg.Done()
			poll = GetPollSubmissionByUsernameWeekAndSeason(collegeTeam.Coach)
		}()
		wg.Wait()

	}

	// Pros
	if len(proID) > 0 && proID != "0" {
		wg.Add(5)
		go func() {
			defer wg.Done()
			proTeam = GetProTeamByTeamID(proID)
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
			proNews = GetAllPHLNewsLogs()
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
			draftPicks = GetAllDraftPicksBySeasonID("")
		}()
		go func() {
			defer wg.Done()
			proGameplan = repository.FindProGameplanRecord(proID)
		}()

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
		CollegeInjuryReport:       injuredCollegePlayers,
		CollegeNews:               collegeNews,
		CollegeNotifications:      collegeNotifications,
		CHLGameplan:               collegeGameplan,
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
		PHLGameplan:               proGameplan,
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
		CollegePoll:               poll,
		OfficialPolls:             officialPolls,
	}
}
