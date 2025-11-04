package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type CrootGenerator struct {
	nameMap           map[string]map[string][]string
	collegePlayerList []structs.CollegePlayer
	teamMap           map[uint]structs.CollegeTeam
	positionList      []string
	CrootList         []structs.Recruit
	GlobalList        []structs.GlobalPlayer
	FacesList         []structs.FaceData
	attributeBlob     map[string]map[string]map[string]map[string]interface{}
	usCrootLocations  map[string][]structs.CrootLocation
	cnCrootLocations  map[string][]structs.CrootLocation
	svCrootLocations  map[string][]structs.CrootLocation
	ruCrootLocations  map[string][]structs.CrootLocation
	faceDataBlob      map[string][]string
	newID             uint
	count             int
	requiredPlayers   int
	cCount            int
	fCount            int
	dCount            int
	gCount            int
	star5             int
	star4             int
	star3             int
	star2             int
	star1             int
	highestOvr        uint8
	lowestOvr         uint8
	pickedEthnicity   string
	caser             cases.Caser
}

func (pg *CrootGenerator) GenerateRecruits() {
	for pg.count < pg.requiredPlayers {
		player, globalPlayer := pg.generatePlayer()
		pg.CrootList = append(pg.CrootList, player)
		pg.GlobalList = append(pg.GlobalList, globalPlayer)
		pg.updateStatistics(player) // A method to update player counts and statistics
		if player.RelativeType == 5 {
			twinPlayer, twinGlobal := pg.generateTwin(&player)
			pg.updateStatistics(twinPlayer)
			pg.CrootList = append(pg.CrootList, twinPlayer)
			pg.GlobalList = append(pg.GlobalList, twinGlobal)
			pg.count++
		}
		pg.count++
		pg.newID++
	}
}

