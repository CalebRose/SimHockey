package managers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	fbsvc "github.com/CalebRose/SimHockey/firebase"
	"github.com/CalebRose/SimHockey/structs"
)

// ─────────────────────────────────────────────
// Forum IDs
// ─────────────────────────────────────────────

// PostGameForumID is the Firestore document ID of the forum category used for
// post-game discussion threads.
const PostGameForumID = "postgame-discussions"

// ─────────────────────────────────────────────
// Post-game discussion threads
// ─────────────────────────────────────────────

// CreatePostGameDiscussionThreadForCHLGame creates a system-generated postgame
// discussion thread for a completed CHL (college hockey) game.
// Intended to be called as a goroutine after the game is saved.
// The operation is idempotent: calling it twice for the same game has no effect.
func CreatePostGameDiscussionThreadForCHLGame(
	game structs.CollegeGame,
	starOne, starTwo, starThree string,
	seasonID uint,
	homeTeamStats structs.CollegeTeamGameStats,
	awayTeamStats structs.CollegeTeamGameStats,
	homePlayerStats []structs.CollegePlayerGameStats,
	awayPlayerStats []structs.CollegePlayerGameStats,
	collegePlayersMap map[uint]structs.CollegePlayer,
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:chl:season%d:game%s", seasonID, gameID)

	title := buildHockeyPostGameThreadTitle(game.BaseGame)
	nodes := buildCHLPostGameNodes(game, homeTeamStats, awayTeamStats, starOne, starTwo, starThree, homePlayerStats, awayPlayerStats, collegePlayersMap)
	bodyText := nodesToPlainText(nodes)
	richBody := buildRichDoc(nodes)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           PostGameForumID + "-simchl",
		ForumPath:         []string{PostGameForumID, "simchl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeGameReference,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedGameID:  gameID,
		ReferencedLeague:  "chl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create CHL postgame thread for game %s: %v", gameID, err)
		return
	}

	log.Printf("ForumManager: created CHL postgame thread %s for game %s (%s)", thread.ID, gameID, title)
}

// CreatePostGameDiscussionThreadForPHLGame creates a system-generated postgame
// discussion thread for a completed PHL (professional hockey) game.
// Intended to be called as a goroutine after the game is saved.
// The operation is idempotent: calling it twice for the same game has no effect.
func CreatePostGameDiscussionThreadForPHLGame(
	game structs.ProfessionalGame,
	starOne, starTwo, starThree string,
	seasonID uint,
	homeTeamStats structs.ProfessionalTeamGameStats,
	awayTeamStats structs.ProfessionalTeamGameStats,
	homePlayerStats []structs.ProfessionalPlayerGameStats,
	awayPlayerStats []structs.ProfessionalPlayerGameStats,
	proPlayersMap map[uint]structs.ProfessionalPlayer,
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:phl:season%d:game%s", seasonID, gameID)

	title := buildHockeyPostGameThreadTitle(game.BaseGame)
	nodes := buildPHLPostGameNodes(game, homeTeamStats, awayTeamStats, starOne, starTwo, starThree, homePlayerStats, awayPlayerStats, proPlayersMap)
	bodyText := nodesToPlainText(nodes)
	richBody := buildRichDoc(nodes)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           PostGameForumID + "-simphl",
		ForumPath:         []string{PostGameForumID, "simphl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeGameReference,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedGameID:  gameID,
		ReferencedLeague:  "phl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create PHL postgame thread for game %s: %v", gameID, err)
		return
	}

	log.Printf("ForumManager: created PHL postgame thread %s for game %s (%s)", thread.ID, gameID, title)
}

// ─────────────────────────────────────────────
// Post-game body builders
// ─────────────────────────────────────────────

