package repository

import (
	"github.com/myproject/api/domain/entity"
	"github.com/myproject/api/domain/value_object"
)

type IUserRepository interface {
	GetUsers() ([]*entity.User, error)
	GetUserByID(id int64) (*entity.User, error)
	GetUserByEmail(email value_object.Email) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) error
	ChangePasswordUser(user *entity.User) (*entity.User, error)
	DeleteUser(user *entity.User) error
}
