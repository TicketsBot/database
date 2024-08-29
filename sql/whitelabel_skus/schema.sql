CREATE TABLE IF NOT EXISTS whitelabel_skus
(
    sku_id                    UUID NOT NULL,
    bots_permitted            int4 NOT NULL,
    servers_per_bot_permitted int4,
    PRIMARY KEY (sku_id),
    FOREIGN KEY (sku_id) REFERENCES skus (id)
);
