package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllRecruitRecords() []structs.Recruit {
	return repository.FindAllRecruits(false, false, false, false, false, "")
}

func GetAllCrootRecords() []structs.Croot {
	recruits := repository.FindAllRecruits(true, false, false, false, true, "")

	var croots []structs.Croot
	for _, recruit := range recruits {
		var croot structs.Croot
		croot.Map(recruit)

		croots = append(croots, croot)
	}

	sort.Sort(structs.ByCrootRank(croots))

	return croots
}

func GetAllUnsignedRecruits() []structs.Recruit {
	return repository.FindAllRecruits(false, true, false, false, true, "")
}

func GetCollegeRecruitByID(id string) structs.Recruit {
	return repository.FindCollegeRecruitRecord(id, false)
}

func GetCrootRecordByID(id string) structs.Croot {
	recruit := repository.FindCollegeRecruitRecord(id, true)
	var croot structs.Croot
	croot.Map(recruit)
	return croot
}

func GetRecruitProfilesWithRecruitByProfileID(id string) []structs.RecruitPlayerProfile {
	return repository.FindRecruitPlayerProfileRecords(id, "", true, false, false)
}

func GetOnlyRecruitProfilesByProfileID(id string) []structs.RecruitPlayerProfile {
	return repository.FindRecruitPlayerProfileRecords(id, "", false, false, false)
}

func GetSignedRecruitsByTeamID(teamID string) []structs.Recruit {
	return repository.FindAllRecruits(false, true, true, true, false, teamID)
}

func GetRecruitProfileMap() map[uint][]structs.RecruitPlayerProfile {
	profiles := repository.FindRecruitPlayerProfileRecords("", "", false, true, false)
	return MakeRecruitProfileMapByRecruitID(profiles)
}

func CreateRecruitingProfileForRecruit(dto structs.CreateRecruitProfileDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()
	teamID := strconv.Itoa(dto.ProfileID)
	recruitEntry := repository.FindRecruitPlayerProfileRecord(strconv.Itoa(dto.RecruitID),
		teamID)

	if recruitEntry.RecruitID != 0 && recruitEntry.ProfileID != 0 {
		// Replace Recruit
		recruitEntry.ToggleRemoveFromBoard()
		repository.SaveRecruitProfileRecord(db, recruitEntry)
		return recruitEntry
	}

	modifier := CalculateModifierTowardsRecruit(dto.PlayerRecruit.PlayerPreferences, dto.Team)

	// Find CTH Value
	state := dto.PlayerRecruit.State
	closeToHome := dto.PlayerRecruit.Country == util.USA && state == dto.Team.State
	isPipeline := false

	if !closeToHome {
		// Check for pipeline
		currentRoster := repository.FindAllCollegePlayers(repository.PlayerQuery{TeamID: teamID})
		piplineMap := make(map[string]int)

		for _, p := range currentRoster {
			key := p.Country
			if p.Country == util.USA || p.Country == util.Canada {
				key = p.State
			}
			if piplineMap[key] > 0 {
				piplineMap[key] = piplineMap[key] + 1
			} else {
				piplineMap[key] = 1
			}
		}

		playerKey := dto.PlayerRecruit.Country
		if playerKey == util.USA || playerKey == util.Canada {
			playerKey = dto.PlayerRecruit.State
		}

		if piplineMap[playerKey] > 6 {
			isPipeline = true
		}
	}

	createRecruitEntry := structs.RecruitPlayerProfile{
		SeasonID:           uint(dto.SeasonID),
		RecruitID:          uint(dto.RecruitID),
		ProfileID:          uint(dto.ProfileID),
		Modifier:           modifier,
		TotalPoints:        0,
		CurrentWeeksPoints: 0,
		SpendingCount:      0,
		Scholarship:        false,
		ScholarshipRevoked: false,
		RemovedFromBoard:   false,
		IsSigned:           false,
		IsHomeState:        closeToHome,
		IsPipelineState:    isPipeline,
	}

	// Create
	repository.CreateRecruitProfileRecord(db, createRecruitEntry)

	return createRecruitEntry
}

