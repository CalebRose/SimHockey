package managers

import (
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
	filePath := filepath.Join(os.Getenv("ROOT"), "data", "phl_teams_test.csv")
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
		colorOne := team[16]
		colorTwo := team[17]
		colorThree := team[18]

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
	repository.CreateProTeamRecordsBatch(db, teams, 20)
	repository.CreateProfessionalLineupRecordsBatch(db, proLineups, 50)
}
