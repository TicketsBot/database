package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Prefix struct {
	GuildId uint64
	Prefix  string
}

func (p Prefix) Schema() string {
	return `CREATE TABLE IF NOT EXISTS prefix("guild_id" int8 NOT NULL UNIQUE, "prefix" varchar(8) NOT NULL, PRIMARY KEY("guild_id"));`
}

func (p *Prefix) Get(db *Database) {
	if err := db.QueryRow(context.Background(), `SELECT "prefix" from prefix WHERE "guild_id" = $1`, p.GuildId).Scan(p.Prefix); err != nil && err != pgx.ErrNoRows {
		db.Logger.Error(err)
	}
}

func (p *Prefix) Set(db *Database) {
	if _, err := db.Exec(context.Background(), `INSERT INTO prefix("guild_id", "prefix") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "prefix" = $2;`, p.GuildId, p.Prefix); err != nil {
		db.Logger.Error(err)
	}
}
