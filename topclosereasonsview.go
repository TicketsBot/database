package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TopCloseReasonsView struct {
	*pgxpool.Pool
}

func newTopCloseReasonsView(db *pgxpool.Pool) *TopCloseReasonsView {
	return &TopCloseReasonsView{
		db,
	}
}

func (v TopCloseReasonsView) Schema() string {
	s := v.schema("top_close_reasons")
	for _, indexSchema := range v.indexes("top_close_reasons") {
		s += "\n"
		s += indexSchema
	}

	return s
}

func (v TopCloseReasonsView) schema(tableName string) string {
	return fmt.Sprintf(`
CREATE MATERIALIZED VIEW IF NOT EXISTS %[1]s
AS
	SELECT * FROM (
		SELECT
			tickets.guild_id,
			tickets.panel_id,
			close_reason.close_reason,
			ROW_NUMBER() OVER (PARTITION BY tickets.guild_id, tickets.panel_id ORDER BY COUNT(*) DESC) AS ranking
		FROM close_reason
		INNER JOIN tickets
		ON close_reason.guild_id = tickets.guild_id AND close_reason.ticket_id = tickets.id
		WHERE "close_reason" != 'Automatically closed due to inactivity'
		GROUP BY tickets.guild_id, tickets.panel_id, close_reason
	) AS top_reasons_inner
	WHERE ranking <= 10
WITH DATA;
`, tableName)
}

func (v TopCloseReasonsView) indexes(tableName string) []string {
	return []string{
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %[1]s_guild_id_panel_id_key ON %[1]s(guild_id, panel_id);", tableName),
	}
}

func (v *TopCloseReasonsView) Refresh() error {
	statements := slice(v.schema("top_close_reasons_new"))
	statements = append(statements, v.indexes("top_close_reasons_new")...)
	statements = append(statements,
		"DROP MATERIALIZED VIEW IF EXISTS top_close_reasons;",
		"ALTER MATERIALIZED VIEW top_close_reasons_new RENAME TO top_close_reasons;",
		"ALTER INDEX top_close_reasons_new_guild_id_panel_id_key RENAME TO top_close_reasons_guild_id_panel_id_key;",
	)

	tx, err := transact(v.Pool, statements...)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (v *TopCloseReasonsView) Get(guildId uint64, panelId *int) ([]string, error) {
	query := `
SELECT "close_reason"
FROM top_close_reasons
WHERE "guild_id" = $1 AND "panel_id" = $2
ORDER BY "ranking" ASC;
`

	rows, err := v.Query(context.Background(), query, guildId, panelId)
	if err != nil {
		return nil, err
	}

	data := make([]string, 10)
	var count int
	for rows.Next() {
		if err := rows.Scan(&data[count]); err != nil {
			return nil, err
		}

		count++
	}

	return data[:count], nil
}
