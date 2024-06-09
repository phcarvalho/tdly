package services

import (
	"database/sql"
)

type Service struct {
	Board  *BoardService
	Item   *ItemService
	Config *ConfigService
	DB     *sql.DB
}

func NewService(db *sql.DB) *Service {
	configService := newConfigService(db)

	return &Service{
		Config: configService,
		Board: &BoardService{
			DB: db,
		},
		Item: &ItemService{
			DB: db,
		},
	}
}
