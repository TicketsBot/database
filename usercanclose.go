package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UsersCanClose struct {
	*pgxpool.Pool
}

func newUsersCanClose(db *pgxpool.Pool) *UsersCanClose {
	return &UsersCanClose{
		db,
	}
}

func (u UsersCanClose) Schema() string {
	return `CREATE TABLE IF NOT EXISTS users_can_close("guild_id" int8 NOT NULL UNIQUE, "users_can_close" bool NOT NULL, PRIMARY KEY("guild_id"));`
}

func (u *UsersCanClose) Get(guildId uint64) (usersCanClose bool, e error) {
	if err := u.QueryRow(context.Background(), `SELECT "users_can_close" from users_can_close WHERE "guild_id" = $1;`, guildId).Scan(&u); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (u *UsersCanClose) Set(guildId, usersCanClose bool) (err error) {
	_, err = u.Exec(context.Background(), `INSERT INTO users_can_close("guild_id", "users_can_close") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "users_can_close" = $2;`, guildId, usersCanClose)
	return
}
