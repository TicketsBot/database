CREATE TABLE IF NOT EXISTS legacy_premium_entitlements
(
    "user_id"    int8         NOT NULL UNIQUE,
    "tier"       int4         NOT NULL,
    "sku_label"  VARCHAR(255) NOT NULL,
    "expires_at" timestamp    NOT NULL,
    PRIMARY KEY ("user_id")
);