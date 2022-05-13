package services

import (
	"github.com/gerdagi/bookstore_users-api/domain/users"
	cryptoutils "github.com/gerdagi/bookstore_users-api/utils/crypto_utils"
	date_utils "github.com/gerdagi/bookstore_users-api/utils/date_utils"
	resterrors "github.com/gerdagi/bookstore_utils-go/rest_errors"
	// "github.com/gerdagi/bookstore_utils-go/resterrors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *resterrors.RestError)
	GetUser(int64) (*users.User, *resterrors.RestError)
	LoginUser(users.LoginRequest) (*users.User, *resterrors.RestError)
	DeleteUser(int64) *resterrors.RestError
	UpdateUser(bool, users.User) (*users.User, *resterrors.RestError)
	SearchUser(string) ([]users.User, *resterrors.RestError)
}

func (s *usersService) SearchUser(status string) ([]users.User, *resterrors.RestError) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) DeleteUser(userId int64) *resterrors.RestError {
	current := &users.User{Id: int32(userId)}
	return current.Delete()
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *resterrors.RestError) {
	current, err := s.GetUser(int64(user.Id))
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, *resterrors.RestError) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Status = users.StatusActive
	user.Password = cryptoutils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, *resterrors.RestError) {
	dao := &users.User{
		Email:    request.Email,
		Password: cryptoutils.GetMd5(request.Password),
	}

	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	return dao, nil
}

func (s *usersService) GetUser(userId int64) (*users.User, *resterrors.RestError) {
	if userId <= 0 {
		return nil, resterrors.NewBadRequestError("user id should be bigger than 0")
	}

	result := &users.User{
		Id: int32(userId),
	}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}
