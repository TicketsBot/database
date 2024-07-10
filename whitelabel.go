package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelBot struct {
	UserId uint64
	BotId  uint64
	Token  string
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
	"token" VARCHAR(84) NOT NULL UNIQUE,
	PRIMARY KEY("user_id")
);
CREATE INDEX IF NOT EXISTS whitelabel_bot_id ON whitelabel("bot_id");
`
}

func (w *WhitelabelBotTable) GetByUserId(ctx context.Context, userId uint64) (res WhitelabelBot, e error) {
	query := `SELECT * FROM whitelabel WHERE "user_id" = $1;`
	if err := w.QueryRow(ctx, query, userId).Scan(&res.UserId, &res.BotId, &res.Token); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (w *WhitelabelBotTable) GetByBotId(ctx context.Context, botId uint64) (res WhitelabelBot, e error) {
	query := `SELECT * FROM whitelabel WHERE "bot_id" = $1;`
	if err := w.QueryRow(ctx, query, botId).Scan(&res.UserId, &res.BotId, &res.Token); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (w *WhitelabelBotTable) Set(ctx context.Context, data WhitelabelBot) (err error) {
	query := `INSERT INTO whitelabel("user_id", "bot_id", "token") VALUES($1, $2, $3) ON CONFLICT("user_id") DO UPDATE SET "bot_id" = $2, "token" = $3;`
	_, err = w.Exec(ctx, query, data.UserId, data.BotId, data.Token)
	return
}

func (w *WhitelabelBotTable) Delete(ctx context.Context, userId uint64) (err error) {
	query := `DELETE FROM whitelabel WHERE "user_id"=$1;`
	_, err = w.Exec(ctx, query, userId)
	return
}

func (w *WhitelabelBotTable) DeleteByToken(ctx context.Context, token string) (err error) {
	query := `DELETE FROM whitelabel WHERE "token"=$1;`
	_, err = w.Exec(ctx, query, token)
	return
}