func (pg *CrootGenerator) generatePlayer() (structs.Recruit, structs.GlobalPlayer) {
	cpLen := len(pg.collegePlayerList) - 1
	relativeType := 0
	relativeID := 0
	coachTeamID := 0
	coachTeamAbbr := ""
	notes := ""
	star := util.GetStarRating(false, false)
	state := ""
	country := pickCountry()
	switch country {
	case util.USA:
		state = util.PickState()
	case util.Canada:
		state = util.PickProvince()
	case util.Sweden:
		state = util.PickSwedishRegion()
	case util.Russia:
		state = util.PickRussianRegion()
	}
	pickedEthnicity := pickLocale(country)
	countryNames := pg.nameMap[pickedEthnicity]
	firstNameList := countryNames["first_names"]
	lastNameList := countryNames["last_names"]
	fName := util.PickFromStringList(firstNameList)
	firstName := pg.caser.String(strings.ToLower(fName))
	lastName := ""
	roof := 150
	relativeRoll := util.GenerateIntFromRange(1, roof)
	relativeIdx := 0
	if relativeRoll == roof {
		relativeType = getRelativeType()
		switch relativeType {
		case 2:
			// Brother of college player
			fmt.Println("BROTHER")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			lastName = cp.LastName
			state = cp.State
			country = cp.Country
			notes = "Brother of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		case 3:
			fmt.Println("COUSIN")
			// Cousin
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 1 {
				lastName = cp.LastName
			} else {
				lName := util.PickFromStringList(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Cousin of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		case 4:
			// Half Brother
			fmt.Println("HALF BROTHER GENERATED")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 3)
			if coinFlip < 3 {
				lastName = cp.LastName
			} else {
				lName := util.PickFromStringList(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Half-Brother of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		case 5:
			// Twin
			relativeType = 5
			relativeID = int(pg.newID)
		default:
			relativeType = 1
		}
	}
	if relativeType == 1 || relativeType == 5 || lastName == "" {
		lName := util.PickFromStringList(lastNameList)
		lastName = pg.caser.String(strings.ToLower(lName))
	}
	if state == "" && country == util.USA {
		state = util.PickState()
	}

	pickedPosition := util.PickPosition()
	crootLocations := pg.usCrootLocations[state]
	switch country {
	case util.Canada:
		crootLocations = pg.cnCrootLocations[state]
	case util.Sweden:
		crootLocations = pg.svCrootLocations[state]
	case util.Russia:
		crootLocations = pg.ruCrootLocations[state]
	}
	player := createRecruit(pickedPosition, "", star, firstName, lastName, pg.attributeBlob, country, state, "", "", crootLocations)
	player.AssignRelativeData(uint(relativeID), uint(relativeType), uint(coachTeamID), coachTeamAbbr, notes)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID:      pg.newID,
		RecruitID:            pg.newID,
		ProfessionalPlayerID: pg.newID,
	}

	skinColor := getSkinColor(country)
	face := getFace(pg.newID, int(player.Weight), skinColor, pg.faceDataBlob)

	pg.FacesList = append(pg.FacesList, face)

	globalPlayer.AssignID(pg.newID)
	player.AssignID(pg.newID)
	return player, globalPlayer
}

func (pg *CrootGenerator) generateTwin(player *structs.Recruit) (structs.Recruit, structs.GlobalPlayer) {
	fmt.Println("TWIN!!")
	// Generate Twin Record
	firstTwinRelativeID := pg.newID
	pg.newID++
	// Twin being generated is secondTwin
	secondTwinRelativeID := pg.newID
	country := pickCountry()
	pickedEthnicity := pickLocale(country)
	countryNames := pg.nameMap[pickedEthnicity]
	firstNameList := countryNames["first_names"]
	twinName := util.PickFromStringList(firstNameList)
	twinN := pg.caser.String(strings.ToLower(twinName))
	twinPosition := util.PickFromStringList(pg.positionList)
	coinFlip := util.GenerateIntFromRange(1, 2)
	stars := util.GetStarRating(false, false)
	if coinFlip == 2 {
		twinPosition = player.Position
		stars = int(player.Stars)
	}
	crootLocations := pg.usCrootLocations[player.State]
	if country == "Canada" {
		crootLocations = pg.cnCrootLocations[player.State]
	}
	twinNotes := "Twin Brother of " + strconv.Itoa(int(player.Stars)) + " Star Recruit " + player.Position + " " + player.FirstName + " " + player.LastName
	twinPlayer := createRecruit(twinPosition, "", stars, twinN, player.LastName, pg.attributeBlob, player.Country, player.State, "", "", crootLocations)
	twinPlayer.AssignRelativeData(uint(firstTwinRelativeID), 4, 0, "", twinNotes)
	twinPlayer.AssignTwinData(player.LastName, player.City, player.State, player.HighSchool)
	notes := "Twin Brother of " + strconv.Itoa(int(twinPlayer.Stars)) + " Star Recruit " + twinPlayer.Position + " " + twinPlayer.FirstName + " " + twinPlayer.LastName
	player.AssignRelativeData(uint(secondTwinRelativeID), 4, 0, "", notes)
	globalTwinPlayer := structs.GlobalPlayer{
		CollegePlayerID:      secondTwinRelativeID,
		RecruitID:            secondTwinRelativeID,
		ProfessionalPlayerID: secondTwinRelativeID,
	}
	globalTwinPlayer.AssignID(secondTwinRelativeID)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID:      firstTwinRelativeID,
		RecruitID:            firstTwinRelativeID,
		ProfessionalPlayerID: firstTwinRelativeID,
	}
	globalPlayer.AssignID(uint(firstTwinRelativeID))
	skinColor := getSkinColor(player.Country)

	face := getFace(secondTwinRelativeID, int(twinPlayer.Weight), skinColor, pg.faceDataBlob)

	pg.FacesList = append(pg.FacesList, face)
	return twinPlayer, globalTwinPlayer
}

func (pg *CrootGenerator) updateStatistics(player structs.Recruit) {
	switch player.Stars {
	case 5:
		pg.star5++
	case 4:
		pg.star4++
	case 3:
		pg.star3++
	case 2:
		pg.star2++
	default:
		pg.star1++
	}
	switch player.Position {
	case "C":
		pg.cCount++
	case "F":
		pg.fCount++
	case "D":
		pg.dCount++
	case "G":
		pg.gCount++
	}

	if player.Overall > pg.highestOvr {
		pg.highestOvr = player.Overall
	}
	if player.Overall < pg.lowestOvr {
		pg.lowestOvr = player.Overall
	}
}

func (pg *CrootGenerator) OutputRecruitStats() {
	// Croot Stats
	fmt.Println("Total Recruit Count: ", pg.count)
	fmt.Println("Total 5 Star  Count: ", pg.star5)
	fmt.Println("Total 4 Star  Count: ", pg.star4)
	fmt.Println("Total 3 Star  Count: ", pg.star3)
	fmt.Println("Total 2 Star  Count: ", pg.star2)
	fmt.Println("Total 1 Star  Count: ", pg.star1)
	fmt.Println("Total C  Count: ", pg.cCount)
	fmt.Println("Total F  Count: ", pg.fCount)
	fmt.Println("Total D  Count: ", pg.dCount)
	fmt.Println("Total G  Count: ", pg.gCount)

	fmt.Println("Highest Recruit Ovr: ", pg.highestOvr)
	fmt.Println("Lowest  Recruit Ovr: ", pg.lowestOvr)
}

func GenerateCroots() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer
	ts := GetTimestamp()

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer
	// fNameMap, lNameMap := getNameMaps()
	generator := CrootGenerator{
		nameMap:           getInternationalNameMap(),
		collegePlayerList: GetAllCollegePlayers(),
		teamMap:           GetCollegeTeamMap(),
		usCrootLocations:  getCrootLocations("HS"),
		cnCrootLocations:  getCrootLocations("CanadianHS"),
		svCrootLocations:  getCrootLocations("SwedenHS"),
		ruCrootLocations:  getCrootLocations("RussianHS"),
		attributeBlob:     getAttributeBlob(),
		positionList:      util.GetPositionList(),
		newID:             lastPlayerRecord.ID + 1,
		requiredPlayers:   util.GenerateIntFromRange(462, 660),
		faceDataBlob:      getFaceDataBlob(),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		highestOvr:        0,
		lowestOvr:         200,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.GlobalPlayer{},
		FacesList:         []structs.FaceData{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}

	// The plan is to ensure that every recruit is signed
	generator.GenerateRecruits()
	// Croot Stats
	generator.OutputRecruitStats()

	repository.CreateHockeyRecruitRecordsBatch(db, generator.CrootList, 500)
	repository.CreateGlobalPlayerRecordsBatch(db, generator.GlobalList, 500)
	repository.CreateFaceRecordsBatch(db, generator.FacesList, 500)
	ts.ToggleGeneratedCroots()
	repository.SaveTimestamp(ts, db)
	// AssignAllRecruitRanks()
}

func GenerateCustomCroots() {
	db := dbprovider.GetInstance().GetDB()
	lastPlayerRecord := repository.FindLatestGlobalPlayerRecord()
	latestID := lastPlayerRecord.ID + 1
	cpList := []structs.Recruit{}
	globalList := []structs.GlobalPlayer{}
	facesList := []structs.FaceData{}

	filePath := filepath.Join(os.Getenv("ROOT"), "data", "gen", "custom_croots.csv")
	playersCSV := util.ReadCSV(filePath)
	pg := CrootGenerator{
		attributeBlob: getAttributeBlob(),
		faceDataBlob:  getFaceDataBlob(),
	}
	for idx, row := range playersCSV {
		if idx == 0 {
			continue
		}

		star := util.GetStarRating(true, false)
		fName := row[0]
		lName := row[1]
		position := row[2]
		archetype := row[3]
		country := row[4]
		state := row[5]
		city := row[6]
		hs := row[7]
		height := util.ConvertStringToInt(row[8])
		weight := util.ConvertStringToInt(row[9])
		crootFor := row[10]
		player := createRecruit(position, archetype, star, fName, lName, pg.attributeBlob, country, state, city, hs, []structs.CrootLocation{})
		player.AssignHeightAndWeight(height, weight)
		player.ToggleCustomCroot(crootFor)

		skinColor := getSkinColor(country)
		face := getFace(latestID, int(player.Weight), skinColor, pg.faceDataBlob)

		facesList = append(facesList, face)

		globalPlayer := structs.GlobalPlayer{
			Model: gorm.Model{
				ID: latestID,
			},
			RecruitID:            latestID,
			CollegePlayerID:      latestID,
			ProfessionalPlayerID: latestID,
		}
		globalPlayer.AssignID(latestID)
		player.AssignID(latestID)
		latestID++

		globalList = append(globalList, globalPlayer)
		cpList = append(cpList, player)

	}
	repository.CreateHockeyRecruitRecordsBatch(db, cpList, 100)
	repository.CreateGlobalPlayerRecordsBatch(db, globalList, 100)
	repository.CreateFaceRecordsBatch(db, facesList, 100)

	AssignAllRecruitRanks()
}

func GenerateWalkonCroots() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer
	ts := GetTimestamp()

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	profiles := repository.FindTeamRecruitingProfiles(false)
	collegePlayers := GetAllCollegePlayers()
	collegeRosterMap := MakeCollegePlayerMapByTeamID(collegePlayers)
	recruits := repository.FindAllRecruits(false, true, true, false, false, "")
	recruitMap := MakeCollegeRecruitMapByTeamID(recruits)
	generator := CrootGenerator{
		nameMap:          getInternationalNameMap(),
		usCrootLocations: getCrootLocations("HS"),
		cnCrootLocations: getCrootLocations("CanadianHS"),
		attributeBlob:    getAttributeBlob(),
		newID:            lastPlayerRecord.ID + 1,
		faceDataBlob:     getFaceDataBlob(),
		caser:            cases.Title(language.English),
	}
	recruitProfileBatchList := []structs.RecruitPlayerProfile{}
	faces := []structs.FaceData{}
	globalPlayerList := []structs.GlobalPlayer{}
	recruitBatchList := []structs.Recruit{}

	rosterLimit := 34
	positionList := []string{Center, Forward, Forward, Defender, Defender, Goalie}
	newID := generator.newID

	for _, team := range profiles {
		teamID := team.TeamID
		collegeRoster := collegeRosterMap[teamID]
		currentRecruits := recruitMap[teamID]
		neededPlayers := rosterLimit - (len(collegeRoster) + len(currentRecruits))
		if neededPlayers < 1 {
			continue
		}

		rand.Shuffle(len(positionList), func(i, j int) {
			positionList[i], positionList[j] = positionList[j], positionList[i]
		})

		count := 0
		for _, pos := range positionList {
			if count >= neededPlayers {
				break
			}
			arch := util.GetArchetype(pos)
			country := team.Country
			pickedEthnicity := pickLocale(country)

			state := pickWalkonState(team.State)
			crootLocations := generator.usCrootLocations[state]
			if team.Country == util.Canada {
				crootLocations = generator.cnCrootLocations[state]
			}

			city, hs := getCityAndHighSchool(crootLocations)

			countryNames := generator.nameMap[pickedEthnicity]
			firstNameList := countryNames["first_names"]
			lastNameList := countryNames["last_names"]
			stars := util.GetStarRating(false, false)
			lastName := util.PickFromStringList(lastNameList)
			casedLastName := generator.caser.String(strings.ToLower(lastName))
			recruit := createRecruit(pos, arch, stars, util.PickFromStringList(firstNameList), casedLastName, generator.attributeBlob, country, state, city, hs, crootLocations)

			recruit.AssignWalkon(team.Team, int(team.ID), newID)

			recruitPlayerRecord := structs.RecruitPlayerProfile{
				ProfileID:   team.ID,
				RecruitID:   newID,
				IsSigned:    true,
				IsLocked:    true,
				SeasonID:    ts.SeasonID,
				TotalPoints: 1,
			}

			playerRecord := structs.GlobalPlayer{
				RecruitID:            newID,
				CollegePlayerID:      newID,
				ProfessionalPlayerID: newID,
			}
			playerRecord.AssignID(newID)
			count++
			skinColor := getSkinColor(pickedEthnicity)

			face := getFace(newID, int(recruit.Weight), skinColor, generator.faceDataBlob)
			faces = append(faces, face)
			globalPlayerList = append(globalPlayerList, playerRecord)
			recruitBatchList = append(recruitBatchList, recruit)
			recruitProfileBatchList = append(recruitProfileBatchList, recruitPlayerRecord)
			newID++
			team.IncreaseCommitCount()

		}
		repository.SaveTeamProfileRecord(db, team)
	}
	repository.CreateHockeyRecruitRecordsBatch(db, recruitBatchList, 500)
	repository.CreateGlobalPlayerRecordsBatch(db, globalPlayerList, 500)
	repository.CreateFaceRecordsBatch(db, faces, 500)
	repository.CreateRecruitProfileRecordsBatch(db, recruitProfileBatchList, 500)
	ts.ToggleGeneratedCroots()
	repository.SaveTimestamp(ts, db)
	// AssignAllRecruitRanks()
}

