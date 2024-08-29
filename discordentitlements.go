package database

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DiscordEntitlements struct {
	*pgxpool.Pool
}

var (
	//go:embed sql/discord_entitlements/schema.sql
	discordEntitlementsSchema string

	//go:embed sql/discord_entitlements/create.sql
	discordEntitlementsCreate string

	//go:embed sql/discord_entitlements/get_entitlement_id.sql
	discordEntitlementsGetEntitlementId string

	//go:embed sql/discord_entitlements/list_all.sql
	discordEntitlementsListAll string
)

func newDiscordEntitlementsTable(db *pgxpool.Pool) *DiscordEntitlements {
	return &DiscordEntitlements{
		db,
	}
}

func (DiscordEntitlements) Schema() string {
	return discordEntitlementsSchema
}

func (e *DiscordEntitlements) Create(ctx context.Context, tx pgx.Tx, discordId uint64, entitlementId uuid.UUID) error {
	_, err := tx.Exec(ctx, discordEntitlementsCreate, discordId, entitlementId)
	return err
}

func (e *DiscordEntitlements) GetEntitlementId(ctx context.Context, tx pgx.Tx, discordId uint64) (*uuid.UUID, error) {
	var entitlementId uuid.UUID
	if err := tx.QueryRow(ctx, discordEntitlementsGetEntitlementId, discordId).Scan(&entitlementId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &entitlementId, nil
}

func (e *DiscordEntitlements) ListAll(ctx context.Context, tx pgx.Tx) (map[uint64]uuid.UUID, error) {
	res := make(map[uint64]uuid.UUID)

	rows, err := tx.Query(ctx, discordEntitlementsListAll)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var discordId uint64
		var entitlementId uuid.UUID
		if err := rows.Scan(&discordId, &entitlementId); err != nil {
			return nil, err
		}

		res[discordId] = entitlementId
	}

	return res, nil
}
