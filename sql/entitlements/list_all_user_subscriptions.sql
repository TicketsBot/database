SELECT entitlements.id, entitlements.user_id, entitlements.source, entitlements.expires_at, skus.id, skus.label, subscription_skus.tier, subscription_skus.priority
FROM entitlements
INNER JOIN skus ON entitlements.sku_id = skus.id
INNER JOIN subscription_skus ON skus.id = subscription_skus.sku_id
WHERE
    entitlements.user_id IS NOT NULL
    AND
    (
        entitlements.expires_at IS NULL OR
        entitlements.expires_at > (NOW() - $1::interval)
    );