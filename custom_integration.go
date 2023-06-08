package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CustomIntegrationTable struct {
	*pgxpool.Pool
}

type CustomIntegration struct {
	Id               int     `json:"id"`
	OwnerId          uint64  `json:"owner_id"`
	HttpMethod       string  `json:"http_method"`
	WebhookUrl       string  `json:"webhook_url"`
	ValidationUrl    *string `json:"validation_url"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ImageUrl         *string `json:"image_url"`
	PrivacyPolicyUrl *string `json:"privacy_policy_url"`
	Public           bool    `json:"public"`
	Approved         bool    `json:"approved"`
}

type CustomIntegrationWithGuildCount struct {
	CustomIntegration
	GuildCount int `json:"guild_count"`
}

type CustomIntegrationWithActive struct {
	CustomIntegrationWithGuildCount
	Active bool `json:"active"`
}

func newCustomIntegrationTable(db *pgxpool.Pool) *CustomIntegrationTable {
	return &CustomIntegrationTable{
		db,
	}
}

func (i CustomIntegrationTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS custom_integrations(
	"id" SERIAL NOT NULL UNIQUE,
	"owner_id" int8 NOT NULL,
	"webhook_url" VARCHAR(255) NOT NULL,
	"validation_url" VARCHAR(255) NULL,
    "http_method" VARCHAR(4) NOT NULL,
	"name" VARCHAR(32) NOT NULL,
	"description" VARCHAR(255) NOT NULL,
	"image_url" VARCHAR(255) NULL,
	"privacy_policy_url" VARCHAR(255) NULL,
	"public" BOOL NOT NULL DEFAULT 'f',
	"approved" BOOL NOT NULL DEFAULT 'f',
	PRIMARY KEY("id")
);
CREATE INDEX IF NOT EXISTS custom_integrations_owner_id ON custom_integrations("owner_id");
`
}

func (i *CustomIntegrationTable) Get(id int) (CustomIntegration, bool, error) {
	query := `SELECT "id", "owner_id", "webhook_url", "validation_url", "http_method", "name", "description", "image_url", "privacy_policy_url", "public", "approved" FROM custom_integrations WHERE "id" = $1;`

	var integration CustomIntegration
	err := i.QueryRow(context.Background(), query, id).Scan(
		&integration.Id,
		&integration.OwnerId,
		&integration.WebhookUrl,
		&integration.ValidationUrl,
		&integration.HttpMethod,
		&integration.Name,
		&integration.Description,
		&integration.ImageUrl,
		&integration.PrivacyPolicyUrl,
		&integration.Public,
		&integration.Approved,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return integration, false, nil
		} else {
			return CustomIntegration{}, false, err
		}
	}

	return integration, true, nil
}

