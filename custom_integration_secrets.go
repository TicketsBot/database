package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationSecretsTable struct {
	*pgxpool.Pool
}

type CustomIntegrationSecret struct {
	Id            int     `json:"id"`
	IntegrationId int     `json:"integration_id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
}

func newCustomIntegrationSecretsTable(db *pgxpool.Pool) *CustomIntegrationSecretsTable {
	return &CustomIntegrationSecretsTable{
		db,
	}
}

func (i CustomIntegrationSecretsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS custom_integration_secrets(
	"id" SERIAL NOT NULL UNIQUE,
	"integration_id" int NOT NULL,
	"name" VARCHAR(32) NOT NULL,
	"description" VARCHAR(255) NULL,
	UNIQUE ("integration_id", "name"),
	FOREIGN KEY("integration_id") REFERENCES custom_integrations("id") ON DELETE CASCADE,
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS custom_integration_secrets_integration_id ON custom_integration_secrets("integration_id");
`
}

func (i *CustomIntegrationSecretsTable) GetByIntegration(integrationId int) ([]CustomIntegrationSecret, error) {
	query := `SELECT "id", "integration_id", "name", "description" FROM custom_integration_secrets WHERE "integration_id" = $1;`

	rows, err := i.Query(context.Background(), query, integrationId)
	if err != nil {
		return nil, err
	}

	var secrets []CustomIntegrationSecret
	for rows.Next() {
		var secret CustomIntegrationSecret
		if err := rows.Scan(&secret.Id, &secret.IntegrationId, &secret.Name, &secret.Description); err != nil {
			return nil, err
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

/// Assume that secrets[].Id is valid for the guild and integration
func (i *CustomIntegrationSecretsTable) CreateOrUpdate(integrationId int, secrets []CustomIntegrationSecret) ([]CustomIntegrationSecret, error) {
	if len(secrets) == 0 {
		query := `DELETE FROM custom_integration_secrets WHERE "integration_id" = $1;`
		_, err := i.Exec(context.Background(), query, integrationId)
		return nil, err
	}

	tx, err := i.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(context.Background()) // Does not matter if commit succeeds

	// Delete existing secrets that are not in the new list
	query := `DELETE FROM custom_integration_secrets WHERE "integration_id" = $1 AND NOT("id" = ANY($2));`

	var ids []int
	for _, secret := range secrets {
		if secret.Id != 0 {
			ids = append(ids, secret.Id)
		}
	}

	array := &pgtype.Int4Array{}
	if err := array.Set(ids); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(context.Background(), query, integrationId, array); err != nil {
		return nil, err
	}

	// Create or update new secrets
	var newSecrets []CustomIntegrationSecret
	for _, secret := range secrets {
		var res CustomIntegrationSecret
		if secret.Id == 0 { // Create
			query := `
INSERT INTO custom_integration_secrets("integration_id", "name", "description")
VALUES ($1, $2, $3)
RETURNING "id", "integration_id", "name", "description";
;`

			err = tx.QueryRow(context.Background(), query, integrationId, secret.Name, secret.Description).Scan(&res.Id, &res.IntegrationId, &res.Name, &res.Description)
		} else { // Update
			query := `
UPDATE custom_integration_secrets
SET "name" = $3, "description" = $4
WHERE "id" = $1 AND "integration_id" = $2;`

			_, err = tx.Exec(context.Background(), query, secret.Id, integrationId, secret.Name, secret.Description)
			res = CustomIntegrationSecret{
				Id:            secret.Id,
				IntegrationId: integrationId,
				Name:          secret.Name,
				Description:   secret.Description,
			}
		}

		if err != nil {
			return nil, err
		}

		newSecrets = append(newSecrets, res)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return newSecrets, nil
}

func (i *CustomIntegrationSecretsTable) Delete(id int) (err error) {
	query := `
DELETE FROM custom_integration_secrets
WHERE "id" = $1;
`

	_, err = i.Exec(context.Background(), query, id)
	return
}
