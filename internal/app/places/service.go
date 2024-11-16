package places

import (
	"context"

	"github.com/bperezgo/rtsp/internal/domain/aggregates/place"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s *Service) GetPlaces(ctx context.Context) ([]place.Place, error) {
	place1, err := place.New("1", "Place 1", "1")

	if err != nil {
		return []place.Place{}, err
	}

	return []place.Place{
		place1,
	}, nil
}
