package graph

import (
	"github.com/bperezgo/rtsp/internal/app/observability"
	"github.com/bperezgo/rtsp/internal/app/places"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	observabilityProvider observability.Tracer
	placesService         places.Service
}

func NewResolver(
	placesService places.Service,
	observabilityProvider observability.Tracer,
) *Resolver {
	return &Resolver{
		placesService:         placesService,
		observabilityProvider: observabilityProvider,
	}
}
