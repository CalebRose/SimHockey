package structs

import "gorm.io/gorm"

type BaseGameplan struct {
	TeamID                  uint
	IsAI                    bool  // True == AI active gameplan, False == Not active
	ForwardShotPreference   uint8 // 1 == Close, 2 == Balanced, 3 == Long Shot
	DefenderShotPreference  uint8 // 1 == Close, 2 == Balanced, 3 == Long Shot
	ForwardCheckPreference  uint8 // 1 == Body, 2 == Balanced, 3 == Stick
	DefenderCheckPreference uint8 // 1 == Body, 2 == Balanced, 3 == Stick
	CenterSortPreference1   uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	CenterSortPreference2   uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	CenterSortPreference3   uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	ForwardSortPreference1  uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	ForwardSortPreference2  uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	ForwardSortPreference3  uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	DefenderSortPreference1 uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	DefenderSortPreference2 uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	DefenderSortPreference3 uint8 // 1 == Overall, 2 == Close Shot Accuracy, 3 == Long Shot Accuracy, 4 == Agility, 5 == Puck Handling, 6 == Strength, 7 == BodyCheck, 8 == StickCheck, 9 == Faceoff
	GoalieSortPreference    uint8 // 1 == Overall, 2 == Goalkeeping, 3 == Goalievision
	LongerPassesEnabled     bool
}

type CollegeGameplan struct {
	gorm.Model
	BaseGameplan
}

type ProGameplan struct {
	gorm.Model
	BaseGameplan
}

func (bg *BaseGameplan) UpdateGameplan(updatedGameplan BaseGameplan) {
	bg = &updatedGameplan
}

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
	AGZShot       int8
	AGZPass       int8
	AGZPassBack   int8
	AGZAgility    int8
	AGZStickCheck int8
	AGZBodyCheck  int8
	// AZ == Attacking Zone
	AZShot       int8
	AZPass       int8
	AZLongPass   int8
	AZAgility    int8
	AZStickCheck int8
	AZBodyCheck  int8
	// N == Neutral
	NPass       int8
	NAgility    int8
	NStickCheck int8
	NBodyCheck  int8
	// DZ Defending Zone
	DZPass       int8
	DZPassBack   int8
	DZAgility    int8
	DZStickCheck int8
	DZBodyCheck  int8
	// DGZ == Defending Goal Zone
	DGZPass       int8
	DGZLongPass   int8
	DGZAgility    int8
	DGZStickCheck int8
	DGZBodyCheck  int8
}
