CREATE TABLE IF NOT EXISTS legacy_premium_entitlement_guilds (
    user_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    entitlement_id UUID NOT NULL UNIQUE,
    PRIMARY KEY (user_id, guild_id),
    FOREIGN KEY (user_id) REFERENCES legacy_premium_entitlements (user_id),
    FOREIGN KEY (entitlement_id) REFERENCES entitlements (id)
);
