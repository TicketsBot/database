package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MentionType uint16

const (
	MentionTypeRole MentionType = 0
	MentionTypeUser MentionType = 1
)

type PanelMention struct {
	PanelMessageId uint64      `json:"message_id"`
	Snowflake      uint64      `json:"snowflake"`
	Type           MentionType `json:"type"`
}

type PanelMentionTable struct {
	*pgxpool.Pool
}

func newPanelMentionTable(db *pgxpool.Pool) *PanelMentionTable {
	return &PanelMentionTable{
		db,
	}
}

func (p PanelMentionTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS panel_mentions(
	"panel_message_id" int8 NOT NULL,
	"snowflake" int8 NOT NULL,
	"type" int2 NOT NULL,
	FOREIGN KEY("panel_message_id") REFERENCES panels("message_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("panel_message_id", "snowflake")
);
CREATE INDEX IF NOT EXISTS panel_mentions_panel_message_id ON panel_mentions("panel_message_id");
`
}

func (p *PanelMentionTable) GetMentions(panelMessageId uint64) (mentions []PanelMention, e error) {
	query := `SELECT * from panel_mentions WHERE "panel_message_id"=$1;`

	rows, err := p.Query(context.Background(), query, panelMessageId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var mention PanelMention
		var mentionType uint16
		if err := rows.Scan(&mention.PanelMessageId, &mention.Snowflake, &mentionType); err != nil {
			continue
		}

		mention.Type = MentionType(mentionType)

		mentions = append(mentions, mention)
	}

	return
}

func (p *PanelMentionTable) Add(panelMessageId, snowflake uint64, mentionType MentionType) (err error) {
	query := `INSERT INTO panel_mentions("panel_message_id", "snowflake", "type") VALUES($1, $2, $3) ON CONFLICT("panel_message_id", "snowflake") DO NOTHING;`
	_, err = p.Exec(context.Background(), query, panelMessageId, snowflake, mentionType)
	return
}

func (p *PanelMentionTable) Delete(panelMessageId, snowflake uint64) (err error) {
	query := `DELETE FROM panel_mentions WHERE "panel_message_id"=$1 AND "snowflake"=$2;`
	_, err = p.Exec(context.Background(), query, panelMessageId, snowflake)
	return
}
