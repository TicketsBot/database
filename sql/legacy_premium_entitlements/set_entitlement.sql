INSERT INTO legacy_premium_entitlements ("user_id", "tier", "sku_label", "sku_id", "is_legacy", "expires_at")
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ("user_id") DO UPDATE SET "tier"       = $2,
                                      "sku_label"  = $3,
                                      "sku_id"     = $4,
                                      "is_legacy"  = $5,
                                      "expires_at" = $6;