package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type EmbedField struct {
	FieldId int    `json:"-"`
	EmbedId int    `json:"-"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Inline  bool   `json:"inline"`
}

type EmbedFieldsTable struct {
	*pgxpool.Pool
}

func newEmbedFieldsTable(db *pgxpool.Pool) *EmbedFieldsTable {
	return &EmbedFieldsTable{
		db,
	}
}

func (s EmbedFieldsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS embed_fields(
	"id" SERIAL NOT NULL UNIQUE,
	"embed_id" int NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"value" TEXT NOT NULL CONSTRAINT value_length CHECK (length(value) <= 1024),
	"inline" BOOL NOT NULL,
	FOREIGN KEY("embed_id") REFERENCES embeds("id") ON DELETE CASCADE,
	PRIMARY KEY("id")
);
`
}

func (s *EmbedFieldsTable) GetField(id int) (field EmbedField, err error) {
	query := `
SELECT 
	"id",
	"embed_id",
	"name",
	"value",
	"inline"
FROM embed_fields
WHERE "id" = $1;
`

	err = s.QueryRow(context.Background(), query, id).Scan(
		&field.FieldId,
		&field.EmbedId,
		&field.Name,
		&field.Value,
		&field.Inline,
	)

	return
}

func (s *EmbedFieldsTable) GetFieldsForEmbed(embedId int) ([]EmbedField, error) {
	query := `
SELECT 
	"id",
	"embed_id",
	"name",
	"value",
	"inline"
FROM embed_fields
WHERE "embed_id" = $1;
`

	rows, err := s.Query(context.Background(), query, embedId)
	if err != nil {
		return nil, err
	}

	var fields []EmbedField
	for rows.Next() {
		var field EmbedField
		if err := rows.Scan(&field.FieldId, &field.EmbedId, &field.Name, &field.Value, &field.Inline); err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// GetAllFieldsForPanels Returns a map of [embed_id][]EmbedField
func (s *EmbedFieldsTable) GetAllFieldsForPanels(guildId uint64) (map[int][]EmbedField, error) {
	query := `
SELECT 
	embed_fields.id,
	embed_fields.embed_id,
	embed_fields.name,
	embed_fields.value,
	embed_fields.inline
FROM embed_fields
INNER JOIN embeds
ON embeds.id = embed_fields.embed_id
INNER JOIN panels
ON panels.welcome_message = embeds.id
WHERE embeds.guild_id = $1
ORDER BY embed_fields.embed_id, embed_fields.id;
`

	rows, err := s.Query(context.Background(), query, guildId)
	if err != nil {
		return nil, err
	}

	fields := make(map[int][]EmbedField)
	for rows.Next() {
		var field EmbedField
		if err := rows.Scan(&field.FieldId, &field.EmbedId, &field.Name, &field.Value, &field.Inline); err != nil {
			return nil, err
		}

		slice, ok := fields[field.EmbedId]
		if !ok {
			slice = make([]EmbedField, 0)
		}

		slice = append(slice, field)
		fields[field.EmbedId] = slice
	}

	return fields, nil
}
