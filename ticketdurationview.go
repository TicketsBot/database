package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type TicketDurationData struct {
	GuildId uint64         `json:"guild_id"`
	AllTime *time.Duration `json:"all_time"`
	Monthly *time.Duration `json:"monthly"`
	Weekly  *time.Duration `json:"weekly"`
}

type TicketDurationView struct {
	*pgxpool.Pool
}

func newTicketDurationView(db *pgxpool.Pool) *TicketDurationView {
	return &TicketDurationView{
		db,
	}
}

func (d TicketDurationView) Schema() string {
	s := d.schema("ticket_duration")
	for _, indexSchema := range d.indexes("ticket_duration") {
		s += "\n"
		s += indexSchema
	}

	return s
}

func (d TicketDurationView) schema(tableName string) string {
	return fmt.Sprintf(`
CREATE MATERIALIZED VIEW IF NOT EXISTS %[1]s
AS
	SELECT
		guild_id,
		AVG(close_time-open_time) AS "all_time",
		AVG(close_time-open_time) FILTER (where "close_time" > NOW() - interval '30d') AS "monthly",
		AVG(close_time-open_time) FILTER (where "close_time" > NOW() - interval '7d') AS "weekly"
	FROM tickets
	WHERE close_time IS NOT NULL
	GROUP BY guild_id
WITH DATA;
`, tableName)
}

func (d TicketDurationView) indexes(tableName string) []string {
	return []string{
		fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_guild_id_key ON %[1]s(guild_id);", tableName),
	}
}

func (d *TicketDurationView) Refresh() error {
	statements := slice(d.schema("ticket_duration_new"))
	statements = append(statements, d.indexes("ticket_duration_new")...)
	statements = append(statements,
		"DROP MATERIALIZED VIEW IF EXISTS ticket_duration;",
		"ALTER MATERIALIZED VIEW ticket_duration_new RENAME TO ticket_duration;",
		"ALTER INDEX ticket_duration_new_guild_id_key RENAME TO ticket_duration_guild_id_key;",
	)

	tx, err := transact(d.Pool, statements...)
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	return tx.Commit(context.Background())
}

func (d *TicketDurationView) Get(guildId uint64) (TicketDurationData, error) {
	query := `
SELECT "guild_id", "all_time", "monthly", "weekly"
FROM ticket_duration
WHERE "guild_id" = $1;
`

	var data TicketDurationData
	if err := d.QueryRow(context.Background(), query, guildId).Scan(&data.GuildId, &data.AllTime, &data.Monthly, &data.Weekly); err != nil {
		if err == pgx.ErrNoRows {
			return TicketDurationData{GuildId: guildId}, nil // Return durations of zero
		} else {
			return TicketDurationData{}, err
		}
	}

	return data, nil
}
