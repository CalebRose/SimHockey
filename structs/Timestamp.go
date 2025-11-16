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
	CHLAttributesUpdated          bool
	SeasonMigrationDone           bool
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
	PreseasonPhase                uint
	IsPreseason                   bool
	SeasonPhase                   uint
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
	if t.IsOffSeason || t.IsDraftTime {
		return
	}
	if t.IsPreseason {
		t.PreseasonPhase++
		if t.PreseasonPhase > 3 {
			t.IsPreseason = !t.IsPreseason
			t.PreseasonPhase = 0
		}
	} else {
		t.WeekID++
		t.Week++
	}
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
	t.IsPreseason = true
}

func (t *Timestamp) MoveUpSeason() {
	t.SeasonID++
	t.Season++
	t.Week = 0
	baseSeason := t.Season - 2000
	multSeason := baseSeason * 100
	t.WeekID = multSeason
	t.Y1Capspace = t.Y2Capspace
	t.Y2Capspace = t.Y3Capspace
	t.Y3Capspace = t.Y4Capspace
	t.Y4Capspace = t.Y5Capspace
	t.Y5Capspace += 5
	t.PreseasonPhase = 1
	t.CrootsGenerated = false
	t.ProgressedCollegePlayers = false
	t.ProgressedProfessionalPlayers = false
	t.CollegeSeasonOver = false
	t.NHLSeasonOver = false
	t.IsDraftTime = true
	t.IsOffSeason = true
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
	t.GamesARan = false
	t.GamesBRan = false
	t.GamesCRan = false
	t.GamesDRan = false
	// t.ToggleRES()
	t.ToggleRecruiting()
	t.ToggleGMActions()

	// Migrate game results ?
}

func (t *Timestamp) ToggleTimeSlot(ts string) {

}

func (t *Timestamp) ToggleGames(matchType string) {
	switch matchType {
	case "A":
		t.GamesARan = true
	case "B":
		t.GamesBRan = true
	case "C":
		t.GamesCRan = true
	case "D":
		t.GamesDRan = true
	}
	if t.IsPreseason {
		t.PreseasonPhase++
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

func (t *Timestamp) GetCurrentGameType(isCollege bool) (int, string) {
	if t.IsPreseason {
		return 1, "1"
	}
	if (t.Week > 17 && isCollege) || (t.Week > 18 && !isCollege) {
		return 3, "3"
	}
	return 2, "2"
}
