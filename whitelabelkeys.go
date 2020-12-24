package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelKeys struct {
	*pgxpool.Pool
}

func newWhitelabelKeys(db *pgxpool.Pool) *WhitelabelKeys {
	return &WhitelabelKeys{
		db,
	}
}

func (w WhitelabelKeys) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_keys(
	"bot_id" int8 UNIQUE NOT NULL,
	"key" varchar(64) NOT NULL,
	FOREIGN KEY("bot_id") REFERENCES whitelabel("bot_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("bot_id")
);
`
}

func (w *WhitelabelKeys) Get(botId uint64) (status string, e error) {
	query := `SELECT "key" FROM whitelabel_keys WHERE "bot_id" = $1;`
	if err := w.QueryRow(context.Background(), query, botId).Scan(&status); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (w *WhitelabelKeys) Set(botId uint64, key string) (err error) {
	query := `INSERT INTO whitelabel_keys("bot_id", "key") VALUES($1, $2) ON CONFLICT("bot_id") DO UPDATE SET "key" = $2;`
	_, err = w.Exec(context.Background(), query, botId, key)
	return
}

func (w *WhitelabelKeys) Delete(botId uint64) (err error) {
	query := `DELETE FROM whitelabel_keys WHERE "bot_id" = $1;`
	_, err = w.Exec(context.Background(), query, botId)
	return
}
