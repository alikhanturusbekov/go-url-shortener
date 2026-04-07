package main

import (
	"fmt"
	shortenerpb "github.com/alikhanturusbekov/go-url-shortener/api/proto"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
	"go.uber.org/zap"
	"net"

	"google.golang.org/grpc"

	"github.com/alikhanturusbekov/go-url-shortener/internal/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/grpcserver"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
)

// setupGRPCServer prepares everything to start GRPC server
func setupGRPCServer(appConfig *config.Config, urlService *service.URLService) (*grpc.Server, net.Listener, error) {
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcserver.AuthUnaryInterceptor([]byte(appConfig.AuthorizationKey)),
		),
	)

	shortenerpb.RegisterShortenerServiceServer(
		grpcSrv,
		grpcserver.NewShortenerServer(urlService),
	)

	grpcListener, err := net.Listen("tcp", appConfig.GRPCAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen on gRPC address: %w", err)
	}

	return grpcSrv, grpcListener, nil
}

// startGRPCServer starts GRPC server
func startGRPCServer(grpcSrv *grpc.Server, grpcListener net.Listener) <-chan error {
	grpcErr := make(chan error, 1)

	go func() {
		logger.Log.Info("starting gRPC server", zap.String("address", grpcListener.Addr().String()))

		if err := grpcSrv.Serve(grpcListener); err != nil {
			grpcErr <- err
			return
		}

		grpcErr <- nil
	}()

	return grpcErr
}
