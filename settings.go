package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO: Migrate all settings to this table
type Settings struct {
	HideClaimButton bool `json:"hide_claim_button"`
}

func defaultSettings() Settings {
	return Settings{
		HideClaimButton: false,
	}
}

type SettingsTable struct {
	*pgxpool.Pool
}

func newSettingsTable(db *pgxpool.Pool) *SettingsTable {
	return &SettingsTable{
		db,
	}
}

func (s SettingsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS settings(
	"guild_id" int8 NOT NULL,
	"hide_claim_button" bool DEFAULT 'f',
	PRIMARY KEY("guild_id")
);
`
}

func (s *SettingsTable) Get(guildId uint64) (Settings, error) {
	query := `
SELECT "hide_claim_button"
FROM settings
WHERE "guild_id" = $1;
`

	var settings Settings
	err := s.QueryRow(context.Background(), query, guildId).Scan(
		&settings.HideClaimButton,
	)

	if err == nil {
		return settings, nil
	} else if err == pgx.ErrNoRows {
		return defaultSettings(), nil
	} else {
		return settings, err
	}
}

func (s *SettingsTable) Set(guildId uint64, settings Settings) (err error) {
	query := `
INSERT INTO settings("guild_id", "hide_claim_button")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "hide_claim_button" = $2;
`

	_, err = s.Exec(context.Background(), query, guildId, settings.HideClaimButton)
	return
}

func (s *SettingsTable) SetHideClaimButton(guildId uint64, hideClaimButton bool) (err error) {
	query := `
INSERT INTO settings("guild_id", "hide_claim_button")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "hide_claim_button" = $2;
`

	_, err = s.Exec(context.Background(), query, guildId, hideClaimButton)
	return
}
