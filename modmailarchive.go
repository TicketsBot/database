package database

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type ModmailArchive struct {
	Uuid      uuid.UUID
	GuildId   uint64
	UserId    uint64
	CloseTime time.Time
}

type ModmailArchiveTable struct {
	*pgxpool.Pool
}

func newModmailArchiveTable(db *pgxpool.Pool) *ModmailArchiveTable {
	return &ModmailArchiveTable{
		db,
	}
}

func (m ModmailArchiveTable) Schema() string {
	return `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS modmail_archive(
	"uuid" uuid NOT NULL UNIQUE,
	"guild_id" int8 NOT NULL,
	"user_id" int8 NOT NULL,
	"close_time" timestamp NOT NULL,
	PRIMARY KEY("uuid")
);`
}

func (m *ModmailArchiveTable) Get(uuid uuid.UUID) (archive ModmailArchive, e error) {
	query := `SELECT * from modmail_archive WHERE "uuid" = $1;`
	if err := m.QueryRow(context.Background(), query, uuid).Scan(&archive.Uuid, &archive.GuildId, &archive.Uuid, &archive.CloseTime); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (m *ModmailArchiveTable) GetByGuild(guildId uint64, limit int, after string) (archives []ModmailArchive, e error) {
	var query string
	if after == "" {
		query = `SELECT * from modmail_archive WHERE "guild_id" = $1 ORDER BY "close_time" DESC LIMIT $2;`
	} else {
		query = `SELECT * from modmail_archive WHERE "guild_id" = $1 AND "close_time" < SELECT("close_time" FROM modmail_archive WHERE "uuid" = $3 AND "guild_id" = $1 LIMIT 1) ORDER BY "close_time" DESC LIMIT $2;`
	}

	rows, err := m.Query(context.Background(), query, guildId, limit, after)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var archive ModmailArchive
		if err := rows.Scan(&archive.Uuid, &archive.GuildId, &archive.UserId, &archive.CloseTime); err != nil {
			e = err
			continue
		}

		archives = append(archives, archive)
	}

	return
}

func (m *ModmailArchiveTable) GetByUser(userId uint64, limit int, after string) (archives []ModmailArchive, e error) {
	var query string
	if after == "" {
		query = `SELECT * from modmail_archive WHERE "user_id" = $1 ORDER BY "close_time" DESC LIMIT $2;`
	} else {
		query = `SELECT * from modmail_archive WHERE "user_id" = $1 AND "close_time" < SELECT("close_time" FROM modmail_archive WHERE "uuid" = $3 AND "user_id" = $1 LIMIT 1) ORDER BY "close_time" DESC LIMIT $2;`
	}

	rows, err := m.Query(context.Background(), query, userId, limit, after)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var archive ModmailArchive
		if err := rows.Scan(&archive.Uuid, &archive.GuildId, &archive.UserId, &archive.CloseTime); err != nil {
			e = err
			continue
		}

		archives = append(archives, archive)
	}

	return
}

func (m *ModmailArchiveTable) Set(archive ModmailArchive) (err error) {
	query := `INSERT INTO modmail_archive("uuid", "guild_id", "user_id", "close_time") VALUES($1, $2, $3, $4) ON CONFLICT("uuid") DO NOTHING;`
	_, err = m.Exec(context.Background(), query, archive.Uuid, archive.GuildId, archive.UserId, archive.CloseTime)
	return
}
