DELETE FROM legacy_premium_entitlements
WHERE "user_id" = $1 AND "sku_label" = $2;