package managers

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

type ScheduledGame struct {
	structs.ProfessionalGame
	Week int
	Slot string // "A", "B", "C" or "D"
}

func (sg *ScheduledGame) ApplyWeekAndSlot(week int, slot string) {
	sg.Week = week
	sg.Slot = slot
}

// PHLScheduleConfig describes the season format.
//   - Each divisional pair plays DivisionalGamesPerPair games (half home, half away).
//   - Each same-conference, other-division pair plays ConferenceGamesPerPair games (half home, half away).
//   - Each team plays NonConferenceOpponents teams from the other conference once, rotating each season.
//   - Weeks 1..ThreeSlotWeeks have slots A-C. The last FourSlotWeeks weeks add a "D" slot reserved for divisional games.
type PHLScheduleConfig struct {
	SeasonIndex            int
	DivisionalGamesPerPair int
	ConferenceGamesPerPair int
	NonConferenceOpponents int
	ThreeSlotWeeks         int
	FourSlotWeeks          int
	MaxPartitionAttempts   int
	MinRoundGapPenalty     int // rounds within this distance that repeat a pairing are penalized when ordering
}

func DefaultPHLScheduleConfig(seasonID uint, maxAttempts int) PHLScheduleConfig {
	return PHLScheduleConfig{
		SeasonIndex:            int(seasonID),
		DivisionalGamesPerPair: 6,
		ConferenceGamesPerPair: 2,
		NonConferenceOpponents: 6,
		ThreeSlotWeeks:         16,
		FourSlotWeeks:          2,
		MaxPartitionAttempts:   maxAttempts,
		MinRoundGapPenalty:     2,
	}
}

const sameWeekPenalty = 1000

type phlMatchup struct {
	home, away int // indexes into the teams slice
}

type phlSlot struct {
	week     int
	slot     string
	divisonO bool // only divisional games (plus the unavoidable same-conference leftover) may be played
}

func GeneratePHLSeasonSchedule() {
	db := dbprovider.GetInstance().GetDB()
	ts := repository.FindTimestamp()
	phlTeams := repository.FindAllProTeams(repository.TeamClauses{LeagueID: "1"})
	schedule, err := GenerateProSeasonSchedule(phlTeams, ts.SeasonID, 2000)
	if err != nil {
		log.Println("Generation Failed: ", err)
		return
	}
	finalUpload := []structs.ProfessionalGame{}
	for _, game := range schedule {
		// Need to get 2027 to 2700. Right now 2027 becomes 2027000, which is not correct.
		seasonKey := util.GetWeekID(ts.SeasonID, uint(game.Week))
		game.AddWeekData(seasonKey, uint(game.Week), game.Slot)
		finalUpload = append(finalUpload, game.ProfessionalGame)
	}
	repository.CreatePHLGamesRecordsBatch(db, finalUpload, 100)
}

func GenerateProSeasonSchedule(
	teams []structs.ProfessionalTeam,
	seasonID uint,
	maxPartitionAttempts int,
) ([]ScheduledGame, error) {
	cfg := DefaultPHLScheduleConfig(seasonID, maxPartitionAttempts)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return generatePHLSchedule(rng, teams, seasonID, cfg)
}

