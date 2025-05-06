package managers

import (
	"math/rand"

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
