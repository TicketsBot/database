SELECT
    exit_survey_responses.question_id,
    form_input.label AS question,
    exit_survey_responses.response
FROM exit_survey_responses
INNER JOIN tickets ON exit_survey_responses.guild_id = tickets.guild_id AND exit_survey_responses.ticket_id = tickets.id
INNER JOIN panels ON tickets.panel_id = panels.panel_id
INNER JOIN form_input on exit_survey_responses.question_id = form_input.id
WHERE exit_survey_responses.guild_id = $1 AND exit_survey_responses.ticket_id = $2 AND exit_survey_responses.form_id = panels.exit_survey_form_id;