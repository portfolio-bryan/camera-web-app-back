package user

import "github.com/bperezgo/rtsp/shared/domain/valueobject"

type User struct {
	id    valueobject.ID
	email valueobject.Email
}

func NewUser(id, email string) (User, error) {
	voID, err := valueobject.NewID(id)

	if err != nil {
		return User{}, err
	}

	voEmail, err := valueobject.NewEmail(email)

	if err != nil {
		return User{}, err
	}

	return User{
		id:    voID,
		email: voEmail,
	}, nil
}
