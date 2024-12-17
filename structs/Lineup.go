package structs

import "gorm.io/gorm"

// Lineup -- Database records for handling each line's strategy for a team. Four records for forwards, three for defenders, two for goalies. Total, nine per team.
type BaseLineup struct {
	TeamID   uint
	Line     uint8
	LineType uint8 // 1== Forward, 2== Defender, 3== Goalie
	LineupPlayerIDs
	Allocations
}

type CollegeLineup struct {
	gorm.Model
	BaseLineup
}

func (c *BaseLineup) MapIDsAndAllocations(ids LineupPlayerIDs, allo Allocations) {
	c.LineupPlayerIDs = ids
	c.Allocations = allo
}

type ProfessionalLineup struct {
	gorm.Model
	BaseLineup
}

type LineupPlayerIDs struct {
	CenterID    uint // Any of the below player IDs will be zero based on LineType
	Forward1ID  uint
	Forward2ID  uint
	Defender1ID uint
	Defender2ID uint
	GoalieID    uint
}

type Allocations struct {
	// AGZ == Attacking Goal Zone
	AGZShot       uint8
	AGZPass       uint8
	AGZStickCheck uint8
	AGZBodyCheck  uint8
	// AZ == Attacking Zone
	AZSlapshot   uint8
	AZWristshot  uint8
	AZPass       uint8
	AZAgility    uint8
	AZStickCheck uint8
	AZBodyCheck  uint8
	// N == Neutral
	NPass       uint8
	NAgility    uint8
	NStickCheck uint8
	NBodyCheck  uint8
	// DZ Defending Zone
	DZPass       uint8
	DZAgility    uint8
	DZStickCheck uint8
	DZBodyCheck  uint8
	// DGZ == Defending Goal Zone
	DGZPass       uint8
	DGZAgility    uint8
	DGZStickCheck uint8
	DGZBodyCheck  uint8
}
