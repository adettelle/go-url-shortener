package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConnector interface {
	Connect() (*sql.DB, error)
}

type DBConnection struct {
	DBParams string
}

func NewDBConnection(params string) *DBConnection {
	return &DBConnection{
		DBParams: params,
	}
}

func connect(dbParams string) (*sql.DB, error) {
	log.Println("Connecting to DB", dbParams)

	db, err := sql.Open("pgx", dbParams)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (dbCon *DBConnection) Connect() (*sql.DB, error) {
	return connect(dbCon.DBParams)
}
