package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO: Migrate all settings to this table
type Settings struct {
	HideClaimButton             bool    `json:"hide_claim_button"`
	DisableOpenCommand          bool    `json:"disable_open_command"`
	ContextMenuPermissionLevel  int     `json:"context_menu_permission_level,string"`
	ContextMenuAddSender        bool    `json:"context_menu_add_sender"`
	ContextMenuPanel            *int    `json:"context_menu_panel"`
	StoreTranscripts            bool    `json:"store_transcripts"`
	UseThreads                  bool    `json:"use_threads"`
	TicketNotificationChannel   *uint64 `json:"ticket_notification_channel,string"`
	ThreadArchiveDuration       int     `json:"thread_archive_duration"`
	OverflowEnabled             bool    `json:"overflow_enabled"`
	OverflowCategoryId          *uint64 `json:"overflow_category_id,string"` // If overflow_enabled and nil, use root
	ExitSurveyFormId            *uint64 `json:"exit_survey_form_id,string"`
	AnonymiseDashboardResponses bool    `json:"anonymise_dashboard_responses"`
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
		TicketNotificationChannel:  nil,
		ThreadArchiveDuration:      10080,
		OverflowEnabled:            false,
		OverflowCategoryId:         nil,
		ExitSurveyFormId:           nil,
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
	"ticket_notification_channel" int8 DEFAULT NULL,
    "thread_archive_duration" int DEFAULT '10080',
	"overflow_enabled" bool DEFAULT 'f',
	"overflow_category_id" int8 DEFAULT NULL,
	"exit_survey_form_id" int4 DEFAULT NULL,
	"anonymise_dashboard_responses" bool DEFAULT 'f',
	FOREIGN KEY("context_menu_panel") REFERENCES panels("panel_id") ON DELETE SET NULL,
	FOREIGN KEY("exit_survey_form_id") REFERENCES forms("form_id") ON DELETE SET NULL,
	PRIMARY KEY("guild_id"),
	CHECK (use_threads = false OR ticket_notification_channel IS NOT NULL)
);
`
}

func (s *SettingsTable) Get(ctx context.Context, guildId uint64) (Settings, error) {
	query := `
SELECT
	"hide_claim_button",
	"disable_open_command",
	"context_menu_permission_level",
	"context_menu_add_sender",
	"context_menu_panel",
	"store_transcripts",
    "use_threads",
	"ticket_notification_channel",
    "thread_archive_duration",
	"overflow_enabled",
	"overflow_category_id",
	"anonymise_dashboard_responses"
FROM settings
WHERE "guild_id" = $1;
`

	var settings Settings
	err := s.QueryRow(ctx, query, guildId).Scan(
		&settings.HideClaimButton,
		&settings.DisableOpenCommand,
		&settings.ContextMenuPermissionLevel,
		&settings.ContextMenuAddSender,
		&settings.ContextMenuPanel,
		&settings.StoreTranscripts,
		&settings.UseThreads,
		&settings.TicketNotificationChannel,
		&settings.ThreadArchiveDuration,
		&settings.OverflowEnabled,
		&settings.OverflowCategoryId,
		&settings.AnonymiseDashboardResponses,
	)

	if err == nil {
		return settings, nil
	} else if err == pgx.ErrNoRows {
		return defaultSettings(), nil
	} else {
		return settings, err
	}
}

func (s *SettingsTable) Set(ctx context.Context, guildId uint64, settings Settings) (err error) {
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
	"ticket_notification_channel",
    "thread_archive_duration",
	"overflow_enabled",
	"overflow_category_id",
	"anonymise_dashboard_responses"
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT("guild_id")
DO UPDATE SET
	"hide_claim_button" = $2,
	"disable_open_command" = $3,
	"context_menu_permission_level" = $4,
	"context_menu_add_sender" = $5,
	"context_menu_panel" = $6,
	"store_transcripts" = $7,
    "use_threads" = $8,
    "ticket_notification_channel" = $9,
    "thread_archive_duration" = $10,
	"overflow_enabled" = $11,
	"overflow_category_id" = $12,
	"anonymise_dashboard_responses" = $13
;
`

	_, err = s.Exec(ctx, query,
		guildId,
		settings.HideClaimButton,
		settings.DisableOpenCommand,
		settings.ContextMenuPermissionLevel,
		settings.ContextMenuAddSender,
		settings.ContextMenuPanel,
		settings.StoreTranscripts,
		settings.UseThreads,
		settings.TicketNotificationChannel,
		settings.ThreadArchiveDuration,
		settings.OverflowEnabled,
		settings.OverflowCategoryId,
		settings.AnonymiseDashboardResponses,
	)

	return
}

func (s *SettingsTable) SetHideClaimButton(ctx context.Context, guildId uint64, hideClaimButton bool) (err error) {
	query := `
INSERT INTO settings("guild_id", "hide_claim_button")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "hide_claim_button" = $2;
`

	_, err = s.Exec(ctx, query, guildId, hideClaimButton)
	return
}

func (s *SettingsTable) SetDisableOpenCommand(ctx context.Context, guildId uint64, disableOpenCommand bool) (err error) {
	query := `
INSERT INTO settings("guild_id", "disable_open_command")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "disable_open_command" = $2;
`

	_, err = s.Exec(ctx, query, guildId, disableOpenCommand)
	return
}

func (s *SettingsTable) SetContextMenuPermissionLevel(ctx context.Context, guildId uint64, permissionLevel int) (err error) {
	query := `
INSERT INTO settings("guild_id", "context_menu_permission_level")
VALUES($1, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "context_menu_permission_level" = $2;
`

	_, err = s.Exec(ctx, query, guildId, permissionLevel)
	return
}

func (s *SettingsTable) SetOverflow(ctx context.Context, guildId uint64, enabled bool, categoryId *uint64) (err error) {
	query := `
INSERT INTO settings("guild_id", "overflow_enabled", "overflow_category_id")
VALUES($1, $2, $3)
ON CONFLICT("guild_id")
DO UPDATE SET "overflow_enabled" = $2, "overflow_category_id" = $3;
`

	_, err = s.Exec(ctx, query, guildId, enabled, categoryId)
	return
}

func (s *SettingsTable) EnableThreads(ctx context.Context, guildId uint64, ticketNotificationChannel uint64) (err error) {
	query := `
INSERT INTO settings("guild_id", "use_threads", "ticket_notification_channel")
VALUES($1, true, $2)
ON CONFLICT("guild_id")
DO UPDATE SET "use_threads" = true, "ticket_notification_channel" = $2;
`

	_, err = s.Exec(ctx, query, guildId, ticketNotificationChannel)
	return
}

func (s *SettingsTable) DisableThreads(ctx context.Context, guildId uint64) (err error) {
	query := `
INSERT INTO settings("guild_id", "use_threads", "ticket_notification_channel")
VALUES($1, false, NULL)
ON CONFLICT("guild_id")
DO UPDATE SET "use_threads" = false, "ticket_notification_channel" = NULL;
`

	_, err = s.Exec(ctx, query, guildId)
	return
}
