CREATE TABLE IF NOT EXISTS multi_server_skus
(
    sku_id            UUID NOT NULL,
    servers_permitted INT  NOT NULL,
    PRIMARY KEY (sku_id),
    FOREIGN KEY (sku_id) REFERENCES skus (id)
);