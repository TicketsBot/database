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

func transact(pool *pgxpool.Pool, statements ...string) (pgx.Tx, error) {
	tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return tx, err
	}

	for _, statement := range statements {
		if _, err := tx.Exec(context.Background(), statement); err != nil {
			return tx, err
		}
	}

	return tx, nil
}
