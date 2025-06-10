package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "./migrate"

func runMigrations(dsn string) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to open DB: %w", err))
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Errorf("failed to set dialect: %w", err))
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		panic(fmt.Errorf("failed to apply migrations: %w", err))
	}
}
