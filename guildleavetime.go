package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type GuildLeaveTime struct {
	*pgxpool.Pool
}

func newGuildLeaveTime(db *pgxpool.Pool) *GuildLeaveTime {
	return &GuildLeaveTime{
		db,
	}
}

func (c GuildLeaveTime) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS guild_leave_time(
	"guild_id" int8 NOT NULL UNIQUE,
	"leave_time" timestamptz NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (c *GuildLeaveTime) GetBefore(before time.Duration) (ids []uint64, e error) {
	query := `
SELECT "guild_id"
FROM guild_leave_time
WHERE "leave_time" < NOW() - $1::interval;
`

	rows, err := c.Query(context.Background(), query, before)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id uint64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return
}

func (c *GuildLeaveTime) Set(guildId uint64) (err error) {
	_, err = c.Exec(context.Background(), `INSERT INTO guild_leave_time("guild_id", "leave_time") VALUES($1, NOW()) ON CONFLICT("guild_id") DO UPDATE SET "leave_time" = NOW();`, guildId, confirm)
	return
}

func (c *GuildLeaveTime) Delete(guildId uint64) (err error) {
	_, err = c.Exec(context.Background(), `DELETE FROM guild_leave_time WHERE "guild_id" = $1;`, guildId)
	return
}

func (c *GuildLeaveTime) DeleteAll(guildIds []uint64) (err error) {
	array := &pgtype.Int8Array{}
	if err = array.Set(guildIds); err != nil {
		return
	}

	_, err = c.Exec(context.Background(), `DELETE FROM guild_leave_time WHERE "guild_id" = ANY($1);`, array)
	return
}
