package managers

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func ConductDraftLottery() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	// Get the pre-lottery base order from standings/series data.
	// Lottery teams (Chances set) come first, sorted worst→best record.
	// Playoff teams (Chances empty) follow, sorted by elimination round then record.
	baseOrder := GetDraftLotteryOrder(db)
	draftMap := GetDraftPickMapForLottery(seasonID)

	// Split into lottery balls and playoff order based on whether Chances is set.
	lotteryBalls := []structs.DraftLottery{}
	playoffOrder := []structs.DraftLottery{}
	for _, entry := range baseOrder {
		if len(entry.Chances) > 0 {
			lotteryBalls = append(lotteryBalls, entry)
		} else {
			playoffOrder = append(playoffOrder, entry)
		}
	}

	// Run weighted lottery for picks 1-4.
	finalRound1 := []structs.DraftLottery{}
	remaining := make([]structs.DraftLottery, len(lotteryBalls))
	copy(remaining, lotteryBalls)

	for i := 0; i < 4; i++ {
		sum := 0
		for _, l := range remaining {
			l.ApplyCurrentChance(i)
			sum += int(l.CurrentChance)
		}
		chance := util.GenerateIntFromRange(1, sum)
		sum2 := 0
		for _, l := range remaining {
			l.ApplyCurrentChance(i)
			sum2 += int(l.CurrentChance)
			if chance < sum2 {
				finalRound1 = append(finalRound1, l)
				remaining = filterLotteryPicks(remaining, l.ID)
				break
			}
		}
	}
	// Picks 5-N: remaining lottery teams keep their base order.
	finalRound1 = append(finalRound1, remaining...)
	// Picks N+1 through 32: playoff teams in elimination order.
	finalRound1 = append(finalRound1, playoffOrder...)

	// Assign Round 1 draft numbers and collect pick records.
	draftPicks := []structs.DraftPick{}
	round1PickMap := make(map[uint]uint) // teamID → R1 pick number (for R2 tiebreaking)
	for idx, team := range finalRound1 {
		pickNum := uint(idx + 1)
		round1PickMap[team.ID] = pickNum
		key := "1 " + strconv.Itoa(int(team.ID))
		pick := draftMap[key]
		pick.AssignDraftNumber(pickNum)
		draftPicks = append(draftPicks, pick)
	}

	// Build Round 2 order: all teams sorted by record (worst first), no playoff
	// consideration. Ties broken by reverse Round 1 pick (higher R1 pick = earlier R2 pick).
	prevSeasonID := strconv.Itoa(int(ts.SeasonID - 1))
	allStandings := repository.FindAllProfessionalStandings(repository.StandingsQuery{SeasonID: prevSeasonID})
	standingsMap := make(map[uint]structs.ProfessionalStandings)
	for _, s := range allStandings {
		if s.TeamID <= 32 {
			standingsMap[s.TeamID] = s
		}
	}

	round2Order := make([]structs.DraftLottery, len(finalRound1))
	copy(round2Order, finalRound1)
	sort.Slice(round2Order, func(i, j int) bool {
		si := standingsMap[round2Order[i].ID]
		sj := standingsMap[round2Order[j].ID]
		if si.TotalWins != sj.TotalWins {
			return si.TotalWins < sj.TotalWins
		}
		return round1PickMap[round2Order[i].ID] > round1PickMap[round2Order[j].ID]
	})

	for idx, team := range round2Order {
		pickNum := uint(idx + 1)
		key := "2 " + strconv.Itoa(int(team.ID))
		pick := draftMap[key]
		pick.AssignDraftNumber(pickNum)
		draftPicks = append(draftPicks, pick)
	}

	sort.Sort(structs.ByDraftNumber(draftPicks))

	for _, pick := range draftPicks {
		fmt.Println("Pick " + strconv.Itoa(int(pick.DraftNumber)) + ": " + pick.OriginalTeam)
		db.Save(&pick)
	}

	// Create a draft lottery forum thread (best-effort).
	season, picks := ts.Season, draftPicks
	go CreateDraftLotteryForumThread(int(season), picks)
}

