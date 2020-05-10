package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Webhook struct {
	Id    uint64
	Token string
}

type WebhookTable struct {
	*pgxpool.Pool
}

func newWebhookTable(db *pgxpool.Pool) *WebhookTable {
	return &WebhookTable{
		db,
	}
}

func (w WebhookTable) Schema() string {
	return `CREATE TABLE IF NOT EXISTS webhooks("guild_id" int8 NOT NULL, "ticket_id" int4 NOT NULL, "webhook_id" int8 NOT NULL UNIQUE, "webhook_token" varchar(100) NOT NULL, PRIMARY KEY("guild_id", "ticket_id"));`
}

func (w *WebhookTable) Get(guildId uint64, ticketId int) (webhook Webhook, e error) {
	query := `SELECT "webhook_id", "webhook_token" from webhooks WHERE "guild_id"=$1 AND "ticket_id"=$2;`
	if err := w.QueryRow(context.Background(), query, guildId, ticketId).Scan(&webhook.Id, &webhook.Token); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (w *WebhookTable) Create(guildId uint64, ticketId int, webhook Webhook) (err error) {
	query := `INSERT INTO webhooks("guild_id", "ticket_id", "webhook_id", "webhook_token") VALUES($1, $2, $3, $4) ON CONFLICT("guild_id", "ticket_id") DO UPDATE SET "webhook_id" = $3, "webhook_token" = $4;`
	_, err = w.Exec(context.Background(), query, guildId, ticketId, webhook.Id, webhook.Token)
	return
}

func (w *WebhookTable) Delete(guildId uint64, ticketId int) (err error) {
	query := `DELETE FROM webhooks WHERE "guild_id"=$1 AND "ticket_id"=$2;`
	_, err = w.Exec(context.Background(), query, guildId, ticketId)
	return
}
