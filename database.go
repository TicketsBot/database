package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	ActiveLanguage      *ActiveLanguage
	ArchiveChannel      *ArchiveChannel
	AutoClose           *AutoCloseTable
	AutoCloseExclude    *AutoCloseExclude
	Blacklist           *Blacklist
	ChannelCategory     *ChannelCategory
	ClaimSettings       *ClaimSettingsTable
	CloseConfirmation   *CloseConfirmation
	DmOnOpen            *DmOnOpen
	FirstResponseTime   *FirstResponseTime
	ModmailArchive      *ModmailArchiveTable
	ModmailEnabled      *ModmailEnabled
	ModmailForcedGuilds *ModmailForcedGuilds
	ModmailSession      *ModmailSessionTable
	ModmailWebhook      *ModmailWebhookTable
	MultiPanels         *MultiPanelTable
	MultiPanelTargets   *MultiPanelTargets
	NamingScheme        *TicketNamingScheme
	Panel               *PanelTable
	Participants        *ParticipantTable
	PanelRoleMentions   *PanelRoleMentions
	PanelUserMention    *PanelUserMention
	Permissions         *Permissions
	PingEveryone        *PingEveryone
	Prefix              *Prefix
	PremiumGuilds       *PremiumGuilds
	PremiumKeys         *PremiumKeys
	RolePermissions     *RolePermissions
	ServerBlacklist     *ServerBlacklist
	Tag                 *Tag
	TicketClaims        *TicketClaims
	TicketLastMessage   *TicketLastMessageTable
	TicketLimit         *TicketLimit
	TicketMembers       *TicketMembers
	Tickets             *TicketTable
	Translations        *Translations
	UsedKeys            *UsedKeys
	UsersCanClose       *UsersCanClose
	UserGuilds          *UserGuildsTable
	Votes               *Votes
	Webhooks            *WebhookTable
	WelcomeMessages     *WelcomeMessages
	Whitelabel          *WhitelabelBotTable
	WhitelabelErrors    *WhitelabelErrors
	WhitelabelGuilds    *WhitelabelGuilds
	WhitelabelStatuses  *WhitelabelStatuses
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		ActiveLanguage:      newActiveLanguage(pool),
		ArchiveChannel:      newArchiveChannel(pool),
		AutoClose:           newAutoCloseTable(pool),
		AutoCloseExclude:    newAutoCloseExclude(pool),
		Blacklist:           newBlacklist(pool),
		ChannelCategory:     newChannelCategory(pool),
		ClaimSettings:       newClaimSettingsTable(pool),
		CloseConfirmation:   newCloseConfirmation(pool),
		DmOnOpen:            newDmOnOpen(pool),
		FirstResponseTime:   newFirstResponseTime(pool),
		ModmailArchive:      newModmailArchiveTable(pool),
		ModmailEnabled:      newModmailEnabled(pool),
		ModmailForcedGuilds: newModmailForcedGuilds(pool),
		ModmailSession:      newModmailSessionTable(pool),
		ModmailWebhook:      newModmailWebhookTable(pool),
		MultiPanels:         newMultiMultiPanelTable(pool),
		MultiPanelTargets:   newMultiPanelTargets(pool),
		NamingScheme:        newTicketNamingScheme(pool),
		Panel:               newPanelTable(pool),
		Participants:        newParticipantTable(pool),
		PanelRoleMentions:   newPanelRoleMentions(pool),
		PanelUserMention:    newPanelUserMention(pool),
		Permissions:         newPermissions(pool),
		PingEveryone:        newPingEveryone(pool),
		Prefix:              newPrefix(pool),
		PremiumGuilds:       newPremiumGuilds(pool),
		PremiumKeys:         newPremiumKeys(pool),
		RolePermissions:     newRolePermissions(pool),
		ServerBlacklist:     newServerBlacklist(pool),
		Tag:                 newTag(pool),
		TicketClaims:        newTicketClaims(pool),
		TicketLastMessage:   newTicketLastMessageTable(pool),
		TicketLimit:         newTicketLimit(pool),
		TicketMembers:       newTicketMembers(pool),
		Tickets:             newTicketTable(pool),
		Translations:        newTranslations(pool),
		UsedKeys:            newUsedKeys(pool),
		UsersCanClose:       newUsersCanClose(pool),
		UserGuilds:          newUserGuildsTable(pool),
		Votes:               newVotes(pool),
		Webhooks:            newWebhookTable(pool),
		WelcomeMessages:     newWelcomeMessages(pool),
		Whitelabel:          newWhitelabelBotTable(pool),
		WhitelabelErrors:    newWhitelabelErrors(pool),
		WhitelabelGuilds:    newWhitelabelGuilds(pool),
		WhitelabelStatuses:  newWhitelabelStatuses(pool),
	}
}

func (d *Database) CreateTables(pool *pgxpool.Pool) {
	mustCreate(pool, d.ActiveLanguage)
	mustCreate(pool, d.ArchiveChannel)
	mustCreate(pool, d.AutoClose)
	mustCreate(pool, d.Blacklist)
	mustCreate(pool, d.ChannelCategory)
	mustCreate(pool, d.ClaimSettings)
	mustCreate(pool, d.CloseConfirmation)
	mustCreate(pool, d.DmOnOpen)
	mustCreate(pool, d.ModmailArchive)
	mustCreate(pool, d.ModmailEnabled)
	mustCreate(pool, d.ModmailSession)
	mustCreate(pool, d.ModmailWebhook)
	mustCreate(pool, d.MultiPanels)
	mustCreate(pool, d.NamingScheme)
	mustCreate(pool, d.Panel)
	mustCreate(pool, d.MultiPanelTargets) // must be created after panels table
	mustCreate(pool, d.PanelRoleMentions)
	mustCreate(pool, d.PanelUserMention)
	mustCreate(pool, d.Permissions)
	mustCreate(pool, d.PingEveryone)
	mustCreate(pool, d.Prefix)
	mustCreate(pool, d.PremiumGuilds)
	mustCreate(pool, d.PremiumKeys)
	mustCreate(pool, d.RolePermissions)
	mustCreate(pool, d.ServerBlacklist)
	mustCreate(pool, d.Tag)
	mustCreate(pool, d.TicketLimit)
	mustCreate(pool, d.Tickets) // Must be created before members table
	mustCreate(pool, d.TicketLastMessage)
	mustCreate(pool, d.Participants) // Must be created after Tickets table
	mustCreate(pool, d.AutoCloseExclude) // Must be created after Tickets table
	mustCreate(pool, d.FirstResponseTime)
	mustCreate(pool, d.TicketMembers)
	mustCreate(pool, d.TicketClaims)
	mustCreate(pool, d.Translations)
	mustCreate(pool, d.UsedKeys)
	mustCreate(pool, d.UsersCanClose)
	mustCreate(pool, d.UserGuilds)
	mustCreate(pool, d.Votes)
	mustCreate(pool, d.Webhooks)
	mustCreate(pool, d.WelcomeMessages)
	mustCreate(pool, d.Whitelabel)
	mustCreate(pool, d.WhitelabelErrors)
	mustCreate(pool, d.WhitelabelGuilds)
	mustCreate(pool, d.WhitelabelStatuses)
	mustCreate(pool, d.ModmailForcedGuilds)
}

func mustCreate(pool *pgxpool.Pool, table Table) {
	if _, err := pool.Exec(context.Background(), table.Schema()); err != nil {
		panic(err)
	}
}
