package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CalebRose/SimHockey/controllers"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/middleware"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/CalebRose/SimHockey/ws"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/checks/sendgrid"
	"github.com/robfig/cron/v3"
)

func InitialMigration() {
	initiate := dbprovider.GetInstance().InitDatabase()
	if !initiate {
		log.Println("Initiate pool failure... Ending application")
		os.Exit(1)
	}
}

func monitorDBForUpdates() {
	var ts structs.Timestamp
	for {
		currentTS := controllers.GetUpdatedTimestamp()
		if currentTS.UpdatedAt.After(ts.UpdatedAt) {
			ts = currentTS
			err := ws.BroadcastTSUpdate(ts)
			if err != nil {
				log.Printf("Error broadcasting timestamp: %v", err)
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func handleRequests() http.Handler {
	myRouter := mux.NewRouter().StrictSlash(true)

	// Handler & Middleware
	loadEnvs()
	origins := os.Getenv("ORIGIN_ALLOWED")
	originsOk := handlers.AllowedOrigins([]string{origins})
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Access-Control-Allow-Origin"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "HEAD"})
	apiRouter := myRouter.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.GzipMiddleware)

	// Health Controls
	HealthCheck := health.New(
		health.Health{
			Version:   "1",
			ReleaseID: "0.0.7-SNAPSHOT",
		},
		sendgrid.Health(),
	)
	myRouter.HandleFunc("/health", HealthCheck.Handler).Methods("GET")

	// Admin
	apiRouter.HandleFunc("/admin/generate/ts/models/", controllers.CreateTSModelsFile).Methods("GET")
	// apiRouter.HandleFunc("/admin/ai/generate/college/lineups/", controllers.RunAICollegeLineups).Methods("GET")
	// apiRouter.HandleFunc("/admin/ai/generate/gameplans/", controllers.CreateGameplans).Methods("GET")
	// apiRouter.HandleFunc("/admin/ai/generate/pro/lineups/", controllers.RunAIProLineups).Methods("GET")
	// apiRouter.HandleFunc("/admin/test/engine/", controllers.TestEngine).Methods("GET")
	// apiRouter.HandleFunc("/admin/show/results/", controllers.ShowGameResults).Methods("GET")
	// apiRouter.HandleFunc("/admin/assign/ranks/", controllers.AssignAllRecruitRanks).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/test/college/teams/", controllers.GenerateCollegeTeams).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/init/college/rosters/", controllers.GenerateInitialRosters).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/college/recruits/", controllers.GenerateCroots).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/phl/schedule/", controllers.GeneratePHLSchedule).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/chl/schedule/", controllers.GenerateCHLSchedule).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/pre/schedule/", controllers.GeneratePreseasonGames).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/test/pro/rosters/", controllers.GenerateProTestData).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/capsheets/", controllers.GenerateCapsheets).Methods("GET")
	// apiRouter.HandleFunc("/admin/generate/fa/preferences/", controllers.AddFAPreferences).Methods("GET")
	// apiRouter.HandleFunc("/admin/run/fa/sync/", controllers.TestFASync).Methods("GET")
	// apiRouter.HandleFunc("/admin/run/fa/sync/", controllers.TestFAOffers).Methods("GET")

	// Bootstrap
	apiRouter.HandleFunc("/bootstrap/{collegeID}/{proID}", controllers.BootstrapHockeyData).Methods("GET")
	apiRouter.HandleFunc("/bootstrap/teams/", controllers.BootstrapTeamData).Methods("GET")

	// Exports
	apiRouter.HandleFunc("/export/pro/players/all", controllers.ExportAllProPlayers).Methods("GET")
	apiRouter.HandleFunc("/export/college/players/all", controllers.ExportAllCollegePlayers).Methods("GET")
	apiRouter.HandleFunc("/export/college/recruits/all", controllers.ExportCHLRecruits).Methods("GET")
	apiRouter.HandleFunc("/export/college/roster/{teamID}", controllers.ExportCollegeRoster).Methods("GET")
	apiRouter.HandleFunc("/export/pro/roster/{teamID}", controllers.ExportProRoster).Methods("GET")
	apiRouter.HandleFunc("/export/stats/chl/{seasonID}/{weekID}/{viewType}/{gameType}", controllers.ExportCHLStatsPageContentForSeason).Methods("GET")
	apiRouter.HandleFunc("/export/stats/phl/{seasonID}/{weekID}/{viewType}/{gameType}", controllers.ExportProStatsPageContent).Methods("GET")

	// Free Agency
	apiRouter.HandleFunc("/phl/freeagency/create/offer", controllers.CreateFreeAgencyOffer).Methods("POST")
	apiRouter.HandleFunc("/phl/freeagency/cancel/offer", controllers.CancelFreeAgencyOffer).Methods("POST")
	apiRouter.HandleFunc("/phl/waiverwire/create/offer", controllers.CreateWaiverWireOffer).Methods("POST")
	apiRouter.HandleFunc("/phl/waiverwire/cancel/offer", controllers.CancelWaiverWireOffer).Methods("POST")

	// Games
	apiRouter.HandleFunc("/games/result/chl/{gameID}", controllers.GetCollegeGameResultsByGameID).Methods("GET")
	apiRouter.HandleFunc("/games/result/phl/{gameID}", controllers.GetProGameResultsByGameID).Methods("GET")
	apiRouter.HandleFunc("/games/export/results/{seasonID}/{weekID}/{timeslot}", controllers.ExportHCKGameResults).Methods("GET")

	// Imports
	// apiRouter.HandleFunc("/admin/import/pro/rosters/", controllers.ImportProRosters).Methods("GET")
	// apiRouter.HandleFunc("/admin/import/chl/team/profiles/", controllers.ImportTeamProfileRecords).Methods("GET")

	// Migrations
	// apiRouter.HandleFunc("/migrate/faces", controllers.MigrateFaceData).Methods("GET")

	// Poll Controls
	apiRouter.HandleFunc("/college/poll/create/", controllers.CreatePollSubmission).Methods("POST")
	apiRouter.HandleFunc("/college/poll/sync", controllers.SyncCollegePoll).Methods("GET")

	// Requests
	// apiRouter.HandleFunc("/admin/import/pro/teams/", controllers.GenerateProTeams).Methods("GET")
	apiRouter.HandleFunc("/admin/requests/hck/", controllers.GetAllHCKRequests).Methods("GET")
	apiRouter.HandleFunc("/requests/view/chl/{teamID}", controllers.ViewCHLTeamUponRequest).Methods("GET")
	apiRouter.HandleFunc("/requests/view/phl/{teamID}", controllers.ViewPHLTeamUponRequest).Methods("GET")
	apiRouter.HandleFunc("/chl/requests/approve", controllers.ApproveCHLTeamRequest).Methods("POST")
	apiRouter.HandleFunc("/phl/requests/approve", controllers.ApprovePHLTeamRequest).Methods("POST")
	apiRouter.HandleFunc("/chl/requests/create", controllers.CreateCHLTeamRequest).Methods("POST")
	apiRouter.HandleFunc("/phl/requests/create", controllers.CreatePHLTeamRequest).Methods("POST")
	apiRouter.HandleFunc("/chl/requests/reject", controllers.RejectCHLTeamRequest).Methods("POST")
	apiRouter.HandleFunc("/phl/requests/reject", controllers.RejectPHLTeamRequest).Methods("POST")

	// Recruiting
	apiRouter.HandleFunc("/recruiting/add/recruit/", controllers.CreateRecruitPlayerProfile).Methods("POST")
	apiRouter.HandleFunc("/recruiting/remove/recruit/", controllers.RemoveRecruitFromBoard).Methods("POST")
	apiRouter.HandleFunc("/recruiting/toggle/scholarship/", controllers.SendScholarshipToRecruit).Methods("POST")
	apiRouter.HandleFunc("/recruiting/scout/attribute/", controllers.ScoutAttribute).Methods("POST")
	apiRouter.HandleFunc("/recruiting/save/board/", controllers.SaveRecruitingBoard).Methods("POST")
	apiRouter.HandleFunc("/recruiting/save/ai/", controllers.ToggleAIBehavior).Methods("POST")

	// Roster Page
	apiRouter.HandleFunc("/chl/roster/cut/{playerID}", controllers.CutCHLPlayerFromRoster).Methods("GET")
	apiRouter.HandleFunc("/chl/roster/redshirt/{playerID}", controllers.RedshirtCHLPlayer).Methods("GET")
	apiRouter.HandleFunc("/chl/roster/promise/{playerID}", controllers.PromiseCHLPlayer).Methods("POST")
	apiRouter.HandleFunc("/phl/roster/cut/{playerID}", controllers.CutPHLPlayerFromRoster).Methods("GET")
	apiRouter.HandleFunc("/phl/roster/affiliate/{playerID}", controllers.SendPHLPlayerToAffiliate).Methods("GET")
	apiRouter.HandleFunc("/phl/roster/tradeblock/{playerID}", controllers.SendPHLPlayerToTradeBlock).Methods("GET")
	apiRouter.HandleFunc("/phl/roster/extend/{playerID}", controllers.ExtendPHLPlayer).Methods("POST")

	// Strategy
	apiRouter.HandleFunc("/chl/strategy/update", controllers.SaveCHLLineups).Methods("POST")
	apiRouter.HandleFunc("/chl/gameplan/update", controllers.SaveCHLGameplan).Methods("POST")
	apiRouter.HandleFunc("/phl/strategy/update", controllers.SavePHLLineups).Methods("POST")
	apiRouter.HandleFunc("/phl/gameplan/update", controllers.SavePHLGameplan).Methods("POST")

	// Stats
	apiRouter.HandleFunc("/statistics/interface/chl/{seasonID}/{weekID}/{viewType}/{gameType}", controllers.GetCHLStatsPageContentForSeason).Methods("GET")
	apiRouter.HandleFunc("/statistics/interface/phl/{seasonID}/{weekID}/{viewType}/{gameType}", controllers.GetProStatsPageContent).Methods("GET")

	// Teams
	apiRouter.HandleFunc("/chl/teams/remove/{teamID}", controllers.RemoveUserFromCollegeTeam).Methods("GET")
	apiRouter.HandleFunc("/phl/teams/remove/user", controllers.RemoveUserFromProTeam).Methods("POST")

	// Trades
	apiRouter.HandleFunc("/trades/phl/preferences/update", controllers.UpdateTradePreferences).Methods("POST")
	apiRouter.HandleFunc("/trades/phl/create/proposal", controllers.CreateTradeProposal).Methods("POST")
	apiRouter.HandleFunc("/trades/phl/proposal/accept/{proposalID}", controllers.AcceptTradeOffer).Methods("GET")
	apiRouter.HandleFunc("/trades/phl/proposal/reject/{proposalID}", controllers.RejectTradeOffer).Methods("GET")
	apiRouter.HandleFunc("/trades/phl/proposal/cancel/{proposalID}", controllers.CancelTradeOffer).Methods("GET")
	apiRouter.HandleFunc("/trades/admin/accept/sync/{proposalID}", controllers.SyncAcceptedTrade).Methods("GET")
	apiRouter.HandleFunc("/trades/admin/veto/sync/{proposalID}", controllers.VetoAcceptedTrade).Methods("GET")
	apiRouter.HandleFunc("/trades/admin/cleanup", controllers.CleanUpRejectedTrades).Methods("GET")

	// Discord
	apiRouter.HandleFunc("/ds/chl/team/{teamID}/", controllers.GetCHLTeamByTeamIDForDiscord).Methods("GET")
	apiRouter.HandleFunc("/ds/phl/team/{teamID}/", controllers.GetPHLTeamByTeamIDForDiscord).Methods("GET")
	apiRouter.HandleFunc("/ds/chl/player/id/{id}", controllers.GetCHLPlayerForDiscord).Methods("GET")
	apiRouter.HandleFunc("/ds/chl/player/name/{firstName}/{lastName}/{abbr}", controllers.GetCHLPlayerByName).Methods("GET")
	apiRouter.HandleFunc("/ds/phl/player/id/{id}", controllers.GetPHLPlayerForDiscord).Methods("GET")
	apiRouter.HandleFunc("/ds/phl/player/name/{firstName}/{lastName}/{abbr}", controllers.GetPHLPlayerByName).Methods("GET")
	apiRouter.HandleFunc("/ds/chl/flex/{teamOneID}/{teamTwoID}/", controllers.CompareCHLTeams).Methods("GET")
	apiRouter.HandleFunc("/ds/phl/flex/{teamOneID}/{teamTwoID}/", controllers.ComparePHLTeams).Methods("GET")
	apiRouter.HandleFunc("/ds/chl/assign/discord/{teamID}/{discordID}", controllers.AssignDiscordIDToCHLTeam).Methods("GET")
	apiRouter.HandleFunc("/ds/phl/assign/discord/{teamID}/{discordID}", controllers.AssignDiscordIDToPHLTeam).Methods("GET")
	apiRouter.HandleFunc("/ds/chl/croots/class/{teamID}/", controllers.GetCHLRecruitingClassByTeamID).Methods("GET")
	apiRouter.HandleFunc("/ds/chl/croot/{id}", controllers.GetCHLRecruitViaDiscord).Methods("GET")

	// Websocket
	myRouter.HandleFunc("/ws", ws.WebSocketHandler)

	// log.Fatal(http.ListenAndServe(":5001", handler))
	return handlers.CORS(originsOk, headersOk, methodsOk)(myRouter)
}

func loadEnvs() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("CANNOT LOAD ENV VARIABLES")
	}
}

func handleCron() *cron.Cron {
	c := cron.New()
	c.AddFunc("0 14 * * *", controllers.SyncFreeAgencyViaCron)
	c.AddFunc("0 12 * * 2,4,6,7", controllers.RunAIGameplanViaCron)

	c.Start()

	return c
}

func main() {
	loadEnvs()
	InitialMigration()
	fmt.Println("Setting up polling...")
	go monitorDBForUpdates()

	fmt.Println("Loading cron...")
	cronJobs := handleCron()
	fmt.Println("Loading Handler Requests.")
	fmt.Println("Hockey Server Initialized.")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handleRequests(),
	}

	go func() {
		fmt.Println("Server starting on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
		}
	}()

	// Create a channel to listen for system interrupts (Ctrl+C, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	fmt.Println("Shutting down server...")

	// Gracefully shutdown the server with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop cron jobs
	cronJobs.Stop()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
