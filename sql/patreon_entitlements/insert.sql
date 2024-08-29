INSERT INTO patreon_entitlements ("entitlement_id", "user_id")
VALUES($1, $2)
ON CONFLICT ("entitlement_id", "user_id") DO NOTHING;