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
	"bot_id" int8 UNIQUE NOT NULL,
	"error" varchar(255) NOT NULL,
	FOREIGN KEY("bot_id") REFERENCES whitelabel("bot_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("error_id")
);
`
}

func (w *WhitelabelErrors) GetRecent(botId uint64, limit int) (errors []string, e error) {
	query := `SELECT "error" FROM whitelabel_statuses WHERE "bot_id" = $1 ORDER BY "error_id" DESC LIMIT $2;`

	rows, err := w.Query(context.Background(), query, botId, limit)
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

func (w *WhitelabelErrors) Append(botId uint64, error string) (err error) {
	query := `INSERT INTO whitelabel_errors("bot_id", "error") VALUES($1, $2);`
	_, err = w.Exec(context.Background(), query, botId, error)
	return
}
