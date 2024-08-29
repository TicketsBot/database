WITH result AS (SELECT *
                FROM whitelabel_allowed_guilds
                WHERE "bot_id" = $1
                    FOR UPDATE)
SELECT COUNT(*)
FROM result;