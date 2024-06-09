package services

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/sqids/sqids-go"
)

type ConfigService struct {
	db        *sql.DB
	mutex     *sync.Mutex
	encoder   *sqids.Sqids
	AppConfig *Config
}

type Config struct {
	Alphabet    string
	LastBoardID uint
}

func newConfigService(db *sql.DB) *ConfigService {
	config := &Config{
		LastBoardID: 0,
	}

	row := db.QueryRow("SELECT alphabet, last_board_id FROM configs")
	if row.Err() != nil {
		log.Fatal(row.Err().Error())
	}

	err := row.Scan(&config.Alphabet, &config.LastBoardID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err.Error())
		}

		letters := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
		rand.Shuffle(len(letters), func(i, j int) {
			letters[i], letters[j] = letters[j], letters[i]
		})
		config.Alphabet = strings.Join(letters, "")

		_, err := db.Exec("INSERT INTO configs(alphabet, last_board_id) VALUES(?, ?)", config.Alphabet, config.LastBoardID)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	encoder, err := sqids.New(sqids.Options{
		Alphabet:  config.Alphabet,
		MinLength: uint8(8),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return &ConfigService{
		db:        db,
		mutex:     &sync.Mutex{},
		encoder:   encoder,
		AppConfig: config,
	}
}

func (m *ConfigService) GetNextID() (string, error) {
	m.mutex.Lock()
	m.AppConfig.LastBoardID += 1

	id, err := m.encoder.Encode([]uint64{uint64(m.AppConfig.LastBoardID)})
	if err != nil {
		m.mutex.Unlock()
		return "", err
	}

	_, err = m.db.Exec("UPDATE configs SET last_board_id = ?", m.AppConfig.LastBoardID)
	if err != nil {
		m.mutex.Unlock()
		return "", err
	}
	m.mutex.Unlock()

	return id, err
}
