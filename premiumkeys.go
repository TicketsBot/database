package database

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PremiumKeys struct {
	*pgxpool.Pool
}

func newPremiumKeys(db *pgxpool.Pool) *PremiumKeys {
	return &PremiumKeys{
		db,
	}
}

func (k PremiumKeys) Schema() string {
	return `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS premium_keys(
	"key" uuid NOT NULL UNIQUE,
	"length" interval NOT NULL,
	"premium_type" int NOT NULL,
	PRIMARY KEY("key")
);`
}

func (k *PremiumKeys) Create(ctx context.Context, key uuid.UUID, length time.Duration, premiumType int) (err error) {
	_, err = k.Exec(ctx, `INSERT INTO premium_keys("key", "length", "premium_type") VALUES($1, $2, $3);`, key, length, premiumType)
	return
}

func (k *PremiumKeys) Delete(ctx context.Context, key uuid.UUID) (length time.Duration, premiumType int, e error) {
	if err := k.QueryRow(ctx, `DELETE from premium_keys WHERE "key" = $1 RETURNING "length", "premium_type";`, key).Scan(&length, &premiumType); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}
