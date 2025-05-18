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

func SaveCollegePollSubmissionRecord(pollRecord structs.CollegePollSubmission, db *gorm.DB) {
	err := db.Save(&pollRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(pollRecord.ID)))
	}
}

func SaveCollegeStandingsRecord(standingsRecord structs.CollegeStandings, db *gorm.DB) {
	err := db.Save(&standingsRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(standingsRecord.ID)))
	}
}

func SaveProfessionalStandingsRecord(standingsRecord structs.ProfessionalStandings, db *gorm.DB) {
	err := db.Save(&standingsRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(standingsRecord.ID)))
	}
}
