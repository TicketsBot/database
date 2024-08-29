package database

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VoteCredits struct {
	*pgxpool.Pool
}

var (
	//go:embed sql/vote_credits/schema.sql
	voteCreditsSchema string

	//go:embed sql/vote_credits/get.sql
	voteCreditsGet string

	//go:embed sql/vote_credits/increment.sql
	voteCreditsIncrement string

	//go:embed sql/vote_credits/delete.sql
	voteCreditsDelete string
)

func newVoteCreditsTable(db *pgxpool.Pool) *VoteCredits {
	return &VoteCredits{
		db,
	}
}

func (VoteCredits) Schema() string {
	return voteCreditsSchema
}

func (v *VoteCredits) Get(ctx context.Context, tx pgx.Tx, userId uint64) (int, error) {
	var credits int
	if err := tx.QueryRow(ctx, voteCreditsGet, userId).Scan(&credits); err != nil {
		// pgx.ErrNoRows is impossible, as we use coalesce in the query
		return 0, err
	}

	return credits, nil
}

func (v *VoteCredits) Increment(ctx context.Context, tx pgx.Tx, userId uint64) error {
	_, err := tx.Exec(ctx, voteCreditsIncrement, userId)
	return err
}

func (v *VoteCredits) Delete(ctx context.Context, tx pgx.Tx, userId uint64) error {
	_, err := tx.Exec(ctx, voteCreditsDelete, userId)
	return err
}
