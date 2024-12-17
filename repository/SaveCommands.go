package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func SaveTimestamp(ts structs.Timestamp, db *gorm.DB) {
	err := db.Save(&ts).Error
	if err != nil {
		log.Panicln("Could not save timestamp")
	}
}

func SaveCollegeHockeyPlayerRecord(playerRecord structs.CollegePlayer, db *gorm.DB) {
	err := db.Save(&playerRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(playerRecord.ID)))
	}
}

func SaveCollegeHockeyRecruitRecord(recruitRecord structs.Recruit, db *gorm.DB) {
	err := db.Save(&recruitRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(recruitRecord.ID)))
	}
}

func SaveCollegeLineupRecord(lineupRecord structs.CollegeLineup, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(lineupRecord.ID)))
	}
}
