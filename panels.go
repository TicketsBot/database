package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Panel struct {
	MessageId      uint64  `json:"message_id,string"`
	ChannelId      uint64  `json:"channel_id,string"`
	GuildId        uint64  `json:"guild_id,string"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Colour         int32   `json:"colour"`
	TargetCategory uint64  `json:"category,string"`
	ReactionEmote  string  `json:"reaction_emote"`
	WelcomeMessage *string `json:"welcome_message"`
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
	return `
CREATE TABLE IF NOT EXISTS panels(
	"message_id" int8 NOT NULL UNIQUE,
	"channel_id" int8 NOT NULL,
	"guild_id" int8 NOT NULL,
	"title" varchar(255) NOT NULL,
	"content" text NOT NULL,
	"colour" int4 NOT NULL,
	"target_category" int8 NOT NULL,
	"reaction_emote" varchar(32) NOT NULL,
	"welcome_message" text,
	PRIMARY KEY("message_id")
);
CREATE INDEX IF NOT EXISTS panels_guild_id ON panels("guild_id");`
}

func (p *PanelTable) Get(messageId uint64) (panel Panel, e error) {
	query := `SELECT * from panels WHERE "message_id" = $1;`

	if err := p.QueryRow(context.Background(), query, messageId).Scan(&panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage); err != nil && err != pgx.ErrNoRows {
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
		if err := rows.Scan(&panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage); err != nil {
			e = err
			continue
		}

		panels = append(panels, panel)
	}

	return
}

func (p *PanelTable) Create(panel Panel) (err error) {
	query := `INSERT INTO panels("message_id", "channel_id", "guild_id", "title", "content", "colour", "target_category", "reaction_emote", "welcome_message") VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT("message_id") DO NOTHING;`
	_, err = p.Exec(context.Background(), query, panel.MessageId, panel.ChannelId, panel.GuildId, panel.Title, panel.Content, panel.Colour, panel.TargetCategory, panel.ReactionEmote, panel.WelcomeMessage)
	return
}

func (p *PanelTable) Update(oldMessageId uint64, panel Panel) (err error) {
	query := `
UPDATE panels
	SET "message_id" = $2,
		"channel_id" = $3,
		"title" = $4,
		"content" = $5,
		"colour" = $6,
		"target_category" = $7,
		"reaction_emote" = $8,
		"welcome_message" = $9
	WHERE
		"message_id" = $1
;`
	_, err = p.Exec(context.Background(), query, oldMessageId, panel.MessageId, panel.ChannelId, panel.Title, panel.Content, panel.Colour, panel.TargetCategory, panel.ReactionEmote, panel.WelcomeMessage)
	return
}

func (p *PanelTable) Delete(messageId uint64) (err error) {
	query := `DELETE FROM panels WHERE "message_id"=$1;`
	_, err = p.Exec(context.Background(), query, messageId)
	return
}
