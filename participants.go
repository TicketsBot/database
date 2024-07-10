package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
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

func (p ParticipantTable) Schema() string {
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

func (p *ParticipantTable) GetParticipants(ctx context.Context, guildId uint64, ticketId int) (participants []uint64, err error) {
	query := `
SELECT "user_id"
FROM participant
WHERE "guild_id" = $1 AND "ticket_id" = $2;
`

	rows, err := p.Query(ctx, query, guildId, ticketId)
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

func (p *ParticipantTable) GetTickets(ctx context.Context, userId uint64) (tickets []Participant, err error) {
	query := `
SELECT "guild_id", "ticket_id"
FROM participant
WHERE "user_id" = $1;
`

	rows, err := p.Query(ctx, query, userId)
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

func (p *ParticipantTable) HasParticipated(ctx context.Context, guildId uint64, ticketId int, userId uint64) (hasParticipated bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM participant
	WHERE "guild_id" = $1 AND "ticket_id" = $2 AND "user_id" = $3
);
`

	err = p.QueryRow(ctx, query, guildId, ticketId, userId).Scan(&hasParticipated)
	return
}

func (p *ParticipantTable) Set(ctx context.Context, guildId uint64, ticketId int, userId uint64) (err error) {
	query := `
INSERT INTO participant("guild_id", "ticket_id", "user_id")
VALUES($1, $2, $3)
ON CONFLICT("guild_id", "ticket_id", "user_id")
DO NOTHING;`

	_, err = p.Exec(ctx, query, guildId, ticketId, userId)
	return
}

func (p *ParticipantTable) Delete(ctx context.Context, guildId uint64, ticketId int, userId uint64) (err error) {
	query := `
DELETE FROM participant
WHERE "guild_id"=$1 AND "ticket_id"=$2 AND "user_id"=$3;
`

	_, err = p.Exec(ctx, query, guildId, ticketId, userId)
	return
}

func (p *ParticipantTable) GetParticipatedCount(ctx context.Context, guildId, userId uint64) (count int, err error) {
	query := `
SELECT COUNT(*)
FROM participant
WHERE "guild_id" = $1 AND "user_id" = $2;
`

	err = p.QueryRow(ctx, query, guildId, userId).Scan(&count)
	return
}

func (p *ParticipantTable) GetParticipatedCountInterval(ctx context.Context, guildId, userId uint64, interval time.Duration) (count int, err error) {
	parsed, err := toInterval(interval)
	if err != nil {
		return 0, err
	}

	query := `
SELECT COUNT(*)
FROM participant
INNER JOIN tickets
ON tickets.guild_id = participant.guild_id AND tickets.id = participant.ticket_id
WHERE participant.guild_id = $1 AND participant.user_id = $2 AND tickets.open_time > NOW() - $3::interval;
`

	err = p.QueryRow(ctx, query, guildId, userId, parsed).Scan(&count)
	return
}

// GetParticipatedGlobalWithTranscript returns all closed tickets with a transcript that the user has participated in or opened
func (p *ParticipantTable) GetParticipatedGlobalWithTranscript(ctx context.Context, userId uint64) ([]Participant, error) {
	query := `
(
	SELECT participant.guild_id, participant.ticket_id
    FROM participant 
    INNER JOIN tickets ON tickets.guild_id = participant.guild_id AND tickets.id = participant.ticket_id
    WHERE participant.user_id = 414075045678284810 AND tickets.has_transcript = true
)
UNION
(
    SELECT tickets.guild_id, tickets.id
    FROM tickets
    WHERE tickets.user_id = 414075045678284810 AND tickets.has_transcript = true
);
`

	rows, err := p.Query(ctx, query, userId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var participants []Participant
	for rows.Next() {
		participant := Participant{
			UserId: userId,
		}

		if err := rows.Scan(&participant.GuildId, &participant.TicketId); err != nil {
			return nil, err
		}

		participants = append(participants, participant)
	}

	return participants, nil
}