func CalculateModifierTowardsRecruit(prefs structs.PlayerPreferences, team structs.CollegeTeam) float32 {
	programMod := calculateMultiplier(uint(team.ProgramPrestige), uint(prefs.ProgramPref))
	professionalDevMod := calculateMultiplier(uint(team.ProfessionalPrestige), uint(prefs.ProfDevPref))
	traditionsMod := calculateMultiplier(uint(team.Traditions), uint(prefs.TraditionsPref))
	facilitiesMod := calculateMultiplier(uint(team.Facilities), uint(prefs.FacilitiesPref))
	atmosphereMod := calculateMultiplier(uint(team.Atmosphere), uint(prefs.AtmospherePref))
	academicsMod := calculateMultiplier(uint(team.Academics), uint(prefs.AcademicsPref))
	conferenceMod := calculateMultiplier(uint(team.ConferencePrestige), uint(prefs.ConferencePref))
	coachMod := calculateMultiplier(uint(team.CoachRating), uint(prefs.CoachPref))
	seasonMod := calculateMultiplier(uint(team.SeasonMomentum), uint(prefs.SeasonMomentumPref))

	return (programMod + professionalDevMod + traditionsMod + facilitiesMod + atmosphereMod + academicsMod + conferenceMod + coachMod + seasonMod) / 9
}

func calculateBaseModifier(attr int) float32 {
	return 1 + float32(attr-5)/5
}

func calculateAdjustmentFactor(teamAttr, playerPref int) float32 {
	return 1 + float32((teamAttr-playerPref)/10)
}

func calculateMultiplier(teamAttr uint, playerPref uint) float32 {
	baseMod := calculateBaseModifier(int(teamAttr))
	adjFactor := calculateAdjustmentFactor(int(teamAttr), int(playerPref))
	return baseMod * adjFactor
}

func SendScholarshipToRecruit(dto structs.UpdateRecruitProfileDto) (structs.RecruitPlayerProfile, structs.RecruitingTeamProfile) {
	db := dbprovider.GetInstance().GetDB()

	teamProfile := repository.FindTeamRecruitingProfile(strconv.Itoa(dto.ProfileID), false, false)

	if teamProfile.ScholarshipsAvailable == 0 {
		log.Panicln("\nTeamId: " + strconv.Itoa(dto.ProfileID) + " does not have any availabe scholarships")
	}

	crootProfile := repository.FindRecruitPlayerProfileRecord(
		strconv.Itoa(dto.RecruitID),
		strconv.Itoa(dto.ProfileID),
	)

	if !crootProfile.Scholarship && !crootProfile.ScholarshipRevoked {
		teamProfile.SubtractScholarshipsAvailable()
		crootProfile.ToggleScholarship(true, false)
	} else {
		teamProfile.ReallocateScholarship()
		crootProfile.ToggleScholarship(false, true)
	}

	repository.SaveRecruitProfileRecord(db, crootProfile)
	repository.SaveTeamProfileRecord(db, teamProfile)

	return crootProfile, teamProfile
}

func RevokeScholarshipFromRecruit(dto structs.UpdateRecruitProfileDto) (structs.RecruitPlayerProfile, structs.RecruitingTeamProfile) {
	db := dbprovider.GetInstance().GetDB()

	teamProfile := repository.FindTeamRecruitingProfile(strconv.Itoa(dto.ProfileID), false, false)

	crootProfile := repository.FindRecruitPlayerProfileRecord(
		strconv.Itoa(dto.RecruitID),
		strconv.Itoa(dto.ProfileID),
	)

	if !crootProfile.Scholarship {
		fmt.Printf("%s", "\nCannot revoke an inexistant scholarship from Recruit "+strconv.Itoa(int(crootProfile.RecruitID)))
		return crootProfile, teamProfile
	}

	// recruitingPointsProfile.ToggleScholarship()
	teamProfile.ReallocateScholarship()

	repository.SaveRecruitProfileRecord(db, crootProfile)
	repository.SaveTeamProfileRecord(db, teamProfile)

	return crootProfile, teamProfile
}

