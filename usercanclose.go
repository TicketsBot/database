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
	return `
CREATE TABLE IF NOT EXISTS users_can_close(
	"guild_id" int8 NOT NULL UNIQUE,
	"users_can_close" bool NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (u *UsersCanClose) Get(ctx context.Context, guildId uint64) (usersCanClose bool, e error) {
	if err := u.QueryRow(ctx, `SELECT "users_can_close" from users_can_close WHERE "guild_id" = $1;`, guildId).Scan(&usersCanClose); err != nil {
		if err == pgx.ErrNoRows {
			usersCanClose = true
		} else {
			e = err
		}
	}

	return
}

func (u *UsersCanClose) Set(ctx context.Context, guildId uint64, usersCanClose bool) (err error) {
	_, err = u.Exec(ctx, `INSERT INTO users_can_close("guild_id", "users_can_close") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "users_can_close" = $2;`, guildId, usersCanClose)
	return
}
