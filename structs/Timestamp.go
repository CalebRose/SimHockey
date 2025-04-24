package structs

import "gorm.io/gorm"

type Timestamp struct {
	gorm.Model
	RunCron                       bool
	RunGames                      bool
	WeekID                        uint
	Week                          uint
	SeasonID                      uint
	Season                        uint
	GamesARan                     bool
	GamesBRan                     bool
	GamesCRan                     bool
	GamesDRan                     bool
	CollegePollRan                bool
	RecruitingSynced              bool
	GMActionsCompleted            bool
	IsOffSeason                   bool
	IsRecruitingLocked            bool
	AIDepthchartsSync             bool
	AIRecruitingBoardsSynced      bool
	CollegeSeasonOver             bool
	NHLSeasonOver                 bool
	CrootsGenerated               bool
	ProgressedCollegePlayers      bool
	ProgressedProfessionalPlayers bool
	TransferPortalPhase           uint
	TransferPortalRound           uint
	IsFreeAgencyLocked            bool
	FreeAgencyRound               uint
	IsDraftTime                   bool
	Y1Capspace                    float64
	Y2Capspace                    float64
	Y3Capspace                    float64
	Y4Capspace                    float64
	Y5Capspace                    float64
	DeadCapLimit                  float64
}

func (t *Timestamp) GetGameDay() string {
	if !t.GamesARan {
		return "A"
	}
	if !t.GamesBRan {
		return "B"
	}
	if !t.GamesCRan {
		return "C"
	}
	return "D"
}

func (t *Timestamp) MoveUpWeek() {
	t.WeekID++
	t.Week++
}

func (t *Timestamp) MoveUpFreeAgencyRound() {
	t.FreeAgencyRound++
	if t.FreeAgencyRound > 10 {
		t.FreeAgencyRound = 0
		t.IsFreeAgencyLocked = true
		t.IsDraftTime = true
	}
}

func (t *Timestamp) DraftIsOver() {
	t.IsDraftTime = false
	t.IsOffSeason = false
	t.IsFreeAgencyLocked = false
}

func (t *Timestamp) MoveUpSeason() {
	t.SeasonID++
	t.Season++
	t.Week = 0
	t.Y1Capspace = t.Y2Capspace
	t.Y2Capspace = t.Y3Capspace
	t.Y3Capspace = t.Y4Capspace
	t.Y4Capspace = t.Y5Capspace
	t.Y5Capspace += 5
}

func (t *Timestamp) ToggleRecruiting() {
	t.RecruitingSynced = false
	t.IsRecruitingLocked = false
}

func (t *Timestamp) ToggleGMActions() {
	t.GMActionsCompleted = !t.GMActionsCompleted
}

func (t *Timestamp) ToggleLockRecruiting() {
	t.IsRecruitingLocked = !t.IsRecruitingLocked
}

func (t *Timestamp) ToggleFALock() {
	t.IsFreeAgencyLocked = !t.IsFreeAgencyLocked
}

func (t *Timestamp) SyncToNextWeek() {
	t.MoveUpWeek()
	// Reset Toggles
	t.AIDepthchartsSync = false
	t.AIRecruitingBoardsSynced = false
	t.RunGames = false
	// t.ToggleRES()
	t.ToggleRecruiting()
	t.ToggleGMActions()

	// Migrate game results ?
}

func (t *Timestamp) ToggleTimeSlot(ts string) {

}

func (t *Timestamp) ToggleGames(matchType string) {
	if matchType == "A" {
		t.GamesARan = true
	} else if matchType == "B" {
		t.GamesBRan = true
	} else if matchType == "C" {
		t.GamesCRan = true
	} else if matchType == "D" {
		t.GamesDRan = true
	}
}

func (t *Timestamp) ToggleRunGames() {
	t.RunGames = !t.RunGames
}

func (t *Timestamp) ToggleAIrecruitingSync() {
	t.AIRecruitingBoardsSynced = !t.AIRecruitingBoardsSynced
}

func (t *Timestamp) ToggleAIDepthCharts() {
	t.AIDepthchartsSync = !t.AIDepthchartsSync
}

func (t *Timestamp) ToggleDraftTime() {
	t.IsDraftTime = !t.IsDraftTime
	// t.IsNBAOffseason = false
}

func (t *Timestamp) TogglePollRan() {
	t.CollegePollRan = !t.CollegePollRan
}

func (t *Timestamp) EndTheCollegeSeason() {
	t.IsOffSeason = true
	t.TransferPortalPhase = 1
	t.CollegeSeasonOver = true
}

func (t *Timestamp) ClosePortal() {
	t.TransferPortalPhase = 0
}

func (t *Timestamp) EnactPromisePhase() {
	t.TransferPortalPhase = 2
}

func (t *Timestamp) EnactPortalPhase() {
	t.TransferPortalPhase = 3
}

func (t *Timestamp) IncrementTransferPortalRound() {
	t.IsRecruitingLocked = false
	if t.TransferPortalRound < 10 {
		t.TransferPortalRound += 1
	}
}

func (t *Timestamp) EndTheProfessionalSeason() {
	t.IsOffSeason = true
	t.FreeAgencyRound = 1
	t.IsDraftTime = false
	t.IsFreeAgencyLocked = true
	t.NHLSeasonOver = true
}

func (t *Timestamp) ToggleGeneratedCroots() {
	t.CrootsGenerated = !t.CrootsGenerated
}

func (t *Timestamp) ToggleCollegeProgression() {
	t.ProgressedCollegePlayers = !t.ProgressedCollegePlayers
}

func (t *Timestamp) ToggleProfessionalProgression() {
	t.ProgressedProfessionalPlayers = !t.ProgressedProfessionalPlayers
	t.IsFreeAgencyLocked = false
	t.IsDraftTime = true
}

// func (t *Timestamp) GetNHLCurrentGameType() (int, string) {
// 	if t.NHLPreseason {
// 		return 1, "1"
// 	}
// 	if t.NHLWeek > 18 {
// 		return 3, "3"
// 	}
// 	return 2, "2"
// }

// func (t *Timestamp) GetCollegeCurrentGameType() (int, string) {
// 	if t.CFBSpringGames {
// 		return 1, "1"
// 	}
// 	if t.CollegeWeek > 14 {
// 		return 3, "3"
// 	}
// 	return 2, "2"
// }
