package auth

import "github.com/bperezgo/rtsp/shared/domain/dto"

type InmemoryRepository struct{}

func (a *InmemoryRepository) GetUser() dto.User {
	return dto.User{
		ID:        "1",
		CompanyID: "3",
	}
}
