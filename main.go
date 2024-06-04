package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/phcarvalho/tdly/internal/services"
)

type application struct {
	DB      *sql.DB
	Service *services.Service
	Router  *http.ServeMux
}

func (a *application) init() {
	db, err := sql.Open("sqlite3", "./data.db?_journal=WAL&_timeout=3000&_fk=true")
	if err != nil {
		log.Fatal(err)
	}

	a.DB = db
	a.Service = services.NewService(db)
	a.Router = a.getServeMux()

	log.Println("starting application")
}

func (a *application) close() {
	log.Println("exiting application")
	a.DB.Close()
}

func main() {
	app := &application{}
	app.init()

	// TODO: handle graceful shutdowns and interrupts
	defer app.close()

	if err := http.ListenAndServe(":4000", app.Router); err != nil {
		log.Fatal(err)
	}
}
