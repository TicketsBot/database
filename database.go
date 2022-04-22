package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	ActiveLanguage     *ActiveLanguage
	ArchiveChannel     *ArchiveChannel
	AutoClose          *AutoCloseTable
	AutoCloseExclude   *AutoCloseExclude
	Blacklist          *Blacklist
	ChannelCategory    *ChannelCategory
	ClaimSettings      *ClaimSettingsTable
	CloseConfirmation  *CloseConfirmation
	CloseReason        *CloseReasonTable
	CloseRequest       *CloseRequestTable
	CustomColours      *CustomColours
	DmOnOpen           *DmOnOpen
	FeedbackEnabled    *FeedbackEnabled
	FirstResponseTime  *FirstResponseTime
	FormInput          *FormInputTable
	Forms              *FormsTable
	GlobalBlacklist    *GlobalBlacklist
	GuildLeaveTime     *GuildLeaveTime
	MultiPanels        *MultiPanelTable
	MultiPanelTargets  *MultiPanelTargets
	NamingScheme       *TicketNamingScheme
	Panel              *PanelTable
	PanelTeams         *PanelTeamsTable
	Participants       *ParticipantTable
	PanelRoleMentions  *PanelRoleMentions
	PanelUserMention   *PanelUserMention
	Permissions        *Permissions
	PingEveryone       *PingEveryone
	Prefix             *Prefix
	PremiumGuilds      *PremiumGuilds
	PremiumKeys        *PremiumKeys
	RolePermissions    *RolePermissions
	ServerBlacklist    *ServerBlacklist
	ServiceRatings     *ServiceRatings
	Settings           *SettingsTable
	SupportTeam        *SupportTeamTable
	SupportTeamMembers *SupportTeamMembersTable
	SupportTeamRoles   *SupportTeamRolesTable
	Tag                *Tag
	TicketClaims       *TicketClaims
	TicketLastMessage  *TicketLastMessageTable
	TicketLimit        *TicketLimit
	TicketMembers      *TicketMembers
	Tickets            *TicketTable
	UsedKeys           *UsedKeys
	UsersCanClose      *UsersCanClose
	UserGuilds         *UserGuildsTable
	Votes              *Votes
	Webhooks           *WebhookTable
	WelcomeMessages    *WelcomeMessages
	Whitelabel         *WhitelabelBotTable
	WhitelabelErrors   *WhitelabelErrors
	WhitelabelGuilds   *WhitelabelGuilds
	WhitelabelKeys     *WhitelabelKeys
	WhitelabelStatuses *WhitelabelStatuses
	WhitelabelUsers    *WhitelabelUsers
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		ActiveLanguage:     newActiveLanguage(pool),
		ArchiveChannel:     newArchiveChannel(pool),
		AutoClose:          newAutoCloseTable(pool),
		AutoCloseExclude:   newAutoCloseExclude(pool),
		Blacklist:          newBlacklist(pool),
		ChannelCategory:    newChannelCategory(pool),
		ClaimSettings:      newClaimSettingsTable(pool),
		CloseConfirmation:  newCloseConfirmation(pool),
		CloseReason:        newCloseReasonTable(pool),
		CloseRequest:       newCloseRequestTable(pool),
		CustomColours:      newCustomColours(pool),
		DmOnOpen:           newDmOnOpen(pool),
		FeedbackEnabled:    newFeedbackEnabled(pool),
		FirstResponseTime:  newFirstResponseTime(pool),
		FormInput:          newFormInputTable(pool),
		Forms:              newFormsTable(pool),
		GlobalBlacklist:    newGlobalBlacklist(pool),
		GuildLeaveTime:     newGuildLeaveTime(pool),
		MultiPanels:        newMultiMultiPanelTable(pool),
		MultiPanelTargets:  newMultiPanelTargets(pool),
		NamingScheme:       newTicketNamingScheme(pool),
		Panel:              newPanelTable(pool),
		PanelTeams:         newPanelTeamsTable(pool),
		Participants:       newParticipantTable(pool),
		PanelRoleMentions:  newPanelRoleMentions(pool),
		PanelUserMention:   newPanelUserMention(pool),
		Permissions:        newPermissions(pool),
		PingEveryone:       newPingEveryone(pool),
		Prefix:             newPrefix(pool),
		PremiumGuilds:      newPremiumGuilds(pool),
		PremiumKeys:        newPremiumKeys(pool),
		RolePermissions:    newRolePermissions(pool),
		ServerBlacklist:    newServerBlacklist(pool),
		ServiceRatings:     newServiceRatings(pool),
		Settings:           newSettingsTable(pool),
		SupportTeam:        newSupportTeamTable(pool),
		SupportTeamMembers: newSupportTeamMembersTable(pool),
		SupportTeamRoles:   newSupportTeamRolesTable(pool),
		Tag:                newTag(pool),
		TicketClaims:       newTicketClaims(pool),
		TicketLastMessage:  newTicketLastMessageTable(pool),
		TicketLimit:        newTicketLimit(pool),
		TicketMembers:      newTicketMembers(pool),
		Tickets:            newTicketTable(pool),
		UsedKeys:           newUsedKeys(pool),
		UsersCanClose:      newUsersCanClose(pool),
		UserGuilds:         newUserGuildsTable(pool),
		Votes:              newVotes(pool),
		Webhooks:           newWebhookTable(pool),
		WelcomeMessages:    newWelcomeMessages(pool),
		Whitelabel:         newWhitelabelBotTable(pool),
		WhitelabelErrors:   newWhitelabelErrors(pool),
		WhitelabelGuilds:   newWhitelabelGuilds(pool),
		WhitelabelKeys:     newWhitelabelKeys(pool),
		WhitelabelStatuses: newWhitelabelStatuses(pool),
		WhitelabelUsers:    newWhitelabelUsers(pool),
	}
}

func (d *Database) CreateTables(pool *pgxpool.Pool) {
	mustCreate(pool,
		d.ActiveLanguage,
		d.ArchiveChannel,
		d.AutoClose,
		d.Blacklist,
		d.ChannelCategory,
		d.ClaimSettings,
		d.CloseConfirmation,
		d.CustomColours,
		d.DmOnOpen,
		d.FeedbackEnabled,
		d.Forms,
		d.FormInput,
		d.GlobalBlacklist,
		d.GuildLeaveTime,
		d.MultiPanels,
		d.NamingScheme,
		d.Panel,
		d.MultiPanelTargets, // must be created after panels table
		d.PanelRoleMentions,
		d.PanelUserMention,
		d.Permissions,
		d.PingEveryone,
		d.Prefix,
		d.PremiumGuilds,
		d.PremiumKeys,
		d.RolePermissions,
		d.ServerBlacklist,
		d.Settings,
		d.SupportTeam,
		d.SupportTeamMembers,
		d.SupportTeamRoles,
		d.PanelTeams, // Must be created after panels & support teams tables
		d.Tag,
		d.TicketLimit,
		d.Tickets, // Must be created before members table
		d.TicketLastMessage,
		d.Participants,     // Must be created after Tickets table
		d.AutoCloseExclude, // Must be created after Tickets table
		d.CloseReason,      // Must be created after Tickets table
		d.CloseRequest,     // Must be created after Tickets table
		d.ServiceRatings,   // Must be created after Tickets table
		d.FirstResponseTime,
		d.TicketMembers,
		d.TicketClaims,
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

func mustCreate(pool *pgxpool.Pool, tables ...Table) {
	for _, table := range tables {
		if _, err := pool.Exec(context.Background(), table.Schema()); err != nil {
			panic(err)
		}
	}
}
