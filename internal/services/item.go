package services

import (
	"database/sql"
	"time"
)

type Item struct {
	ID          int
	BoardID     int
	Text        string
	CompletedAt time.Time
}

type ItemService struct {
	DB *sql.DB
}

func (m *ItemService) GetByID(id int) (*Item, error) {
	stmt := "SELECT id, board_id, text, completed_at FROM items" +
		"\nWHERE id = ?"

	row := m.DB.QueryRow(stmt, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	var item Item
	var completedAt sql.NullTime
	err = row.Scan(&item.ID, &item.BoardID, &item.Text, &completedAt)
	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		item.CompletedAt = completedAt.Time
	}

	return &item, nil
}

func (m *ItemService) GetByBoardID(boardID int) ([]*Item, error) {
	stmt := "SELECT id, board_id, text, completed_at FROM items" +
		"\nWHERE board_id = ?"

	items := []*Item{}
	rows, err := m.DB.Query(stmt, boardID)
	if err != nil {
		return items, err
	}

	var completedAt sql.NullTime
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.BoardID, &item.Text, &completedAt)
		if err != nil {
			return items, err
		}

		if completedAt.Valid {
			item.CompletedAt = completedAt.Time
		}

		items = append(items, &item)
	}

	return items, nil
}

func (m *ItemService) Insert(boardID int, text string) (int, error) {
	stmt := "INSERT INTO items (board_id, text) VALUES (?, ?)"

	res, err := m.DB.Exec(stmt, boardID, text)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ItemService) DeleteByID(id int) error {
	stmt := "DELETE FROM items" +
		"\nWHERE id = ?"

	_, err := m.DB.Exec(stmt, id)
	return err
}

func (m *ItemService) ToggleByID(id int) error {
	stmt := "UPDATE items" +
		"\nSET completed_at = (CASE " +
		"\n\tWHEN completed_at IS NULL THEN CURRENT_TIMESTAMP ELSE NULL" +
		"\nEND)" +
		"\nWHERE id = ?"

	_, err := m.DB.Exec(stmt, id)
	return err
}