func buildCHLPostGameNodes(
	game structs.CollegeGame,
	homeTeamStats structs.CollegeTeamGameStats,
	awayTeamStats structs.CollegeTeamGameStats,
	starOne, starTwo, starThree string,
	homePlayerStats []structs.CollegePlayerGameStats,
	awayPlayerStats []structs.CollegePlayerGameStats,
	collegePlayersMap map[uint]structs.CollegePlayer,
) []map[string]interface{} {
	nodes := buildHockeyPostGameNodes(
		game.AwayTeam, game.HomeTeam,
		int(game.AwayTeamScore), int(game.HomeTeamScore),
		int(game.AwayTeamShootoutScore), int(game.HomeTeamShootoutScore),
		game.IsOvertime, game.IsShootout,
		game.Arena, game.City, game.State, game.Country,
		awayTeamStats.BaseTeamStats, homeTeamStats.BaseTeamStats,
		starOne, starTwo, starThree,
	)
	awayRows := toCHLPlayerStatRows(awayPlayerStats, collegePlayersMap)
	homeRows := toCHLPlayerStatRows(homePlayerStats, collegePlayersMap)
	nodes = appendHockeyPlayerStatTables(nodes, game.AwayTeam, awayRows, game.HomeTeam, homeRows)
	nodes = append(nodes, rtParagraph("Postgame discussion is open. Share your reactions below."))
	return nodes
}

func buildPHLPostGameNodes(
	game structs.ProfessionalGame,
	homeTeamStats structs.ProfessionalTeamGameStats,
	awayTeamStats structs.ProfessionalTeamGameStats,
	starOne, starTwo, starThree string,
	homePlayerStats []structs.ProfessionalPlayerGameStats,
	awayPlayerStats []structs.ProfessionalPlayerGameStats,
	proPlayersMap map[uint]structs.ProfessionalPlayer,
) []map[string]interface{} {
	nodes := buildHockeyPostGameNodes(
		game.AwayTeam, game.HomeTeam,
		int(game.AwayTeamScore), int(game.HomeTeamScore),
		int(game.AwayTeamShootoutScore), int(game.HomeTeamShootoutScore),
		game.IsOvertime, game.IsShootout,
		game.Arena, game.City, game.State, game.Country,
		awayTeamStats.BaseTeamStats, homeTeamStats.BaseTeamStats,
		starOne, starTwo, starThree,
	)
	awayRows := toPHLPlayerStatRows(awayPlayerStats, proPlayersMap)
	homeRows := toPHLPlayerStatRows(homePlayerStats, proPlayersMap)
	nodes = appendHockeyPlayerStatTables(nodes, game.AwayTeam, awayRows, game.HomeTeam, homeRows)
	nodes = append(nodes, rtParagraph("Postgame discussion is open. Share your reactions below."))
	return nodes
}

