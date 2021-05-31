package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type WhitelabelErrors struct {
	*pgxpool.Pool
}

func newWhitelabelErrors(db *pgxpool.Pool) *WhitelabelErrors {
	return &WhitelabelErrors{
		db,
	}
}

type WhitelabelError struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func (w WhitelabelErrors) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_errors(
	"error_id" serial,
	"user_id" int8 NOT NULL,
	"error" varchar(255) NOT NULL,
	"error_time" timestamptz NOT NULL,
	PRIMARY KEY("error_id")
);
`
}

func (w *WhitelabelErrors) GetRecent(userId uint64, limit int) (errors []WhitelabelError, e error) {
	query := `SELECT "error", "error_time" FROM whitelabel_errors WHERE "user_id" = $1 ORDER BY "error_id" DESC LIMIT $2;`

	rows, err := w.Query(context.Background(), query, userId, limit)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var error WhitelabelError
		if e = rows.Scan(&error.Message, &error.Time); e != nil {
			continue
		}

		errors = append(errors, error)
	}

	return
}

func (w *WhitelabelErrors) Append(userId uint64, error string) (err error) {
	query := `INSERT INTO whitelabel_errors("user_id", "error", "error_time") VALUES($1, $2, NOW());`
	_, err = w.Exec(context.Background(), query, userId, error)
	return
}
