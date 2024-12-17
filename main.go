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
	apiRouter.HandleFunc("/admin/ai/generate/lineups/", controllers.RunAILineups).Methods("GET")
	apiRouter.HandleFunc("/admin/test/engine/", controllers.TestEngine).Methods("GET")
	apiRouter.HandleFunc("/admin/generate/test/rosters/", controllers.GenerateTestData).Methods("GET")

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
