package database

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type LegacyPremiumEntitlement struct {
	UserId    uint64    `json:"user_id"`
	TierId    int32     `json:"tier_id"`
	SkuLabel  string    `json:"sku_label"`
	ExpiresAt time.Time `json:"expires_at"`
}

type LegacyPremiumEntitlements struct {
	*pgxpool.Pool
}

func newLegacyPremiumEntitlement(db *pgxpool.Pool) *LegacyPremiumEntitlements {
	return &LegacyPremiumEntitlements{
		db,
	}
}

var (
	//go:embed sql/legacy_premium_entitlements/schema.sql
	legacyPremiumEntitlementsSchema string

	//go:embed sql/legacy_premium_entitlements/list_expired.sql
	legacyPremiumEntitlementsListAll string

	//go:embed sql/legacy_premium_entitlements/get_guild_tier.sql
	legacyPremiumEntitlementsGetGuildTier string

	//go:embed sql/legacy_premium_entitlements/get_user_entitlement.sql
	legacyPremiumEntitlementGetUserEntitlement string

	//go:embed sql/legacy_premium_entitlements/set_entitlement.sql
	legacyPremiumEntitlementsSet string

	//go:embed sql/legacy_premium_entitlements/delete_entitlement.sql
	legacyPremiumEntitlementsDelete string
)

func (e LegacyPremiumEntitlements) Schema() string {
	return legacyPremiumEntitlementsSchema
}

func (e *LegacyPremiumEntitlements) ListAll(ctx context.Context, tx pgx.Tx) ([]LegacyPremiumEntitlement, error) {
	rows, err := tx.Query(ctx, legacyPremiumEntitlementsListAll)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var entitlements []LegacyPremiumEntitlement
	for rows.Next() {
		var entitlement LegacyPremiumEntitlement
		if err := rows.Scan(
			&entitlement.UserId,
			&entitlement.TierId,
			&entitlement.SkuLabel,
			&entitlement.ExpiresAt,
		); err != nil {
			return nil, err
		}

		entitlements = append(entitlements, entitlement)
	}

	return entitlements, nil
}

func (e *LegacyPremiumEntitlements) GetGuildTier(ctx context.Context, guildId, ownerId uint64, gracePeriod time.Duration) (int32, bool, error) {
	var tier int32
	if err := e.QueryRow(ctx, legacyPremiumEntitlementsGetGuildTier, guildId, ownerId, gracePeriod).Scan(&tier); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, false, nil
		}

		return -1, false, err
	}

	return tier, true, nil
}

func (e *LegacyPremiumEntitlements) GetUserTier(ctx context.Context, userId uint64, gracePeriod time.Duration) (*LegacyPremiumEntitlement, error) {
	var entitlement LegacyPremiumEntitlement
	if err := e.QueryRow(ctx, legacyPremiumEntitlementGetUserEntitlement, userId, gracePeriod).Scan(
		&entitlement.UserId,
		&entitlement.TierId,
		&entitlement.SkuLabel,
		&entitlement.ExpiresAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &entitlement, nil
}

func (e *LegacyPremiumEntitlements) SetEntitlement(ctx context.Context, tx pgx.Tx, entitlement LegacyPremiumEntitlement) error {
	_, err := tx.Exec(ctx, legacyPremiumEntitlementsSet,
		entitlement.UserId,
		entitlement.TierId,
		entitlement.SkuLabel,
		entitlement.ExpiresAt,
	)

	return err
}

func (e *LegacyPremiumEntitlements) Delete(ctx context.Context, tx pgx.Tx, userId uint64, skuLabel string) error {
	_, err := tx.Exec(ctx, legacyPremiumEntitlementsDelete, userId, skuLabel)
	return err
}
