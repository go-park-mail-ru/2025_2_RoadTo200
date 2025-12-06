package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.TransactionRepository = (*TransactionRepository)(nil)

type TransactionRepository struct {
	pool interfaces.PgxIface
}

func NewTransactionRepository(pool interfaces.PgxIface) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (tr *TransactionRepository) Begin(ctx context.Context) error {
	mctx, ok := ctx.(*httpserver.Context)
	if !ok {
		return ctx.Err()
	}
	tx, err := tr.pool.Begin(ctx)
	if err != nil {
		return err
	}
	mctx.SetConn(tx)
	return nil
}

func (tr *TransactionRepository) Commit(ctx context.Context) error {
	mctx, ok := ctx.(*httpserver.Context)
	if !ok {
		return ctx.Err()
	}
	tx, ok := mctx.Conn().(*pgx.Tx)
	if !ok {
		return ctx.Err()
	}
	err := (*tx).Commit(ctx)
	mctx.SetConn(nil)
	return err
}

func (tr *TransactionRepository) Rollback(ctx context.Context) error {
	mctx, ok := ctx.(*httpserver.Context)
	if !ok {
		return ctx.Err()
	}
	tx, ok := mctx.Conn().(*pgx.Tx)
	if !ok {
		return ctx.Err()
	}
	err := (*tx).Rollback(ctx)
	mctx.SetConn(nil)
	return err
}
