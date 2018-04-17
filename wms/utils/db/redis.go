package db

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool     *redis.Pool
	redisServer   string
	redisPassword string
	maxIdle       int
	maxActive     int
	idleTimeout   int
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if redisPassword != "" {
				fmt.Println(redisPassword)
				if _, err = c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func InitRedisPool() {
	redisServer = beego.AppConfig.String("redis::redis_server")
	redisPassword = beego.AppConfig.String("redis::redis_pwd")
	maxIdle = 3
	maxActive = 5
	idleTimeout = 240

	RedisPool = newPool(redisServer, redisPassword)

	return
}

func String(reply interface{}, err1 error) (value string, err2 error) {
	value, err2 = redis.String(reply, err1)
	return
}

func Bool(reply interface{}, err1 error) (value bool, err2 error) {
	value, err2 = redis.Bool(reply, err1)
	return
}
