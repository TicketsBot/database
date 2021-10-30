package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO: Migrate all settings to this table
type Settings struct {
	HideClaimButton            bool `json:"hide_claim_button"`
	DisableOpenCommand         bool `json:"disable_open_command"`
	ContextMenuPermissionLevel int  `json:"context_menu_permission_level,string"`
	ContextMenuAddSender       bool `json:"context_menu_add_sender"`
	ContextMenuPanel           *int `json:"context_menu_panel"`
	StoreTranscripts           bool `json:"store_transcripts"`
	UseThreads                 bool `json:"use_threads"`
	ThreadArchiveDuration      int  `json:"thread_archive_duration"`
}

func defaultSettings() Settings {
	return Settings{
		HideClaimButton:            false,
		DisableOpenCommand:         false,
		ContextMenuPermissionLevel: 0,
		ContextMenuAddSender:       true,
		ContextMenuPanel:           nil,
		StoreTranscripts:           true,
		UseThreads:                 false,
		ThreadArchiveDuration:      10080,
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
	"context_menu_permission_level" int DEFAULT '0',
	"context_menu_add_sender" bool DEFAULT 't',
	"context_menu_panel" int DEFAULT NULL,
	"store_transcripts" bool DEFAULT 't',
    "use_threads" bool DEFAULT 'f',
    "thread_archive_duration" int DEFAULT '10080',
	FOREIGN KEY("context_menu_panel") REFERENCES panels("panel_id") ON DELETE SET NULL,
	PRIMARY KEY("guild_id")
);
`
}

// TODO: GetSome func
func (s *SettingsTable) Get(guildId uint64) (Settings, error) {
	query := `
SELECT
	"hide_claim_button",
	"disable_open_command",
	"context_menu_permission_level",
	"context_menu_add_sender",
	"context_menu_panel",
	"store_transcripts",
    "use_threads",
    "thread_archive_duration"
FROM settings
WHERE "guild_id" = $1;
`

	var settings Settings
	err := s.QueryRow(context.Background(), query, guildId).Scan(
		&settings.HideClaimButton,
		&settings.DisableOpenCommand,
		&settings.ContextMenuPermissionLevel,
		&settings.ContextMenuAddSender,
		&settings.ContextMenuPanel,
		&settings.StoreTranscripts,
		&settings.UseThreads,
        &settings.ThreadArchiveDuration,
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
INSERT INTO settings(
	"guild_id",
	"hide_claim_button",
	"disable_open_command",
	"context_menu_permission_level",
	"context_menu_add_sender",
	"context_menu_panel",
	"store_transcripts",
    "use_threads",
    "thread_archive_duration"
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT("guild_id")
DO UPDATE SET
	"hide_claim_button" = $2,
	"disable_open_command" = $3,
	"context_menu_permission_level" = $4,
	"context_menu_add_sender" = $5,
	"context_menu_panel" = $6,
	"store_transcripts" = $7,
    "use_threads" = $8,
    "thread_archive_duration" = $9;
;
`

	_, err = s.Exec(context.Background(), query,
		guildId,
		settings.HideClaimButton,
		settings.DisableOpenCommand,
		settings.ContextMenuPermissionLevel,
		settings.ContextMenuAddSender,
		settings.ContextMenuPanel,
		settings.StoreTranscripts,
		settings.UseThreads,
        settings.ThreadArchiveDuration,
	)

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

func (s *SettingsTable) SetContextMenuPermissionLevel(guildId uint64, permissionLevel int) (err error) {
	query := `
INSERT INTO settings("guild_id", "context_menu_permission_level")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "context_menu_permission_level" = $2;
`

	_, err = s.Exec(context.Background(), query, guildId, permissionLevel)
	return
}
