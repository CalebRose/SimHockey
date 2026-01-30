package managers

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func ImportDraftPicksForSeason(seasonID, season uint) {
	db := dbprovider.GetInstance().GetDB()

	proTeams := GetAllProfessionalTeams()

	rand.Shuffle(len(proTeams), func(i, j int) {
		proTeams[i], proTeams[j] = proTeams[j], proTeams[i]
	})

	rounds := 7
	picksToUpload := []structs.DraftPick{}

	for i := 1; i <= rounds; i++ {
		for idx, team := range proTeams {
			pick := structs.DraftPick{
				SeasonID:    seasonID,
				Season:      season,
				DrafteeID:   0,
				DraftRound:  uint(i),
				DraftNumber: uint(idx) + 1,
				TeamID:      team.ID,
				Team:        team.Abbreviation,
				DraftValue:  float64(rounds) - (float64(i) - 1),
			}

			picksToUpload = append(picksToUpload, pick)
		}
	}

	repository.CreateDraftPickRecordsBatch(db, picksToUpload, 50)
}

func GetAllDraftPicksBySeasonID(id string) []structs.DraftPick {
	return repository.FindDraftPicks(id)
}

func GetDraftBootstrap() structs.ProDraftPageResponse {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	var wg sync.WaitGroup
	wg.Add(5)
	var (
		draftees         []structs.DraftablePlayer
		collegePlayers   []structs.CollegePlayer
		warRooms         []structs.ProWarRoom
		draftPicks       []structs.DraftPick
		draftPickFormat  map[uint][]structs.DraftPick
		scoutingProfiles []structs.ScoutingProfile
	)

	// GetProWarRoomByTeamID
	go func() {
		defer wg.Done()
		warRooms = repository.FindProWarRooms()
	}()

	// GetProDrafteesForDraftPage
	go func() {
		defer wg.Done()
		draftees = repository.FindAllDraftablePlayers(repository.PlayerQuery{})
	}()

	go func() {
		defer wg.Done()
		collegePlayers = repository.FindAllCollegePlayers(repository.PlayerQuery{})
	}()

	// GetAllProWarRooms
	go func() {
		defer wg.Done()
		scoutingProfiles = repository.FindScoutingProfiles(repository.ScoutProfileQuery{})
	}()

	// GetAllCurrentSeasonDraftPicksForDraftRoom
	go func() {
		defer wg.Done()
		draftPicks = repository.FindDraftPicks(seasonID)
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	warRoomMap := MakeProWarRoomMapByTeamID(warRooms)
	collegeDraftablePlayers := MakeDraftablePlayerList(collegePlayers)
	cleanedDraftablePlayers := MakeDraftablePlayerListWithGrades(draftees)

	finalDraftablePlayerList := append(cleanedDraftablePlayers, collegeDraftablePlayers...)

	sort.Slice(finalDraftablePlayerList, func(i, j int) bool {
		// Sort by overall desc and then age desc
		if finalDraftablePlayerList[i].Overall == finalDraftablePlayerList[j].Overall {
			return finalDraftablePlayerList[i].Age > finalDraftablePlayerList[j].Age
		}
		return finalDraftablePlayerList[i].Overall > finalDraftablePlayerList[j].Overall
	})

	for _, pick := range draftPicks {
		roundIdx := uint(pick.DraftRound)
		if roundIdx > 0 {
			draftPickFormat[roundIdx] = append(draftPickFormat[roundIdx], pick)
		} else {
			log.Panicln("Invalid round to insert pick!")
		}

	}

	res := structs.ProDraftPageResponse{
		WarRoomMap:       warRoomMap,
		DraftablePlayers: finalDraftablePlayerList,
		DraftPicks:       draftPickFormat,
		ScoutingProfiles: scoutingProfiles,
	}

	return res
}

func CreateScoutingProfile(dto structs.ScoutingProfileDTO) structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := repository.FindScoutingProfile(repository.ScoutProfileQuery{PlayerID: strconv.Itoa(int(dto.PlayerID)), TeamID: strconv.Itoa(int(dto.TeamID))})

	// If Recruit Already Exists
	if scoutProfile.PlayerID > 0 && scoutProfile.TeamID > 0 {
		scoutProfile.ReplaceOnBoard()
		repository.SaveScoutingProfileRecord(scoutProfile, db)
		return scoutProfile
	}

	newScoutingProfile := structs.ScoutingProfile{
		PlayerID:         dto.PlayerID,
		TeamID:           dto.TeamID,
		ShowCount:        0,
		RemovedFromBoard: false,
	}

	repository.CreateScoutingProfileRecord(newScoutingProfile, db)

	return newScoutingProfile
}

