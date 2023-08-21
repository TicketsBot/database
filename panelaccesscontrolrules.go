package database

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AccessControlAction string

const (
	AccessControlActionAllow AccessControlAction = "allow"
	AccessControlActionDeny  AccessControlAction = "deny"
)

var ErrNoRuleMatched = errors.New("no rule matched")

type PanelAccessControlRule struct {
	RoleId uint64              `json:"role_id,string"`
	Action AccessControlAction `json:"action"`
}

type PanelAccessControlRules struct {
	*pgxpool.Pool
}

func newPanelAccessControlRules(db *pgxpool.Pool) *PanelAccessControlRules {
	return &PanelAccessControlRules{
		db,
	}
}

var (
	//go:embed sql/panel_access_control_rules/schema.sql
	panelAccessControlRulesSchema string

	//go:embed sql/panel_access_control_rules/delete_rules.sql
	panelAccessControlRulesDelete string

	//go:embed sql/panel_access_control_rules/get_all.sql
	panelAccessControlRulesGetAll string

	//go:embed sql/panel_access_control_rules/get_all_for_guild.sql
	panelAccessControlRulesGetAllForGuild string

	//go:embed sql/panel_access_control_rules/get_first_matched.sql
	panelAccessControlRulesGetFirstMatched string

	//go:embed sql/panel_access_control_rules/insert.sql
	panelAccessControlRulesInsert string
)

func (p PanelAccessControlRules) Schema() string {
	return panelAccessControlRulesSchema
}

func (p *PanelAccessControlRules) GetAll(ctx context.Context, panelId int) ([]PanelAccessControlRule, error) {
	rows, err := p.Query(ctx, panelAccessControlRulesGetAll, panelId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rules := make([]PanelAccessControlRule, 0, 10)
	for rows.Next() {
		var rule PanelAccessControlRule
		if err := rows.Scan(&rule.RoleId, &rule.Action); err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules[:], nil
}

// GetAllForGuild returns a map[panel_id][]rules
func (p *PanelAccessControlRules) GetAllForGuild(ctx context.Context, guildId uint64) (map[int][]PanelAccessControlRule, error) {
	rows, err := p.Query(ctx, panelAccessControlRulesGetAllForGuild, guildId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rules := make(map[int][]PanelAccessControlRule)
	for rows.Next() {
		var panelId int
		var rule PanelAccessControlRule
		if err := rows.Scan(&panelId, &rule.RoleId, &rule.Action); err != nil {
			return nil, err
		}

		rules[panelId] = append(rules[panelId], rule)
	}

	return rules, nil
}

func (p *PanelAccessControlRules) GetFirstMatched(ctx context.Context, panelId int, userRoles []uint64) (uint64, AccessControlAction, error) {
	idArray := &pgtype.Int8Array{}
	if err := idArray.Set(userRoles); err != nil {
		return 0, "", err
	}

	var roleId uint64
	var action AccessControlAction
	if err := p.QueryRow(ctx, panelAccessControlRulesGetFirstMatched, panelId, idArray).Scan(&roleId, &action); err != nil {
		if err == pgx.ErrNoRows {
			return 0, "", ErrNoRuleMatched
		}

		return 0, "", err
	}

	return roleId, action, nil
}

func (p *PanelAccessControlRules) Replace(ctx context.Context, panelId int, rules []PanelAccessControlRule) error {
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := p.ReplaceWithTx(ctx, tx, panelId, rules); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PanelAccessControlRules) ReplaceWithTx(ctx context.Context, tx pgx.Tx, panelId int, rules []PanelAccessControlRule) error {
	// Remove existing rules
	if _, err := tx.Exec(ctx, panelAccessControlRulesDelete, panelId); err != nil {
		return err
	}

	// Add each rule
	for position, rule := range rules {
		if _, err := tx.Exec(ctx, panelAccessControlRulesInsert, panelId, rule.RoleId, position, rule.Action); err != nil {
			return err
		}
	}

	return nil
}
