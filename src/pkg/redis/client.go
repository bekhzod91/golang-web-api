package redis

import (
	"context"
	"fmt"
	"github.com/myproject/api/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(config config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:      addr,
		Username:  config.RedisUser,
		Password:  config.RedisPassword,
		DB:        config.RedisDatabase,
		TLSConfig: config.RedisTLSConfig,
	})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		return nil, fmt.Errorf("failed to connection redis: %w", err)
	}

	return client, nil
}
