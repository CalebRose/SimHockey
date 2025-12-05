package repository

import (
	"errors"
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

// Base query for free agents
func FindAllFreeAgents(isFreeAgent, isWaiverWire, includeOffers, includeSeasonStats bool) []structs.ProfessionalPlayer {
	db := dbprovider.GetInstance().GetDB()

	var proPlayers []structs.ProfessionalPlayer

	query := db.Model(&proPlayers)

	if includeOffers {
		query = query.Preload("Offers", func(db *gorm.DB) *gorm.DB {
			return db.Order("contract_value DESC").Where("is_active = true")
		})
	}

	if includeSeasonStats {
		ts := FindTimestamp()
		seasonID := ts.SeasonID
		if ts.IsOffSeason && ts.SeasonID > 1 {
			seasonID = ts.SeasonID - 1
		}
		query = query.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ?", strconv.Itoa(int(seasonID)))
		})
	}

	if isFreeAgent {
		query = query.Where("is_free_agent = ?", isFreeAgent)
	}

	if isWaiverWire {
		query = query.Preload("Contract", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = true")
		}).Where("is_waived = ?", isWaiverWire)
	}

	if err := query.Order("overall desc").Where("team_id = ?", "0").Find(&proPlayers).Error; err != nil {
		return []structs.ProfessionalPlayer{}
	}

	return proPlayers
}

func FindAffiliatePlayers(teamID, affiliateTeamID string, includeOffers, includeSeasonStats bool) []structs.ProfessionalPlayer {
	db := dbprovider.GetInstance().GetDB()

	var proPlayers []structs.ProfessionalPlayer

	query := db.Model(&proPlayers)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if len(affiliateTeamID) > 0 {
		query = query.Where("affiliate_team_id = ?", affiliateTeamID)
	}

	if includeOffers {
		query = query.Preload("Offers", func(db *gorm.DB) *gorm.DB {
			return db.Order("contract_value DESC").Where("is_active = true")
		})
	}

	if includeSeasonStats {
		ts := FindTimestamp()
		seasonID := ts.SeasonID
		if ts.IsOffSeason && ts.SeasonID > 1 {
			seasonID = ts.SeasonID - 1
		}
		query = query.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ?", strconv.Itoa(int(seasonID)))
		})
	}

	if err := query.Order("overall desc").Where("is_affiliate_player = ?", true).Find(&proPlayers).Error; err != nil {
		return []structs.ProfessionalPlayer{}
	}

	return proPlayers
}

func FindFreeAgentOfferRecord(teamID, playerID, offerID string, onlyActive bool) structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := structs.FreeAgencyOffer{}

	query := db.Model(&offers)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if len(playerID) > 0 {
		query = query.Where("player_id = ?", playerID)
	}

	if len(offerID) > 0 {
		query = query.Where("id = ?", offerID)
	}

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return structs.FreeAgencyOffer{}
	}

	return offers
}

func FindWaiverWireOfferRecord(teamID, playerID, offerID string, onlyActive bool) structs.WaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := structs.WaiverOffer{}

	query := db.Model(&offers)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if len(playerID) > 0 {
		query = query.Where("player_id = ?", playerID)
	}

	if len(offerID) > 0 {
		query = query.Where("id = ?", offerID)
	}

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return structs.WaiverOffer{}
	}

	return offers
}

func FindAllFreeAgentOffers(teamID, playerID, offerID string, onlyActive bool) []structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.FreeAgencyOffer{}

	query := db.Model(&offers)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if len(playerID) > 0 {
		query = query.Where("player_id = ?", playerID)
	}

	if len(offerID) > 0 {
		query = query.Where("id = ?", offerID)
	}

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.FreeAgencyOffer{}
	}

	return offers
}

func FindAllWaiverWireOffers(teamID, playerID, offerID string, onlyActive bool) []structs.WaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.WaiverOffer{}

	query := db.Model(&offers)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if len(playerID) > 0 {
		query = query.Where("player_id = ?", playerID)
	}

	if len(offerID) > 0 {
		query = query.Where("id = ?", offerID)
	}

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.WaiverOffer{}
	}

	return offers
}

func FindAllProContracts(onlyActive bool) []structs.ProContract {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.ProContract{}

	query := db.Model(&offers)

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.ProContract{}
	}

	return offers
}

func FindAllProExtensions(onlyActive bool) []structs.ExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.ExtensionOffer{}

	query := db.Model(&offers)

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.ExtensionOffer{}
	}

	return offers
}

func CreateProContractRecord(db *gorm.DB, contract structs.ProContract) error {
	if err := db.Create(&contract).Error; err != nil {
		return err
	}

	return nil
}

func CreateProContractRecordsBatch(db *gorm.DB, contracts []structs.ProContract, batchSize int) error {
	total := len(contracts)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)
		if err := db.CreateInBatches(contracts[i:end], batchSize).Error; err != nil {
			return err
		}
	}

	return nil
}

func FindLatestFreeAgentOfferID(db *gorm.DB) uint {
	var latestOffer structs.FreeAgencyOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func FindLatestWaiverOfferID(db *gorm.DB) uint {
	var latestOffer structs.WaiverOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func SaveFreeAgentOfferRecord(offerRecord structs.FreeAgencyOffer, db *gorm.DB) {
	err := db.Save(&offerRecord).Error
	if err != nil {
		log.Panicln("Could not save waiver " + strconv.Itoa(int(offerRecord.ID)))
	}
}

func SaveWaiverRecord(waiverRecord structs.WaiverOffer, db *gorm.DB) {
	err := db.Save(&waiverRecord).Error
	if err != nil {
		log.Panicln("Could not save waiver " + strconv.Itoa(int(waiverRecord.ID)))
	}
}

func DeleteWaiverRecord(waiverRecord structs.WaiverOffer, db *gorm.DB) {
	err := db.Delete(&waiverRecord).Error
	if err != nil {
		log.Panicln("Could not delete waiver " + strconv.Itoa(int(waiverRecord.ID)))
	}
}

func DeleteExtensionRecord(extensionRecord structs.ExtensionOffer, db *gorm.DB) {
	err := db.Delete(&extensionRecord).Error
	if err != nil {
		log.Panicln("Could not delete extension " + strconv.Itoa(int(extensionRecord.ID)))
	}
}

func CreateExtensionRecord(extensionRecord structs.ExtensionOffer, db *gorm.DB) {
	err := db.Create(&extensionRecord).Error
	if err != nil {
		log.Panicln("Could not save extension " + strconv.Itoa(int(extensionRecord.ID)))
	}
}

func SaveExtensionRecord(extensionRecord structs.ExtensionOffer, db *gorm.DB) {
	err := db.Save(&extensionRecord).Error
	if err != nil {
		log.Panicln("Could not save extension " + strconv.Itoa(int(extensionRecord.ID)))
	}
}
