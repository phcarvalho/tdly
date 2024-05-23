package models

import (
	"database/sql"
)

type Board struct {
	ID    int
	Title string
}

type BoardModel struct {
	DB *sql.DB
}

func (m *BoardModel) GetByID(id int) (*Board, error) {
	var board Board
	stmt := "SELECT id, title FROM boards" +
		"\nWHERE id = ?"

	row := m.DB.QueryRow(stmt, id)

	err := row.Err()
	if err != nil {
		return nil, err
	}

	err = row.Scan(&board.ID, &board.Title)
	if err != nil {
		return nil, err
	}

	return &board, nil
}

func (m *BoardModel) Insert(title string) (int, error) {
	stmt := "INSERT INTO boards(title) VALUES(?)"

	res, err := m.DB.Exec(stmt, title)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
