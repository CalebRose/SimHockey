package managers

import (
	"fmt"
	"os"
	"path/filepath"

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
	proTeams := repository.FindAllProTeams()

	for _, team := range collegeTeams {
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

	for _, team := range proTeams {
		standings := structs.ProfessionalStandings{
			BaseStandings: structs.BaseStandings{
				TeamID:       team.ID,
				TeamName:     team.TeamName,
				SeasonID:     ts.SeasonID,
				Season:       ts.Season,
				LeagueID:     1,
				ConferenceID: uint(team.ConferenceID),
			},
			DivisionID: uint(team.DivisionID),
		}

		db.Create(&standings)
	}
}

func ImportTeamRecruitingProfiles() {
	db := dbprovider.GetInstance().GetDB()

	teams := repository.FindAllCollegeTeams()

	for _, team := range teams {
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
