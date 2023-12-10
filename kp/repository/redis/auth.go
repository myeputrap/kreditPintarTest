package redis

import (
	"context"
	"goKreditPintar/domain"

	"github.com/go-redis/redis/v8"
)

type redisAuthRepository struct {
	Conn *redis.Client
}

func NewRedisAuthRepository(Conn *redis.Client) domain.AuthRedisRepository {
	return &redisAuthRepository{Conn}
}

func (r *redisAuthRepository) SetRedis(ctx context.Context, key string, value string) (err error) {
	_, err = r.Conn.Set(ctx, key, value, 0).Result()
	return err
}

func (r *redisAuthRepository) GetRedis(ctx context.Context, key string) (value string, err error) {
	value, err = r.Conn.Get(ctx, key).Result()
	return
}
