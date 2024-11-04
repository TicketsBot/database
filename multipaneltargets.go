package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MultiPanelTargets struct {
	*pgxpool.Pool
}

func newMultiPanelTargets(db *pgxpool.Pool) *MultiPanelTargets {
	return &MultiPanelTargets{
		db,
	}
}

func (p MultiPanelTargets) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS multi_panel_targets(
	"multi_panel_id" int4 NOT NULL,
	"panel_id" int NOT NULL,
	FOREIGN KEY("multi_panel_id") REFERENCES multi_panels("id") ON DELETE CASCADE,
	FOREIGN KEY ("panel_id") REFERENCES panels("panel_id") ON DELETE CASCADE,
	PRIMARY KEY("multi_panel_id", "panel_id")
);
CREATE INDEX IF NOT EXISTS multi_panel_targets_multi_panel_id ON multi_panel_targets("multi_panel_id");
`
}

func (p *MultiPanelTargets) GetPanels(ctx context.Context, multiPanelId int) (panels []Panel, e error) {
	query := `
SELECT
	panels.panel_id,
	panels.message_id,
	panels.channel_id,
	panels.guild_id,
	panels.title,
	panels.content,
	panels.colour,
	panels.target_category,
	panels.emoji_name,
	panels.emoji_id,
	panels.welcome_message,
	panels.default_team,
	panels.custom_id,
	panels.image_url,
	panels.thumbnail_url,
	panels.button_style,
	panels.button_label,
	panels.form_id,
	panels.naming_scheme,
	panels.force_disabled,
	panels.disabled,
	panels.exit_survey_form_id,
	panels.pending_category
FROM multi_panel_targets
INNER JOIN panels
ON panels.panel_id = multi_panel_targets.panel_id
WHERE "multi_panel_id" = $1;`

	rows, err := p.Query(ctx, query, multiPanelId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var panel Panel
		if err := rows.Scan(panel.fieldPtrs()...); err != nil {
			return nil, err
		}

		panels = append(panels, panel)
	}

	return
}

func (p *MultiPanelTargets) GetMultiPanels(ctx context.Context, panelId int) ([]MultiPanel, error) {
	query := `
SELECT
	multi_panels.id,
	multi_panels.message_id,
	multi_panels.channel_id,
	multi_panels.guild_id,
	multi_panels.select_menu,
	multi_panels.select_menu_placeholder,
	multi_panels.embed
FROM multi_panel_targets
INNER JOIN multi_panels
ON multi_panels.id = multi_panel_targets.multi_panel_id
WHERE multi_panel_targets.panel_id = $1;
`

	rows, err := p.Query(ctx, query, panelId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var multiPanels []MultiPanel
	for rows.Next() {
		var multiPanel MultiPanel
		var embedRaw *string
		err := rows.Scan(
			&multiPanel.Id,
			&multiPanel.MessageId,
			&multiPanel.ChannelId,
			&multiPanel.GuildId,
			&multiPanel.SelectMenu,
			&multiPanel.SelectMenuPlaceholder,
			&embedRaw,
		)

		if err != nil {
			return nil, err
		}

		if embedRaw != nil {
			if err := json.Unmarshal([]byte(*embedRaw), &multiPanel.Embed); err != nil {
				return nil, err
			}
		}

		multiPanels = append(multiPanels, multiPanel)
	}

	return multiPanels, nil
}

func (p *MultiPanelTargets) Insert(ctx context.Context, multiPanelId, panelId int) (err error) {
	query := `
INSERT INTO multi_panel_targets("multi_panel_id", "panel_id")
VALUES ($1, $2) 
ON CONFLICT("multi_panel_id", "panel_id") DO NOTHING;
`

	_, err = p.Exec(ctx, query, multiPanelId, panelId)
	return
}

func (p *MultiPanelTargets) DeleteAll(ctx context.Context, multiPanelId int) (err error) {
	query := `
DELETE FROM
	multi_panel_targets
WHERE
	"multi_panel_id"=$1
;`

	_, err = p.Exec(ctx, query, multiPanelId)
	return
}

func (p *MultiPanelTargets) Delete(ctx context.Context, multiPanelId, panelId int) (err error) {
	query := `
DELETE FROM
	multi_panel_targets
WHERE
	"multi_panel_id"=$1
	AND
	"panel_id" = $2
;`

	_, err = p.Exec(ctx, query, multiPanelId, panelId)
	return
}
