INSERT INTO legacy_premium_entitlements ("user_id", "tier", "sku_label", "expires_at")
VALUES ($1, $2, $3, $4)
ON CONFLICT ("user_id") DO UPDATE SET "tier"       = $2,
                                      "sku_label"  = $3,
                                      "expires_at" = $4;