func GetDraftLotteryOrder(db *gorm.DB) []structs.DraftLottery {
	lotteryOrder := []structs.DraftLottery{}

	proTeams := repository.FindAllProTeams(repository.TeamClauses{})
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID - 1))
	proSeries := repository.FindProSeriesRecords(seasonID)
	proGames := repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID})

	// Build full team name map (City + Nickname) for NBA teams only (ID <= 32)
	teamNameMap := make(map[uint]string)
	for _, t := range proTeams {
		if t.ID <= 24 {
			teamNameMap[t.ID] = t.TeamName + " " + t.Mascot
		}
	}

	// Load previous season standings for NBA teams (ID <= 32)
	allStandings := repository.FindAllProfessionalStandings(repository.StandingsQuery{SeasonID: seasonID})
	standingsMap := make(map[uint]structs.ProfessionalStandings)
	proStandings := []structs.ProfessionalStandings{}
	for _, s := range allStandings {
		if s.TeamID <= 24 {
			standingsMap[s.TeamID] = s
			proStandings = append(proStandings, s)
		}
	}

	// Build head-to-head win map from regular season games only.
	// Exclude playoff games, play-in games (week 18), and international teams.
	headToHead := make(map[uint]map[uint]int)
	for _, g := range proGames {
		if g.IsPlayoffGame || !g.GameComplete {
			continue
		}
		if g.HomeTeamID > 24 || g.AwayTeamID > 24 {
			continue
		}
		if headToHead[g.HomeTeamID] == nil {
			headToHead[g.HomeTeamID] = make(map[uint]int)
		}
		if headToHead[g.AwayTeamID] == nil {
			headToHead[g.AwayTeamID] = make(map[uint]int)
		}
		if g.HomeTeamWin {
			headToHead[g.HomeTeamID][g.AwayTeamID]++
		} else if g.AwayTeamWin {
			headToHead[g.AwayTeamID][g.HomeTeamID]++
		}
	}

	// Count playoff series wins per team to determine the round they were eliminated.
	// 0 wins = Round 1 loser, 1 = Round 2 loser, 2 = Conference Finals loser,
	// 3 = Finals loser, 4 = Champion.
	seriesWins := make(map[uint]int)
	for _, s := range proSeries {
		if !s.IsPlayoffGame || s.IsInternational || !s.SeriesComplete {
			continue
		}
		if s.HomeTeamSeriesWin && s.HomeTeamID <= 24 {
			seriesWins[s.HomeTeamID]++
		} else if !s.HomeTeamSeriesWin && s.AwayTeamID <= 24 {
			seriesWins[s.AwayTeamID]++
		}
	}

	// Collect all NBA teams that appeared in a playoff series
	type playoffEntry struct {
		TeamID   uint
		TeamName string
		Wins     int
		Standing structs.ProfessionalStandings
	}
	playoffMap := make(map[uint]playoffEntry)
	for _, s := range proSeries {
		if !s.IsPlayoffGame || s.IsInternational || !s.SeriesComplete {
			continue
		}
		for _, id := range []uint{s.HomeTeamID, s.AwayTeamID} {
			if id > 24 {
				continue
			}
			name := teamNameMap[id]
			if name == "" {
				if id == s.HomeTeamID {
					name = s.HomeTeam
				} else {
					name = s.AwayTeam
				}
			}
			playoffMap[id] = playoffEntry{
				TeamID:   id,
				TeamName: name,
				Wins:     seriesWins[id],
				Standing: standingsMap[id],
			}
		}
	}

	playoffTeamIDs := make(map[uint]bool)
	for id := range playoffMap {
		playoffTeamIDs[id] = true
	}

	// --- BOTTOM 16: Lottery teams ---
	// Any NBA team (ID <= 24) not in a playoff series is a lottery team.
	// This naturally includes play-in losers (they never appear in NBASeries).
	lotteryStandings := []structs.ProfessionalStandings{}
	for _, s := range proStandings {
		if !playoffTeamIDs[s.TeamID] {
			lotteryStandings = append(lotteryStandings, s)
		}
	}

	// Sort: worst record first (fewest wins).
	// Tiebreaker 1: head-to-head wins between the tied teams.
	// Tiebreaker 2: point differential (lower = picks earlier).
	// Tiebreaker 3: random coin flip.
	sort.Slice(lotteryStandings, func(i, j int) bool {
		si, sj := lotteryStandings[i], lotteryStandings[j]
		if si.TotalWins != sj.TotalWins {
			return si.TotalWins < sj.TotalWins
		}
		h2hI := headToHead[si.TeamID][sj.TeamID]
		h2hJ := headToHead[sj.TeamID][si.TeamID]
		if h2hI != h2hJ {
			return h2hI < h2hJ
		}
		if si.GoalsFor-si.GoalsAgainst != sj.GoalsFor-sj.GoalsAgainst {
			return (si.GoalsFor - si.GoalsAgainst) < (sj.GoalsFor - sj.GoalsAgainst)
		}
		return rand.Intn(2) == 0
	})

	// Assign lottery ball chances based on position (1 = worst team, 20 = best lottery team)
	for idx, s := range lotteryStandings {
		chances := util.Get20TeamLotteryChances(idx + 1)
		name := teamNameMap[s.TeamID]
		if name == "" {
			name = s.TeamName
		}
		lotteryOrder = append(lotteryOrder, structs.DraftLottery{
			ID:      s.TeamID,
			Team:    name,
			Chances: chances,
		})
	}

	// --- TOP 16: Playoff teams ---
	// Grouped by series wins (proxy for elimination round).
	// Within each round group, sort by regular season record (worst first).
	sortPlayoffGroup := func(group []playoffEntry) {
		sort.Slice(group, func(i, j int) bool {
			si, sj := group[i].Standing, group[j].Standing
			if si.TotalWins != sj.TotalWins {
				return si.TotalWins < sj.TotalWins
			}
			h2hI := headToHead[group[i].TeamID][group[j].TeamID]
			h2hJ := headToHead[group[j].TeamID][group[i].TeamID]
			if h2hI != h2hJ {
				return h2hI < h2hJ
			}
			return (si.GoalsFor - si.GoalsAgainst) < (sj.GoalsFor - sj.GoalsAgainst)
		})
	}

	roundGroups := make(map[int][]playoffEntry)
	for _, entry := range playoffMap {
		if entry.TeamID > 32 {
			continue
		}
		roundGroups[entry.Wins] = append(roundGroups[entry.Wins], entry)
	}

	// Append in order: R1 losers (0 wins) → Semis losers (1) → CF losers (2) → Finals loser (3) → Champion (4)
	for wins := 0; wins <= 4; wins++ {
		group := roundGroups[wins]
		sortPlayoffGroup(group)
		for _, entry := range group {
			lotteryOrder = append(lotteryOrder, structs.DraftLottery{
				ID:      entry.TeamID,
				Team:    entry.TeamName,
				Chances: []uint{},
			})
		}
	}

	return lotteryOrder
}

func GetDraftPickMapForLottery(seasonID string) map[string]structs.DraftPick {
	draftPicks := repository.FindDraftPicks(seasonID)
	draftMap := make(map[string]structs.DraftPick)

	for _, pick := range draftPicks {
		if pick.ID == 0 {
			continue
		}
		keyString := strconv.Itoa(int(pick.DraftRound)) + " " + strconv.Itoa(int(pick.OriginalTeamID))
		draftMap[keyString] = pick
	}
	return draftMap
}

func filterLotteryPicks(list []structs.DraftLottery, id uint) []structs.DraftLottery {
	newList := []structs.DraftLottery{}
	for _, l := range list {
		if l.ID != id {
			newList = append(newList, l)
		}
	}
	return newList
}
