CREATE TABLE IF NOT EXISTS exit_survey_responses(
    "guild_id" int8 NOT NULL,
    "ticket_id" int4 NOT NULL,
    "form_id" int4,
    "question_id" int4,
    "response" TEXT,
    FOREIGN KEY ("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
    FOREIGN KEY ("form_id") REFERENCES forms("form_id") ON DELETE CASCADE,
    FOREIGN KEY ("question_id") REFERENCES form_input("id") ON DELETE CASCADE,
    PRIMARY KEY ("guild_id", "ticket_id", "question_id")
);

CREATE INDEX IF NOT EXISTS exit_survey_responses_guild_id ON exit_survey_responses("guild_id");
CREATE INDEX IF NOT EXISTS exit_survey_responses_form_id ON exit_survey_responses("form_id");