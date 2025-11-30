package interfaces

import (
	"context"
)

type TransactionRepository interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
