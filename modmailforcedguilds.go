package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ModmailForcedGuilds struct {
	*pgxpool.Pool
}

func newModmailForcedGuilds(db *pgxpool.Pool) *ModmailForcedGuilds {
	return &ModmailForcedGuilds{
		db,
	}
}

func (m ModmailForcedGuilds) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS modmail_forced_guilds(
	"bot_id" int8 UNIQUE NOT NULL,
	"guild_id" int8 NOT NULL,
	FOREIGN KEY("bot_id") REFERENCES whitelabel("bot_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("bot_id")
);
`
}

func (m *ModmailForcedGuilds) Get(botId uint64) (guildId uint64, e error) {
	query := `SELECT "guild_id" FROM modmail_forced_guilds WHERE "bot_id" = $1;`
	if err := m.QueryRow(context.Background(), query, botId).Scan(&guildId); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (m *ModmailForcedGuilds) Set(botId, guildId uint64) (err error) {
	query := `INSERT INTO modmail_forced_guilds("bot_id", "guild_id") VALUES($1, $2) ON CONFLICT("bot_id") DO UPDATE SET "guild_id" = $2;`
	_, err = m.Exec(context.Background(), query, botId, guildId)
	return
}

func (m *ModmailForcedGuilds) Delete(botId uint64) (err error) {
	query := `DELETE FROM modmail_forced_guilds WHERE "bot_id"=$1;`
	_, err = m.Exec(context.Background(), query, botId)
	return
}
