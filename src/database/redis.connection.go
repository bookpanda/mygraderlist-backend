package database

import (
	"errors"

	"github.com/bookpanda/mygraderlist-backend/src/config"
	"github.com/go-redis/redis/v8"
)

func InitRedisConnect(conf *config.Redis) (cache *redis.Client, err error) {
	cache = redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		DB:       3,
		Username: "",
		Password: conf.Password,
	})

	if cache == nil {
		return nil, errors.New("Cannot connect to redis server")
	}

	return
}
