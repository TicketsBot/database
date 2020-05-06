package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Votes struct {
	*pgxpool.Pool
}

func newVotes(db *pgxpool.Pool) *Votes {
	return &Votes{
		db,
	}
}

func (v Votes) Schema() string {
	return `CREATE TABLE IF NOT EXISTS votes("user_id" int8 NOT NULL UNIQUE, "vote_time" timestamp NOT NULL, PRIMARY KEY("user_id"));`
}

func (v *Votes) Get(userId uint64) (voteTime time.Time, e error) {
	query := `SELECT "vote_time" from votes WHERE "user_id" = $1`

	if err := v.QueryRow(context.Background(), query, userId).Scan(&voteTime); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (v *Votes) Set(userId uint64) (err error) {
	query := `INSERT INTO votes("guild_id", "vote_time") VALUES($1, NOW()) ON CONFLICT("guild_id") DO UPDATE SET "vote_time" = NOW();`
	_, err = v.Exec(context.Background(), query, userId)
	return
}
