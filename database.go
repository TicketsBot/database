package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const defaultTransactionTimeout = time.Second * 3

type Database struct {
	pool                          *pgxpool.Pool
	ActiveLanguage                *ActiveLanguage
	ArchiveChannel                *ArchiveChannel
	ArchiveMessages               *ArchiveMessages
	AutoClose                     *AutoCloseTable
	AutoCloseExclude              *AutoCloseExclude
	Blacklist                     *Blacklist
	BotStaff                      *BotStaff
	ChannelCategory               *ChannelCategory
	ClaimSettings                 *ClaimSettingsTable
	CloseConfirmation             *CloseConfirmation
	CloseReason                   *CloseMetadataTable
	CloseRequest                  *CloseRequestTable
	CustomIntegrations            *CustomIntegrationTable
	CustomIntegrationGuildCounts  *CustomIntegrationGuildCountsView
	CustomIntegrationGuilds       *CustomIntegrationGuildsTable
	CustomIntegrationHeaders      *CustomIntegrationHeadersTable
	CustomIntegrationPlaceholders *CustomIntegrationPlaceholdersTable
	CustomIntegrationSecretValues *CustomIntegrationSecretValuesTable
	CustomIntegrationSecrets      *CustomIntegrationSecretsTable
	CustomColours                 *CustomColours
	DmOnOpen                      *DmOnOpen
	EmbedFields                   *EmbedFieldsTable
	Embeds                        *EmbedsTable
	ExitSurveyResponses           *ExitSurveyResponses
	FeedbackEnabled               *FeedbackEnabled
	FirstResponseTime             *FirstResponseTime
	FirstResponseTimeGuildView    *FirstResponseTimeGuildView
	FormInput                     *FormInputTable
	Forms                         *FormsTable
	GlobalBlacklist               *GlobalBlacklist
	GuildLeaveTime                *GuildLeaveTime
	GuildMetadata                 *GuildMetadataTable
	MultiPanels                   *MultiPanelTable
	MultiPanelTargets             *MultiPanelTargets
	NamingScheme                  *TicketNamingScheme
	OnCall                        *OnCall
	Panel                         *PanelTable
	PanelTeams                    *PanelTeamsTable
	Participants                  *ParticipantTable
	PanelRoleMentions             *PanelRoleMentions
	PanelUserMention              *PanelUserMention
	Permissions                   *Permissions
	PingEveryone                  *PingEveryone
	Prefix                        *Prefix
	PremiumGuilds                 *PremiumGuilds
	PremiumKeys                   *PremiumKeys
	RoleBlacklist                 *RoleBlacklist
	RolePermissions               *RolePermissions
	ServerBlacklist               *ServerBlacklist
	ServiceRatings                *ServiceRatings
	Settings                      *SettingsTable
	StaffOverride                 *StaffOverride
	SupportTeam                   *SupportTeamTable
	SupportTeamMembers            *SupportTeamMembersTable
	SupportTeamRoles              *SupportTeamRolesTable
	Tag                           *TagsTable
	TicketClaims                  *TicketClaims
	TicketDurationView            *TicketDurationView
	TicketLastMessage             *TicketLastMessageTable
	TicketLimit                   *TicketLimit
	TicketMembers                 *TicketMembers
	TicketPermissions             *TicketPermissionsTable
	Tickets                       *TicketTable
	TopCloseReasonsView           *TopCloseReasonsView
	UsedKeys                      *UsedKeys
	UsersCanClose                 *UsersCanClose
	UserGuilds                    *UserGuildsTable
	Votes                         *Votes
	Webhooks                      *WebhookTable
	WelcomeMessages               *WelcomeMessages
	Whitelabel                    *WhitelabelBotTable
	WhitelabelErrors              *WhitelabelErrors
	WhitelabelGuilds              *WhitelabelGuilds
	WhitelabelKeys                *WhitelabelKeys
	WhitelabelStatuses            *WhitelabelStatuses
	WhitelabelUsers               *WhitelabelUsers
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	db := &Database{
		pool:                          pool,
		ActiveLanguage:                newActiveLanguage(pool),
		ArchiveChannel:                newArchiveChannel(pool),
		ArchiveMessages:               newArchiveMessages(pool),
		AutoClose:                     newAutoCloseTable(pool),
		AutoCloseExclude:              newAutoCloseExclude(pool),
		Blacklist:                     newBlacklist(pool),
		BotStaff:                      newBotStaff(pool),
		ChannelCategory:               newChannelCategory(pool),
		ClaimSettings:                 newClaimSettingsTable(pool),
		CloseConfirmation:             newCloseConfirmation(pool),
		CloseReason:                   newCloseReasonTable(pool),
		CloseRequest:                  newCloseRequestTable(pool),
		CustomIntegrations:            newCustomIntegrationTable(pool),
		CustomIntegrationGuildCounts:  newCustomIntegrationGuildCountsView(pool),
		CustomIntegrationGuilds:       newCustomIntegrationGuildsTable(pool),
		CustomIntegrationHeaders:      newCustomIntegrationHeadersTable(pool),
		CustomIntegrationPlaceholders: newCustomIntegrationPlaceholdersTable(pool),
		CustomIntegrationSecretValues: newCustomIntegrationSecretValuesTable(pool),
		CustomIntegrationSecrets:      newCustomIntegrationSecretsTable(pool),
		CustomColours:                 newCustomColours(pool),
		DmOnOpen:                      newDmOnOpen(pool),
		EmbedFields:                   newEmbedFieldsTable(pool),
		Embeds:                        newEmbedsTable(pool),
		ExitSurveyResponses:           newExitSurveyResponses(pool),
		FeedbackEnabled:               newFeedbackEnabled(pool),
		FirstResponseTime:             newFirstResponseTime(pool),
		FirstResponseTimeGuildView:    newFirstResponseTimeGuildView(pool),
		FormInput:                     newFormInputTable(pool),
		Forms:                         newFormsTable(pool),
		GlobalBlacklist:               newGlobalBlacklist(pool),
		GuildLeaveTime:                newGuildLeaveTime(pool),
		GuildMetadata:                 newGuildMetadataTable(pool),
		MultiPanels:                   newMultiMultiPanelTable(pool),
		MultiPanelTargets:             newMultiPanelTargets(pool),
		NamingScheme:                  newTicketNamingScheme(pool),
		OnCall:                        newOnCall(pool),
		Panel:                         newPanelTable(pool),
		PanelTeams:                    newPanelTeamsTable(pool),
		Participants:                  newParticipantTable(pool),
		PanelRoleMentions:             newPanelRoleMentions(pool),
		PanelUserMention:              newPanelUserMention(pool),
		Permissions:                   newPermissions(pool),
		PingEveryone:                  newPingEveryone(pool),
		Prefix:                        newPrefix(pool),
		PremiumGuilds:                 newPremiumGuilds(pool),
		PremiumKeys:                   newPremiumKeys(pool),
		RoleBlacklist:                 newRoleBlacklist(pool),
		RolePermissions:               newRolePermissions(pool),
		ServerBlacklist:               newServerBlacklist(pool),
		ServiceRatings:                newServiceRatings(pool),
		Settings:                      newSettingsTable(pool),
		StaffOverride:                 newStaffOverride(pool),
		SupportTeam:                   newSupportTeamTable(pool),
		SupportTeamMembers:            newSupportTeamMembersTable(pool),
		SupportTeamRoles:              newSupportTeamRolesTable(pool),
		Tag:                           newTag(pool),
		TicketClaims:                  newTicketClaims(pool),
		TicketDurationView:            newTicketDurationView(pool),
		TicketLastMessage:             newTicketLastMessageTable(pool),
		TicketLimit:                   newTicketLimit(pool),
		TicketMembers:                 newTicketMembers(pool),
		TicketPermissions:             newTicketPermissionsTable(pool),
		Tickets:                       newTicketTable(pool),
		TopCloseReasonsView:           newTopCloseReasonsView(pool),
		UsedKeys:                      newUsedKeys(pool),
		UsersCanClose:                 newUsersCanClose(pool),
		UserGuilds:                    newUserGuildsTable(pool),
		Votes:                         newVotes(pool),
		Webhooks:                      newWebhookTable(pool),
		WelcomeMessages:               newWelcomeMessages(pool),
		Whitelabel:                    newWhitelabelBotTable(pool),
		WhitelabelErrors:              newWhitelabelErrors(pool),
		WhitelabelGuilds:              newWhitelabelGuilds(pool),
		WhitelabelKeys:                newWhitelabelKeys(pool),
		WhitelabelStatuses:            newWhitelabelStatuses(pool),
		WhitelabelUsers:               newWhitelabelUsers(pool),
	}

	return db
}