func GetRecruitFromRecruitsList(id uint, recruits []structs.RecruitPlayerProfile) structs.RecruitPlayerProfile {
	var recruit structs.RecruitPlayerProfile

	for i := range recruits {
		if recruits[i].RecruitID == id {
			recruit = recruits[i]
			break
		}
	}

	return recruit
}

func GetRecruitingTeamBoardByTeamID(teamID string) structs.SimTeamBoardResponse {

	profile := repository.FindTeamRecruitingProfile(teamID, true, false)

	var teamProfileResponse structs.SimTeamBoardResponse
	var crootProfiles []structs.CrootProfile

	for i := 0; i < len(profile.Recruits); i++ {
		var crootProfile structs.CrootProfile
		var croot structs.Croot

		croot.Map(profile.Recruits[i].Recruit)

		crootProfile.Map(profile.Recruits[i], croot)

		crootProfiles = append(crootProfiles, crootProfile)
	}

	sort.Sort(structs.ByCrootProfileTotal(crootProfiles))

	teamProfileResponse.Map(profile, crootProfiles)

	return teamProfileResponse
}

func GetRecruitingClassByTeamID(teamID string) structs.SimTeamBoardResponse {

	profile := repository.FindTeamRecruitingProfile(teamID, false, true)

	var teamProfileResponse structs.SimTeamBoardResponse
	var crootProfiles []structs.CrootProfile

	for i := range profile.Recruits {
		var crootProfile structs.CrootProfile
		var croot structs.Croot

		croot.Map(profile.Recruits[i].Recruit)

		crootProfile.Map(profile.Recruits[i], croot)

		crootProfiles = append(crootProfiles, crootProfile)
	}

	sort.Sort(structs.ByCrootProfileTotal(crootProfiles))

	teamProfileResponse.Map(profile, crootProfiles)

	return teamProfileResponse
}

func GetAllTeamProfilesForSync() []structs.RecruitingTeamProfile {
	return repository.FindTeamRecruitingProfiles(false)
}

func RemoveRecruitFromBoard(updateRecruitPointsDto structs.UpdateRecruitProfileDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := repository.FindRecruitPlayerProfileRecord(
		strconv.Itoa(updateRecruitPointsDto.RecruitID),
		strconv.Itoa(updateRecruitPointsDto.ProfileID))

	if recruitProfile.RemovedFromBoard {
		log.Panicln("Recruit has already been removed from Team Recruiting Board.")
	}

	recruitProfile.ToggleRemoveFromBoard()
	repository.SaveRecruitProfileRecord(db, recruitProfile)

	return recruitProfile
}

func UpdateRecruitingProfile(updateRecruitingBoardDto structs.UpdateRecruitingBoardDTO) structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var teamID = strconv.Itoa(updateRecruitingBoardDto.TeamID)

	var teamProfile = repository.FindTeamRecruitingProfile(teamID, false, false)

	var recruitProfiles = repository.FindRecruitPlayerProfileRecords(teamID, "", false, false, false)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	var currentPoints float32 = 0

	for i := range recruitProfiles {
		updatedRecruit := GetRecruitFromRecruitsList(recruitProfiles[i].RecruitID, updatedRecruits)
		currentPoints += updatedRecruit.CurrentWeeksPoints
		teamProfile.AllocateSpentPoints(currentPoints)
		if recruitProfiles[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {
			if teamProfile.SpentPoints <= teamProfile.WeeklyPoints {
				recruitProfiles[i].AllocateCurrentWeekPoints(updatedRecruit.CurrentWeeksPoints)
				fmt.Println("Saving recruit " + strconv.Itoa(int(recruitProfiles[i].RecruitID)))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(int(teamProfile.TeamID)) + " than what is allowed.")
			}
			// Save Recruit Profile
			repository.SaveRecruitProfileRecord(db, recruitProfiles[i])
		}
	}

	// Save team recruiting profile
	repository.SaveTeamProfileRecord(db, teamProfile)

	return teamProfile
}