func GenerateInitialRosters() {
	db := dbprovider.GetInstance().GetDB()
	lastPlayerRecord := repository.FindLatestGlobalPlayerRecord()
	latestID := lastPlayerRecord.ID + 1
	cpList := []structs.CollegePlayer{}
	globalList := []structs.GlobalPlayer{}
	// filePath := filepath.Join(os.Getenv("ROOT"), "data", "gen", "init_roster.csv")
	// playersCSV := util.ReadCSV(filePath)
	teams := GetAllCollegeTeams()
	generator := CrootGenerator{
		nameMap:           getInternationalNameMap(),
		collegePlayerList: GetAllCollegePlayers(),
		teamMap:           GetCollegeTeamMap(),
		usCrootLocations:  getCrootLocations("HS"),
		cnCrootLocations:  getCrootLocations("CanadianHS"),
		svCrootLocations:  getCrootLocations("SwedenHS"),
		ruCrootLocations:  getCrootLocations("RussianHS"),
		attributeBlob:     getAttributeBlob(),
		positionList:      util.GetPositionList(),
		newID:             1,
		requiredPlayers:   util.GenerateIntFromRange(6400, 6601),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		highestOvr:        0,
		lowestOvr:         200,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.GlobalPlayer{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}
	for _, team := range teams {
		if team.ID < 73 || team.LeagueID == 2 {
			continue
		}
		teamID := team.ID
		queue := getCollegeGenList(teamID)

		for _, dto := range queue {
			year := dto.Year
			pos := dto.Pos
			age := 18 + year - 1
			p, _ := generator.createInitialPlayer(pos, false, age)
			cp := structs.CollegePlayer{
				BasePlayer:     p.BasePlayer,
				BasePotentials: p.BasePotentials,
				BaseInjuryData: p.BaseInjuryData,
				Year:           1,
			}
			cp.AssignTeam(team.ID, team.Abbreviation, 1)
			cp.AssignID(latestID)

			for j := cp.Age; j < uint8(age); j++ {
				cp = ProgressCollegePlayer(cp, "1", []structs.CollegePlayerGameStats{})
			}
			cpList = append(cpList, cp)

			globalPlayer := structs.GlobalPlayer{
				Model: gorm.Model{
					ID: latestID,
				},
				RecruitID:            latestID,
				CollegePlayerID:      latestID,
				ProfessionalPlayerID: latestID,
			}
			latestID++

			globalList = append(globalList, globalPlayer)

		}
	}
	repository.CreateCollegeHockeyPlayerRecordsBatch(db, cpList, 100)
	repository.CreateGlobalPlayerRecordsBatch(db, globalList, 100)

}

func GenerateInitialCHLRosters() {
	db := dbprovider.GetInstance().GetDB()
	lastPlayerRecord := repository.FindLatestGlobalPlayerRecord()
	latestID := lastPlayerRecord.ID + 1
	cpList := []structs.CollegePlayer{}
	globalList := []structs.GlobalPlayer{}
	// filePath := filepath.Join(os.Getenv("ROOT"), "data", "gen", "init_roster.csv")
	// playersCSV := util.ReadCSV(filePath)
	teams := GetAllCanadianCHLTeams()
	generator := CrootGenerator{
		nameMap:           getInternationalNameMap(),
		collegePlayerList: GetAllCollegePlayers(),
		teamMap:           GetCollegeTeamMap(),
		usCrootLocations:  getCrootLocations("HS"),
		cnCrootLocations:  getCrootLocations("CanadianHS"),
		svCrootLocations:  getCrootLocations("SwedenHS"),
		ruCrootLocations:  getCrootLocations("RussianHS"),
		attributeBlob:     getAttributeBlob(),
		positionList:      util.GetPositionList(),
		newID:             1,
		requiredPlayers:   util.GenerateIntFromRange(6400, 6601),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		highestOvr:        0,
		lowestOvr:         200,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.GlobalPlayer{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}
	for _, team := range teams {
		if team.LeagueID == 1 {
			continue
		}
		teamID := team.ID
		queue := getCollegeGenList(teamID)

		for _, dto := range queue {
			year := dto.Year
			pos := dto.Pos
			age := 16 + year - 1
			p, _ := generator.createInitialPlayer(pos, true, age)
			cp := structs.CollegePlayer{
				BasePlayer:     p.BasePlayer,
				BasePotentials: p.BasePotentials,
				BaseInjuryData: p.BaseInjuryData,
				Year:           1,
			}
			cp.AssignTeam(team.ID, team.Abbreviation, 1)
			cp.AssignID(latestID)

			for j := uint8(18); j < uint8(age); j++ {
				cp = ProgressCollegePlayer(cp, "1", []structs.CollegePlayerGameStats{})
			}
			cp.ResetCHLCollegeYear()
			cpList = append(cpList, cp)

			globalPlayer := structs.GlobalPlayer{
				Model: gorm.Model{
					ID: latestID,
				},
				RecruitID:            latestID,
				CollegePlayerID:      latestID,
				ProfessionalPlayerID: latestID,
			}
			latestID++

			globalList = append(globalList, globalPlayer)

		}
	}
	repository.CreateCollegeHockeyPlayerRecordsBatch(db, cpList, 100)
	repository.CreateGlobalPlayerRecordsBatch(db, globalList, 100)

}

func GenerateInitialProPool() {
	db := dbprovider.GetInstance().GetDB()

	lastPlayerRecord := repository.FindLatestGlobalPlayerRecord()
	latestID := lastPlayerRecord.ID + 1
	proTeams := repository.FindAllProTeams(repository.TeamClauses{LeagueID: "1"})
	ageList := []int{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7, 8, 8, 9}
	collegeTeams := GetAllCollegeTeams()
	proList := []structs.ProfessionalPlayer{}
	globalList := []structs.GlobalPlayer{}
	positionNeeds := []string{"C", "F", "F", "D", "D", "G"}
	generator := CrootGenerator{
		nameMap:           getInternationalNameMap(),
		collegePlayerList: GetAllCollegePlayers(),
		teamMap:           GetCollegeTeamMap(),
		usCrootLocations:  getCrootLocations("HS"),
		cnCrootLocations:  getCrootLocations("CanadianHS"),
		svCrootLocations:  getCrootLocations("SwedenHS"),
		ruCrootLocations:  getCrootLocations("RussianHS"),
		attributeBlob:     getAttributeBlob(),
		positionList:      util.GetPositionList(),
		newID:             lastPlayerRecord.ID + 1,
		requiredPlayers:   util.GenerateIntFromRange(6400, 6601),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		highestOvr:        0,
		lowestOvr:         200,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.GlobalPlayer{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}
	positionNeedIdx := 0
	for idx := range proTeams {
		if idx > 24 {
			break
		}
		for _, age := range ageList {
			pos := positionNeeds[positionNeedIdx]
			positionNeedIdx++
			if positionNeedIdx > len(positionNeeds)-1 {
				positionNeedIdx = 0
			}
			p, _ := generator.createInitialPlayer(pos, false, age)
			cp := structs.CollegePlayer{
				BasePlayer:     p.BasePlayer,
				BasePotentials: p.BasePotentials,
				BaseInjuryData: p.BaseInjuryData,
				Year:           1,
			}
			teamIdx := util.GenerateIntFromRange(0, len(collegeTeams)-1)
			team := collegeTeams[teamIdx]
			cp.AssignTeam(team.ID, team.Abbreviation, 1)

			for range 4 {
				cp = ProgressCollegePlayer(cp, "1", []structs.CollegePlayerGameStats{})
			}
			// Age of player is now 22
			pro := structs.ProfessionalPlayer{
				BasePlayer:     cp.BasePlayer,
				BasePotentials: cp.BasePotentials,
				BaseInjuryData: cp.BaseInjuryData,
				Year:           1,
				CollegeID:      uint(cp.TeamID),
			}
			pro.AssignTeam(0, "", 0)
			pro.AssignID(latestID)

			// pro.AssignTeam(proTeam.ID, proTeam.Abbreviation)

			for range age {
				pro = ProgressProPlayer(pro, "1", []structs.ProfessionalPlayerGameStats{})
			}

			proList = append(proList, pro)

			globalPlayer := structs.GlobalPlayer{
				Model: gorm.Model{
					ID: latestID,
				},
				RecruitID:            latestID,
				CollegePlayerID:      latestID,
				ProfessionalPlayerID: latestID,
			}
			latestID++

			globalList = append(globalList, globalPlayer)
		}

	}

	repository.CreateProHockeyPlayerRecordsBatch(db, proList, 500)
	repository.CreateGlobalPlayerRecordsBatch(db, globalList, 500)
}

func createRecruit(position, arch string, stars int, firstName, lastName string, blob map[string]map[string]map[string]map[string]interface{}, country, state, cit, hs string, hsBlob []structs.CrootLocation) structs.Recruit {
	age := 18
	city, highSchool := cit, hs
	if state != "" && len(hsBlob) > 0 {
		city, highSchool = getCityAndHighSchool(hsBlob)
	}
	archetype := util.GetArchetype(position)
	if len(arch) > 0 {
		archetype = arch
	}
	height := getAttributeValue(position, archetype, stars, "Height", blob)
	weight := getAttributeValue(position, archetype, stars, "Weight", blob)
	agility := getAttributeValue(position, archetype, stars, util.Agility, blob)
	agility = getNationalityValueModifier(country, util.Agility, agility)
	faceoffs := getAttributeValue(position, archetype, stars, "Faceoffs", blob)
	longShotAccuracy := getAttributeValue(position, archetype, stars, "LongShotAccuracy", blob)
	longShotPower := getAttributeValue(position, archetype, stars, util.LongShotPower, blob)
	longShotPower = getNationalityValueModifier(country, util.LongShotPower, longShotPower)
	closeShotAccuracy := getAttributeValue(position, archetype, stars, "CloseShotAccuracy", blob)
	closeShotPower := getAttributeValue(position, archetype, stars, "CloseShotPower", blob)
	closeShotPower = getNationalityValueModifier(country, util.CloseShotPower, closeShotPower)
	oneTimer := getAttributeValue(position, archetype, stars, "OneTimer", blob)
	passing := getAttributeValue(position, archetype, stars, "Passing", blob)
	puckHandling := getAttributeValue(position, archetype, stars, util.PuckHandling, blob)
	puckHandling = getNationalityValueModifier(country, util.PuckHandling, puckHandling)
	strength := getAttributeValue(position, archetype, stars, util.Strength, blob)
	strength = getNationalityValueModifier(country, util.Strength, strength)
	bodychecking := getAttributeValue(position, archetype, stars, "BodyChecking", blob)
	bodychecking = getNationalityValueModifier(country, util.BodyChecking, bodychecking)
	stickChecking := getAttributeValue(position, archetype, stars, "StickChecking", blob)
	stickChecking = getNationalityValueModifier(country, util.StickChecking, stickChecking)
	shotBlocking := getAttributeValue(position, archetype, stars, "ShotBlocking", blob)
	goalkeeping := getAttributeValue(position, archetype, stars, "Goalkeeping", blob)
	goalieVision := getAttributeValue(position, archetype, stars, "GoalieVision", blob)
	goalieReboundControl := getAttributeValue(position, archetype, stars, "GoalieReboundControl", blob)
	injury := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	stamina := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	injuryDeviation := util.GenerateIntFromRange(1, 20)
	disciplineDeviation := util.GenerateIntFromRange(1, 20)
	discipline := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	aggression := int(util.GenerateNormalizedIntFromMeanStdev(50, 15))
	aggression = getNationalityValueModifier(country, "Aggression", int(aggression))
	clutch := getClutchValue()
	if archetype == util.Enforcer {
		aggression += util.GenerateIntFromRange(10, 30)
	}
	if aggression > 100 {
		aggression = 100
	}
	personality := util.GetPersonality()

	program := util.GenerateNormalizedIntFromRange(1, 9)
	profDevelopment := util.GenerateNormalizedIntFromRange(1, 9)
	traditions := util.GenerateNormalizedIntFromRange(1, 9)
	facilities := util.GenerateNormalizedIntFromRange(1, 9)
	atmosphere := util.GenerateNormalizedIntFromRange(1, 9)
	academics := util.GenerateNormalizedIntFromRange(1, 9)
	conferencePrestige := util.GenerateNormalizedIntFromRange(1, 9)
	coachPref := util.GenerateNormalizedIntFromRange(1, 9)
	seasonMomentumPref := util.GenerateNormalizedIntFromRange(1, 9)
	playtime := util.GenerateNormalizedIntFromRange(1, 9)
	competitiveness := util.GenerateNormalizedIntFromRange(1, 9)

	agilityPotential := util.GeneratePotential(position, archetype, Agility)
	agilityPotential = uint8(getNationalityPotentialModifier(country, util.Agility, int(agilityPotential)))
	faceoffsPotential := util.GeneratePotential(position, archetype, Faceoffs)
	closeShotAccuracyPotential := util.GeneratePotential(position, archetype, CloseShotAccuracy)
	closeShotPowerPotential := util.GeneratePotential(position, archetype, CloseShotPower)
	closeShotPowerPotential = uint8(getNationalityPotentialModifier(country, util.CloseShotPower, int(closeShotPowerPotential)))
	longShotAccuracyPotential := util.GeneratePotential(position, archetype, LongShotAccuracy)
	longShotPowerPotential := util.GeneratePotential(position, archetype, LongShotPower)
	longShotPowerPotential = uint8(getNationalityPotentialModifier(country, util.LongShotPower, int(longShotPowerPotential)))
	passingPotential := util.GeneratePotential(position, archetype, Passing)
	puckHandlingPotential := util.GeneratePotential(position, archetype, PuckHandling)
	puckHandlingPotential = uint8(getNationalityPotentialModifier(country, util.PuckHandling, int(puckHandlingPotential)))
	strengthPotential := util.GeneratePotential(position, archetype, Strength)
	strengthPotential = uint8(getNationalityPotentialModifier(country, util.Strength, int(strengthPotential)))
	bodyCheckingPotential := util.GeneratePotential(position, archetype, BodyChecking)
	bodyCheckingPotential = uint8(getNationalityPotentialModifier(country, util.BodyChecking, int(bodyCheckingPotential)))
	stickCheckingPotential := util.GeneratePotential(position, archetype, StickChecking)
	stickCheckingPotential = uint8(getNationalityPotentialModifier(country, util.StickChecking, int(stickCheckingPotential)))
	shotBlockingPotential := util.GeneratePotential(position, archetype, ShotBlocking)
	goalkeepingPotential := util.GeneratePotential(position, archetype, Goalkeeping)
	goalieVisionPotential := util.GeneratePotential(position, archetype, GoalieVision)
	goalieReboundPotential := util.GeneratePotential(position, archetype, GoalieRebound)

	potentials := structs.BasePotentials{
		AgilityPotential:           agilityPotential,
		FaceoffsPotential:          faceoffsPotential,
		CloseShotAccuracyPotential: closeShotAccuracyPotential,
		CloseShotPowerPotential:    closeShotPowerPotential,
		LongShotAccuracyPotential:  longShotAccuracyPotential,
		LongShotPowerPotential:     longShotPowerPotential,
		PassingPotential:           passingPotential,
		PuckHandlingPotential:      puckHandlingPotential,
		StrengthPotential:          strengthPotential,
		BodyCheckingPotential:      bodyCheckingPotential,
		StickCheckingPotential:     stickCheckingPotential,
		ShotBlockingPotential:      shotBlockingPotential,
		GoalkeepingPotential:       goalkeepingPotential,
		GoalieVisionPotential:      goalieVisionPotential,
		GoalieReboundPotential:     goalieReboundPotential,
	}

	injuryData := structs.BaseInjuryData{
		Regression: uint8(util.GenerateNormalizedIntFromRange(1, 3)),
		DecayRate:  float32(util.GenerateFloatFromRange(0.15, 0.65)),
	}

	starVal := stars
	if starVal == 6 {
		starVal = 5
	}

	basePlayer := structs.BasePlayer{
		FirstName:            firstName,
		LastName:             lastName,
		Position:             position,
		Archetype:            archetype,
		Age:                  uint8(age),
		Stars:                uint8(starVal),
		Height:               uint8(height),
		Weight:               uint16(weight),
		Stamina:              uint8(stamina),
		InjuryRating:         uint8(injury),
		Agility:              uint8(agility),
		Faceoffs:             uint8(faceoffs),
		LongShotAccuracy:     uint8(longShotAccuracy),
		LongShotPower:        uint8(longShotPower),
		CloseShotAccuracy:    uint8(closeShotAccuracy),
		CloseShotPower:       uint8(closeShotPower),
		OneTimer:             uint8(oneTimer),
		Passing:              uint8(passing),
		PuckHandling:         uint8(puckHandling),
		Strength:             uint8(strength),
		BodyChecking:         uint8(bodychecking),
		StickChecking:        uint8(stickChecking),
		ShotBlocking:         uint8(shotBlocking),
		Goalkeeping:          uint8(goalkeeping),
		GoalieVision:         uint8(goalieVision),
		GoalieReboundControl: uint8(goalieReboundControl),
		Personality:          personality,
		Discipline:           uint8(discipline),
		City:                 city,
		HighSchool:           highSchool,
		State:                state,
		Country:              country,
		PlayerPreferences: structs.PlayerPreferences{
			ProgramPref:        uint8(program),
			ProfDevPref:        uint8(profDevelopment),
			TraditionsPref:     uint8(traditions),
			FacilitiesPref:     uint8(facilities),
			AtmospherePref:     uint8(atmosphere),
			AcademicsPref:      uint8(academics),
			ConferencePref:     uint8(conferencePrestige),
			CoachPref:          uint8(coachPref),
			SeasonMomentumPref: uint8(seasonMomentumPref),
		},

		PlaytimePreference:  uint8(playtime),
		Competitiveness:     uint8(competitiveness),
		Clutch:              int8(clutch),
		Aggression:          uint8(aggression),
		InjuryDeviation:     uint8(injuryDeviation),
		DisciplineDeviation: uint8(disciplineDeviation),
		PrimeAge:            uint8(util.GetPrimeAge(position, archetype)),
		PlayerMorale:        100,
		HasProgressed:       false,
	}

	basePlayer.GetOverall()

	return structs.Recruit{
		BasePlayer:     basePlayer,
		BasePotentials: potentials,
		IsSigned:       false,
		BaseInjuryData: injuryData,
	}
}

func (pg *CrootGenerator) createInitialPlayer(position string, isCHL bool, age int) (structs.Recruit, structs.GlobalPlayer) {
	cpLen := len(pg.collegePlayerList) - 1
	relativeType := 0
	relativeID := 0
	coachTeamID := 0
	coachTeamAbbr := ""
	notes := ""
	star := util.GetStarRating(false, isCHL)
	state := ""
	country := pickCountry()
	if isCHL {
		country = util.Canada
	}
	switch country {
	case util.USA:
		state = util.PickState()
	case util.Canada:
		state = util.PickProvince()
	case util.Sweden:
		state = util.PickSwedishRegion()
	case util.Russia:
		state = util.PickRussianRegion()
	}
	pickedEthnicity := pickLocale(country)
	countryNames := pg.nameMap[pickedEthnicity]
	firstNameList := countryNames["first_names"]
	lastNameList := countryNames["last_names"]
	fName := util.PickFromStringList(firstNameList)
	firstName := pg.caser.String(strings.ToLower(fName))
	lastName := ""
	roof := 100
	relativeRoll := util.GenerateIntFromRange(1, roof)
	relativeIdx := 0
	if relativeRoll == roof && cpLen > 0 && !isCHL {
		relativeType = getRelativeType()
		if relativeType == 2 {
			// Brother of college player
			fmt.Println("BROTHER")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			lastName = cp.LastName
			state = cp.State
			country = cp.Country
			notes = "Brother of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 3 && cpLen > 0 {
			fmt.Println("COUSIN")
			// Cousin
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 1 {
				lastName = cp.LastName
			} else {
				lName := util.PickFromStringList(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Cousin of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 4 && cpLen > 0 {
			// Half Brother
			fmt.Println("HALF BROTHER GENERATED")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 3)
			if coinFlip < 3 {
				lastName = cp.LastName
			} else {
				lName := util.PickFromStringList(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Half-Brother of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 5 && cpLen > 0 {
			// Twin
			relativeType = 5
			relativeID = int(pg.newID)
		} else {
			relativeType = 1
		}
	}
	if relativeType == 1 || relativeType == 5 || lastName == "" {
		lName := util.PickFromStringList(lastNameList)
		if isCHL {
			lastName = lName
		} else {
			lastName = pg.caser.String(strings.ToLower(lName))
		}
	}
	if state == "" && country == util.USA {
		state = util.PickState()
	}

	crootLocations := pg.usCrootLocations[state]
	switch country {
	case util.Canada:
		crootLocations = pg.cnCrootLocations[state]
	case util.Sweden:
		crootLocations = pg.svCrootLocations[state]
	case util.Russia:
		crootLocations = pg.ruCrootLocations[state]
	}

	player := createRecruit(position, "", star, firstName, lastName, pg.attributeBlob, country, state, "", "", crootLocations)
	player.AssignAge(age)
	player.AssignRelativeData(uint(relativeID), uint(relativeType), uint(coachTeamID), coachTeamAbbr, notes)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID:      pg.newID,
		RecruitID:            pg.newID,
		ProfessionalPlayerID: pg.newID,
	}

	globalPlayer.AssignID(pg.newID)
	return player, globalPlayer
}

func getCityAndHighSchool(schools []structs.CrootLocation) (string, string) {
	if len(schools) == 0 {
		fmt.Println("NO SCHOOLS?!")
		return "", ""
	}
	randInt := util.GenerateIntFromRange(0, len(schools)-1)

	return schools[randInt].City, schools[randInt].HighSchool
}

func getRelativeType() int {
	roll := util.GenerateIntFromRange(1, 1000)
	// Brother of existing player
	if roll < 600 {
		return 2
	}
	// Cousin of existing player
	if roll < 800 {
		return 3
	}
	// Half brother of existing player
	if roll < 850 {
		return 4
	}
	// Twin
	if roll < 900 {
		return 5
	}
	// Best friend of another recruit
	if roll < 925 {
		return 8
	}
	// Best friend of a college player
	if roll < 950 {
		return 8
	}
	return 8
}

func pickLocale(country string) string {
	countryMap := map[string][]string{
		"USA":                {"en_US"},
		"England":            {"en_GB", "en_US"},
		"Scotland":           {"en_GB", "en_IE"},
		"Austria":            {"de_AT"},
		"Canada":             {"fr_CA", "en_CA"},
		"Ireland":            {"en_IE"},
		"Wales":              {"en_GB", "en_IE"},
		"Spain":              {"es_ES"},
		"Malta":              {"es_ES"},
		"Italy":              {"it_IT"},
		"Portugal":           {"pt_PT"},
		"France":             {"fr_FR", "fr_CA"},
		"Switzerland":        {"fr_FR", "de_DE"},
		"Andorra":            {"fr_FR", "es_ES"},
		"Germany":            {"de_AT", "de_CH", "de_DE"},
		"Belgium":            {"nl_BE", "nl_NL", "fr_FR"},
		"Netherlands":        {"nl_BE", "nl_NL", "de_DE"},
		"Lithuania":          {"lt_LT"},
		"Latvia":             {"lv_LV", "lt_LT"},
		"Poland":             {"pl_PL"},
		"Finland":            {"sv_SE", "fi_FI"},
		"Denmark":            {"dk_DK", "no_NO"},
		"Sweden":             {"sv_SE", "no_NO"},
		"Iceland":            {"sv_SE", "no_NO"},
		"Norway":             {"no_NO"},
		"Bulgaria":           {"bg_BG", "ro_RO"},
		"Serbia":             {"bs_BA", "sl_SI", "ro_RO", "bg_BG"},
		"Croatia":            {"hu_HU", "sl_SI", "hr_HR"},
		"Hungary":            {"sl_SI", "hu_HU"},
		"Bosnia":             {"bs_BA", "ro_RO", "sl_SI"},
		"Czech Republic":     {"cs_CZ", "bg_BG"},
		"Slovakia":           {"cs_CZ"},
		"Estonia":            {"et_EE", "lt_LT"},
		"Kosovo":             {"sl_SI", "ro_RO"},
		"Montenegro":         {"sl_SI", "ro_RO"},
		"Romania":            {"sl_SI", "ru_RU", "ro_RO", "bg_BG"},
		"Moldova":            {"uk_UA", "ru_RU", "ro_RO"},
		"Slovenia":           {"sl_SI", "ro_RO", "bg_BG"},
		"Cyprus":             {"el_GR", "tr_TR"},
		"Turkey":             {"tr_TR"},
		"Greece":             {"el_GR", "tr_TR"},
		"Albania":            {"el_GR"},
		"North Macedonia":    {"el_GR"},
		"Mexico":             {"es_MX"},
		"Argentina":          {"es_MX"},
		"Brazil":             {"es_MX", "pt_BR"},
		"China":              {"zh_CN"},
		"HK":                 {"zh_CN"},
		"Japan":              {"ja_JP"},
		"South Korea":        {"ko_KR"},
		"Taiwan":             {"zh_TW"},
		"Philippines":        {"en_PH", "es_ES"},
		"Indonesia":          {"ms_MY", "zh_CN"},
		"Malaysia":           {"ms_MY", "vi_VN", "th_TH", "zh_CN"},
		"Singapore":          {"zh_CN", "th_TH"},
		"Laos":               {"zh_CN", "vi_VN"},
		"Myanmar":            {"zh_CN", "th_TH"},
		"Cambodia":           {"zh_CN", "vi_VN"},
		"Thailand":           {"en_TH"},
		"Vietnam":            {"vi_VN"},
		"Papua New Guinea":   {"en_PH", "en_NZ"},
		"Solomon Islands":    {"en_PH", "en_NZ"},
		"New Caledonia":      {"en_PH", "en_NZ"},
		"Fiji":               {"en_PH", "en_NZ"},
		"French Polynesia":   {"en_PH", "en_NZ"},
		"Vanuatu":            {"en_PH", "en_NZ"},
		"Australia":          {"en_AU"},
		"New Zealand":        {"en_NZ"},
		"Chile":              {"es_MX"},
		"Colombia":           {"es_MX"},
		"Guatemala":          {"es_MX"},
		"Dominican Republic": {"es_MX"},
		"The Bahamas":        {"es_MX"},
		"El Salvador":        {"es_MX"},
		"Belize":             {"es_MX"},
		"Honduras":           {"es_MX"},
		"Trinidad":           {"es_MX"},
		"French Guiana":      {"es_MX", "fr_FR"},
		"Jamaica":            {"es_MX", "zu_ZA"},
		"Haiti":              {"es_MX", "zu_ZA"},
		"Costa Rica":         {"es_MX"},
		"Nicaragua":          {"es_MX"},
		"Panama":             {"es_MX"},
		"Cuba":               {"es_MX"},
		"Puerto Rico":        {"es_MX"},
		"Venezuela":          {"es_MX"},
		"Guyana":             {"es_MX"},
		"Peru":               {"es_MX"},
		"Paraguay":           {"es_MX"},
		"Sierra Leone":       {"es_MX"},
		"Uruguay":            {"pt_PT", "es_MX", "pt_BR"},
		"Azerbaijan":         {"uk_UA", "hy_AM", "az_AZ"},
		"Georgia":            {"uk_UA", "hy_AM", "az_AZ"},
		"Armenia":            {"hy_AM", "az_AZ"},
		"Ukraine":            {"uk_UA"},
		"Russia":             {"ru_RU"},
		"Belarus":            {"ru_RU"},
		"Tajikistan":         {"ar_SA", "ru_RU"},
		"Kyrgyzstan":         {"zh_CN", "ru_RU"},
		"Kazakhstan":         {"tr_TR", "ru_RU"},
		"Uzbekistan":         {"tr_TR", "ru_RU"},
		"Turkmenistan":       {"ar_SA", "ru_RU"},
		"Mongolia":           {"ru_RU", "zh_CN"},
		"Nepal":              {"zh_CN"},
		"Bangladesh":         {"en_IN"},
		"India":              {"en_IN"},
		"Pakistan":           {"id_ID", "en_IN"},
		"Ethiopia":           {"sa_SA", "zu_ZA"},
		"Chad":               {"sa_SA"},
		"Senegal":            {"sa_SA"},
		"Algeria":            {"sa_SA", "ar_EG"},
		"Togo":               {"sa_SA"},
		"Cameroon":           {"sa_SA"},
		"Eritrea":            {"sa_SA"},
		"Liberia":            {"sa_SA"},
		"Libya":              {"sa_SA", "ar_EG"},
		"Tanzania":           {"sa_SA"},
		"Guinea":             {"sa_SA"},
		"The Gambia":         {"sa_SA"},
		"Mali":               {"sa_SA"},
		"Niger":              {"sa_SA"},
		"Nigeria":            {"sa_SA"},
		"Benin":              {"sa_SA"},
		"Gabon":              {"sa_SA"},
		"Angola":             {"sa_SA"},
		"Malawi":             {"sa_SA"},
		"Namibia":            {"sa_SA"},
		"Botswana":           {"sa_SA"},
		"South Africa":       {"sa_SA"},
		"Zimbabwe":           {"sa_SA"},
		"Mozambique":         {"sa_SA"},
		"Madagascar":         {"sa_SA"},
		"Kenya":              {"sa_SA"},
		"Somalia":            {"sa_SA"},
		"Djibouti":           {"sa_SA"},
		"Sudan":              {"sa_SA"},
		"Rwanda":             {"sa_SA"},
		"Uganda":             {"sa_SA"},
		"DRC":                {"sa_SA"},
		"South Sudan":        {"sa_SA"},
		"Burundi":            {"sa_SA"},
		"Ivory Coast":        {"sa_SA"},
		"Tunisia":            {"sa_SA", "ar_EG"},
		"Zambia":             {"sa_SA"},
		"Morocco":            {"ar_EG", "sa_SA"},
		"Egypt":              {"ar_EG"},
		"Palestine":          {"ar_EG", "ar_SA"},
		"Israel":             {"ar_JO"},
		"Jordan":             {"ar_JO"},
		"Lebanon":            {"ar_EG", "ar_SA", "ar_JO"},
		"Iraq":               {"ar_EG", "ar_SA"},
		"Iran":               {"ar_EG", "ar_SA"},
		"Saudi Arabia":       {"ar_EG", "ar_SA"},
		"Kuwait":             {"ar_EG", "ar_SA"},
		"Syria":              {"ar_EG", "ar_SA"},
		"Bahrain":            {"ar_EG", "ar_SA"},
		"Qatar":              {"ar_EG", "ar_SA"},
		"UAE":                {"ar_EG", "ar_SA"},
		"Yemen":              {"ar_EG", "ar_SA"},
	}
	selectedCountryCodes := countryMap[country]
	if len(selectedCountryCodes) == 0 {
		fmt.Println(country)
	}
	code := util.PickFromStringList(countryMap[country])
	return code
}

func pickCountry() string {
	countries := []util.Locale{
		{Name: "USA", Weight: 62},
		{Name: "Canada", Weight: 50},
		{Name: "Sweden", Weight: 20},
		{Name: "Russia", Weight: 20},
		{Name: "Finland", Weight: 8},
		{Name: "Czech Republic", Weight: 6},
		{Name: "Slovakia", Weight: 5},
		{Name: "Germany", Weight: 3},
		{Name: "Switzerland", Weight: 2},
		{Name: "Latvia", Weight: 1},
		{Name: "Norway", Weight: 1},      // For smaller hockey nations
		{Name: "Denmark", Weight: 1},     // For smaller hockey nations
		{Name: "Netherlands", Weight: 1}, // For smaller hockey nations
		{Name: "Belarus", Weight: 1},     // For smaller hockey nations
		{Name: "Ukraine", Weight: 1},     // For smaller hockey nations
	}

	totalWeight := 0
	for _, country := range countries {
		totalWeight += country.Weight
	}

	randomWeight := util.GenerateIntFromRange(0, totalWeight)
	for _, country := range countries {
		if randomWeight < country.Weight {
			return country.Name
		}
		randomWeight -= country.Weight
	}
	return util.USA
}

func getInternationalNameMap() map[string]map[string][]string {
	path := filepath.Join(os.Getenv("ROOT"), "data", "unique_male_names_by_country.json")
	content := util.ReadJson(path)
	var payload map[string]map[string][]string

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}

func getAttributeValue(pos string, arch string, star int, attr string, blob map[string]map[string]map[string]map[string]interface{}) int {
	starStr := strconv.Itoa(star)
	switch pos {
	case "C":
		switch arch {
		case "Enforcer":
			switch attr {
			case "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Passing", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "LongShotPower", "CloseShotPower", "BodyChecking", "StickChecking", "GoalieVision", "Goalkeeping", "GoalieReboundControl":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Grinder":
			switch attr {
			case "OneTimer", "Passing", "Agility":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "LongShotPower", "CloseShotPower", "LongShotAccuracy", "CloseShotAccuracy", "PuckHandling", "GoalieVision", "Goalkeeping", "GoalieReboundControl", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Playmaker":
			switch attr {
			case "LongShotPower", "CloseShotPower", "BodyChecking", "StickChecking", "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Agility":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl", "Strength", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Power":
			switch attr {
			case "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Agility", "PuckHandling", "Passing":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl", "LongShotPower", "BodyChecking", "StickChecking", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Sniper":
			switch attr {
			case "CloseShotAccuracy", "OneTimer", "Agility", "PuckHandling":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl", "CloseShotPower", "BodyChecking", "StickChecking", "ShotBlocking", "Strength":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Two-Way":
			switch attr {
			case "LongShotPower", "CloseShotPower", "BodyChecking", "PuckHandling", "StickChecking", "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Agility", "Strength", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	case "F":
		switch arch {
		case "Enforcer":
			switch attr {
			case "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Passing", "Faceoffs", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "LongShotPower", "CloseShotPower", "BodyChecking", "StickChecking", "GoalieVision", "Goalkeeping", "GoalieReboundControl":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Grinder":
			switch attr {
			case "OneTimer", "Passing", "Agility", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "LongShotPower", "CloseShotPower", "LongShotAccuracy", "CloseShotAccuracy", "ShotBlocking", "PuckHandling", "GoalieVision", "Goalkeeping", "GoalieReboundControl":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Playmaker":
			switch attr {
			case "LongShotPower", "CloseShotPower", "BodyChecking", "StickChecking", "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Agility", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl", "Strength", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Power":
			switch attr {
			case "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Agility", "PuckHandling", "Passing", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl", "LongShotPower", "BodyChecking", "StickChecking", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Sniper":
			switch attr {
			case "CloseShotAccuracy", "OneTimer", "Agility", "PuckHandling", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl", "CloseShotPower", "BodyChecking", "StickChecking", "ShotBlocking", "Strength":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Two-Way":
			switch attr {
			case "LongShotPower", "CloseShotPower", "BodyChecking", "PuckHandling", "StickChecking", "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Agility", "Strength", "Faceoffs", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "GoalieVision", "Goalkeeping", "GoalieReboundControl":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	case "D":
		switch arch {
		case "Defensive":
			switch attr {
			case "OneTimer", "Passing", "Agility":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "LongShotAccuracy", "CloseShotAccuracy", "PuckHandling", "LongShotPower", "CloseShotPower", "StickChecking", "GoalieVision", "Goalkeeping", "GoalieReboundControl", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Enforcer":
			switch attr {
			case "LongShotAccuracy", "CloseShotAccuracy", "OneTimer", "Passing", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "PuckHandling", "LongShotPower", "CloseShotPower", "StickChecking", "GoalieVision", "Goalkeeping", "GoalieReboundControl", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Two-Way":
			switch attr {
			case "OneTimer", "Strength", "LongShotAccuracy", "CloseShotAccuracy", "ShotBlocking":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "PuckHandling", "LongShotPower", "CloseShotPower", "StickChecking", "GoalieVision", "Goalkeeping", "Agility", "GoalieReboundControl", "Faceoffs":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Offensive":
			switch attr {
			case "OneTimer", "Agility", "LongShotAccuracy", "CloseShotAccuracy":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "StickChecking", "BodyChecking", "Strength", "GoalieVision", "Goalkeeping", "ShotBlocking", "GoalieReboundControl", "Faceoffs", "CloseShotPower":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	case "G":
		switch arch {
		case "Stand-Up":
			switch attr {
			case "GoalieReboundControl", "Passing":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "StickChecking", "BodyChecking", "LongShotAccuracy", "LongShotPower", "Goalkeeping", "Agility", "ShotBlocking", "CloseShotAccuracy", "CloseShotPower", "Faceoffs", "PuckHandling", "OneTimer":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Hybrid":
			switch attr {
			case "GoalieReboundControl", "Passing", "Goalkeeping", "GoalieVision", "Agility", "Strength":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "StickChecking", "BodyChecking", "LongShotAccuracy", "LongShotPower", "ShotBlocking", "CloseShotAccuracy", "CloseShotPower", "Faceoffs", "PuckHandling", "OneTimer":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		case "Butterfly":
			switch attr {
			case "GoalieReboundControl", "Passing":
				return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
			case "StickChecking", "BodyChecking", "LongShotAccuracy", "LongShotPower", "GoalieVision", "Strength", "ShotBlocking", "CloseShotAccuracy", "CloseShotPower", "Faceoffs", "PuckHandling", "OneTimer":
				return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	}
	return util.GenerateIntFromRange(5, 15)
}

func getNationalityValueModifier(country, attr string, value int) int {
	if country == util.USA || country == util.Canada {
		return value
	}

	switch country {
	case util.Russia:
		if attr == "Aggression" {
			return value + 5
		} else if attr == util.CloseShotPower || attr == util.Strength || attr == util.BodyChecking {
			return value + 2
		} else if (attr == util.LongShotPower || attr == util.Agility || attr == util.StickChecking) && value > 2 {
			return value - 2
		}
	case util.Sweden:
		if (attr == util.CloseShotPower || attr == util.Agility) && value > 2 {
			return value - 2
		} else if attr == util.LongShotPower {
			return value + 2
		}
	case "Finland":
		if attr == "Aggression" {
			return value - 5
		} else if (attr == util.CloseShotPower || attr == util.Strength || attr == util.BodyChecking) && value > 2 {
			return value - 2
		} else if attr == util.LongShotPower || attr == util.Agility || attr == util.PuckHandling || attr == util.StickChecking {
			return value + 2
		}
	}

	return value
}

func getNationalityPotentialModifier(country, attr string, value int) int {
	if country == util.USA || country == util.Canada {
		return value
	}

	switch country {
	case util.Russia:
		switch attr {
		case util.CloseShotPower, util.Strength, util.BodyChecking:
			return value + 5
		case util.LongShotPower, util.Agility, util.StickChecking:
			return value - 5
		}
	case util.Sweden:
		switch attr {
		case util.CloseShotPower, util.Agility:
			return value - 5
		case util.LongShotPower:
			return value + 5
		}
	case "Finland":
		switch attr {
		case util.CloseShotPower, util.Strength, util.BodyChecking:
			return value - 5
		case util.LongShotPower, util.Agility, util.PuckHandling, util.StickChecking:
			return value + 5
		}
	}

	return value
}

func getValueFromInterfaceRange(star string, starMap map[string]interface{}) int {
	// Check if the key exists in the map
	u, exists := starMap[star]
	if !exists {
		fmt.Printf("Key '%s' not found in starMap.\n", star)
		return 0 // Return a default value
	}

	// Check if the value can be asserted as a slice of interfaces
	minMax, ok := u.([]interface{})
	if !ok {
		fmt.Printf("Value for key '%s' is not a slice of interfaces.\n", star)
		return 0 // Return a default value
	}

	// Ensure the slice has at least two elements
	if len(minMax) < 2 {
		fmt.Printf("Value for key '%s' does not have enough elements (expected at least 2).\n", star)
		return 0 // Return a default value
	}

	// Check if the first element is a float64
	min, ok := minMax[0].(float64)
	if !ok {
		fmt.Printf("First element of '%s' is not a float64.\n", star)
		return 0 // Return a default value
	}

	// Check if the second element is a float64
	max, ok := minMax[1].(float64)
	if !ok {
		fmt.Printf("Second element of '%s' is not a float64.\n", star)
		return 0 // Return a default value
	}

	// Generate a random value in the range [min, max]
	return util.GenerateIntFromRange(int(min), int(max))
}

func getClutchValue() int {
	clutchNum := util.GenerateIntFromRange(1, 1000)
	if clutchNum < 145 {
		return -1
	} else if clutchNum < 845 {
		return 0
	} else if clutchNum < 990 {
		return 1
	}
	return 2
}

func getCrootLocations(locale string) map[string][]structs.CrootLocation {
	path := filepath.Join(os.Getenv("ROOT"), "data", locale+".json")

	content := util.ReadJson(path)

	var payload map[string][]structs.CrootLocation
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getAttributeBlob() map[string]map[string]map[string]map[string]interface{} {
	path := filepath.Join(os.Getenv("ROOT"), "data", "attributeBlob.json")

	content := util.ReadJson(path)

	var payload map[string]map[string]map[string]map[string]interface{}
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getCollegeGenList(id uint) []structs.CollegeGenObj {
	return []structs.CollegeGenObj{{Year: 4, Pos: "C"}, {Year: 4, Pos: "F"}, {Year: 4, Pos: "F"}, {Year: 4, Pos: "D"}, {Year: 4, Pos: "D"}, {Year: 4, Pos: "G"}, {Year: 4, Pos: util.PickFromStringList([]string{"C", "F", "D"})},
		{Year: 3, Pos: "C"}, {Year: 3, Pos: "F"}, {Year: 3, Pos: "F"}, {Year: 3, Pos: "D"}, {Year: 3, Pos: "D"}, {Year: 3, Pos: "G"}, {Year: 3, Pos: util.PickFromStringList([]string{"C", "F", "D"})},
		{Year: 2, Pos: "C"}, {Year: 2, Pos: "F"}, {Year: 2, Pos: "F"}, {Year: 2, Pos: "D"}, {Year: 2, Pos: "D"}, {Year: 2, Pos: "G"}, {Year: 2, Pos: util.PickFromStringList([]string{"C", "F", "D"})},
		{Year: 1, Pos: "C"}, {Year: 1, Pos: "F"}, {Year: 1, Pos: "F"}, {Year: 1, Pos: "D"}, {Year: 1, Pos: "D"}, {Year: 1, Pos: "G"}, {Year: 1, Pos: util.PickFromStringList([]string{"C", "F", "D"})}}
}

func pickWalkonState(state string) string {
	switch state {
	case "AL":
		return util.PickFromStringList([]string{"AL", "LA", "MS", "TN", "GA", "FL"})
	case "AR":
		return util.PickFromStringList([]string{"AR", "LA", "MO", "TN", "TX"})
	case "AZ":
		return util.PickFromStringList([]string{"AZ", "NM", "CA"})
	case "CA":
		return util.PickFromStringList([]string{"CA", "AZ", "HI"})
	case "CO":
		return util.PickFromStringList([]string{"CO", "KS", "UT", "WY"})
	case "CT":
		return util.PickFromStringList([]string{"CT", "NY", "NJ", "RI"})
	case "DC":
		return util.PickFromStringList([]string{"DC", "MD", "VA"})
	case "FL":
		return util.PickFromStringList([]string{"FL", "GA", "AL"})
	case "GA":
		return util.PickFromStringList([]string{"GA", "FL", "SC", "AL"})
	case "HI":
		return util.PickFromStringList([]string{"HI"})
	case "IA":
		return util.PickFromStringList([]string{"IA", "MN", "WI", "NE"})
	case "ID":
		return util.PickFromStringList([]string{"ID", "WA", "UT"})
	case "IN":
		return util.PickFromStringList([]string{"IN", "IL", "OH", "MI", "AK"})
	case "IL":
		return util.PickFromStringList([]string{"IL", "IN", "WI", "MI"})
	case "KS":
		return util.PickFromStringList([]string{"KS", "MO", "NE"})
	case "KY":
		return util.PickFromStringList([]string{"KY", "OH", "TN"})
	case "LA":
		return util.PickFromStringList([]string{"LA", "TX", "MS"})
	case "MA":
		return util.PickFromStringList([]string{"MA", "CT", "RI", "NH", "VT", "ME"})
	case "MD":
		return util.PickFromStringList([]string{"DC", "MD", "VA", "DE"})
	case "MI":
		return util.PickFromStringList([]string{"MI", "OH", "IN"})
	case "MN":
		return util.PickFromStringList([]string{"MN", "WI", "IA"})
	case "MO":
		return util.PickFromStringList([]string{"MO", "AR", "KS"})
	case "MS":
		return util.PickFromStringList([]string{"MS", "LA", "AL"})
	case "MT":
		return util.PickFromStringList([]string{"MT", "ID", "WY"})
	case "NC":
		return util.PickFromStringList([]string{"NC", "SC", "VA"})
	case "ND":
		return util.PickFromStringList([]string{"ND", "SD", "MN"})
	case "NE":
		return util.PickFromStringList([]string{"NE", "KS", "SD", "IA"})
	case "NH":
		return util.PickFromStringList([]string{"NH", "VT", "ME", "MA"})
	case "NJ":
		return util.PickFromStringList([]string{"NJ", "DE", "NY", "CT", "PA"})
	case "NM":
		return util.PickFromStringList([]string{"NM", "AZ", "TX"})
	case "NV":
		return util.PickFromStringList([]string{"NV", "UT", "CA"})
	case "NY":
		return util.PickFromStringList([]string{"NY", "NJ", "PA", "CT"})
	case "OH":
		return util.PickFromStringList([]string{"OH", "KY", "MI", "PA"})
	case "OK":
		return util.PickFromStringList([]string{"OK", "TX", "KS", "AR"})
	case "OR":
		return util.PickFromStringList([]string{"OR", "WA", "CA"})
	case "PA":
		return util.PickFromStringList([]string{"PA", "NJ", "DE", "OH", "WV"})
	case "RI":
		return util.PickFromStringList([]string{"RI", "MA", "CT", "NY"})
	case "SC":
		return util.PickFromStringList([]string{"SC", "NC", "GA"})
	case "SD":
		return util.PickFromStringList([]string{"SD", "ND", "MN", "NE"})
	case "TN":
		return util.PickFromStringList([]string{"TN", "KY", "GA", "AL", "AR"})
	case "TX":
		return util.PickFromStringList([]string{"TX"})
	case "UT":
		return util.PickFromStringList([]string{"UT", "CO", "ID", "AZ"})
	case "VA":
		return util.PickFromStringList([]string{"VA", "WV", "DC", "MD"})
	case "WA":
		return util.PickFromStringList([]string{"WA", "OR", "ID", "AK"})
	case "WI":
		return util.PickFromStringList([]string{"WI", "MN", "IL", "MI"})
	case "WV":
		return util.PickFromStringList([]string{"WV", "PA", "VA"})
	case "WY":
		return util.PickFromStringList([]string{"WY", "CO", "UT", "MO", "ID"})
	case "BC":
		return util.PickFromStringList([]string{"BC", "AB", "YK"})
	}

	return "AK"
}
