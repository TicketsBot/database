package database

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ModmailWebhook struct {
	Uuid         uuid.UUID
	WebhookId    uint64
	WebhookToken string
}

type ModmailWebhookTable struct {
	*pgxpool.Pool
}

func newModmailWebhookTable(db *pgxpool.Pool) *ModmailWebhookTable {
	return &ModmailWebhookTable{
		db,
	}
}

func (m ModmailWebhookTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS modmail_webhooks(
	"uuid" uuid NOT NULL REFERENCES modmail_sessions("uuid"),
	"webhook_id" int8 NOT NULL UNIQUE,
	"webhook_token" varchar(100) NOT NULL,
	PRIMARY KEY("uuid")
);
`
}

func (m *ModmailWebhookTable) Get(uuid uuid.UUID) (webhook ModmailWebhook, e error) {
	query := `SELECT * from modmail_webhooks WHERE "uuid" = $1;`
	if err := m.QueryRow(context.Background(), query, uuid).Scan(&webhook.Uuid, &webhook.WebhookId, &webhook.WebhookToken); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (m *ModmailWebhookTable) Create(webhook ModmailWebhook) (err error) {
	query := `INSERT INTO modmail_webhooks("uuid", "webhook_id", "webhook_token") VALUES($1, $2, $3) ON CONFLICT("uuid") DO NOTHING;`
	_, err = m.Exec(context.Background(), query, webhook.Uuid, webhook.WebhookId, webhook.WebhookToken)
	return
}

func (m *ModmailWebhookTable) Delete(uuid uuid.UUID) (err error) {
	query := `DELETE FROM modmail_webhooks WHERE "uuid" = $1;`
	_, err = m.Exec(context.Background(), query, uuid)
	return
}
