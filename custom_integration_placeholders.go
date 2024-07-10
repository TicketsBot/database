package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationPlaceholdersTable struct {
	*pgxpool.Pool
}

type CustomIntegrationPlaceholder struct {
	Id            int    `json:"id"`
	IntegrationId int    `json:"integration_id"`
	Name          string `json:"name"`
	JsonPath      string `json:"json_path"`
}

func newCustomIntegrationPlaceholdersTable(db *pgxpool.Pool) *CustomIntegrationPlaceholdersTable {
	return &CustomIntegrationPlaceholdersTable{
		db,
	}
}

func (i CustomIntegrationPlaceholdersTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS custom_integration_placeholders(
	"id" SERIAL NOT NULL UNIQUE,
	"integration_id" int NOT NULL,
	"name" VARCHAR(32) NOT NULL,
	"json_path" VARCHAR(255) NOT NULL,
	UNIQUE("integration_id", "name"),
	FOREIGN KEY("integration_id") REFERENCES custom_integrations("id") ON DELETE CASCADE,
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS custom_integration_placeholders_integration_id ON custom_integration_placeholders("integration_id");
`
}

func (i *CustomIntegrationPlaceholdersTable) GetByIntegration(ctx context.Context, integrationId int) ([]CustomIntegrationPlaceholder, error) {
	query := `SELECT "id", "integration_id", "name", "json_path" FROM custom_integration_placeholders WHERE "integration_id" = $1;`

	rows, err := i.Query(ctx, query, integrationId)
	if err != nil {
		return nil, err
	}

	var placeholders []CustomIntegrationPlaceholder
	for rows.Next() {
		var placeholder CustomIntegrationPlaceholder
		if err := rows.Scan(&placeholder.Id, &placeholder.IntegrationId, &placeholder.Name, &placeholder.JsonPath); err != nil {
			return nil, err
		}
		placeholders = append(placeholders, placeholder)
	}

	return placeholders, nil
}

func (i *CustomIntegrationPlaceholdersTable) GetAllForOwnedIntegrations(ctx context.Context, ownerId uint64) (map[int][]CustomIntegrationPlaceholder, error) {
	query := `
SELECT placeholders.id, placeholders.integration_id, placeholders.name, placeholders.json_path
FROM custom_integration_placeholders AS placeholders 
INNER JOIN custom_integrations ON placeholders.integration_id = custom_integrations.id
WHERE custom_integrations.owner_id = $1;`

	rows, err := i.Query(ctx, query, ownerId)
	if err != nil {
		return nil, err
	}

	placeholders := make(map[int][]CustomIntegrationPlaceholder)
	for rows.Next() {
		var placeholder CustomIntegrationPlaceholder
		if err := rows.Scan(&placeholder.IntegrationId, &placeholder.Id, &placeholder.Name, &placeholder.JsonPath); err != nil {
			return nil, err
		}

		slice, ok := placeholders[placeholder.IntegrationId]
		if !ok {
			slice = make([]CustomIntegrationPlaceholder, 0)
		}

		placeholders[placeholder.IntegrationId] = append(slice, placeholder)
	}

	return placeholders, nil
}

func (i *CustomIntegrationPlaceholdersTable) GetAllActivatedInGuild(ctx context.Context, guildId uint64) ([]CustomIntegrationPlaceholder, error) {
	query := `
SELECT placeholders.id, placeholders.integration_id, placeholders.name, placeholders.json_path
FROM custom_integration_placeholders AS placeholders
INNER JOIN custom_integrations integrations ON placeholders.integration_id = integrations.id
INNER JOIN custom_integration_guilds guilds ON integrations.id = guilds.integration_id
WHERE guilds.guild_id = $1;
`

	rows, err := i.Query(ctx, query, guildId)
	if err != nil {
		return nil, err
	}

	var placeholders []CustomIntegrationPlaceholder
	for rows.Next() {
		var placeholder CustomIntegrationPlaceholder
		if err := rows.Scan(&placeholder.Id, &placeholder.IntegrationId, &placeholder.Name, &placeholder.JsonPath); err != nil {
			return nil, err
		}

		placeholders = append(placeholders, placeholder)
	}

	return placeholders, nil
}

// / Only Name and JsonPath are used
func (i *CustomIntegrationPlaceholdersTable) Set(ctx context.Context, integrationId int, placeholders []CustomIntegrationPlaceholder) ([]CustomIntegrationPlaceholder, error) {
	tx, err := i.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // Does not matter if commit succeeds

	// Delete existing placeholders
	query := `DELETE FROM custom_integration_placeholders WHERE "integration_id" = $1;`
	if _, err := tx.Exec(ctx, query, integrationId); err != nil {
		return nil, err
	}

	var newPlaceholders []CustomIntegrationPlaceholder
	for _, placeholder := range placeholders {
		query := `
INSERT INTO custom_integration_placeholders("integration_id", "name", "json_path")
VALUES ($1, $2, $3)
RETURNING "id", "integration_id", "name", "json_path";
;`

		var res CustomIntegrationPlaceholder
		err := tx.QueryRow(ctx, query, integrationId, placeholder.Name, placeholder.JsonPath).Scan(
			&placeholder.Id,
			&placeholder.IntegrationId,
			&placeholder.Name,
			&placeholder.JsonPath,
		)

		if err != nil {
			return nil, err
		}

		newPlaceholders = append(newPlaceholders, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return newPlaceholders, nil
}

func (i *CustomIntegrationPlaceholdersTable) Delete(ctx context.Context, id int) (err error) {
	query := `
DELETE FROM custom_integration_placeholders
WHERE "id" = $1;
`

	_, err = i.Exec(ctx, query, id)
	return
}
