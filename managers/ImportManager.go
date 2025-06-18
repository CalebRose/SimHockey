package managers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func ImportCollegeTeams() {
	db := dbprovider.GetInstance().GetDB()
	filePath := filepath.Join(os.Getenv("ROOT"), "data", "simchl_teams.csv")
	teamsCSV := util.ReadCSV(filePath)
	teams := []structs.CollegeTeam{}
	arenas := []structs.Arena{}
	collegeLineups := []structs.CollegeLineup{}
	for idx, team := range teamsCSV {
		if idx < 2 {
			continue
		}
		id := util.ConvertStringToInt(team[0])
		teamName := team[1]
		mascot := team[2]
		abbr := team[3]
		conferenceId := util.ConvertStringToInt(team[4])
		conference := team[5]
		coach := team[6]
		city := team[7]
		state := team[8]
		country := team[9]
		arenaId := util.ConvertStringToInt(team[10])
		arenaName := team[11]
		capacity := util.ConvertStringToInt(team[12])
		record := 0
		firstYear := 2025
		discordId := ""
		colorOne := team[16]
		colorTwo := team[17]
		colorThree := team[18]
		program := util.ConvertStringToInt(team[25])
		profDev := util.ConvertStringToInt(team[26])
		traditions := util.ConvertStringToInt(team[27])
		facilities := util.ConvertStringToInt(team[28])
		atmosphere := util.ConvertStringToInt(team[29])
		academicPrestige := util.ConvertStringToInt(team[31])

		arena := structs.Arena{
			Name:             arenaName,
			City:             city,
			State:            state,
			Country:          country,
			Capacity:         uint16(capacity),
			RecordAttendance: 0,
		}

		team := structs.CollegeTeam{
			BaseTeam: structs.BaseTeam{
				TeamName:         teamName,
				Mascot:           mascot,
				ConferenceID:     uint8(conferenceId),
				Conference:       conference,
				City:             city,
				State:            state,
				Country:          country,
				Abbreviation:     abbr,
				ArenaID:          uint16(arenaId),
				Coach:            coach,
				RecordAttendance: uint16(record),
				ArenaCapacity:    uint16(capacity),
				FirstPlayed:      uint16(firstYear),
				DiscordID:        discordId,
				ColorOne:         colorOne,
				ColorTwo:         colorTwo,
				ColorThree:       colorThree,
			},
			ProfileAttributes: structs.ProfileAttributes{
				ProgramPrestige:      uint8(program),
				ProfessionalPrestige: uint8(profDev),
				Traditions:           uint8(traditions),
				Atmosphere:           uint8(atmosphere),
				Facilities:           uint8(facilities),
				Academics:            uint8(academicPrestige),
				ConferencePrestige:   5,
				CoachRating:          5,
				SeasonMomentum:       5,
			},
		}
		arenas = append(arenas, arena)
		teams = append(teams, team)

		for i := 1; i < 5; i++ {
			forwardLineup := structs.CollegeLineup{
				BaseLineup: structs.BaseLineup{
					TeamID:   uint(id),
					LineType: 1,
					Line:     uint8(i),
					Allocations: structs.Allocations{
						AGZShot:       20,
						AGZPass:       20,
						AGZStickCheck: 20,
						AGZBodyCheck:  20,
						AZShot:        20,
						AZPass:        20,
						AZAgility:     20,
						AZStickCheck:  20,
						AZBodyCheck:   20,
						NPass:         20,
						NAgility:      20,
						NStickCheck:   20,
						NBodyCheck:    20,
						DZPass:        20,
						DZAgility:     20,
						DZStickCheck:  20,
						DZBodyCheck:   20,
						DGZPass:       20,
						DGZAgility:    20,
						DGZStickCheck: 20,
						DGZBodyCheck:  20,
					},
				},
			}
			collegeLineups = append(collegeLineups, forwardLineup)
		}
		for i := 1; i < 4; i++ {
			defenderLineup := structs.CollegeLineup{
				BaseLineup: structs.BaseLineup{
					TeamID:   uint(id),
					LineType: 2,
					Line:     uint8(i),
					Allocations: structs.Allocations{
						AGZShot:       20,
						AGZPass:       20,
						AGZStickCheck: 20,
						AGZBodyCheck:  20,
						AZShot:        20,
						AZPass:        20,
						AZAgility:     20,
						AZStickCheck:  20,
						AZBodyCheck:   20,
						NPass:         20,
						NAgility:      20,
						NStickCheck:   20,
						NBodyCheck:    20,
						DZPass:        20,
						DZAgility:     20,
						DZStickCheck:  20,
						DZBodyCheck:   20,
						DGZPass:       20,
						DGZAgility:    20,
						DGZStickCheck: 20,
						DGZBodyCheck:  20,
					},
				},
			}
			collegeLineups = append(collegeLineups, defenderLineup)
		}
		for i := 1; i < 3; i++ {
			goalieLineup := structs.CollegeLineup{
				BaseLineup: structs.BaseLineup{
					TeamID:   uint(id),
					LineType: 3,
					Line:     uint8(i),
					Allocations: structs.Allocations{
						AGZShot:       20,
						AGZPass:       20,
						AGZStickCheck: 20,
						AGZBodyCheck:  20,
						AZShot:        20,
						AZPass:        20,
						AZAgility:     20,
						AZStickCheck:  20,
						AZBodyCheck:   20,
						NPass:         20,
						NAgility:      20,
						NStickCheck:   20,
						NBodyCheck:    20,
						DZPass:        20,
						DZAgility:     20,
						DZStickCheck:  20,
						DZBodyCheck:   20,
						DGZPass:       20,
						DGZAgility:    20,
						DGZStickCheck: 20,
						DGZBodyCheck:  20,
					},
				},
			}
			collegeLineups = append(collegeLineups, goalieLineup)
		}
	}

	repository.CreateArenaRecordsBatch(db, arenas, 30)
	repository.CreateCollegeTeamRecordsBatch(db, teams, 30)
	repository.CreateCollegeLineupRecordsBatch(db, collegeLineups, 50)

	GenerateInitialRosters()
}

