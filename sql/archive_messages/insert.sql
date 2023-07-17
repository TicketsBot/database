INSERT INTO archive_messages (guild_id, ticket_id, channel_id, message_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (guild_id, ticket_id) DO UPDATE SET
    channel_id = excluded.channel_id,
    message_id = excluded.message_id;