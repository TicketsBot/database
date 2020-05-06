package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UsedKeys struct {
	*pgxpool.Pool
}

func newUsedKeys(db *pgxpool.Pool) *UsedKeys {
	return &UsedKeys{
		db,
	}
}

func (k UsedKeys) Schema() string {
	return `CREATE TABLE IF NOT EXISTS used_keys("key" varchar(36) NOT NULL UNIQUE, "guild_id" int8 NOT NULL, "activated_by" int8 NOT NULL, PRIMARY KEY("key"));`
}

func (k *UsedKeys) Set(key string, guildId, userId uint64) (err error) {
	_, err = k.Exec(context.Background(), `INSERT INTO used_keys("key", "guild_id", "activated_by") VALUES($1, $2, $3) ON CONFLICT("key") DO NOTHING;`, key, guildId, userId)
	return
}
