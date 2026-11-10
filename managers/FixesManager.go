package managers

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func AddCompensationPicksMidRound() {
	db := dbprovider.GetInstance().GetDB()

	draftPicksToUpload := []structs.DraftPick{
		{Model: gorm.Model{
			ID: 1345,
		},
			DraftRound:     4,
			DraftNumber:    29,
			TeamID:         1,
			Team:           "MONT",
			OriginalTeamID: 1,
			OriginalTeam:   "MONT",
			DraftValue:     0,
			SeasonID:       3,
			Season:         2027,
			IsCompensation: true,
		},
		{Model: gorm.Model{
			ID: 1346,
		},
			DraftRound:     4,
			DraftNumber:    30,
			TeamID:         17,
			Team:           "CALG",
			OriginalTeamID: 17,
			OriginalTeam:   "CALG",
			DraftValue:     0,
			SeasonID:       3,
			Season:         2027,
			IsCompensation: true,
		},
	}

	repository.CreateDraftPickRecordsBatch(db, draftPicksToUpload, 2)
}
