package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Prefix struct {
	*pgxpool.Pool
}

func newPrefix(db *pgxpool.Pool) *Prefix {
	return &Prefix{
		db,
	}
}

func (t Prefix) Schema() string {
	return `CREATE TABLE IF NOT EXISTS prefix("guild_id" int8 NOT NULL UNIQUE, "prefix" varchar(8) NOT NULL, PRIMARY KEY("guild_id"));`
}

func (t *Prefix) Get(guildId uint64) (prefix string, e error) {
	query := `SELECT "prefix" from prefix WHERE "guild_id" = $1;`
	if err := t.QueryRow(context.Background(), query, guildId).Scan(&prefix); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (t *Prefix) Set(guildId uint64, prefix string) (err error) {
	query := `INSERT INTO prefix("guild_id", "prefix") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "prefix" = $2;`
	_, err = t.Exec(context.Background(), query, guildId, prefix)
	return
}
