WITH tiers AS (
    SELECT subscription_skus.tier, subscription_skus.priority
    FROM entitlements
    INNER JOIN skus ON entitlements.sku_id = skus.id
    INNER JOIN subscription_skus ON skus.id = subscription_skus.sku_id
    WHERE (
            entitlements.expires_at IS NULL OR
            entitlements.expires_at > (NOW() - $3::interval)
          ) AND
          entitlements.guild_id = $1 AND
          (entitlements.source != 'voting' OR $4 = true)

    UNION ALL

    SELECT subscription_skus.tier, subscription_skus.priority
    FROM entitlements
    INNER JOIN skus ON entitlements.sku_id = skus.id
    INNER JOIN subscription_skus ON skus.id = subscription_skus.sku_id
    LEFT OUTER JOIN permissions ON permissions.user_id = entitlements.user_id AND permissions.guild_id = $1
    WHERE (
            entitlements.expires_at IS NULL OR
            entitlements.expires_at > (NOW() - $3::interval)
        ) AND
        entitlements.guild_id IS NULL AND
        entitlements.user_id IS NOT NULL AND
        subscription_skus.is_global = true AND
        (entitlements.source != 'voting' OR $4 = true) AND
        (
            entitlements.user_id = $2
                OR
            (entitlements.user_id = permissions.user_id AND permissions.admin = 't' AND permissions.guild_id = $1)
        )
), sorted AS (
    SELECT tier FROM tiers
    ORDER BY priority DESC
)
SELECT DISTINCT tier FROM sorted;