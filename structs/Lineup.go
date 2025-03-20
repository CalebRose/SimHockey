package structs

import "gorm.io/gorm"

type UpdateLineupsDTO struct {
	CHLTeamID         uint
	CHLLineups        []CollegeLineup
	CHLShootoutLineup CollegeShootoutLineup
	CollegePlayers    []CollegePlayer
	PHLTeamID         uint
	PHLLineups        []ProfessionalLineup
	PHLShootoutLineup ProfessionalShootoutLineup
	ProPlayers        []ProfessionalPlayer
}

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

type CollegeShootoutLineup struct {
	ShootoutPlayerIDs
}

type ProfessionalShootoutLineup struct {
	ShootoutPlayerIDs
}

type ShootoutPlayerIDs struct {
	ID               uint
	TeamID           uint
	Shooter1ID       uint
	Shooter1ShotType uint8 // 1 == Close, 2 == Long
	Shooter2ID       uint
	Shooter2ShotType uint8 // 1 == Close, 2 == Long
	Shooter3ID       uint
	Shooter3ShotType uint8 // 1 == Close, 2 == Long
	Shooter4ID       uint
	Shooter4ShotType uint8 // 1 == Close, 2 == Long
	Shooter5ID       uint
	Shooter5ShotType uint8 // 1 == Close, 2 == Long
	Shooter6ID       uint
	Shooter6ShotType uint8 // 1 == Close, 2 == Long
}

func (s *ShootoutPlayerIDs) AssignIDs(s1, s2, s3, s4, s5, s6 uint) {
	s.Shooter1ID = s1
	s.Shooter2ID = s2
	s.Shooter3ID = s3
	s.Shooter4ID = s4
	s.Shooter5ID = s5
	s.Shooter6ID = s6
}

func (s *ShootoutPlayerIDs) AssignShotTypes(st1, st2, st3, st4, st5, st6 uint8) {
	s.Shooter1ShotType = st1
	s.Shooter2ShotType = st2
	s.Shooter3ShotType = st3
	s.Shooter4ShotType = st4
	s.Shooter5ShotType = st5
	s.Shooter6ShotType = st6
}

type Allocations struct {
	// AGZ == Attacking Goal Zone
	AGZShot       uint8
	AGZPass       uint8
	AGZPassBack   uint8
	AGZAgility    uint8
	AGZStickCheck uint8
	AGZBodyCheck  uint8
	// AZ == Attacking Zone
	AZShot       uint8
	AZPass       uint8
	AZLongPass   uint8
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
	DZPassBack   uint8
	DZAgility    uint8
	DZStickCheck uint8
	DZBodyCheck  uint8
	// DGZ == Defending Goal Zone
	DGZPass       uint8
	DGZLongPass   uint8
	DGZAgility    uint8
	DGZStickCheck uint8
	DGZBodyCheck  uint8
}
