package models

import (
	"database/sql"
)

type Item struct {
	ID        int
	Text      string
	Completed bool
}

type ItemModel struct {
	DB *sql.DB
}

func (m *ItemModel) GetByID(id int) (*Item, error) {
	stmt := "SELECT id, text, completed FROM items" +
		"\nWHERE id = ?"

	row := m.DB.QueryRow(stmt, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	var item Item
	err = row.Scan(&item.ID, &item.Text, &item.Completed)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (m *ItemModel) GetByBoardID(boardID int) ([]*Item, error) {
	stmt := "SELECT id, text, completed FROM items" +
		"\nWHERE board_id = ?"

	items := []*Item{}
	rows, err := m.DB.Query(stmt, boardID)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.Text, &item.Completed)
		if err != nil {
			return items, err
		}

		items = append(items, &item)
	}

	return items, nil
}

func (m *ItemModel) Insert(boardID int, text string) (int, error) {
	stmt := "INSERT INTO items (board_id, text)" +
		"VALUES (?, ?)"

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

func (m *ItemModel) ToggleByID(id int) error {
	stmt := "UPDATE items" +
		"\nSET completed = NOT completed" +
		"\nWHERE id = ?"

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
