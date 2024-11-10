package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServerBlacklist struct {
	*pgxpool.Pool
}

func newServerBlacklist(db *pgxpool.Pool) *ServerBlacklist {
	return &ServerBlacklist{
		db,
	}
}

func (b ServerBlacklist) Schema() string {
	return `CREATE TABLE IF NOT EXISTS server_blacklist("guild_id" int8 NOT NULL UNIQUE, PRIMARY KEY("guild_id"));`
}

func (b *ServerBlacklist) IsBlacklisted(ctx context.Context, guildId uint64) (bool, *string, error) {
	query := `SELECT "reason" FROM server_blacklist WHERE "guild_id" = $1;`

	var reason *string
	if err := b.QueryRow(ctx, query, guildId).Scan(&reason); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, nil
		} else {
			return false, nil, err
		}
	}

	return true, reason, nil
}

func (b *ServerBlacklist) ListAll(ctx context.Context) ([]uint64, error) {
	query := `SELECT "guild_id" FROM server_blacklist;`

	rows, err := b.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var guilds []uint64
	for rows.Next() {
		var guildId uint64
		if err := rows.Scan(&guildId); err != nil {
			return nil, err
		}

		guilds = append(guilds, guildId)
	}

	return guilds, nil
}

func (b *ServerBlacklist) Add(ctx context.Context, guildId uint64, reason *string) (err error) {
	query := `INSERT INTO server_blacklist("guild_id", "reason") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "reason" = $2`
	_, err = b.Exec(ctx, query, guildId, reason)
	return
}

func (b *ServerBlacklist) Delete(ctx context.Context, guildId uint64) (err error) {
	_, err = b.Exec(ctx, `DELETE FROM server_blacklist WHERE "guild_id" = $1;`, guildId)
	return
}
