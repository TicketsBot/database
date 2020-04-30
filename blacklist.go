package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Blacklist struct {
	*pgxpool.Pool
}

func newBlacklist(db *pgxpool.Pool) *Blacklist {
	return &Blacklist{
		db,
	}
}

func (b Blacklist) Schema() string {
	return `CREATE TABLE IF NOT EXISTS blacklist("guild_id" int8 NOT NULL, "user_id" int8 NOT NULL, PRIMARY KEY("guild_id", "user_id");`
}

func (b *Blacklist) IsBlacklisted(guildId, userId uint64) (exists bool, e error) {
	query := `SELECT EXISTS(SELECT 1 FROM blacklist WHERE "guild_id"=$1 AND "user_id"=$2);`
	if err := b.QueryRow(context.Background(), query, guildId, userId).Scan(&exists); err != nil {
		e = err
	}

	return
}

func (b *Blacklist) Add(guildId, userId uint64) (err error) {
	// on conflict, user is already blacklisted
	query := `INSERT INTO blacklist("guild_id", "user_id") VALUES($1, $2) ON CONFLICT DO NOTHING;`
	_, err = b.Exec(context.Background(), query, guildId, userId)
	return
}

func (b *Blacklist) Remove(guildId, userId uint64) (err error) {
	query := `DELETE FROM blacklist WHERE "guild_id"=$1 AND "user_id"=$2;`
	_, err = b.Exec(context.Background(), query, guildId, userId)
	return
}
