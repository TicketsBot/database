package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NamingScheme string

const (
	Id       NamingScheme = "id"
	Username NamingScheme = "username"
)

type TicketNamingScheme struct {
	*pgxpool.Pool
}

func newTicketNamingScheme(db *pgxpool.Pool) *TicketNamingScheme {
	return &TicketNamingScheme{
		db,
	}
}

func (t TicketNamingScheme) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS naming_scheme(
	"guild_id" int8 NOT NULL UNIQUE,
	"naming_scheme" varchar(16) NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (t *TicketNamingScheme) Get(ctx context.Context, guildId uint64) (ns NamingScheme, e error) {
	query := `SELECT "naming_scheme" from naming_scheme WHERE "guild_id" = $1`

	var namingScheme string
	if err := t.QueryRow(ctx, query, guildId).Scan(&namingScheme); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	if namingScheme == "" {
		ns = Id
	} else {
		ns = NamingScheme(namingScheme)
	}

	return
}

func (t *TicketNamingScheme) Set(ctx context.Context, guildId uint64, scheme NamingScheme) (err error) {
	query := `INSERT INTO naming_scheme("guild_id", "naming_scheme") VALUES($1, $2) ON CONFLICT("guild_id") DO UPDATE SET "naming_scheme" = $2;`
	_, err = t.Exec(ctx, query, guildId, scheme)
	return
}