func (d *Database) BeginTx() (pgx.Tx, error) {
	return d.pool.Begin(context.Background())
}

func (d *Database) CreateTables(pool *pgxpool.Pool) {
	mustCreate(pool,
		d.ActiveLanguage,
		d.ArchiveChannel,
		d.AutoClose,
		d.Blacklist,
		d.BotStaff,
		d.ChannelCategory,
		d.ClaimSettings,
		d.CloseConfirmation,
		d.CustomIntegrations,
		d.CustomIntegrationGuilds,
		d.CustomIntegrationGuildCounts,
		d.CustomIntegrationHeaders,
		d.CustomIntegrationPlaceholders,
		d.CustomIntegrationSecrets,
		d.CustomIntegrationSecretValues,
		d.CustomColours,
		d.DmOnOpen,
		d.Embeds,
		d.EmbedFields, // depends on embeds
		d.FeedbackEnabled,
		d.Forms,
		d.FormInput,
		d.GlobalBlacklist,
		d.GuildLeaveTime,
		d.GuildMetadata,
		d.MultiPanels,
		d.NamingScheme,
		d.OnCall,
		d.Panel,
		d.MultiPanelTargets, // must be created after panels table
		d.PanelRoleMentions,
		d.PanelUserMention,
		d.Permissions,
		d.PingEveryone,
		d.Prefix,
		d.PremiumGuilds,
		d.PremiumKeys,
		d.RoleBlacklist,
		d.RolePermissions,
		d.ServerBlacklist,
		d.Settings,
		d.StaffOverride,
		d.SupportTeam,
		d.SupportTeamMembers,
		d.SupportTeamRoles,
		d.PanelTeams, // Must be created after panels & support teams tables
		d.Tag,
		d.TicketLimit,
		d.TicketPermissions,
		d.Tickets,             // Must be created before members table
		d.TicketLastMessage,   // Must be created after Tickets table
		d.TicketDurationView,  // Must be created after Tickets table
		d.Participants,        // Must be created after Tickets table
		d.AutoCloseExclude,    // Must be created after Tickets table
		d.CloseReason,         // Must be created after Tickets table
		d.CloseRequest,        // Must be created after Tickets table
		d.ServiceRatings,      // Must be created after Tickets table
		d.ExitSurveyResponses, // Must be created after Tickets table
		d.ArchiveMessages,     // Must be created after Tickets table
		d.FirstResponseTime,
		d.FirstResponseTimeGuildView,
		d.TicketMembers,
		d.TicketClaims,
		d.TopCloseReasonsView,
		d.UsedKeys,
		d.UsersCanClose,
		d.UserGuilds,
		d.Votes,
		d.Webhooks,
		d.WelcomeMessages,
		d.Whitelabel,
		d.WhitelabelErrors,
		d.WhitelabelGuilds,
		d.WhitelabelKeys,
		d.WhitelabelStatuses,
		d.WhitelabelUsers,
	)
}

func (d *Database) Views() []View {
	return []View{
		d.TicketDurationView,
		d.FirstResponseTimeGuildView,
		d.TopCloseReasonsView,
		d.CustomIntegrationGuildCounts,
	}
}

func mustCreate(pool *pgxpool.Pool, tables ...Table) {
	for _, table := range tables {
		if _, err := pool.Exec(context.Background(), table.Schema()); err != nil {
			panic(err)
		}
	}
}