// buildHockeyPostGameNodes constructs the ordered list of ProseMirror content nodes
// for both CHL and PHL post-game forum posts.
func buildHockeyPostGameNodes(
	awayTeam, homeTeam string,
	awayScore, homeScore int,
	awayShootoutScore, homeShootoutScore int,
	isOvertime, isShootout bool,
	arena, city, state, country string,
	away, home structs.BaseTeamStats,
	starOne, starTwo, starThree string,
) []map[string]interface{} {
	nodes := []map[string]interface{}{}

	// ── Final score ──────────────────────────────────────────────────────────
	suffix := ""
	if isShootout {
		suffix = " (SO)"
	} else if isOvertime {
		suffix = " (OT)"
	}
	nodes = append(nodes, rtBoldParagraph(fmt.Sprintf(
		"FINAL%s: %s %d, %s %d",
		suffix, awayTeam, awayScore, homeTeam, homeScore,
	)))

	// ── Period scoring table ──────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Scoring by Period"))
	awayTotal := int(away.Period1Score) + int(away.Period2Score) + int(away.Period3Score) + int(away.OTScore)
	homeTotal := int(home.Period1Score) + int(home.Period2Score) + int(home.Period3Score) + int(home.OTScore)
	periodHeaders := []string{"Team", "P1", "P2", "P3", "OT", "Total"}
	periodRows := [][]string{
		{awayTeam, fmt.Sprintf("%d", away.Period1Score), fmt.Sprintf("%d", away.Period2Score), fmt.Sprintf("%d", away.Period3Score), fmt.Sprintf("%d", away.OTScore), fmt.Sprintf("%d", awayTotal)},
		{homeTeam, fmt.Sprintf("%d", home.Period1Score), fmt.Sprintf("%d", home.Period2Score), fmt.Sprintf("%d", home.Period3Score), fmt.Sprintf("%d", home.OTScore), fmt.Sprintf("%d", homeTotal)},
	}
	if isShootout {
		periodHeaders = []string{"Team", "P1", "P2", "P3", "OT", "SO", "Total"}
		periodRows = [][]string{
			{awayTeam, fmt.Sprintf("%d", away.Period1Score), fmt.Sprintf("%d", away.Period2Score), fmt.Sprintf("%d", away.Period3Score), fmt.Sprintf("%d", away.OTScore), fmt.Sprintf("%d", awayShootoutScore), fmt.Sprintf("%d", awayTotal)},
			{homeTeam, fmt.Sprintf("%d", home.Period1Score), fmt.Sprintf("%d", home.Period2Score), fmt.Sprintf("%d", home.Period3Score), fmt.Sprintf("%d", home.OTScore), fmt.Sprintf("%d", homeShootoutScore), fmt.Sprintf("%d", homeTotal)},
		}
	}
	nodes = append(nodes, rtTableNode(periodHeaders, periodRows))

	// ── Venue ─────────────────────────────────────────────────────────────────
	venueStr := arena
	if city != "" || state != "" || country != "" {
		location := city
		if state != "" {
			if location != "" {
				location += ", " + state
			} else {
				location = state
			}
		}
		if country != "" && country != "USA" && country != "US" {
			if location != "" {
				location += ", " + country
			} else {
				location = country
			}
		}
		if location != "" {
			venueStr += " — " + location
		}
	}
	if venueStr != "" {
		nodes = append(nodes, rtParagraph("Venue: "+venueStr))
	}

	// ── Three Stars ───────────────────────────────────────────────────────────
	if starOne != "" || starTwo != "" || starThree != "" {
		nodes = append(nodes, rtHeading(3, "Three Stars"))
		if starOne != "" {
			nodes = append(nodes, rtParagraph("⭐⭐⭐ "+starOne))
		}
		if starTwo != "" {
			nodes = append(nodes, rtParagraph("⭐⭐ "+starTwo))
		}
		if starThree != "" {
			nodes = append(nodes, rtParagraph("⭐ "+starThree))
		}
	}

	// ── Offense table ─────────────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Offense"))
	nodes = append(nodes, rtTableNode(
		[]string{"Team", "Goals", "Shots", "PP", "SH", "OT Goals"},
		[][]string{
			{awayTeam, fmt.Sprintf("%d", away.GoalsFor), fmt.Sprintf("%d", away.Shots), fmt.Sprintf("%d", away.PowerPlayGoals), fmt.Sprintf("%d", away.ShorthandedGoals), fmt.Sprintf("%d", away.OvertimeGoals)},
			{homeTeam, fmt.Sprintf("%d", home.GoalsFor), fmt.Sprintf("%d", home.Shots), fmt.Sprintf("%d", home.PowerPlayGoals), fmt.Sprintf("%d", home.ShorthandedGoals), fmt.Sprintf("%d", home.OvertimeGoals)},
		},
	))

	// ── Goaltending table ─────────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Goaltending"))
	nodes = append(nodes, rtTableNode(
		[]string{"Team", "Saves", "Shots Against", "SV%"},
		[][]string{
			{awayTeam, fmt.Sprintf("%d", away.Saves), fmt.Sprintf("%d", away.ShotsAgainst), fmt.Sprintf("%.3f", away.SavePercentage)},
			{homeTeam, fmt.Sprintf("%d", home.Saves), fmt.Sprintf("%d", home.ShotsAgainst), fmt.Sprintf("%.3f", home.SavePercentage)},
		},
	))

	return nodes
}

