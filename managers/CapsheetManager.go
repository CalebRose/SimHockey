package managers

import (
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
