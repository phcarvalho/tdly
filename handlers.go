package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/phcarvalho/tdly/internal/services"
)

type boardPageData struct {
	Board *services.Board
	Items []*services.Item
}

func (a *application) getServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", a.handleHomePage)
	mux.HandleFunc("GET /b/{id}", a.handleBoardPage)
	mux.HandleFunc("POST /boards", a.handleBoardCreate)
	mux.HandleFunc("POST /items", a.handleItemCreate)
	mux.HandleFunc("POST /items/{id}/toggle", a.handleItemToggle)
	mux.HandleFunc("DELETE /items/{id}/", a.handleItemDelete)

	return mux
}

func getBoardID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("board")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
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
	id, err := a.Service.Config.GetNextID()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = a.Service.Board.Insert(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/b/%s", id), http.StatusSeeOther)
}

func (a *application) handleBoardPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	board, err := a.Service.Board.GetByID(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	cookie := http.Cookie{
		Name:     "board",
		Value:    id,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	items, err := a.Service.Item.GetByBoardID(board.ID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := boardPageData{
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
	boardID, err := getBoardID(r)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

func (a *application) handleItemDelete(w http.ResponseWriter, r *http.Request) {
	boardID, err := getBoardID(r)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	itemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Item id '%s' is invalid", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	item, err := a.Service.Item.GetByID(itemID)
	if err != nil || item.BoardID != boardID {
		http.Error(w, fmt.Sprintf("Item %d not found", itemID), http.StatusNotFound)
		return
	}

	err = a.Service.Item.DeleteByID(itemID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *application) handleItemToggle(w http.ResponseWriter, r *http.Request) {
	boardID, err := getBoardID(r)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	itemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Item id '%s' is invalid", r.PathValue("id")), http.StatusBadRequest)
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
