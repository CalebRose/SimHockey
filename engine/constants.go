package engine

const HomeGoal string = "Home Goal"
const HomeZone string = "Home Zone"
const NeutralZone string = "Neutral Zone"
const AwayZone string = "Away Zone"
const AwayGoal string = "Away Goal"
const Defender string = "D"
const Forward string = "F"
const Center string = "C"
const Goalie string = "G"
const Rebound string = "Rebound"
const Defense string = "Defense"
const ShotBlock string = "ShotBlock"
const Faceoff string = "Faceoff"
const Pass string = "Pass"
const EasyReq float64 = 8
const BaseReq float64 = 10
const DiffReq float64 = 14
const CritSuccess int = 20
const CritFail int = 1
const Heads int = 1
const Tails int = 1
const ModifierFactor float64 = 1.3 // Adjust as needed for your testing
const ScaleFactor float64 = 1.7    // Adjust as needed for your testing
const MinorPenalty string = "Minor Penalty"
const MajorPenalty string = "Major Penalty"
const Misconduct string = "Misconduct"
const GameMisconduct string = "Game Misconduct"
const MatchPenalty string = "Match Penalty"
const BodyCheck string = "BodyCheck"
const StickCheck string = "StickCheck"
const General string = "General"
const Fight string = "Fight"
const ShootoutMomenumModifier float64 = 0.375
const RegularPeriodTime uint16 = 1200
const OvertimePeriodTime uint16 = 300
const MaxTimeOnClock uint16 = 65000

// Event Constants
const FaceoffID uint8 = 1
const PhysDefenseCheckID uint8 = 2
const DexDefenseCheckID uint8 = 3
const PassCheckID uint8 = 4
const AgilityCheckID uint8 = 5
const WristshotCheckID uint8 = 6
const SlapshotCheckID uint8 = 7
const PenaltyCheckID uint8 = 8

// Zone IDs
const HomeGoalZoneID uint8 = 9
const HomeZoneID uint8 = 10
const NeutralZoneID uint8 = 11
const AwayZoneID uint8 = 12
const AwayGoalZoneID uint8 = 13

// Outcome IDs
const DefenseTakesPuckID uint8 = 14
const CarrierKeepsPuckID uint8 = 15
const DefenseStopAgilityID uint8 = 16
const OffenseMovesUpID uint8 = 17
const GeneralPenaltyID uint8 = 18
const FightPenaltyID uint8 = 20
const InterceptedPassID uint8 = 21
const ReceivedPassID uint8 = 22
const HomeFaceoffWinID uint8 = 23
const AwayFaceoffWinID uint8 = 24
const InAccurateShotID uint8 = 25
const ShotBlockedID uint8 = 26
const GoalieSaveID uint8 = 27
const GoalieReboundID uint8 = 28
const ShotOnGoalID uint8 = 29

// PenaltyIDs
