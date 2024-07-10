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
	"panel_id" int NOT NULL,
	"should_mention_user" bool NOT NULL,
	FOREIGN KEY("panel_id") REFERENCES panels("panel_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("panel_id")
);
`
}

func (p *PanelUserMention) ShouldMentionUser(ctx context.Context, panelId int) (shouldMention bool, e error) {
	query := `SELECT "should_mention_user" from panel_user_mentions WHERE "panel_id"=$1;`

	if err := p.QueryRow(ctx, query, panelId).Scan(&shouldMention); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelUserMention) Set(ctx context.Context, panelId int, shouldMentionUser bool) error {
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := p.SetWithTx(ctx, tx, panelId, shouldMentionUser); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PanelUserMention) SetWithTx(ctx context.Context, tx pgx.Tx, panelId int, shouldMentionUser bool) (err error) {
	query := `INSERT INTO panel_user_mentions("panel_id", "should_mention_user") VALUES($1, $2) ON CONFLICT("panel_id") DO UPDATE SET "should_mention_user" = $2;`
	_, err = tx.Exec(ctx, query, panelId, shouldMentionUser)
	return
}
