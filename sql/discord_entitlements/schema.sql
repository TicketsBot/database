CREATE TABLE IF NOT EXISTS discord_entitlements (
    discord_id int8 NOT NULL,
    entitlement_id UUID NOT NULL,
    PRIMARY KEY (discord_id),
    FOREIGN KEY (entitlement_id) REFERENCES entitlements(id) ON DELETE CASCADE
);