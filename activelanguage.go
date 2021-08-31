package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ActiveLanguage struct {
	*pgxpool.Pool
}

func newActiveLanguage(db *pgxpool.Pool) *ActiveLanguage {
	return &ActiveLanguage{
		db,
	}
}

func (l ActiveLanguage) Schema() string {
	return `CREATE TABLE IF NOT EXISTS active_language("guild_id" int8 NOT NULL UNIQUE, "language" varchar(8) NOT NULL, PRIMARY KEY("guild_id"));`
}

func (c *ActiveLanguage) Get(guildId uint64) (language string, e error) {
	if err := c.QueryRow(context.Background(), `SELECT "language" from active_language WHERE "guild_id" = $1`, guildId).Scan(&language); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (c *ActiveLanguage) Set(guildId uint64, language string) (err error) {
	_, err = c.Exec(context.Background(), `INSERT INTO active_language("guild_id", "language") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "language" = $2;`, guildId, language)
	return
}