func generatePHLSchedule(rng *rand.Rand, teams []structs.ProfessionalTeam, seasonID uint, cfg PHLScheduleConfig) ([]ScheduledGame, error) {
	n := len(teams)
	if n == 0 || n%2 != 0 {
		return nil, fmt.Errorf("need an even, non-zero number of teams, got %d", n)
	}
	slots := buildPHLSlots(cfg)
	totalRounds := len(slots)

	matchups, err := buildPHLMatchups(teams, cfg)
	if err != nil {
		return nil, err
	}
	if err := validatePHLMatchups(n, matchups, totalRounds); err != nil {
		return nil, err
	}

	divisionOf := make([]int, n)
	confOf := make([]int, n)
	for i, t := range teams {
		divisionOf[i] = int(t.DivisionID)
		confOf[i] = int(t.ConferenceID)
	}
	// In a divisional-only round, an odd-sized division leaves one team without a divisional
	// opponent; those leftover teams may only be paired with a same-conference team from another division.
	anyAllowed := func(i, j int) bool { return true }
	divisionalRound := func(i, j int) bool {
		return divisionOf[i] == divisionOf[j] || confOf[i] == confOf[j]
	}

	attempts := cfg.MaxPartitionAttempts
	if attempts < 1 {
		attempts = 1
	}
	var rounds [][]phlMatchup // indexed like slots
	ordered := false
	for attempt := 0; attempt < attempts && !ordered; attempt++ {
		rounds, err = partitionPHLRounds(rng, n, matchups, slots, anyAllowed, divisionalRound)
		if err != nil {
			continue
		}
		rounds, ordered = orderPHLRounds(rng, n, rounds, slots, cfg.MinRoundGapPenalty)
		if !ordered {
			err = fmt.Errorf("could not avoid the same pair playing twice in one week")
		}
	}
	if err != nil {
		return nil, fmt.Errorf("could not build schedule after %d attempts: %w", attempts, err)
	}

	var final []ScheduledGame
	for idx, s := range slots {
		for _, m := range rounds[idx] {
			g := ScheduledGame{ProfessionalGame: generateProGameRecord(teams[m.home], teams[m.away], seasonID)}
			g.ApplyWeekAndSlot(s.week, s.slot)
			final = append(final, g)
		}
	}
	return final, nil
}

// buildPHLSlots lists every round in chronological order.
func buildPHLSlots(cfg PHLScheduleConfig) []phlSlot {
	var slots []phlSlot
	for w := 1; w <= cfg.ThreeSlotWeeks; w++ {
		for _, s := range []string{"A", "B", "C"} {
			slots = append(slots, phlSlot{week: w, slot: s})
		}
	}
	for w := cfg.ThreeSlotWeeks + 1; w <= cfg.ThreeSlotWeeks+cfg.FourSlotWeeks; w++ {
		for _, s := range []string{"A", "B", "C"} {
			slots = append(slots, phlSlot{week: w, slot: s})
		}
		slots = append(slots, phlSlot{week: w, slot: "D", divisonO: true})
	}
	return slots
}

