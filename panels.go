package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ALTER TABLE panels ADD COLUMN default_team bool NOT NULL DEFAULT 't';
type Panel struct {
	PanelId             int     `json:"panel_id"`
	MessageId           uint64  `json:"message_id,string"`
	ChannelId           uint64  `json:"channel_id,string"`
	GuildId             uint64  `json:"guild_id,string"`
	Title               string  `json:"title"`
	Content             string  `json:"content"`
	Colour              int32   `json:"colour"`
	TargetCategory      uint64  `json:"category_id,string"`
	EmojiName           *string `json:"emoji_name"`
	EmojiId             *uint64 `json:"emoji_id,string"`
	WelcomeMessageEmbed *int    `json:"welcome_message_embed"`
	WithDefaultTeam     bool    `json:"default_team"`
	CustomId            string  `json:"-"`
	ImageUrl            *string `json:"image_url,omitempty"`
	ThumbnailUrl        *string `json:"thumbnail_url,omitempty"`
	ButtonStyle         int     `json:"button_style"`
	ButtonLabel         string  `json:"button_label"`
	FormId              *int    `json:"form_id"`
	NamingScheme        *string `json:"naming_scheme"`
}

type PanelWithWelcomeMessage struct {
	Panel
	WelcomeMessage *CustomEmbed
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
	"emoji_name" varchar(32) DEFAULT NULL,
	"emoji_id" int8 DEFAULT NULL,
	"welcome_message" int NULL,
	"default_team" bool NOT NULL,
	"custom_id" varchar(100) NOT NULL,
	"image_url" varchar(255),
	"thumbnail_url" varchar(255),
	"button_style" int2 DEFAULT 1,
	"button_label" varchar(80) NOT NULL,
	"form_id" int DEFAULT NULL,
	"naming_scheme" varchar(100) DEFAULT NULL,
	FOREIGN KEY ("welcome_message") REFERENCES embeds("id") ON DELETE SET NULL,
	FOREIGN KEY ("form_id") REFERENCES forms("form_id"),
	PRIMARY KEY("panel_id")
);
CREATE INDEX IF NOT EXISTS panels_guild_id ON panels("guild_id");
CREATE INDEX IF NOT EXISTS panels_message_id ON panels("message_id");
CREATE INDEX IF NOT EXISTS panels_form_id ON panels("form_id");
CREATE INDEX IF NOT EXISTS panels_guild_id_form_id ON panels("guild_id", "form_id");
CREATE INDEX IF NOT EXISTS panels_custom_id ON panels("custom_id");`
}

func (p *PanelTable) Get(messageId uint64) (panel Panel, e error) {
	query := `
SELECT
	panel_id,
	message_id,
	channel_id,
	guild_id,
	title,
	content,
	colour,
	target_category,
	emoji_name,
	emoji_id,
	welcome_message,
	default_team,
	custom_id,
	image_url,
	thumbnail_url,
	button_style,
	button_label,
	form_id,
	naming_scheme
FROM panels
WHERE "message_id" = $1;
`

	if err := p.QueryRow(context.Background(), query, messageId).Scan(
		&panel.PanelId,
		&panel.MessageId,
		&panel.ChannelId,
		&panel.GuildId,
		&panel.Title,
		&panel.Content,
		&panel.Colour,
		&panel.TargetCategory,
		&panel.EmojiName,
		&panel.EmojiId,
		&panel.WelcomeMessageEmbed,
		&panel.WithDefaultTeam,
		&panel.CustomId,
		&panel.ImageUrl,
		&panel.ThumbnailUrl,
		&panel.ButtonStyle,
		&panel.ButtonLabel,
		&panel.FormId,
		&panel.NamingScheme,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelTable) GetById(panelId int) (panel Panel, e error) {
	query := `
SELECT
	panel_id,
	message_id,
	channel_id,
	guild_id,
	title,
	content,
	colour,
	target_category,
	emoji_name,
	emoji_id,
	welcome_message,
	default_team,
	custom_id,
	image_url,
	thumbnail_url,
	button_style,
	button_label,
	form_id,
	naming_scheme
FROM panels
WHERE "panel_id" = $1;
`

	if err := p.QueryRow(context.Background(), query, panelId).Scan(
		&panel.PanelId,
		&panel.MessageId,
		&panel.ChannelId,
		&panel.GuildId,
		&panel.Title,
		&panel.Content,
		&panel.Colour,
		&panel.TargetCategory,
		&panel.EmojiName,
		&panel.EmojiId,
		&panel.WelcomeMessageEmbed,
		&panel.WithDefaultTeam,
		&panel.CustomId,
		&panel.ImageUrl,
		&panel.ThumbnailUrl,
		&panel.ButtonStyle,
		&panel.ButtonLabel,
		&panel.FormId,
		&panel.NamingScheme,
	); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PanelTable) GetByCustomId(guildId uint64, customId string) (panel Panel, ok bool, e error) {
	query := `
