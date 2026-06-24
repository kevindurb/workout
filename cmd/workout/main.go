package main

import (
	"database/sql"
	"log"
	"net/http"

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

	log.Fatal(http.ListenAndServe(":1337", a.Routes()))
}
