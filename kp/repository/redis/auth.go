package redis

import (
	"context"
	"errors"
	"goKreditPintar/domain"
	"time"

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

func (r *redisAuthRepository) HSetRedis(ctx context.Context, key string, fields map[string]interface{}, expTime int64) error {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}

	pipe := r.Conn.Pipeline()
	defer pipe.Close()

	pipe.HSet(ctx, key, args...)
	if expTime > 0 {
		expiration := time.Duration(expTime) * time.Second
		pipe.Expire(ctx, key, expiration)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisAuthRepository) HDelRedis(ctx context.Context, key string) error {
	client := r.Conn.WithContext(ctx)

	_, err := client.Del(ctx, key).Result()
	return err
}

func (r *redisAuthRepository) HExistsRedis(ctx context.Context, key string) (bool, error) {
	client := r.Conn.WithContext(ctx)

	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *redisAuthRepository) HGetRedis(ctx context.Context, key string) (respon domain.Session, err error) {
	data := r.Conn.HGetAll(ctx, key)
	res, err := data.Result()
	if err != nil {
		return
	}

	if len(res) == 0 {
		err = errors.New("not found")
		return
	}

	err = data.Scan(&respon)
	if err != nil {
		return
	}

	return
}
