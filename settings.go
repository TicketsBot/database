package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO: Migrate all settings to this table
type Settings struct {
	HideClaimButton    bool `json:"hide_claim_button"`
	DisableOpenCommand bool `json:"disable_open_command"`
}

func defaultSettings() Settings {
	return Settings{
		HideClaimButton:    false,
		DisableOpenCommand: false,
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
	"disable_open_command" bool DEFAULT 'f',
	PRIMARY KEY("guild_id")
);
`
}

func (s *SettingsTable) Get(guildId uint64) (Settings, error) {
	query := `
SELECT "hide_claim_button", "disable_open_command"
FROM settings
WHERE "guild_id" = $1;
`

	var settings Settings
	err := s.QueryRow(context.Background(), query, guildId).Scan(
		&settings.HideClaimButton,
		&settings.DisableOpenCommand,
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
INSERT INTO settings("guild_id", "hide_claim_button", "disable_open_command")
VALUES($1, $2, $3)
ON CONFLICT("guild_id")
DO UPDATE SET
	"hide_claim_button" = $2,
	"disable_open_command" = $3;
`

	_, err = s.Exec(context.Background(), query, guildId, settings.HideClaimButton, settings.DisableOpenCommand)
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

func (s *SettingsTable) SetDisableOpenCommand(guildId uint64, disableOpenCommand bool) (err error) {
	query := `
INSERT INTO settings("guild_id", "disable_open_command")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "disable_open_command" = $2;
`

	_, err = s.Exec(context.Background(), query, guildId, disableOpenCommand)
	return
}
