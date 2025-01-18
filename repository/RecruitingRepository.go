package repository

import (
	"errors"
	"log"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func FindAllRecruits(includeProfiles, includeIsSigned, isSigned, orderByOverall bool, teamID string) []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit

	query := db.Model(&recruits)

	if orderByOverall {
		query.Order("overall desc")
	}

	if includeProfiles {
		query = query.Preload("RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
			return db.Order("total_points DESC")
		})
	}

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if includeIsSigned {
		query = query.Where("is_signed = ?", isSigned)
	}

	if err := query.Find(&recruits).Error; err != nil {
		return []structs.Recruit{}
	}

	return recruits
}

func FindCollegeRecruitRecord(id string, includePlayerProfiles bool) structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var recruits structs.Recruit

	query := db.Model(&recruits)

	if includePlayerProfiles {
		query = query.Preload("RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
			return db.Order("total_points DESC").Where("total_points > ?", "0")
		})
	}

	if len(id) > 0 {
		query = query.Where("id = ?", id)
	}

	if err := query.Find(&recruits).Error; err != nil {
		return structs.Recruit{}
	}

	return recruits
}

func FindRecruitPlayerProfileRecords(profileID string, includeRecruit, orderByOverall bool) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	query := db.Model(&croots)

	if includeRecruit {
		query = query.Preload("Recruit")
	}

	if len(profileID) > 0 {
		query = query.Where("profile_id = ?", profileID)
	}

	if orderByOverall {
		query = query.Order("overall desc")
	}

	if err := query.Find(&croots).Error; err != nil {
		return []structs.RecruitPlayerProfile{}
	}

	return croots
}

func FindRecruitProfileRecordsForAIPointSync(ProfileID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	err := db.Preload("Recruit", func(db *gorm.DB) *gorm.DB {
		return db.Order("stars DESC")
	}).Where("profile_id = ? AND removed_from_board = ?", ProfileID, false).Order("total_points DESC").Find(&croots).Error
	if err != nil {
		log.Fatal(err)
	}

	return croots
}

func FindRecruitPlayerProfileRecord(recruitID, profileID string) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croot structs.RecruitPlayerProfile
	err := db.Where("recruit_id = ? and profile_id = ?", recruitID, profileID).Find(&croot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.RecruitPlayerProfile{}
		} else {
			log.Fatal(err)
		}
	}
	return croot
}

func CreateRecruitProfileRecord(db *gorm.DB, record structs.RecruitPlayerProfile) error {
	if err := db.Create(record).Error; err != nil {
		return err
	}

	return nil
}

func SaveRecruitProfileRecord(db *gorm.DB, record structs.RecruitPlayerProfile) error {
	record.Recruit = structs.Recruit{}
	if err := db.Save(record).Error; err != nil {
		return err
	}

	return nil
}
