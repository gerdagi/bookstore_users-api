package users

import (
	"fmt"

	"github.com/gerdagi/bookstore_users-api/utils/errors"
)

// dao stands for data access object
// only database access sould be in dao not anather place

var (
	usersDB = make(map[int64]*User)
)

// method da *point kullanmamızın nedeni asıl user üzinde (method da)
// güncelle yapmamız, *pointer kullnamadığımızda copy üzerinde güncelleme yapıyoruz
func (user *User) Get() *errors.RestError {
	result := usersDB[int64(user.Id)]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestError {
	currentUser := usersDB[int64(user.Id)]
	if currentUser != nil {
		if currentUser.Email == user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s allready registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d allready exists", user.Id))
	}

	// save
	usersDB[int64(user.Id)] = user
	return nil
}
