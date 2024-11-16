package main

import (
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bperezgo/rtsp/config"
	"github.com/bperezgo/rtsp/graph"
	"github.com/bperezgo/rtsp/internal/app/observability"
	"github.com/bperezgo/rtsp/internal/app/places"
	"github.com/bperezgo/rtsp/shared/platform/apm"
	"github.com/bperezgo/rtsp/shared/platform/middlewares"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// Defining the Graphql handler
func graphqlHandler(tracerProvider trace.TracerProvider) gin.HandlerFunc {
	getPlacesService := places.NewService()
	observabilityProvider := observability.New(tracerProvider)

	resolver := graph.NewResolver(
		getPlacesService, observabilityProvider,
	)
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	server.SetErrorPresenter(middlewares.ErrorPresenter)
	// server.Use(middlewares.XTracer)
	// graphql.HandlerExtension
	server.Use(apm.Middleware(apm.WithTracerProvider(tracerProvider)))

	return func(c *gin.Context) {
		server.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal("error loading .env file", err)
	}

	ctx := context.Background()

	c := config.GetConfig()
	port := c.ServerPort

	// wrappedHandler := otelhttp.NewHandler(handler, "hello")

	honeycombProvider := apm.NewHoneycombTracerProvider(ctx, apm.HoneycombOptions{
		Name: c.Otel.ServiceName,
	})

	defer honeycombProvider.Shutdown(ctx)

	r := gin.Default()
	// r.Use(middlewares.Tracer())
	r.Use(middlewares.Cors())
	r.Use(middlewares.GinContextToContextMiddleware())
	r.Use(middlewares.MetadataMiddleware())
	r.Use(middlewares.Logging())
	r.POST("/query", graphqlHandler(honeycombProvider.TracerProvider()))
	r.GET("/", playgroundHandler())
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("error running server", err)
	}
}
