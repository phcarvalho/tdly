package services

import "database/sql"

type Service struct {
	Board *BoardService
	Item  *ItemService
	DB    *sql.DB
}

func NewService(db *sql.DB) *Service {
	service := &Service{
		Board: &BoardService{
			DB: db,
		},
		Item: &ItemService{
			DB: db,
		},
	}

	return service
}
