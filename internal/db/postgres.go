package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"book-forum/internal/config"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *PostgresDB) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDB) DB() *sql.DB {
	return p.db
}
