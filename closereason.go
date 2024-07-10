package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CloseMetadata struct {
	Reason   *string `json:"reason"`
	ClosedBy *uint64 `json:"closed_by,string"` // Null if auto-closed
}

type CloseMetadataTable struct {
	*pgxpool.Pool
}

func newCloseReasonTable(db *pgxpool.Pool) *CloseMetadataTable {
	return &CloseMetadataTable{
		db,
	}
}

func (c CloseMetadataTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS close_reason(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"close_reason" TEXT,
	"closed_by" int8,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (c *CloseMetadataTable) Get(ctx context.Context, guildId uint64, ticketId int) (CloseMetadata, bool, error) {
	query := `
SELECT "close_reason", "closed_by"
FROM close_reason
WHERE "guild_id" = $1 AND "ticket_id" = $2;
`

	var data CloseMetadata
	if err := c.QueryRow(ctx, query, guildId, ticketId).Scan(&data.Reason, &data.ClosedBy); err != nil {
		if err == pgx.ErrNoRows {
			return CloseMetadata{}, false, nil
		} else {
			return CloseMetadata{}, false, err
		}
	}

	return data, true, nil
}

func (c *CloseMetadataTable) GetMulti(ctx context.Context, guildId uint64, ticketIds []int) (map[int]CloseMetadata, error) {
	query := `
SELECT "ticket_id", "close_reason", "closed_by"
FROM close_reason
WHERE "guild_id" = $1 AND "ticket_id" = ANY($2);
`

	array := &pgtype.Int4Array{}
	if err := array.Set(ticketIds); err != nil {
		return nil, err
	}

	rows, err := c.Query(ctx, query, guildId, ticketIds)
	if err != nil {
		return nil, err
	}

	ticketMetadata := make(map[int]CloseMetadata)
	for rows.Next() {
		var ticketId int
		var data CloseMetadata
		if err := rows.Scan(&ticketId, &data.Reason, &data.ClosedBy); err != nil {
			return nil, err
		}

		ticketMetadata[ticketId] = data
	}

	return ticketMetadata, nil
}

func (c *CloseMetadataTable) Set(ctx context.Context, guildId uint64, ticketId int, data CloseMetadata) (err error) {
	query := `
INSERT INTO close_reason("guild_id", "ticket_id", "close_reason", "closed_by")
VALUES($1, $2, $3, $4)
ON CONFLICT("guild_id", "ticket_id") DO UPDATE SET "close_reason" = $3, "closed_by" = $4;
`

	_, err = c.Exec(ctx, query, guildId, ticketId, data.Reason, data.ClosedBy)
	return
}

func (c *CloseMetadataTable) Delete(ctx context.Context, guildId uint64, ticketId int) (err error) {
	query := `DELETE FROM close_reason WHERE "guild_id"=$1 AND "ticket_id"=$2;`
	_, err = c.Exec(ctx, query, guildId, ticketId)
	return
}
