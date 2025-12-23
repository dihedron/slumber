package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// API is the command to start the API server.
type API struct {
	// GRPCPort is the port for the gRPC server.
	GRPCPort int `short:"g" long:"grpc-port" description:"The port for the gRPC server" default:"9090"`
	// RESTPort is the port for the REST gateway.
	RESTPort int `short:"r" long:"rest-port" description:"The port for the REST gateway" default:"8080"`
}

// Execute is the main entry point for the api command.
func (cmd *API) Execute(args []string) error {
	slog.Info("starting API server", "grpc_port", cmd.GRPCPort, "rest_port", cmd.RESTPort)

	osClient, err := NewOpenStackClient()
	if err != nil {
		return fmt.Errorf("failed to create OpenStack client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Start gRPC server
	grpcServer := grpc.NewServer()
	RegisterVMControlServer(grpcServer, NewServer(osClient))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cmd.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port: %w", err)
	}

	go func() {
		slog.Info("gRPC server listening", "address", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	// 2. Start REST gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = RegisterVMControlHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", cmd.GRPCPort), opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	restAddr := fmt.Sprintf(":%d", cmd.RESTPort)
	restServer := &http.Server{
		Addr:    restAddr,
		Handler: mux,
	}

	go func() {
		slog.Info("REST gateway listening", "address", restAddr)
		if err := restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("REST gateway failed", "error", err)
		}
	}()

	// Wait for interruption
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	slog.Info("shutting down API server")
	grpcServer.GracefulStop()
	restServer.Shutdown(ctx)

	return nil
}
