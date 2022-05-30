package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type FirstResponseTimeData struct {
	GuildId uint64         `json:"guild_id"`
	AllTime *time.Duration `json:"all_time"`
	Monthly *time.Duration `json:"monthly"`
	Weekly  *time.Duration `json:"weekly"`
}

type FirstResponseTimeGuildView struct {
	*pgxpool.Pool
}

func newFirstResponseTimeGuildView(db *pgxpool.Pool) *FirstResponseTimeGuildView {
	return &FirstResponseTimeGuildView{
		db,
	}
}

func (d FirstResponseTimeGuildView) Schema() string {
	s := d.schema("first_response_time_guild_view")
	for _, indexSchema := range d.indexes("first_response_time_guild_view") {
		s += "\n"
		s += indexSchema
	}

	return s
}

func (d FirstResponseTimeGuildView) schema(tableName string) string {
	return fmt.Sprintf(`
CREATE MATERIALIZED VIEW IF NOT EXISTS %[1]s
AS
	SELECT
		first_response_time.guild_id,
		AVG(first_response_time.response_time) AS "all_time",
		AVG(first_response_time.response_time) FILTER (WHERE tickets.open_time > NOW() - INTERVAL '30d') AS "monthly",
		AVG(first_response_time.response_time) FILTER (WHERE tickets.open_time > NOW() - INTERVAL '7d') AS "weekly"
	FROM first_response_time
	INNER JOIN tickets
	ON first_response_time.guild_id = tickets.guild_id AND first_response_time.ticket_id = tickets.id
	GROUP BY first_response_time.guild_id
WITH DATA;
`, tableName)
}

func (d FirstResponseTimeGuildView) indexes(tableName string) []string {
	return []string{
		fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_guild_id_key ON %[1]s(guild_id);", tableName),
	}
}

func (d *FirstResponseTimeGuildView) Refresh() error {
	statements := slice(d.schema("first_response_time_guild_view_new"))
	statements = append(statements, d.indexes("first_response_time_guild_view_new")...)
	statements = append(statements,
		"DROP MATERIALIZED VIEW IF EXISTS first_response_time_guild_view;",
		"ALTER MATERIALIZED VIEW first_response_time_guild_view_new RENAME TO first_response_time_guild_view;",
		"ALTER INDEX first_response_time_guild_view_new_guild_id_key RENAME TO first_response_time_guild_view_guild_id_key;",
	)

	tx, err := transact(d.Pool, statements...)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (d *FirstResponseTimeGuildView) Get(guildId uint64) (FirstResponseTimeData, error) {
	query := `
SELECT "guild_id", "all_time", "monthly", "weekly"
FROM first_response_time_guild_view
WHERE "guild_id" = $1;
`

	var data FirstResponseTimeData
	if err := d.QueryRow(context.Background(), query, guildId).Scan(&data); err != nil {
		if err != pgx.ErrNoRows {
			return FirstResponseTimeData{GuildId: guildId}, nil // Return durations of zero
		} else {
			return FirstResponseTimeData{}, err
		}
	}

	return data, nil
}
