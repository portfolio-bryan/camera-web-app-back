package apm

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/bperezgo/rtsp/shared/domain/observability"
	"github.com/google/uuid"

	"go.opentelemetry.io/otel/codes"
)

const (
	extensionName   = "OpenTelemetry"
	complexityLimit = "ComplexityLimit"
)

// Tracer is a GraphQL extension that traces GraphQL requests.
type Tracer struct {
	complexityExtensionName     string
	tracer                      observability.Tracer
	requestVariablesBuilderFunc RequestVariablesBuilderFunc
	shouldCreateSpanFromFields  FieldsPredicateFunc
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

	// TODO: Define logic to define the OperationName with SetOperationName function
	opName := operationName(ctx)
	ctx, span := a.tracer.Start(ctx, opName)
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	oc := graphql.GetOperationContext(ctx)
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
	xTracerID := opCtx.Headers.Get(observability.XTracerIDHeader)

	if xTracerID == "" {
		xTracerID = uuid.New().String()
		opCtx.Headers.Set(observability.XTracerIDHeader, xTracerID)
	}

	ctx = context.WithValue(ctx, observability.XTracerIDCtxKey, xTracerID)

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
		panic("TracerProvider is required")
	}
	if cfg.RequestVariablesBuilder == nil {
		cfg.RequestVariablesBuilder = RequestVariables
	}
	if cfg.ShouldCreateSpanFromFields == nil {
		cfg.ShouldCreateSpanFromFields = alwaysTrue()
	}

	tracer := cfg.TracerProvider.Tracer()

	return Tracer{
		tracer:                      tracer,
		requestVariablesBuilderFunc: cfg.RequestVariablesBuilder,
		shouldCreateSpanFromFields:  cfg.ShouldCreateSpanFromFields,
	}

}

// alwaysTrue returns a FieldsPredicateFunc that always returns true.
func alwaysTrue() FieldsPredicateFunc {
	return func(_ *graphql.FieldContext) bool {
		return true
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
