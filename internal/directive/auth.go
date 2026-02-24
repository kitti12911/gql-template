package directive

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
)

func Auth(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
	slog.InfoContext(ctx, "auth directive called")
	return next(ctx)
}
