package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
)

func GetConn(ctx context.Context, pool interfaces.PgxIface) (interfaces.PgxIface, error) {
	mctx := ctx.(*httpserver.Context)
	if mctx.Conn() == nil {
		return pool, nil
	}
	res, ok := mctx.Conn().(interfaces.PgxIface)
	if !ok {
		return nil, fmt.Errorf("error create connection")
	}
	return res, nil
}

func GetSession(ctx context.Context) (logger.Log, error) {
	mctx := ctx.(*httpserver.Context)
	if mctx.Logger() == nil {
		return nil, errors.New("no logger settup")
	}
	res, ok := mctx.Logger().(logger.Log)
	if !ok {
		return nil, fmt.Errorf("error create connection")
	}
	return res, nil
}
