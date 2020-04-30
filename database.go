package database

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	ArchiveChannel  *ArchiveChannel
	Blacklist       *Blacklist
	ChannelCategory *ChannelCategory
	Prefix          *Prefix
	Tag             *Tag
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		ArchiveChannel:  newArchiveChannel(pool),
		Blacklist:       newBlacklist(pool),
		ChannelCategory: newChannelCategory(pool),
		Prefix:          newPrefix(pool),
		Tag:             newTag(pool),
	}
}
