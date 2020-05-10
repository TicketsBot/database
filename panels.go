package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Panel struct {
	MessageId      uint64
	ChannelId      uint64
	GuildId        uint64
	Title          string
	Content        string
	Colour         int32
	TargetCategory uint64
	ReactionEmote  string
}

type PanelTable struct {
	*pgxpool.Pool
}

func newPanelTable(db *pgxpool.Pool) *PanelTable {
	return &PanelTable{
		db,
	}
}

func (p PanelTable) Schema() string {
	return `CREATE TABLE IF NOT EXISTS panels(
"message_id" int8 NOT NULL UNIQUE,
"channel_id" int8 NOT NULL,
"guild_id" int8 NOT NULL,
"title" varchar(255) NOT NULL,
"content" text NOT NULL,
"colour" int4 NOT NULL,
"target_category" int8 NOT NULL,
"reaction_emote" varchar(32) NOT NULL,
PRIMARY KEY("message_id"));
CREATE INDEX IF NOT EXISTS panels_guild_id ON panels("guild_id");`
}

func (p *PanelTable) Get(messageId uint64) (panel Panel, e error) {
	query := `SELECT * from panels WHERE "message_id" = $1;`

	if err := p.QueryRow(context.Background(), query, messageId).Scan(&panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelTable) GetByGuild(guildId uint64) (panels []Panel, e error) {
	query := `SELECT * from panels WHERE "guild_id" = $1;`

	rows, err := p.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var panel Panel
		if err := rows.Scan(&panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote); err != nil {
			e = err
			continue
		}

		panels = append(panels, panel)
	}

	return
}

func (p *PanelTable) Create(panel Panel) (err error) {
	query := `INSERT INTO panels("message_id", "channel_id", "guild_id", "title", "content", "colour", "target_category", "reaction_emote") VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT("message_id") DO NOTHING;`
	_, err = p.Exec(context.Background(), query, panel.MessageId, panel.ChannelId, panel.GuildId, panel.Title, panel.Content, panel.Colour, panel.TargetCategory, panel.ReactionEmote)
	return
}

func (p *PanelTable) Delete(messageId uint64) (err error) {
	query := `DELETE FROM panels WHERE "message_id"=$1;`
	_, err = p.Exec(context.Background(), query, messageId)
	return
}
