package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ArchiveChannel struct {
	*pgxpool.Pool
}

func newArchiveChannel(db *pgxpool.Pool) *ArchiveChannel {
	return &ArchiveChannel{
		db,
	}
}

func (c ArchiveChannel) Schema() string {
	return `CREATE TABLE IF NOT EXISTS archive_channel("guild_id" int8 NOT NULL UNIQUE, "channel_id" int8 NOT NULL UNIQUE, PRIMARY KEY("guild_id"));`
}

func (c *ArchiveChannel) Get(guildId uint64) (archiveChannel uint64, e error) {
	if err := c.QueryRow(context.Background(), `SELECT "channel_id" from archive_channel WHERE "guild_id" = $1`, guildId).Scan(&archiveChannel); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (c *ArchiveChannel) Set(guildId, archiveChannel uint64) (err error) {
	_, err = c.Exec(context.Background(), `INSERT INTO archive_channel("guild_id", "channel_id") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "channel_id" = $2;`, guildId, archiveChannel)
	return
}
