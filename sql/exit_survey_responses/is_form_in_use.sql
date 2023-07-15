SELECT EXISTS (
    SELECT 1
    FROM exit_survey_responses
    WHERE guild_id = $1 AND form_id = $2
);