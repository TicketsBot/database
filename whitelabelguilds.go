package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelGuilds struct {
	*pgxpool.Pool
}

func newWhitelabelGuilds(db *pgxpool.Pool) *WhitelabelGuilds {
	return &WhitelabelGuilds{
		db,
	}
}

func (w WhitelabelGuilds) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_guilds(
	"bot_id" int8 NOT NULL,
	"guild_id" int8 NOT NULL,
	FOREIGN KEY("bot_id") REFERENCES whitelabel("bot_id"),
	PRIMARY KEY("bot_id", "guild_id")
);`
}

func (w *WhitelabelGuilds) GetGuilds(botId uint64) (guilds []uint64, e error) {
	query := `SELECT "guild_id" from whitelabel_guilds WHERE "bot_id"=$1;`

	rows, err := w.Query(context.Background(), query, botId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			continue
		}

		guilds = append(guilds, id)
	}

	return
}

func (w *WhitelabelGuilds) GetBotByGuild(guildId uint64) (botId uint64, found bool, e error) {
	query := `SELECT "bot_id" from whitelabel_guilds WHERE "guild_id"=$1 LIMIT 1;`

	if err := w.QueryRow(context.Background(), query, guildId).Scan(&botId); err != nil {
		if err == pgx.ErrNoRows {
			found = false
		} else {
			e = err
		}
		return
	}

	found = true

	return
}

func (w *WhitelabelGuilds) Add(botId, guildId uint64) (err error) {
	query := `INSERT INTO whitelabel_guilds("bot_id", "guild_id") VALUES($1, $2) ON CONFLICT("bot_id", "guild_id") DO NOTHING;`
	_, err = w.Exec(context.Background(), query, botId, guildId)
	return
}

func (w *WhitelabelGuilds) Delete(botId, guildId uint64) (err error) {
	query := `DELETE FROM whitelabel_guilds WHERE "bot_id"=$1 AND "guild_id"=$2;`
	_, err = w.Exec(context.Background(), query, botId, guildId)
	return
}
