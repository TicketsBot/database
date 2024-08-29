INSERT INTO discord_entitlements(discord_id, entitlement_id)
VALUES ($1, $2)
ON CONFLICT ("discord_id") DO UPDATE SET "entitlement_id" = $2;