func GetRecruitingClassSizeForTeams() {
	db := dbprovider.GetInstance().GetDB()
	profiles := GetAllTeamProfilesForSync()

	for _, team := range profiles {
		count := 0

		players := repository.FindAllCollegePlayers(
			repository.PlayerQuery{TeamID: strconv.Itoa(int(team.ID))})

		rosterSize := len(players)

		for _, player := range players {
			if (player.Year == 4 && !player.IsRedshirt) || (player.Year == 5 && player.IsRedshirt) && player.Stars > 0 {
				count++
			}
		}

		rosterMinusGrads := rosterSize - count

		if rosterMinusGrads+25 > 105 {
			count = 105 - rosterMinusGrads
		} else if rosterMinusGrads+25 < 85 {
			count = 85 - rosterMinusGrads
		} else {
			count = 25
		}

		team.SetRecruitingClassSize(uint8(count))

		repository.SaveTeamProfileRecord(db, team)
	}
}

// SaveAIBehavior -- Toggle whether a Team will use AI recruiting or not
func SaveAIBehavior(profile structs.RecruitingTeamProfile) {
	db := dbprovider.GetInstance().GetDB()
	TeamID := strconv.Itoa(int(profile.TeamID))
	recruitingProfile := repository.FindTeamRecruitingProfile(TeamID, false, false)
	recruitingProfile.UpdateAIBehavior(profile.IsAI, profile.AIStarMax, profile.AIStarMin, profile.AIMinThreshold, profile.AIMaxThreshold)
	repository.SaveTeamProfileRecord(db, recruitingProfile)
}

func ScoutAttribute(dto structs.ScoutAttributeDTO) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitID := strconv.Itoa(int(dto.RecruitID))
	profileID := strconv.Itoa(int(dto.ProfileID))

	teamProfile := repository.FindTeamRecruitingProfile(profileID, false, false)

	recruitProfile := repository.FindRecruitPlayerProfileRecord(recruitID, profileID)

	if teamProfile.ID == 0 || recruitProfile.ID == 0 {
		log.Panic("ERROR: IDs PROVIDED DON'T LINE UP")
	}

	if teamProfile.WeeklyScoutingPoints == 0 {
		log.Panic("ERROR: User doesn't have enough scouting points")
	}

	recruitProfile.ApplyScoutingAttribute(dto.Attribute)

	teamProfile.SubtractScoutingPoints()

	repository.SaveTeamProfileRecord(db, teamProfile)
	repository.SaveRecruitProfileRecord(db, recruitProfile)

	return recruitProfile
}

func ScoutPortalAttribute(dto structs.ScoutAttributeDTO) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	playerID := strconv.Itoa(int(dto.RecruitID))
	profileID := strconv.Itoa(int(dto.ProfileID))

	teamProfile := repository.FindTeamRecruitingProfile(profileID, false, false)

	portalProfile := repository.FindTransferPortalProfileRecord(repository.TransferPortalQuery{CollegePlayerID: playerID, ProfileID: profileID})

	if teamProfile.ID == 0 || portalProfile.ID == 0 {
		return portalProfile
	}

	if teamProfile.WeeklyScoutingPoints == 0 {
		return portalProfile
	}

	portalProfile.ApplyScoutingAttribute(dto.Attribute)

	teamProfile.SubtractScoutingPoints()

	repository.SaveTeamProfileRecord(db, teamProfile)
	repository.SaveTransferPortalProfileRecord(portalProfile, db)

	return portalProfile
}
