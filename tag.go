package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Tag struct {
	*pgxpool.Pool
}

func newTag(db *pgxpool.Pool) *Tag {
	return &Tag{
		db,
	}
}

func (t Tag) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS tags(
	"guild_id" int8 NOT NULL,
	"tag_id" varchar(16) NOT NULL,
	"content" text NOT NULL,
	PRIMARY KEY("guild_id", "tag_id")
);`
}

func (t *Tag) Get(guildId uint64, tagId string) (content string, e error) {
	query := `SELECT "content" from tags WHERE "guild_id"=$1 AND LOWER("tag_id")=LOWER($2);`
	if err := t.QueryRow(context.Background(), query, guildId, tagId).Scan(&content); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (t *Tag) GetTagIds(guildId uint64) (ids []string, e error) {
	query := `SELECT LOWER("tag_id") from tags WHERE "guild_id"=$1;`
	rows, err := t.Query(context.Background(), query, guildId)
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

func (t *Tag) GetByGuild(guildId uint64) (tags map[string]string, e error) {
	tags = make(map[string]string)

	query := `SELECT LOWER("tag_id"), "content" from tags WHERE "guild_id"=$1;`
	rows, err := t.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var id, content string
		if err := rows.Scan(&id, &content); err != nil {
			e = err
			continue
		}

		tags[id] = content
	}

	return
}

func (t *Tag) GetStartingWith(guildId uint64, prefix string, limit int) (tagIds []string, e error) {
	query := `SELECT LOWER("tag_id"), "content" from tags WHERE "guild_id"=$1 AND "tag_id" LIKE '%$2' LIMIT $3;`
	rows, err := t.Query(context.Background(), query, guildId, prefix, limit)
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

		tagIds = append(tagIds, id)
	}

	return
}

func (t *Tag) Set(guildId uint64, tagId, content string) (err error) {
	query := `INSERT INTO tags("guild_id", "tag_id", "content") VALUES($1, LOWER($2), $3) ON CONFLICT("guild_id", "tag_id") DO UPDATE SET "content"=$3;`
	_, err = t.Exec(context.Background(), query, guildId, tagId, content)
	return
}

func (t *Tag) Delete(guildId uint64, tagId string) (err error) {
	query := `DELETE FROM tags WHERE "guild_id"=$1 AND "tag_id"=LOWER($2);`
	_, err = t.Exec(context.Background(), query, guildId, tagId)
	return
}
