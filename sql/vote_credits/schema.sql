CREATE TABLE IF NOT EXISTS vote_credits
(
    user_id int8 NOT NULL,
    credits int4 NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id)
);
