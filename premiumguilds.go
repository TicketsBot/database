package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PremiumGuilds struct {
	*pgxpool.Pool
}

func newPremiumGuilds(db *pgxpool.Pool) *PremiumGuilds {
	return &PremiumGuilds{
		db,
	}
}

func (p PremiumGuilds) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS premium_guilds(
	"guild_id" int8 NOT NULL UNIQUE,
	"expiry" timestamp NOT NULL,
	PRIMARY KEY("guild_id")
);`
}

func (p *PremiumGuilds) IsPremium(guildId uint64) (bool, error) {
	expiry, err := p.GetExpiry(guildId)
	if err != nil {
		return false, err
	}

	return expiry.After(time.Now()), nil
}

func (p *PremiumGuilds) GetExpiry(guildId uint64) (expiry time.Time, e error) {
	if err := p.QueryRow(context.Background(), `SELECT "expiry" from premium_guilds WHERE "guild_id" = $1;`, guildId).Scan(&expiry); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}

func (p *PremiumGuilds) Add(guildId uint64, interval time.Duration) (err error) {
	_, err = p.Exec(context.Background(), `INSERT INTO premium_guilds("guild_id", "expiry") VALUES($1, NOW() + $2) ON CONFLICT("guild_id") DO UPDATE SET "expiry" = premium_guilds.expiry + $2;`, guildId, interval)
	return
}
