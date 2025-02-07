package dbstorage

import (
	"context"
	"database/sql"
	"log"

	"github.com/adettelle/go-url-shortener/internal/helpers"
	"github.com/adettelle/go-url-shortener/internal/storage"
)

// DBStorage - это имплементация (или реализация) интерфейса Storage
type DBStorage struct {
	Ctx context.Context
	DB  *sql.DB
}

func NewURLRepo(ctx context.Context, db *sql.DB) *DBStorage {
	return &DBStorage{Ctx: ctx, DB: db}
}

func (s *DBStorage) GetAddress(name string) (string, error) {
	sqlStatement := "SELECT original_url where short_url = $1"
	row := s.DB.QueryRowContext(s.Ctx, sqlStatement, name)

	// переменная для чтения результата
	var originalURL string

	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s *DBStorage) AddAddress(originalURL string) (string, error) {
	log.Println("Writing to DB")

	if originalURL == "" {
		return "", &storage.EmptyAddressError{}
	}

	// // TODO повтор
	// rangeStart := 2
	// rangeEnd := 10
	// offset := rangeEnd - rangeStart
	// randLength := storage.SeededRand.Intn(offset) + rangeStart

	randString, err := helpers.StringWithCharset()
	if err != nil {
		return "", err
	}

	sqlStatement := `insert into url_mapping (short_url, original_url) 
		values ($1, $2)` // TODO on conflict ?????

	_, err = s.DB.ExecContext(s.Ctx, sqlStatement, randString, originalURL)
	if err != nil {
		log.Println("error in adding url:", err)
		return "", err
	}

	log.Println("Saved")
	return randString, nil
}

func (s *DBStorage) Finalize() error {
	return nil
}
