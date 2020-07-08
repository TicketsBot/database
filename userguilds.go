package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserGuild struct {
	GuildId         uint64
	Name            string
	Owner           bool
	UserPermissions int32
	Icon            string
}

type UserGuildsTable struct {
	*pgxpool.Pool
}

func newUserGuildsTable(db *pgxpool.Pool) *UserGuildsTable {
	return &UserGuildsTable{
		db,
	}
}

func (u UserGuildsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS user_guilds(
	"user_id" int8 NOT NULL,
	"guild_id" int8 NOT NULL,
	"name" varchar(100) NOT NULL,
	"owner" bool NOT NULL,
	"permissions" int4 NOT NULL,
	"icon" varchar(34),
	PRIMARY KEY("user_id", "guild_id")
);`
}

func (u *UserGuildsTable) Get(userId uint64) (guilds []UserGuild, e error) {
	query := `SELECT "guild_id", "name", "owner", "permissions", "icon" FROM user_guilds WHERE "user_id" = $1;`

	rows, err := u.Query(context.Background(), query, userId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
		return
	}

	for rows.Next() {
		var guild UserGuild
		if err := rows.Scan(&guild.GuildId, &guild.Name, &guild.Owner, &guild.UserPermissions, &guild.Icon); err != nil {
			e = err
			continue
		}

		guilds = append(guilds, guild)
	}

	return
}

func (u *UserGuildsTable) Set(userId uint64, guilds []UserGuild) (err error) {
	// create slice of guild ids
	var guildIds []uint64
	for _, guild := range guilds {
		guildIds = append(guildIds, guild.GuildId)
	}

	guildIdArray := &pgtype.Int8Array{}
	if err = guildIdArray.Set(guildIds); err != nil {
		return
	}

	batch := &pgx.Batch{}

	batch.Queue(`DELETE FROM user_guilds WHERE "user_id" = $1 AND NOT ("guild_id" = ANY($2));`, userId, guildIdArray)

	for _, guild := range guilds {
		query := `INSERT INTO user_guilds("user_id", "guild_id", "name", "owner", "permissions", "icon") VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT("user_id", "guild_id") DO UPDATE SET "name" = $3, "owner" = $4, "permissions" = $5, "icon" = $6;`
		batch.Queue(query, userId, guild.GuildId, guild.Name, guild.Owner, guild.UserPermissions, guild.Icon)
	}

	br := u.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err = br.Exec()

	return
}
