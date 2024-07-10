package database

import (
	"context"
	"github.com/jackc/pgx/pgtype"
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

func (p *WhitelabelUsers) IsPremium(ctx context.Context, userId uint64) (bool, error) {
	expiry, err := p.GetExpiry(ctx, userId)
	if err != nil {
		return false, err
	}

	return expiry.After(time.Now()), nil
}

func (p *WhitelabelUsers) AnyPremium(ctx context.Context, userIds []uint64) (bool, error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM whitelabel_users
	WHERE "user_id" = ANY($1) AND "expiry" > NOW()
);
`

	userIdArray := &pgtype.Int8Array{}
	if err := userIdArray.Set(userIds); err != nil {
		return false, err
	}

	var res bool
	if err := p.QueryRow(ctx, query, userIdArray).Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}

func (p *WhitelabelUsers) GetExpiry(ctx context.Context, userId uint64) (expiry time.Time, e error) {
	if err := p.QueryRow(ctx, `SELECT "expiry" from whitelabel_users WHERE "user_id" = $1;`, userId).Scan(&expiry); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *WhitelabelUsers) Add(ctx context.Context, userId uint64, interval time.Duration) (err error) {
	query := `
INSERT INTO whitelabel_users("user_id", "expiry")
VALUES($1, NOW() + $2)
ON CONFLICT("user_id")
DO UPDATE SET "expiry" = CASE WHEN whitelabel_users.expiry < NOW()
	THEN NOW() + $2
	ELSE whitelabel_users.expiry + $2
END;`

	_, err = p.Exec(ctx, query, userId, interval)
	return
}
