package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PanelUserMention struct {
	*pgxpool.Pool
}

func newPanelUserMention(db *pgxpool.Pool) *PanelUserMention {
	return &PanelUserMention{
		db,
	}
}

func (p PanelUserMention) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS panel_user_mentions(
	"panel_message_id" int8 NOT NULL,
	"should_mention_user" bool NOT NULL,
	FOREIGN KEY("panel_message_id") REFERENCES panels("message_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("panel_message_id")
);
`
}

func (p *PanelUserMention) ShouldMentionUser(panelMessageId uint64) (shouldMention bool, e error) {
	query := `SELECT "should_mention_user" from panel_user_mentions WHERE "panel_message_id"=$1;`

	if err := p.QueryRow(context.Background(), query, panelMessageId).Scan(&shouldMention); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelUserMention) Set(panelMessageId uint64, shouldMentionUser bool) (err error) {
	query := `INSERT INTO panel_user_mentions("panel_message_id", "should_mention_user") VALUES($1, $2) ON CONFLICT("panel_message_id") DO UPDATE SET "should_mention_user" = $2;`
	_, err = p.Exec(context.Background(), query, panelMessageId, shouldMentionUser)
	return
}
