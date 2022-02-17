package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type TicketLastMessageTable struct {
	*pgxpool.Pool
}

type TicketLastMessage struct {
	LastMessageId   *uint64
	LastMessageTime *time.Time
}

func newTicketLastMessageTable(db *pgxpool.Pool) *TicketLastMessageTable {
	return &TicketLastMessageTable{
		db,
	}
}

func (m TicketLastMessageTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS ticket_last_message(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"last_message_id" int8,
	"last_message_time" timestamptz,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (m *TicketLastMessageTable) Get(guildId uint64, ticketId int) (lastMessage TicketLastMessage, e error) {
	query := `SELECT "last_message_id", "last_message_time" FROM ticket_last_message WHERE "guild_id" = $1 AND "ticket_id" = $2;`
	if err := m.QueryRow(context.Background(), query, guildId, ticketId).Scan(&lastMessage.LastMessageId, &lastMessage.LastMessageTime); err != nil && err != pgx.ErrNoRows { // defaults to nil if no rows
		e = err
	}

	return
}

func (m *TicketLastMessageTable) Set(guildId uint64, ticketId int, messageId uint64) (err error) {
	query := `
INSERT INTO ticket_last_message("guild_id", "ticket_id", "last_message_id", "last_message_time")
VALUES($1, $2, $3, NOW()) ON CONFLICT("guild_id", "ticket_id")
DO UPDATE SET "last_message_id" = $3, "last_message_time" = NOW();`

	_, err = m.Exec(context.Background(), query, guildId, ticketId, messageId)
	return
}

func (m *TicketLastMessageTable) Delete(guildId uint64, ticketId int) (err error) {
	query := `DELETE FROM ticket_last_message WHERE "guild_id"=$1 AND "ticket_id"=$2;`
	_, err = m.Exec(context.Background(), query, guildId, ticketId)
	return
}
