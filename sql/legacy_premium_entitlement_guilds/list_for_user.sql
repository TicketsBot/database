SELECT "user_id", "guild_id", "entitlement_id"
FROM legacy_premium_entitlement_guilds
WHERE "user_id" = $1
FOR UPDATE;