SELECT
	panel_id,
	message_id,
	channel_id,
	guild_id,
	title,
	content,
	colour,
	target_category,
	emoji_name,
	emoji_id,
	welcome_message,
	default_team,
	custom_id,
	image_url,
	thumbnail_url,
	button_style,
	button_label,
	form_id,
	naming_scheme
FROM panels
WHERE "guild_id" = $1 AND "custom_id" = $2;
`

	err := p.QueryRow(context.Background(), query, guildId, customId).Scan(
		&panel.PanelId,
		&panel.MessageId,
		&panel.ChannelId,
		&panel.GuildId,
		&panel.Title,
		&panel.Content,
		&panel.Colour,
		&panel.TargetCategory,
		&panel.EmojiName,
		&panel.EmojiId,
		&panel.WelcomeMessageEmbed,
		&panel.WithDefaultTeam,
		&panel.CustomId,
		&panel.ImageUrl,
		&panel.ThumbnailUrl,
		&panel.ButtonStyle,
		&panel.ButtonLabel,
		&panel.FormId,
		&panel.NamingScheme,
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

func (p *PanelTable) GetByFormId(guildId uint64, formId int) (panel Panel, ok bool, e error) {
	query := `
SELECT
	panel_id,
	message_id,
	channel_id,
	guild_id,
	title,
	content,
	colour,
	target_category,
	emoji_name,
	emoji_id,,
	welcome_message,
	default_team,
	custom_id,
	image_url,
	thumbnail_url,
	button_style,
	button_label,
	form_id,
	naming_scheme
FROM panels
WHERE "guild_id" = $1 AND "form_id" = $2;
`

	err := p.QueryRow(context.Background(), query, guildId, formId).Scan(
		&panel.PanelId,
		&panel.MessageId,
		&panel.ChannelId,
		&panel.GuildId,
		&panel.Title,
		&panel.Content,
		&panel.Colour,
		&panel.TargetCategory,
		&panel.EmojiName,
		&panel.EmojiId,
		&panel.WelcomeMessageEmbed,
		&panel.WithDefaultTeam,
		&panel.CustomId,
		&panel.ImageUrl,
		&panel.ThumbnailUrl,
		&panel.ButtonStyle,
		&panel.ButtonLabel,
		&panel.FormId,
		&panel.NamingScheme,
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

func (p *PanelTable) GetByFormCustomId(guildId uint64, customId string) (panel Panel, ok bool, e error) {
	query := `
SELECT
	panels.panel_id,
	panels.message_id,
	panels.channel_id,
	panels.guild_id,
	panels.title,
	panels.content,
	panels.colour,
	panels.target_category,
	panels.emoji_name,
	panels.emoji_id,
	panels.welcome_message,
	panels.default_team,
	panels.custom_id,
	panels.image_url,
	panels.thumbnail_url,
	panels.button_style,
	panels.button_label,
	panels.form_id,
	panels.naming_scheme
FROM panels
INNER JOIN forms
ON forms.form_id = panels.form_id
WHERE forms.guild_id = $1 AND forms.form_id = $2;
`

	err := p.QueryRow(context.Background(), query, guildId, customId).Scan(
		&panel.PanelId,
		&panel.MessageId,
		&panel.ChannelId,
		&panel.GuildId,
		&panel.Title,
		&panel.Content,
		&panel.Colour,
		&panel.TargetCategory,
		&panel.EmojiName,
		&panel.EmojiId,
		&panel.WelcomeMessageEmbed,
		&panel.WithDefaultTeam,
		&panel.CustomId,
		&panel.ImageUrl,
		&panel.ThumbnailUrl,
		&panel.ButtonStyle,
		&panel.ButtonLabel,
		&panel.FormId,
		&panel.NamingScheme,
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
SELECT
	panel_id,
	message_id,
	channel_id, 
	guild_id,
	title,
	content,
	colour,
	target_category,
	emoji_name,
	emoji_id,
	welcome_message,
	default_team,
	custom_id,
	image_url,
	thumbnail_url,
	button_style,
	button_label,
	form_id,
	naming_scheme
FROM panels
WHERE "guild_id" = $1
ORDER BY "panel_id" ASC;`

	rows, err := p.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var panel Panel
		err := rows.Scan(
			&panel.PanelId,
			&panel.MessageId,
			&panel.ChannelId,
			&panel.GuildId,
			&panel.Title,
			&panel.Content,
			&panel.Colour,
			&panel.TargetCategory,
			&panel.EmojiName,
			&panel.EmojiId,
			&panel.WelcomeMessageEmbed,
			&panel.WithDefaultTeam,
			&panel.CustomId,
			&panel.ImageUrl,
			&panel.ThumbnailUrl,
			&panel.ButtonStyle,
			&panel.ButtonLabel,
			&panel.FormId,
			&panel.NamingScheme,
		)

		if err != nil {
			return nil, err
		}

		panels = append(panels, panel)
	}

	return
}

