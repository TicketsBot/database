SELECT "id", "guild_id", "user_id", "sku_id", "source", "expires_at"
FROM entitlements
WHERE "source" = $1;