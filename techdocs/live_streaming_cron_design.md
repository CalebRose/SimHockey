# Live Game Streaming — Technical Design Document

**Language:** Go  
**Status:** Proposed  
**Scope:** CHL & PHL leagues

---

## Overview

This document describes the design for a background cron job that streams up to 8 simultaneous games in real time, as if a user had switched on a TV channel mid-broadcast. The API is the only entity that reads from or writes to the relevant Firebase collections. Clients read game metadata from Firebase once on page load, then source all play-by-play data from the API — keeping Firebase reads minimal and eliminating client-side writes entirely.

---

## Goals

- Stream exactly 8 games concurrently at all times (or as many as are available if fewer than 8 remain).
- Compute a deterministic `StreamStartTime` and `StreamEndTime` per game from the play-by-play data so any client joining mid-stream can calculate the current play without polling.
- When a game ends, dequeue the next unplayed game and begin streaming it, maintaining the 8-game ceiling.
- Minimize Firebase reads and writes. Firebase stores only a lightweight registry of active streams; play-by-play is served exclusively from the API.
- Ensure no client ever writes to the live stream collections.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Cron Job (Go goroutine, fires every N seconds)         │
│                                                         │
│  StreamScheduler                                        │
│  ├── ActiveSlots [8]GameStream                          │
│  ├── GameQueue  []PendingGame                           │
│  └── tick() → advance plays, rotate completed games     │
└────────────────────────┬────────────────────────────────┘
                         │ writes (batch, infrequent)
                         ▼
                    Firebase Firestore
                    ┌─────────────────────────────┐
                    │  live_chl_games / live_phl_games    │
                    │  (one doc per active game,   │
                    │   ~8 docs max at any time)   │
                    └─────────────────────────────┘
                         │ read once on page load
                         ▼
                      Client (browser / app)
                         │
                         │ GET /api/live-plays/:gameID
                         ▼
                      API (Go handler)
                      └── returns full PbP slice for game
                           from DB, no Firebase read
```

---

## Firebase Schema Changes

### `live_chl_games` / `live_phl_games` (one document per active stream)

This collection is already defined (`LiveGameRecord`). Two fields need to be added:

```go
// In firebase/types.go — extend LiveGameRecord
type LiveGameRecord struct {
    GameID          uint      `firestore:"GameID"`
    HomeTeamID      uint      `firestore:"HomeTeamID"`
    AwayTeamID      uint      `firestore:"AwayTeamID"`
    HomeTeam        string    `firestore:"HomeTeam"`
    AwayTeam        string    `firestore:"AwayTeam"`
    League          string    `firestore:"League"`
    StreamStartTime time.Time `firestore:"StreamStartTime"`
    StreamEndTime   time.Time `firestore:"StreamEndTime"`   // NEW
    TotalPlays      int       `firestore:"TotalPlays"`      // NEW
    IsRevealed      bool      `firestore:"IsRevealed"`
}
```

`StreamEndTime` = `StreamStartTime` + sum of all `SecondsConsumed` across the game's play-by-play.  
`TotalPlays` allows the client to validate the index it calculates without fetching the full PbP list.

No other Firebase collections are touched by this feature.

---

## Server-Side Design

### Data Structures

```go
// managers/StreamScheduler.go (new file)

// GameStream represents one active streaming slot.
type GameStream struct {
    GameID          uint
    League          string          // "chl" or "phl"
    StartTime       time.Time
    EndTime         time.Time
    TotalSeconds    int             // sum of SecondsConsumed across all plays
    IsComplete      bool
}

// StreamScheduler manages the 8 concurrent slots and the waiting queue.
type StreamScheduler struct {
    mu          sync.Mutex
    ActiveSlots [8]*GameStream      // nil slot = available
    Queue       []PendingGame       // ordered list of games awaiting a slot
    League      string
}

