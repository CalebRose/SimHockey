package managers

import (
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

// GetTimestamp -- Get the Timestamp
func GetTimestamp() structs.Timestamp {
	return repository.FindTimestamp()
}
