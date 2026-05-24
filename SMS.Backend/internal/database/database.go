package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/lib/pq"
)

//go:embed migrations/001_init.sql
var migration001 string

//go:embed migrations/002_add_datasource_to_boxes.sql
var migration002 string

//go:embed migrations/003_boxes_content_to_description.sql
var migration003 string

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	for _, sql := range []string{migration001, migration002, migration003} {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
