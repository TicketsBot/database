package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OnCall struct {
	*pgxpool.Pool
}

func newOnCall(db *pgxpool.Pool) *OnCall {
	return &OnCall{
		db,
	}
}

func (b OnCall) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS on_call(
	"guild_id" int8 NOT NULL,
	"user_id" int8 NOT NULL,
	"is_on_call" bool NOT NULL,
	PRIMARY KEY("guild_id", "user_id")
);`
}

func (b *OnCall) IsOnCall(ctx context.Context, guildId, userId uint64) (bool, error) {
	query := `SELECT "is_on_call" FROM on_call WHERE "guild_id" = $1 AND "user_id" = $2;`

	var onCall bool
	if err := b.QueryRow(ctx, query, guildId, userId).Scan(&onCall); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return onCall, nil
}

func (b *OnCall) GetUsersOnCall(ctx context.Context, guildId uint64) ([]uint64, error) {
	query := `SELECT "user_id" FROM on_call WHERE "guild_id" = $1 AND "is_on_call" = true;`

	rows, err := b.Query(ctx, query, guildId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var users []uint64
	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}

		users = append(users, userId)
	}

	return users, nil
}

func (b *OnCall) GetOnCallCount(ctx context.Context, guildId uint64) (count int, err error) {
	query := `SELECT COUNT(1) FROM on_call WHERE "guild_id" = $1;`

	err = b.QueryRow(ctx, query, guildId).Scan(&count)
	return
}

func (b *OnCall) Toggle(ctx context.Context, guildId, userId uint64) (onCall bool, err error) {
	query := `
INSERT INTO on_call("guild_id", "user_id", "is_on_call") 
VALUES($1, $2, true)
ON CONFLICT ("guild_id", "user_id") 
DO UPDATE SET "is_on_call" = NOT on_call.is_on_call
RETURNING "is_on_call";`

	err = b.QueryRow(ctx, query, guildId, userId).Scan(&onCall)
	return
}

func (b *OnCall) Remove(ctx context.Context, guildId, userId uint64) (err error) {
	query := `DELETE FROM on_call WHERE "guild_id" = $1 AND "user_id" = $2;`
	_, err = b.Exec(ctx, query, guildId, userId)
	return
}
