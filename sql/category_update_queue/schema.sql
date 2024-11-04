CREATE TYPE ticket_status AS ENUM ('OPEN', 'PENDING', 'CLOSED');

CREATE TABLE IF NOT EXISTS category_update_queue (
    guild_id INT8 NOT NULL,
    ticket_id INT8 NOT NULL,
    new_status ticket_status NOT NULL,
    status_changed_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (guild_id, ticket_id),
    FOREIGN KEY (guild_id, ticket_id) REFERENCES tickets(guild_id, id) ON DELETE CASCADE
);