// buildPHLMatchups creates every game (with its host) that makes up the season.
func buildPHLMatchups(teams []structs.ProfessionalTeam, cfg PHLScheduleConfig) ([]phlMatchup, error) {
	if cfg.DivisionalGamesPerPair%2 != 0 || cfg.ConferenceGamesPerPair%2 != 0 {
		return nil, fmt.Errorf("divisional and conference games per pair must be even for an even home/away split")
	}
	confSet := map[int]bool{}
	for _, t := range teams {
		confSet[int(t.ConferenceID)] = true
	}
	if len(confSet) != 2 {
		return nil, fmt.Errorf("expected exactly 2 conferences, found %d", len(confSet))
	}
	confIDs := make([]int, 0, 2)
	for id := range confSet {
		confIDs = append(confIDs, id)
	}
	sort.Ints(confIDs)

	// position of each team within its conference, ordered by team ID
	confTeams := map[int][]int{}
	for i, t := range teams {
		confTeams[int(t.ConferenceID)] = append(confTeams[int(t.ConferenceID)], i)
	}
	pos := make([]int, len(teams))
	for _, idxs := range confTeams {
		sort.Slice(idxs, func(a, b int) bool { return teams[idxs[a]].ID < teams[idxs[b]].ID })
		for p, ti := range idxs {
			pos[ti] = p
		}
	}
	size := len(confTeams[confIDs[0]])
	if len(confTeams[confIDs[1]]) != size {
		return nil, fmt.Errorf("conferences must be the same size (%d vs %d)", size, len(confTeams[confIDs[1]]))
	}
	if cfg.NonConferenceOpponents > size {
		return nil, fmt.Errorf("cannot rotate %d non-conference opponents among %d teams", cfg.NonConferenceOpponents, size)
	}

	var out []phlMatchup
	addPair := func(a, b, games int) {
		for x := 0; x < games; x++ {
			if x%2 == 0 {
				out = append(out, phlMatchup{home: a, away: b})
			} else {
				out = append(out, phlMatchup{home: b, away: a})
			}
		}
	}

	// The window of offsets advances by NonConferenceOpponents each season. Offsets are
	// consecutive, so parity alternates and every team hosts exactly half of its block.
	start := (cfg.SeasonIndex * cfg.NonConferenceOpponents) % size
	if start < 0 {
		start += size
	}
	parity := cfg.SeasonIndex % 2
	if parity < 0 {
		parity += 2
	}

	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			a, b := teams[i], teams[j]
			switch {
			case a.ConferenceID == b.ConferenceID && a.DivisionID == b.DivisionID:
				addPair(i, j, cfg.DivisionalGamesPerPair)
			case a.ConferenceID == b.ConferenceID:
				addPair(i, j, cfg.ConferenceGamesPerPair)
			default:
				first, second := i, j // first is in the lower conference ID
				if int(a.ConferenceID) != confIDs[0] {
					first, second = j, i
				}
				k := ((pos[second]-pos[first])%size + size) % size
				if (k-start+size)%size >= cfg.NonConferenceOpponents {
					continue
				}
				if (k+parity)%2 == 0 {
					out = append(out, phlMatchup{home: first, away: second})
				} else {
					out = append(out, phlMatchup{home: second, away: first})
				}
			}
		}
	}
	return out, nil
}

func validatePHLMatchups(n int, matchups []phlMatchup, totalRounds int) error {
	games := make([]int, n)
	home := make([]int, n)
	for _, m := range matchups {
		games[m.home]++
		games[m.away]++
		home[m.home]++
	}
	for i := 0; i < n; i++ {
		if games[i] != totalRounds {
			return fmt.Errorf("team index %d has %d games, expected %d", i, games[i], totalRounds)
		}
		if home[i]*2 != totalRounds {
			return fmt.Errorf("team index %d has %d home games out of %d", i, home[i], games[i])
		}
	}
	return nil
}

// partitionPHLRounds splits the games into one perfect matching per slot. Divisional-only
// slots are filled first, then every round is a randomized perfect matching of the remaining
// games. Pairs that must appear in every remaining round are forced, and each choice is checked
// so that the leftover games still contain a perfect matching.
func partitionPHLRounds(
	rng *rand.Rand,
	n int,
	matchups []phlMatchup,
	slots []phlSlot,
	anyAllowed, divisionalAllowed func(i, j int) bool,
) ([][]phlMatchup, error) {
	pairGames := make(map[int][]phlMatchup)
	rem := make([][]int, n)
	for i := range rem {
		rem[i] = make([]int, n)
	}
	for _, m := range matchups {
		a, b := m.home, m.away
		if a > b {
			a, b = b, a
		}
		pairGames[a*n+b] = append(pairGames[a*n+b], m)
		rem[a][b]++
		rem[b][a]++
	}

	// process divisional-only slots first, then the rest; remember each slot's index
	var order []int
	for i, s := range slots {
		if s.divisonO {
			order = append(order, i)
		}
	}
	for i, s := range slots {
		if !s.divisonO {
			order = append(order, i)
		}
	}
	allowedFor := func(slotIdx int) func(i, j int) bool {
		if slots[slotIdx].divisonO {
			return divisionalAllowed
		}
		return anyAllowed
	}

	rounds := make([][]phlMatchup, len(slots))
	for step, slotIdx := range order {
		roundsLeft := len(order) - step
		var nextAllowed func(i, j int) bool
		if step+1 < len(order) {
			nextAllowed = allowedFor(order[step+1])
		}
		mate := pickPHLRound(rng, n, rem, roundsLeft, allowedFor(slotIdx), nextAllowed)
		if mate == nil {
			return nil, fmt.Errorf("round %d of %d: no valid matching", step+1, len(order))
		}
		var games []phlMatchup
		for i := 0; i < n; i++ {
			j := mate[i]
			if i > j {
				continue
			}
			key := i*n + j
			list := pairGames[key]
			games = append(games, list[0])
			pairGames[key] = list[1:]
			rem[i][j]--
			rem[j][i]--
		}
		rounds[slotIdx] = games
	}
	return rounds, nil
}

