package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/structs"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func CreateTSModelsFile(w http.ResponseWriter, r *http.Request) {
	converter := typescriptify.New().
		Add(structs.BasePlayer{}).
		Add(structs.GameResultsResponse{}).
		Add(structs.SearchStatsResponse{}).
		Add(structs.PlayerPreferences{}).
		Add(structs.BasePlayerProgressions{}).
		Add(structs.BasePotentials{}).
		Add(structs.BaseInjuryData{}).
		Add(structs.CollegePlayer{}).
		Add(structs.ProfessionalPlayer{}).
		Add(structs.BaseTeam{}).
		Add(structs.CollegeTeam{}).
		Add(structs.ProfessionalTeam{}).
		Add(structs.ProfessionalTeamFranchise{}).
		Add(structs.Arena{}).
		Add(structs.BootstrapData{}).
		Add(structs.CollegeStandings{}).
		Add(structs.Recruit{}).
		Add(structs.RecruitPlayerProfile{}).
		Add(structs.RecruitingTeamProfile{}).
		Add(structs.Croot{}).
		Add(structs.BaseRecruitingGrades{}).
		Add(structs.LeadingTeams{}).
		Add(structs.CreateRecruitProfileDto{}).
		Add(structs.UpdateRecruitProfileDto{}).
		Add(structs.ScoutAttributeDTO{}).
		Add(structs.CrootProfile{}).
		Add(structs.SimTeamBoardResponse{}).
		Add(structs.UpdateRecruitingBoardDTO{}).
		Add(structs.RecruitPointAllocation{}).
		Add(structs.RecruitingOdds{}).
		Add(structs.CollegeGame{}).
		Add(structs.CollegeLineup{}).
		Add(structs.CollegeShootoutLineup{}).
		Add(structs.ProfessionalStandings{}).
		Add(structs.ProfessionalGame{}).
		Add(structs.ProfessionalLineup{}).
		Add(structs.ProfessionalShootoutLineup{}).
		Add(structs.ProCapsheet{}).
		Add(structs.ProContract{}).
		Add(structs.FreeAgencyOffer{}).
		Add(structs.WaiverOffer{}).
		Add(structs.ExtensionOffer{}).
		Add(structs.TeamRequest{}).
		Add(structs.CollegeTeamRequest{}).
		Add(structs.ProTeamRequest{}).
		Add(structs.CollegePollSubmission{}).
		Add(structs.CollegePollOfficial{}).
		Add(structs.PollDataResponse{}).
		Add(structs.ExtensionOffer{}).
		Add(structs.FreeAgencyOfferDTO{}).
		Add(structs.WaiverOffer{}).
		Add(structs.WaiverOfferDTO{}).
		Add(structs.ShootoutPlayerIDs{}).
		Add(structs.Allocations{}).
		Add(structs.BaseLineup{}).
		Add(structs.TeamRecordResponse{}).
		Add(structs.TopPlayer{}).
		Add(structs.InboxResponse{}).
		Add(structs.BasePlayerStats{}).
		Add(structs.BaseTeamStats{}).
		Add(structs.TeamSeasonStats{}).
		Add(structs.CollegePlayerSeasonStats{}).
		Add(structs.CollegePlayerGameStats{}).
		Add(structs.CollegeTeamSeasonStats{}).
		Add(structs.CollegeTeamGameStats{}).
		Add(structs.ProfessionalPlayerSeasonStats{}).
		Add(structs.ProfessionalPlayerGameStats{}).
		Add(structs.ProfessionalTeamSeasonStats{}).
		Add(structs.ProfessionalTeamGameStats{}).
		Add(structs.NewsLog{}).
		Add(structs.Notification{}).
		Add(structs.Timestamp{}).
		Add(structs.TeamRequestsResponse{})
	err := converter.ConvertToFile("ts/models.ts")
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode("Models ran!")
}
