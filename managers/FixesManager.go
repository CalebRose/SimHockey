package managers

import (
	util "github.com/CalebRose/SimHockey/_util"
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

func AddPreferenceValuesToRecruits() {
	db := dbprovider.GetInstance().GetDB()

	recruits := repository.FindAllRecruits(false, false, false, true, false, "")
	for _, recruit := range recruits {
		campusLife := util.GenerateIntFromRange(1, 9)
		media := util.GenerateIntFromRange(1, 9)
		recruit.CampusLifePref = uint8(campusLife)
		recruit.MediaPref = uint8(media)
		repository.SaveCollegeHockeyRecruitRecord(recruit, db)
	}
}

func FixModifierValuesForRecruits() {
	db := dbprovider.GetInstance().GetDB()

	recruits := repository.FindAllRecruits(false, false, false, true, false, "")
	recruitProfiles := repository.FindRecruitPlayerProfileRecords("", "", false, false, false)
	recruitProfilesMap := MakeRecruitProfileMapByRecruitID(recruitProfiles)
	recruitingTeamProfiles := repository.FindTeamRecruitingProfiles(false)
	recruitingTeamProfilesMap := MakeTeamProfileMap(recruitingTeamProfiles)
	collegeTeams := repository.FindAllCollegeTeams(repository.TeamClauses{})
	collegeTeamsMap := MakeCollegeTeamMap(collegeTeams)

	for _, recruit := range recruits {
		profiles := recruitProfilesMap[recruit.ID]
		if len(profiles) == 0 {
			continue
		}
		for _, profile := range profiles {
			teamProfile := recruitingTeamProfilesMap[profile.ProfileID]
			if teamProfile == nil {
				continue
			}
			// Apply any necessary fixes to the recruit based on the profile and teamProfile
			team := collegeTeamsMap[teamProfile.ID]

			modifier := CalculateModifierTowardsRecruit(recruit.PlayerPreferences, team)
			profile.Modifier = modifier

			repository.SaveRecruitProfileRecord(db, profile)
		}
	}
}
