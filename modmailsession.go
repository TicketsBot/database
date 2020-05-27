package database

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ModmailSession struct {
	Uuid             uuid.UUID
	GuildId          uint64
	BotId            uint64
	UserId           uint64
	StaffChannelId   uint64
	WelcomeMessageId uint64
}

type ModmailSessionTable struct {
	*pgxpool.Pool
}

func newModmailSessionTable(db *pgxpool.Pool) *ModmailSessionTable {
	return &ModmailSessionTable{
		db,
	}
}

func (m ModmailSessionTable) Schema() string {
	return `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS modmail_sessions(
	"uuid" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4 (),
	"guild_id" int8 NOT NULL,
	"bot_id" int8 NOT NULL,
	"user_id" int8 NOT NULL,
	"staff_channel" int8 NOT NULL UNIQUE,
	"welcome_message_id" int8 NOT NULL UNIQUE,
	UNIQUE("bot_id", "user_id"),
	PRIMARY KEY("uuid")
);
`
}

func (m *ModmailSessionTable) GetByUser(botId, userId uint64) (session ModmailSession, e error) {
	query := `SELECT * from modmail_sessions WHERE "bot_id" = $1 AND "user_id" = $2;`
	if err := m.QueryRow(context.Background(), query, botId, userId).Scan(&session.Uuid, &session.GuildId, &session.UserId, &session.StaffChannelId, &session.WelcomeMessageId); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (m *ModmailSessionTable) GetByChannel(channelId uint64) (session ModmailSession, e error) {
	query := `SELECT * from modmail_sessions WHERE "staff_channel" = $1;`
	if err := m.QueryRow(context.Background(), query, channelId).Scan(&session.Uuid, &session.GuildId, &session.UserId, &session.StaffChannelId, &session.WelcomeMessageId); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (m *ModmailSessionTable) Create(session ModmailSession) (uuid uuid.UUID, err error) {
	query := `INSERT INTO modmail_sessions("guild_id", "user_id", "staff_channel", "welcome_message_id") VALUES($1, $2, $3, $4) ON CONFLICT("uuid") DO NOTHING RETURNING "uuid";`
	err = m.QueryRow(context.Background(), query, session.GuildId, session.UserId, session.StaffChannelId, session.WelcomeMessageId).Scan(&uuid)
	return
}

func (m *ModmailSessionTable) DeleteByUser(botId, userId uint64) (err error) {
	query := `DELETE FROM modmail_sessions WHERE "bot_id" = $1 AND "user_id" = $2;`
	_, err = m.Exec(context.Background(), query, botId, userId)
	return
}
