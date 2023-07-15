CREATE TABLE IF NOT EXISTS archive_messages (
    guild_id int8 NOT NULL,
    ticket_id int4 NOT NULL,
    channel_id int8 NOT NULL,
    message_id int8 NOT NULL,
    FOREIGN KEY (guild_id, ticket_id) REFERENCES tickets(guild_id, id) ON DELETE CASCADE,
    PRIMARY KEY (guild_id, ticket_id)
);