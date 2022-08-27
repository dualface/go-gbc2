package redisback

import (
	"time"

	goredislib "github.com/go-redis/redis/v8"
)

type (
	statusRedisBack struct {
		rdb *goredislib.Client
	}
)

func newStatusRedisBack(conf *RedisConf) *statusRedisBack {
	return &statusRedisBack{
		rdb: newRedisConn(conf),
	}
}

func (bk *statusRedisBack) Add(st uint8, b []byte) (id string, err error) {
	return "", nil
}

func (bk *statusRedisBack) Save(id string, st uint8, b []byte) error {
	return nil
}

func (bk *statusRedisBack) Load(id string) (st uint8, b []byte, err error) {
	return 0, nil, nil
}

func (bk *statusRedisBack) SetStatus(id string, st uint8) error {
	return nil
}

func (bk *statusRedisBack) GetStatus(id string) (uint8, error) {
	return 0, nil
}

func (bk *statusRedisBack) Del(id string) error {
	return nil
}

func (bk *statusRedisBack) KeepAlive(id string, expire time.Duration) error {
	return nil
}
