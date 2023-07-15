SELECT EXISTS (
    SELECT 1
    FROM exit_survey_responses
    WHERE guild_id = $1 AND ticket_id = $2
)