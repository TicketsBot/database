package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"time"
)

type Ticket struct {
	Id               int        `json:"id"`
	GuildId          uint64     `json:"guild_id"`
	ChannelId        *uint64    `json:"channel_id"`
	UserId           uint64     `json:"user_id"`
	Open             bool       `json:"open"`
	OpenTime         time.Time  `json:"open_time"`
	WelcomeMessageId *uint64    `json:"welcome_message_id"`
	PanelId          *int       `json:"panel_id"`
	HasTranscript    bool       `json:"has_transcript"`
	CloseTime        *time.Time `json:"close_time"`
	IsThread         bool       `json:"is_thread"`
	JoinMessageId    *uint64    `json:"join_message_id"`
	NotesThreadId    *uint64    `json:"notes_thread_id"`
}

type TicketQueryOptions struct {
	Id      int       `json:"id"`
	GuildId uint64    `json:"guild_id"`
	UserIds []uint64  `json:"user_ids"`
	Open    *bool     `json:"open"`
	PanelId int       `json:"panel_id"`
	Rating  int       `json:"rating"`
	Order   OrderType `json:"order_type"`
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
}

type OrderType string

const (
	OrderTypeNone       OrderType = ""
	OrderTypeAscending  OrderType = "ASC"
	OrderTypeDescending OrderType = "DESC"
)

func (o TicketQueryOptions) HasWhereClause() bool {
	return o.Id == 0 &&
		o.GuildId == 0 &&
		len(o.UserIds) == 0 &&
		o.Open == nil &&
		o.Rating != 0
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
	"has_transcript" bool NOT NULL DEFAULT 'f',
	"close_time" timestamptz DEFAULT NULL,
    "is_thread" bool NOT NULL DEFAULT 'f',
    "join_message_id" int8 DEFAULT NULL,
    "notes_thread_id" int8 DEFAULT NULL,
	FOREIGN KEY("panel_id") REFERENCES panels("panel_id") ON DELETE SET NULL ON UPDATE CASCADE,
	PRIMARY KEY("id", "guild_id")
);
CREATE INDEX IF NOT EXISTS tickets_channel_id ON tickets("channel_id");
CREATE INDEX IF NOT EXISTS tickets_panel_id ON tickets("panel_id");
`
}

func (t *TicketTable) Create(guildId, userId uint64, isThread bool, panelId *uint64) (id int, err error) {
	query := `
INSERT INTO tickets("id", "guild_id", "user_id", "open", "open_time", "is_thread", "panel_id")
VALUES(
       (SELECT COALESCE(MAX("id"), 0) + 1 FROM tickets WHERE "guild_id" = $1), 
       $1, $2, true, NOW(), $3, $4
)
RETURNING "id";`

	err = t.QueryRow(context.Background(), query, guildId, userId, isThread, panelId).Scan(&id)
	return
}

func (t *TicketTable) SetChannelId(guildId uint64, ticketId int, channelId uint64) (err error) {
	query := `
UPDATE tickets
SET "channel_id" = $1
WHERE "guild_id" = $2 AND "id" = $3;`

	_, err = t.Exec(context.Background(), query, channelId, guildId, ticketId)
	return
}

func (t *TicketTable) SetMessageIds(guildId uint64, ticketId int, welcomeMessageId uint64, joinMessageId *uint64) (err error) {
	query := `
