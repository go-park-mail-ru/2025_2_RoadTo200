package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"google.golang.org/grpc"
)

// ContextInjectorInterceptor внедряет кастомный контекст в gRPC вызовы
func ContextInjectorInterceptor(logger logger.Log) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Trace("Create context")
		// Создаем новый кастомный контекст
		customCtx := httpserver.NewContext(ctx, logger)

		// Запускаем обработчик с кастомным контекстом
		return handler(customCtx, req)
	}
}

// StreamContextInjectorInterceptor для потоковых вызовов
func StreamContextInjectorInterceptor(logger logger.Log) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logger.Trace("Create stream context")
		// Создаем кастомный контекст
		customCtx := httpserver.NewContext(ss.Context(), logger)

		// Обертка потока с кастомным контекстом
		wrappedStream := &serverStreamWrapper{
			ServerStream: ss,
			ctx:          customCtx,
		}

		return handler(srv, wrappedStream)
	}
}

// serverStreamWrapper оборачивает ServerStream для использования кастомного контекста
type serverStreamWrapper struct {
	grpc.ServerStream
	ctx *httpserver.Context
}

func (w *serverStreamWrapper) Context() context.Context {
	return w.ctx
}

// ToContext пытается преобразовать стандартный контекст в наш кастомный
func ToContext(ctx context.Context) (*httpserver.Context, bool) {
	if customCtx, ok := ctx.(*httpserver.Context); ok {
		return customCtx, true
	}
	return nil, false
}

// MustContext преобразует контекст или паникует
func MustContext(ctx context.Context) *httpserver.Context {
	if customCtx, ok := ctx.(*httpserver.Context); ok {
		return customCtx
	}
	panic("context is not of type *grpc.Context")
}
