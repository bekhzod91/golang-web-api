package repository

import (
	"github.com/myproject/api/domain/entity"
	"github.com/myproject/api/domain/value_object"
	"time"
)

type ITokenRepository interface {
	GetUserByToken(token value_object.Token) (*entity.User, error)
	CreateUserToken(token value_object.Token, expiration time.Duration, user *entity.User) error
	RevokeUserToken(token value_object.Token) error
}
