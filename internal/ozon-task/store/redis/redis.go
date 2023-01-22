package redis

import (
	"context"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	"github.com/go-redis/redis/v8"
)

const (
	storeName = "redis"
)

type Redis struct {
	rdb *redis.Client
	*links
}

func (r *Redis) Close() error {
	return r.rdb.Close()
}

func (r *Redis) Name() string {
	return storeName
}

func Provide(conf *config.Redis) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password},
	)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Redis{
		rdb:   rdb,
		links: &links{rdb: rdb},
	}, nil
}
