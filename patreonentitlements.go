package database

import (
	"context"
	_ "embed"
	"github.com/TicketsBot/common/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PatreonEntitlements struct {
	*pgxpool.Pool
}

func newPatreonEntitlements(db *pgxpool.Pool) *PatreonEntitlements {
	return &PatreonEntitlements{
		db,
	}
}

var (
	//go:embed sql/patreon_entitlements/schema.sql
	patreonEntitlementsSchema string

	//go:embed sql/patreon_entitlements/insert.sql
	patreonEntitlementsInsert string

	//go:embed sql/patreon_entitlements/list_by_user.sql
	patreonEntitlementsListByUser string

	//go:embed sql/patreon_entitlements/delete.sql
	patreonEntitlementsDelete string

	//go:embed sql/patreon_entitlements/delete_by_user.sql
	patreonEntitlementsDeleteByUser string
)

func (e PatreonEntitlements) Schema() string {
	return patreonEntitlementsSchema
}

func (e *PatreonEntitlements) Insert(ctx context.Context, tx pgx.Tx, entitlementId uuid.UUID, userId uint64) error {
	_, err := tx.Exec(ctx, patreonEntitlementsInsert, entitlementId, userId)
	return err
}

func (e *PatreonEntitlements) ListByUser(ctx context.Context, tx pgx.Tx, userId uint64) ([]model.Entitlement, error) {
	rows, err := tx.Query(ctx, patreonEntitlementsListByUser, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var entitlements []model.Entitlement
	for rows.Next() {
		var entitlement model.Entitlement
		if err := rows.Scan(
			&entitlement.Id,
			&entitlement.GuildId,
			&entitlement.UserId,
			&entitlement.SkuId,
			&entitlement.Source,
			&entitlement.ExpiresAt,
		); err != nil {
			return nil, err
		}

		entitlements = append(entitlements, entitlement)
	}

	return entitlements, nil
}

func (e *PatreonEntitlements) Delete(ctx context.Context, tx pgx.Tx, entitlementId uuid.UUID) error {
	_, err := tx.Exec(ctx, patreonEntitlementsDelete, entitlementId)
	return err
}

func (e *PatreonEntitlements) DeleteByUser(ctx context.Context, tx pgx.Tx, userId uint64) error {
	_, err := tx.Exec(ctx, patreonEntitlementsDeleteByUser, userId)
	return err
}
