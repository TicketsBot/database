CREATE TABLE IF NOT EXISTS whitelabel_allowed_guilds (
    bot_id int8 NOT NULL,
    guild_id int8 NOT NULL,
    PRIMARY KEY (bot_id, guild_id),
    FOREIGN KEY (bot_id) REFERENCES whitelabel (bot_id) ON DELETE CASCADE
);