func (i *CustomIntegrationTable) GetAll(ids []int) ([]CustomIntegration, error) {
	query := `
SELECT "id", "owner_id", "webhook_url", "validation_url", "http_method", "name", "description", "image_url", "privacy_policy_url", "public", "approved"
FROM custom_integrations
WHERE "id" = ANY($1);`

	idArray := &pgtype.Int4Array{}
	if err := idArray.Set(ids); err != nil {
		return nil, err
	}

	rows, err := i.Query(context.Background(), query, idArray)
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
			&integration.ValidationUrl,
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

func (i *CustomIntegrationTable) GetOwnedCount(userId uint64) (count int, err error) {
	query := `SELECT COUNT(*) FROM custom_integrations WHERE "owner_id" = $1;`
	err = i.QueryRow(context.Background(), query, userId).Scan(&count)
	return
}

func (i *CustomIntegrationTable) GetAllOwned(ownerId uint64) ([]CustomIntegrationWithGuildCount, error) {
	query := `
SELECT
	integrations.id,
	integrations.owner_id,
	integrations.webhook_url,
	integrations.validation_url,
	integrations.http_method,
	integrations.name,
	integrations.description,
	integrations.image_url,
	integrations.privacy_policy_url,
	integrations.public,
	integrations.approved,
	COALESCE(counts.count, 0) AS guild_count
FROM custom_integrations AS integrations
LEFT OUTER JOIN custom_integration_guild_counts counts ON integrations.id = counts.integration_id
WHERE "owner_id" = $1;`

	rows, err := i.Query(context.Background(), query, ownerId)
	if err != nil {
		return nil, err
	}

	var integrations []CustomIntegrationWithGuildCount
	for rows.Next() {
		var integration CustomIntegrationWithGuildCount
		err := rows.Scan(
			&integration.Id,
			&integration.OwnerId,
			&integration.WebhookUrl,
			&integration.ValidationUrl,
			&integration.HttpMethod,
			&integration.Name,
			&integration.Description,
			&integration.ImageUrl,
			&integration.PrivacyPolicyUrl,
			&integration.Public,
			&integration.Approved,
			&integration.GuildCount,
		)

		if err != nil {
			return nil, err
		}

		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (i *CustomIntegrationGuildsTable) GetAvailableIntegrationsWithActive(guildId, userId uint64, limit, offset int) ([]CustomIntegrationWithActive, error) {
	query := `
WITH active AS (
	SELECT integration_id
	FROM custom_integration_guilds
	WHERE guild_id=$1
)
SELECT
	integrations.id,
	integrations.owner_id,
	integrations.webhook_url,
	integrations.validation_url,
	integrations.http_method,
	integrations.name,
	integrations.description,
	integrations.image_url,
	integrations.privacy_policy_url,
	integrations.public,
	integrations.approved,
	COALESCE(counts.count, 0) AS guild_count,
	CASE WHEN active.integration_id IS NOT NULL THEN TRUE ELSE FALSE END AS added
FROM custom_integrations as integrations
LEFT OUTER JOIN active ON active.integration_id = integrations.id
LEFT OUTER JOIN custom_integration_guild_counts counts ON integrations.id = counts.integration_id
WHERE active.integration_id IS NOT NULL OR 
	((integrations.public = 't' AND integrations.approved = 't') OR integrations.owner_id = $2)
ORDER BY active.integration_id NULLS LAST, guild_count DESC
LIMIT $3 OFFSET $4;
`

	rows, err := i.Query(context.Background(), query, guildId, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	var integrations []CustomIntegrationWithActive
	for rows.Next() {
		var integration CustomIntegrationWithActive
		err := rows.Scan(
			&integration.Id,
			&integration.OwnerId,
			&integration.WebhookUrl,
			&integration.ValidationUrl,
			&integration.HttpMethod,
			&integration.Name,
			&integration.Description,
			&integration.ImageUrl,
			&integration.PrivacyPolicyUrl,
			&integration.Public,
			&integration.Approved,
			&integration.GuildCount,
			&integration.Active,
		)

		if err != nil {
			return nil, err
		}

		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (i *CustomIntegrationGuildsTable) CanActivate(integrationId int, userId uint64) (canActivate bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM custom_integrations AS integrations
	WHERE integrations.id = $1 AND ((integrations.public = 't' AND integrations.approved = 't') OR integrations.owner_id = $2)
);
`

	err = i.QueryRow(context.Background(), query, integrationId, userId).Scan(&canActivate)
	return
}

func (i *CustomIntegrationTable) Create(ownerId uint64, webhookUrl string, validationUrl *string, httpMethod, name, description string, imageUrl, privacyPolicyUrl *string) (CustomIntegration, error) {
	query := `
INSERT INTO custom_integrations("owner_id", "webhook_url", "validation_url", "http_method", "name", "description", "image_url", "privacy_policy_url", "public", "approved")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'f', 'f')
RETURNING "id";
;`

	integration := CustomIntegration{
		OwnerId:          ownerId,
		Name:             name,
		WebhookUrl:       webhookUrl,
		ValidationUrl:    validationUrl,
		HttpMethod:       httpMethod,
		Description:      description,
		ImageUrl:         imageUrl,
		PrivacyPolicyUrl: privacyPolicyUrl,
		Public:           false,
		Approved:         false,
	}

	if err := i.QueryRow(context.Background(), query, ownerId, webhookUrl, validationUrl, httpMethod, name, description, imageUrl, privacyPolicyUrl).Scan(&integration.Id); err != nil {
		return CustomIntegration{}, err
	}

	return integration, nil
}

func (i *CustomIntegrationTable) SetPublic(integrationId int) (err error) {
	query := `UPDATE custom_integrations SET "public" = TRUE WHERE "id" = $1;`
	_, err = i.Exec(context.Background(), query, integrationId)
	return
}

func (i *CustomIntegrationTable) Update(integration CustomIntegration) (err error) {
	query := `
UPDATE custom_integrations
SET
	"webhook_url" = $2,
	"validation_url" = $3,
	"http_method" = $4,
	"name" = $5,
	"description" = $6,
	"image_url" = $7,
	"privacy_policy_url" = $8,
	"public" = $9,
	"approved" = $10
WHERE "id" = $1;
`

	_, err = i.Exec(
		context.Background(),
		query,
		integration.Id,
		integration.WebhookUrl,
		integration.ValidationUrl,
		integration.HttpMethod,
		integration.Name,
		integration.Description,
		integration.ImageUrl,
		integration.PrivacyPolicyUrl,
		integration.Public,
		integration.Approved,
	)

	return
}

func (i *CustomIntegrationTable) Delete(id int) (err error) {
	query := `
DELETE FROM custom_integrations
WHERE "id" = $1;
`

	_, err = i.Exec(context.Background(), query, id)
	return
}