func ImportProTeams() {
	db := dbprovider.GetInstance().GetDB()
	filePath := filepath.Join(os.Getenv("ROOT"), "data", "phl_teams.csv")
	teamsCSV := util.ReadCSV(filePath)
	teams := []structs.ProfessionalTeam{}
	arenas := []structs.Arena{}
	proLineups := []structs.ProfessionalLineup{}
	for idx, team := range teamsCSV {
		if idx < 1 {
			continue
		}
		id := util.ConvertStringToInt(team[0])
		teamName := team[1]
		mascot := team[2]
		abbr := team[3]
		conferenceId := util.ConvertStringToInt(team[4])
		conference := team[5]
		divisionID := util.ConvertStringToInt(team[6])
		division := team[7]
		owner := ""
		gm := ""
		coach := ""
		scout := ""
		marketing := ""
		city := team[13]
		state := team[14]
		country := team[15]
		arenaId := util.ConvertStringToInt(team[16])
		arenaName := team[17]
		capacity := util.ConvertStringToInt(team[18])
		record := 0
		firstYear := 2025
		discordId := ""
		colorOne := team[22]
		colorTwo := team[23]
		colorThree := team[24]

		arena := structs.Arena{
			Name:             arenaName,
			City:             city,
			State:            state,
			Country:          country,
			Capacity:         uint16(capacity),
			RecordAttendance: 0,
		}

		team := structs.ProfessionalTeam{
			Owner:      owner,
			GM:         gm,
			Scout:      scout,
			Marketing:  marketing,
			DivisionID: uint8(divisionID),
			Division:   division,
			BaseTeam: structs.BaseTeam{
				TeamName:         teamName,
				Mascot:           mascot,
				ConferenceID:     uint8(conferenceId),
				Conference:       conference,
				City:             city,
				State:            state,
				Country:          country,
				Abbreviation:     abbr,
				ArenaID:          uint16(arenaId),
				Coach:            coach,
				RecordAttendance: uint16(record),
				ArenaCapacity:    uint16(capacity),
				FirstPlayed:      uint16(firstYear),
				DiscordID:        discordId,
				ColorOne:         colorOne,
				ColorTwo:         colorTwo,
				ColorThree:       colorThree,
			},
		}
		arenas = append(arenas, arena)
		teams = append(teams, team)

		for i := 1; i < 5; i++ {
			forwardLineup := structs.ProfessionalLineup{
				BaseLineup: structs.BaseLineup{
					TeamID:   uint(id),
					LineType: 1,
					Line:     uint8(i),
					Allocations: structs.Allocations{
						AGZShot:       20,
						AGZPass:       20,
						AGZStickCheck: 20,
						AGZBodyCheck:  20,
						AZShot:        20,
						AZPass:        20,
						AZAgility:     20,
						AZStickCheck:  20,
						AZBodyCheck:   20,
						NPass:         20,
						NAgility:      20,
						NStickCheck:   20,
						NBodyCheck:    20,
						DZPass:        20,
						DZAgility:     20,
						DZStickCheck:  20,
						DZBodyCheck:   20,
						DGZPass:       20,
						DGZAgility:    20,
						DGZStickCheck: 20,
						DGZBodyCheck:  20,
					},
				},
			}
			proLineups = append(proLineups, forwardLineup)
		}
		for i := 1; i < 4; i++ {
			defenderLineup := structs.ProfessionalLineup{
				BaseLineup: structs.BaseLineup{
					TeamID:   uint(id),
					LineType: 2,
					Line:     uint8(i),
					Allocations: structs.Allocations{
						AGZShot:       20,
						AGZPass:       20,
						AGZStickCheck: 20,
						AGZBodyCheck:  20,
						AZShot:        20,
						AZPass:        20,
						AZAgility:     20,
						AZStickCheck:  20,
						AZBodyCheck:   20,
						NPass:         20,
						NAgility:      20,
						NStickCheck:   20,
						NBodyCheck:    20,
						DZPass:        20,
						DZAgility:     20,
						DZStickCheck:  20,
						DZBodyCheck:   20,
						DGZPass:       20,
						DGZAgility:    20,
						DGZStickCheck: 20,
						DGZBodyCheck:  20,
					},
				},
			}
			proLineups = append(proLineups, defenderLineup)
		}
		for i := 1; i < 3; i++ {
			goalieLineup := structs.ProfessionalLineup{
				BaseLineup: structs.BaseLineup{
					TeamID:   uint(id),
					LineType: 3,
					Line:     uint8(i),
					Allocations: structs.Allocations{
						AGZShot:       20,
						AGZPass:       20,
						AGZStickCheck: 20,
						AGZBodyCheck:  20,
						AZShot:        20,
						AZPass:        20,
						AZAgility:     20,
						AZStickCheck:  20,
						AZBodyCheck:   20,
						NPass:         20,
						NAgility:      20,
						NStickCheck:   20,
						NBodyCheck:    20,
						DZPass:        20,
						DZAgility:     20,
						DZStickCheck:  20,
						DZBodyCheck:   20,
						DGZPass:       20,
						DGZAgility:    20,
						DGZStickCheck: 20,
						DGZBodyCheck:  20,
					},
				},
			}
			proLineups = append(proLineups, goalieLineup)
		}
	}

	repository.CreateArenaRecordsBatch(db, arenas, 20)
	repository.CreateProTeamRecordsBatch(db, teams, 24)
	repository.CreateProfessionalLineupRecordsBatch(db, proLineups, 50)
}

