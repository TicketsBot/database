package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type GuildMetadata struct {
	OnCallRole *uint64 `json:"on_call_role_id"`
}

func defaultGuildMetadata() GuildMetadata {
	return GuildMetadata{
		OnCallRole: nil,
	}
}

type GuildMetadataTable struct {
	*pgxpool.Pool
}

func newGuildMetadataTable(db *pgxpool.Pool) *GuildMetadataTable {
	return &GuildMetadataTable{
		db,
	}
}

func (s GuildMetadataTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS guild_metadata(
	"guild_id" int8 NOT NULL,
	"on_call_role" int8 DEFAULT NULL,
	PRIMARY KEY("guild_id")
);
`
}

func (s *GuildMetadataTable) Get(guildId uint64) (GuildMetadata, error) {
	query := `
SELECT
	"on_call_role"
FROM guild_metadata
WHERE "guild_id" = $1;
`

	var metadata GuildMetadata
	err := s.QueryRow(context.Background(), query, guildId).Scan(&metadata.OnCallRole)

	if err == nil {
		return metadata, nil
	} else if err == pgx.ErrNoRows {
		return defaultGuildMetadata(), nil
	} else {
		return GuildMetadata{}, err
	}
}

func (s *GuildMetadataTable) Set(guildId uint64, metadata GuildMetadata) (err error) {
	query := `
INSERT INTO guild_metadata(
	"guild_id",
	"on_call_role"
)
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET
    "on_call_role" = $2;
`

	_, err = s.Exec(context.Background(), query,
		guildId,
		metadata.OnCallRole,
	)

	return
}

func (s *GuildMetadataTable) SetOnCallRole(guildId uint64, roleId *uint64) (err error) {
	query := `
INSERT INTO guild_metadata("guild_id", "on_call_role")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "on_call_role" = $2;
`

	_, err = s.Exec(context.Background(), query, guildId, roleId)
	return
}
