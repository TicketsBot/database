package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type TicketClaims struct {
	*pgxpool.Pool
}

func newTicketClaims(db *pgxpool.Pool) *TicketClaims {
	return &TicketClaims{
		db,
	}
}

func (c TicketClaims) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS ticket_claims(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"user_id" int8 NOT NULL,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (c *TicketClaims) Get(guildId uint64, ticketId int) (userId uint64, e error) {
	query := `SELECT "user_id" FROM ticket_claims WHERE "guild_id" = $1 AND "ticket_id" = $2;`
	if err := c.QueryRow(context.Background(), query, guildId, ticketId).Scan(&userId); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (c *TicketClaims) Set(guildId uint64, ticketId int, userId uint64) (err error) {
	query := `INSERT INTO ticket_claims("guild_id", "ticket_id", "user_id") VALUES($1, $2, $3) ON CONFLICT("guild_id", "ticket_id") DO UPDATE SET "user_id" = $3;`
	_, err = c.Exec(context.Background(), query, guildId, ticketId, userId)
	return
}

func (c *TicketClaims) Delete(guildId uint64, ticketId int) (err error) {
	query := `DELETE FROM ticket_claims WHERE "guild_id"=$1 AND "ticket_id"=$2;`
	_, err = c.Exec(context.Background(), query, guildId, ticketId)
	return
}

// stats
func (c *TicketClaims) GetClaimedSinceCount(guildId, userId uint64, interval time.Duration) (count int, e error) {
	query := `
SELECT COUNT(*)
FROM ticket_claims
INNER JOIN tickets
ON ticket_claims.guild_id = tickets.guild_id AND ticket_claims.ticket_id = tickets.id
WHERE ticket_claims.guild_id = $1 AND ticket_claims.user_id = $2 AND tickets.open_time > NOW() - $3;`

	if err := c.QueryRow(context.Background(), query, guildId, userId, interval).Scan(&count); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}
