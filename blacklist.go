package database

import "context"

type Blacklist struct {
	GuildId uint64
	UserId  uint64
}

func (b Blacklist) Schema() string {
	return `CREATE TABLE IF NOT EXISTS blacklist("guild_id" int8 NOT NULL, "user_id" int8 NOT NULL, PRIMARY KEY("guild_id", "user_id");`
}

func (b *Blacklist) IsBlacklisted(db *Database) bool {
	var exists bool
	if err := db.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM blacklist WHERE "guild_id"=$1 AND "user_id"=$2);`, b.GuildId, b.UserId).Scan(&exists); err != nil {
		db.Logger.Error(err)
	}

	return exists
}

func (b *Blacklist) Add(db *Database) {
	// on conflict, user is already blacklisted
	if _, err := db.Exec(context.Background(), `INSERT INTO blacklist("guild_id", "user_id") VALUES($1, $2) ON CONFLICT DO NOTHING;`, b.GuildId, b.UserId); err != nil {
		db.Logger.Error(err)
	}
}

func (b *Blacklist) Remove(db *Database) {
	if _, err := db.Exec(context.Background(), `DELETE FROM blacklist WHERE "guild_id"=$1 AND "user_id"=$2;`, b.GuildId, b.UserId); err != nil {
		db.Logger.Error(err)
	}
}
