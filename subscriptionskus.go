package database

import (
	"context"
	_ "embed"
	"errors"
	"github.com/TicketsBot/common/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SubscriptionSkus struct {
	*pgxpool.Pool
}

var (
	//go:embed sql/subscription_skus/schema.sql
	subscriptionSkusSchema string

	//go:embed sql/subscription_skus/get.sql
	subscriptionSkusGet string

	//go:embed sql/subscription_skus/search.sql
	subscriptionSkusSearch string
)

func newSubscriptionSkusTable(db *pgxpool.Pool) *SubscriptionSkus {
	return &SubscriptionSkus{
		db,
	}
}

func (SubscriptionSkus) Schema() string {
	return subscriptionSkusSchema
}

func (e *SubscriptionSkus) GetSku(ctx context.Context, tx pgx.Tx, skuId uuid.UUID) (*model.SubscriptionSku, error) {
	var sku model.SubscriptionSku
	if err := tx.QueryRow(ctx, subscriptionSkusGet, skuId).Scan(
		&sku.Id,
		&sku.Label,
		&sku.SkuType,
		&sku.Tier,
		&sku.Priority,
		&sku.IsGlobal,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &sku, nil
}

func (e *SubscriptionSkus) Search(ctx context.Context, label string, limit int) ([]model.SubscriptionSku, error) {
	rows, err := e.Query(ctx, subscriptionSkusSearch, label, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var skus []model.SubscriptionSku
	for rows.Next() {
		var sku model.SubscriptionSku
		if err := rows.Scan(&sku.Id, &sku.Label, &sku.SkuType, &sku.Tier, &sku.Priority, &sku.IsGlobal); err != nil {
			return nil, err
		}

		skus = append(skus, sku)
	}

	return skus, nil
}
