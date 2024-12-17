package controllers

import (
	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
)

func GetUpdatedTimestamp() structs.Timestamp {
	ts := managers.GetTimestamp()
	return ts
}
