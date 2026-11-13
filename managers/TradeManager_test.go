package managers

import (
	"errors"
	"strings"
	"testing"

	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func TestDraftTradeOptions(t *testing.T) {
	player := structs.TradeOption{Model: gorm.Model{ID: 10}, TeamID: 1, PlayerID: 99}
	pick := structs.TradeOption{TeamID: 2, DraftPickID: 55}
	proposal := structs.TradeProposal{
		TeamTradeOptions:          []structs.TradeOption{player},
		RecepientTeamTradeOptions: []structs.TradeOption{player, pick},
	}
	options, err := draftTradeOptions(proposal)
	if err != nil || len(options) != 2 || options[0].PlayerID != 99 || options[1].DraftPickID != 55 {
		t.Fatalf("unexpected merged options: %+v, %v", options, err)
	}
	proposal.RecepientTeamTradeOptions[0].SalaryPercentage = 50
	if _, err := draftTradeOptions(proposal); err == nil {
		t.Fatal("expected conflicting duplicate options to be rejected")
	}
	proposal.TeamTradeOptions = []structs.TradeOption{{TeamID: 1, PlayerID: 99}}
	proposal.RecepientTeamTradeOptions = []structs.TradeOption{player}
	options, err = draftTradeOptions(proposal)
	if err != nil || len(options) != 1 || options[0].ID != 10 {
		t.Fatalf("expected persisted duplicate to be retained: %+v, %v", options, err)
	}
	proposal.TeamTradeOptions = []structs.TradeOption{{PlayerID: 1, DraftPickID: 2}}
	if _, err := draftTradeOptions(proposal); err == nil {
		t.Fatal("an option cannot specify both a player and a pick")
	}
	if err := SyncAcceptedDraftTrade(structs.TradeProposal{TeamID: 1, RecepientTeamID: 2}); !errors.Is(err, ErrInvalidDraftTrade) {
		t.Errorf("expected empty draft trade to be rejected, got %v", err)
	}
}

func TestValidateTradeTeams(t *testing.T) {
	proposal := structs.TradeProposal{TeamID: 1, RecepientTeamID: 2}
	teams := map[uint]structs.ProfessionalTeam{
		1: {Model: gorm.Model{ID: 1}},
		2: {Model: gorm.Model{ID: 2}},
	}
	capsheets := map[uint]structs.ProCapsheet{
		1: {Model: gorm.Model{ID: 11}},
		2: {Model: gorm.Model{ID: 12}},
	}
	if err := validateTradeTeams(proposal, teams, capsheets); err != nil {
		t.Fatal(err)
	}
	delete(capsheets, 2)
	if err := validateTradeTeams(proposal, teams, capsheets); err == nil {
		t.Error("trade with a missing capsheet must not be processed")
	}
}

func TestTradeReportIdentityAndMessage(t *testing.T) {
	teams := map[uint]structs.ProfessionalTeam{
		1: {BaseTeam: structs.BaseTeam{Abbreviation: "AAA"}},
		2: {BaseTeam: structs.BaseTeam{Abbreviation: "BBB"}},
	}
	proposal := structs.TradeProposal{Model: gorm.Model{ID: 17}, TeamID: 1, RecepientTeamID: 2}
	details := []string{"Player One (player 8): AAA to BBB", "2030 round 2 pick (pick 9): BBB to AAA"}
	id, message, err := tradeReportIdentityAndMessage(proposal, teams, details, false)
	if err != nil || id != "trade-proposal-17" {
		t.Fatalf("unexpected saved trade ID %q: %v", id, err)
	}
	for _, part := range []string{"#17", "AAA (team 1)", "BBB (team 2)", details[0], details[1]} {
		if !strings.Contains(message, part) {
			t.Errorf("message %q does not contain %q", message, part)
		}
	}
	id, _, err = tradeReportIdentityAndMessage(proposal, teams, details, true)
	if err != nil || id != "draft-trade-17" {
		t.Fatalf("unexpected draft trade ID %q: %v", id, err)
	}
	proposal.ID = 0
	first, _, err := tradeReportIdentityAndMessage(proposal, teams, details, true)
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := tradeReportIdentityAndMessage(proposal, teams, details, true)
	if err != nil || first == second || !strings.HasPrefix(first, "draft-trade-") {
		t.Fatalf("drafts without IDs need distinct report IDs: %q, %q, %v", first, second, err)
	}
	proposal.RecepientTeamID = 3
	if _, _, err := tradeReportIdentityAndMessage(proposal, teams, details, false); err == nil {
		t.Error("expected missing team to be rejected")
	}
}
