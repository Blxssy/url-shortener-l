package model

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	DB    *sql.DB
	Cache map[string]string
	UseDB bool
}

func NewStorage(db *sql.DB, useDB bool) *Storage {
	if useDB {
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		shortURL VARCHAR(10) PRIMARY KEY,
		longURL TEXT NOT NULL 
		)`)

		if err != nil {
			log.Fatal(err)
		}

		return &Storage{
			DB:    db,
			UseDB: true,
		}
	} else {
		storage := make(map[string]string)
		return &Storage{
			Cache: storage,
			UseDB: false,
		}
	}
}

func (s *Storage) GetLongURL(shortURL string) (string, error) {
	var longURL string
	if s.UseDB {
		err := s.DB.QueryRow("SELECT longURL FROM urls WHERE shortURL = $1", shortURL).Scan(&longURL)
		if err != nil {
			return "", err
		}
	} else {
		var found bool
		longURL, found = s.Cache[shortURL]
		if !found {
			return "", errors.New("not found")
		}
	}
	return longURL, nil
}

func (s *Storage) PutURL(shortURL string, longURL string) error {
	if s.UseDB {
		_, err := s.DB.Exec("INSERT INTO urls (shortURL, longURL) VALUES ($1, $2)", shortURL, longURL)
		if err != nil {
			return errors.New("internal server error")
		}
	} else {
		s.Cache[shortURL] = longURL
	}
	return nil
}
