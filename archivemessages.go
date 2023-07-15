package database

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ArchiveMessage struct {
	ChannelId uint64 `json:"channel_id,string"`
	MessageId uint64 `json:"message_id,string"`
}

type ArchiveMessages struct {
	*pgxpool.Pool
}

func newArchiveMessages(db *pgxpool.Pool) *ArchiveMessages {
	return &ArchiveMessages{
		db,
	}
}

var (
	//go:embed sql/archive_messages/schema.sql
	archiveMessagesSchema string

	//go:embed sql/archive_messages/insert.sql
	archiveMessagesInsert string

	//go:embed sql/archive_messages/get.sql
	archiveMessagesGet string
)

func (a *ArchiveMessages) Schema() string {
	return archiveMessagesSchema
}

func (a *ArchiveMessages) Set(guildId uint64, ticketId int, channelId, messageId uint64) error {
	_, err := a.Exec(context.Background(), archiveMessagesInsert, guildId, ticketId, channelId, messageId)
	return err
}

func (a *ArchiveMessages) Get(guildId uint64, ticketId int) (ArchiveMessage, bool, error) {
	var data ArchiveMessage
	err := a.QueryRow(context.Background(), archiveMessagesGet, guildId, ticketId).Scan(&data.ChannelId, &data.MessageId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ArchiveMessage{}, false, nil
		} else {
			return ArchiveMessage{}, false, err
		}
	}

	return data, true, nil
}
