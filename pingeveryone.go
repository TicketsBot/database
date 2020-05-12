package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PingEveryone struct {
	*pgxpool.Pool
}

func newPingEveryone(db *pgxpool.Pool) *PingEveryone {
	return &PingEveryone{
		db,
	}
}

func (p PingEveryone) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS ping_everyone(
	"guild_id" int8 NOT NULL UNIQUE,
	"ping_everyone" bool NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (p *PingEveryone) Get(guildId uint64) (pingEveryone bool, e error) {
	if err := p.QueryRow(context.Background(), `SELECT "ping_everyone" from ping_everyone WHERE "guild_id" = $1;`, guildId).Scan(&pingEveryone); err != nil {
		if err == pgx.ErrNoRows {
			pingEveryone = true
		} else {
			e = err
		}
	}

	return
}

func (p *PingEveryone) Set(guildId uint64, pingEveryone bool) (err error) {
	_, err = p.Exec(context.Background(), `INSERT INTO pingEveryone("guild_id", "ping_everyone") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "ping_everyone" = $2;`, guildId, pingEveryone)
	return
}