// formatPeriods returns a single line showing per-period scoring for a team.
func formatPeriods(team string, s structs.BaseTeamStats, shootoutScore int, isShootout bool) string {
	line := fmt.Sprintf("  %-20s  P1: %2d  P2: %2d  P3: %2d",
		team, s.Period1Score, s.Period2Score, s.Period3Score)
	if s.OTScore > 0 {
		line += fmt.Sprintf("  OT: %2d", s.OTScore)
	}
	if isShootout && shootoutScore > 0 {
		line += fmt.Sprintf("  SO: %2d", shootoutScore)
	}
	line += fmt.Sprintf("  TOTAL: %2d", int(s.Period1Score)+int(s.Period2Score)+int(s.Period3Score)+int(s.OTScore))
	return line
}

// buildHockeyPostGameThreadTitle returns a forum-ready title for a hockey game.
func buildHockeyPostGameThreadTitle(game structs.BaseGame) string {
	season := game.SeasonID + 2024
	if game.GameTitle != "" {
		return fmt.Sprintf("[%d] Week %d%s: %s: %s vs %s", season, game.Week, game.GameDay, game.GameTitle, game.AwayTeam, game.HomeTeam)
	}
	return fmt.Sprintf("[%d] Week %d%s: %s at %s", season, game.Week, game.GameDay, game.AwayTeam, game.HomeTeam)
}

// ─────────────────────────────────────────────
// Player name helpers
// ─────────────────────────────────────────────

// LookupCollegeStarName returns "FirstName LastName (Position)" for a college
// player ID, or an empty string if the player is not found.
func LookupCollegeStarName(id uint, playerMap map[uint]structs.CollegePlayer) string {
	if id == 0 {
		return ""
	}
	if p, ok := playerMap[id]; ok {
		return fmt.Sprintf("%s %s (%s)", p.FirstName, p.LastName, p.Position)
	}
	return ""
}

// LookupProStarName returns "FirstName LastName (Position)" for a professional
// player ID, or an empty string if the player is not found.
func LookupProStarName(id uint, playerMap map[uint]structs.ProfessionalPlayer) string {
	if id == 0 {
		return ""
	}
	if p, ok := playerMap[id]; ok {
		return fmt.Sprintf("%s %s (%s)", p.FirstName, p.LastName, p.Position)
	}
	return ""
}

// ─────────────────────────────────────────────
// Player stat helpers
// ─────────────────────────────────────────────

// hockeyPlayerStatRow pairs a display label and position with a player's per-game stats.
type hockeyPlayerStatRow struct {
	Label    string
	Position string
	structs.BasePlayerStats
}

// toCHLPlayerStatRows converts a slice of college player game stats into hockeyPlayerStatRows.
func toCHLPlayerStatRows(stats []structs.CollegePlayerGameStats, playerMap map[uint]structs.CollegePlayer) []hockeyPlayerStatRow {
	rows := make([]hockeyPlayerStatRow, 0, len(stats))
	for _, s := range stats {
		p, ok := playerMap[s.PlayerID]
		if !ok {
			continue
		}
		label := fmt.Sprintf("[%d] %s %s %s %s", p.ID, p.Team, p.Position, p.FirstName, p.LastName)
		rows = append(rows, hockeyPlayerStatRow{Label: label, Position: p.Position, BasePlayerStats: s.BasePlayerStats})
	}
	return rows
}

// toPHLPlayerStatRows converts a slice of professional player game stats into hockeyPlayerStatRows.
func toPHLPlayerStatRows(stats []structs.ProfessionalPlayerGameStats, playerMap map[uint]structs.ProfessionalPlayer) []hockeyPlayerStatRow {
	rows := make([]hockeyPlayerStatRow, 0, len(stats))
	for _, s := range stats {
		p, ok := playerMap[s.PlayerID]
		if !ok {
			continue
		}
		label := fmt.Sprintf("[%d] %s %s %s %s", p.ID, p.Team, p.Position, p.FirstName, p.LastName)
		rows = append(rows, hockeyPlayerStatRow{Label: label, Position: p.Position, BasePlayerStats: s.BasePlayerStats})
	}
	return rows
}

