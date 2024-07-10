package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BotStaff struct {
	*pgxpool.Pool
}

func newBotStaff(db *pgxpool.Pool) *BotStaff {
	return &BotStaff{
		db,
	}
}

func (s BotStaff) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS bot_staff(
	"user_id" int8 NOT NULL UNIQUE,
	PRIMARY KEY("user_id")
);`
}

func (s *BotStaff) IsStaff(ctx context.Context, userId uint64) (isStaff bool, err error) {
	query := `
SELECT EXISTS (
	SELECT 1
	FROM bot_staff
	where "user_id" = $1
);
`

	err = s.QueryRow(ctx, query, userId).Scan(&isStaff)
	return
}

func (s *BotStaff) GetAll(ctx context.Context) ([]uint64, error) {
	query := `SELECT "user_id" FROM bot_staff;`

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userIds []uint64
	for rows.Next() {
		var userId uint64
		if err = rows.Scan(&userId); err != nil {
			return nil, err
		}

		userIds = append(userIds, userId)
	}

	return userIds, nil
}

func (s *BotStaff) Add(ctx context.Context, userId uint64) (err error) {
	query := `
INSERT INTO bot_staff("user_id")
VALUES($1)
ON CONFLICT("user_id") DO NOTHING;
`

	_, err = s.Exec(ctx, query, userId)
	return
}

func (s *BotStaff) Delete(ctx context.Context, userId uint64) (err error) {
	query := `
DELETE FROM bot_staff
WHERE "user_id" = $1;`

	_, err = s.Exec(ctx, query, userId)
	return
}
