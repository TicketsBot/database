package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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
	return `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS used_keys(
	"key" uuid NOT NULL UNIQUE,
	"guild_id" int8 NOT NULL,
	"activated_by" int8 NOT NULL,
	PRIMARY KEY("key")
);`
}

func (k *UsedKeys) Set(ctx context.Context, tx pgx.Tx, key uuid.UUID, guildId, userId uint64) (err error) {
	_, err = tx.Exec(ctx, `INSERT INTO used_keys("key", "guild_id", "activated_by") VALUES($1, $2, $3) ON CONFLICT("key") DO NOTHING;`, key, guildId, userId)
	return
}
