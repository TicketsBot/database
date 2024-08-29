CREATE TYPE premium_source AS ENUM ('discord', 'patreon', 'voting', 'key');

CREATE TABLE IF NOT EXISTS entitlements
(
    id         UUID DEFAULT gen_random_uuid(),
    guild_id   int8 DEFAULT NULL,
    user_id    int8,
    sku_id     UUID           NOT NULL,
    source     premium_source NOT NULL,
    expires_at timestamptz,
    PRIMARY KEY (id),
    UNIQUE NULLS NOT DISTINCT (guild_id, user_id, sku_id, source),
    FOREIGN KEY (sku_id) REFERENCES skus (id)
);
