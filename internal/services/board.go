package services

import (
	"database/sql"
)

type Board struct {
	ID    int
	Title string
}

type BoardService struct {
	DB *sql.DB
}

func (m *BoardService) GetByID(id int) (*Board, error) {
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

func (m *BoardService) Insert(title string) (int, error) {
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
