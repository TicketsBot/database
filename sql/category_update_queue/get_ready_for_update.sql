WITH cte AS (
    DELETE FROM category_update_queue
        WHERE status_changed_at < NOW() - $1::INTERVAL
        RETURNING guild_id, ticket_id, new_status
)
SELECT cte.guild_id, cte.ticket_id, cte.new_status, tickets.channel_id, tickets.panel_id
FROM cte
INNER JOIN tickets ON cte.guild_id = tickets.guild_id AND cte.ticket_id = tickets.ticket_id;