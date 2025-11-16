package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
)

// GetTimeStamp
func GetCurrentTimestamp(w http.ResponseWriter, r *http.Request) {

	timestamp := managers.GetTimestamp()

	json.NewEncoder(w).Encode(timestamp)
}

func TestEngine(w http.ResponseWriter, r *http.Request) {
	managers.RunGames()

	json.NewEncoder(w).Encode("Game ran!")
}

func GenerateCollegeTeams(w http.ResponseWriter, r *http.Request) {
	managers.ImportCollegeTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateProTeams(w http.ResponseWriter, r *http.Request) {
	managers.ImportProTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateInitialRosters(w http.ResponseWriter, r *http.Request) {
	managers.GenerateInitialRosters()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCHLRosters(w http.ResponseWriter, r *http.Request) {
	managers.GenerateInitialCHLRosters()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateProTestData(w http.ResponseWriter, r *http.Request) {
	managers.GenerateInitialProPool()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func RunAICollegeLineups(w http.ResponseWriter, r *http.Request) {
	managers.RunLineupsForAICollegeTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func RunAIProLineups(w http.ResponseWriter, r *http.Request) {
	managers.RunLineupsForAIProTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCroots(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCroots()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateWalkons(w http.ResponseWriter, r *http.Request) {
	managers.GenerateWalkonCroots()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func ImportProRosters(w http.ResponseWriter, r *http.Request) {
	managers.ImportProRosters()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func ImportTeamProfileRecords(w http.ResponseWriter, r *http.Request) {
	managers.ImportTeamRecruitingProfiles()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCapsheets(w http.ResponseWriter, r *http.Request) {
	managers.AllocateCapsheets()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func AddFAPreferences(w http.ResponseWriter, r *http.Request) {
	managers.AddFAPreferences()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GeneratePHLSchedule(w http.ResponseWriter, r *http.Request) {
	managers.ImportPHLSeasonSchedule()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCHLSchedule(w http.ResponseWriter, r *http.Request) {
	managers.ImportCHLSchedule()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCHLTourneySchedule(w http.ResponseWriter, r *http.Request) {
	managers.ImportCHLTourneySchedule()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GeneratePreseasonGames(w http.ResponseWriter, r *http.Request) {
	managers.GeneratePreseasonGames()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func ShowGameResults(w http.ResponseWriter, r *http.Request) {
	managers.ShowGames()
	json.NewEncoder(w).Encode("Game results revealed!")
}

func SyncToNextWeek(w http.ResponseWriter, r *http.Request) {
	managers.MoveUpWeek()
	json.NewEncoder(w).Encode("Game results revealed!")
}

func AssignAllRecruitRanks(w http.ResponseWriter, r *http.Request) {
	managers.AssignAllRecruitRanks()
	json.NewEncoder(w).Encode("Game results revealed!")
}

func GenerateCustomCroots(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCustomCroots()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func FillAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.FillAIRecruitingBoards()
	fmt.Println(w, "AI Teams Successfully filled boards.")
}

func SyncAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.ResetAIBoardsForCompletedTeams()
	managers.AllocatePointsToAIBoards()
	fmt.Println(w, "AI teams successfully spent points.")
}

func SyncRecruiting(w http.ResponseWriter, r *http.Request) {
	managers.SyncCollegeRecruiting()
	json.NewEncoder(w).Encode("Recruiting Sync Complete")
}

func UpdateTeamRankings(w http.ResponseWriter, r *http.Request) {
	managers.UpdateTeamRanks()
	json.NewEncoder(w).Encode("Recruiting Sync Complete")
}

func GenerateProPlayoffGames(w http.ResponseWriter, r *http.Request) {
	// db := dbprovider.GetInstance().GetDB()
	// ts := managers.GetTimestamp()
	// managers.GenerateProPlayoffGames(db, ts)
	json.NewEncoder(w).Encode("Pro Playoff Games Generated")
}

func ProgressCollegePlayers(w http.ResponseWriter, r *http.Request) {
	managers.CollegeProgressionMain()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func ProgressProPlayers(w http.ResponseWriter, r *http.Request) {
	managers.ProfessionalProgressionMain()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func FixHistoricCollegePlayerData(w http.ResponseWriter, r *http.Request) {
	managers.FixMigrationPlayerData()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func UpdateCHLTeamValues(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func SeasonMigration(w http.ResponseWriter, r *http.Request) {
	managers.HandlePostSeasonMigration()
	json.NewEncoder(w).Encode("Data Generated ran!")
}
