package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Form struct {
	Id       int    `json:"form_id"`
	GuildId  uint64 `json:"guild_id,string"`
	Title    string `json:"title"`
	CustomId string `json:"custom_id"`
}

type FormsTable struct {
	*pgxpool.Pool
}

func newFormsTable(db *pgxpool.Pool) *FormsTable {
	return &FormsTable{
		db,
	}
}

func (f FormsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS forms(
	"form_id" SERIAL NOT NULL UNIQUE,
	"guild_id" int8 NOT NULL,
	"title" VARCHAR(255) NOT NULL,
    "custom_id" VARCHAR(100) UNIQUE NOT NULL,
	PRIMARY KEY("form_id")
);
CREATE INDEX IF NOT EXISTS forms_guild_id ON forms("guild_id");
`
}

func (f *FormsTable) Get(ctx context.Context, formId int) (form Form, ok bool, e error) {
	query := `SELECT "form_id", "guild_id", "title", "custom_id" FROM forms WHERE "form_id" = $1;`

	err := f.QueryRow(ctx, query, formId).Scan(&form.Id, &form.GuildId, &form.Title, &form.CustomId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Form{}, false, nil
		} else {
			return Form{}, false, err
		}
	}

	return form, true, nil
}

func (f *FormsTable) GetForms(ctx context.Context, guildId uint64) (forms []Form, e error) {
	query := `SELECT "form_id", "guild_id", "title", "custom_id" FROM forms WHERE "guild_id" = $1;`

	rows, err := f.Query(ctx, query, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var form Form
		if err := rows.Scan(&form.Id, &form.GuildId, &form.Title, &form.CustomId); err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return
}

func (f *FormsTable) Create(ctx context.Context, guildId uint64, title, customId string) (int, error) {
	query := `
INSERT INTO forms("guild_id", "title", "custom_id")
VALUES($1, $2, $3)
RETURNING "form_id";
`

	var id int
	if err := f.QueryRow(ctx, query, guildId, title, customId).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (f *FormsTable) UpdateTitle(ctx context.Context, formId int, title string) (err error) {
	query := `UPDATE forms SET "title" = $1 WHERE "form_id" = $2;`
	_, err = f.Exec(ctx, query, title, formId)
	return
}

func (f *FormsTable) Delete(ctx context.Context, formId int) (err error) {
	query := `DELETE FROM forms WHERE "form_id" = $1;`
	_, err = f.Exec(ctx, query, formId)
	return
}
