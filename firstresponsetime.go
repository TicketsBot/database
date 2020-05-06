package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type FirstResponseTime struct {
	*pgxpool.Pool
}

func newFirstResponseTime(db *pgxpool.Pool) *FirstResponseTime {
	return &FirstResponseTime{
		db,
	}
}

func (f FirstResponseTime) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS first_response_time("guild_id" int8 NOT NULL, "ticket_id" int4 NOT NULL, "response_time" interval NOT NULL, PRIMARY KEY("guild_id", "ticket_id"));
CREATE INDEX CONCURRENTLY IF NOT EXISTS first_response_time_guild_id ON first_response_time("guild_id");
`
}

func (f *FirstResponseTime) HasResponse(guildId uint64, ticketId int) (hasResponse bool, e error) {
	query := `SELECT EXISTS(SELECT 1 FROM first_response_time WHERE "guild_id" = $1 AND "ticket_id" = $2);`
	if err := f.QueryRow(context.Background(), query, guildId, ticketId).Scan(&hasResponse); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (f *FirstResponseTime) GetAverage(guildId uint64, interval time.Duration) (responseTime time.Duration, e error) {
	query := `SELECT AVG(first_response_time.response_time) FROM first_response_time, tickets WHERE tickets.open_time > current_timestamp - $1 AND first_response_time.guild_id = $2;`
	if err := f.QueryRow(context.Background(), query, interval, guildId).Scan(&responseTime); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (f *FirstResponseTime) GetAverageAllTime(guildId uint64) (responseTime time.Duration, e error) {
	query := `SELECT AVG(first_response_time.response_time) FROM first_response_time, tickets WHERE first_response_time.guild_id = $1;`
	if err := f.QueryRow(context.Background(), query, guildId).Scan(&responseTime); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (f *FirstResponseTime) Set(guildId uint64, ticketId int, responseTime time.Duration) (err error) {
	query := `INSERT INTO first_response_time("guild_id", "ticket_id", "response_time") VALUES($1, $2, $3) ON CONFLICT("guild_id", "ticket_id") DO NOTHING;`
	_, err = f.Exec(context.Background(), query, guildId, ticketId, responseTime)
	return
}
