package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllRecruitRecords() []structs.Recruit {
	return repository.FindAllRecruits(false, false, false, false, "")
}

func GetAllCrootRecords() []structs.Croot {
	recruits := repository.FindAllRecruits(true, false, false, false, "")

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
	return repository.FindAllRecruits(false, true, false, false, "")
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
	return repository.FindRecruitPlayerProfileRecords(id, true, false)
}

func GetOnlyRecruitProfilesByProfileID(id string) []structs.RecruitPlayerProfile {
	return repository.FindRecruitPlayerProfileRecords(id, false, false)
}

func GetSignedRecruitsByTeamID(teamID string) []structs.Recruit {
	return repository.FindAllRecruits(false, true, true, true, teamID)
}

func GetRecruitProfileMap() map[uint][]structs.RecruitPlayerProfile {
	profiles := repository.FindRecruitPlayerProfileRecords("", false, true)
	return MakeRecruitProfileMapByRecruitID(profiles)
}

func CreateRecruitingProfileForRecruit(recruitPointsDto structs.CreateRecruitProfileDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitEntry := repository.FindRecruitPlayerProfileRecord(strconv.Itoa(recruitPointsDto.RecruitID),
		strconv.Itoa(recruitPointsDto.ProfileID))

	if recruitEntry.RecruitID != 0 && recruitEntry.ProfileID != 0 {
		// Replace Recruit
		recruitEntry.ToggleRemoveFromBoard()
		repository.SaveRecruitProfileRecord(db, recruitEntry)
		return recruitEntry
	}

	modifier := CalculateModifierTowardsRecruit(recruitPointsDto.PlayerRecruit, recruitPointsDto.Team)

	createRecruitEntry := structs.RecruitPlayerProfile{
		SeasonID:           uint(recruitPointsDto.SeasonID),
		RecruitID:          uint(recruitPointsDto.RecruitID),
		ProfileID:          uint(recruitPointsDto.ProfileID),
		Modifier:           modifier,
		TotalPoints:        0,
		CurrentWeeksPoints: 0,
		SpendingCount:      0,
		Scholarship:        false,
		ScholarshipRevoked: false,
		RemovedFromBoard:   false,
		IsSigned:           false,
	}

	// Create
	repository.CreateRecruitProfileRecord(db, createRecruitEntry)

	return createRecruitEntry
}

func CalculateModifierTowardsRecruit(recruit structs.Recruit, team structs.CollegeTeam) float32 {
	programMod := CalculateMultiplier(uint(team.ProgramPrestige), uint(recruit.ProgramPref))
	professionalDevMod := CalculateMultiplier(uint(team.ProfessionalPrestige), uint(recruit.ProfDevPref))
	traditionsMod := CalculateMultiplier(uint(team.Traditions), uint(recruit.TraditionsPref))
	facilitiesMod := CalculateMultiplier(uint(team.Facilities), uint(recruit.FacilitiesPref))
	atmosphereMod := CalculateMultiplier(uint(team.Atmosphere), uint(recruit.AtmospherePref))
	academicsMod := CalculateMultiplier(uint(team.Academics), uint(recruit.AcademicsPref))
	conferenceMod := CalculateMultiplier(uint(team.ConferencePrestige), uint(recruit.ConferencePref))
	coachMod := CalculateMultiplier(uint(team.CoachRating), uint(recruit.CoachPref))
	seasonMod := CalculateMultiplier(uint(team.SeasonMomentum), uint(recruit.SeasonMomentumPref))

	return (programMod + professionalDevMod + traditionsMod + facilitiesMod + atmosphereMod + academicsMod + conferenceMod + coachMod + seasonMod) / 9
}

func CalculateBaseModifier(attr uint) float32 {
	return 1 + float32((attr)-5)/5
}

func CalculateAdjustmentFactor(teamAttr, playerPref uint) float32 {
	return 1 + float32((teamAttr-playerPref)/10)
}

func CalculateMultiplier(teamAttr uint, playerPref uint) float32 {
	baseMod := CalculateBaseModifier(teamAttr)
	adjFactor := CalculateAdjustmentFactor(teamAttr, playerPref)
	return baseMod * adjFactor
}
