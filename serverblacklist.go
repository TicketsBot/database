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

func (b *ServerBlacklist) IsBlacklisted(ctx context.Context, guildId uint64) (blacklisted bool, err error) {
	err = b.QueryRow(ctx, `SELECT EXISTS (SELECT 1 from server_blacklist WHERE "guild_id" = $1);`, guildId).Scan(&blacklisted)
	return
}

func (b *ServerBlacklist) Add(ctx context.Context, guildId uint64) (err error) {
	_, err = b.Exec(ctx, `INSERT INTO server_blacklist("guild_id") VALUES($1) ON CONFLICT("guild_id") DO NOTHING;`, guildId)
	return
}

func (b *ServerBlacklist) Delete(ctx context.Context, guildId uint64) (err error) {
	_, err = b.Exec(ctx, `DELETE FROM server_blacklist WHERE "guild_id" = $1;`, guildId)
	return
}
