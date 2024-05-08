package middleware

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func TraceUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(semconv.RPCSystemGRPC)
	span.SetAttributes(semconv.RPCServiceKey.String("we-wallet-store"))
	span.SetAttributes(semconv.RPCMethodKey.String(info.FullMethod))
	resp, err := handler(ctx, req)
	if err != nil {
		code := status.Code(err)
		span.SetStatus(codes.Error, code.String())
	}
	return resp, err
	
}
