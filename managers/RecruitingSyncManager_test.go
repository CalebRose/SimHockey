package managers

import (
	"strings"
	"testing"

	"github.com/CalebRose/SimHockey/firebase"
	fbsvc "github.com/CalebRose/SimHockey/firebase"
	"github.com/CalebRose/SimHockey/structs"
)

func TestMissedRecruitingReports(t *testing.T) {
	profiles := []structs.RecruitingTeamProfile{
		{TeamID: 11, Team: "Wolves", Recruiter: "coach", IsUserTeam: true},
		{TeamID: 12, Team: "Bears", IsUserTeam: true, TotalCommitments: 5, RecruitClassSize: 5},
		{TeamID: 13, Team: "Lions", IsUserTeam: true, SpentPoints: 1},
		{TeamID: 14, Team: "Tigers", IsUserTeam: false},
	}
	ts := structs.Timestamp{Week: 3, WeekID: 203, Phase: 7}
	reports := firebase.MissedRecruitingReports(profiles, ts)
	if len(reports) != 2 {
		t.Fatalf("got %d reports, want 2", len(reports))
	}
	report := reports[0]
	if report.TeamID != 11 || report.League != "chl" || report.WeekID != 203 || report.Phase != 7 {
		t.Errorf("incorrect report context: %+v", report)
	}
	if report.ResponseType != fbsvc.AdminReportMissedRecruiting || report.ReportType != "missed_recruiting" || report.Status != "pending" || report.ReportedBy != "system" {
		t.Errorf("incorrect report classification: %+v", report)
	}
	if !strings.Contains(report.Message, "coach") || !strings.Contains(report.Message, "Wolves") || !strings.Contains(report.Message, "week 3") {
		t.Errorf("report message lacks context: %q", report.Message)
	}
	if reports[1].TeamID != 12 || !strings.Contains(reports[1].Message, "Bears") {
		t.Errorf("full recruiting class should still be reported: %+v", reports[1])
	}
}
