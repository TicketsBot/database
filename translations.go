package database

import (
	"context"
	"github.com/TicketsBot/database/translations"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Translations struct {
	*pgxpool.Pool
}

func newTranslations(db *pgxpool.Pool) *Translations {
	return &Translations{
		db,
	}
}

func (t Translations) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS translations(
	"language" varchar(8) NOT NULL,
	"message_id" int4 NOT NULL,
	"content" text,
	PRIMARY KEY("language", "message_id")
);`
}

func (t *Translations) Get(language database.Language, id database.MessageId) (content string, e error) {
	if err := t.QueryRow(context.Background(), `SELECT "content" from translations WHERE "language" = $1 AND "message_id" = $2`, language, id).Scan(&content); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (t *Translations) GetAll() (messages map[database.Language]map[database.MessageId]string, e error) {
	rows, err := t.Query(context.Background(), "SELECT * FROM translations;")
	if err != nil {
		e = err
		return
	}

	messages = make(map[database.Language]map[database.MessageId]string)

	for rows.Next() {
		var language database.Language
		var messageId database.MessageId
		var value string

		if err := rows.Scan(&language, &messageId, &value); err != nil {
			e = err
			continue
		}

		if messages[language] == nil {
			messages[language] = make(map[database.MessageId]string)
		}

		messages[language][messageId] = value
	}

	return
}
