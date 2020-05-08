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
	"user_id" int8 NOT NULL UNIQUE,
	"staff_channel" int8 NOT NULL UNIQUE,
	"welcome_message_id" int8 NOT NULL UNIQUE,
	PRIMARY KEY("uuid")
);
`
}

func (m *ModmailSessionTable) GetByUser(userId uint64) (session ModmailSession, e error) {
	query := `SELECT * from modmail_sessions WHERE "user_id" = $1;`
	if err := m.QueryRow(context.Background(), query, userId).Scan(&session.Uuid, &session.GuildId, &session.UserId, &session.StaffChannelId, &session.WelcomeMessageId); err != nil && err != pgx.ErrNoRows {
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

func (m *ModmailSessionTable) Create(session ModmailSession) (err error) {
	query := `INSERT INTO modmail_sessions("guild_id", "user_id", "staff_channel", "welcome_message_id") VALUES($1, $2, $3, $4) ON CONFLICT("uuid") DO NOTHING;`
	_, err = m.Exec(context.Background(), query, session.GuildId, session.UserId, session.StaffChannelId, session.WelcomeMessageId)
	return
}

func (m *ModmailSessionTable) DeleteByUser(userId uint64) (err error) {
	query := `DELETE FROM modmail_sessions WHERE "user_id" = $1;`
	_, err = m.Exec(context.Background(), query, userId)
	return
}
