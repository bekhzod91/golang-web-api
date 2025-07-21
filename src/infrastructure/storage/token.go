package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/hzmat24/api/domain/entity"
	"github.com/hzmat24/api/domain/repository"
	"github.com/hzmat24/api/domain/value_object"
)

var ctx = context.Background()

type tokenRepository struct {
	redisClient *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) repository.ITokenRepository {
	return &tokenRepository{redisClient: redisClient}
}

func (r *tokenRepository) GetUserByToken(token value_object.Token) (*entity.User, error) {
	data, err := r.redisClient.Get(ctx, token.String()).Bytes()

	if err != nil {
		return nil, errors.New("token not found")
	}

	user := entity.User{}
	err = user.UnmarshalBinary(data)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *tokenRepository) CreateUserToken(token value_object.Token, expiration time.Duration, user *entity.User) error {
	data, err := user.MarshalBinary()
	if err != nil {
		return fmt.Errorf("CreateTokenToUser: %w", err)
	}

	return r.redisClient.Set(ctx, token.String(), data, expiration).Err()
}

func (r *tokenRepository) RevokeUserToken(token value_object.Token) error {
	return r.redisClient.Del(ctx, token.String()).Err()
}
