package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServerBlacklist struct {
	*pgxpool.Pool
}

func newServerBlacklist(db *pgxpool.Pool) *ServerBlacklist {
	return &ServerBlacklist{
		db,
	}
}

func (b ServerBlacklist) Schema() string {
	return `CREATE TABLE IF NOT EXISTS server_blacklist("guild_id" int8 NOT NULL UNIQUE, PRIMARY KEY("guild_id"));`
}

func (b *ServerBlacklist) IsBlacklisted(guildId uint64) (blacklisted bool, e error) {
	var count int

	if err := b.QueryRow(context.Background(), `SELECT COUNT(*) from server_blacklist WHERE "guild_id" = $1;`, guildId).Scan(&count); err != nil {
		e = err
	}

	return count > 0, e
}

func (b *ServerBlacklist) Add(guildId uint64) (err error) {
	_, err = b.Exec(context.Background(), `INSERT INTO server_blacklist("guild_id") VALUES($1) ON CONFLICT("guild_id") DO NOTHING;`, guildId)
	return
}

func (b *ServerBlacklist) Delete(guildId uint64) (err error) {
	_, err = b.Exec(context.Background(), `DELETE FROM server_blacklist WHERE "guild_id" = $1;`, guildId)
	return
}
