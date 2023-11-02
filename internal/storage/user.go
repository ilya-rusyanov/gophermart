package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

func (u *User) FindUser(ctx context.Context, login entities.Login) error {
	row := u.db.QueryRowContext(ctx,
		"SELECT username FROM users WHERE username = $1", login)
	var user string
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.ErrNotFound
		}

		return fmt.Errorf("SQL error: %w", err)
	}

	return nil
}

func (u *User) AddCredentials(ctx context.Context, creds entities.AuthCredentials) error {
	_, err := u.db.ExecContext(ctx,
		`INSERT INTO users (username, password)
VALUES ($1, $2)`, creds.Login, creds.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entities.ErrLoginConflict
		}

		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}

func (u *User) FindCredentials(ctx context.Context, creds entities.AuthCredentials) error {
	row := u.db.QueryRowContext(ctx,
		`SELECT username FROM users WHERE username = $1 AND password = $2`,
		creds.Login, creds.Password)
	var user string
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.ErrNotFound
		}

		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}
