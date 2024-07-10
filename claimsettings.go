package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ClaimSettings struct {
	SupportCanView bool `json:"support_can_view"`
	SupportCanType bool `json:"support_can_type"`
}

var defaultClaimSettings = ClaimSettings{
	SupportCanView: true,
	SupportCanType: false,
}

type ClaimSettingsTable struct {
	*pgxpool.Pool
}

func newClaimSettingsTable(db *pgxpool.Pool) *ClaimSettingsTable {
	return &ClaimSettingsTable{
		db,
	}
}

func (c ClaimSettingsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS claim_settings(
	"guild_id" int8 NOT NULL,
	"support_can_view" bool NOT NULL,
	"support_can_type" bool NOT NULL,
	PRIMARY KEY("guild_id")
);
`
}

func (c *ClaimSettingsTable) Get(ctx context.Context, guildId uint64) (settings ClaimSettings, e error) {
	query := `SELECT "support_can_view", "support_can_type" FROM claim_settings WHERE "guild_id" = $1;`
	if err := c.QueryRow(ctx, query, guildId).Scan(&settings.SupportCanView, &settings.SupportCanType); err != nil {
		if err == pgx.ErrNoRows {
			settings = defaultClaimSettings
		} else {
			e = err
		}
	}

	return
}

func (c *ClaimSettingsTable) Set(ctx context.Context, guildId uint64, settings ClaimSettings) (err error) {
	query := `
INSERT INTO claim_settings("guild_id", "support_can_view", "support_can_type") VALUES($1, $2, $3)
	ON CONFLICT("guild_id") DO UPDATE SET
	"support_can_view" = $2,
	"support_can_type" = $3;`

	_, err = c.Exec(ctx, query, guildId, settings.SupportCanView, settings.SupportCanType)
	return
}
