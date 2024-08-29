SELECT skus.id, skus.label, skus.type, subscription_skus.tier, subscription_skus.priority, subscription_skus.is_global
FROM skus
INNER JOIN subscription_skus ON subscription_skus.sku_id = skus.id
WHERE skus.label ILIKE '%' || $1 || '%'
LIMIT $2;