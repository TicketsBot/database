package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type CustomEmbed struct {
	Id            int        `json:"-"`
	GuildId       uint64     `json:"-"`
	Title         *string    `json:"title,omitempty"`
	Description   *string    `json:"description,omitempty"`
	Url           *string    `json:"url,omitempty"`
	Colour        uint32     `json:"colour,omitempty"`
	AuthorName    *string    `json:"author_name,omitempty"`
	AuthorIconUrl *string    `json:"author_icon_url,omitempty"`
	AuthorUrl     *string    `json:"author_url,omitempty"`
	ImageUrl      *string    `json:"image_url,omitempty"`
	ThumbnailUrl  *string    `json:"thumbnail_url,omitempty"`
	FooterText    *string    `json:"footer_text,omitempty"`
	FooterIconUrl *string    `json:"footer_icon_url,omitempty"`
	Timestamp     *time.Time `json:"timestamp,omitempty"`
}

type CustomEmbedWithFields struct {
	*CustomEmbed
	Fields []EmbedField `json:"fields,omitempty"`
}

type EmbedsTable struct {
	*pgxpool.Pool
}

func newEmbedsTable(db *pgxpool.Pool) *EmbedsTable {
	return &EmbedsTable{
		db,
	}
}

func (s EmbedsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS embeds(
	"id" SERIAL NOT NULL UNIQUE,
	"guild_id" int8 NOT NULL,
	"title" VARCHAR(255) NULL,
	"description" TEXT NULL CONSTRAINT description_length CHECK (length(description) <= 4096),
	"url" VARCHAR(255) NULL,
	"colour" int4 NOT NULL CONSTRAINT colour_range CHECK (colour >= 0 AND colour <= 16777215),
	"author_name" VARCHAR(255) NULL,
	"author_icon_url" VARCHAR(255) NULL,
	"author_url" VARCHAR(255) NULL,
	"image_url" VARCHAR(255) NULL,
	"thumbnail_url" VARCHAR(255) NULL,
	"footer_text" TEXT NULL CONSTRAINT footer_text_length CHECK (length(footer_text) <= 2048),
	"footer_icon_url" VARCHAR(255) NULL,
	"timestamp" TIMESTAMP NULL,
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS embeds_guild_id ON embeds("guild_id");
`
}

func (s *EmbedsTable) GetEmbed(id int) (embed CustomEmbed, err error) {
	query := `
SELECT 
	"id",
	"guild_id",
	"title",
	"description",
	"url",
	"colour",
	"author_name",
	"author_icon_url",
	"author_url",
	"image_url",
	"thumbnail_url",
	"footer_text",
	"footer_icon_url",
	"timestamp"
FROM embeds
WHERE "id" = $1;
`

	err = s.QueryRow(context.Background(), query, id).Scan(
		&embed.Id,
		&embed.GuildId,
		&embed.Title,
		&embed.Description,
		&embed.Url,
		&embed.Colour,
		&embed.AuthorName,
		&embed.AuthorIconUrl,
		&embed.AuthorUrl,
		&embed.ImageUrl,
		&embed.ThumbnailUrl,
		&embed.FooterText,
		&embed.FooterIconUrl,
		&embed.Timestamp,
	)

	return
}

func (s *EmbedsTable) Create(embed *CustomEmbed) (id int, err error) {
	query := `
INSERT INTO embeds(
	"guild_id",
	"title",
	"description",
	"url",
	"colour",
	"author_name",
	"author_icon_url",
	"author_url",
	"image_url",
	"thumbnail_url",
	"footer_text", 
	"footer_icon_url",
	"timestamp"
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING "id";
`

	err = s.QueryRow(
		context.Background(),
		query,
		embed.GuildId,
		embed.Title,
		embed.Description,
		embed.Url,
		embed.Colour,
		embed.AuthorName,
		embed.AuthorIconUrl,
		embed.AuthorUrl,
		embed.ImageUrl,
		embed.ThumbnailUrl,
		embed.FooterText,
		embed.FooterIconUrl,
		embed.Timestamp,
	).Scan(&id)
	return
}

func (s *EmbedsTable) CreateWithFields(embed *CustomEmbed, fields []EmbedField) (int, error) {
	tx, err := s.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(context.Background())

	id, err := s.CreateWithFieldsTx(tx, embed, fields)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *EmbedsTable) CreateWithFieldsTx(tx pgx.Tx, embed *CustomEmbed, fields []EmbedField) (int, error) {
	query := `
INSERT INTO embeds(
	"guild_id",
	"title",
	"description",
	"url",
	"colour",
	"author_name",
	"author_icon_url",
	"author_url",
	"image_url",
	"thumbnail_url",
	"footer_text", 
	"footer_icon_url",
	"timestamp"
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING "id";
`

	// Create actual embed
	var embedId int
	err := tx.QueryRow(
		context.Background(),
		query,
		embed.GuildId,
		embed.Title,
		embed.Description,
		embed.Url,
		embed.Colour,
		embed.AuthorName,
		embed.AuthorIconUrl,
		embed.AuthorUrl,
		embed.ImageUrl,
		embed.ThumbnailUrl,
		embed.FooterText,
		embed.FooterIconUrl,
		embed.Timestamp,
	).Scan(&embedId)
	if err != nil {
		return 0, err
	}

	// Create fields
	for _, field := range fields {
		query := `
INSERT INTO embed_fields(
	"embed_id",
	"name",
	"value",
	"inline"
)
VALUES($1, $2, $3, $4);
`
		_, err = tx.Exec(context.Background(), query, embedId, field.Name, field.Value, field.Inline)
		if err != nil {
			return 0, err
		}
	}

	return embedId, nil
}

func (s *EmbedsTable) Update(embed *CustomEmbed) error {
	query := `
UPDATE embeds
SET
	"title" = $2,
	"description" = $3,
	"url" = $4,
	"colour" = $5,
	"author_name" = $6,
	"author_icon_url" = $7,
	"author_url" = $8,
	"image_url" = $9,
	"thumbnail_url" = $10,
	"footer_text" = $11, 
	"footer_icon_url" = $12,
	"timestamp" = $13
WHERE "id" = $1;
`

	_, err := s.Exec(
		context.Background(),
		query,
		embed.Id,
		embed.Title,
		embed.Description,
		embed.Url,
		embed.Colour,
		embed.AuthorName,
		embed.AuthorIconUrl,
		embed.AuthorUrl,
		embed.ImageUrl,
		embed.ThumbnailUrl,
		embed.FooterText,
		embed.FooterIconUrl,
		embed.Timestamp,
	)

	return err
}

func (s *EmbedsTable) UpdateWithFields(embed *CustomEmbed, fields []EmbedField) error {
	tx, err := s.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	if err := s.UpdateWithFieldsTx(tx, embed, fields); err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (s *EmbedsTable) UpdateWithFieldsTx(tx pgx.Tx, embed *CustomEmbed, fields []EmbedField) error {
	query := `
UPDATE embeds
SET
	"title" = $2,
	"description" = $3,
	"url" = $4,
	"colour" = $5,
	"author_name" = $6,
	"author_icon_url" = $7,
	"author_url" = $8,
	"image_url" = $9,
	"thumbnail_url" = $10,
	"footer_text" = $11, 
	"footer_icon_url" = $12,
	"timestamp" = $13
WHERE "id" = $1;
`

	tx, err := s.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	// Update actual embed
	_, err = tx.Exec(
		context.Background(),
		query,
		embed.Id,
		embed.Title,
		embed.Description,
		embed.Url,
		embed.Colour,
		embed.AuthorName,
		embed.AuthorIconUrl,
		embed.AuthorUrl,
		embed.ImageUrl,
		embed.ThumbnailUrl,
		embed.FooterText,
		embed.FooterIconUrl,
		embed.Timestamp,
	)

	// Delete and recreate fields
	if _, err := tx.Exec(context.Background(), `DELETE FROM embed_fields WHERE embed_id = $1;`, embed.Id); err != nil {
		return err
	}

	for _, field := range fields {
		query := `
INSERT INTO embed_fields(
	"embed_id",
	"name",
	"value",
	"inline"
)
VALUES($1, $2, $3, $4);
`
		_, err = tx.Exec(context.Background(), query, embed.Id, field.Name, field.Value, field.Inline)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EmbedsTable) Delete(id int) (err error) {
	query := `
DELETE FROM embeds
WHERE "id" = $1;`

	_, err = s.Exec(context.Background(), query, id)
	return
}

func (s *EmbedsTable) DeleteTx(tx pgx.Tx, id int) (err error) {
	query := `
DELETE FROM embeds
WHERE "id" = $1;`

	_, err = tx.Exec(context.Background(), query, id)
	return
}
