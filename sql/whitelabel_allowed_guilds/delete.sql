DELETE FROM whitelabel_allowed_guilds
WHERE "bot_id" = $1 AND "guild_id" = $2;