package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TicketPermissions struct {
	AttachFiles  bool `json:"attach_files"`
	EmbedLinks   bool `json:"embed_links"`
	AddReactions bool `json:"add_reactions"`
}

type TicketPermissionsTable struct {
	*pgxpool.Pool
}

func newTicketPermissionsTable(db *pgxpool.Pool) *TicketPermissionsTable {
	return &TicketPermissionsTable{
		db,
	}
}

func (c TicketPermissionsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS ticket_permissions(
	"guild_id" int8 NOT NULL,
	"attach_files" bool NOT NULL DEFAULT 't',
	"embed_links" bool NOT NULL DEFAULT 't',
	"add_reactions" bool NOT NULL DEFAULT 't',
	PRIMARY KEY("guild_id")
);
`
}

func (c *TicketPermissionsTable) Get(ctx context.Context, guildId uint64) (TicketPermissions, error) {
	query := `
SELECT "attach_files", "embed_links", "add_reactions"
FROM ticket_permissions
WHERE "guild_id" = $1;`

	var permissions TicketPermissions
	err := c.QueryRow(ctx, query, guildId).Scan(
		&permissions.AttachFiles,
		&permissions.EmbedLinks,
		&permissions.AddReactions,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return TicketPermissions{
				AttachFiles:  true,
				EmbedLinks:   true,
				AddReactions: true,
			}, nil
		} else {
			return TicketPermissions{}, err
		}
	}

	return permissions, nil
}

func (c *TicketPermissionsTable) Set(ctx context.Context, guildId uint64, permissions TicketPermissions) (err error) {
	query := `
INSERT INTO ticket_permissions("guild_id", "attach_files", "embed_links", "add_reactions")
VALUES($1, $2, $3, $4)
ON CONFLICT("guild_id") DO UPDATE SET "attach_files" = $2, "embed_links" = $3, "add_reactions" = $4;`

	_, err = c.Exec(ctx, query, guildId, permissions.AttachFiles, permissions.EmbedLinks, permissions.AddReactions)
	return
}

func (c *TicketPermissionsTable) Delete(ctx context.Context, guildId uint64) error {
	query := `DELETE FROM ticket_permissions WHERE "guild_id"=$1;`
	_, err := c.Exec(ctx, query, guildId)
	return err
}
