package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type TicketDurationData struct {
	GuildId uint64        `json:"guild_id"`
	AllTime time.Duration `json:"all_time"`
	Monthly time.Duration `json:"monthly"`
	Weekly  time.Duration `json:"weekly"`
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
	return d.schema("ticket_duration")
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
CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_guild_id_key ON ticket_duration(guild_id);
`, tableName)
}

func (d *TicketDurationView) Refresh() error {
	tx, err := transact(d.Pool,
		d.schema("ticket_duration_new"),
		"DROP MATERIALIZED VIEW ticket_duration;",
		"ALTER MATERIALIZED VIEW ticket_duration_new RENAME TO ticket_duration;",
		"ALTER INDEX ticket_duration_new_guild_id_key RENAME TO ticket_duration_guild_id_key;",
	)

	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (d *TicketDurationView) Get(guildId uint64) (TicketDurationData, error) {
	query := `SELECT "guild_id", "all_time", "monthly", "weekly" FROM ticket_duration;`

	var data TicketDurationData
	if err := d.QueryRow(context.Background(), query, guildId).Scan(&data); err != nil {
		if err != pgx.ErrNoRows {
			return TicketDurationData{GuildId: guildId}, nil // Return durations of zero
		} else {
			return TicketDurationData{}, err
		}
	}

	return data, nil
}
