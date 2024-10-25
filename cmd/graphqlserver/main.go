package main

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bperezgo/rtsp/config"
	"github.com/bperezgo/rtsp/graph"
	"github.com/bperezgo/rtsp/shared/platform/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	resolver := graph.NewResolver()
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	server.SetErrorPresenter(middlewares.ErrorPresenter)
	// server.Use(apm.Middleware())

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

	c := config.GetConfig()
	port := c.ServerPort

	// use otelconfig to set up OpenTelemetry SDK
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry()
	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}
	defer otelShutdown()

	otel.Tracer("rtsp")

	// wrappedHandler := otelhttp.NewHandler(handler, "hello")

	r := gin.Default()
	// r.Use(middlewares.Tracer())
	r.Use(middlewares.Cors())
	r.Use(middlewares.GinContextToContextMiddleware())
	r.Use(middlewares.MetadataMiddleware())
	r.Use(middlewares.Logging())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("error running server", err)
	}
}
