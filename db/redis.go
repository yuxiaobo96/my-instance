package db

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

// redis配置
var (
	RedisMaxIdle   = 5
	RedisMaxActive = 15
)

const RedisIdleTimeoutSec = 24

func NewRedisPool(redisURL string) *redis.Pool {
	log.WithField("URL", redisURL).Info("NewRedisPool")
	return &redis.Pool{
		MaxIdle:     RedisMaxIdle,
		MaxActive:   RedisMaxActive,
		IdleTimeout: RedisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				log.WithError(err).Error("redis connection error")
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.WithError(err).Error("ping redis error")
				return err
			}
			return nil
		},
	}
}