func ImportProRosters() {
	db := dbprovider.GetInstance().GetDB()
	filePath := filepath.Join(os.Getenv("ROOT"), "data", "init_pro_rosters.csv")
	rosterCSV := util.ReadCSV(filePath)
	teams := repository.FindAllProTeams()
	proPlayers := repository.FindAllProPlayers(repository.PlayerQuery{})
	proPlayerMap := MakeProfessionalPlayerMap(proPlayers)
	contracts := []structs.ProContract{}

	teamMap := make(map[string]structs.ProfessionalTeam)

	for _, team := range teams {
		teamMap[team.Abbreviation] = team
	}

	for idx, row := range rosterCSV {
		if idx < 1 {
			continue
		}
		playerID := util.ConvertStringToInt(row[0])
		player := proPlayerMap[uint(playerID)]
		teamAbbr := row[1]
		team := teamMap[teamAbbr]
		salary := row[2]
		salaryNum := util.ConvertStringToFloat(salary)
		if team.ID == 0 {
			fmt.Println("ERROR!")
		}
		player.AssignTeam(team.ID, team.Abbreviation)
		contract := structs.ProContract{
			PlayerID:       player.ID,
			TeamID:         team.ID,
			OriginalTeamID: team.ID,
			ContractLength: 3,
			Y1BaseSalary:   float32(salaryNum),
			Y2BaseSalary:   float32(salaryNum),
			Y3BaseSalary:   float32(salaryNum),
			ContractType:   "Auction",
			IsActive:       true,
			ContractValue:  float32(salaryNum) * 3,
		}
		contracts = append(contracts, contract)
		repository.SaveProPlayerRecord(player, db)
	}

	for _, c := range contracts {
		repository.CreateProContractRecord(db, c)
	}
}

