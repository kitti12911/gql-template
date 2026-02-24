package server

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type OtelTracer struct {
	tracer oteltrace.Tracer
}

func NewOtelTracer() *OtelTracer {
	return &OtelTracer{
		tracer: otel.Tracer("graphql"),
	}
}

var _ graphql.HandlerExtension = &OtelTracer{}
var _ graphql.OperationInterceptor = &OtelTracer{}
var _ graphql.ResponseInterceptor = &OtelTracer{}
var _ graphql.FieldInterceptor = &OtelTracer{}

func (t *OtelTracer) ExtensionName() string {
	return "OpenTelemetryTracer"
}

func (t *OtelTracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (t *OtelTracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)

	operationType := "query"
	if oc.Operation != nil && oc.Operation.Operation != "" {
		operationType = string(oc.Operation.Operation)
	}

	operationName := "anonymous"
	if oc.OperationName != "" {
		operationName = oc.OperationName
	}

	spanName := fmt.Sprintf("graphql.%s %s", operationType, operationName)
	ctx, span := t.tracer.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		attribute.String("graphql.operation.type", operationType),
		attribute.String("graphql.operation.name", operationName),
	)

	return next(ctx)
}

func (t *OtelTracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	span := oteltrace.SpanFromContext(ctx)
	response := next(ctx)

	if response != nil && len(response.Errors) > 0 {
		span.SetAttributes(attribute.Int("graphql.errors.count", len(response.Errors)))

		for i, err := range response.Errors {
			span.AddEvent(fmt.Sprintf("graphql.error.%d", i), oteltrace.WithAttributes(
				attribute.String("error.message", err.Message),
				attribute.String("error.path", fmt.Sprintf("%v", err.Path)),
			))
		}
	}

	return response
}

func (t *OtelTracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	if !fc.IsResolver {
		return next(ctx)
	}

	spanName := fmt.Sprintf("%s.%s", fc.Object, fc.Field.Name)
	ctx, span := t.tracer.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindInternal))
	defer span.End()

	span.SetAttributes(
		attribute.String("graphql.field.object", fc.Object),
		attribute.String("graphql.field.name", fc.Field.Name),
	)

	start := time.Now()
	res, err := next(ctx)
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("graphql.field.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.SetAttributes(
			attribute.String("error.message", err.Error()),
			attribute.Bool("graphql.field.error", true),
		)
	}

	return res, err
}
