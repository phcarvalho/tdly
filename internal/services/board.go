package services

import (
	"database/sql"
)

type Board struct {
	ID    string
	Title string
}

type BoardService struct {
	DB *sql.DB
}

func (m *BoardService) GetByID(id string) (*Board, error) {
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

func (m *BoardService) Insert(id string) error {
	stmt := "INSERT INTO boards(id) VALUES (?)"

	_, err := m.DB.Exec(stmt, id)

	return err
}

func (m *BoardService) Edit(token string, title string) error {
	stmt := "UPDATE boards SET title = ? WHERE id = ?"

	_, err := m.DB.Exec(stmt, title, token)

	return err
}
