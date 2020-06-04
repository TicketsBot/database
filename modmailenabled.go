package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ModmailEnabled struct {
	*pgxpool.Pool
}

func newModmailEnabled(db *pgxpool.Pool) *ModmailEnabled {
	return &ModmailEnabled{
		db,
	}
}

func (m ModmailEnabled) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS modmail_enabled(
	"guild_id" int8 NOT NULL UNIQUE,
	"enabled" bool NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (m *ModmailEnabled) Get(guildId uint64) (enabled bool, e error) {
	if err := m.QueryRow(context.Background(), `SELECT "enabled" from modmail_enabled WHERE "guild_id" = $1;`, guildId).Scan(&enabled); err != nil {
		if err == pgx.ErrNoRows {
			enabled = true
		} else {
			e = err
		}
	}

	return
}

func (m *ModmailEnabled) Set(guildId uint64, enabled bool) (err error) {
	_, err = m.Exec(context.Background(), `INSERT INTO modmail_enabled("guild_id", "enabled") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "enabled" = $2;`, guildId, enabled)
	return
}
