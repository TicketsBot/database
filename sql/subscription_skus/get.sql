SELECT skus.id, skus.label, skus.type, subscription_skus.tier, subscription_skus.priority, subscription_skus.is_global
FROM skus
INNER JOIN subscription_skus ON skus.id = subscription_skus.sku_id
WHERE skus.id = $1;