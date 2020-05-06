package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WelcomeMessages struct {
	*pgxpool.Pool
}

func newWelcomeMessages(db *pgxpool.Pool) *WelcomeMessages {
	return &WelcomeMessages{
		db,
	}
}

func (w WelcomeMessages) Schema() string {
	return `CREATE TABLE IF NOT EXISTS welcome_messages("guild_id" int8 NOT NULL UNIQUE, "welcome_message" text NOT NULL, PRIMARY KEY("guild_id"));`
}

func (w *WelcomeMessages) Get(guildId uint64) (welcomeMessage string, e error) {
	query := `SELECT "welcome_message" from welcome_message WHERE "guild_id" = $1`

	if err := w.QueryRow(context.Background(), query, guildId).Scan(&welcomeMessage); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (w *WelcomeMessages) Set(guildId uint64, welcomeMessage string) (err error) {
	query := `INSERT INTO welcome_messages("guild_id", "welcome_message") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "welcome_message" = $2;`
	_, err = w.Exec(context.Background(), query, guildId, welcomeMessage)
	return
}
