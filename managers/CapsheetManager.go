package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllProCapsheets() []structs.ProCapsheet {
	return repository.FindCapsheetRecords()
}

func GetProCapsheetMap() map[uint]structs.ProCapsheet {
	capsheets := GetAllProCapsheets()
	return MakeCapsheetMap(capsheets)
}

func AllocateCapsheets() {
	db := dbprovider.GetInstance().GetDB()
	teams := repository.FindAllProTeams(repository.TeamClauses{})
	proCapSheetMap := GetProCapsheetMap()
	contractMap := GetContractMap()

	for _, team := range teams {
		TeamID := strconv.Itoa(int(team.ID))
		players := repository.FindAllProPlayers(repository.PlayerQuery{TeamID: TeamID})
		playerMap := MakeProfessionalPlayerMap(players)
		contracts := []structs.ProContract{}

		for _, p := range players {
			contract := contractMap[p.ID]
			contracts = append(contracts, contract)
		}

		Capsheet := proCapSheetMap[team.ID]

		if Capsheet.ID == 0 {
			Capsheet.AssignCapsheet(team.ID)
		}

		Capsheet.ResetCapsheet()

		sort.Slice(contracts, func(i, j int) bool {
			return contracts[i].Y1BaseSalary > contracts[j].Y1BaseSalary
		})
		window := 22
		for idx, contract := range contracts {
			if idx > window {
				break
			}
			player := playerMap[contract.PlayerID]
			if contract.IsCut || player.IsAffiliatePlayer {
				window += 1
				continue
			}

			Capsheet.AddContractToCapsheet(contract)
		}

		repository.SaveProCapsheetRecord(Capsheet, db)
	}
}
