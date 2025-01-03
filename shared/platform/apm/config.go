package apm

import (
	"github.com/99designs/gqlgen/graphql"
	sharedob "github.com/bperezgo/rtsp/shared/domain/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type FieldsPredicateFunc func(ctx *graphql.FieldContext) bool

type SpanKindSelectorFunc func(operationName string) trace.SpanKind

// config is used to configure the mongo tracer.
type config struct {
	TracerProvider             sharedob.TracerProvider
	ComplexityExtensionName    string
	RequestVariablesBuilder    RequestVariablesBuilderFunc
	ShouldCreateSpanFromFields FieldsPredicateFunc
}

// RequestVariablesBuilderFunc is the signature of the function
// used to build the request variables attributes.
type RequestVariablesBuilderFunc func(requestVariables map[string]interface{}) []attribute.KeyValue

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider sharedob.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.TracerProvider = provider
	})
}

// WithComplexityExtensionName specifies complexity extension name.
func WithComplexityExtensionName(complexityExtensionName string) Option {
	return optionFunc(func(cfg *config) {
		cfg.ComplexityExtensionName = complexityExtensionName
	})
}

// WithRequestVariablesAttributesBuilder allows specifying a custom function
// to handle the building of the attributes for the variables.
func WithRequestVariablesAttributesBuilder(builder RequestVariablesBuilderFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.RequestVariablesBuilder = builder
	})
}

// WithoutVariables allows disabling the variables attributes.
func WithoutVariables() Option {
	return optionFunc(func(cfg *config) {
		cfg.RequestVariablesBuilder = func(_ map[string]interface{}) []attribute.KeyValue {
			return nil
		}
	})
}

// WithCreateSpanFromFields allows specifying a custom function
// to handle the creation or not of spans regarding the GraphQL context fields.
func WithCreateSpanFromFields(predicate FieldsPredicateFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.ShouldCreateSpanFromFields = predicate
	})
}
