package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CloseReasonTable struct {
	*pgxpool.Pool
}

func newCloseReasonTable(db *pgxpool.Pool) *CloseReasonTable {
	return &CloseReasonTable{
		db,
	}
}

func (c CloseReasonTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS close_reason(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"close_reason" TEXT,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (c *CloseReasonTable) Get(guildId uint64, ticketId int) (reason string, ok bool, e error) {
	query := `
SELECT "close_reason"
FROM close_reason
WHERE "guild_id" = $1 AND "ticket_id" = $2;
`
	if err := c.QueryRow(context.Background(), query, guildId, ticketId).Scan(&reason); err != nil {
		if err != pgx.ErrNoRows {
			e = err
		}
	} else {
		ok = true
	}

	return
}

func (c *CloseReasonTable) GetCommon(guildId uint64, prefix string, limit int) ([]string, error) {
	query := `
SELECT "close_reason"
FROM close_reason
WHERE "guild_id" = $1 AND "close_reason" like '$2' || '%'
GROUP BY "close_reason"
ORDER BY COUNT(*) DESC
LIMIT $3;
`

	rows, err := c.Query(context.Background(), query, guildId, prefix, limit)
	if err != nil {
		return nil, err
	}

	var reasons []string
	for rows.Next() {
		var reason string
		if err := rows.Scan(&reason); err != nil {
			return nil, err
		}

		reasons = append(reasons, reason)
	}

	return reasons, nil
}

func (c *CloseReasonTable) GetMulti(guildId uint64, ticketIds []int) (map[int]string, error) {
	query := `
SELECT "ticket_id", "close_reason"
FROM close_reason
WHERE "guild_id" = $1 AND "ticket_id" = ANY($2);
`

	array := &pgtype.Int4Array{}
	if err := array.Set(ticketIds); err != nil {
		return nil, err
	}

	rows, err := c.Query(context.Background(), query, guildId, ticketIds)
	if err != nil {
		return nil, err
	}

	reasons := make(map[int]string)
	for rows.Next() {
		var ticketId int
		var reason string
		if err := rows.Scan(&ticketId, &reason); err != nil {
			return nil, err
		}

		reasons[ticketId] = reason
	}

	return reasons, nil
}

func (c *CloseReasonTable) Set(guildId uint64, ticketId int, reason string) (err error) {
	query := `
INSERT INTO close_reason("guild_id", "ticket_id", "close_reason")
VALUES($1, $2, $3)
ON CONFLICT("guild_id", "ticket_id") DO UPDATE SET "close_reason" = $3;
`

	_, err = c.Exec(context.Background(), query, guildId, ticketId, reason)
	return
}

func (c *CloseReasonTable) Delete(guildId uint64, ticketId int) (err error) {
	query := `DELETE FROM close_reason WHERE "guild_id"=$1 AND "ticket_id"=$2;`
	_, err = c.Exec(context.Background(), query, guildId, ticketId)
	return
}
