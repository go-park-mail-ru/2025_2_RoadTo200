package postgres

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
)

func GetConn(ctx context.Context, pool interfaces.PoolIface) (interfaces.PgxIface, error) {
	mctx := ctx.(*httpserver.Context)
	if mctx.Conn() == nil {
		conn, err := pool.Acquire(mctx)
		if err != nil {
			return nil, err
		}
		mctx.SetConn(conn)
	}
	res, ok := mctx.Conn().(interfaces.PgxIface)
	if !ok {
		return nil, fmt.Errorf("error create connection")
	}
	return res, nil
}
