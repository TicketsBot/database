package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelBot struct {
	UserId    uint64
	BotId     uint64
	PublicKey string
	Token     string
}

type WhitelabelBotTable struct {
	*pgxpool.Pool
}

func newWhitelabelBotTable(db *pgxpool.Pool) *WhitelabelBotTable {
	return &WhitelabelBotTable{
		db,
	}
}

func (w WhitelabelBotTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel(
	"user_id" int8 UNIQUE NOT NULL,
	"bot_id" int8 UNIQUE NOT NULL,
	"public_key" CHAR(64) NOT NULL,
	"token" VARCHAR(84) NOT NULL UNIQUE,
	PRIMARY KEY("user_id")
);
CREATE INDEX IF NOT EXISTS whitelabel_bot_id ON whitelabel("bot_id");
`
}

func (w *WhitelabelBotTable) GetByUserId(ctx context.Context, userId uint64) (WhitelabelBot, error) {
	query := `SELECT "user_id", "bot_id", "public_key", "token" FROM whitelabel WHERE "user_id" = $1;`

	var bot WhitelabelBot
	if err := w.QueryRow(ctx, query, userId).Scan(&bot.UserId, &bot.BotId, &bot.PublicKey, &bot.Token); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return WhitelabelBot{}, err
	}

	return bot, nil
}

func (w *WhitelabelBotTable) GetByBotId(ctx context.Context, botId uint64) (WhitelabelBot, error) {
	query := `SELECT "user_id", "bot_id", "public_key", "token" FROM whitelabel WHERE "bot_id" = $1;`

	var bot WhitelabelBot
	if err := w.QueryRow(ctx, query, botId).Scan(&bot.UserId, &bot.BotId, &bot.PublicKey, &bot.Token); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return WhitelabelBot{}, err
	}

	return bot, nil
}

func (w *WhitelabelBotTable) Set(ctx context.Context, data WhitelabelBot) error {
	query := `
INSERT INTO whitelabel("user_id", "bot_id", "public_key", "token")
VALUES($1, $2, $3, $4)
ON CONFLICT("user_id") DO UPDATE SET "bot_id" = $2, "public_key" = $3, "token" = $4;`
	_, err := w.Exec(ctx, query, data.UserId, data.BotId, data.PublicKey, data.Token)
	return err
}

func (w *WhitelabelBotTable) Delete(ctx context.Context, userId uint64) error {
	query := `DELETE FROM whitelabel WHERE "user_id"=$1;`
	_, err := w.Exec(ctx, query, userId)
	return err
}

func (w *WhitelabelBotTable) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM whitelabel WHERE "token"=$1;`
	_, err := w.Exec(ctx, query, token)
	return err
}
