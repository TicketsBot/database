package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"time"
)

type Ticket struct {
	Id               int       `json:"id"`
	GuildId          uint64    `json:"guild_id"`
	ChannelId        *uint64   `json:"channel_id"`
	UserId           uint64    `json:"user_id"`
	Open             bool      `json:"open"`
	OpenTime         time.Time `json:"open_time"`
	WelcomeMessageId *uint64   `json:"welcome_message_id"`
	PanelId          *int      `json:"panel_id"`
}

type TicketQueryOptions struct {
	Id      int      `json:"id"`
	GuildId uint64   `json:"guild_id"`
	UserIds []uint64 `json:"user_ids"`
	Open    *bool    `json:"open"`
	Before  int      `json:"before"`
	After   int      `json:"after"`
	Limit   int      `json:"limit"`
}

func (o TicketQueryOptions) HasWhereClause() bool {
	return o.Id == 0 &&
		o.GuildId == 0 &&
		len(o.UserIds) == 0 &&
		o.Open == nil &&
		o.Before == 0 &&
		o.After == 0
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
	"panel_id" int,
	FOREIGN KEY("panel_id") REFERENCES panels("panel_id") ON DELETE SET NULL ON UPDATE CASCADE,
	PRIMARY KEY("id", "guild_id")
);
CREATE INDEX IF NOT EXISTS tickets_channel_id ON tickets("channel_id");
CREATE INDEX IF NOT EXISTS tickets_panel_id ON tickets("panel_id");
`
}

func (t *TicketTable) Create(guildId, userId uint64) (id int, err error) {
	query := `INSERT INTO tickets("id", "guild_id", "user_id", "open", "open_time") VALUES((SELECT COALESCE(MAX("id"), 0) + 1 FROM tickets WHERE "guild_id" = $1), $1, $2, true, NOW()) RETURNING "id";`
	err = t.QueryRow(context.Background(), query, guildId, userId).Scan(&id)
	return
}

func (t *TicketTable) SetTicketProperties(guildId uint64, ticketId int, channelId, welcomeMessageId uint64, panelId *int) (err error) {
	query := `UPDATE tickets SET "channel_id" = $1, "welcome_message_id" = $2, "panel_id" = $3 WHERE "guild_id" = $4 AND "id" = $5;`
	_, err = t.Exec(context.Background(), query, channelId, welcomeMessageId, panelId, guildId, ticketId)
	return
}

func (t *TicketTable) Get(ticketId int, guildId uint64) (ticket Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
FROM tickets
WHERE "id" = $1 AND "guild_id" = $2;`

	if err := t.QueryRow(context.Background(), query, ticketId, guildId).Scan(
		&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}
	return
}

func (t *TicketTable) GetByOptions(options TicketQueryOptions) (tickets []Ticket, e error) {
	query, args, err := options.BuildQuery()
	if err != nil {
		return nil, err
	}

	rows, err := t.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ticket Ticket
		err = rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId,
		)

		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (o TicketQueryOptions) BuildQuery() (query string, args []interface{}, _err error) {
	query = `SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id FROM tickets`

	if !o.HasWhereClause() {
		query += " WHERE "
	}

	var needsAnd bool

	if o.Id != 0 {
		args = append(args, o.Id)
		query += fmt.Sprintf(`"id" = $%d`, len(args))
		needsAnd = true
	}

	if o.GuildId != 0 {
		if needsAnd {
			query += " AND "
		}

		args = append(args, o.GuildId)
		query += fmt.Sprintf(`"guild_id" = $%d`, len(args))
		needsAnd = true
	}

	if len(o.UserIds) > 0 {
		if needsAnd {
			query += " AND "
		}

		userIdArray := &pgtype.Int8Array{}
		if err := userIdArray.Set(o.UserIds); err != nil {
			return "", nil, err
		}

		args = append(args, userIdArray)
		query += fmt.Sprintf(`"user_id" = ANY($%d)`, len(args))
		needsAnd = true
	}

	if o.Open != nil {
		if needsAnd {
			query += " AND "
		}

		args = append(args, *o.Open)
		query += fmt.Sprintf(`"open" = $%d`, len(args))
		needsAnd = true
	}

	if o.Before != 0 {
		if needsAnd {
			query += " AND "
		}

		args = append(args, o.Before)
		query += fmt.Sprintf(`"id" < $%d`, len(args))
		needsAnd = true
	}

	if o.After != 0 {
		if needsAnd {
			query += " AND "
		}

		args = append(args, o.After)
		query += fmt.Sprintf(`"id" > $%d`, len(args))
		needsAnd = true
	}

	query += ` ORDER BY "id" DESC `

	if o.Limit != 0 {
		args = append(args, o.Limit)
		query += fmt.Sprintf(` LIMIT $%d`, len(args))
	}

	query += ";"
	return
}

