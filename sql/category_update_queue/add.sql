INSERT INTO category_update_queue (guild_id, ticket_id, new_status, status_changed_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (guild_id, ticket_id) DO UPDATE
SET new_status = EXCLUDED.new_status, status_changed_at = EXCLUDED.status_changed_at;