UPDATE tickets
SET "welcome_message_id" = $1, "join_message_id" = $2
WHERE "guild_id" = $3 AND "id" = $4;`

	_, err = t.Exec(context.Background(), query, welcomeMessageId, joinMessageId, guildId, ticketId)
	return
}

func (t *TicketTable) Get(ticketId int, guildId uint64) (ticket Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
FROM tickets
WHERE "id" = $1 AND "guild_id" = $2;`

	if err := t.QueryRow(context.Background(), query, ticketId, guildId).Scan(
		&ticket.Id,
		&ticket.GuildId,
		&ticket.ChannelId,
		&ticket.UserId,
		&ticket.Open,
		&ticket.OpenTime,
		&ticket.WelcomeMessageId,
		&ticket.PanelId,
		&ticket.HasTranscript,
		&ticket.CloseTime,
		&ticket.IsThread,
		&ticket.JoinMessageId,
		&ticket.NotesThreadId,
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		)

		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (o TicketQueryOptions) BuildQuery() (query string, args []interface{}, _err error) {
	query = `
SELECT tickets.id,
	tickets.guild_id,
	tickets.channel_id,
	tickets.user_id,
	tickets.open,
	tickets.open_time,
	tickets.welcome_message_id,
	tickets.panel_id,
	tickets.has_transcript,
	tickets.close_time,
	tickets.is_thread,
	tickets.join_message_id,
	tickets.notes_thread_id
FROM tickets`

	if o.Rating != 0 {
		query += " INNER JOIN service_ratings ON tickets.guild_id = service_ratings.guild_id AND tickets.id = service_ratings.ticket_id "
	}

	if !o.HasWhereClause() {
		query += " WHERE "
	}

	var needsAnd bool

	if o.Id != 0 {
		args = append(args, o.Id)
		query += fmt.Sprintf(`tickets.id = $%d`, len(args))
		needsAnd = true
	}

	if o.GuildId != 0 {
		if needsAnd {
			query += " AND "
		}

		args = append(args, o.GuildId)
		query += fmt.Sprintf(`tickets.guild_id = $%d`, len(args))
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
		query += fmt.Sprintf(`tickets.user_id = ANY($%d)`, len(args))
		needsAnd = true
	}

	if o.Open != nil {
		if needsAnd {
			query += " AND "
		}

		args = append(args, *o.Open)
		query += fmt.Sprintf(`tickets.open = $%d`, len(args))
		needsAnd = true
	}

	if o.PanelId > 0 {
		if needsAnd {
			query += " AND "
		}

		args = append(args, o.PanelId)
		query += fmt.Sprintf(`tickets.panel_id = $%d`, len(args))
		needsAnd = true
	}

	if o.Rating > 0 {
		if needsAnd {
			query += " AND "
		}

		args = append(args, o.Rating)
		query += fmt.Sprintf(`service_ratings.rating = $%d`, len(args))
		needsAnd = true
	}

	// Cannot use prepared statement for this value
	if o.Order == OrderTypeAscending || o.Order == OrderTypeDescending {
		query += fmt.Sprintf(` ORDER BY "id" %s `, o.Order)
	}

	if o.Limit != 0 {
		args = append(args, o.Limit)
		query += fmt.Sprintf(` LIMIT $%d `, len(args))
	}

	if o.Offset != 0 {
		args = append(args, o.Offset)
		query += fmt.Sprintf(` OFFSET $%d `, len(args))
	}

	query += ";"
	return
}

func (t *TicketTable) GetByChannel(ctx context.Context, channelId uint64) (Ticket, bool, error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
FROM tickets
WHERE "channel_id" = $1;`

	var ticket Ticket
	if err := t.QueryRow(ctx, query, channelId).Scan(
		&ticket.Id,
		&ticket.GuildId,
		&ticket.ChannelId,
		&ticket.UserId,
		&ticket.Open,
		&ticket.OpenTime,
		&ticket.WelcomeMessageId,
		&ticket.PanelId,
		&ticket.HasTranscript,
		&ticket.CloseTime,
		&ticket.IsThread,
		&ticket.JoinMessageId,
		&ticket.NotesThreadId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Ticket{}, false, nil
		} else {
			return Ticket{}, false, err
		}
	}

	return ticket, true, nil
}

func (t *TicketTable) GetByChannelAndGuild(ctx context.Context, channelId, guildId uint64) (ticket Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
FROM tickets
WHERE "channel_id" = $1 AND "guild_id" = $2;`

	if err := t.QueryRow(ctx, query, channelId, guildId).Scan(
		&ticket.Id,
		&ticket.GuildId,
		&ticket.ChannelId,
		&ticket.UserId,
		&ticket.Open,
		&ticket.OpenTime,
		&ticket.WelcomeMessageId,
		&ticket.PanelId,
		&ticket.HasTranscript,
		&ticket.CloseTime,
		&ticket.IsThread,
		&ticket.JoinMessageId,
		&ticket.NotesThreadId,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}
	return
}

