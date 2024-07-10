package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type StaffOverride struct {
	*pgxpool.Pool
}

func newStaffOverride(db *pgxpool.Pool) *StaffOverride {
	return &StaffOverride{
		db,
	}
}

func (s StaffOverride) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS staff_override(
	"guild_id" int8 NOT NULL UNIQUE,
	"expires" timestamptz NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (s *StaffOverride) HasActiveOverride(ctx context.Context, guildId uint64) (bool, error) {
	query := `
SELECT "expires" 
FROM staff_override
WHERE "guild_id" = $1;
`

	var expires time.Time
	err := s.QueryRow(ctx, query, guildId).Scan(&expires)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return expires.After(time.Now()), nil
}

func (s *StaffOverride) Set(ctx context.Context, guildId uint64, expires time.Time) (err error) {
	query := `
INSERT INTO staff_override("guild_id", "expires")
VALUES($1, $2)
ON CONFLICT("guild_id") DO UPDATE SET "expires" = $2;`

	_, err = s.Exec(ctx, query, guildId, expires)
	return
}

func (s *StaffOverride) Delete(ctx context.Context, guildId uint64) (err error) {
	query := `
DELETE FROM staff_override
WHERE "guild_id" = $1;`

	_, err = s.Exec(ctx, query, guildId)
	return
}
