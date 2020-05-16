package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	ArchiveChannel    *ArchiveChannel
	Blacklist         *Blacklist
	ChannelCategory   *ChannelCategory
	ClaimSettings     *ClaimSettingsTable
	CloseConfirmation *CloseConfirmation
	DmOnOpen          *DmOnOpen
	FirstResponseTime *FirstResponseTime
	ModmailArchive    *ModmailArchiveTable
	ModmailSession    *ModmailSessionTable
	ModmailWebhook    *ModmailWebhookTable
	NamingScheme      *TicketNamingScheme
	Panel             *PanelTable
	Permissions       *Permissions
	PingEveryone      *PingEveryone
	Prefix            *Prefix
	PremiumGuilds     *PremiumGuilds
	PremiumKeys       *PremiumKeys
	RolePermissions   *RolePermissions
	Tag               *Tag
	TicketClaims      *TicketClaims
	TicketLimit       *TicketLimit
	TicketMembers     *TicketMembers
	Tickets           *TicketTable
	UsedKeys          *UsedKeys
	UsersCanClose     *UsersCanClose
	UserGuilds        *UserGuildsTable
	Votes             *Votes
	Webhooks          *WebhookTable
	WelcomeMessages   *WelcomeMessages
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		ArchiveChannel:    newArchiveChannel(pool),
		Blacklist:         newBlacklist(pool),
		ChannelCategory:   newChannelCategory(pool),
		ClaimSettings:     newClaimSettingsTable(pool),
		CloseConfirmation: newCloseConfirmation(pool),
		DmOnOpen:          newDmOnOpen(pool),
		FirstResponseTime: newFirstResponseTime(pool),
		ModmailArchive:    newModmailArchiveTable(pool),
		ModmailSession:    newModmailSessionTable(pool),
		ModmailWebhook:    newModmailWebhookTable(pool),
		NamingScheme:      newTicketNamingScheme(pool),
		Panel:             newPanelTable(pool),
		Permissions:       newPermissions(pool),
		PingEveryone:      newPingEveryone(pool),
		Prefix:            newPrefix(pool),
		PremiumGuilds:     newPremiumGuilds(pool),
		PremiumKeys:       newPremiumKeys(pool),
		RolePermissions:   newRolePermissions(pool),
		Tag:               newTag(pool),
		TicketClaims:      newTicketClaims(pool),
		TicketLimit:       newTicketLimit(pool),
		TicketMembers:     newTicketMembers(pool),
		Tickets:           newTicketTable(pool),
		UsedKeys:          newUsedKeys(pool),
		UsersCanClose:     newUsersCanClose(pool),
		UserGuilds:        newUserGuildsTable(pool),
		Votes:             newVotes(pool),
		Webhooks:          newWebhookTable(pool),
		WelcomeMessages:   newWelcomeMessages(pool),
	}
}

func (d *Database) CreateTables(pool *pgxpool.Pool) {
	mustCreate(pool, d.ArchiveChannel)
	mustCreate(pool, d.Blacklist)
	mustCreate(pool, d.ChannelCategory)
	mustCreate(pool, d.ClaimSettings)
	mustCreate(pool, d.DmOnOpen)
	mustCreate(pool, d.ModmailArchive)
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
	mustCreate(pool, d.FirstResponseTime)
	mustCreate(pool, d.TicketMembers)
	mustCreate(pool, d.TicketClaims)
	mustCreate(pool, d.UsedKeys)
	mustCreate(pool, d.UsersCanClose)
	mustCreate(pool, d.UserGuilds)
	mustCreate(pool, d.Votes)
	mustCreate(pool, d.Webhooks)
	mustCreate(pool, d.WelcomeMessages)
}

func mustCreate(pool *pgxpool.Pool, table Table) {
	if _, err := pool.Exec(context.Background(), table.Schema()); err != nil {
		panic(err)
	}
}