func pickPHLRound(
	rng *rand.Rand,
	n int,
	rem [][]int,
	roundsLeft int,
	allowed, nextAllowed func(i, j int) bool,
) []int {
	forced := make([]int, n)
	for i := range forced {
		forced[i] = -1
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j && rem[i][j] >= roundsLeft {
				if forced[i] != -1 || !allowed(i, j) {
					return nil
				}
				forced[i] = j
			}
		}
	}

	for attempt := 0; attempt < 60; attempt++ {
		adj := make([][]int, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j || rem[i][j] == 0 || !allowed(i, j) {
					continue
				}
				if (forced[i] != -1 && forced[i] != j) || (forced[j] != -1 && forced[j] != i) {
					continue
				}
				adj[i] = append(adj[i], j)
			}
			rng.Shuffle(len(adj[i]), func(a, b int) { adj[i][a], adj[i][b] = adj[i][b], adj[i][a] })
			// favor pairs with the most games left so they don't get stranded in the final rounds
			row := rem[i]
			sort.SliceStable(adj[i], func(a, b int) bool { return row[adj[i][a]] > row[adj[i][b]] })
		}
		mate := findPerfectMatching(rng, n, adj)
		if mate == nil {
			continue
		}
		if nextAllowed == nil {
			return mate
		}

		// lookahead: the games left after this round must still allow a perfect matching
		for i := 0; i < n; i++ {
			if i < mate[i] {
				rem[i][mate[i]]--
				rem[mate[i]][i]--
			}
		}
		next := make([][]int, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && rem[i][j] > 0 && nextAllowed(i, j) {
					next[i] = append(next[i], j)
				}
			}
		}
		ok := findPerfectMatching(rng, n, next) != nil
		for i := 0; i < n; i++ {
			if i < mate[i] {
				rem[i][mate[i]]++
				rem[mate[i]][i]++
			}
		}
		if ok {
			return mate
		}
	}
	return nil
}

// orderPHLRounds shuffles non-"D" rounds among their slots and hill-climbs. A pair meeting twice in
// the same week is a hard rule (reported via the bool); nearby repeats are only discouraged.
func orderPHLRounds(rng *rand.Rand, n int, rounds [][]phlMatchup, slots []phlSlot, gap int) ([][]phlMatchup, bool) {
	var free []int
	for i, s := range slots {
		if !s.divisonO {
			free = append(free, i)
		}
	}
	rng.Shuffle(len(free), func(a, b int) {
		rounds[free[a]], rounds[free[b]] = rounds[free[b]], rounds[free[a]]
	})

	pairSets := func() [][]bool {
		sets := make([][]bool, len(rounds))
		for r, games := range rounds {
			sets[r] = make([]bool, n*n)
			for _, g := range games {
				a, b := g.home, g.away
				if a > b {
					a, b = b, a
				}
				sets[r][a*n+b] = true
			}
		}
		return sets
	}
	cost := func(sets [][]bool) int {
		total := 0
		for r := range rounds {
			for d := 1; r+d < len(rounds); d++ {
				sameWeek := slots[r].week == slots[r+d].week
				if d > gap && !sameWeek {
					break
				}
				for _, g := range rounds[r] {
					a, b := g.home, g.away
					if a > b {
						a, b = b, a
					}
					if sets[r+d][a*n+b] {
						if sameWeek {
							total += sameWeekPenalty
						} else {
							total += gap - d + 1
						}
					}
				}
			}
		}
		return total
	}

	if len(free) < 2 {
		return rounds, cost(pairSets()) < sameWeekPenalty
	}
	best := cost(pairSets())
	for iter := 0; iter < 20000 && best > 0; iter++ {
		a, b := free[rng.Intn(len(free))], free[rng.Intn(len(free))]
		if a == b {
			continue
		}
		rounds[a], rounds[b] = rounds[b], rounds[a]
		if c := cost(pairSets()); c <= best {
			best = c
		} else {
			rounds[a], rounds[b] = rounds[b], rounds[a]
		}
	}
	return rounds, best < sameWeekPenalty
}

