package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type AutoCloseTable struct {
	*pgxpool.Pool
}

type AutoCloseSettings struct {
	Enabled                 bool           `json:"enabled"`
	SinceOpenWithNoResponse *time.Duration `json:"since_open_with_no_response"`
	SinceLastMessage        *time.Duration `json:"since_last_message"`
	OnUserLeave             *bool          `json:"on_user_leave"`
}

func newAutoCloseTable(db *pgxpool.Pool) *AutoCloseTable {
	return &AutoCloseTable{
		db,
	}
}

func (a AutoCloseTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS auto_close(
	"guild_id" int8 NOT NULL,
	"enabled" bool NOT NULL,
	"since_open_with_no_response" interval,
	"since_last_message" interval,
	"on_user_leave" bool,
	PRIMARY KEY("guild_id")
);
`
}

func (a *AutoCloseTable) Get(ctx context.Context, guildId uint64) (settings AutoCloseSettings, e error) {
	query := `SELECT "enabled", "since_open_with_no_response", "since_last_message", "on_user_leave" FROM auto_close WHERE "guild_id" = $1;`
	if err := a.QueryRow(ctx, query, guildId).Scan(&settings.Enabled, &settings.SinceOpenWithNoResponse, &settings.SinceLastMessage, &settings.OnUserLeave); err != nil && err != pgx.ErrNoRows { // defaults to nil if no rows
		e = err
	}

	return
}

func (a *AutoCloseTable) Set(ctx context.Context, guildId uint64, settings AutoCloseSettings) (err error) {
	query := `
INSERT INTO
	auto_close("guild_id", "enabled", "since_open_with_no_response", "since_last_message", "on_user_leave")
VALUES
	($1, $2, $3, $4, $5)
ON CONFLICT("guild_id") DO
	UPDATE SET
		"enabled" = $2,
		"since_open_with_no_response" = $3,
		"since_last_message" = $4,
		"on_user_leave" = $5
;`

	_, err = a.Exec(ctx, query, guildId, settings.Enabled, settings.SinceOpenWithNoResponse, settings.SinceLastMessage, settings.OnUserLeave)
	return
}

func (a *AutoCloseTable) Reset(ctx context.Context, guildId uint64) (err error) {
	query := `
UPDATE auto_close
SET since_open_with_no_response = NULL, since_last_message = NULL
WHERE "guild_id" = $1;
`

	_, err = a.Exec(ctx, query, guildId)
	return
}

func (a *AutoCloseTable) Delete(ctx context.Context, guildId uint64) (err error) {
	query := `
DELETE FROM auto_close
WHERE "guild_id" = $1;
`

	_, err = a.Exec(ctx, query, guildId)
	return
}
