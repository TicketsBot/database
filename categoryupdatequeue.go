package database

import (
	"context"
	_ "embed"
	"github.com/TicketsBot/common/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type CategoryUpdateQueue struct {
	*pgxpool.Pool
}

type CategoryUpdateQueueItem struct {
	GuildId   uint64
	TicketId  int
	NewStatus model.TicketStatus
	ChannelId *uint64
	PanelId   *int
}

var (
	//go:embed sql/category_update_queue/schema.sql
	categoryUpdateQueueSchema string

	//go:embed sql/category_update_queue/add.sql
	categoryUpdateQueueAdd string

	//go:embed sql/category_update_queue/get_ready_for_update.sql
	categoryUpdateQueueGetReadyForUpdate string
)

func newCategoryUpdateQueueTable(db *pgxpool.Pool) *CategoryUpdateQueue {
	return &CategoryUpdateQueue{
		db,
	}
}

func (CategoryUpdateQueue) Schema() string {
	return categoryUpdateQueueSchema
}

func (q *CategoryUpdateQueue) Add(ctx context.Context, guildId uint64, ticketId int, newStatus model.TicketStatus) error {
	_, err := q.Exec(ctx, categoryUpdateQueueAdd, guildId, ticketId, newStatus)
	return err
}

func (q *CategoryUpdateQueue) GetReadyForUpdate(ctx context.Context, delayInterval time.Duration) ([]CategoryUpdateQueueItem, error) {
	rows, err := q.Query(ctx, categoryUpdateQueueGetReadyForUpdate, delayInterval)
	if err != nil {
		return nil, err
	}

	var items []CategoryUpdateQueueItem
	for rows.Next() {
		var item CategoryUpdateQueueItem
		if err := rows.Scan(&item.GuildId, &item.TicketId, &item.NewStatus, &item.ChannelId, &item.PanelId); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
