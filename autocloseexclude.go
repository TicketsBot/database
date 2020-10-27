package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AutoCloseExclude struct {
	*pgxpool.Pool
}

func newAutoCloseExclude(db *pgxpool.Pool) *AutoCloseExclude {
	return &AutoCloseExclude{
		db,
	}
}

func (a AutoCloseExclude) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS auto_close_exclude(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (a *AutoCloseExclude) IsExcluded(guildId uint64, ticketId int) (excluded bool, e error) {
	query := `
SELECT COUNT(*)
FROM auto_close_exclude
WHERE "guild_id" = $1 AND "ticket_id" = $2
;
`

	var count int
	if err := a.QueryRow(context.Background(), query, guildId, ticketId).Scan(&count); err != nil {
		e = err
	}

	excluded = count > 0

	return
}

func (a *AutoCloseExclude) Exclude(guildId uint64, ticketId int) (err error) {
	query := `
INSERT INTO auto_close_exclude("guild_id", "ticket_id")
VALUES ($1, $2)
ON CONFLICT("guild_id", "ticket_id") DO NOTHING
;
`

	_, err = a.Exec(context.Background(), query, guildId, ticketId)
	return
}
