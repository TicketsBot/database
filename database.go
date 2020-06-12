package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	ArchiveChannel     *ArchiveChannel
	AutoClose          *AutoCloseTable
	Blacklist          *Blacklist
	ChannelCategory    *ChannelCategory
	ClaimSettings      *ClaimSettingsTable
	CloseConfirmation  *CloseConfirmation
	DmOnOpen           *DmOnOpen
	FirstResponseTime  *FirstResponseTime
	ModmailArchive     *ModmailArchiveTable
	ModmailEnabled     *ModmailEnabled
	ModmailSession     *ModmailSessionTable
	ModmailWebhook     *ModmailWebhookTable
	NamingScheme       *TicketNamingScheme
	Panel              *PanelTable
	Permissions        *Permissions
	PingEveryone       *PingEveryone
	Prefix             *Prefix
	PremiumGuilds      *PremiumGuilds
	PremiumKeys        *PremiumKeys
	RolePermissions    *RolePermissions
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
	WhitelabelStatuses *WhitelabelStatuses
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		ArchiveChannel:     newArchiveChannel(pool),
		AutoClose:          newAutoCloseTable(pool),
		Blacklist:          newBlacklist(pool),
		ChannelCategory:    newChannelCategory(pool),
		ClaimSettings:      newClaimSettingsTable(pool),
		CloseConfirmation:  newCloseConfirmation(pool),
		DmOnOpen:           newDmOnOpen(pool),
		FirstResponseTime:  newFirstResponseTime(pool),
		ModmailArchive:     newModmailArchiveTable(pool),
		ModmailEnabled:     newModmailEnabled(pool),
		ModmailSession:     newModmailSessionTable(pool),
		ModmailWebhook:     newModmailWebhookTable(pool),
		NamingScheme:       newTicketNamingScheme(pool),
		Panel:              newPanelTable(pool),
		Permissions:        newPermissions(pool),
		PingEveryone:       newPingEveryone(pool),
		Prefix:             newPrefix(pool),
		PremiumGuilds:      newPremiumGuilds(pool),
		PremiumKeys:        newPremiumKeys(pool),
		RolePermissions:    newRolePermissions(pool),
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
		WhitelabelStatuses: newWhitelabelStatuses(pool),
	}
}

func (d *Database) CreateTables(pool *pgxpool.Pool) {
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
	mustCreate(pool, d.NamingScheme)
	mustCreate(pool, d.Panel)
	mustCreate(pool, d.Permissions)
	mustCreate(pool, d.PingEveryone)
	mustCreate(pool, d.Prefix)
	mustCreate(pool, d.PremiumGuilds)
	mustCreate(pool, d.PremiumKeys)
	mustCreate(pool, d.RolePermissions)
	mustCreate(pool, d.Tag)
	mustCreate(pool, d.TicketLimit)
	mustCreate(pool, d.Tickets) // Must be created before members table
	mustCreate(pool, d.TicketLastMessage)
	mustCreate(pool, d.FirstResponseTime)
	mustCreate(pool, d.TicketMembers)
	mustCreate(pool, d.TicketClaims)
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
}

func mustCreate(pool *pgxpool.Pool, table Table) {
	if _, err := pool.Exec(context.Background(), table.Schema()); err != nil {
		panic(err)
	}
}
