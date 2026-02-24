package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"gql-template/graph"
	"gql-template/internal/directive"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(port int, resolver *graph.Resolver) *HTTPServer {
	cfg := graph.Config{
		Resolvers: resolver,
	}

	cfg.Directives.Auth = directive.Auth

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(cfg))
	srv.Use(NewOtelTracer())

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	mux.Handle("/query", srv)

	return &HTTPServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

func (s *HTTPServer) Start() error {
	slog.Info("HTTP server listening", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		slog.WarnContext(ctx, "HTTP server shutdown error", "error", err)
	}
}
