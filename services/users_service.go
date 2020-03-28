package services

import (
	"microservices/bookstore_users-api/domain/users"
	"microservices/bookstore_users-api/utils"
	"microservices/bookstore_users-api/utils/date_utils"
	"microservices/bookstore_users-api/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	UpdateUser(users.User, bool) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
}

type usersService struct{}

func (s *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(user users.User, isPartial bool) (*users.User, *errors.RestErr) {
	current, err := s.GetUser(user.Id)
	if err != nil {
		return nil, err
	}
	if err := user.Validate(); err != nil {
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
	if updateErr := current.Update(); updateErr != nil {
		return nil, updateErr
	}
	return current, nil
}

func (s *usersService) GetUser(userID int64) (*users.User, *errors.RestErr) {
	user := &users.User{Id: userID}
	if err := user.Get(); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) DeleteUser(userID int64) *errors.RestErr {
	user := &users.User{Id: userID}
	return user.Delete()
}

func (s *usersService) SearchUser(status string) (users.Users, *errors.RestErr) {
	user := &users.User{}
	return user.FindByStatus(status)
}
