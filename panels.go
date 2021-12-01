package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ALTER TABLE panels ADD COLUMN default_team bool NOT NULL DEFAULT 't';
type Panel struct {
	PanelId         int     `json:"panel_id"`
	MessageId       uint64  `json:"message_id,string"`
	ChannelId       uint64  `json:"channel_id,string"`
	GuildId         uint64  `json:"guild_id,string"`
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	Colour          int32   `json:"colour"`
	TargetCategory  uint64  `json:"category_id,string"`
	ReactionEmote   string  `json:"emote"`
	WelcomeMessage  *string `json:"welcome_message"`
	WithDefaultTeam bool    `json:"default_team"`
	CustomId        string  `json:"-"`
	ImageUrl        *string `json:"image_url,omitempty"`
	ThumbnailUrl    *string `json:"thumbnail_url,omitempty"`
	ButtonStyle     int     `json:"button_style"`
	FormId          *int    `json:"form_id"`
}

type PanelTable struct {
	*pgxpool.Pool
}

func newPanelTable(db *pgxpool.Pool) *PanelTable {
	return &PanelTable{
		db,
	}
}

// TODO: Make custom_id unique
func (p PanelTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS panels(
	"panel_id" SERIAL NOT NULL UNIQUE,
	"message_id" int8 NOT NULL UNIQUE,
	"channel_id" int8 NOT NULL,
	"guild_id" int8 NOT NULL,
	"title" varchar(255) NOT NULL,
	"content" text NOT NULL,
	"colour" int4 NOT NULL,
	"target_category" int8 NOT NULL,
	"reaction_emote" varchar(32) NOT NULL,
	"welcome_message" text,
	"default_team" bool NOT NULL,
	"custom_id" varchar(100) NOT NULL,
	"image_url" varchar(255),
	"thumbnail_url" varchar(255),
	"button_style" int2 DEFAULT 1,
	"form_id" int DEFAULT NULL,
	FOREIGN KEY ("form_id") REFERENCES forms("form_id"),
	PRIMARY KEY("panel_id")
);
CREATE INDEX IF NOT EXISTS panels_guild_id ON panels("guild_id");
CREATE INDEX IF NOT EXISTS panels_message_id ON panels("message_id");
CREATE INDEX IF NOT EXISTS panels_custom_id ON panels("custom_id");`
}

func (p *PanelTable) Get(messageId uint64) (panel Panel, e error) {
	query := `
SELECT panel_id, message_id, channel_id, guild_id, title, content, colour, target_category, reaction_emote, welcome_message, default_team, custom_id, image_url, thumbnail_url, button_style, form_id
FROM panels
WHERE "message_id" = $1;
`

	if err := p.QueryRow(context.Background(), query, messageId).Scan(
		&panel.PanelId, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage, &panel.WithDefaultTeam, &panel.CustomId, &panel.ImageUrl, &panel.ThumbnailUrl, &panel.ButtonStyle, &panel.FormId,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelTable) GetById(panelId int) (panel Panel, e error) {
	query := `
SELECT panel_id, message_id, channel_id, guild_id, title, content, colour, target_category, reaction_emote, welcome_message, default_team, custom_id, image_url, thumbnail_url, button_style, form_id
FROM panels
WHERE "panel_id" = $1;
`

	if err := p.QueryRow(context.Background(), query, panelId).Scan(
		&panel.PanelId, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage, &panel.WithDefaultTeam, &panel.CustomId, &panel.ImageUrl, &panel.ThumbnailUrl, &panel.ButtonStyle, &panel.FormId,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelTable) GetByCustomId(guildId uint64, customId string) (panel Panel, ok bool, e error) {
	query := `
SELECT panel_id, message_id, channel_id, guild_id, title, content, colour, target_category, reaction_emote, welcome_message, default_team, custom_id, image_url, thumbnail_url, button_style, form_id
FROM panels
WHERE "guild_id" = $1 AND "custom_id" = $2;
`

	err := p.QueryRow(context.Background(), query, guildId, customId).Scan(
		&panel.PanelId, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage, &panel.WithDefaultTeam, &panel.CustomId, &panel.ImageUrl, &panel.ThumbnailUrl, &panel.ButtonStyle, &panel.FormId,
	)

	switch err {
	case nil:
		ok = true
	case pgx.ErrNoRows:
	default:
		e = err
	}

	return
}

func (p *PanelTable) GetByGuild(guildId uint64) (panels []Panel, e error) {
	query := `
SELECT panel_id, message_id, channel_id, guild_id, title, content, colour, target_category, reaction_emote, welcome_message, default_team, custom_id, image_url, thumbnail_url, button_style, form_id
FROM panels
WHERE "guild_id" = $1;`

	rows, err := p.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var panel Panel
		if err := rows.Scan(
			&panel.PanelId, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.TargetCategory, &panel.ReactionEmote, &panel.WelcomeMessage, &panel.WithDefaultTeam, &panel.CustomId, &panel.ImageUrl, &panel.ThumbnailUrl, &panel.ButtonStyle, &panel.FormId,
		); err != nil {
			e = err
			continue
		}

		panels = append(panels, panel)
	}

	return
}

func (p *PanelTable) Create(panel Panel) (panelId int, err error) {
	query := `
INSERT INTO panels("message_id", "channel_id", "guild_id", "title", "content", "colour", "target_category", "reaction_emote", "welcome_message", "default_team", "custom_id", "image_url", "thumbnail_url", "button_style", "form_id")
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
ON CONFLICT("message_id") DO NOTHING
RETURNING "panel_id";`

	err = p.QueryRow(context.Background(), query, panel.MessageId, panel.ChannelId, panel.GuildId, panel.Title, panel.Content, panel.Colour, panel.TargetCategory, panel.ReactionEmote, panel.WelcomeMessage, panel.WithDefaultTeam, panel.CustomId, panel.ImageUrl, panel.ThumbnailUrl, panel.ButtonStyle, panel.FormId).Scan(&panelId)
	return
}

func (p *PanelTable) Update(panel Panel) (err error) {
	query := `
UPDATE panels
	SET "message_id" = $2,
		"channel_id" = $3,
		"title" = $4,
		"content" = $5,
		"colour" = $6,
		"target_category" = $7,
		"reaction_emote" = $8,
		"welcome_message" = $9,
		"default_team" = $10,
		"custom_id" = $11,
		"image_url" = $12,
		"thumbnail_url" = $13,
		"button_style" = $14
		"form_id" = $15
	WHERE
		"panel_id" = $1
;`
	_, err = p.Exec(context.Background(), query, panel.PanelId, panel.MessageId, panel.ChannelId, panel.Title, panel.Content, panel.Colour, panel.TargetCategory, panel.ReactionEmote, panel.WelcomeMessage, panel.WithDefaultTeam, panel.CustomId, panel.ImageUrl, panel.ThumbnailUrl, panel.ButtonStyle, panel.FormId)
	return
}

func (p *PanelTable) UpdateMessageId(panelId int, messageId uint64) (err error) {
	query := `
UPDATE panels
SET "message_id" = $1
WHERE "panel_id" = $2;
`

	_, err = p.Exec(context.Background(), query, messageId, panelId)
	return
}

func (p *PanelTable) Delete(panelId int) (err error) {
	query := `DELETE FROM panels WHERE "panel_id"=$1;`
	_, err = p.Exec(context.Background(), query, panelId)
	return
}
