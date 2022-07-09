package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationGuildsTable struct {
	*pgxpool.Pool
}

func newCustomIntegrationGuildsTable(db *pgxpool.Pool) *CustomIntegrationGuildsTable {
	return &CustomIntegrationGuildsTable{
		db,
	}
}

func (i CustomIntegrationGuildsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS custom_integration_guilds(
	"integration_id" int NOT NULL,
	"guild_id" int8 NOT NULL,
	FOREIGN KEY("integration_id") REFERENCES custom_integrations("id") ON DELETE CASCADE,
	PRIMARY KEY("integration_id", "guild_id")
);
CREATE INDEX IF NOT EXISTS custom_integration_guilds_guild_id ON custom_integration_guilds("guild_id");
`
}

func (i *CustomIntegrationGuildsTable) GetGuildIntegrations(guildId uint64) ([]CustomIntegration, error) {
	query := `
SELECT integrations.id, integrations.owner_id, integrations.webhook_url, integrations.http_method, integrations.name, integrations.description, integrations.image_url, integrations.privacy_policy_url, integrations.public, integrations.approved
FROM custom_integration_guilds
INNER JOIN custom_integrations AS integrations ON custom_integration_guilds.integration_id = integrations.id
WHERE custom_integration_guilds.guild_id = $1;
`

	rows, err := i.Query(context.Background(), query, guildId)
	if err != nil {
		return nil, err
	}

	var integrations []CustomIntegration
	for rows.Next() {
		var integration CustomIntegration
		err := rows.Scan(
			&integration.Id,
			&integration.OwnerId,
			&integration.WebhookUrl,
			&integration.HttpMethod,
			&integration.Name,
			&integration.Description,
			&integration.ImageUrl,
			&integration.PrivacyPolicyUrl,
			&integration.Public,
			&integration.Approved,
		)

		if err != nil {
			return nil, err
		}

		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (i *CustomIntegrationGuildsTable) GetGuildIntegrationCount(guildId uint64) (count int, err error) {
	query := `SELECT COUNT(*) FROM custom_integration_guilds WHERE "guild_id" = $1;`
	err = i.QueryRow(context.Background(), query, guildId).Scan(&count)
	return
}

func (i *CustomIntegrationGuildsTable) IsActive(integrationId int, guildId uint64) (isActive bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM custom_integration_guilds
	WHERE "integration_id" = $1 AND "guild_id" = $2
);`

	err = i.QueryRow(context.Background(), query, integrationId, guildId).Scan(&isActive)
	return
}

func (i *CustomIntegrationGuildsTable) AddToGuild(integrationId int, guildId uint64) error {
	query := `
INSERT INTO custom_integration_guilds("integration_id", "guild_id")
VALUES ($1, $2)
ON CONFLICT ("integration_id", "guild_id") DO NOTHING;`

	_, err := i.Exec(context.Background(), query, integrationId, guildId)
	return err
}

func (i *CustomIntegrationGuildsTable) AddToGuildWithSecrets(integrationId int, guildId uint64, secrets map[int]string) error {
	tx, err := i.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	// Add integration to guild
	{
		query := `
INSERT INTO custom_integration_guilds("integration_id", "guild_id")
VALUES ($1, $2)
ON CONFLICT ("integration_id", "guild_id") DO NOTHING;`

		if _, err := tx.Exec(context.Background(), query, integrationId, guildId); err != nil {
			return err
		}
	}

	// Add secrets to guild
	for secretId, value := range secrets {
		query := `
INSERT INTO custom_integration_secret_values("secret_id", "integration_id", "guild_id", "value")
VALUES($1, $2, $3, $4);
`

		if _, err := tx.Exec(context.Background(), query, secretId, integrationId, guildId, value); err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (i *CustomIntegrationGuildsTable) RemoveFromGuild(integrationId int, guildId uint64) (err error) {
	query := `
DELETE FROM custom_integration_guilds
WHERE "integration_id" = $1 AND "guild_id" = $2;`

	_, err = i.Exec(context.Background(), query, integrationId, guildId)
	return
}
