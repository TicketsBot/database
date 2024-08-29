CREATE TABLE IF NOT EXISTS patreon_entitlements (
    "entitlement_id" UUID NOT NULL,
    "user_id" int8 NOT NULL,
    PRIMARY KEY ("entitlement_id"),
    UNIQUE ("entitlement_id", "user_id"), -- For use in ON CONFLICT
    UNIQUE ("user_id"),
    FOREIGN KEY ("user_id") REFERENCES legacy_premium_entitlements ("user_id"),
    FOREIGN KEY ("entitlement_id") REFERENCES entitlements ("id")
);