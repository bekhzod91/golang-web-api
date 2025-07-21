package storage

import (
	"database/sql"
	"github.com/redis/go-redis/v9"

	"github.com/hzmat24/api/domain/repository"
)

type IStorage interface {
	Token() repository.ITokenRepository
	User() repository.IUserRepository
	Role() repository.IRoleRepository
}

type storage struct {
	redisClient          *redis.Client
	sharedPostgresClient *sql.DB
	tenantPostgresClient *sql.DB
	tokenRepository      repository.ITokenRepository
	userRepository       repository.IUserRepository
	roleRepository       repository.IRoleRepository
}

func NewStorage(
	redisClient *redis.Client,
	sharedPostgresClient *sql.DB,
	tenantPostgresClient *sql.DB,
) IStorage {
	return storage{
		redisClient:          redisClient,
		sharedPostgresClient: sharedPostgresClient,
		tenantPostgresClient: tenantPostgresClient,
		tokenRepository:      NewTokenRepository(redisClient),
		userRepository:       NewUserRepository(tenantPostgresClient),
		roleRepository:       NewRoleRepository(tenantPostgresClient),
	}
}

func (s storage) Token() repository.ITokenRepository {
	return s.tokenRepository
}

func (s storage) User() repository.IUserRepository {
	return s.userRepository
}

func (s storage) Role() repository.IRoleRepository {
	return s.roleRepository
}