func ImportStandingsForNewSeason() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	collegeTeams := repository.FindAllCollegeTeams()
	// proTeams := repository.FindAllProTeams()

	for _, team := range collegeTeams {
		if team.ID < 67 {
			continue
		}
		standings := structs.CollegeStandings{
			BaseStandings: structs.BaseStandings{
				TeamID:       team.ID,
				TeamName:     team.TeamName,
				SeasonID:     ts.SeasonID,
				Season:       ts.Season,
				LeagueID:     1,
				ConferenceID: uint(team.ConferenceID),
			},
		}

		db.Create(&standings)
	}

	// for _, team := range proTeams {
	// 	standings := structs.ProfessionalStandings{
	// 		BaseStandings: structs.BaseStandings{
	// 			TeamID:       team.ID,
	// 			TeamName:     team.TeamName,
	// 			SeasonID:     ts.SeasonID,
	// 			Season:       ts.Season,
	// 			LeagueID:     1,
	// 			ConferenceID: uint(team.ConferenceID),
	// 		},
	// 		DivisionID: uint(team.DivisionID),
	// 	}

	// 	db.Create(&standings)
	// }
}

func ImportTeamRecruitingProfiles() {
	db := dbprovider.GetInstance().GetDB()

	teams := repository.FindAllCollegeTeams()

	for _, team := range teams {
		if team.ID < 67 {
			continue
		}
		teamProfile := structs.RecruitingTeamProfile{
			TeamID:                team.ID,
			Team:                  team.Abbreviation,
			State:                 team.State,
			Country:               team.Country,
			ScholarshipsAvailable: 30,
			WeeklyPoints:          50,
			WeeklyScoutingPoints:  30,
			SpentPoints:           0,
			TotalCommitments:      0,
			RecruitClassSize:      7,
			PortalReputation:      100,
			IsUserTeam:            team.IsUserCoached,
			AIMinThreshold:        uint8(util.GenerateIntFromRange(4, 7)),
			AIMaxThreshold:        uint8(util.GenerateIntFromRange(8, 14)),
			AIStarMin:             uint8(util.GenerateIntFromRange(1, 2)),
			AIStarMax:             uint8(util.GenerateIntFromRange(3, 5)),
			Recruiter:             team.Coach,
		}

		db.Create(&teamProfile)
	}
}

func AddFAPreferences() {
	db := dbprovider.GetInstance().GetDB()

	players := repository.FindAllProPlayers(repository.PlayerQuery{})

	for _, p := range players {
		marketDR := util.GenerateIntFromRange(1, 100)
		competeDR := util.GenerateIntFromRange(1, 100)
		financeDR := util.GenerateIntFromRange(1, 100)
		m := uint8(1)
		c := uint8(1)
		f := uint8(1)
		if marketDR > 70 {
			// Generate new market preference
			marketList := []uint8{util.MarketLarge, util.MarketNotLarge, util.MarketSmall, util.MarketNotSmall, util.MarketLoyal}
			if p.Country == util.USA || p.Country == util.Canada {
				marketList = append(marketList, util.MarketCTH)
			} else if p.Country != util.USA && p.Country != util.Canada {
				marketList = append(marketList, util.MarketCountrymen)
			}

			maxNumber := len(marketList) - 1
			index := util.GenerateIntFromRange(0, maxNumber)
			m = marketList[index]
		}

		if competeDR > 70 {
			// Generate new compete preference
			competeList := []uint8{util.CompetitiveFirstLine, util.CompetitiveSecondLine}
			if p.Age <= 24 {
				competeList = append(competeList, util.CompetitiveMentorship)
			}
			if p.Age >= 28 {
				competeList = append(competeList, util.CompetitiveVeteranMentor)
			}

			maxNumber := len(competeList) - 1
			index := util.GenerateIntFromRange(0, maxNumber)
			c = competeList[index]
		}

		if financeDR > 70 {
			// Generate new finance preference
			competeList := []uint8{util.FinancialShort, util.FinancialLong, util.FinancialLargeAAV}

			maxNumber := len(competeList) - 1
			index := util.GenerateIntFromRange(0, maxNumber)
			c = competeList[index]
		}

		p.AssignPreferences(m, c, f)
		repository.SaveProPlayerRecord(p, db)
	}
}

