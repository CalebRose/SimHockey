package managers

import (
	"context"
	"fmt"
	"log"
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
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:chl:season%d:game%s", seasonID, gameID)

	title := buildHockeyPostGameThreadTitle(game.BaseGame)
	paragraphs := buildCHLPostGameParagraphs(game, homeTeamStats, awayTeamStats, starOne, starTwo, starThree)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

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
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:phl:season%d:game%s", seasonID, gameID)

	title := buildHockeyPostGameThreadTitle(game.BaseGame)
	paragraphs := buildPHLPostGameParagraphs(game, homeTeamStats, awayTeamStats, starOne, starTwo, starThree)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

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

func buildCHLPostGameParagraphs(
	game structs.CollegeGame,
	homeTeamStats structs.CollegeTeamGameStats,
	awayTeamStats structs.CollegeTeamGameStats,
	starOne, starTwo, starThree string,
) []string {
	return buildHockeyPostGameParagraphs(
		game.AwayTeam, game.HomeTeam,
		int(game.AwayTeamScore), int(game.HomeTeamScore),
		int(game.AwayTeamShootoutScore), int(game.HomeTeamShootoutScore),
		game.IsOvertime, game.IsShootout,
		game.Arena, game.City, game.State, game.Country,
		awayTeamStats.BaseTeamStats, homeTeamStats.BaseTeamStats,
		starOne, starTwo, starThree,
	)
}

func buildPHLPostGameParagraphs(
	game structs.ProfessionalGame,
	homeTeamStats structs.ProfessionalTeamGameStats,
	awayTeamStats structs.ProfessionalTeamGameStats,
	starOne, starTwo, starThree string,
) []string {
	return buildHockeyPostGameParagraphs(
		game.AwayTeam, game.HomeTeam,
		int(game.AwayTeamScore), int(game.HomeTeamScore),
		int(game.AwayTeamShootoutScore), int(game.HomeTeamShootoutScore),
		game.IsOvertime, game.IsShootout,
		game.Arena, game.City, game.State, game.Country,
		awayTeamStats.BaseTeamStats, homeTeamStats.BaseTeamStats,
		starOne, starTwo, starThree,
	)
}

// buildHockeyPostGameParagraphs constructs the ordered list of paragraph strings
// for both CHL and PHL post-game forum posts.
func buildHockeyPostGameParagraphs(
	awayTeam, homeTeam string,
	awayScore, homeScore int,
	awayShootoutScore, homeShootoutScore int,
	isOvertime, isShootout bool,
	arena, city, state, country string,
	away, home structs.BaseTeamStats,
	starOne, starTwo, starThree string,
) []string {
	paras := []string{}

	// ── Final score ──────────────────────────────────────────────────────────
	suffix := ""
	if isShootout {
		suffix = " (SO)"
	} else if isOvertime {
		suffix = " (OT)"
	}
	paras = append(paras, fmt.Sprintf(
		"FINAL%s: %s %d, %s %d",
		suffix, awayTeam, awayScore, homeTeam, homeScore,
	))

	// ── Period-by-period scoring ──────────────────────────────────────────────
	awayPeriods := formatPeriods(awayTeam, away, awayShootoutScore, isShootout)
	homePeriods := formatPeriods(homeTeam, home, homeShootoutScore, isShootout)
	paras = append(paras, "SCORING BY PERIOD:\n"+awayPeriods)
	paras = append(paras, "\n"+homePeriods)
	paras = append(paras, "Note: OT and SO scoring is included in the final score but not the period breakdown.")

	// ── Offensive stats ───────────────────────────────────────────────────────
	offLines := []string{
		"OFFENSE:\n",
		fmt.Sprintf("  %-20s  Goals: %2d   Shots: %3d   PP: %d   SH: %d   OT: %d\n",
			awayTeam,
			away.GoalsFor, away.Shots,
			away.PowerPlayGoals, away.ShorthandedGoals, away.OvertimeGoals,
		),
		fmt.Sprintf("  %-20s  Goals: %2d   Shots: %3d   PP: %d   SH: %d   OT: %d\n",
			homeTeam,
			home.GoalsFor, home.Shots,
			home.PowerPlayGoals, home.ShorthandedGoals, home.OvertimeGoals,
		),
	}
	paras = append(paras, strings.Join(offLines, "\n"))

	// ── Goaltending ──────────────────────────────────────────────────────────
	gtLines := []string{
		"GOALTENDING:\n",
		fmt.Sprintf("  %-20s  Saves: %3d / %3d   SV%%: %.3f\n",
			awayTeam,
			away.Saves, away.ShotsAgainst, away.SavePercentage,
		),
		fmt.Sprintf("  %-20s  Saves: %3d / %3d   SV%%: %.3f\n",
			homeTeam,
			home.Saves, home.ShotsAgainst, home.SavePercentage,
		),
	}
	paras = append(paras, strings.Join(gtLines, "\n"))

	// ── Venue ────────────────────────────────────────────────────────────────
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
		paras = append(paras, "VENUE: "+venueStr)
	}

	// ── Three Stars ───────────────────────────────────────────────────────────
	if starOne != "" || starTwo != "" || starThree != "" {
		starsLines := []string{"THREE STARS:\n"}
		if starOne != "" {
			starsLines = append(starsLines, "  ⭐⭐⭐ "+starOne+"\n")
		}
		if starTwo != "" {
			starsLines = append(starsLines, "  ⭐⭐  "+starTwo+"\n")
		}
		if starThree != "" {
			starsLines = append(starsLines, "  ⭐     "+starThree)
		}
		paras = append(paras, strings.Join(starsLines, "\n"))
	}

	// ── Discussion prompt ─────────────────────────────────────────────────────
	paras = append(paras, "Postgame discussion is open. Share your reactions below.")

	return paras
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