func (t *TicketTable) GetByChannel(channelId uint64) (ticket Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
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
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
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
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
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

func (t *TicketTable) GetClosedByAnyBefore(guildId uint64, userIds []uint64, before, limit int) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
FROM tickets
WHERE "guild_id" = $1 "user_id" = ANY($2) AND "open" = false AND "id" < $3
ORDER BY "id" DESC
LIMIT $4;`

	userIdArray := &pgtype.Int8Array{}
	if err := userIdArray.Set(userIds); err != nil {
		return nil, err
	}

	rows, err := t.Query(context.Background(), query, guildId, userIdArray, before, limit)
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

type TicketWithCloseReason struct {
	Ticket
	CloseReason *string `json:"close_reason"`
}

func (t *TicketTable) GetClosedByAnyBeforeWithCloseReason(guildId uint64, userIds []uint64, before, limit int) (tickets []TicketWithCloseReason, e error) {
	query := `
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, close_reason.close_reason
FROM tickets
LEFT JOIN close_reason
ON tickets.id = close_reason.ticket_id AND tickets.guild_id = close_reason.guild_id
WHERE tickets.guild_id = $1 AND tickets.user_id = ANY($2) AND tickets.open = false AND tickets.id < $3
ORDER BY tickets.id DESC
LIMIT $4;`

	userIdArray := &pgtype.Int8Array{}
	if err := userIdArray.Set(userIds); err != nil {
		return nil, err
	}

	if before <= 0 {
		before = math.MaxInt32
	}

	rows, err := t.Query(context.Background(), query, guildId, userIdArray, before, limit)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket TicketWithCloseReason
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId, &ticket.CloseReason,
		); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetClosedByAnyAfterWithCloseReason(guildId uint64, userIds []uint64, after, limit int) (tickets []TicketWithCloseReason, e error) {
	query := `
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, close_reason.close_reason
FROM tickets
LEFT JOIN close_reason
ON tickets.id = close_reason.ticket_id AND tickets.guild_id = close_reason.guild_id
WHERE tickets.guild_id = $1 AND tickets.user_id = ANY($2) AND tickets.open = false AND tickets.id > $3
ORDER BY tickets.id ASC
LIMIT $4;`

	userIdArray := &pgtype.Int8Array{}
	if err := userIdArray.Set(userIds); err != nil {
		return nil, err
	}

	rows, err := t.Query(context.Background(), query, guildId, userIdArray, after, limit)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket TicketWithCloseReason
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId, &ticket.CloseReason,
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
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
FROM tickets
WHERE "guild_id" = $1 AND "open" = true
ORDER BY id DESC;`

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
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
FROM tickets
WHERE "guild_id" = $1 AND "open" = false AND "id" < $3
ORDER BY "id" DESC LIMIT $2;`

	if before <= 0 {
		before = math.MaxInt32
	}

	rows, err := t.Query(context.Background(), query, guildId, limit, before)
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

func (t *TicketTable) GetGuildClosedTicketsBeforeWithCloseReason(guildId uint64, limit, before int) (tickets []TicketWithCloseReason, e error) {
	query := `
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, close_reason.close_reason
FROM tickets
LEFT JOIN close_reason
ON tickets.id = close_reason.ticket_id AND tickets.guild_id = close_reason.guild_id
WHERE tickets.guild_id = $1 AND tickets.open = false AND tickets.id < $3
ORDER BY tickets.id DESC LIMIT $2;`

	if before <= 0 {
		before = math.MaxInt32
	}

	rows, err := t.Query(context.Background(), query, guildId, limit, before)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket TicketWithCloseReason
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId, &ticket.CloseReason,
		); err != nil {
			e = err
			continue
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildClosedTicketsAfterWithCloseReason(guildId uint64, limit, after int) (tickets []TicketWithCloseReason, e error) {
	query := `
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, close_reason.close_reason
FROM tickets
LEFT JOIN close_reason
ON tickets.id = close_reason.ticket_id AND tickets.guild_id = close_reason.guild_id
WHERE tickets.guild_id = $1 AND tickets.open = false AND tickets.id > $3
ORDER BY tickets.id ASC LIMIT $2;`

	rows, err := t.Query(context.Background(), query, guildId, limit, after)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var ticket TicketWithCloseReason
		if err := rows.Scan(
			&ticket.Id, &ticket.GuildId, &ticket.ChannelId, &ticket.UserId, &ticket.Open, &ticket.OpenTime, &ticket.WelcomeMessageId, &ticket.PanelId, &ticket.CloseReason,
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

	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id
FROM tickets
WHERE "guild_id" = $1 AND "user_id" = ANY($2) AND "open" = false AND "id" < $4
ORDER BY "id" DESC LIMIT $3;
`

	if before <= 0 {
		before = math.MaxInt32
	}

	rows, err := t.Query(context.Background(), query, guildId, userIds, limit, before)
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
