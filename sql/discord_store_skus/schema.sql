CREATE TABLE IF NOT EXISTS discord_store_skus
(
    discord_id int8 NOT NULL,
    sku_id     UUID NOT NULL,
    PRIMARY KEY (discord_id),
    FOREIGN KEY (sku_id) REFERENCES skus (id)
);