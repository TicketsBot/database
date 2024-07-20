package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Tag struct {
	Id                   string
	GuildId              uint64
	Content              *string
	Embed                *CustomEmbedWithFields
	ApplicationCommandId *uint64
}

type TagsTable struct {
	*pgxpool.Pool
	repository *Database
}

func newTag(db *pgxpool.Pool) *TagsTable {
	return &TagsTable{
		Pool: db,
	}
}

func (t TagsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS tags(
	"tag_id" varchar(16) NOT NULL,
	"guild_id" int8 NOT NULL,
	"content" text DEFAULT NULL CONSTRAINT content_length CHECK (length(content) <= 4096),
	"embed" JSONB DEFAULT NULL,
	"application_command_id" int8 DEFAULT NULL,
	PRIMARY KEY("guild_id", "tag_id")
);
CREATE INDEX IF NOT EXISTS tags_guild_id_idx ON tags("guild_id");
`
}

func (t *TagsTable) Exists(ctx context.Context, guildId uint64, tagId string) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM tags WHERE "guild_id" = $1 AND LOWER("tag_id") = LOWER($2));`
	err = t.QueryRow(ctx, query, guildId, tagId).Scan(&exists)
	return
}

func (t *TagsTable) Get(ctx context.Context, guildId uint64, tagId string) (Tag, bool, error) {
	query := `
SELECT LOWER(tag_id), "guild_id", "content", "embed", "application_command_id"
FROM tags
WHERE "guild_id" = $1 AND LOWER("tag_id") = LOWER($2);
`

	var tag Tag
	var embedRaw *string
	err := t.QueryRow(ctx, query, guildId, tagId).Scan(
		&tag.Id,
		&tag.GuildId,
		&tag.Content,
		&embedRaw,
		&tag.ApplicationCommandId,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return Tag{}, false, nil
		} else {
			return Tag{}, false, err
		}
	}

	if embedRaw != nil {
		if err := json.UnmarshalFromString(*embedRaw, &tag.Embed); err != nil {
			return Tag{}, false, err
		}
	}

	return tag, true, nil
}

func (t *TagsTable) GetTagIds(ctx context.Context, guildId uint64) (ids []string, e error) {
	query := `SELECT LOWER("tag_id") from tags WHERE "guild_id"=$1;`
	rows, err := t.Query(ctx, query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			e = err
			continue
		}

		ids = append(ids, id)
	}

	return
}

func (t *TagsTable) GetByApplicationCommandId(ctx context.Context, guildId, applicationCommandId uint64) (Tag, bool, error) {
	query := `
SELECT LOWER(tags.tag_id), tags.guild_id, tags.content, tags.embed, tags.application_command_id
FROM tags
WHERE "guild_id" = $1 AND "application_command_id" = $2;
`

	var tag Tag
	var embedRaw *string
	if err := t.QueryRow(ctx, query, guildId, applicationCommandId).Scan(&tag.Id, &tag.GuildId, &tag.Content, &embedRaw, &tag.ApplicationCommandId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tag{}, false, nil
		}

		return Tag{}, false, err
	}

	if embedRaw != nil {
		if err := json.UnmarshalFromString(*embedRaw, &tag.Embed); err != nil {
			return Tag{}, false, err
		}
	}

	return tag, true, nil
}

func (t *TagsTable) GetByGuild(ctx context.Context, guildId uint64) (map[string]Tag, error) {
	query := `
SELECT LOWER(tags.tag_id), tags.guild_id, tags.content, tags.embed, tags.application_command_id
FROM tags
WHERE "guild_id" = $1;`

	rows, err := t.Query(ctx, query, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tags := make(map[string]Tag)
	for rows.Next() {
		var tag Tag
		var embedRaw *string
		if err := rows.Scan(&tag.Id, &tag.GuildId, &tag.Content, &embedRaw, &tag.ApplicationCommandId); err != nil {
			return nil, err
		}

		if embedRaw != nil {
			if err := json.UnmarshalFromString(*embedRaw, &tag.Embed); err != nil {
				return nil, err
			}
		}

		tags[tag.Id] = tag
	}

	return tags, nil
}

func (t *TagsTable) GetTagCount(ctx context.Context, guildId uint64) (count int, err error) {
	query := `SELECT COUNT(*) FROM tags WHERE "guild_id" = $1;`
	err = t.QueryRow(ctx, query, guildId).Scan(&count)
	return
}

func (t *TagsTable) GetStartingWith(ctx context.Context, guildId uint64, prefix string, limit int) (tagIds []string, e error) {
	query := `SELECT LOWER("tag_id") FROM tags WHERE "guild_id"=$1 AND "tag_id" LIKE $2 || '%' LIMIT $3;`
	rows, err := t.Query(ctx, query, guildId, prefix, limit)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		tagIds = append(tagIds, id)
	}

	return
}

func (t *TagsTable) Set(ctx context.Context, tag Tag) error {
	query := `
INSERT INTO tags("tag_id", "guild_id", "content", "embed", "application_command_id")
VALUES(LOWER($1), $2, $3, $4, $5)
ON CONFLICT("tag_id", "guild_id") DO
UPDATE SET "content" = $3, "embed" = $4, "application_command_id" = $5;`

	var embedRaw *string
	if tag.Embed != nil {
		tmp, err := json.MarshalToString(tag.Embed)
		if err != nil {
			return err
		}

		embedRaw = &tmp
	}

	_, err := t.Exec(ctx, query, tag.Id, tag.GuildId, tag.Content, embedRaw, tag.ApplicationCommandId)
	return err
}

func (t *TagsTable) Delete(ctx context.Context, guildId uint64, tagId string) (err error) {
	query := `
DELETE FROM tags 
WHERE "guild_id" = $1 AND "tag_id" = LOWER($2);`

	_, err = t.Exec(ctx, query, guildId, tagId)
	return
}
