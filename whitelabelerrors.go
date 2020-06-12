package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelErrors struct {
	*pgxpool.Pool
}

func newWhitelabelErrors(db *pgxpool.Pool) *WhitelabelErrors {
	return &WhitelabelErrors{
		db,
	}
}

func (w WhitelabelErrors) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_errors(
	"error_id" serial,
	"user_id" int8 NOT NULL,
	"error" varchar(255) NOT NULL,
	PRIMARY KEY("error_id")
);
`
}

func (w *WhitelabelErrors) GetRecent(userId uint64, limit int) (errors []string, e error) {
	query := `SELECT "error" FROM whitelabel_errors WHERE "user_id" = $1 ORDER BY "error_id" DESC LIMIT $2;`

	rows, err := w.Query(context.Background(), query, userId, limit)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var message string
		if e = rows.Scan(&message); e != nil {
			continue
		}

		errors = append(errors, message)
	}

	return
}

func (w *WhitelabelErrors) Append(userId uint64, error string) (err error) {
	query := `INSERT INTO whitelabel_errors("user_id", "error") VALUES($1, $2);`
	_, err = w.Exec(context.Background(), query, userId, error)
	return
}
