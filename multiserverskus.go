package database

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MultiServerSkus struct {
	*pgxpool.Pool
}

var (
	//go:embed sql/multi_server_skus/schema.sql
	multiServerSkusSchema string

	//go:embed sql/multi_server_skus/get_permitted_server_count.sql
	multiServerSkusGetPermittedServerCount string
)

func newMultiServerSkusTable(db *pgxpool.Pool) *MultiServerSkus {
	return &MultiServerSkus{
		db,
	}
}

func (MultiServerSkus) Schema() string {
	return multiServerSkusSchema
}

func (m *MultiServerSkus) GetPermittedServerCount(ctx context.Context, tx pgx.Tx, skuId uuid.UUID) (int, bool, error) {
	var count int
	if err := tx.QueryRow(ctx, multiServerSkusGetPermittedServerCount, skuId).Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, nil
		}

		return 0, false, err
	}

	return count, true, nil
}
