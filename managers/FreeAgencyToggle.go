package managers

import (
	"github.com/CalebRose/SimHockey/dbprovider"
)

func ToggleFreeAgency() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	ts.IsFreeAgencyLocked = false
    ts.FreeAgencyRound = 1
	db.Save(&ts)
}
