package managers

// StreamGameQueueItem is the wire shape simsn-live consumes to build its
// own PendingGame queue without touching SimHockey's database directly.
type StreamGameQueueItem struct {
	GameID       uint   `json:"gameID"`
	HomeTeamID   uint   `json:"homeTeamID"`
	AwayTeamID   uint   `json:"awayTeamID"`
	HomeTeam     string `json:"homeTeam"`
	AwayTeam     string `json:"awayTeam"`
	IsUserGame   bool   `json:"isUserGame"`
	HomeTeamRank int    `json:"homeTeamRank"`
	AwayTeamRank int    `json:"awayTeamRank"`
	Arena        string `json:"arena"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	TotalSeconds int    `json:"totalSeconds"`
	TotalPlays   int    `json:"totalPlays"`
}

// BuildStreamGameQueue returns the queued, unrevealed games for a league
// (chl|phl), ordered user-games-first, with PbP totals attached.
func BuildStreamGameQueue(league, weekID, seasonID, gameDay string, isPreseason bool) []StreamGameQueueItem {
	if league == "chl" {
		return buildCHLStreamQueue(weekID, seasonID, gameDay, isPreseason)
	}
	return buildPHLStreamQueue(weekID, seasonID, gameDay, isPreseason)
}

func buildCHLStreamQueue(weekID, seasonID, gameDay string, isPreseason bool) []StreamGameQueueItem {
	var userGames, aiGames []StreamGameQueueItem

	games := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
	teamMap := GetCollegeTeamMap()
	for _, g := range games {
		if !g.GameComplete || g.IsRevealed {
			continue
		}
		homeTeam := teamMap[g.HomeTeamID]
		awayTeam := teamMap[g.AwayTeamID]
		item := StreamGameQueueItem{
			GameID:       g.ID,
			HomeTeamID:   g.HomeTeamID,
			AwayTeamID:   g.AwayTeamID,
			HomeTeam:     homeTeam.Abbreviation,
			AwayTeam:     awayTeam.Abbreviation,
			IsUserGame:   homeTeam.IsUserCoached || awayTeam.IsUserCoached,
			HomeTeamRank: int(g.HomeTeamRank),
			AwayTeamRank: int(g.AwayTeamRank),
			Arena:        g.Arena,
			City:         g.City,
			State:        g.State,
			Country:      g.Country,
			TotalSeconds: loadTotalSeconds(g.ID, true),
			TotalPlays:   loadTotalPlays(g.ID, true),
		}
		if item.IsUserGame {
			userGames = append(userGames, item)
		} else {
			aiGames = append(aiGames, item)
		}
	}
	return append(userGames, aiGames...)
}

func buildPHLStreamQueue(weekID, seasonID, gameDay string, isPreseason bool) []StreamGameQueueItem {
	var userGames, aiGames []StreamGameQueueItem

	games := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
	teamMap := GetProTeamMap()
	for _, g := range games {
		if !g.GameComplete || g.IsRevealed {
			continue
		}
		homeTeam := teamMap[g.HomeTeamID]
		awayTeam := teamMap[g.AwayTeamID]
		isUser := homeTeam.Owner != "" || awayTeam.Owner != "" ||
			homeTeam.GM != "" || awayTeam.GM != ""
		item := StreamGameQueueItem{
			GameID:       g.ID,
			HomeTeamID:   g.HomeTeamID,
			AwayTeamID:   g.AwayTeamID,
			HomeTeam:     homeTeam.Abbreviation,
			AwayTeam:     awayTeam.Abbreviation,
			IsUserGame:   isUser,
			HomeTeamRank: int(g.HomeTeamRank),
			AwayTeamRank: int(g.AwayTeamRank),
			Arena:        g.Arena,
			City:         g.City,
			State:        g.State,
			Country:      g.Country,
			TotalSeconds: loadTotalSeconds(g.ID, false),
			TotalPlays:   loadTotalPlays(g.ID, false),
		}
		if item.IsUserGame {
			userGames = append(userGames, item)
		} else {
			aiGames = append(aiGames, item)
		}
	}
	return append(userGames, aiGames...)
}
