package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

func toInterval(duration time.Duration) (interval pgtype.Interval, err error) {
	err = interval.Set(duration)
	return
}

func transact(ctx context.Context, pool *pgxpool.Pool, statements ...string) (pgx.Tx, error) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return tx, err
	}

	for _, statement := range statements {
		if _, err := tx.Exec(ctx, statement); err != nil {
			return tx, err
		}
	}

	return tx, nil
}

func slice[T any](v ...T) []T {
	return v
}

func ptr[T any](v T) *T {
	return &v
}
