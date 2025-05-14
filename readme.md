# SimHockey API

## Project Overview

SimHockey is a Golang-based HTTP API designed to simulate hockey leagues (College Hockey League \[CHL] and Professional Hockey League \[PHL]). It provides endpoints for league data management, free agency workflows, trade synchronization, statistics export, and real-time updates via WebSockets and Discord integrations. Connects to (SimSN Interface 2.0)[https://github.com/CalebRose/simsn-interface-v2].

## Repository Structure

```
.
├── main.go               # application entrypoint and router setup
├── .env                  # environment variables (loaded via godotenv)
├── .gitignore            # files to ignore in Git
├── controllers/          # HTTP handler functions organized by domain
├── dbprovider/           # database initialization and connection pooling (SSH support)
├── middleware/           # HTTP middleware (CORS, logging, etc.)
├── engine/               # core simulation logic (game events, penalties, scoring, etc.)
├── managers/             # business logic orchestrating league operations (schedules, drafting, trading, progression)
├── repository/           # data access layer for reading/writing domain models
├── structs/              # domain model definitions (teams, players, trades, timestamps)
├── ws/                   # WebSocket server and broadcast logic for timestamp updates
├── ts/                   # TypeScript model definitions for testing/frontend
├── data/                 # seed data and JSON assets (face data, name lists)
├── _util/                # helper utilities (CSV import/export, common helpers)
├── test_results/         # example outputs and test artifacts
└── readme.md             # this file
```

### Highlights of Key Directories

- **controllers/**: Contains files like `AdminController.go`, `BootstrapController.go`, `ExportController.go`, `FreeAgencyController.go`, `TradeController.go`, `WebsocketController.go`, etc., each grouping related HTTP handlers.
- **engine/**: Implements the core game engine (`events.go`, `penalties.go`, `basegame.go`, etc.), defining how matches, scoring, and penalties are simulated.
- **managers/**: High-level orchestration for operations such as drafting (`DraftManagers.go`), free agency (`FreeAgencyManager.go`), trading (`TradeManager.go`), progression of players (`ProgressionManager.go`), and more.
- **dbprovider/**: Sets up the database connection, supporting SSH tunnels and environment-based configuration.
- **repository/**: CRUD operations against the database models defined in `structs/`.
- **ws/**: WebSocket entrypoint (`WebsocketController.go`) and broadcasting logic (`broadcast.go`) for pushing timestamp updates to clients.

## API Endpoints

### Health Check

- **GET** `/health`
  Returns basic service health metadata (version, release ID).

### Admin

- **GET** `/api/admin/generate/ts/models/` — Generate TypeScript model definitions from Go structs
- **GET** `/api/admin/run/fa/sync/` — Trigger a test free agency synchronization

### Bootstrap

- **GET** `/api/bootstrap/{collegeID}/{proID}` — Seed initial hockey data for a college and pro team

### Export

- **GET** `/api/export/pro/players/all` — Export all professional player data
- **GET** `/api/export/college/players/all` — Export all college player data
- **GET** `/api/export/stats/chl/{seasonID}/{weekID}/{gameType}` — Export CHL stats page content for a given season/week
- **GET** `/api/export/stats/phl/{seasonID}/{weekID}/{gameType}` — Export PHL stats page content for a given season/week

### Free Agency

- **POST** `/api/phl/freeagency/create/offer` — Create a free agency offer for a PHL player

### Trades

- **GET** `/api/trades/admin/accept/sync/{proposalID}` — Sync an accepted trade proposal
- **GET** `/api/trades/admin/veto/sync/{proposalID}` — Veto an accepted trade proposal
- **GET** `/api/trades/admin/cleanup` — Clean up rejected trade records

### Discord (DS)

- **GET** `/api/ds/chl/team/{teamID}/` — Fetch CHL team info for Discord bots
- **GET** `/api/ds/phl/team/{teamID}/` — Fetch PHL team info for Discord bots
- **GET** `/api/ds/chl/player/id/{id}` — Fetch CHL player by ID for Discord bots
- **GET** `/api/ds/chl/player/name/{firstName}/{lastName}/{abbr}` — Fetch CHL player by name for Discord bots
- **GET** `/api/ds/phl/player/id/{id}` — Fetch PHL player by ID for Discord bots
- **GET** `/api/ds/phl/player/name/{firstName}/{lastName}/{abbr}` — Fetch PHL player by name for Discord bots
- **GET** `/api/ds/chl/flex/{teamOneID}/{teamTwoID}/` — Compare CHL teams (flex comparison) via Discord
- **GET** `/api/ds/phl/flex/{teamOneID}/{teamTwoID}/` — Compare PHL teams (flex comparison) via Discord
- **GET** `/api/ds/chl/assign/discord/{teamID}/{discordID}` — Assign a Discord ID to a CHL team
- **GET** `/api/ds/phl/assign/discord/{teamID}/{discordID}` — Assign a Discord ID to a PHL team
- **GET** `/api/ds/chl/croots/class/{teamID}/` — Fetch CHL recruiting class by team ID for Discord
- **GET** `/api/ds/chl/croot/{id}` — Fetch an individual CHL recruit for Discord

### WebSocket

- **GET** `/ws` — Upgrade HTTP connection to WebSocket for real-time timestamp broadcasts

### Cronjobs

> _Details on cron jobs will be added at a later date._

## Getting Started

> _Instructions for cloning, building, and running the API will be added soon._
