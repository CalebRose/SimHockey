package firebase

import (
	"context"
	"math"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/CalebRose/SimHockey/structs"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestIncompleteRosterReport(t *testing.T) {
	ts := structs.Timestamp{WeekID: 102, Phase: 4}
	for _, league := range []string{"CHL", "PHL"} {
		report, err := IncompleteRosterReport(league, "Missing: 1 goalie", 42, ts)
		if err != nil {
			t.Fatal(err)
		}
		if report.League != strings.ToLower(league) || report.TeamID != 42 ||
			report.ReportType != "incomplete_roster" || report.ResponseType != AdminReportBrokenGameplan ||
			report.Status != "pending" || report.WeekID != 102 || report.Phase != 4 ||
			report.ReportedBy != "system" || report.Message != "Missing: 1 goalie" {
			t.Errorf("unexpected report: %+v", report)
		}
	}
	if strconv.IntSize == 64 {
		tooLarge := uint64(math.MaxInt64) + 1
		if _, err := IncompleteRosterReport("CHL", "missing players", uint(tooLarge), ts); err == nil {
			t.Error("expected an error for an ID outside Firestore's integer range")
		}
	}
}

func TestTradeTransactionReport(t *testing.T) {
	report, err := TradeTransactionReport(42, "Trade processed", structs.Timestamp{WeekID: 102, Phase: 4})
	if err != nil {
		t.Fatal(err)
	}
	if report.League != "phl" || report.TeamID != 42 || report.ReportType != "trade_transaction" ||
		report.ResponseType != AdminReportTradeTransaction || report.Status != "pending" ||
		report.WeekID != 102 || report.Phase != 4 || report.ReportedBy != "system" ||
		report.Message != "Trade processed" {
		t.Errorf("unexpected trade report: %+v", report)
	}
	if strconv.IntSize == 64 {
		tooLarge := uint64(math.MaxInt64) + 1
		if _, err := TradeTransactionReport(uint(tooLarge), "trade", structs.Timestamp{}); err == nil {
			t.Error("expected an error for an ID outside Firestore's integer range")
		}
	}
}

type reportFirestoreServer struct {
	firestorepb.UnimplementedFirestoreServer
	requests chan *firestorepb.CommitRequest
}

func (s *reportFirestoreServer) Commit(_ context.Context, req *firestorepb.CommitRequest) (*firestorepb.CommitResponse, error) {
	s.requests <- req
	name := req.Writes[0].GetUpdate().GetName()
	if strings.HasSuffix(name, "/exists") {
		return nil, status.Error(codes.AlreadyExists, "already exists")
	}
	if strings.HasSuffix(name, "/failure") {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	return &firestorepb.CommitResponse{
		WriteResults: []*firestorepb.WriteResult{{UpdateTime: timestamppb.Now()}},
	}, nil
}

func TestCreateAdminReport(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	fake := &reportFirestoreServer{requests: make(chan *firestorepb.CommitRequest, 3)}
	firestorepb.RegisterFirestoreServer(server, fake)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	client, err := firestore.NewClient(ctx, "test-project", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })

	report := AdminReport{
		League: "chl", TeamID: 42, ReportType: "missed_recruiting",
		Message: "No points spent", ResponseType: AdminReportMissedRecruiting,
		Status: "pending", Phase: 4, WeekID: 102, ReportedBy: "system",
	}
	if err := createAdminReport(ctx, client, "new", report); err != nil {
		t.Fatal(err)
	}
	write := (<-fake.requests).Writes[0]
	doc := write.GetUpdate()
	if !strings.HasSuffix(doc.GetName(), "/adminReports/new") {
		t.Errorf("unexpected document name: %s", doc.GetName())
	}
	if write.GetCurrentDocument() == nil || write.GetCurrentDocument().GetExists() {
		t.Error("report must be created with an exists=false precondition")
	}
	fields := doc.GetFields()
	if fields["teamId"].GetIntegerValue() != 42 || fields["responseType"].GetIntegerValue() != 1 ||
		fields["weekId"].GetIntegerValue() != 102 || fields["phase"].GetIntegerValue() != 4 ||
		fields["status"].GetStringValue() != "pending" {
		t.Errorf("unexpected report fields: %v", fields)
	}
	if fields["createdAt"].GetTimestampValue() == nil || fields["updatedAt"].GetTimestampValue() == nil {
		t.Error("report timestamps are missing")
	}
	if _, ok := fields["resolvedAt"]; ok {
		t.Error("unresolved report should not have a resolvedAt timestamp")
	}

	if err := createAdminReport(ctx, client, "exists", report); err != nil {
		t.Errorf("existing report should not be overwritten: %v", err)
	}
	<-fake.requests
	if err := createAdminReport(ctx, client, "failure", report); status.Code(err) != codes.PermissionDenied {
		t.Errorf("expected permission error, got %v", err)
	}
	<-fake.requests

	rosterReport, err := IncompleteRosterReport("PHL", "Missing: 1 goalie", 42, structs.Timestamp{WeekID: 102, Phase: 4})
	if err != nil {
		t.Fatal(err)
	}
	if err := createAdminReport(ctx, client, "incomplete-roster-phl-42-102-4", rosterReport); err != nil {
		t.Fatal(err)
	}
	rosterDoc := (<-fake.requests).Writes[0].GetUpdate()
	if !strings.HasSuffix(rosterDoc.GetName(), "/adminReports/incomplete-roster-phl-42-102-4") ||
		rosterDoc.Fields["league"].GetStringValue() != "phl" ||
		rosterDoc.Fields["reportType"].GetStringValue() != "incomplete_roster" ||
		rosterDoc.Fields["responseType"].GetIntegerValue() != 2 ||
		rosterDoc.Fields["message"].GetStringValue() != "Missing: 1 goalie" {
		t.Errorf("unexpected roster report document: %v", rosterDoc)
	}

	tradeReport, err := TradeTransactionReport(42, "Trade processed", structs.Timestamp{WeekID: 102, Phase: 4})
	if err != nil {
		t.Fatal(err)
	}
	if err := createAdminReport(ctx, client, "trade-proposal-17", tradeReport); err != nil {
		t.Fatal(err)
	}
	tradeDoc := (<-fake.requests).Writes[0].GetUpdate()
	if !strings.HasSuffix(tradeDoc.GetName(), "/adminReports/trade-proposal-17") ||
		tradeDoc.Fields["reportType"].GetStringValue() != "trade_transaction" ||
		tradeDoc.Fields["responseType"].GetIntegerValue() != 3 ||
		tradeDoc.Fields["teamId"].GetIntegerValue() != 42 {
		t.Errorf("unexpected trade report document: %v", tradeDoc)
	}
}

func TestCreateAdminReportRejectsInvalidID(t *testing.T) {
	for _, id := range []string{"", "folder/document"} {
		if err := CreateAdminReport(context.Background(), id, AdminReport{}); err == nil {
			t.Errorf("expected an error for document ID %q", id)
		}
	}
}