// appendHockeyPlayerStatTables appends per-player stat sections for forwards, defenders,
// and goalies to the node list. Within each section rows are sorted by team then by
// primary stat descending.
func appendHockeyPlayerStatTables(
	nodes []map[string]interface{},
	awayTeam string, awayRows []hockeyPlayerStatRow,
	homeTeam string, homeRows []hockeyPlayerStatRow,
) []map[string]interface{} {
	allRows := append(awayRows, homeRows...)

	// ── Forwards (Centers & Wingers) ─────────────────────────────────────────
	var fwdRows [][]string
	var forwards []hockeyPlayerStatRow
	for _, r := range allRows {
		if r.Position == "C" || r.Position == "F" {
			forwards = append(forwards, r)
		}
	}
	sort.Slice(forwards, func(i, j int) bool {
		if forwards[i].TeamID != forwards[j].TeamID {
			return forwards[i].TeamID < forwards[j].TeamID
		}
		if forwards[i].Goals != forwards[j].Goals {
			return forwards[i].Goals > forwards[j].Goals
		}
		return forwards[i].Points > forwards[j].Points
	})
	for _, r := range forwards {
		shotPct := "0.0%"
		if r.Shots > 0 {
			shotPct = fmt.Sprintf("%.1f%%", float32(r.Goals)/float32(r.Shots)*100)
		}
		fwdRows = append(fwdRows, []string{
			r.Label,
			fmt.Sprintf("%d", r.Goals),
			fmt.Sprintf("%d", r.Assists),
			fmt.Sprintf("%d", r.Points),
			fmt.Sprintf("%+d", r.PlusMinus),
			fmt.Sprintf("%d", r.PenaltyMinutes),
			fmt.Sprintf("%d", r.PowerPlayGoals),
			fmt.Sprintf("%d", r.ShorthandedGoals),
			fmt.Sprintf("%d", r.Shots),
			shotPct,
		})
	}
	if len(fwdRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Forwards"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "G", "A", "PTS", "+/-", "PIM", "PPG", "SHG", "Shots", "S%"},
			fwdRows,
		))
	}

	// ── Defenders ─────────────────────────────────────────────────────────────
	var defRows [][]string
	var defenders []hockeyPlayerStatRow
	for _, r := range allRows {
		if r.Position == "D" {
			defenders = append(defenders, r)
		}
	}
	sort.Slice(defenders, func(i, j int) bool {
		if defenders[i].TeamID != defenders[j].TeamID {
			return defenders[i].TeamID < defenders[j].TeamID
		}
		if defenders[i].Points != defenders[j].Points {
			return defenders[i].Points > defenders[j].Points
		}
		return defenders[i].PlusMinus > defenders[j].PlusMinus
	})
	for _, r := range defenders {
		defRows = append(defRows, []string{
			r.Label,
			fmt.Sprintf("%d", r.Goals),
			fmt.Sprintf("%d", r.Assists),
			fmt.Sprintf("%d", r.Points),
			fmt.Sprintf("%+d", r.PlusMinus),
			fmt.Sprintf("%d", r.PenaltyMinutes),
			fmt.Sprintf("%d", r.ShotsBlocked),
			fmt.Sprintf("%d", r.BodyChecks),
		})
	}
	if len(defRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Defenders"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "G", "A", "PTS", "+/-", "PIM", "ShotBlk", "Checks"},
			defRows,
		))
	}

	// ── Goalies ───────────────────────────────────────────────────────────────
	var goalieRows [][]string
	var goalies []hockeyPlayerStatRow
	for _, r := range allRows {
		if r.Position == "G" {
			goalies = append(goalies, r)
		}
	}
	sort.Slice(goalies, func(i, j int) bool {
		if goalies[i].TeamID != goalies[j].TeamID {
			return goalies[i].TeamID < goalies[j].TeamID
		}
		return goalies[i].Saves > goalies[j].Saves
	})
	for _, r := range goalies {
		svPct := ".000"
		if r.ShotsAgainst > 0 {
			svPct = fmt.Sprintf("%.3f", float32(r.Saves)/float32(r.ShotsAgainst))
		}
		goalieRows = append(goalieRows, []string{
			r.Label,
			fmt.Sprintf("%d", r.Saves),
			fmt.Sprintf("%d", r.ShotsAgainst),
			fmt.Sprintf("%d", r.GoalsAgainst),
			svPct,
			fmt.Sprintf("%d", r.Shutouts),
			fmt.Sprintf("%d-%d", r.GoalieWins, r.GoalieLosses),
		})
	}
	if len(goalieRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Goalies"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "Saves", "SA", "GA", "SV%", "SO", "W-L"},
			goalieRows,
		))
	}

	_ = awayTeam
	_ = homeTeam
	return nodes
}

