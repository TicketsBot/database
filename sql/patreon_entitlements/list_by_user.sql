SELECT entitlements.id, entitlements.guild_id, entitlements.user_id, entitlements.sku_id, entitlements.source, entitlements.expires_at
FROM entitlements
INNER JOIN patreon_entitlements ON entitlements.id = patreon_entitlements.entitlement_id
WHERE patreon_entitlements.user_id = $1 AND entitlements.source = 'patreon'
FOR UPDATE;