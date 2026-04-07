package grpcserver

import (
	"context"

	"github.com/alikhanturusbekov/go-url-shortener/pkg/authorization"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthUnaryInterceptor reads authorization from gRPC metadata
func AuthUnaryInterceptor(jwtKey []byte) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		if values := md.Get("authorization"); len(values) > 0 {
			claims, err := authorization.ParseToken(values[0], jwtKey)
			if err == nil && claims.UserID != "" {
				ctx = authorization.WithUserID(ctx, claims.UserID)
				return handler(ctx, req)
			}
		}

		userID := authorization.NewUserID()
		token, err := authorization.CreateToken(userID, jwtKey)
		if err != nil {
			return nil, err
		}

		if err := grpc.SetHeader(ctx, metadata.Pairs("authorization", token)); err != nil {
			return nil, err
		}

		ctx = authorization.WithUserID(ctx, userID)
		return handler(ctx, req)
	}
}
