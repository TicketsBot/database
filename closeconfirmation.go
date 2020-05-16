package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CloseConfirmation struct {
	*pgxpool.Pool
}

func newCloseConfirmation(db *pgxpool.Pool) *CloseConfirmation {
	return &CloseConfirmation{
		db,
	}
}

func (c CloseConfirmation) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS close_confirmation(
	"guild_id" int8 NOT NULL UNIQUE,
	"confirm" bool NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (c *CloseConfirmation) Get(guildId uint64) (confirm bool, e error) {
	if err := c.QueryRow(context.Background(), `SELECT "confirm" from close_confirmation WHERE "guild_id" = $1;`, guildId).Scan(&confirm); err != nil {
		if err == pgx.ErrNoRows {
			confirm = true
		} else {
			e = err
		}
	}

	return
}

func (c *CloseConfirmation) Set(guildId uint64, confirm bool) (err error) {
	_, err = c.Exec(context.Background(), `INSERT INTO close_confirmation("guild_id", "confirm") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "confirm" = $2;`, guildId, confirm)
	return
}
