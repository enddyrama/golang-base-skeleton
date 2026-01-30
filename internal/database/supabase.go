package database

import (
	"base-skeleton/config"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

func NewSupabase(cfg *config.Config) (*sql.DB, error) {
	if cfg.DBSupabase == "" {
		return nil, errors.New("DB_SUPABASE connection string is empty")
	}

	db, err := sql.Open("postgres", cfg.DBSupabase)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
