package managers

import (
	"os"
	"path/filepath"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func ImportTeams() {
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
		reputation := 100
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
			ProgramPrestige:      uint8(program),
			ProfessionalPrestige: uint8(profDev),
			Traditions:           uint8(traditions),
			Atmosphere:           uint8(atmosphere),
			Facilities:           uint8(facilities),
			AcademicPrestige:     uint8(academicPrestige),
			ConferencePrestige:   5,
			CoachReputation:      uint8(reputation),
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
						AZSlapshot:    20,
						AZWristshot:   20,
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
						AZSlapshot:    20,
						AZWristshot:   20,
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
						AZSlapshot:    20,
						AZWristshot:   20,
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

	GenerateTestRosters()
}
