package services

import (
	"github.com/gerdagi/bookstore_users-api/domain/users"
	"github.com/gerdagi/bookstore_users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestError) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(userId int64) (*users.User, *errors.RestError) {
	if userId <= 0 {
		return nil, errors.NewBadRequestError("user id should be bigger than 0")
	}

	result := &users.User{
		Id: int32(userId),
	}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}
