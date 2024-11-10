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

func (b *GlobalBlacklist) IsBlacklisted(ctx context.Context, userId uint64) (blacklisted bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1 FROM global_blacklist WHERE "user_id" = $1
);
`

	err = b.QueryRow(ctx, query, userId).Scan(&blacklisted)
	return
}

func (b *GlobalBlacklist) ListAll(ctx context.Context) (users []uint64, err error) {
	query := `SELECT "user_id" FROM global_blacklist;`

	rows, err := b.Query(ctx, query)
	if err != nil {
		return
	}

	for rows.Next() {
		var userId uint64
		if err = rows.Scan(&userId); err != nil {
			return
		}

		users = append(users, userId)
	}

	return
}

func (b *GlobalBlacklist) Add(ctx context.Context, userId uint64) (err error) {
	_, err = b.Exec(ctx, `INSERT INTO global_blacklist("user_id") VALUES($1) ON CONFLICT("user_id") DO NOTHING;`, userId)
	return
}

func (b *GlobalBlacklist) Delete(ctx context.Context, userId uint64) (err error) {
	_, err = b.Exec(ctx, `DELETE FROM global_blacklist WHERE "user_id" = $1;`, userId)
	return
}
