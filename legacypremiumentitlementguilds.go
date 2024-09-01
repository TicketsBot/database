package database

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LegacyPremiumEntitlementGuildRecord struct {
	UserId        uint64    `json:"user_id"`
	GuildId       uint64    `json:"guild_id"`
	EntitlementId uuid.UUID `json:"entitlement_id"`
}

type LegacyPremiumEntitlementGuilds struct {
	*pgxpool.Pool
}

var (
	//go:embed sql/legacy_premium_entitlement_guilds/schema.sql
	legacyPremiumEntitlementGuildsSchema string

	//go:embed sql/legacy_premium_entitlement_guilds/list_for_user.sql
	legacyPremiumEntitlementGuildsListForUser string

	//go:embed sql/legacy_premium_entitlement_guilds/insert.sql
	legacyPremiumEntitlementGuildsInsert string

	//go:embed sql/legacy_premium_entitlement_guilds/delete.sql
	legacyPremiumEntitlementGuildsDelete string

	//go:embed sql/legacy_premium_entitlement_guilds/delete_by_entitlement.sql
	legacyPremiumEntitlementGuildsDeleteByEntitlement string
)

func newLegacyPremiumEntitlementGuildsTable(db *pgxpool.Pool) *LegacyPremiumEntitlementGuilds {
	return &LegacyPremiumEntitlementGuilds{
		db,
	}
}

func (LegacyPremiumEntitlementGuilds) Schema() string {
	return legacyPremiumEntitlementGuildsSchema
}

func (g *LegacyPremiumEntitlementGuilds) ListForUser(ctx context.Context, tx pgx.Tx, userId uint64) ([]LegacyPremiumEntitlementGuildRecord, error) {
	rows, err := tx.Query(ctx, legacyPremiumEntitlementGuildsListForUser, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	records := make([]LegacyPremiumEntitlementGuildRecord, 0)
	for rows.Next() {
		var record LegacyPremiumEntitlementGuildRecord
		if err := rows.Scan(&record.UserId, &record.GuildId, &record.EntitlementId); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func (g *LegacyPremiumEntitlementGuilds) Insert(ctx context.Context, tx pgx.Tx, userId, guildId uint64, entitlementId uuid.UUID) error {
	_, err := tx.Exec(ctx, legacyPremiumEntitlementGuildsInsert, userId, guildId, entitlementId)
	return err
}

func (g *LegacyPremiumEntitlementGuilds) Delete(ctx context.Context, tx pgx.Tx, userId, guildId uint64) error {
	_, err := tx.Exec(ctx, legacyPremiumEntitlementGuildsDelete, userId, guildId)
	return err
}

func (g *LegacyPremiumEntitlementGuilds) DeleteByEntitlement(ctx context.Context, tx pgx.Tx, entitlementId uuid.UUID) error {
	_, err := tx.Exec(ctx, legacyPremiumEntitlementGuildsDeleteByEntitlement, entitlementId)
	return err
}
