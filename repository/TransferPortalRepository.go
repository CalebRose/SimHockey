package repository

import (
	"log"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type TransferPortalQuery struct {
	ID               string
	TeamID           string
	ProfileID        string
	CollegePlayerID  string
	RemovedFromBoard string
	OrderByPoints    bool
	IsActive         string
}

func FindTransferPortalProfileRecord(clauses TransferPortalQuery) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles structs.TransferPortalProfile

	query := db.Model(&profiles)

	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.RemovedFromBoard) > 0 {
		isRemoved := clauses.RemovedFromBoard == "Y"
		query = query.Where("removed_from_board = ?", isRemoved)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&profiles).Error; err != nil {
		return structs.TransferPortalProfile{}
	}

	return profiles
}

func FindTransferPortalProfileRecords(clauses TransferPortalQuery) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	query := db.Model(&profiles)

	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.RemovedFromBoard) > 0 {
		isRemoved := clauses.RemovedFromBoard == "Y"
		query = query.Where("removed_from_board = ?", isRemoved)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&profiles).Error; err != nil {
		return []structs.TransferPortalProfile{}
	}

	return profiles
}

func FindCollegePromiseRecord(clauses TransferPortalQuery) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := structs.CollegePromise{}

	query := db.Model(&p)

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if err := query.Find(&p).Error; err != nil {
		return structs.CollegePromise{}
	}
	return p
}

func FindCollegePromiseRecords(clauses TransferPortalQuery) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := []structs.CollegePromise{}

	query := db.Model(&p)

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.ID)
	}

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.IsActive) > 0 {
		isActive := clauses.IsActive == "Y"
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&p).Error; err != nil {
		return []structs.CollegePromise{}
	}
	return p
}

// Create
func CreateCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateTransferPortalProfileRecord(profile structs.TransferPortalProfile, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&profile).Error
	if err != nil {
		log.Panicln("Could not create new college recruit record")
	}
}

func CreateTransferPortalProfileRecordsBatch(db *gorm.DB, profiles []structs.TransferPortalProfile, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

// Save
func SaveCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Player Record
	err := db.Save(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func SaveTransferPortalProfileRecord(profile structs.TransferPortalProfile, db *gorm.DB) {
	// Save College Player Record
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

// Delete
func DeleteCollegePromise(promise structs.CollegePromise, db *gorm.DB) {
	err := db.Delete(&promise).Error
	if err != nil {
		log.Panicln("Could not delete old college promise record.")
	}
}
