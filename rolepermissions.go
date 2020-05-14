package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RolePermissions struct {
	*pgxpool.Pool
}

func newRolePermissions(db *pgxpool.Pool) *RolePermissions {
	return &RolePermissions{
		db,
	}
}

func (p RolePermissions) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS role_permissions(
	"guild_id" int8 NOT NULL,
	"role_id" int8 NOT NULL,
	"support" bool NOT NULL,
	"admin" bool NOT NULL,
	PRIMARY KEY("role_id")
);
CREATE INDEX IF NOT EXISTS role_permissions_guild_id ON role_permissions("guild_id");
`
}

func (p *RolePermissions) IsSupport(roleId uint64) (support bool, e error) {
	var admin bool

	if err := p.QueryRow(context.Background(), `SELECT "support", "admin" from role_permissions WHERE "role_id" = $1;`, roleId).Scan(&support, &admin); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	if admin {
		support = true
	}

	return
}

func (p *RolePermissions) IsAdmin(roleId uint64) (admin bool, e error) {
	if err := p.QueryRow(context.Background(), `SELECT "admin" from role_permissions WHERE "role_id" = $1;`, roleId).Scan(&admin); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *RolePermissions) GetAdminRoles(guildId uint64) (adminRoles []uint64, e error) {
	rows, err := p.Query(context.Background(), `SELECT "role_id" from role_permissions WHERE "guild_id" = $1 AND "admin" = true;`, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			e = err
			continue
		}

		adminRoles = append(adminRoles, roleId)
	}

	return
}

func (p *RolePermissions) GetSupportRoles(guildId uint64) (supportRoles []uint64, e error) {
	rows, err := p.Query(context.Background(), `SELECT "role_id" from role_permissions WHERE "guild_id" = $1 AND ("admin" = true OR "support" = true);`, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			e = err
			continue
		}

		supportRoles = append(supportRoles, roleId)
	}

	return
}

func (p *RolePermissions) GetSupportRolesOnly(guildId uint64) (supportRoles []uint64, e error) {
	query := `SELECT "role_id" from role_permissions WHERE "guild_id" = $1 AND "admin" = false AND "support" = true;`
	rows, err := p.Query(context.Background(), query, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			e = err
			continue
		}

		supportRoles = append(supportRoles, roleId)
	}

	return
}

func (p *RolePermissions) AddAdmin(guildId, roleId uint64) (err error) {
	query := `INSERT INTO role_permissions("guild_id", "role_id", "support", "admin") VALUES($1, $2, true, true) ON CONFLICT("role_id") DO UPDATE SET "admin" = true, "support" = true;`
	_, err = p.Exec(context.Background(), query, guildId, roleId)
	return
}

func (p *RolePermissions) AddSupport(guildId, roleId uint64) (err error) {
	query := `INSERT INTO role_permissions("guild_id", "role_id", "support", "admin") VALUES($1, $2, true, false) ON CONFLICT("role_id") DO UPDATE SET "admin" = false, "support" = true;`
	_, err = p.Exec(context.Background(), query, guildId, roleId)
	return
}

func (p *RolePermissions) RemoveAdmin(guildId, roleId uint64) (err error) {
	query := `UPDATE role_permissions SET "admin" = false WHERE "guild_id" = $1 AND "role_id" = $2;`
	_, err = p.Exec(context.Background(), query, guildId, roleId)
	return
}

func (p *RolePermissions) RemoveSupport(guildId, roleId uint64) (err error) {
	query := `UPDATE role_permissions SET "admin" = false, "support" = false WHERE "guild_id" = $1 AND "role_id" = $2;`
	_, err = p.Exec(context.Background(), query, guildId, roleId)
	return
}

