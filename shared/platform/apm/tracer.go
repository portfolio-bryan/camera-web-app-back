package apm

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/google/uuid"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName      = "github.com/ravilushqa/otelgqlgen"
	extensionName   = "OpenTelemetry"
	complexityLimit = "ComplexityLimit"
)

// Tracer is a GraphQL extension that traces GraphQL requests.
type Tracer struct {
	complexityExtensionName     string
	tracer                      oteltrace.Tracer
	requestVariablesBuilderFunc RequestVariablesBuilderFunc
	shouldCreateSpanFromFields  FieldsPredicateFunc
	spanKindSelector            SpanKindSelectorFunc
}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.OperationInterceptor
} = Tracer{}

// ExtensionName returns the extension name.
func (a Tracer) ExtensionName() string {
	return extensionName
}

// Validate checks if the extension is configured properly.
func (a Tracer) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

// InterceptResponse intercepts the incoming request.
func (a Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if !graphql.HasOperationContext(ctx) {
		return next(ctx)
	}

	opName := operationName(ctx)
	spanKind := a.spanKindSelector(opName)
	xID := xTracerID(ctx)
	ctx, span := a.tracer.Start(ctx, opName, oteltrace.WithSpanKind(spanKind), oteltrace.WithAttributes(
		XTracerIDHeader(xID),
	))
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	oc := graphql.GetOperationContext(ctx)
	span.SetAttributes(
		RequestQuery(oc.RawQuery),
	)

	span.SetAttributes(
		RequestQuery(oc.RawQuery),
	)
	complexityExtension := a.complexityExtensionName
	if complexityExtension == "" {
		complexityExtension = complexityLimit
	}
	complexityStats, ok := oc.Stats.GetExtension(complexityExtension).(*extension.ComplexityStats)
	if !ok {
		// complexity extension is not used
		complexityStats = &extension.ComplexityStats{}
	}

	if complexityStats.ComplexityLimit > 0 {
		span.SetAttributes(
			RequestComplexityLimit(int64(complexityStats.ComplexityLimit)),
			RequestOperationComplexity(int64(complexityStats.Complexity)),
		)
	}

	if a.requestVariablesBuilderFunc != nil {
		span.SetAttributes(a.requestVariablesBuilderFunc(oc.Variables)...)
	}

	resp := next(ctx)
	if resp != nil && len(resp.Errors) > 0 {
		span.SetStatus(codes.Error, resp.Errors.Error())
		span.RecordError(fmt.Errorf("graphql response errors: %v", resp.Errors.Error()))
		span.SetAttributes(ResolverErrors(resp.Errors)...)
	} else {
		span.SetStatus(codes.Ok, "Finished successfully")
	}

	return resp
}

func (a Tracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	// Validation of the X-Tracer-ID header is done in the middleware.
	opCtx := graphql.GetOperationContext(ctx)
	xTracerID := opCtx.Headers.Get(xTracerIDHeader)

	if xTracerID == "" {
		opCtx.Headers.Set(xTracerIDHeader, uuid.New().String())
	}

	return next(ctx)
}

// Middleware sets up a handler to start tracing the incoming
// requests.  The service parameter should describe the name of the
// (virtual) server handling the request. extension parameter may be empty string.
func Middleware(opts ...Option) Tracer {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	if cfg.RequestVariablesBuilder == nil {
		cfg.RequestVariablesBuilder = RequestVariables
	}
	if cfg.ShouldCreateSpanFromFields == nil {
		cfg.ShouldCreateSpanFromFields = alwaysTrue()
	}
	if cfg.SpanKindSelectorFunc == nil {
		cfg.SpanKindSelectorFunc = alwaysServer()
	}

	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(otelcontrib.Version()),
	)

	return Tracer{
		tracer:                      tracer,
		requestVariablesBuilderFunc: cfg.RequestVariablesBuilder,
		shouldCreateSpanFromFields:  cfg.ShouldCreateSpanFromFields,
		spanKindSelector:            cfg.SpanKindSelectorFunc,
	}

}

// alwaysTrue returns a FieldsPredicateFunc that always returns true.
func alwaysTrue() FieldsPredicateFunc {
	return func(_ *graphql.FieldContext) bool {
		return true
	}
}

func alwaysServer() SpanKindSelectorFunc {
	return func(_ string) oteltrace.SpanKind {
		return oteltrace.SpanKindServer
	}
}

func operationName(ctx context.Context) string {
	opContext := graphql.GetOperationContext(ctx)
	if opName := opContext.OperationName; opName != "" {
		return opName
	}
	if opContext.Operation != nil && opContext.Operation.Name != "" {
		return opContext.Operation.Name
	}
	return GetOperationName(ctx)
}

func xTracerID(ctx context.Context) string {
	opCtx := graphql.GetOperationContext(ctx)
	return opCtx.Headers.Get("X-Tracer-Id")
}

type operationNameCtxKey struct{}

// SetOperationName adds the operation name to the context so that the interceptors can use it.
// It will replace the operation name if it already exists in the context.
// example:
//
//		ctx = otelgqlgen.SetOperationName(r.Context(), "my-operation")
//	 	r = r.WithContext(ctx)
func SetOperationName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, operationNameCtxKey{}, name)
}

// GetOperationName gets the operation name from the context.
func GetOperationName(ctx context.Context) string {
	if oc, _ := ctx.Value(operationNameCtxKey{}).(string); oc != "" {
		return oc
	}
	return "nameless-operation"
}
