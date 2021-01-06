package database

import (
	"context"
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
