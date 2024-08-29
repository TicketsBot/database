package database

import (
	"context"
	_ "embed"
	"errors"
	"github.com/TicketsBot/common/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DiscordStoreSkus struct {
	*pgxpool.Pool
}

var (
	//go:embed sql/discord_store_skus/schema.sql
	discordStoreSkusSchema string

	//go:embed sql/discord_store_skus/get_sku.sql
	discordStoreSkusGetSku string
)

func newDiscordStoreSkusTable(db *pgxpool.Pool) *DiscordStoreSkus {
	return &DiscordStoreSkus{
		db,
	}
}

func (DiscordStoreSkus) Schema() string {
	return discordStoreSkusSchema
}

func (e *DiscordStoreSkus) GetSku(ctx context.Context, discordId uint64) (*model.Sku, error) {
	var sku model.Sku
	if err := e.QueryRow(ctx, discordStoreSkusGetSku, discordId).Scan(&sku.Id, &sku.Label, &sku.SkuType); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &sku, nil
}
