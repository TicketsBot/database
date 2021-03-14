package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type WhitelabelUsers struct {
	*pgxpool.Pool
}

func newWhitelabelUsers(db *pgxpool.Pool) *WhitelabelUsers {
	return &WhitelabelUsers{
		db,
	}
}

func (p WhitelabelUsers) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS whitelabel_users(
	"user_id" int8 NOT NULL UNIQUE,
	"expiry" timestamp NOT NULL,
	PRIMARY KEY("user_id")
);`
}

func (p *WhitelabelUsers) IsPremium(guildId uint64) (bool, error) {
	expiry, err := p.GetExpiry(guildId)
	if err != nil {
		return false, err
	}

	return expiry.After(time.Now()), nil
}

func (p *WhitelabelUsers) GetExpiry(userId uint64) (expiry time.Time, e error) {
	if err := p.QueryRow(context.Background(), `SELECT "expiry" from whitelabel_users WHERE "user_id" = $1;`, userId).Scan(&expiry); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *WhitelabelUsers) Add(userId uint64, interval time.Duration) (err error) {
	query := `
INSERT INTO whitelabel_users("user_id", "expiry")
VALUES($1, NOW() + $2)
ON CONFLICT("user_id") DO
UPDATE SET "expiry" = whitelabel_users.expiry + $2;`

	_, err = p.Exec(context.Background(), query, userId, interval)
	return
}
