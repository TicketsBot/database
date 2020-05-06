package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TicketMembers struct {
	*pgxpool.Pool
}

func newTicketMembers(db *pgxpool.Pool) *TicketMembers {
	return &TicketMembers{
		db,
	}
}

func (m TicketMembers) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS ticket_members("guild_id" int8 NOT NULL REFERENCES tickets("guild_id"), "ticket_id" int4 NOT NULL REFERENCES tickets("id"), "user_id" int8 NOT NULL, PRIMARY KEY("guild_id", "ticket_id", "user_id"));
CREATE INDEX CONCURRENTLY IF NOT EXISTS ticket_members_guild_ticket ON ticket_members("guild_id", "ticket_id");
`
}

func (m *TicketMembers) Get(guildId uint64, ticketId int) (members []uint64, e error) {
	query := `SELECT "user_id" FROM ticket_members WHERE "guild_id" = $1 AND "ticket_id" = $2;`
	rows, err := m.Query(context.Background(), query, guildId, ticketId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			e = err
			continue
		}

		members = append(members, userId)
	}

	return
}

func (m *TicketMembers) Add(guildId uint64, ticketId int, userId uint64) (err error) {
	query := `INSERT INTO ticket_members("guild_id", "ticket_id", "user_id") VALUES($1, $2, $3) ON CONFLICT("guild_id", "ticket_id", "user_id") DO NOTHING;`
	_, err = m.Exec(context.Background(), query, guildId, ticketId, userId)
	return
}

func (m *TicketMembers) Delete(guildId uint64, ticketId int, userId uint64) (err error) {
	query := `DELETE FROM ticket_members WHERE "guild_id"=$1 AND "ticket_id"=$2 AND "user_id"=$3;`
	_, err = m.Exec(context.Background(), query, guildId, ticketId, userId)
	return
}
