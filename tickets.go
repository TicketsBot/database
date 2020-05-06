package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Ticket struct {
	Id               int
	GuildId          uint64
	ChannelId        uint64
	UserId           uint64
	Open             bool
	OpenTime         time.Time
	WelcomeMessageId uint64
}

type TicketTable struct {
	*pgxpool.Pool
}

func newTicketTable(db *pgxpool.Pool) *TicketTable {
	return &TicketTable{
		db,
	}
}

func (t TicketTable) Schema() string {
	return `CREATE TABLE IF NOT EXISTS tickets(
"id" int4 NOT NULL,
"guild_id" int8 NOT NULL,
"channel_id" int8 UNIQUE,
"user_id" int8 NOT NULL,
"open" bool NOT NULL,
"open_time" timestamp NOT NULL,
"welcome_message_id" int8,
PRIMARY KEY("id", "guild_id"));
CREATE INDEX CONCURRENTLY IF NOT EXISTS tickets_channel_id ON tickets("channel_id");
`
}

func (t *TicketTable) Create(guildId, userId uint64) (id int, err error) {
	query := `INSERT INTO tickets("id", "guild_id", "user_id", "open", "open_time") VALUES((SELECT COALESCE(MAX("id"), 0) + 1 FROM tickets WHERE "guild_id" = $1), $1, $2, true, NOW()) RETURNING "id";`
	err = t.QueryRow(context.Background(), query, guildId, userId).Scan(&id)
	return
}

func (t *TicketTable) SetTicketProperties(guildId uint64, ticketId int, channelId, welcomeMessageId uint64) (err error) {
	query := `UPDATE tickets SET "channel_id" = $1 AND "welcome_message_id" = $2 WHERE "guild_id" = $3 AND "id" = $4;`
	_, err = t.Exec(context.Background(), query, channelId, welcomeMessageId, guildId, ticketId)
	return
}

func (t *TicketTable) Get(ticketId int, guildId uint64) (ticket Ticket, e error) {
	query := `SELECT * FROM tickets WHERE "id" = $1 AND "guild_id" = $2;`
	if err := t.QueryRow(context.Background(), query, ticketId, guildId).Scan(&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId); err != nil && err != pgx.ErrNoRows {
		e = err
	}
	return
}

func (t *TicketTable) GetByChannel(channelId uint64) (ticket Ticket, e error) {
	query := `SELECT * FROM tickets WHERE "channel_id" = $1;`
	if err := t.QueryRow(context.Background(), query, channelId).Scan(&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId); err != nil && err != pgx.ErrNoRows {
		e = err
	}
	return
}

func (t *TicketTable) GetAllByUser(userId uint64) (tickets []Ticket, e error) {
	query := `SELECT * FROM tickets WHERE "user_id" = $1;`

	rows, err := t.Query(context.Background(), query, userId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetOpenByUser(userId uint64) (tickets []Ticket, e error) {
	query := `SELECT * FROM tickets WHERE "user_id" = $1 AND "open" = true;`

	rows, err := t.Query(context.Background(), query, userId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildOpenTickets(guildId uint64) (tickets []Ticket, e error) {
	query := `SELECT * FROM tickets WHERE "guild_id" = $1 AND "open" = true;`

	rows, err := t.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetTotalTicketCount(guildId uint64) (count int, e error) {
	query := `SELECT COUNT(*) FROM tickets WHERE "guild_id" = $1;`
	if err := t.QueryRow(context.Background(), query, guildId).Scan(&count); err != nil {
		e = err
	}
	return
}

func (t *TicketTable) Close(ticketId int, guildId uint64) (err error) {
	query := `UPDATE tickets SET "open"=false WHERE "id"=$1 AND "guild_id"=$2;`
	_, err = t.Exec(context.Background(), query, ticketId, guildId)
	return
}
