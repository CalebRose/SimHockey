package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func GetCHLTeamByTeamIDForDiscord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetCHLTeamDataForDiscord(teamID)
	json.NewEncoder(w).Encode(team)
}

func GetPHLTeamByTeamIDForDiscord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetPHLTeamDataForDiscord(teamID)
	json.NewEncoder(w).Encode(team)
}

func GetCHLPlayerForDiscord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCollegePlayerViaDiscord(id)

	json.NewEncoder(w).Encode(player)
}

func GetCHLPlayerByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	teamID := vars["abbr"]

	if len(firstName) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCollegePlayerByNameViaDiscord(firstName, lastName, teamID)

	json.NewEncoder(w).Encode(player)
}

func GetPHLPlayerForDiscord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetProPlayerViaDiscord(id)

	json.NewEncoder(w).Encode(player)
}

func GetPHLPlayerByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	teamID := vars["abbr"]

	if len(firstName) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetProPlayerByNameViaDiscord(firstName, lastName, teamID)

	json.NewEncoder(w).Encode(player)
}

func CompareCHLTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOneID := vars["teamOneID"]
	if len(teamOneID) == 0 {
		panic("User did not provide teamID")
	}

	teamTwoID := vars["teamTwoID"]
	if len(teamTwoID) == 0 {
		panic("User did not provide teamID")
	}

	res := managers.CompareCHLTeams(teamOneID, teamTwoID)

	json.NewEncoder(w).Encode(res)
}

func ComparePHLTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOneID := vars["teamOneID"]
	if len(teamOneID) == 0 {
		panic("User did not provide teamID")
	}

	teamTwoID := vars["teamTwoID"]
	if len(teamTwoID) == 0 {
		panic("User did not provide teamID")
	}

	res := managers.CompareCHLTeams(teamOneID, teamTwoID)

	json.NewEncoder(w).Encode(res)
}

func AssignDiscordIDToCHLTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	discordID := vars["discordID"]
	if len(teamID) == 0 {
		panic("User did not provide conference name")
	}

	managers.AssignDiscordIDToCHLTeam(teamID, discordID)
}

func AssignDiscordIDToPHLTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	discordID := vars["discordID"]
	if len(teamID) == 0 {
		panic("User did not provide conference name")
	}

	managers.AssignDiscordIDToPHLTeam(teamID, discordID)
}

func GetCHLRecruitingClassByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	recruitingProfile := managers.GetRecruitingClassByTeamID(teamID)

	json.NewEncoder(w).Encode(recruitingProfile)
}

func GetCHLRecruitViaDiscord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		panic("User did not provide a first name")
	}

	recruit := managers.GetCollegeRecruitViaDiscord(id)

	json.NewEncoder(w).Encode(recruit)
}
