package controllers

import (
	"fmt"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func CutCHLPlayerFromRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["playerID"]

	managers.CutCHLPlayer(PlayerID)

	fmt.Println(w, "NFL Player Cut from Roster")
}

func RedshirtCHLPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["playerID"]

	managers.RedshirtCHLPlayer(PlayerID)

	fmt.Println(w, "NFL Player Cut from Roster")
}

func PromiseCHLPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["playerID"]

	fmt.Println(w, "Implement promises for "+PlayerID+"!")
}

func CutPHLPlayerFromRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["playerID"]

	managers.CutProPlayer(PlayerID)

	fmt.Println(w, "NFL Player Cut from Roster")
}

func SendPHLPlayerToAffiliate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["playerID"]

	managers.SendPHLPlayerToAffiliate(PlayerID)
}

func SendPHLPlayerToTradeBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["playerID"]

	managers.SendPHLPlayerToTradeBlock(PlayerID)
}

func ExtendPHLPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FINISH THIS")
}
