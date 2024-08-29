SELECT EXISTS (
    (
        SELECT 1
        FROM whitelabel_allowed_guilds
        WHERE "bot_id" = $1 AND "guild_id" = $2
     )
    UNION
    (
        SELECT 1
        FROM whitelabel
        INNER JOIN entitlements ON entitlements.id = whitelabel.entitlement_id
        INNER JOIN skus ON skus.id = entitlements.sku_id
        INNER JOIN whitelabel_skus ON whitelabel_skus.sku_id = skus.id
        WHERE whitelabel.bot_id = $1 AND whitelabel_skus.servers_per_bot_permitted IS NULL
    )
);