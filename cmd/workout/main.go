package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kevindurb/planner/internal/app"
	"github.com/kevindurb/planner/internal/database"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "database.db?"+
		"_pragma=journal_mode(WAL)"+
		"&_pragma=foreign_keys(ON)"+
		"&_pragma=busy_timeout(5000)"+
		"&_pragma=synchronous(NORMAL)",
	)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer db.Close()

	database.Up(db)

	a := app.New(db)

	srv := &http.Server{
		Addr:    "0.0.0.0:1337",
		Handler: a.Routes(),
	}

	go func() {
		log.Printf("Listening on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down…")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Forced shutdown: %s", err)
	}

	log.Println("Server stopped.")
}
