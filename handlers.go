package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/phcarvalho/tdly/internal/services"
)

type boardData struct {
	Board *services.Board
	Items []*services.Item
}

func (a *application) getServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", a.handleHomePage)
	mux.HandleFunc("GET /boards/{id}", a.handleBoardPage)
	mux.HandleFunc("POST /boards", a.handleBoardCreate)

	return mux
}

func (a *application) handleHomePage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (a *application) handleBoardCreate(w http.ResponseWriter, r *http.Request) {
	id, err := a.Service.Board.Insert("My new board")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/boards/%d", id), http.StatusSeeOther)
}

func (a *application) handleBoardPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idNumber, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Id should be a number", http.StatusBadRequest)
		return
	}

	board, err := a.Service.Board.GetByID(idNumber)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/board.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := boardData{
		Board: board,
		Items: []*services.Item{},
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
