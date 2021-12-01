package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FormInput struct {
	Id          int     `json:"id"`
	FormId      int     `json:"form_id"`
	CustomId    string  `json:"custom_id"`
	Style       uint8   `json:"style"`
	Label       string  `json:"label"`
	Placeholder *string `json:"placeholder,omitempty"`
}

type FormInputTable struct {
	*pgxpool.Pool
}

func newFormInputTable(db *pgxpool.Pool) *FormInputTable {
	return &FormInputTable{
		db,
	}
}

func (f FormInputTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS form_input(
	"id" SERIAL NOT NULL UNIQUE,
	"form_id" int NOT NULL,
    "custom_id" VARCHAR(100) NOT NULL,
    "style" int2 NOT NULL,
    "label" VARCHAR(255) NOT NULL,
    "placeholder" VARCHAR(100) NULL,
	FOREIGN KEY("form_id") REFERENCES forms("form_id"),
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS form_input_form_id ON form_input("form_id");
`
}

func (f *FormInputTable) Get(id int) (input FormInput, ok bool, e error) {
	query := `SELECT "id", "form_id", "custom_id", "style", "label", "placeholder" FROM form_input WHERE "id" = $1;`

	err := f.QueryRow(context.Background(), query, id).Scan(&input.Id, &input.FormId, &input.CustomId, &input.Style, &input.Label, &input.Placeholder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return FormInput{}, false, nil
		} else {
			return FormInput{}, false, err
		}
	}

	return
}

func (f *FormInputTable) GetInputs(formId int) (inputs []FormInput, e error) {
	query := `SELECT "id", "form_id", "custom_id", "style", "label", "placeholder" FROM form_input WHERE "form_id" = $1;`

	rows, err := f.Query(context.Background(), query, formId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var input FormInput
		if err := rows.Scan(&input.Id, &input.FormId, &input.CustomId, &input.Style, &input.Label, &input.Placeholder); err != nil {
			return nil, err
		}

		inputs = append(inputs, input)
	}

	return
}

// Form ID -> Form Input
func (f *FormInputTable) GetInputsForGuild(guildId uint64) (inputs map[int]FormInput, e error) {
	query := `
SELECT form_input.id, form_input.form_id, form_input.custom_id, form_input.style, form_input.label, form_input.placeholder
FROM form_input 
INNER JOIN forms ON form_input.form_id = forms.form_id
WHERE forms.guild_id = $1;
`

	rows, err := f.Query(context.Background(), query, guildId)
	if err != nil {
		return nil, err
	}

	inputs = make(map[int]FormInput)
	for rows.Next() {
		var input FormInput
		if err := rows.Scan(&input.Id, &input.FormId, &input.CustomId, &input.Style, &input.Label, &input.Placeholder); err != nil {
			return nil, err
		}

		inputs[input.FormId] = input
	}

	return
}

func (f *FormInputTable) Update(input FormInput) (err error) {
	query := `
UPDATE form_input
SET "style" = $2,
	"label"= $3,
    "placeholder" = $4
WHERE "id" = $1;
`

	_, err = f.Exec(context.Background(), query, input.Id, input.Style, input.Label, input.Placeholder)
	return
}

func (f *FormInputTable) Delete(formInputId, formId int) (err error) {
	query := `DELETE FROM form_input WHERE "id" = $1 AND "form_id" = $2;`
	_, err = f.Exec(context.Background(), query, formInputId, formId)
	return
}
