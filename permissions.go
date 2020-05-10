package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Permissions struct {
	*pgxpool.Pool
}

func newPermissions(db *pgxpool.Pool) *Permissions {
	return &Permissions{
		db,
	}
}

func (p Permissions) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS permissions("guild_id" int8 NOT NULL, "user_id" int8 NOT NULL, "support" bool NOT NULL, "admin" bool NOT NULL, PRIMARY KEY("guild_id", "user_id"));
CREATE INDEX IF NOT EXISTS permissions_guild_id ON permissions("guild_id");
`
}

func (p *Permissions) IsSupport(guildId, userId uint64) (support bool, e error) {
	var admin bool

	if err := p.QueryRow(context.Background(), `SELECT "support", "admin" from permissions WHERE "guild_id" = $1 AND "user_id" = $2;`, guildId, userId).Scan(&support, &admin); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	if admin {
		support = true
	}

	return
}

func (p *Permissions) IsAdmin(guildId, userId uint64) (admin bool, e error) {
	if err := p.QueryRow(context.Background(), `SELECT "admin" from permissions WHERE "guild_id" = $1 AND "user_id" = $2;`, guildId, userId).Scan(&admin); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *Permissions) GetAdmins(guildId uint64) (admins []uint64, e error) {
	rows, err := p.Query(context.Background(), `SELECT "user_id" from permissions WHERE "guild_id" = $1 AND "admin" = true;`, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
	}

	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			e = err
			continue
		}

		admins = append(admins, userId)
	}

	return
}

func (p *Permissions) GetSupport(guildId uint64) (support []uint64, e error) {
	rows, err := p.Query(context.Background(), `SELECT "user_id" from permissions WHERE "guild_id" = $1 AND ("admin" = true OR "support" = true);`, guildId)
	defer rows.Close()
	if err != nil && err != pgx.ErrNoRows {
		e = err
	}

	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			e = err
			continue
		}

		support = append(support, userId)
	}

	return
}

func (p *Permissions) AddAdmin(guildId, userId uint64) (err error) {
	query := `INSERT INTO permissions("guild_id", "user_id", "support", "admin") VALUES($1, $2, true, true) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "admin" = true, "support" = true;`
	_, err = p.Exec(context.Background(), query, guildId, userId)
	return
}

func (p *Permissions) AddSupport(guildId, userId uint64) (err error) {
	query := `INSERT INTO permissions("guild_id", "user_id", "support", "admin") VALUES($1, $2, true, false) ON CONFLICT("guild_id", "user_id") DO UPDATE SET "admin" = false, "support" = true;`
	_, err = p.Exec(context.Background(), query, guildId, userId)
	return
}

func (p *Permissions) RemoveAdmin(guildId, userId uint64) (err error) {
	query := `UPDATE permissions SET "admin" = false WHERE "guild_id" = $1 AND "user_id" = $2;`
	_, err = p.Exec(context.Background(), query, guildId, userId)
	return
}

func (p *Permissions) RemoveSupport(guildId, userId uint64) (err error) {
	query := `UPDATE permissions SET "admin" = false, "support" = false WHERE "guild_id" = $1 AND "user_id" = $2;`
	_, err = p.Exec(context.Background(), query, guildId, userId)
	return
}

