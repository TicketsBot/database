DELETE FROM dashboard_users
WHERE "last_seen" < NOW() - $1::INTERVAL;