// ─────────────────────────────────────────────
// Transfer portal helpers
// ─────────────────────────────────────────────

// TransferIntentionsSummary bundles counters produced by the transfer
// intentions run so they can be passed to the forum-thread creator.
type TransferIntentionsSummary struct {
	Season                 int
	TransferCount          int
	FreshmanCount          int
	RedshirtFreshmanCount  int
	SophomoreCount         int
	RedshirtSophomoreCount int
	JuniorCount            int
	RedshirtJuniorCount    int
	SeniorCount            int
	RedshirtSeniorCount    int
	LowCount               int
	MediumCount            int
	HighCount              int
}

// CreateTransferIntentionsForumThread creates a system-generated thread
// summarising the transfer intentions run for the given season.
// The operation is idempotent: calling it twice for the same season has no effect.
func CreateTransferIntentionsForumThread(summary TransferIntentionsSummary) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCHL: Season %d Transfer Intentions", summary.Season)
	eventKey := fmt.Sprintf("transfer_intentions_thread:chl:season%d", summary.Season)

	paragraphs := buildTransferIntentionsParagraphs(summary)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simchl",
		ForumPath:         []string{"media", "simchl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "chl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer intentions thread for season %d: %v", summary.Season, err)
		return
	}

	log.Printf("ForumManager: created transfer intentions thread %s for season %d", thread.ID, summary.Season)
}

func buildTransferIntentionsParagraphs(s TransferIntentionsSummary) []string {
	var paras []string

	paras = append(paras,
		fmt.Sprintf(
			"Transfer season is underway for Season %d. A total of %d players have announced their intention to enter the transfer portal. Teams have one week to submit promises to retain their players.",
			s.Season, s.TransferCount,
		),
	)

	paras = append(paras,
		fmt.Sprintf(
			"Class breakdown — Freshmen: %d | RS Freshmen: %d | Sophomores: %d | RS Sophomores: %d | Juniors: %d | RS Juniors: %d | Seniors: %d | RS Seniors: %d.",
			s.FreshmanCount, s.RedshirtFreshmanCount,
			s.SophomoreCount, s.RedshirtSophomoreCount,
			s.JuniorCount, s.RedshirtJuniorCount,
			s.SeniorCount, s.RedshirtSeniorCount,
		),
	)

	paras = append(paras,
		fmt.Sprintf(
			"Transfer likeliness — Low: %d | Medium: %d | High: %d.",
			s.LowCount, s.MediumCount, s.HighCount,
		),
	)

	paras = append(paras,
		"Which transfers are you keeping an eye on this season? Share your thoughts below!",
	)

	return paras
}

// ─────────────────────────────────────────────
// Transfer portal open thread
// ─────────────────────────────────────────────

// CreateTransferPortalOpenForumThread creates a system-generated thread
// announcing the transfer portal is open for the given season, with one entry
// per player that officially entered the portal.
// playerLabels should be built before WillTransfer() clears each player's team.
// The operation is idempotent: calling it twice for the same season has no effect.
func CreateTransferPortalOpenForumThread(season int, playerLabels []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCHL: Season %d Transfer Portal is Now Open", season)
	eventKey := fmt.Sprintf("transfer_portal_open:chl:season%d", season)

	paragraphs := buildTransferPortalOpenParagraphs(season, playerLabels)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simchl",
		ForumPath:         []string{"media", "simchl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "chl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer portal open thread for season %d: %v", season, err)
		return
	}

	log.Printf("ForumManager: created transfer portal open thread %s for season %d", thread.ID, season)
}

