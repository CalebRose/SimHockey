package managers

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func GeneratePreseasonGames() {
	db := dbprovider.GetInstance().GetDB()

	collegeTeamMap := GetCollegeTeamMap()
	proTeamMap := GetProTeamMap()

	collegeGames := GetCollegeGamesForPreseason(collegeTeamMap)
	if collegeGames == nil {
		return
	}
	proGames := GetProGamesForPreseason(proTeamMap)

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := repository.CreateCHLGamesRecordsBatch(tx, collegeGames, 20); err != nil {
			return fmt.Errorf("save college preseason games: %w", err)
		}
		if err := repository.CreatePHLGamesRecordsBatch(tx, proGames, 20); err != nil {
			return fmt.Errorf("save professional preseason games: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("could not save preseason schedules: %v", err)
	}
}

func GetCollegeGamesForPreseason(teamMap map[uint]structs.CollegeTeam) []structs.CollegeGame {
	ts := GetTimestamp()
	existingGames, err := repository.FindCollegeGamesWithError(repository.GamesClauses{
		WeekID:      strconv.Itoa(int(ts.WeekID)),
		SeasonID:    strconv.Itoa(int(ts.SeasonID)),
		IsPreseason: true,
	})
	if err != nil {
		log.Printf("could not load college preseason games: %v", err)
		return nil
	}
	games, err := generateCollegePreseasonSchedule(ts, teamMap, existingGames)
	if err != nil {
		log.Printf("could not generate college preseason schedule: %v", err)
		return nil
	}
	return games
}

func generateCollegePreseasonSchedule(ts structs.Timestamp, teamMap map[uint]structs.CollegeTeam, existingGames []structs.CollegeGame) ([]structs.CollegeGame, error) {
	timeSlots := []string{"A", "B", "C"}
	teamIDs := make([]uint, 0, len(teamMap))
	played := make(map[uint]map[uint]bool)
	occupied := make(map[string]map[uint]bool)
	for _, slot := range timeSlots {
		occupied[slot] = make(map[uint]bool)
	}
	for id := range teamMap {
		teamIDs = append(teamIDs, id)
		played[id] = make(map[uint]bool)
	}
	rand.Shuffle(len(teamIDs), func(i, j int) { teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i] })

	for _, game := range existingGames {
		if slots, ok := occupied[game.GameDay]; ok {
			for _, id := range []uint{game.HomeTeamID, game.AwayTeamID} {
				if id == 0 {
					continue
				}
				if slots[id] {
					return nil, fmt.Errorf("team %d is already scheduled more than once in slot %s", id, game.GameDay)
				}
				slots[id] = true
			}
		}
		if game.HomeTeamID != 0 && game.AwayTeamID != 0 {
			if played[game.HomeTeamID] != nil {
				played[game.HomeTeamID][game.AwayTeamID] = true
			}
			if played[game.AwayTeamID] != nil {
				played[game.AwayTeamID][game.HomeTeamID] = true
			}
		}
	}

	games := make([]structs.CollegeGame, 0)
	const maxSearchNodes = 200000
	searchNodes := 0
	var schedule func() bool
	schedule = func() bool {
		searchNodes++
		if searchNodes > maxSearchNodes {
			return false
		}

		var selectedTeam uint
		var selectedSlot string
		var candidates []uint
		for _, slot := range timeSlots {
			for _, teamID := range teamIDs {
				if occupied[slot][teamID] {
					continue
				}
				var eligible []uint
				for _, opponentID := range teamIDs {
					if opponentID == teamID || occupied[slot][opponentID] || played[teamID][opponentID] {
						continue
					}
					eligible = append(eligible, opponentID)
				}
				if len(eligible) == 0 {
					return false
				}
				if selectedSlot == "" || len(eligible) < len(candidates) {
					selectedTeam, selectedSlot, candidates = teamID, slot, eligible
				}
			}
		}
		if selectedSlot == "" {
			return true
		}

		for _, opponentID := range candidates {
			occupied[selectedSlot][selectedTeam], occupied[selectedSlot][opponentID] = true, true
			played[selectedTeam][opponentID], played[opponentID][selectedTeam] = true, true
			games = append(games, generateCollegeGame(ts.SeasonID, ts.WeekID, ts.Week, selectedTeam, opponentID, selectedSlot, "", teamMap, true))
			if schedule() {
				return true
			}
			games = games[:len(games)-1]
			delete(occupied[selectedSlot], selectedTeam)
			delete(occupied[selectedSlot], opponentID)
			delete(played[selectedTeam], opponentID)
			delete(played[opponentID], selectedTeam)
			if searchNodes > maxSearchNodes {
				break
			}
		}
		return false
	}

	if !schedule() {
		return nil, fmt.Errorf("unable to complete all preseason slots without repeat opponents (searched %d states)", searchNodes)
	}
	return games, nil
}

func GetProGamesForPreseason(teamMap map[uint]structs.ProfessionalTeam) []structs.ProfessionalGame {
	ts := GetTimestamp()
	playedGameReference := make(map[uint]map[uint]bool)
	games := []structs.ProfessionalGame{}
	gameDay := "A"
	teamIDs := make([]uint, 0, len(teamMap))
	for id := range teamMap {
		teamIDs = append(teamIDs, id)
		playedGameReference[id] = make(map[uint]bool)
	}

	for round := 1; round <= 3; round++ {
		var pairings [][2]uint
		const maxTries = 500
		for tries := 0; tries < maxTries; tries++ {
			// shuffle
			rand.Shuffle(len(teamIDs), func(i, j int) {
				teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i]
			})

			// attempt to pair
			var err error
			pairings, err = createProGamePairings(teamIDs, teamMap, playedGameReference)
			if err == nil {
				// success!
				break
			}
		}
		if len(pairings)*2 != len(teamIDs) {
			log.Fatalf("couldn’t find a full pairing on round %d after %d tries", round, maxTries)
		}

		// generate the Game objects
		for _, pair := range pairings {
			game := generateProfessionalGame(
				ts.SeasonID, ts.WeekID,
				ts.Week,
				pair[0], pair[1],
				gameDay, teamMap, true,
			)
			games = append(games, game)
		}

		// rotate gameDay A→B→C
		switch gameDay {
		case "A":
			gameDay = "B"
		case "B":
			gameDay = "C"
		}
	}
	return games
}
