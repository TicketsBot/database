SELECT "user_id", "tier", "sku_label", "sku_id", "is_legacy", "expires_at"
FROM legacy_premium_entitlements
WHERE "user_id" = $1 AND "expires_at" > (NOW() - $2::interval);