func RemovePlayerFromScoutBoard(id string) {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := repository.FindScoutingProfile(repository.ScoutProfileQuery{ID: id})

	scoutProfile.RemoveFromBoard()

	repository.SaveScoutingProfileRecord(scoutProfile, db)
}

func GetScoutingDataByPlayerID(id string) structs.ScoutingDataResponse {
	ts := GetTimestamp()
	lastSeasonID := ts.SeasonID - 1
	lastSeasonIDSTR := strconv.Itoa(int(lastSeasonID))

	collegePlayerRecord := repository.FindCollegePlayer(repository.PlayerQuery{ID: id})

	seasonStats := repository.FindCollegePlayerSeasonStatRecord(id, lastSeasonIDSTR, "2")
	teamID := strconv.Itoa(int(collegePlayerRecord.TeamID))
	collegeStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: lastSeasonIDSTR, TeamID: teamID})
	standing := structs.CollegeStandings{}
	if len(collegeStandings) > 0 {
		standing = collegeStandings[0]
	}

	return structs.ScoutingDataResponse{
		DrafteeSeasonStats: seasonStats,
		TeamStandings:      standing,
	}
}

func RevealScoutingAttribute(dto structs.RevealAttributeDTO) bool {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := repository.FindScoutingProfile(repository.ScoutProfileQuery{ID: strconv.Itoa(int(dto.ScoutProfileID))})

	if scoutProfile.ID == 0 {
		return false
	}

	scoutProfile.RevealAttribute(dto.Attribute)

	warRoom := repository.FindProWarRoomRecord(repository.ScoutProfileQuery{TeamID: strconv.Itoa(int(dto.TeamID))})

	if warRoom.ID == 0 || warRoom.SpentPoints >= warRoom.ScoutingPoints || warRoom.SpentPoints+dto.Points > warRoom.ScoutingPoints {
		return false
	}

	warRoom.AddToSpentPoints(dto.Points)
	repository.SaveScoutingProfileRecord(scoutProfile, db)
	repository.SaveProWarRoomRecord(warRoom, db)
	return true
}

