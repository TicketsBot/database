package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MultiPanel struct {
	Id             int     `json:"id"`
	MessageId      uint64  `json:"message_id,string"`
	ChannelId      uint64  `json:"channel_id,string"`
	GuildId        uint64  `json:"guild_id,string"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Colour         int     `json:"colour"`
	SelectMenu     bool    `json:"select_menu"`
	SelectMenuText *string `json:"select_menu_text"`
	ImageUrl       *string `json:"image_url,omitempty"`
	ThumbnailUrl   *string `json:"thumbnail_url,omitempty"`
}

type MultiPanelTable struct {
	*pgxpool.Pool
}

func newMultiMultiPanelTable(db *pgxpool.Pool) *MultiPanelTable {
	return &MultiPanelTable{
		db,
	}
}

func (MultiPanelTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS multi_panels(
	"id" SERIAL NOT NULL,
	"message_id" int8 NOT NULL,
	"channel_id" int8 NOT NULL,
	"guild_id" int8 NOT NULL,
	"title" varchar(255) NOT NULL,
	"content" text NOT NULL,
	"colour" int4 NOT NULL,
	"select_menu" bool DEFAULT 'f',
	"image_url" varchar(255) DEFAULT NULL,
	"thumbnail_url" varchar(255) DEFAULT NULL,
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS multi_panels_guild_id ON multi_panels("guild_id");
CREATE INDEX IF NOT EXISTS multi_panels_message_id ON multi_panels("message_id");`
}

func (p *MultiPanelTable) Get(ctx context.Context, id int) (MultiPanel, bool, error) {
	query := `
SELECT
	"id", "message_id", "channel_id", "guild_id", "title", "content", "colour", "select_menu", "image_url", "thumbnail_url"
FROM
	multi_panels
WHERE
	"id" = $1
;`

	var panel MultiPanel
	err := p.QueryRow(ctx, query, id).Scan(
		&panel.Id, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.SelectMenu, &panel.ImageUrl, &panel.ThumbnailUrl,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return MultiPanel{}, false, nil
		} else {
			return MultiPanel{}, false, err
		}
	}

	return panel, true, nil
}

func (p *MultiPanelTable) GetByMessageId(ctx context.Context, messageId uint64) (MultiPanel, bool, error) {
	query := `
SELECT
	"id", "message_id", "channel_id", "guild_id", "title", "content", "colour", "select_menu", "image_url", "thumbnail_url"
FROM
	multi_panels
WHERE
	"message_id" = $1
;`

	var panel MultiPanel
	err := p.QueryRow(ctx, query, messageId).Scan(
		&panel.Id, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.SelectMenu, &panel.ImageUrl, &panel.ThumbnailUrl,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return MultiPanel{}, false, nil
		} else {
			return MultiPanel{}, false, err
		}
	}

	return panel, true, nil
}

func (p *MultiPanelTable) GetByGuild(ctx context.Context, guildId uint64) ([]MultiPanel, error) {
	query := `
SELECT "id", "message_id", "channel_id", "guild_id", "title", "content", "colour", "select_menu", "image_url", "thumbnail_url"
FROM multi_panels
WHERE "guild_id" = $1;
`

	rows, err := p.Query(ctx, query, guildId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var panels []MultiPanel
	for rows.Next() {
		var panel MultiPanel
		err := rows.Scan(
			&panel.Id, &panel.MessageId, &panel.ChannelId, &panel.GuildId, &panel.Title, &panel.Content, &panel.Colour, &panel.SelectMenu, &panel.ImageUrl, &panel.ThumbnailUrl,
		)

		if err != nil {
			return nil, err
		}

		panels = append(panels, panel)
	}

	return panels, nil
}

func (p *MultiPanelTable) Create(ctx context.Context, panel MultiPanel) (multiPanelId int, err error) {
	query := `
INSERT INTO
	multi_panels("message_id", "channel_id", "guild_id", "title", "content", "colour", "select_menu", "image_url", "thumbnail_url")
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
	"id"
;
`

	err = p.QueryRow(ctx, query,
		panel.MessageId, panel.ChannelId, panel.GuildId, panel.Title, panel.Content, panel.Colour, panel.SelectMenu, panel.ImageUrl, panel.ThumbnailUrl,
	).Scan(&multiPanelId)

	return
}

func (p *MultiPanelTable) Update(ctx context.Context, multiPanelId int, multiPanel MultiPanel) (err error) {
	query := `
UPDATE multi_panels
	SET "message_id" = $2,
		"channel_id" = $3,
		"title" = $4,
		"content" = $5,
		"colour" = $6,
		"select_menu" = $7,
        "image_url" = $8,
        "thumbnail_url" = $9
	WHERE
		"id" = $1
;`
	_, err = p.Exec(ctx, query,
		multiPanelId, multiPanel.MessageId, multiPanel.ChannelId, multiPanel.Title, multiPanel.Content, multiPanel.Colour, multiPanel.SelectMenu, multiPanel.ImageUrl, multiPanel.ThumbnailUrl,
	)

	return
}

func (p *MultiPanelTable) UpdateMessageId(ctx context.Context, multiPanelId int, messageId uint64) (err error) {
	query := `
UPDATE multi_panels
SET "message_id" = $1
WHERE "id" = $2;
`

	_, err = p.Exec(ctx, query, messageId, multiPanelId)
	return
}

func (p *MultiPanelTable) Delete(ctx context.Context, guildId uint64, multiPanelId int) (success bool, err error) {
	query := `
WITH deleted AS (
	DELETE FROM
		multi_panels
	WHERE
		"guild_id" = $1
		AND
		"id" = $2
	RETURNING *
)

SELECT
	COUNT(*)
FROM
	deleted
;
`

	var count int
	err = p.QueryRow(ctx, query, guildId, multiPanelId).Scan(&count)
	success = count > 0

	return
}
