package managers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

// College Functions
func GetAllCollegePlayers() []structs.CollegePlayer {
	return repository.FindAllCollegePlayers(repository.PlayerQuery{})
}

func GetAllHistoricCollegePlayers() []structs.HistoricCollegePlayer {
	return repository.FindAllHistoricCollegePlayers()
}

func GetCollegePlayersByTeamID(TeamID string) []structs.CollegePlayer {
	return repository.FindCollegePlayersByTeamID(TeamID)
}

func GetCollegePlayersMap() map[uint]structs.CollegePlayer {
	players := repository.FindAllCollegePlayers(repository.PlayerQuery{})
	return MakeCollegePlayerMap(players)
}

func GetCollegePlayerMapByTeamID(TeamID string) map[uint]structs.CollegePlayer {
	players := repository.FindCollegePlayersByTeamID(TeamID)
	return MakeCollegePlayerMap(players)
}

func GetCollegePlayerByID(id string) structs.CollegePlayer {
	return repository.FindCollegePlayer(repository.PlayerQuery{ID: id})
}

func GetAllCollegePlayersMapByTeam() map[uint][]structs.CollegePlayer {
	players := repository.FindAllCollegePlayers(repository.PlayerQuery{})
	return MakeCollegePlayerMapByTeamID(players)
}

// Professional Functions
func GetAllProPlayers() []structs.ProfessionalPlayer {
	return repository.FindAllProPlayers(repository.PlayerQuery{})
}

func GetAllRetiredPlayers() []structs.RetiredPlayer {
	return repository.FindAllHistoricProPlayers()
}

func GetProPlayersByTeamID(TeamID string) []structs.ProfessionalPlayer {
	return repository.FindAllProPlayers(repository.PlayerQuery{TeamID: TeamID})
}

func GetProPlayersMap() map[uint]structs.ProfessionalPlayer {
	players := repository.FindAllProPlayers(repository.PlayerQuery{})
	return MakeProfessionalPlayerMap(players)
}

func GetProPlayerMapByTeamID(TeamID string) map[uint]structs.ProfessionalPlayer {
	players := repository.FindAllProPlayers(repository.PlayerQuery{TeamID: TeamID})
	return MakeProfessionalPlayerMap(players)
}

func GetAllProPlayersMapByTeam() map[uint][]structs.ProfessionalPlayer {
	players := repository.FindAllProPlayers(repository.PlayerQuery{})
	return MakeProfessionalPlayerMapByTeamID(players)
}

func RecoverPlayers() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()

	for _, p := range collegePlayers {
		if !p.IsInjured {
			continue
		}
		p.RecoveryCheck()
		repository.SaveCollegeHockeyPlayerRecord(p, db)
	}

	for _, p := range collegePlayers {
		if p.Position != Goalie || p.GoalieStamina == 100 {
			continue
		}
		p.RecoverGoalieStamina()
		repository.SaveCollegeHockeyPlayerRecord(p, db)
	}

	proPlayers := GetAllProPlayers()

	for _, p := range proPlayers {
		if !p.IsInjured {
			continue
		}

		p.RecoveryCheck()
		repository.SaveProPlayerRecord(p, db)
	}

	for _, p := range proPlayers {
		if p.Position != Goalie || p.GoalieStamina == 100 {
			continue
		}
		p.RecoverGoalieStamina()
		repository.SaveProPlayerRecord(p, db)
	}

}

func GetExtensionOfferByPlayerID(playerID string) structs.ExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.ExtensionOffer{}

	err := db.Where("player_id = ?", playerID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func GetLatestExtensionOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.ExtensionOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func CreateExtensionOffer(offer structs.ExtensionOffer) structs.ExtensionOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	extensionOffer := GetExtensionOfferByPlayerID(strconv.Itoa(int(offer.PlayerID)))
	player := repository.FindProPlayer(strconv.Itoa(int(offer.PlayerID)))
	team := repository.FindProTeamRecord(strconv.Itoa(int(offer.TeamID)))
	extensionOffer.CalculateOffer(offer)

	// If the owning team is sending an offer to a player
	if extensionOffer.ID == 0 {
		id := GetLatestExtensionOfferInDB(db)
		extensionOffer.AssignID(id)
		repository.CreateExtensionRecord(extensionOffer, db)
		fmt.Println("Creating Extension Offer!")
		message := team.TeamName + " have offered a " + strconv.Itoa(offer.ContractLength) + " year contract extension for " + player.Position + " " + player.FirstName + " " + player.LastName + "."
		CreateNewsLog("PHL", message, "Free Agency", int(player.TeamID), ts, true)
	} else {
		fmt.Println("Updating Extension Offer!")
		repository.SaveExtensionRecord(extensionOffer, db)
	}

	return extensionOffer
}

func CancelExtensionOffer(offer structs.ExtensionOffer) structs.ExtensionOffer {
	db := dbprovider.GetInstance().GetDB()
	playerID := strconv.Itoa(int(offer.PlayerID))
	freeAgentOffer := GetExtensionOfferByPlayerID(playerID)
	freeAgentOffer.CancelOffer()
	repository.SaveExtensionRecord(freeAgentOffer, db)
	return freeAgentOffer
}
