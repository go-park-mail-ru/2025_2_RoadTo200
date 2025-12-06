package httpserver

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
)

var _ context.Context = (*Context)(nil)

type Context struct {
	context context.Context
	logger  logger.Log
	conn    interface{}
}

func NewContext(ctx context.Context, logger logger.Log) *Context {
	return &Context{
		context: ctx,
		logger:  logger,
	}
}

func (ctx *Context) Logger() logger.Log {
	return ctx.logger
}

func (ctx *Context) Conn() interface{} {
	return ctx.conn
}

func (ctx *Context) SetConn(conn interface{}) {
	ctx.conn = conn
}

func (ctx *Context) Value(key any) any {
	return ctx.context.Value(key)
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.context.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

func (ctx *Context) Err() error {
	return ctx.context.Err()
}
