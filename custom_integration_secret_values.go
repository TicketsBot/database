package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationSecretValuesTable struct {
	*pgxpool.Pool
}

type SecretWithValue struct {
	CustomIntegrationSecret
	Value string `json:"value"`
}

func newCustomIntegrationSecretValuesTable(db *pgxpool.Pool) *CustomIntegrationSecretValuesTable {
	return &CustomIntegrationSecretValuesTable{
		db,
	}
}

func (i CustomIntegrationSecretValuesTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS custom_integration_secret_values(
	"secret_id" SERIAL NOT NULL UNIQUE,
	"integration_id" int NOT NULL,
    "guild_id" int8 NOT NULL,
	"value" VARCHAR(255) NOT NULL,
    FOREIGN KEY("integration_id") REFERENCES custom_integrations("id") ON DELETE CASCADE,
	FOREIGN KEY("secret_id") REFERENCES custom_integration_secrets("id") ON DELETE CASCADE,
	FOREIGN KEY("integration_id", "guild_id") REFERENCES custom_integration_guilds("integration_id", "guild_id") ON DELETE CASCADE,
	PRIMARY KEY("secret_id", "guild_id")
);
CREATE INDEX IF NOT EXISTS custom_integration_secret_values_integration_id_idx ON custom_integration_secret_values("integration_id");
CREATE INDEX IF NOT EXISTS custom_integration_secret_values_guild_id_idx ON custom_integration_secret_values("guild_id");
`
}

func (i *CustomIntegrationSecretValuesTable) Get(integrationId int, guildId uint64) (map[CustomIntegrationSecret]string, error) {
	query := `
SELECT values.secret_id, values.integration_id, secrets.name, values.value
FROM custom_integration_secret_values AS values 
INNER JOIN custom_integration_secrets AS secrets ON secrets.id = values.secret_id
WHERE values.integration_id = $1 AND values.guild_id = $2;`

	rows, err := i.Query(context.Background(), query, integrationId, guildId)
	if err != nil {
		return nil, err
	}

	data := make(map[CustomIntegrationSecret]string)
	for rows.Next() {
		var secret CustomIntegrationSecret
		var value string
		if err := rows.Scan(&secret.Id, &secret.IntegrationId, &secret.Name, &value); err != nil {
			return nil, err
		}

		data[secret] = value
	}

	return data, nil
}

// GetAll integration_id -> SecretWithValue
func (i *CustomIntegrationSecretValuesTable) GetAll(guildId uint64, integrationIds []int) (map[int][]SecretWithValue, error) {
	query := `
SELECT values.secret_id, values.integration_id, secrets.name, values.value
FROM custom_integration_secret_values AS values 
INNER JOIN custom_integration_secrets AS secrets ON secrets.id = values.secret_id
WHERE values.integration_id = ANY($1) AND values.guild_id = $2;`

	idArray := &pgtype.Int4Array{}
	if err := idArray.Set(integrationIds); err != nil {
		return nil, err
	}

	rows, err := i.Query(context.Background(), query, idArray, guildId)
	if err != nil {
		return nil, err
	}

	data := make(map[int][]SecretWithValue)
	for rows.Next() {
		var secret CustomIntegrationSecret
		var value string
		if err := rows.Scan(&secret.Id, &secret.IntegrationId, &secret.Name, &value); err != nil {
			return nil, err
		}

		if _, ok := data[secret.IntegrationId]; !ok {
			data[secret.IntegrationId] = []SecretWithValue{}
		}

		data[secret.IntegrationId] = append(data[secret.IntegrationId], SecretWithValue{
			CustomIntegrationSecret: secret,
			Value:                   value,
		})
	}

	return data, nil
}

func (i *CustomIntegrationSecretValuesTable) UpdateAll(guildId uint64, integrationId int, secrets map[int]string) error {
	tx, err := i.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	for secretId, secretValue := range secrets {
		// Must upsert, in case the secret was created after the integration was activated
		query := `
INSERT INTO custom_integration_secret_values(secret_id, integration_id, guild_id, value)
VALUES ($1, $2, $3, $4)
ON CONFLICT(secret_id, guild_id) DO UPDATE SET value = $4;`

		_, err := tx.Exec(context.Background(), query, secretId, integrationId, guildId, secretValue)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}
