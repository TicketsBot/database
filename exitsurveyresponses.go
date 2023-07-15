package database

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ExitSurveyResponse struct {
	GuildId   uint64             `json:"guild_id,string"`
	TicketId  int                `json:"ticket_id"`
	Responses []QuestionResponse `json:"responses"`
}

type QuestionResponse struct {
	QuestionId *int    `json:"question_id"`
	Question   *string `json:"question"`
	Response   string  `json:"response"`
}

type ExitSurveyResponses struct {
	*pgxpool.Pool
}

func newExitSurveyResponses(db *pgxpool.Pool) *ExitSurveyResponses {
	return &ExitSurveyResponses{
		db,
	}
}

var (
	//go:embed sql/exit_survey_responses/schema.sql
	exitSurveyResponsesSchema string

	//go:embed sql/exit_survey_responses/add_responses.sql
	exitSurveyResponsesAdd string

	//go:embed sql/exit_survey_responses/get_response_single.sql
	exitSurveyResponsesGetSingle string

	//go:embed sql/exit_survey_responses/is_form_in_use.sql
	exitSurveyResponsesIsInUse string

	//go:embed sql/exit_survey_responses/has_response.sql
	exitSurveyHasResponse string
)

func (e *ExitSurveyResponses) Schema() string {
	return exitSurveyResponsesSchema
}

func (e *ExitSurveyResponses) AddResponses(guildId uint64, ticketId int, formId int, responses map[int]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTransactionTimeout)
	defer cancel()

	tx, err := e.Begin(ctx)
	if err != nil {
		return err
	}

	for questionId, response := range responses {
		_, err = tx.Exec(ctx, exitSurveyResponsesAdd, guildId, ticketId, formId, questionId, response)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (e *ExitSurveyResponses) GetResponses(guildId uint64, ticketId int) (ExitSurveyResponse, error) {
	rows, err := e.Query(context.Background(), exitSurveyResponsesGetSingle, guildId, ticketId)
	if err != nil {
		return ExitSurveyResponse{}, err
	}

	var responses []QuestionResponse
	for rows.Next() {
		var response QuestionResponse

		if err := rows.Scan(&response.QuestionId, &response.Question, &response.Response); err != nil {
			return ExitSurveyResponse{}, err
		}

		responses = append(responses, response)
	}

	return ExitSurveyResponse{
		GuildId:   guildId,
		TicketId:  ticketId,
		Responses: responses,
	}, nil
}

func (e *ExitSurveyResponses) IsFormInUse(guildId uint64, formId int) (bool, error) {
	var inUse bool
	err := e.QueryRow(context.Background(), exitSurveyResponsesIsInUse, guildId, formId).Scan(&inUse)
	return inUse, err
}

func (e *ExitSurveyResponses) HasResponse(guildId uint64, formId int) (bool, error) {
	var hasResponse bool
	err := e.QueryRow(context.Background(), exitSurveyHasResponse, guildId, formId).Scan(&hasResponse)
	return hasResponse, err
}