// findPerfectMatching runs Edmonds' blossom algorithm. It returns nil when no perfect matching exists.
func findPerfectMatching(rng *rand.Rand, n int, adj [][]int) []int {
	match := make([]int, n)
	p := make([]int, n)
	base := make([]int, n)
	used := make([]bool, n)
	inBlossom := make([]bool, n)
	for i := range match {
		match[i] = -1
	}

	lca := func(a, b int) int {
		seen := make([]bool, n)
		for {
			a = base[a]
			seen[a] = true
			if match[a] == -1 {
				break
			}
			a = p[match[a]]
		}
		for {
			b = base[b]
			if seen[b] {
				return b
			}
			b = p[match[b]]
		}
	}
	markPath := func(v, b, child int) {
		for base[v] != b {
			inBlossom[base[v]] = true
			inBlossom[base[match[v]]] = true
			p[v] = child
			child = match[v]
			v = p[match[v]]
		}
	}
	findPath := func(root int) int {
		for i := 0; i < n; i++ {
			used[i] = false
			p[i] = -1
			base[i] = i
		}
		used[root] = true
		queue := []int{root}
		for qh := 0; qh < len(queue); qh++ {
			v := queue[qh]
			for _, to := range adj[v] {
				if base[v] == base[to] || match[v] == to {
					continue
				}
				if to == root || (match[to] != -1 && p[match[to]] != -1) {
					cur := lca(v, to)
					for i := range inBlossom {
						inBlossom[i] = false
					}
					markPath(v, cur, to)
					markPath(to, cur, v)
					for i := 0; i < n; i++ {
						if inBlossom[base[i]] {
							base[i] = cur
							if !used[i] {
								used[i] = true
								queue = append(queue, i)
							}
						}
					}
				} else if p[to] == -1 {
					p[to] = v
					if match[to] == -1 {
						return to
					}
					next := match[to]
					used[next] = true
					queue = append(queue, next)
				}
			}
		}
		return -1
	}

	for _, v := range rng.Perm(n) {
		if match[v] != -1 {
			continue
		}
		u := findPath(v)
		if u == -1 {
			return nil
		}
		for u != -1 {
			pv := p[u]
			next := match[pv]
			match[u] = pv
			match[pv] = u
			u = next
		}
	}
	return match
}

func generateProGameRecord(teamA structs.ProfessionalTeam, teamB structs.ProfessionalTeam, seasonID uint) structs.ProfessionalGame {
	return structs.ProfessionalGame{
		BaseGame: structs.BaseGame{
			SeasonID:     seasonID,
			HomeTeamID:   teamA.ID,
			HomeTeam:     teamA.Abbreviation,
			AwayTeamID:   teamB.ID,
			AwayTeam:     teamB.Abbreviation,
			ArenaID:      uint(teamA.ArenaID),
			Arena:        teamA.Arena,
			City:         teamA.City,
			State:        teamA.State,
			Country:      teamA.Country,
			IsConference: teamA.ConferenceID == teamB.ConferenceID,
			LeagueID:     uint(teamA.LeagueID),
		},
		IsDivisional: teamA.DivisionID == teamB.DivisionID,
	}
}