func ImportCHLSchedule() {
	db := dbprovider.GetInstance().GetDB()
	filePath := filepath.Join(os.Getenv("ROOT"), "data", "gen", "2025_chl_schedule.csv")
	gamesCSV := util.ReadCSV(filePath)
	ts := GetTimestamp()
	games := []structs.CollegeGame{}
	collegeTeamMap := make(map[string]structs.CollegeTeam)
	chlTeams := GetAllCollegeTeams()

	for _, team := range chlTeams {
		collegeTeamMap[team.Abbreviation] = team
	}
	baseSeason := (ts.Season - 2000) * 100

	for idx, row := range gamesCSV {
		if idx < 1 {
			continue
		}
		week := util.ConvertStringToInt(row[3])
		day := row[4]
		homeTeam := collegeTeamMap[row[5]]
		awayTeam := collegeTeamMap[row[6]]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		conferenceGame := util.ConvertStringToBool(row[7])

		game := structs.CollegeGame{
			BaseGame: structs.BaseGame{
				SeasonID:      ts.SeasonID,
				WeekID:        baseSeason + uint(week),
				Week:          week,
				HomeTeamID:    homeTeam.ID,
				HomeTeam:      homeTeam.Abbreviation,
				HomeTeamCoach: homeTeam.Coach,
				AwayTeamID:    awayTeam.ID,
				AwayTeam:      awayTeam.Abbreviation,
				AwayTeamCoach: awayTeam.Coach,
				City:          homeTeam.City,
				State:         homeTeam.State,
				Country:       homeTeam.Country,
				ArenaID:       uint(homeTeam.ArenaID),
				Arena:         homeTeam.Arena,
				GameDay:       day,
				IsConference:  conferenceGame,
			},
		}
		games = append(games, game)
	}

	repository.CreateCHLGamesRecordsBatch(db, games, 100)
}

type ScheduledGame struct {
	structs.ProfessionalGame
	Round int
	Slot  string // "A", "B", or "C"
}

func (sg *ScheduledGame) ApplyRoundAndSlot(round int, slot string) {
	sg.Round = round
	sg.Slot = slot
}

func ImportPHLSeasonSchedule() {
	db := dbprovider.GetInstance().GetDB()
	ts := repository.FindTimestamp()
	phlTeams := repository.FindAllProTeams()
	schedule, err := GenerateSeasonSchedule(phlTeams, ts.SeasonID, 2000)
	if err != nil {
		log.Println("Generation Failed: ", err)
		return
	}
	finalUpload := []structs.ProfessionalGame{}
	for _, game := range schedule {
		game.AddWeekData((ts.Season*1000 + uint(game.Round)), uint(game.Round), game.Slot)
		finalUpload = append(finalUpload, game.ProfessionalGame)

	}
	repository.CreatePHLGamesRecordsBatch(db, finalUpload, 100)
}

