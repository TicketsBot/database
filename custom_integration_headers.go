package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationHeadersTable struct {
	*pgxpool.Pool
}

type CustomIntegrationHeader struct {
	Id            int    `json:"id"`
	IntegrationId int    `json:"integration_id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
}

func newCustomIntegrationHeadersTable(db *pgxpool.Pool) *CustomIntegrationHeadersTable {
	return &CustomIntegrationHeadersTable{
		db,
	}
}

func (i CustomIntegrationHeadersTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS custom_integration_headers(
	"id" SERIAL NOT NULL UNIQUE,
	"integration_id" int NOT NULL,
	"name" VARCHAR(32) NOT NULL,
	"value" VARCHAR(255) NOT NULL,
	UNIQUE("integration_id", "name"),
	FOREIGN KEY("integration_id") REFERENCES custom_integrations("id") ON DELETE CASCADE,
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS custom_integration_headers_integration_id ON custom_integration_headers("integration_id");
`
}

func (i *CustomIntegrationHeadersTable) GetByIntegration(ctx context.Context, integrationId int) ([]CustomIntegrationHeader, error) {
	query := `SELECT "id", "integration_id", "name", "value" FROM custom_integration_headers WHERE "integration_id" = $1;`

	rows, err := i.Query(ctx, query, integrationId)
	if err != nil {
		return nil, err
	}

	var headers []CustomIntegrationHeader
	for rows.Next() {
		var header CustomIntegrationHeader
		if err := rows.Scan(&header.Id, &header.IntegrationId, &header.Name, &header.Value); err != nil {
			return nil, err
		}
		headers = append(headers, header)
	}

	return headers, nil
}

// GetAll integration_id -> []CustomIntegrationHeader
func (i *CustomIntegrationHeadersTable) GetAll(ctx context.Context, integrationIds []int) (map[int][]CustomIntegrationHeader, error) {
	query := `SELECT "id", "integration_id", "name", "value" FROM custom_integration_headers WHERE "integration_id" = ANY($1);`

	idArray := &pgtype.Int4Array{}
	if err := idArray.Set(integrationIds); err != nil {
		return nil, err
	}

	rows, err := i.Query(ctx, query, idArray)
	if err != nil {
		return nil, err
	}

	headers := make(map[int][]CustomIntegrationHeader)
	for rows.Next() {
		var header CustomIntegrationHeader
		if err := rows.Scan(&header.Id, &header.IntegrationId, &header.Name, &header.Value); err != nil {
			return nil, err
		}

		if _, ok := headers[header.IntegrationId]; !ok {
			headers[header.IntegrationId] = []CustomIntegrationHeader{}
		}

		headers[header.IntegrationId] = append(headers[header.IntegrationId], header)
	}

	return headers, nil
}

// Assumes that all header IDs are valid for the integration
func (i *CustomIntegrationHeadersTable) CreateOrUpdate(ctx context.Context, integrationId int, headers []CustomIntegrationHeader) ([]CustomIntegrationHeader, error) {
	// The array check is weird if headers is empty
	if len(headers) == 0 {
		query := `DELETE FROM custom_integration_headers WHERE "integration_id" = $1;`
		_, err := i.Exec(ctx, query, integrationId)
		return nil, err
	}

	tx, err := i.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // Does not matter if commit succeeds

	query := `DELETE FROM custom_integration_headers WHERE "integration_id" = $1 AND NOT ("id" = ANY($2));`

	var ids []int
	for _, header := range headers {
		if header.Id != 0 {
			ids = append(ids, header.Id)
		}
	}

	array := &pgtype.Int4Array{}
	if err := array.Set(ids); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, query, integrationId, array); err != nil {
		return nil, err
	}

	// Create or update new secrets
	var newHeaders []CustomIntegrationHeader
	for _, header := range headers {
		var res CustomIntegrationHeader
		if header.Id == 0 { // Create
			query := `
INSERT INTO custom_integration_headers( "integration_id", "name", "value")
VALUES ($1, $2, $3)
RETURNING "id", "integration_id", "name", "value";
;`

			err = tx.QueryRow(ctx, query, integrationId, header.Name, header.Value).Scan(
				&res.Id,
				&res.IntegrationId,
				&res.Name,
				&res.Value,
			)
		} else { // Update
			query := `
UPDATE custom_integration_headers
SET "name" = $3, "value" = $4
WHERE "id" = $1 AND "integration_id" = $2;`

			_, err = tx.Exec(ctx, query, header.Id, integrationId, header.Name, header.Value)
			res = CustomIntegrationHeader{
				Id:            header.Id,
				IntegrationId: integrationId,
				Name:          header.Name,
				Value:         header.Value,
			}
		}

		if err != nil {
			return nil, err
		}

		newHeaders = append(newHeaders, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return newHeaders, nil
}

func (i *CustomIntegrationHeadersTable) Delete(ctx context.Context, id int) (err error) {
	query := `
DELETE FROM custom_integration_headers
WHERE "id" = $1;
`

	_, err = i.Exec(ctx, query, id)
	return
}
