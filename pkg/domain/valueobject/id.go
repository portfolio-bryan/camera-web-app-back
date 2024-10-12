package valueobject

import (
	"fmt"

	"github.com/bperezgo/rtsp/pkg/domain/errors"
	"github.com/google/uuid"
)

var (
	ErrInvalidIDCode errors.ErrorCode = "invalid_id"
)

type ErrInvalidID struct {
	id string
}

func NewErrInvalidID(id string) ErrInvalidID {
	return ErrInvalidID{
		id: id,
	}
}

func (e *ErrInvalidID) Code() errors.ErrorCode {
	return ErrInvalidIDCode
}

func (e *ErrInvalidID) Error() string {
	return fmt.Sprintf("the id '%s' is not valid", e.id)
}

type ID struct {
	value string
}

func NewID(value string) (ID, error) {
	err := uuid.Validate(value)

	if err != nil {
	}
	return ID{
		value: value,
	}, nil
}
