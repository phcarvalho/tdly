package main

import (
	"database/sql"
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
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS boards (" +
		"\n\tid INTEGER PRIMARY KEY AUTOINCREMENT," +
		"\n\ttitle TEXT," +
		"\n\tcreated_at DATETIME DEFAULT CURRENT_TIMESTAMP" +
		"\n)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS items (" +
		"\n\tid INTEGER PRIMARY KEY AUTOINCREMENT," +
		"\n\tboard_id INT NOT NULL," +
		"\n\ttext TEXT NOT NULL," +
		"\n\tcompleted_at DATE," +
		"\n\tcreated_at DATE DEFAULT CURRENT_TIMESTAMP," +
		"\n\n\tFOREIGN KEY(board_id) REFERENCES boards(id)" +
		"\n)")
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
	log.Print("DB Created")
}