func (p *PanelTable) GetByGuildWithWelcomeMessage(guildId uint64) (panels []PanelWithWelcomeMessage, e error) {
	query := `
SELECT
	panels.panel_id,
	panels.message_id,
	panels.channel_id, 
	panels.guild_id,
	panels.title,
	panels.content,
	panels.colour,
	panels.target_category,
	panels.emoji_name,
	panels.emoji_id,
	panels.welcome_message,
	panels.default_team,
	panels.custom_id,
	panels.image_url,
	panels.thumbnail_url,
	panels.button_style,
	panels.button_label,
	panels.form_id,
	panels.naming_scheme,
	embeds.id,
	embeds.guild_id,
	embeds.title,
	embeds.description,
	embeds.url,
	embeds.colour,
	embeds.author_name,
	embeds.author_icon_url,
	embeds.author_url,
	embeds.image_url,
	embeds.thumbnail_url,
	embeds.footer_text,
	embeds.footer_icon_url,
	embeds.timestamp
FROM panels
LEFT JOIN embeds
ON panels.welcome_message = embeds.id
WHERE panels.guild_id = $1
ORDER BY panels.panel_id ASC;`

	rows, err := p.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var panel Panel
		var embed CustomEmbed

		// Can't scan missing values into non-nullable fields
		var embedId *int
		var embedGuildId *uint64
		var embedColour *uint32

		err := rows.Scan(
			&panel.PanelId,
			&panel.MessageId,
			&panel.ChannelId,
			&panel.GuildId,
			&panel.Title,
			&panel.Content,
			&panel.Colour,
			&panel.TargetCategory,
			&panel.EmojiName,
			&panel.EmojiId,
			&panel.WelcomeMessageEmbed,
			&panel.WithDefaultTeam,
			&panel.CustomId,
			&panel.ImageUrl,
			&panel.ThumbnailUrl,
			&panel.ButtonStyle,
			&panel.ButtonLabel,
			&panel.FormId,
			&panel.NamingScheme,
			&embedId,
			&embedGuildId,
			&embed.Title,
			&embed.Description,
			&embed.Url,
			&embedColour,
			&embed.AuthorName,
			&embed.AuthorIconUrl,
			&embed.AuthorUrl,
			&embed.ImageUrl,
			&embed.ThumbnailUrl,
			&embed.FooterText,
			&embed.FooterIconUrl,
			&embed.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		var embedPtr *CustomEmbed
		if embedId != nil {
			embed.Id = *embedId
			embed.GuildId = *embedGuildId
			embed.Colour = *embedColour

			embedPtr = &embed
		}

		panels = append(panels, PanelWithWelcomeMessage{
			Panel:          panel,
			WelcomeMessage: embedPtr,
		})
	}

	return
}

func (p *PanelTable) Create(panel Panel) (panelId int, err error) {
	query := `
INSERT INTO panels(
	"message_id",
	"channel_id",
	"guild_id",
	"title",
	"content",
	"colour",
	"target_category",
	"emoji_name",
	"emoji_id",
	"welcome_message",
	"default_team",
	"custom_id",
	"image_url",
	"thumbnail_url",
	"button_style",
	"button_label",
	"form_id",
	"naming_scheme"
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
ON CONFLICT("message_id") DO NOTHING
RETURNING "panel_id";`

	err = p.QueryRow(context.Background(), query,
		panel.MessageId,
		panel.ChannelId,
		panel.GuildId,
		panel.Title,
		panel.Content,
		panel.Colour,
		panel.TargetCategory,
		panel.EmojiName,
		panel.EmojiId,
		panel.WelcomeMessageEmbed,
		panel.WithDefaultTeam,
		panel.CustomId,
		panel.ImageUrl,
		panel.ThumbnailUrl,
		panel.ButtonStyle,
		panel.ButtonLabel,
		panel.FormId,
		panel.NamingScheme,
	).Scan(&panelId)

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
		"emoji_name" = $8,
		"emoji_id" = $9,
		"welcome_message" = $10,
		"default_team" = $11,
		"custom_id" = $12,
		"image_url" = $13,
		"thumbnail_url" = $14,
		"button_style" = $15,
		"button_label" = $16,
		"form_id" = $17,
		"naming_scheme" = $18
	WHERE
		"panel_id" = $1
;`
	_, err = p.Exec(context.Background(), query,
		panel.PanelId,
		panel.MessageId,
		panel.ChannelId,
		panel.Title,
		panel.Content,
		panel.Colour,
		panel.TargetCategory,
		panel.EmojiName,
		panel.EmojiId,
		panel.WelcomeMessageEmbed,
		panel.WithDefaultTeam,
		panel.CustomId,
		panel.ImageUrl,
		panel.ThumbnailUrl,
		panel.ButtonStyle,
		panel.ButtonLabel,
		panel.FormId,
		panel.NamingScheme,
	)
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
