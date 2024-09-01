DELETE FROM
legacy_premium_entitlement_guilds
WHERE "user_id" = $1 AND "guild_id" = $2;