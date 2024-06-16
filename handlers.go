package main

import (
	"fmt"
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

func (a *application) handleHomePage(w http.ResponseWriter, r *http.Request) {
	ts := a.templateCache["home.html"]

	err := ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (a *application) handleBoardCreate(w http.ResponseWriter, r *http.Request) {
	id, err := a.service.Config.GetNextID()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = a.service.Board.Insert(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/b/%s", id), http.StatusSeeOther)
}

func (a *application) handleBoardPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	board, err := a.service.Board.GetByID(id)
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

	items, err := a.service.Item.GetByBoardID(board.ID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := boardPageData{
		Board: board,
		Items: items,
	}

	ts := a.templateCache["board.html"]

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (a *application) handleItemCreate(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("board")
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	boardID := cookie.Value

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Item text should not be blank", http.StatusBadRequest)
		return
	}

	itemID, err := a.service.Item.Insert(boardID, text)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	item, err := a.service.Item.GetByID(itemID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ts := a.templateCache["board.html"]

	err = ts.ExecuteTemplate(w, "item", item)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (a *application) handleItemDelete(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("board")
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	boardID := cookie.Value

	itemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Item id '%s' is invalid", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	item, err := a.service.Item.GetByID(itemID)
	if err != nil || item.BoardID != boardID {
		http.Error(w, fmt.Sprintf("Item %d not found", itemID), http.StatusNotFound)
		return
	}

	err = a.service.Item.DeleteByID(itemID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *application) handleItemToggle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("board")
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	boardID := cookie.Value

	itemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Item id '%s' is invalid", r.PathValue("id")), http.StatusBadRequest)
		return
	}

	item, err := a.service.Item.GetByID(itemID)
	if err != nil || item.BoardID != boardID {
		http.Error(w, fmt.Sprintf("Item %d not found", itemID), http.StatusNotFound)
		return
	}

	err = a.service.Item.ToggleByID(itemID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	item, err = a.service.Item.GetByID(itemID)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ts := a.templateCache["board.html"]

	err = ts.ExecuteTemplate(w, "item", item)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
