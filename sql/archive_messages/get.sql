SELECT channel_id, message_id
FROM archive_messages
WHERE guild_id = $1 AND ticket_id = $2;