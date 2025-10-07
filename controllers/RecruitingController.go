package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
)

// CreateRecruitPlayerProfile
func CreateRecruitPlayerProfile(w http.ResponseWriter, r *http.Request) {

	var recruitPointsDto structs.CreateRecruitProfileDto
	err := json.NewDecoder(r.Body).Decode(&recruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingProfile := managers.CreateRecruitingProfileForRecruit(recruitPointsDto)

	json.NewEncoder(w).Encode(recruitingProfile)
}

func SendScholarshipToRecruit(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitProfileDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile, _ := managers.SendScholarshipToRecruit(updateRecruitPointsDto)
	json.NewEncoder(w).Encode(recruitingPointsProfile)
}

func RemoveRecruitFromBoard(w http.ResponseWriter, r *http.Request) {
	var updateRecruitPointsDto structs.UpdateRecruitProfileDto
	err := json.NewDecoder(r.Body).Decode(&updateRecruitPointsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recruitingPointsProfile := managers.RemoveRecruitFromBoard(updateRecruitPointsDto)

	json.NewEncoder(w).Encode(recruitingPointsProfile)
}

func SaveRecruitingBoard(w http.ResponseWriter, r *http.Request) {
	var updateRecruitingBoardDto structs.UpdateRecruitingBoardDTO
	err := json.NewDecoder(r.Body).Decode(&updateRecruitingBoardDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ts := managers.GetTimestamp()

	if ts.IsRecruitingLocked {
		http.Error(w, "Recruiting is locked!", http.StatusNotAcceptable)
		return
	}

	result := make(chan structs.RecruitingTeamProfile)

	go func() {
		recruitingProfile := managers.UpdateRecruitingProfile(updateRecruitingBoardDto)
		result <- recruitingProfile
	}()

	crootProfile := <-result
	close(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(crootProfile)
}

func ToggleAIBehavior(w http.ResponseWriter, r *http.Request) {
	var updateRecruitingBoardDto structs.UpdateRecruitingBoardDTO
	err := json.NewDecoder(r.Body).Decode(&updateRecruitingBoardDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.SaveAIBehavior(updateRecruitingBoardDto.Profile)

	json.NewEncoder(w).Encode("AI Behavior Switched For Team " + updateRecruitingBoardDto.Profile.Team)
}

func ScoutAttribute(w http.ResponseWriter, r *http.Request) {
	var scoutAttributeDto structs.ScoutAttributeDTO
	err := json.NewDecoder(r.Body).Decode(&scoutAttributeDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := managers.ScoutAttribute(scoutAttributeDto)

	json.NewEncoder(w).Encode(profile)
}

func ScoutPortalAttribute(w http.ResponseWriter, r *http.Request) {
	var scoutAttributeDto structs.ScoutAttributeDTO
	err := json.NewDecoder(r.Body).Decode(&scoutAttributeDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := managers.ScoutPortalAttribute(scoutAttributeDto)

	json.NewEncoder(w).Encode(profile)
}
