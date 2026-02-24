package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	libconfig "github.com/kitti12911/lib-util/config"
	"github.com/kitti12911/lib-util/logger"
	"github.com/kitti12911/lib-util/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	examplev1 "gql-template/gen/grpc/example/v1"
	"gql-template/graph"
	"gql-template/internal/config"
	"gql-template/internal/server"
)

func main() {
	cfg, err := libconfig.Load[config.Config]("config.yml")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}

	opts := []logger.Option{
		logger.WithLevel(cfg.Logging.Level),
	}

	if cfg.Logging.ServiceName != "" {
		opts = append(opts, logger.WithServiceName(cfg.Logging.ServiceName))
	}

	if cfg.Logging.AddSource {
		opts = append(opts, logger.WithSource())
	}

	if cfg.Logging.EnableTrace {
		opts = append(opts, logger.WithTrace())
	}

	logger.New(opts...)
	ctx := context.Background()

	if cfg.CollectorEndpoint != "" {
		tp, err := tracing.New(ctx, cfg.ServiceName, fmt.Sprintf("%s:%d", cfg.CollectorEndpoint, cfg.CollectorPort))
		if err != nil {
			slog.Error("failed to init tracing", "error", err)
			os.Exit(1)
		}

		defer tracing.Shutdown(ctx, tp)
	}

	exampleConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.ExampleService.Host, cfg.ExampleService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("failed to connect to example service", "error", err)
		os.Exit(1)
	}
	defer exampleConn.Close()

	slog.Info("connected to example service", "addr", fmt.Sprintf("%s:%d", cfg.ExampleService.Host, cfg.ExampleService.Port))

	resolver := &graph.Resolver{
		ExampleClient: examplev1.NewExampleServiceClient(exampleConn),
	}

	srv := server.NewHTTPServer(cfg.Port, resolver)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
	defer cancel()

	srv.Stop(shutdownCtx)

	slog.Info("server stopped")
}
