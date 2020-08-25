package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DmOnOpen struct {
	*pgxpool.Pool
}

func newDmOnOpen(db *pgxpool.Pool) *DmOnOpen {
	return &DmOnOpen{
		db,
	}
}

func (d DmOnOpen) Schema() string {
	return `CREATE TABLE IF NOT EXISTS dm_on_open("guild_id" int8 NOT NULL UNIQUE, "dm_on_open" bool NOT NULL, PRIMARY KEY("guild_id"));`
}

func (d *DmOnOpen) Get(guildId uint64) (dmOnOpen bool, e error) {
	if err := d.QueryRow(context.Background(), `SELECT "dm_on_open" from dm_on_open WHERE "guild_id" = $1;`, guildId).Scan(&dmOnOpen); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (d *DmOnOpen) Set(guildId uint64, dmOnOpen bool) (err error) {
	_, err = d.Exec(context.Background(), `INSERT INTO dm_on_open("guild_id", "dm_on_open") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "dm_on_open" = $2;`, guildId, dmOnOpen)
	return
}
