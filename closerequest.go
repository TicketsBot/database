package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type CloseRequest struct {
	GuildId   uint64
	TicketId  int
	UserId    uint64
	MessageId uint64
	CloseAt   *time.Time
	Reason    *string
}

type CloseRequestTable struct {
	*pgxpool.Pool
}

func newCloseRequestTable(db *pgxpool.Pool) *CloseRequestTable {
	return &CloseRequestTable{
		db,
	}
}

func (c CloseRequestTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS close_request(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"user_id" int8 NOT NULL,
	"message_id" int8 DEFAULT NULL,
	"close_at" timestamptz,
	"close_reason" VARCHAR(255),
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (c *CloseRequestTable) Get(guildId uint64, ticketId int) (CloseRequest, bool, error) {
	query := `
SELECT "guild_id", "ticket_id", "user_id", "message_id", "close_at", "close_reason"
FROM close_request
WHERE "guild_id" = $1 AND "ticket_id" = $2;
`

	var request CloseRequest
	err := c.QueryRow(context.Background(), query, guildId, ticketId).
		Scan(&request.GuildId, &request.TicketId, &request.UserId, &request.MessageId, &request.CloseAt, &request.Reason)

	if err == nil {
		return request, true, nil
	} else if err == pgx.ErrNoRows {
		return request, false, nil
	} else {
		return request, false, err
	}
}

func (c *CloseReasonTable) GetCloseable() ([]CloseRequest, error) {
	query := `
SELECT close_request.guild_id, close_request.ticket_id, close_request.user_id, close_request.close_at, close_request.close_reason
FROM close_request
INNER JOIN tickets
	ON tickets.guild_id = close_request.guild_id AND tickets.id = close_request.ticket_id
LEFT JOIN auto_close_exclude exclude
	ON close_request.guild_id = exclude.guild_id and close_request.ticket_id = exclude.ticket_id
WHERE
	close_request.close_at < NOW()
	AND
	exclude.guild_id IS NULL
	AND
	tickets.open
;
`

	rows, err := c.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var requests []CloseRequest
	for rows.Next() {
		var request CloseRequest
		if err := rows.Scan(&request.GuildId, &request.TicketId, &request.UserId, &request.MessageId, &request.CloseAt, &request.Reason); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (c *CloseRequestTable) Cleanup() (err error) {
	query := `
DELETE
FROM close_request 
USING tickets
WHERE NOT tickets.open;
`
	_, err = c.Exec(context.Background(), query)
	return
}

func (c *CloseRequestTable) Set(request CloseRequest) (err error) {
	query := `
INSERT INTO close_request("guild_id", "ticket_id", "user_id", "close_at", "close_reason")
VALUES($1, $2, $3, $4, $5)
ON CONFLICT("guild_id", "ticket_id") DO UPDATE 
SET "user_id" = $3, "close_at" = $4, "close_reason" = $5;
`

	_, err = c.Exec(context.Background(), query, request.GuildId, request.TicketId, request.UserId, request.CloseAt, request.Reason)
	return
}

func (c *CloseRequestTable) SetMessageId(guildId uint64, ticketId int, messageId uint64) (err error) {
	query := `
UPDATE close_request
SET "message_id" = $1
WHERE "guild_id" = $2 AND "ticket_id" = $3;
`

	_, err = c.Exec(context.Background(), query, messageId, guildId, ticketId)
	return
}

func (c *CloseRequestTable) Delete(guildId uint64, ticketId int) (uint64, error) {
	query := `
DELETE
FROM close_request
WHERE "guild_id" = $1 AND "ticket_id" = $2
RETURNING message_id
;`

	var temp *uint64
	if err := c.QueryRow(context.Background(), query, guildId, ticketId).Scan(&temp); err != nil {
		return 0, err
	}

	if temp == nil {
		return 0, nil
	} else {
		return *temp, nil
	}
}
