package redisback

import (
	"context"
	"gbc2/common/log"
	"gbc2/common/store"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	goredislib "github.com/go-redis/redis/v8"
)

type (
	// redis backend
	RedisBack struct {
	}

	// configuration of redis backend
	RedisBackConf struct {
		Registry *RegistryRedisBackConf `yaml:"registry"`
		Status   *RedisConf             `yaml:"status"`
		SnowNode *SnowNodeConf          `yaml:"snow_node"`
	}

	// configuration of redis
	RedisConf struct {
		Addr     string        `yaml:"addr"`
		Username string        `yaml:"username"`
		Password string        `yaml:"password"`
		DB       int           `yaml:"db"`
		Timeout  time.Duration `yaml:"timeout"`
	}

	// configuration of registry backend
	RegistryRedisBackConf struct {
		Redis *RedisConf `yaml:"redis"`

		ServiceExpire    time.Duration `yaml:"service_expire"`
		ServiceKey       string        `yaml:"service_key"`
		ServiceByNameKey string        `yaml:"service_by_name_key"`
	}

	// configuration of snowflake node
	SnowNodeConf struct {
		KeyPrefix string        `yaml:"key_prefix"`
		Expire    time.Duration `yaml:"expire"`
	}
)

var (
	// configuration of redis backend
	conf RedisBackConf

	// instance of redis backend
	back RedisBack

	// instance of snowflake node
	snow *snowflake.Node

	// instance of registry redis backend
	registry *registryRedisBack

	// instance of status redis backend
	status *statusRedisBack
)

// setup redis backend
func SetupStoreRedisBackend(cf *RedisBackConf) {
	conf = *cf
	snow = newSnowflake(conf.SnowNode, context.Background(), registry.rdb)
	registry = newRegistryRedisBack(conf.Registry)
	status = newStatusRedisBack(conf.Status)
}

func (bk *RedisBack) GenID() string {
	return snow.Generate().Base58()
}

func (bk *RedisBack) GenIDInt() int64 {
	return snow.Generate().Int64()
}

func (bk *RedisBack) Registry() store.RegistryStore {
	return registry
}

func (bk *RedisBack) Status() store.StatusStore {
	return status
}

// create snowflake node
func newSnowflake(conf *SnowNodeConf, stop context.Context, rdb *goredislib.Client) *snowflake.Node {
	nodeMax := 2 ^ int64(snowflake.NodeBits)
	var id int64
	var key string

	for {
		id = time.Now().UnixMicro() % nodeMax
		key = conf.KeyPrefix + strconv.FormatInt(id, 10)
		ok, err := rdb.Exists(rdb.Context(), key).Result()
		if err != nil {
			log.L.Fatalf("gen snowflake node id failed(1), %v", err)
		}
		if ok == 0 {
			break
		}
	}

	expire := conf.Expire * time.Millisecond
	rdb.Set(rdb.Context(), key, id, expire)

	node, err := snowflake.NewNode(id)
	if err != nil {
		log.L.Fatalf("gen snowflake node id failed(2), %v", err)
	}

	// update expire time of node id
	go func() {
		ticker := time.NewTicker(expire / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rdb.Expire(rdb.Context(), key, expire)
			case <-stop.Done():
				return
			}
		}
	}()

	log.L.Infof("create snowflake node '%d' success", id)
	return node
}
