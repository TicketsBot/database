package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelStatuses struct {
	*pgxpool.Pool
}

func newWhitelabelStatuses(db *pgxpool.Pool) *WhitelabelStatuses {
	return &WhitelabelStatuses{
		db,
	}
}

func (w WhitelabelStatuses) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_statuses(
	"bot_id" int8 UNIQUE NOT NULL,
	"status" varchar(255) NOT NULL,
	FOREIGN KEY("bot_id") REFERENCES whitelabel("bot_id") ON DELETE CASCADE,
	PRIMARY KEY("user_id")
);
CREATE INDEX IF NOT EXISTS whitelabel_bot_id ON whitelabel("bot_id");
`
}

func (w *WhitelabelStatuses) Get(botId uint64) (status string, e error) {
	query := `SELECT "status" FROM whitelabel_statuses WHERE "bot_id" = $1;`
	if err := w.QueryRow(context.Background(), query, botId).Scan(&status); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (w *WhitelabelStatuses) Set(botId uint64, status string) (err error) {
	query := `INSERT INTO whitelabel_statuses("bot_id", "status") VALUES($1, $2) ON CONFLICT("bot_id") DO UPDATE SET "status" = $2;`
	_, err = w.Exec(context.Background(), query, botId, status)
	return
}

func (w *WhitelabelStatuses) Delete(botId uint64) (err error) {
	query := `DELETE FROM whitelabel_statuses WHERE "bot_id"=$1;`
	_, err = w.Exec(context.Background(), query, botId)
	return
}
