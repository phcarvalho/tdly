package main

import (
	"database/sql"
	"github.com/phcarvalho/tdly/internal/models"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase(url string) (*sql.DB, error) {
	if url == "" {
		url = "./data/data.db"
	}

	db, err := sql.Open("sqlite3", url+"?_journal=WAL&_timeout=3000&_fk=true")
	if err != nil {
		return nil, err
	}

	err = setupDatabase(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupDatabase(db *sql.DB) error {
	boardStmt := "CREATE TABLE IF NOT EXISTS boards (" +
		"\n\tid INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"\n\ttitle TEXT NOT NULL," +
		"\n\tcreated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" +
		"\n)"
	_, err := db.Exec(boardStmt)
	if err != nil {
		return err
	}

	itemStmt := "CREATE TABLE IF NOT EXISTS items (" +
		"\n\tid INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ," +
		"\n\tboard_id INTEGER NOT NULL," +
		"\n\ttext TEXT NOT NULL," +
		"\n\tcompleted BOOLEAN NOT NULL DEFAULT FALSE," +
		"\n\tcreated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"\n\n\tFOREIGN KEY(board_id) REFERENCES boards(id) ON DELETE CASCADE" +
		"\n)"
	_, err = db.Exec(itemStmt)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := initDatabase("./data/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	boardModel := &models.BoardModel{DB: db}
	itemModel := &models.ItemModel{DB: db}

	boardID, err := boardModel.Insert("My new board")
	if err != nil {
		log.Fatal(err)
	}

	board, err := boardModel.GetByID(boardID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Board %d: %s\n", board.ID, board.Title)
	log.Println("++++++++++++++++++++++")

	_, err = itemModel.Insert(boardID, "That's my first item")
	if err != nil {
		log.Fatal(err)
	}
	_, err = itemModel.Insert(boardID, "That's my second item")
	if err != nil {
		log.Fatal(err)
	}
	id, err := itemModel.Insert(boardID, "That's my third item")
	if err != nil {
		log.Fatal(err)
	}

	err = itemModel.ToggleByID(id)
	if err != nil {
		log.Fatal(err)
	}

	items, err := itemModel.GetByBoardID(boardID)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		mark := " "
		if item.Completed {
			mark = "x"
		}

		log.Printf("- [%s] %s (%d)\n", mark, item.Text, item.ID)
	}
}
