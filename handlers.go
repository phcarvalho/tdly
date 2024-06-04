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
	mux.HandleFunc("GET /b/{id}", a.handleBoardPage)
	mux.HandleFunc("POST /boards", a.handleBoardCreate)
	mux.HandleFunc("POST /boards/{id}/items", a.handleItemCreate)
	mux.HandleFunc("POST /boards/{id}/items/{itemID}/toggle", a.handleItemToggle)

	return mux
}

func (a *application) handleHomePage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
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

	http.Redirect(w, r, fmt.Sprintf("/b/%d", id), http.StatusSeeOther)
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

	items, err := a.Service.Item.GetByBoardID(board.ID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := boardData{
		Board: board,
		Items: items,
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/board.html",
		"./ui/html/partials/item.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (a *application) handleItemCreate(w http.ResponseWriter, r *http.Request) {
	boardID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Board id '%s' is invalid", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Item text should not be blank", http.StatusBadRequest)
		return
	}

	itemID, err := a.Service.Item.Insert(boardID, text)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	item, err := a.Service.Item.GetByID(itemID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ts, err := template.ParseFiles("./ui/html/partials/item.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "item", item)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (a *application) handleItemToggle(w http.ResponseWriter, r *http.Request) {
	boardID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Board id '%s' is invalid", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Item id '%s' is invalid", r.PathValue("itemID")), http.StatusBadRequest)
		return
	}

	item, err := a.Service.Item.GetByID(itemID)
	if err != nil || item.BoardID != boardID {
		http.Error(w, fmt.Sprintf("Item %d not found", itemID), http.StatusNotFound)
		return
	}

	err = a.Service.Item.ToggleByID(itemID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	item, err = a.Service.Item.GetByID(itemID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ts, err := template.ParseFiles("./ui/html/partials/item.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "item", item)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
