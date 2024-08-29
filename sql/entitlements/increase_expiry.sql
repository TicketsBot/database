INSERT INTO entitlements(guild_id, user_id, sku_id, source, expires_at)
VALUES($1, $2, $3, $4, NOW() + $5::INTERVAL)
ON CONFLICT (guild_id, user_id, sku_id, source) DO UPDATE SET expires_at = excluded.expires_at + $5::INTERVAL;