func (t *TicketTable) GetAllByUser(guildId, userId uint64) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
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
		if err := rows.Scan(
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			e = err
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetTotalCountByUser(guildId, userId uint64) (int, error) {
	query := `
SELECT COUNT(id)
FROM tickets
WHERE "guild_id" = $1 AND "user_id" = $2;`

	var count int
	if err := t.QueryRow(context.Background(), query, guildId, userId).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (t *TicketTable) GetOpenByUser(guildId, userId uint64) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetOpenCountByUser(guildId, userId uint64) (int, error) {
	query := `
SELECT COUNT(id)
FROM tickets
WHERE "user_id" = $1 AND "open" = true AND "guild_id" = $2;`

	var count int
	if err := t.QueryRow(context.Background(), query, userId, guildId).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (t *TicketTable) GetClosedByUserPrefixed(guildId, userId uint64, prefix string, limit int) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
FROM tickets
WHERE "user_id" = $1 AND "open" = false AND "guild_id" = $2 AND id::TEXT LIKE $3 || '%'
ORDER BY "id" DESC
LIMIT $4;`

	rows, err := t.Query(context.Background(), query, userId, guildId, prefix, limit)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetClosedByAnyBefore(guildId uint64, userIds []uint64, before, limit int) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
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
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, tickets.has_transcript, tickets.close_time, tickets.is_thread, tickets.join_message_id, tickets.notes_thread_id, close_reason.close_reason
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
			&ticket.CloseReason,
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
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, tickets.has_transcript, tickets.close_time, tickets.is_thread, tickets.join_message_id, tickets.notes_thread_id, close_reason.close_reason
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.CloseReason,
			&ticket.NotesThreadId,
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
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildOpenTicketsExcludeThreads(guildId uint64) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
FROM tickets
WHERE "guild_id" = $1 AND "open" = true AND "is_thread" = false
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildClosedTickets(guildId uint64, limit, before int) (tickets []Ticket, e error) {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetGuildClosedTicketsBeforeWithCloseReason(guildId uint64, limit, before int) (tickets []TicketWithCloseReason, e error) {
	query := `
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, tickets.has_transcript, tickets.close_time, tickets.is_thread, tickets.join_message_id, tickets.notes_thread_id, close_reason.close_reason
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
			&ticket.CloseReason,
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
SELECT tickets.id, tickets.guild_id, tickets.channel_id, tickets.user_id, tickets.open, tickets.open_time, tickets.welcome_message_id, tickets.panel_id, tickets.has_transcript, tickets.close_time, tickets.is_thread, tickets.join_message_id, tickets.notes_thread_id, close_reason.close_reason
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
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
			&ticket.CloseReason,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return
}

func (t *TicketTable) GetMemberClosedTickets(guildId uint64, userIds []uint64, limit, before int) ([]Ticket, error) {
	// create array of user IDs
	array := &pgtype.Int8Array{}
	if err := array.Set(userIds); err != nil {
		return nil, err
	}

	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id
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
		return nil, err
	}

	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
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
	query := `UPDATE tickets SET "open"=false, "close_time"=NOW() WHERE "id"=$1 AND "guild_id"=$2;`
	_, err = t.Exec(context.Background(), query, ticketId, guildId)
	return
}

func (t *TicketTable) CloseByChannel(channelId uint64) (err error) {
	query := `UPDATE tickets SET "open" = false, "close_time" = NOW() WHERE "channel_id" = $1;`
	_, err = t.Exec(context.Background(), query, channelId)
	return
}

func (t *TicketTable) SetHasTranscript(guildId uint64, ticketId int, hasTranscript bool) (err error) {
	query := `UPDATE tickets SET "has_transcript" = $3 WHERE "guild_id" = $1 AND "id" = $2;`
	_, err = t.Exec(context.Background(), query, guildId, ticketId, hasTranscript)
	return
}

func (t *TicketTable) SetPanelId(guildId uint64, ticketId, panelId int) (err error) {
	query := `UPDATE tickets SET "panel_id" = $3 WHERE "guild_id" = $1 AND "id" = $2;`
	_, err = t.Exec(context.Background(), query, guildId, ticketId, panelId)
	return
}

func (t *TicketTable) SetOpen(guildId uint64, ticketId int) (err error) {
	query := `UPDATE tickets SET "open" = TRUE, "close_time" = NULL WHERE "guild_id" = $1 AND "id" = $2;`
	_, err = t.Exec(context.Background(), query, guildId, ticketId)
	return
}

func (t *TicketTable) SetJoinMessageId(guildId uint64, ticketId int, joinMessageId *uint64) (err error) {
	query := `UPDATE tickets SET "join_message_id" = $3 WHERE "guild_id" = $1 AND "id" = $2;`
	_, err = t.Exec(context.Background(), query, guildId, ticketId, joinMessageId)
	return
}

func (t *TicketTable) SetNotesThreadId(guildId uint64, ticketId int, notesThreadId uint64) error {
	query := `UPDATE tickets SET "notes_thread_id" = $3 WHERE "guild_id" = $1 AND "id" = $2;`

	_, err := t.Exec(context.Background(), query, guildId, ticketId, notesThreadId)
	return err
}
