package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Logger interface {
	Infof(string, ...any)
}

func MustInit(ctx context.Context, logger Logger, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	err = migrate(ctx, db)
	if err != nil {
		panic(err)
	}
	return db
}

func migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXITS users
(username text PRIMARY KEY, password text NOT NULL)`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}