func buildTransferPortalOpenParagraphs(season int, playerLabels []string) []string {
	var paras []string

	count := len(playerLabels)
	paras = append(paras,
		fmt.Sprintf(
			"The SimCHL Transfer Portal is now open for Season %d. A total of %d player(s) have officially entered the portal and are seeking a new home.",
			season, count,
		),
	)

	if count > 0 {
		paras = append(paras, "The following players have entered the transfer portal:")
		for _, label := range playerLabels {
			paras = append(paras, label)
		}
	}

	paras = append(paras, "Which players are you targeting this transfer portal cycle? Share your thoughts below!")

	return paras
}

// ─────────────────────────────────────────────
// Transfer portal sync thread
// ─────────────────────────────────────────────

// CreateTransferPortalSyncForumThread creates a system-generated thread
// summarising the signings from a single transfer portal sync round.
// signings is a list of human-readable labels for every player that signed.
// The operation is idempotent: calling it twice for the same season/round has no
// effect.
func CreateTransferPortalSyncForumThread(season, round int, signings []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCHL: Season %d Transfer Portal — Round %d Results", season, round)
	eventKey := fmt.Sprintf("transfer_portal_sync:chl:season%d:round%d", season, round)

	paragraphs := buildTransferPortalSyncParagraphs(season, round, signings)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simchl",
		ForumPath:         []string{"media", "simchl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "chl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer portal sync thread for season %d round %d: %v", season, round, err)
		return
	}

	log.Printf("ForumManager: created transfer portal sync thread %s for season %d round %d", thread.ID, season, round)
}

func buildTransferPortalSyncParagraphs(season, round int, signings []string) []string {
	var paras []string

	count := len(signings)
	if count == 0 {
		paras = append(paras,
			fmt.Sprintf(
				"Transfer Portal Round %d is complete for Season %d. No players signed with new programs this round.",
				round, season,
			),
		)
	} else {
		paras = append(paras,
			fmt.Sprintf(
				"Transfer Portal Round %d results are in for Season %d. A total of %d player(s) have signed with new programs this round.",
				round, season, count,
			),
		)
		for _, label := range signings {
			paras = append(paras, label)
		}
	}

	paras = append(paras, "Discuss the latest transfer portal news below!")

	return paras
}

// ─────────────────────────────────────────────
// Recruiting sync thread
// ─────────────────────────────────────────────

// CreateRecruitingSyncForumThread creates a system-generated weekly thread
// listing every recruit that signed with a program during the sync.
// signings is a list of human-readable labels built at commit time.
// The operation is idempotent: calling it twice for the same season/week has no
// effect.
func CreateRecruitingSyncForumThread(season, week int, signings []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCHL: Season %d Week %d Recruiting Commitments", season, week)
	eventKey := fmt.Sprintf("recruiting_sync:chl:season%d:week%d", season, week)

	paragraphs := buildRecruitingSyncParagraphs(season, week, signings)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simchl",
		ForumPath:         []string{"media", "simchl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "chl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create recruiting sync thread for season %d week %d: %v", season, week, err)
		return
	}

	log.Printf("ForumManager: created recruiting sync thread %s for season %d week %d", thread.ID, season, week)
}

func buildRecruitingSyncParagraphs(season, week int, signings []string) []string {
	var paras []string

	count := len(signings)
	if count == 0 {
		paras = append(paras,
			fmt.Sprintf(
				"Week %d recruiting is complete for Season %d. No recruits signed with a program this week.",
				week, season,
			),
		)
	} else {
		paras = append(paras,
			fmt.Sprintf(
				"Week %d recruiting results are in for Season %d. A total of %d recruit(s) have committed to a program this week.",
				week, season, count,
			),
		)
		for _, label := range signings {
			paras = append(paras, label)
		}
	}

	paras = append(paras, "React to the latest commitments and discuss your team\u2019s recruiting class below!")

	return paras
}

// ─────────────────────────────────────────────
// Rich text helpers
// ─────────────────────────────────────────────