func ExportDraftedPlayers(picks []structs.DraftPick) bool {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	newProPlayerRecords := []structs.ProfessionalPlayer{}
	newContractRecords := []structs.ProContract{}

	collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{})
	draftableCollegePlayerMap := MakeCollegePlayerMap(collegePlayers)
	proPlayers := repository.FindAllProPlayers(repository.PlayerQuery{})
	proPlayerMap := MakeProfessionalPlayerMap(proPlayers)

	for _, pick := range picks {
		if pick.IsVoid {
			continue
		}
		playerId := strconv.Itoa(int(pick.SelectedPlayerID))
		draftee := repository.FindDraftablePlayerRecord(repository.ScoutProfileQuery{PlayerID: playerId})

		if draftee.DraftablePlayerType == 0 {
			// Update college player record
			collegePlayer := draftableCollegePlayerMap[pick.DrafteeID]
			collegePlayer.AssignDraftedTeamData(pick)
			repository.SaveCollegeHockeyPlayerRecord(collegePlayer, db)
		} else if draftee.DraftablePlayerType == 1 {
			draftee.AssignDraftedTeam(pick.DraftNumber, pick.ID, pick.TeamID, pick.Team)
			proPlayer := structs.ProfessionalPlayer{
				Model:          draftee.Model,
				BasePlayer:     draftee.BasePlayer, // Assuming BasePlayer fields are common
				BasePotentials: draftee.BasePotentials,
				DraftedTeamID:  uint8(pick.TeamID),
				DraftedTeam:    pick.Team,
				DraftedRound:   uint8(pick.DraftRound),
				DraftedPick:    uint16(pick.DraftNumber),
				DraftedYearID:  pick.SeasonID,
				CollegeID:      draftee.CollegeID,
				Year:           1,
			}
			proPlayer.AssignDraftedTeam(pick.TeamID, pick.Team, 1)
			playerReference := proPlayerMap[proPlayer.ID]
			if playerReference.ID > 0 {
				continue
			}
			year1Salary := util.GetDrafteeSalary(pick.DraftNumber, 1, pick.DraftRound, true)
			year2Salary := util.GetDrafteeSalary(pick.DraftNumber, 2, pick.DraftRound, true)
			year3Salary := util.GetDrafteeSalary(pick.DraftNumber, 3, pick.DraftRound, true)
			year4Salary := util.GetDrafteeSalary(pick.DraftNumber, 4, pick.DraftRound, true)
			yearsRemaining := 4
			contract := structs.ProContract{
				PlayerID:       proPlayer.ID,
				TeamID:         uint(proPlayer.TeamID),
				OriginalTeamID: uint(proPlayer.TeamID),
				ContractLength: yearsRemaining,
				ContractType:   "Rookie",
				Y1BaseSalary:   year1Salary,
				Y2BaseSalary:   year2Salary,
				Y3BaseSalary:   year3Salary,
				Y4BaseSalary:   year4Salary,
				IsActive:       true,
			}
			newContractRecords = append(newContractRecords, contract)
			newProPlayerRecords = append(newProPlayerRecords, proPlayer)
			repository.SaveDraftablePlayerRecord(draftee, db)
		}
	}

	draftablePlayers := repository.FindAllDraftablePlayers(repository.PlayerQuery{})
	// Move all undrafted players as UDFAs

	for _, draftee := range draftablePlayers {
		if draftee.DraftPickID > 0 {
			continue
		}

		proPlayer := structs.ProfessionalPlayer{
			Model:          draftee.Model,
			BasePlayer:     draftee.BasePlayer, // Assuming BasePlayer fields are common
			BasePotentials: draftee.BasePotentials,
			CollegeID:      draftee.CollegeID,
			Year:           1,
		}
		playerReference := proPlayerMap[proPlayer.ID]
		if playerReference.ID > 0 {
			continue
		}
		newProPlayerRecords = append(newProPlayerRecords, proPlayer)
	}

	repository.CreateProHockeyPlayerRecordsBatch(db, newProPlayerRecords, 250)
	repository.CreateProContractRecordsBatch(db, newContractRecords, 250)

	ts.DraftIsOver()
	repository.SaveTimestamp(ts, db)
	return true
}

func BringUpCollegePlayerToPros(pickID string) bool {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if !ts.IsOffSeason {
		return false
	}
	draftPick := repository.FindDraftPickRecord(pickID)
	if draftPick.ID == 0 || draftPick.DrafteeID == 0 || draftPick.SeasonID >= ts.Season {
		return false
	}

	collegePlayer := repository.FindCollegePlayer(repository.PlayerQuery{ID: strconv.Itoa(int(draftPick.DrafteeID))})

	proPlayer := structs.ProfessionalPlayer{
		Model:          collegePlayer.Model,
		BasePlayer:     collegePlayer.BasePlayer, // Assuming BasePlayer fields are common
		BasePotentials: collegePlayer.BasePotentials,
		DraftedTeamID:  uint8(draftPick.TeamID),
		DraftedTeam:    draftPick.Team,
		DraftedRound:   uint8(draftPick.DraftRound),
		DraftedPick:    uint16(draftPick.DraftNumber),
		DraftedYearID:  draftPick.SeasonID,
		CollegeID:      uint(collegePlayer.TeamID),
		Year:           1,
	}
	proPlayer.AssignDraftedTeam(draftPick.TeamID, draftPick.Team, 1)

	repository.CreateProHockeyPlayerRecordsBatch(db, []structs.ProfessionalPlayer{proPlayer}, 1)

	// Create Historic College Player Record
	historicRecord := structs.HistoricCollegePlayer{
		CollegePlayer: collegePlayer,
	}
	repository.CreateHistoricCollegePlayerRecordsBatch(db, []structs.HistoricCollegePlayer{historicRecord}, 1)

	// Delete College Player Record
	repository.DeleteCollegeHockeyPlayerRecord(db, collegePlayer)

	return true
}
