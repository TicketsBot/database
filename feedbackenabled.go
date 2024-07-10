package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FeedbackEnabled struct {
	*pgxpool.Pool
}

func newFeedbackEnabled(db *pgxpool.Pool) *FeedbackEnabled {
	return &FeedbackEnabled{
		db,
	}
}

func (FeedbackEnabled) Schema() string {
	return `CREATE TABLE IF NOT EXISTS feedback_enabled("guild_id" int8 NOT NULL UNIQUE, "feedback_enabled" bool NOT NULL, PRIMARY KEY("guild_id"));`
}

func (f *FeedbackEnabled) Get(ctx context.Context, guildId uint64) (feedbackEnabled bool, e error) {
	if err := f.QueryRow(ctx, `SELECT "feedback_enabled" from feedback_enabled WHERE "guild_id" = $1;`, guildId).Scan(&feedbackEnabled); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (f *FeedbackEnabled) Set(ctx context.Context, guildId uint64, feedbackEnabled bool) (err error) {
	_, err = f.Exec(ctx, `INSERT INTO feedback_enabled("guild_id", "feedback_enabled") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "feedback_enabled" = $2;`, guildId, feedbackEnabled)
	return
}