func GenerateSeasonSchedule(
	teams []structs.ProfessionalTeam,
	seasonID uint,
	maxPartitionAttempts int,
) ([]ScheduledGame, error) {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	exceptionsMap := make(map[string]uint)

	exceptionsMap["CHI-MNS"] = 4
	exceptionsMap["COL-KCS"] = 4
	exceptionsMap["CALG-EDM"] = 4
	exceptionsMap["DET-ATL"] = 4
	exceptionsMap["FLA-NASH"] = 4
	exceptionsMap["PHI-PIT"] = 4
	exceptionsMap["CBJ-NYR"] = 4
	exceptionsMap["TOR-OTT"] = 4
	exceptionsMap["MONT-QUE"] = 4
	exceptionsMap["SJ-CAL"] = 4
	exceptionsMap["ANA-VGK"] = 4
	exceptionsMap["SEA-VAN"] = 4

	var master []ScheduledGame
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			teamA := teams[i]
			teamB := teams[j]
			if teamA.ID == teamB.ID {
				continue
			}
			numGames := 0
			if teamA.DivisionID == teamB.DivisionID {
				key := teamA.Abbreviation + "-" + teamB.Abbreviation
				secondKey := teamB.Abbreviation + "-" + teamA.Abbreviation
				if exceptionsMap[key] > 0 || exceptionsMap[secondKey] > 0 {
					numGames = 4
				} else {
					numGames = 3
				}
			} else {
				numGames = 2
			}
			half := numGames / 2
			extra := numGames % 2
			for x := 0; x < half+extra; x++ {
				g := generateProGameRecord(teamA, teamB, seasonID)
				master = append(master, ScheduledGame{ProfessionalGame: g})
			}
			// Team B hosts half
			for x := 0; x < half; x++ {
				g := generateProGameRecord(teamB, teamA, seasonID)
				master = append(master, ScheduledGame{ProfessionalGame: g})
			}
		}
	}
	totalRounds := 52
	counts := make(map[uint]int)
	for _, g := range master {
		counts[g.HomeTeamID]++
		counts[g.AwayTeamID]++
	}
	for _, t := range teams {
		if counts[t.ID] != totalRounds {
			log.Printf("⚠️ team %d has %d games (expected %d)\n",
				t.ID, counts[t.ID], totalRounds)
		}
	}

	rng.Shuffle(len(master), func(i, j int) {
		master[i], master[j] = master[j], master[i]
	})

	teamIDs := make([]uint, len(teams))
	for i, t := range teams {
		teamIDs[i] = t.ID
	}
	matchesPerRound := len(teams) / 2 // 24 teams → 12 games per round

	rounds, err := greedyPartition(rng, master, totalRounds, matchesPerRound, maxPartitionAttempts)
	if err != nil {
		return nil, err
	}
	var final []ScheduledGame
	for rIdx, roundGames := range rounds {
		roundNum := rIdx + 1
		var week int
		var slot string

		if roundNum < 52 {
			week = (roundNum + 2) / 3 // ceil(round/3)
			switch roundNum % 3 {
			case 1:
				slot = "A"
			case 2:
				slot = "B"
			case 0:
				slot = "C"
			}
		} else {
			week = 18
			slot = "A"
		}

		for _, g := range roundGames {
			g.ApplyRoundAndSlot(week, slot)
			final = append(final, g)
		}
	}

	return final, nil
}

func greedyPartition(
	rng *rand.Rand,
	allGames []ScheduledGame,
	totalRounds, matchesPerRound, perRoundMaxAttempts int,
) ([][]ScheduledGame, error) {
	// start with full pool
	remaining := append([]ScheduledGame(nil), allGames...)
	rounds := make([][]ScheduledGame, 0, totalRounds)

	for round := 1; round <= totalRounds; round++ {
		var thisRound []ScheduledGame
		success := false

		// keep an immutable snapshot for retries
		orig := append([]ScheduledGame(nil), remaining...)

		for attempt := 1; attempt <= perRoundMaxAttempts; attempt++ {
			// copy the snapshot
			pool := append([]ScheduledGame(nil), orig...)
			// shuffle the copy
			rng.Shuffle(len(pool), func(i, j int) {
				pool[i], pool[j] = pool[j], pool[i]
			})

			used := make(map[uint]bool, matchesPerRound*2)
			candidate := make([]ScheduledGame, 0, matchesPerRound)

			// greedily pull non‑conflicting games
			for i := 0; i < len(pool) && len(candidate) < matchesPerRound; {
				g := pool[i]
				if !used[g.HomeTeamID] && !used[g.AwayTeamID] {
					used[g.HomeTeamID] = true
					used[g.AwayTeamID] = true
					candidate = append(candidate, g)
					// remove from pool
					pool = append(pool[:i], pool[i+1:]...)
				} else {
					i++
				}
			}

			if len(candidate) == matchesPerRound {
				// success: lock in this round and remaining pool
				thisRound = candidate
				remaining = pool
				success = true
				break
			}
			// else: try again with fresh snapshot
		}

		if !success {
			return nil, fmt.Errorf(
				"round %d: failed to fill after %d attempts (got %d/%d games)",
				round, perRoundMaxAttempts, len(thisRound), matchesPerRound,
			)
		}

		rounds = append(rounds, thisRound)
	}

	if len(remaining) != 0 {
		return nil, fmt.Errorf(
			"leftover games after %d rounds: %d",
			totalRounds, len(remaining),
		)
	}
	return rounds, nil
}

