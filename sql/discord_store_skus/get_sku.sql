SELECT skus.id, skus.label, skus.type
FROM discord_store_skus
INNER JOIN skus ON skus.id = discord_store_skus.sku_id
WHERE discord_store_skus.discord_id = $1;