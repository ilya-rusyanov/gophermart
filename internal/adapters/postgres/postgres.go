package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Logger interface {
	Infof(string, ...any)
}

func MustInit(ctx context.Context, logger Logger, dsn string, maxUsernameLen int) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	err = migrate(ctx, db, maxUsernameLen)
	if err != nil {
		panic(err)
	}
	return db
}

func migrate(ctx context.Context, db *sql.DB, maxUserNameLen int) error {
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users
(username varchar(`+strconv.Itoa(maxUserNameLen)+`) PRIMARY KEY,
password text NOT NULL,
balance numeric NOT NULL DEFAULT 0,
withdrawn numeric NOT NULL DEFAULT 0)`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.ExecContext(ctx,
		`DO $$ BEGIN
    CREATE TYPE order_state AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;`)
	if err != nil {
		return fmt.Errorf("failed to create order state enumeration: %w", err)
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS orders
(id text PRIMARY KEY, username text NOT NULL, upload_time timestamptz NOT NULL,
state order_state NOT NULL, value numeric)`)
	if err != nil {
		return fmt.Errorf("failed to create order table: %w", err)
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS withdrawals
(id text PRIMARY KEY, username text NOT NULL, upload_time timestamptz NOT NULL,
value numeric NOT NULL)`)
	if err != nil {
		return fmt.Errorf("failed to create withdrawals table: %w", err)
	}

	return nil
}
