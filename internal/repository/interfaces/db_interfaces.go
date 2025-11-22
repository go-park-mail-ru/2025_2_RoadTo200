package interfaces

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PgxIface interface {
	pgx.Tx
	Begin(context.Context) (pgx.Tx, error)
	Close()
}
