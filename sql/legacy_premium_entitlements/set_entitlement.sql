INSERT INTO legacy_premium_entitlements ("user_id", "tier", "sku_label", "is_legacy", "expires_at")
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ("user_id") DO UPDATE SET "tier"           = $2,
                                      "sku_label"      = $3,
                                      "is_legacy"      = $4,
                                      "expires_at"     = $5;