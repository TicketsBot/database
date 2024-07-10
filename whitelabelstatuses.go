package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WhitelabelStatuses struct {
	*pgxpool.Pool
}

func newWhitelabelStatuses(db *pgxpool.Pool) *WhitelabelStatuses {
	return &WhitelabelStatuses{
		db,
	}
}

func (w WhitelabelStatuses) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_statuses(
	"bot_id" int8 UNIQUE NOT NULL,
	"status" varchar(255) NOT NULL,
	"status_type" int2 NOT NULL DEFAULT 2,
	FOREIGN KEY("bot_id") REFERENCES whitelabel("bot_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("bot_id")
);
`
}

// Get Returns (status, status_type, exists, error)
func (w *WhitelabelStatuses) Get(ctx context.Context, botId uint64) (string, int16, bool, error) {
	query := `SELECT "status", "status_type" FROM whitelabel_statuses WHERE "bot_id" = $1;`

	var status string
	var statusType int16
	if err := w.QueryRow(ctx, query, botId).Scan(&status, &statusType); err != nil {
		if err == pgx.ErrNoRows {
			return "", 0, false, nil
		} else {
			return "", 0, false, err
		}
	}

	return status, statusType, true, nil
}

func (w *WhitelabelStatuses) Set(ctx context.Context, botId uint64, status string, statusType int16) (err error) {
	query := `
INSERT INTO whitelabel_statuses("bot_id", "status", "status_type")
VALUES($1, $2, $3)
ON CONFLICT("bot_id") DO UPDATE SET "status" = $2, "status_type" = $3;`

	_, err = w.Exec(ctx, query, botId, status, statusType)
	return
}

func (w *WhitelabelStatuses) Delete(ctx context.Context, botId uint64) (err error) {
	query := `DELETE FROM whitelabel_statuses WHERE "bot_id"=$1;`
	_, err = w.Exec(ctx, query, botId)
	fmt.Println()
	return
}
