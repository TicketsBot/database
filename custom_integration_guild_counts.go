package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationGuildCountsView struct {
	*pgxpool.Pool
}

func newCustomIntegrationGuildCountsView(db *pgxpool.Pool) *CustomIntegrationGuildCountsView {
	return &CustomIntegrationGuildCountsView{
		db,
	}
}

func (v CustomIntegrationGuildCountsView) Schema() string {
	s := v.schema("custom_integration_guild_counts")
	for _, indexSchema := range v.indexes("custom_integration_guild_counts") {
		s += "\n"
		s += indexSchema
	}

	return s
}

func (v CustomIntegrationGuildCountsView) schema(tableName string) string {
	return fmt.Sprintf(`
CREATE MATERIALIZED VIEW IF NOT EXISTS %[1]s
AS
	SELECT integration_id, COUNT(*) AS COUNT
	FROM custom_integration_guilds
	GROUP BY integration_id
WITH DATA;
`, tableName)
}

func (v CustomIntegrationGuildCountsView) indexes(tableName string) []string {
	return []string{
		fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_integration_id_key ON %[1]s(integration_id);", tableName),
	}
}

func (v *CustomIntegrationGuildCountsView) Refresh(ctx context.Context) error {
	statements := slice(v.schema("custom_integration_guild_counts_new"))
	statements = append(statements, v.indexes("custom_integration_guild_counts_new")...)
	statements = append(statements,
		"DROP MATERIALIZED VIEW IF EXISTS custom_integration_guild_counts;",
		"ALTER MATERIALIZED VIEW custom_integration_guild_counts_new RENAME TO custom_integration_guild_counts;",
		"ALTER INDEX custom_integration_guild_counts_new_integration_id_key RENAME TO custom_integration_guild_counts_integration_id_key;",
	)

	tx, err := transact(ctx, v.Pool, statements...)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	return tx.Commit(ctx)
}
