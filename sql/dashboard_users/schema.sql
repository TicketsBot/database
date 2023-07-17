CREATE TABLE IF NOT EXISTS dashboard_users (
    user_id int8 NOT NULL,
    last_seen timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS dashboard_users_user_id_idx ON dashboard_users(user_id);