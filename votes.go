package database

import (
	"context"
	"github.com/jackc/pgx/pgtype"
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
	return `
CREATE TABLE IF NOT EXISTS votes(
	"user_id" int8 NOT NULL UNIQUE,
	"vote_time" timestamp NOT NULL,
	PRIMARY KEY("user_id")
);`
}

func (v *Votes) Get(ctx context.Context, userId uint64) (voteTime time.Time, e error) {
	query := `SELECT "vote_time" from votes WHERE "user_id" = $1`

	if err := v.QueryRow(ctx, query, userId).Scan(&voteTime); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (v *Votes) Any(ctx context.Context, userIds ...uint64) (bool, error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM votes
	WHERE "user_id" = ANY($1) AND vote_time > NOW() - INTERVAL '24 hours'
);
`

	userIdArray := &pgtype.Int8Array{}
	if err := userIdArray.Set(userIds); err != nil {
		return false, err
	}

	var res bool
	if err := v.QueryRow(ctx, query, userIdArray).Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}

func (v *Votes) Set(ctx context.Context, userId uint64) (err error) {
	query := `INSERT INTO votes("user_id", "vote_time") VALUES($1, NOW()) ON CONFLICT("user_id") DO UPDATE SET "vote_time" = NOW();`
	_, err = v.Exec(ctx, query, userId)
	return
}
