package database

import (
	"context"
	"errors"
	"github.com/google/uuid"
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
	"sku_id" UUID NOT NULL,
	"generated_at" TIMESTAMPTZ NOT NULL,
	PRIMARY KEY("key"),
	FOREIGN KEY("sku_id") REFERENCES skus("id")
);`
}

func (k *PremiumKeys) Create(ctx context.Context, key uuid.UUID, length time.Duration, skuId uuid.UUID) (err error) {
	_, err = k.Exec(ctx, `INSERT INTO premium_keys("key", "length", "sku_id", "generated_at") VALUES($1, $2, $3, NOW());`, key, length, skuId)
	return
}

func (k *PremiumKeys) Delete(ctx context.Context, tx pgx.Tx, key uuid.UUID) (time.Duration, uuid.UUID, bool, error) {
	var length time.Duration
	var skuId uuid.UUID

	query := `DELETE from premium_keys WHERE "key" = $1 RETURNING "length", "sku_id";`
	if err := tx.QueryRow(ctx, query, key).Scan(&length, &skuId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, uuid.Nil, false, nil
		}

		return 0, uuid.Nil, false, err
	}

	return length, skuId, true, nil
}
