package redisback

import (
	"gbc2/common/log"
	"time"

	goredislib "github.com/go-redis/redis/v8"
)

// connect to redis
func newRedisConn(c *RedisConf) *goredislib.Client {
	db := goredislib.NewClient(&goredislib.Options{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	}).WithTimeout(c.Timeout * time.Millisecond)

	if err := db.Ping(db.Context()).Err(); err != nil {
		log.L.Fatalf("ping redis failed, %v", err)
	}
	log.L.Infof("connect redis '%s:%d' ok", c.Addr, c.DB)

	return db
}
