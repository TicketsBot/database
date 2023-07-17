INSERT INTO dashboard_users (user_id, last_seen)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET last_seen = excluded.last_seen;