package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ParticipantTable struct {
	*pgxpool.Pool
}

type Participant struct {
	GuildId  uint64
	TicketId int
	UserId   uint64
}

func newParticipantTable(db *pgxpool.Pool) *ParticipantTable {
	return &ParticipantTable{
		db,
	}
}

func (m ParticipantTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS participant(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"user_id" int8 NOT NULL,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id", "user_id")
);
`
}

func (m *ParticipantTable) GetParticipants(guildId uint64, ticketId int) (participants []uint64, err error) {
	query := `
SELECT "user_id"
FROM participant
WHERE "guild_id" = $1 AND "ticket_id" = $2;
`

	rows, err := m.Query(context.Background(), query, guildId, ticketId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userId uint64

		if err = rows.Scan(&userId); err != nil {
			return
		}

		participants = append(participants, userId)
	}

	return
}

func (m *ParticipantTable) GetTickets(userId uint64) (tickets []Participant, err error) {
	query := `
SELECT "guild_id", "ticket_id"
FROM participant
WHERE "user_id" = $1;
`

	rows, err := m.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		participant := Participant{
			UserId: userId,
		}

		if err = rows.Scan(&participant.GuildId, &participant.TicketId); err != nil {
			return
		}

		tickets = append(tickets, participant)
	}

	return
}

func (m *ParticipantTable) Set(guildId uint64, ticketId int, userId uint64) (err error) {
	query := `
INSERT INTO participant("guild_id", "ticket_id", "user_id")
VALUES($1, $2, $3)
ON CONFLICT("guild_id", "ticket_id", "user_id")
DO NOTHING;`

	_, err = m.Exec(context.Background(), query, guildId, ticketId, userId)
	return
}

func (m *ParticipantTable) Delete(guildId uint64, ticketId int, userId uint64) (err error) {
	query := `
DELETE FROM participant
WHERE "guild_id"=$1 AND "ticket_id"=$2 AND "user_id"=$3;
`

	_, err = m.Exec(context.Background(), query, guildId, ticketId, userId)
	return
}
