package place

import (
	"github.com/bperezgo/rtsp/internal/domain/dto"
	"github.com/bperezgo/rtsp/shared/domain/valueobject"
)

type PlaceName struct {
	value string
}

func NewPlaceName(value string) (PlaceName, error) {
	return PlaceName{value: value}, nil
}

type Place struct {
	id     valueobject.ID
	name   PlaceName
	userID valueobject.ID
}

func New(id string, name string, userID string) (Place, error) {
	voID, err := valueobject.NewID(id)
	if err != nil {
		return Place{}, err
	}

	voName, err := NewPlaceName(name)
	if err != nil {
		return Place{}, err
	}

	voUserID, err := valueobject.NewID(userID)
	if err != nil {
		return Place{}, err
	}

	return Place{id: voID, name: voName, userID: voUserID}, nil
}

func (p Place) ToDTO() dto.Place {
	return dto.Place{
		ID:     p.id.Value(),
		Name:   p.name.value,
		UserID: p.userID.Value(),
	}
}
