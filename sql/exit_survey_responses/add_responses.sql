INSERT INTO exit_survey_responses (guild_id, ticket_id, form_id, question_id, response)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (guild_id, ticket_id, question_id)
DO UPDATE SET response = $5;