// buildRichPostBody converts a slice of paragraph strings into a ProseMirror
// document compatible with the frontend's RichTextDocument interface.
func buildRichPostBody(paragraphs []string) map[string]interface{} {
	content := make([]map[string]interface{}, 0, len(paragraphs))
	for _, p := range paragraphs {
		content = append(content, map[string]interface{}{
			"type": "paragraph",
			"content": []map[string]interface{}{
				{"type": "text", "text": p},
			},
		})
	}
	return map[string]interface{}{
		"type":    "doc",
		"content": content,
	}
}

// buildRichDoc wraps a slice of content nodes into a top-level ProseMirror doc.
func buildRichDoc(nodes []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":    "doc",
		"content": nodes,
	}
}

// rtParagraph creates a plain paragraph node.
func rtParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

// rtBoldParagraph creates a paragraph with bold-marked text.
func rtBoldParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{
				"type":  "text",
				"text":  text,
				"marks": []map[string]interface{}{{"type": "bold"}},
			},
		},
	}
}

// rtHeading creates a heading node at the given level (1–6).
func rtHeading(level int, text string) map[string]interface{} {
	return map[string]interface{}{
		"type":  "heading",
		"attrs": map[string]interface{}{"level": level, "textAlign": "left"},
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

// rtTableCell creates a single table header or data cell wrapping text in a paragraph.
func rtTableCell(text string, isHeader bool) map[string]interface{} {
	cellType := "tableCell"
	if isHeader {
		cellType = "tableHeader"
	}
	return map[string]interface{}{
		"type":  cellType,
		"attrs": map[string]interface{}{"colspan": 1, "rowspan": 1, "colwidth": nil},
		"content": []map[string]interface{}{
			{
				"type":  "paragraph",
				"attrs": map[string]interface{}{"textAlign": nil},
				"content": []map[string]interface{}{
					{"type": "text", "text": text},
				},
			},
		},
	}
}

// rtTableNode builds a TipTap-compatible table node from header strings and row data.
// The first row is rendered as tableHeader cells; all subsequent rows as tableCell.
func rtTableNode(headers []string, rows [][]string) map[string]interface{} {
	tableRows := []map[string]interface{}{}

	headerCells := make([]map[string]interface{}, len(headers))
	for i, h := range headers {
		headerCells[i] = rtTableCell(h, true)
	}
	tableRows = append(tableRows, map[string]interface{}{
		"type":    "tableRow",
		"content": headerCells,
	})

	for _, row := range rows {
		cells := make([]map[string]interface{}, len(row))
		for i, cell := range row {
			cells[i] = rtTableCell(cell, false)
		}
		tableRows = append(tableRows, map[string]interface{}{
			"type":    "tableRow",
			"content": cells,
		})
	}

	return map[string]interface{}{
		"type":    "table",
		"content": tableRows,
	}
}

// nodesToPlainText extracts readable plain text from rich nodes for the bodyText field.
func nodesToPlainText(nodes []map[string]interface{}) string {
	var lines []string
	for _, node := range nodes {
		switch node["type"] {
		case "paragraph", "heading":
			if text := extractInlineText(node); text != "" {
				lines = append(lines, text)
			}
		case "table":
			if rows, ok := node["content"].([]map[string]interface{}); ok {
				for _, row := range rows {
					if cells, ok := row["content"].([]map[string]interface{}); ok {
						var cellTexts []string
						for _, cell := range cells {
							cellTexts = append(cellTexts, extractCellPlainText(cell))
						}
						lines = append(lines, strings.Join(cellTexts, "  |  "))
					}
				}
			}
		}
	}
	return strings.Join(lines, "\n\n")
}

func extractInlineText(node map[string]interface{}) string {
	if content, ok := node["content"].([]map[string]interface{}); ok {
		var texts []string
		for _, child := range content {
			if t, ok := child["text"].(string); ok {
				texts = append(texts, t)
			}
		}
		return strings.Join(texts, "")
	}
	return ""
}

func extractCellPlainText(cell map[string]interface{}) string {
	if content, ok := cell["content"].([]map[string]interface{}); ok {
		for _, para := range content {
			return extractInlineText(para)
		}
	}
	return ""
}