func generateProGameRecord(teamA structs.ProfessionalTeam, teamB structs.ProfessionalTeam, seasonID uint) structs.ProfessionalGame {
	return structs.ProfessionalGame{
		BaseGame: structs.BaseGame{
			SeasonID:     seasonID,
			HomeTeamID:   teamA.ID,
			HomeTeam:     teamA.Abbreviation,
			AwayTeamID:   teamB.ID,
			AwayTeam:     teamB.Abbreviation,
			ArenaID:      uint(teamA.ArenaID),
			Arena:        teamA.Arena,
			City:         teamA.City,
			State:        teamA.State,
			Country:      teamA.Country,
			IsConference: teamA.ConferenceID == teamB.ConferenceID,
		},
		IsDivisional: teamA.DivisionID == teamB.DivisionID,
	}
}

func FixSeasonStatTables() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))
	_, collegeGameType := ts.GetCurrentGameType(true)
	_, proGameType := ts.GetCurrentGameType(true)
	collegePlayerSeasonStatMap := make(map[uint]*structs.CollegePlayerSeasonStats)
	proPlayerSeasonStatMap := make(map[uint]*structs.ProfessionalPlayerSeasonStats)
	collegeTeamSeasonStatMap := make(map[uint]*structs.CollegeTeamSeasonStats)
	proTeamSeasonStatMap := make(map[uint]*structs.ProfessionalTeamSeasonStats)

	collegePlayerGameStats := repository.FindCollegePlayerGameStatsRecords(seasonId, collegeGameType, "")
	collegeTeamGameStats := repository.FindCollegeTeamGameStatsRecords(seasonId, collegeGameType)
	proPlayerGameStats := repository.FindProPlayerGameStatsRecords(seasonId, proGameType, "")
	proTeamGameStats := repository.FindProTeamGameStatsRecords(seasonId, proGameType)

	for _, stat := range collegePlayerGameStats {
		if stat.GameType == 1 {
			continue
		}
		if _, ok := collegePlayerSeasonStatMap[stat.PlayerID]; !ok {
			collegePlayerSeasonStatMap[stat.PlayerID] = &structs.CollegePlayerSeasonStats{}
		}
		collegePlayerSeasonStatMap[stat.PlayerID].AddStatsToSeasonRecord(stat.BasePlayerStats)
	}

	for _, stat := range collegeTeamGameStats {
		if stat.GameType == 1 {
			continue
		}
		if _, ok := collegeTeamSeasonStatMap[stat.TeamID]; !ok {
			collegeTeamSeasonStatMap[stat.TeamID] = &structs.CollegeTeamSeasonStats{}
		}
		collegeTeamSeasonStatMap[stat.TeamID].BaseTeamStats.AddStatsToSeasonRecord(stat.BaseTeamStats)
		collegeTeamSeasonStatMap[stat.TeamID].TeamSeasonStats.AddStatsToSeasonRecord(stat.BaseTeamStats, false, false)
	}

	for _, stat := range proPlayerGameStats {
		if stat.GameType == 1 {
			continue
		}
		if _, ok := proPlayerSeasonStatMap[stat.PlayerID]; !ok {
			proPlayerSeasonStatMap[stat.PlayerID] = &structs.ProfessionalPlayerSeasonStats{}
		}
		proPlayerSeasonStatMap[stat.PlayerID].AddStatsToSeasonRecord(stat.BasePlayerStats)
	}

	for _, stat := range proTeamGameStats {
		if stat.GameType == 1 {
			continue
		}
		if _, ok := proTeamSeasonStatMap[stat.TeamID]; !ok {
			proTeamSeasonStatMap[stat.TeamID] = &structs.ProfessionalTeamSeasonStats{}
		}
		proTeamSeasonStatMap[stat.TeamID].BaseTeamStats.AddStatsToSeasonRecord(stat.BaseTeamStats)
		proTeamSeasonStatMap[stat.TeamID].TeamSeasonStats.AddStatsToSeasonRecord(stat.BaseTeamStats, false, false)
	}

	for key := range collegePlayerSeasonStatMap {
		value := collegePlayerSeasonStatMap[key]
		repository.SaveCollegePlayerSeasonStatsRecord(*value, db)
	}

	for key := range collegeTeamSeasonStatMap {
		value := collegeTeamSeasonStatMap[key]
		repository.SaveCollegeTeamSeasonStatsRecord(*value, db)
	}

	for key := range proPlayerSeasonStatMap {
		value := proPlayerSeasonStatMap[key]
		repository.SaveProPlayerSeasonStatsRecord(*value, db)
	}

	for key := range proTeamSeasonStatMap {
		value := proTeamSeasonStatMap[key]
		repository.SaveProTeamSeasonStatsRecord(*value, db)
	}
}