// PendingGame is a lightweight descriptor for a game waiting to stream.
type PendingGame struct {
    GameID       uint
    TotalSeconds int
}
```

### Computing StreamStartTime and StreamEndTime

When a game is loaded into a slot:

```go
func computeStreamTimes(plays []structs.CollegePlayByPlay) (start, end time.Time, totalSecs int) {
    for _, p := range plays {
        totalSecs += int(p.SecondsConsumed)
    }
    start = time.Now().UTC()
    end = start.Add(time.Duration(totalSecs) * time.Second)
    return
}
```

This is deterministic: the game clock runs at real-time (1 simulated second = 1 wall-clock second). Because every play has a concrete `SecondsConsumed` value in the `PbP` struct, the endpoint time is exact.

### Cron Tick Logic

```go
// Called by the cron job on every tick (recommended: every 5 seconds).
func (s *StreamScheduler) Tick(ctx context.Context) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UTC()
    slotsFreed := 0

    // 1. Mark completed games and free their slots.
    for i, slot := range s.ActiveSlots {
        if slot != nil && !slot.IsComplete && now.After(slot.EndTime) {
            slot.IsComplete = true
            s.ActiveSlots[i] = nil
            slotsFreed++
            // Batch-delete (or mark IsRevealed) in Firebase.
            go firebase.SetGameRevealed(ctx, slot.GameID, slot.League)
        }
    }

    // 2. Fill freed slots from the queue.
    for i, slot := range s.ActiveSlots {
        if slot != nil {
            continue
        }
        next, ok := s.dequeue()
        if !ok {
            break
        }
        plays := loadPlays(next.GameID, s.League)
        start, end, totalSecs := computeStreamTimes(plays)
        s.ActiveSlots[i] = &GameStream{
            GameID:       next.GameID,
            League:       s.League,
            StartTime:    start,
            EndTime:      end,
            TotalSeconds: totalSecs,
        }
        // Write the new slot to Firebase.
        go firebase.UploadLiveGame(ctx, buildLiveGameRecord(s.ActiveSlots[i], plays), s.League)
    }
}
```

### Initializing the Queue

On cron startup (or after RunGames completes), load all complete, unrevealed games for the current week and sort them by GameID (or any stable ordering). Populate `Queue` with the full list, then call `Tick` immediately to fill the initial 8 slots.

```go
func (s *StreamScheduler) InitQueue(weekID, seasonID, gameDay string, isPreseason bool) {
    games := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
    s.mu.Lock()
    defer s.mu.Unlock()
    for _, g := range games {
        if !g.GameComplete || g.IsRevealed {
            continue
        }
        secs := loadTotalSeconds(g.ID, s.League)
        s.Queue = append(s.Queue, PendingGame{GameID: g.ID, TotalSeconds: secs})
    }
}
```

`loadTotalSeconds` queries the PbP table once per game and sums `SecondsConsumed`. This is a single DB read per game, done once at queue-init time.

### Cron Registration

In `managers/SchedulerManager.go` (or wherever your existing crons live):

```go
func StartLiveStreamingCron(league string) {
    scheduler := &StreamScheduler{League: league}
    ts := GetTimestamp()
    gameDay := ts.GetGameDay()
    scheduler.InitQueue(
        strconv.Itoa(int(ts.WeekID)),
        strconv.Itoa(int(ts.SeasonID)),
        gameDay,
        ts.IsPreseason,
    )

    ctx := context.Background()
    scheduler.Tick(ctx) // fill initial slots immediately

    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            scheduler.Tick(ctx)
        }
    }()
}
```

Call `StartLiveStreamingCron("chl")` and `StartLiveStreamingCron("phl")` from `CronController.go` (or your bootstrap path) after the game run is complete.

---

## Client-Side Design

### Step 1 — Page Load: Fetch Active Games from Firebase (one read)

The client reads from `live_chl_games` (or `live_phl_games`) to get the registry of currently active streams. This is the **only** Firebase read for this feature.

Each document gives the client:

- `GameID`, `HomeTeam`, `AwayTeam` — for display
- `StreamStartTime`, `StreamEndTime` — for computing the current play
- `TotalPlays` — for bounds checking

### Step 2 — Fetch Play-by-Play from the API (not Firebase)

The client makes a single GET request per game:

```
GET /api/live-plays/:league/:gameID
```

This returns the complete ordered `[]PlayByPlayResponse` slice for that game from the database. No Firebase read occurs. The response is cacheable for the duration of the stream because the play list is immutable once a game is complete.

#### OPTIONALLY - Check for the fetch data call for all Bulk play by play data (GameManager.go)

Check for GetBulkPlayByPlayData in SimHockey/Managers/GameManager.go for the fetch. Once all play by play data has been retrieved,

Once done fetching for this call, do NOT place the play by play data in firebase. DO NOT. Move on to step 3

### Step 3 — Compute the Current Play (client-side math, no polling)

Given the full PbP list and the two timestamps from Firebase, the client determines which play is "on screen right now":

```typescript
function getCurrentPlayIndex(
  plays: PlayByPlayResponse[],
  streamStartTime: Date,
): number {
  const elapsedSeconds = (Date.now() - streamStartTime.getTime()) / 1000;
  let accumulated = 0;
  for (let i = 0; i < plays.length; i++) {
    accumulated += plays[i].secondsConsumed;
    if (accumulated >= elapsedSeconds) {
      return i;
    }
  }
  return plays.length - 1; // game is over
}
```

The client advances the displayed play using a local `setInterval` that increments by `secondsConsumed` for each play — no network request needed between plays. This is the "turn on the TV mid-broadcast" experience: the user always joins at whatever point in the game wall-clock time dictates.

### Step 4 — Refreshing the Active Game List

After a game ends (`Date.now() > streamEndTime`), the client re-reads the Firebase games collection to discover which game replaced it. This is the only subsequent Firebase read, and it happens at most once per completed game (roughly every few minutes per slot).

---

## API Endpoint

### `GET /api/live-plays/:league/:gameID`

Handler location: `controllers/LiveScoreboardController.go`

```go
func GetLivePlays(c *gin.Context) {
    league := c.Param("league")   // "chl" or "phl"
    gameID := c.Param("gameID")

    var response []structs.PlayByPlayResponse
    if league == "chl" {
        plays := managers.GetCHLPlayByPlaysByGameID(gameID)
        // reuse existing GenerateCHLPlayByPlayResponse, isStream=true
        response = managers.GenerateCHLPlayByPlayResponse(plays, ...)
    } else {
        plays := managers.GetPHLPlayByPlaysByGameID(gameID)
        response = managers.GeneratePHLPlayByPlayResponse(plays, ...)
    }

    c.JSON(http.StatusOK, response)
}
```

This endpoint is **read-only**, **stateless**, and touches **no Firebase** resources.

---

## Firebase Read/Write Budget

| Operation                                  | Who      | When                 | Count                    |
| ------------------------------------------ | -------- | -------------------- | ------------------------ |
| Read `live_chl_games`                      | Client   | Page load            | 1 per session            |
| Read `live_chl_games`                      | Client   | After a game ends    | 1 per slot rotation      |
| Write `live_chl_games` (set)               | API cron | New game enters slot | 1 per rotation           |
| Write `live_chl_games` (update IsRevealed) | API cron | Game completes       | 1 per rotation           |
| Delete stale records                       | API cron | On next RunGames     | 1 batch at session start |

With 8 slots and typical game durations of ~45–60 minutes, slot rotations happen at most once every ~45 minutes per slot. Daily Firebase write volume from this feature is in the dozens, not thousands.

---

## Comparison to Current Approach

| Concern                | Current approach         | New approach                       |
| ---------------------- | ------------------------ | ---------------------------------- |
| Who writes to Firebase | API + Client             | API only                           |
| Play-by-play source    | Firebase                 | API (DB)                           |
| Reads per user session | Many (live updates)      | 1–2 (at load + per rotation)       |
| Payload size per read  | Large (all plays in doc) | Lightweight (8 game metadata docs) |
| Client writes          | Present                  | Eliminated                         |

The existing `UploadLivePlays` / `LivePlaysRecord` pattern (storing the full play list in Firestore) is retired. Plays come from the API; Firebase is purely a scheduling registry.

---

## Implementation Checklist

- [ ] Extend `LiveGameRecord` with `StreamEndTime` and `TotalPlays` fields in `firebase/types.go`
- [ ] Add `firebase.UploadLiveGame` (single-game variant) to `firebase/live_service.go`
- [ ] Create `managers/StreamScheduler.go` with `StreamScheduler`, `GameStream`, `PendingGame`
- [ ] Implement `InitQueue`, `Tick`, `computeStreamTimes`, `loadTotalSeconds`
- [ ] Wire `StartLiveStreamingCron` into `main.go` or bootstrap path
- [ ] Add `GET /api/live-plays/:league/:gameID` route and handler
- [ ] Update `LiveScoreboardController.go` to expose the new endpoint
- [ ] Remove any existing client-side Firebase write paths for the live collections
- [ ] Update client to use `StreamStartTime`/`StreamEndTime` math instead of polling

---

## Open Questions

- **Tick interval:** 5 seconds is conservative. Because `EndTime` is computed deterministically, the tick only needs to run frequently enough to catch game completions within a reasonable window. A 30-second tick is likely sufficient.

#### Answer

-If 5 seconds is conservative, we can interval faster if needed.

- **Game ordering:** Should the queue prioritize user-coached matchups (matching the existing `streamType` logic in `GetCHLPlayByPlayStreamData`)? If so, partition the queue accordingly before filling slots.

#### Answer

Yes, we should prioritize user-coached matchups when possible.

- **PHL vs CHL schedulers:** Run as two independent `StreamScheduler` instances, or a single scheduler that manages both leagues. Two independent instances is simpler and avoids cross-league slot contention.

#### Answer

Yes, there are two cron jobs setup in CronController.go taht we can use for independently setting up the jobs & the stream.

- **Error handling on slot fill:** If `loadPlays` returns an empty slice (e.g., PbP not yet persisted), skip that game and try the next in the queue rather than occupying a slot with a broken stream.

#### Answer

Yes, please do this.
