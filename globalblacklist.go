package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type GlobalBlacklist struct {
	*pgxpool.Pool
}

func newGlobalBlacklist(db *pgxpool.Pool) *GlobalBlacklist {
	return &GlobalBlacklist{
		db,
	}
}

func (b GlobalBlacklist) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS global_blacklist(
	"user_id" int8 NOT NULL UNIQUE,
	PRIMARY KEY("user_id")
);
`
}

func (b *GlobalBlacklist) IsBlacklisted(userId uint64) (blacklisted bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1 FROM global_blacklist WHERE "user_id" = $1
);
`

	err = b.QueryRow(context.Background(), query, userId).Scan(&blacklisted)
	return
}

func (b *GlobalBlacklist) Add(userId uint64) (err error) {
	_, err = b.Exec(context.Background(), `INSERT INTO global_blacklist("user_id") VALUES($1) ON CONFLICT("user_id") DO NOTHING;`, userId)
	return
}

func (b *GlobalBlacklist) Delete(userId uint64) (err error) {
	_, err = b.Exec(context.Background(), `DELETE FROM global_blacklist WHERE "user_id" = $1;`, userId)
	return
}
