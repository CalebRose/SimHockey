package firebase

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/CalebRose/SimHockey/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MissedRecruitingReports(profiles []structs.RecruitingTeamProfile, ts structs.Timestamp) []AdminReport {
	reports := make([]AdminReport, 0)
	if uint64(ts.WeekID) > math.MaxInt64 || uint64(ts.Phase) > math.MaxInt64 {
		log.Printf("RecruitingSyncManager: cannot report missed recruiting: week ID or phase exceeds Firestore integer range")
		return reports
	}
	for _, profile := range profiles {
		if !profile.IsUserTeam || profile.SpentPoints != 0 {
			continue
		}
		if uint64(profile.TeamID) > math.MaxInt64 {
			log.Printf("RecruitingSyncManager: cannot report missed recruiting for team %d: ID exceeds Firestore integer range", profile.TeamID)
			continue
		}
		subject := profile.Team
		if profile.Recruiter != "" {
			subject = fmt.Sprintf("%s (%s)", profile.Recruiter, profile.Team)
		}
		reports = append(reports, AdminReport{
			League:       "chl",
			TeamID:       int64(profile.TeamID),
			ReportType:   "missed_recruiting",
			Message:      fmt.Sprintf("%s spent no recruiting points in week %d.", subject, ts.Week),
			ResponseType: AdminReportMissedRecruiting,
			Status:       "pending",
			Phase:        int64(ts.Phase),
			WeekID:       int64(ts.WeekID),
			ReportedBy:   "system",
		})
	}
	return reports
}

func IncompleteRosterReport(league, message string, teamID uint, ts structs.Timestamp) (AdminReport, error) {
	if uint64(teamID) > math.MaxInt64 || uint64(ts.WeekID) > math.MaxInt64 || uint64(ts.Phase) > math.MaxInt64 {
		return AdminReport{}, fmt.Errorf("firebase: incomplete roster report identifiers exceed Firestore integer range")
	}
	return AdminReport{
		League:       strings.ToLower(league),
		TeamID:       int64(teamID),
		ReportType:   "incomplete_roster",
		Message:      message,
		ResponseType: AdminReportBrokenGameplan,
		Status:       "pending",
		Phase:        int64(ts.Phase),
		WeekID:       int64(ts.WeekID),
		ReportedBy:   "system",
	}, nil
}

func TradeTransactionReport(teamID uint, message string, ts structs.Timestamp) (AdminReport, error) {
	if uint64(teamID) > math.MaxInt64 || uint64(ts.WeekID) > math.MaxInt64 || uint64(ts.Phase) > math.MaxInt64 {
		return AdminReport{}, fmt.Errorf("firebase: trade report identifiers exceed Firestore integer range")
	}
	return AdminReport{
		League:       "phl",
		TeamID:       int64(teamID),
		ReportType:   "trade_transaction",
		Message:      message,
		ResponseType: AdminReportTradeTransaction,
		Status:       "pending",
		Phase:        int64(ts.Phase),
		WeekID:       int64(ts.WeekID),
		ReportedBy:   "system",
	}, nil
}

// CreateAdminReport creates a report without overwriting an existing one.
func CreateAdminReport(ctx context.Context, id string, report AdminReport) error {
	if id == "" || strings.Contains(id, "/") {
		return fmt.Errorf("firebase: invalid admin report document ID %q", id)
	}
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	return createAdminReport(ctx, client, id, report)
}

func createAdminReport(ctx context.Context, client *firestore.Client, id string, report AdminReport) error {
	now := time.Now().UTC()
	report.CreatedAt = now
	report.UpdatedAt = now

	_, err := client.Collection("adminReports").Doc(id).Create(ctx, report)
	if status.Code(err) == codes.AlreadyExists {
		return nil
	}
	if err != nil {
		return fmt.Errorf("firebase: failed to create admin report %s: %w", id, err)
	}
	return nil
}
