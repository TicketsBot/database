package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TicketLimit struct {
	*pgxpool.Pool
}

func newTicketLimit(db *pgxpool.Pool) *TicketLimit {
	return &TicketLimit{
		db,
	}
}

func (t TicketLimit) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS ticket_limit(
	"guild_id" int8 NOT NULL UNIQUE,
	"limit" int2 NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (t *TicketLimit) Get(ctx context.Context, guildId uint64) (limit uint8, e error) {
	query := `SELECT "limit" from ticket_limit WHERE "guild_id" = $1;`
	if err := t.QueryRow(ctx, query, guildId).Scan(&limit); err != nil {
		if err == pgx.ErrNoRows {
			limit = 5
		} else {
			e = err
		}
	}

	return
}

func (t *TicketLimit) Set(ctx context.Context, guildId uint64, limit uint8) (err error) {
	query := `INSERT INTO ticket_limit("guild_id", "limit") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "limit" = $2;`
	_, err = t.Exec(ctx, query, guildId, limit)
	return
}
