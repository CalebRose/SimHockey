package repository

import (
	"errors"
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func FindAllRecruits(includeProfiles, includeIsSigned, isSigned, orderByOverall, forRecruitingPage bool, teamID string) []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit

	query := db.Model(&recruits)

	if orderByOverall {
		query = query.Order("overall desc")
	}

	if forRecruitingPage {
		query = query.Order("stars desc").Order("composite_rank desc")
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

	if err := query.Order("composite_rank desc").Find(&recruits).Error; err != nil {
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

func FindRecruitPlayerProfileRecords(profileID, recruitID string, includeRecruit, orderByOverall, removeFromBoard bool) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	query := db.Model(&croots)

	if includeRecruit {
		query = query.Preload("Recruit")
	}

	if len(profileID) > 0 {
		query = query.Where("profile_id = ?", profileID)
	}

	if len(recruitID) > 0 {
		query = query.Where("recruit_id = ?", recruitID)
	}

	if removeFromBoard {
		query = query.Where("removed_from_board = ?", false)
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
	if err := db.Create(&record).Error; err != nil {
		return err
	}

	return nil
}

func SaveRecruitProfileRecord(db *gorm.DB, record structs.RecruitPlayerProfile) error {
	record.Recruit = structs.Recruit{}
	if err := db.Save(&record).Error; err != nil {
		return err
	}

	return nil
}

func FindTeamRecruitingProfiles(aiOnly bool) []structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.RecruitingTeamProfile

	query := db.Model(&profiles)

	if aiOnly {
		query = query.Where("is_ai = ?", true)
	}

	if err := query.Find(&profiles).Error; err != nil {
		return []structs.RecruitingTeamProfile{}
	}

	return profiles
}

func FindTeamRecruitingProfile(teamID string, includeRecruit, includeSignedRecruitsOnly bool) structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	query := db.Model(&profile)

	if len(teamID) > 0 {
		query = query.Where("id = ?", teamID)
	}

	if includeRecruit {
		query = query.Preload("Recruits.Recruit.RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
			return db.Order("total_points DESC").Where("total_points > 0")
		})
	} else if includeSignedRecruitsOnly {
		query = query.Preload("Recruits.Recruit", func(db *gorm.DB) *gorm.DB {
			return db.Order("total_points DESC").Where("team_id = ? AND is_signed = true", teamID)
		})
	}

	if err := query.Find(&profile).Error; err != nil {
		return structs.RecruitingTeamProfile{}
	}

	return profile
}

func CreateTeamProfileRecord(db *gorm.DB, record structs.RecruitingTeamProfile) error {
	if err := db.Create(&record).Error; err != nil {
		return err
	}

	return nil
}

func SaveTeamProfileRecord(db *gorm.DB, record structs.RecruitingTeamProfile) error {
	record.Recruits = nil
	if err := db.Save(&record).Error; err != nil {
		return err
	}

	return nil
}

func CreatePointAllocationRecord(db *gorm.DB, record structs.RecruitPointAllocation) error {
	if err := db.Create(&record).Error; err != nil {
		return err
	}

	return nil
}

func SaveCollegeHockeyRecruitRecord(recruitRecord structs.Recruit, db *gorm.DB) {
	err := db.Save(&recruitRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(recruitRecord.ID)))
	}
}
