package repository

import (
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func CreateHockeyRecruitRecordsBatch(db *gorm.DB, croots []structs.Recruit, batchSize int) error {
	total := len(croots)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(croots[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegeHockeyPlayerRecordsBatch(db *gorm.DB, players []structs.CollegePlayer, batchSize int) error {
	total := len(players)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(players[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProHockeyPlayerRecordsBatch(db *gorm.DB, players []structs.ProfessionalPlayer, batchSize int) error {
	total := len(players)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(players[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateGlobalPlayerRecordsBatch(db *gorm.DB, croots []structs.GlobalPlayer, batchSize int) error {
	total := len(croots)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(croots[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegeTeamRecordsBatch(db *gorm.DB, teams []structs.CollegeTeam, batchSize int) error {
	total := len(teams)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(teams[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProTeamRecordsBatch(db *gorm.DB, teams []structs.ProfessionalTeam, batchSize int) error {
	total := len(teams)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(teams[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateArenaRecordsBatch(db *gorm.DB, teams []structs.Arena, batchSize int) error {
	total := len(teams)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(teams[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegePollRecord(db *gorm.DB, poll structs.CollegePollOfficial) error {
	if err := db.Create(poll).Error; err != nil {
		return err
	}

	return nil
}

func CreateCollegePollSubmissionRecord(db *gorm.DB, poll structs.CollegePollSubmission) error {
	if err := db.Create(poll).Error; err != nil {
		return err
	}

	return nil
}
