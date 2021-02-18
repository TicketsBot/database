package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// ALTER TABLE tickets ADD COLUMN panel_id int8 DEFAULT NULL;
// ALTER TABLE tickets ADD CONSTRAINT fk_panel_id FOREIGN KEY(panel_id) REFERENCES panels("message_id") ON DELETE SET NULL ON UPDATE CASCADE;
type Ticket struct {
	Id               int
	GuildId          uint64
	ChannelId        *uint64
	UserId           uint64
	Open             bool
	OpenTime         time.Time
	WelcomeMessageId *uint64
	PanelId          *uint64
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
	return `
CREATE TABLE IF NOT EXISTS tickets(
	"id" int4 NOT NULL,
	"guild_id" int8 NOT NULL,
	"channel_id" int8 UNIQUE,
	"user_id" int8 NOT NULL,
	"open" bool NOT NULL,
	"open_time" timestamptz NOT NULL,
	"welcome_message_id" int8,
	"panel_id" int8,
	FOREIGN KEY("panel_id") REFERENCES panels("message_id") ON DELETE SET NULL ON UPDATE CASCADE,
	PRIMARY KEY("id", "guild_id")
);
CREATE INDEX IF NOT EXISTS tickets_channel_id ON tickets("channel_id");
`
}

func (t *TicketTable) Create(guildId, userId uint64) (id int, err error) {
	query := `INSERT INTO tickets("id", "guild_id", "user_id", "open", "open_time") VALUES((SELECT COALESCE(MAX("id"), 0) + 1 FROM tickets WHERE "guild_id" = $1), $1, $2, true, NOW()) RETURNING "id";`
	err = t.QueryRow(context.Background(), query, guildId, userId).Scan(&id)
	return
}

func (t *TicketTable) SetTicketProperties(guildId uint64, ticketId int, channelId, welcomeMessageId uint64, panelId *uint64) (err error) {
	query := `UPDATE tickets SET "channel_id" = $1, "welcome_message_id" = $2, "panel_id" = $3 WHERE "guild_id" = $4 AND "id" = $5;`
	_, err = t.Exec(context.Background(), query, channelId, welcomeMessageId, panelId, guildId, ticketId)
	return
}

func (t *TicketTable) Get(ticketId int, guildId uint64) (ticket Ticket, e error) {
	query := `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "id" = $1 AND "guild_id" = $2;`

	if err := t.QueryRow(context.Background(), query, ticketId, guildId).Scan(
		&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}
	return
}

func (t *TicketTable) GetByChannel(channelId uint64) (ticket Ticket, e error) {
	query := `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "channel_id" = $1;`

	if err := t.QueryRow(context.Background(), query, channelId).Scan(
		&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}
	return
}

func (t *TicketTable) GetAllByUser(guildId, userId uint64) (tickets []Ticket, e error) {
	query := `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "guild_id" = $1 AND "user_id" = $2;`

	rows, err := t.Query(context.Background(), query, guildId, userId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetOpenByUser(guildId, userId uint64) (tickets []Ticket, e error) {
	query := `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "user_id" = $1 AND "open" = true AND "guild_id" = $2;`

	rows, err := t.Query(context.Background(), query, userId, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
		); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildOpenTickets(guildId uint64) (tickets []Ticket, e error) {
	query := `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "guild_id" = $1 AND "open" = true;`

	rows, err := t.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
		); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildClosedTickets(guildId uint64, limit, before int) (tickets []Ticket, e error) {
	var query string
	var args []interface{}
	if before == 0 {
		query = `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "guild_id" = $1 AND "open" = false
ORDER BY "id" DESC LIMIT $2;`

		args = []interface{}{guildId, limit}
	} else {
		query = `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "guild_id" = $1 AND "open" = false AND "id" < $3
ORDER BY "id" DESC LIMIT $2;`

		args = []interface{}{guildId, limit, before}
	}

	rows, err := t.Query(context.Background(), query, args...)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
		); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetMemberClosedTickets(guildId uint64, userIds []uint64, limit, before int) (tickets []Ticket, e error) {
	// create array of user IDs
	array := &pgtype.Int8Array{}
	if e = array.Set(userIds); e != nil {
		return
	}

	var query string
	var args []interface{}
	if before == 0 {
		query = `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "guild_id" = $1 AND "user_id" = ANY($2) AND "open" = false
ORDER BY "id" DESC LIMIT $3;`

		args = []interface{}{guildId, array, limit}
	} else {
		query = `
SELECT ticket.id, ticket.guild_id, ticket.channel_id, ticket.user_id, ticket.open, ticket.open_time, ticket.welcome_message_id, ticket.panel_id
FROM tickets
WHERE "guild_id" = $1 AND "user_id" = ANY($2) AND "open" = false AND "id" < $4
ORDER BY "id" DESC LIMIT $3;`

		args = []interface{}{guildId, array, limit, before}
	}

	rows, err := t.Query(context.Background(), query, args...)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
		); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetTotalTicketCountInterval(guildId uint64, interval time.Duration) (count int, e error) {
	parsed, err := toInterval(interval)
	if err != nil {
		return 0, err
	}

	query := `SELECT COUNT(*) FROM tickets WHERE "guild_id" = $1 AND tickets.open_time > NOW() - $2::interval;`
	if err := t.QueryRow(context.Background(), query, guildId, parsed).Scan(&count); err != nil {
		e = err
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

func (t *TicketTable) CloseByChannel(channelId uint64) (err error) {
	query := `UPDATE tickets SET "open" = false WHERE "channel_id" = $1;`
	_, err = t.Exec(context.Background(), query, channelId)
	return
}
