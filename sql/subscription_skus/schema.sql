CREATE TYPE premium_tier AS ENUM ('premium', 'whitelabel');

CREATE TABLE IF NOT EXISTS subscription_skus
(
    sku_id    UUID         NOT NULL,
    tier      premium_tier NOT NULL,
    priority  INT          NOT NULL,
    is_global BOOLEAN      NOT NULL DEFAULT FALSE,
    PRIMARY KEY (sku_id),
    FOREIGN KEY (sku_id) REFERENCES skus (id)
);
