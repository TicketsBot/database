package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RoleBlacklist struct {
	*pgxpool.Pool
}

func newRoleBlacklist(db *pgxpool.Pool) *RoleBlacklist {
	return &RoleBlacklist{
		db,
	}
}

func (b RoleBlacklist) Schema() string {
	return `CREATE TABLE IF NOT EXISTS role_blacklist("guild_id" int8 NOT NULL, "role_id" int8 NOT NULL, PRIMARY KEY("guild_id", "role_id"));`
}

func (b *RoleBlacklist) IsBlacklisted(guildId, roleId uint64) (exists bool, e error) {
	query := `SELECT EXISTS(SELECT 1 FROM role_blacklist WHERE "guild_id"=$1 AND "role_id"=$2);`
	if err := b.QueryRow(context.Background(), query, guildId, roleId).Scan(&exists); err != nil {
		e = err
	}

	return
}

func (b *RoleBlacklist) GetBlacklistedRoles(guildId uint64) (roles []uint64, e error) {
	query := `SELECT "role_id" FROM role_blacklist WHERE "guild_id" = $1;`

	rows, err := b.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			e = err
			continue
		}

		roles = append(roles, roleId)
	}

	return
}

func (b *RoleBlacklist) Add(guildId, roleId uint64) (err error) {
	// on conflict, role is already role_blacklist
	query := `INSERT INTO role_blacklist("guild_id", "role_id") VALUES($1, $2) ON CONFLICT DO NOTHING;`
	_, err = b.Exec(context.Background(), query, guildId, roleId)
	return
}

func (b *RoleBlacklist) Remove(guildId, roleId uint64) (err error) {
	query := `DELETE FROM role_blacklist WHERE "guild_id"=$1 AND "role_id"=$2;`
	_, err = b.Exec(context.Background(), query, guildId, roleId)
	return
}
