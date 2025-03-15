package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

func New(addr string, maxOpenConns int, maxIdleConns int, maxIdleTimeout string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	maxIdleTimeoutParsed, err := time.ParseDuration(maxIdleTimeout)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(maxIdleTimeoutParsed)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
