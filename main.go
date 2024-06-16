package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/phcarvalho/tdly/internal/services"
)

type application struct {
	db            *sql.DB
	service       *services.Service
	router        *http.ServeMux
	templateCache map[string]*template.Template
}

func (a *application) init() {
	db, err := sql.Open("sqlite3", "./data.db?_journal=WAL&_timeout=3000&_fk=true")
	if err != nil {
		log.Fatal(err)
	}

	a.db = db
	a.service = services.NewService(db)
	a.router = a.getServeMux()
	a.templateCache, err = newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting application")
}

func (a *application) close() {
	log.Println("exiting application")
	a.db.Close()
}

func main() {
	app := &application{}
	app.init()

	// TODO: handle graceful shutdowns and interrupts
	defer app.close()

	if err := http.ListenAndServe(":4000", app.router); err != nil {
		log.Fatal(err)
	}
}
