package redisback

import (
	"fmt"
	"gbc2/common/store"
	"time"

	goredislib "github.com/go-redis/redis/v8"
)

type (
	registryRedisBack struct {
		conf RegistryRedisBackConf
		rdb  *goredislib.Client
	}
)

const (
	svIDField   = "_ID_"
	svNameField = "_NAME_"
)

func newRegistryRedisBack(conf *RegistryRedisBackConf) *registryRedisBack {
	return &registryRedisBack{
		conf: *conf,
		rdb:  newRedisConn(conf.Redis),
	}
}

func (bk *registryRedisBack) Add(name string, fields store.ServiceFields) (id string, err error) {
	var key string
	for {
		id = snow.Generate().Base58()
		key = bk.conf.ServiceKey + id

		ok, err := bk.rdb.Exists(bk.rdb.Context(), key).Result()
		if err != nil {
			err = fmt.Errorf("check service id exist failed, %v", err)
			return "", err
		}

		if ok == 0 {
			break
		}
	}

	fields[svIDField] = id
	fields[svNameField] = name

	ctx := bk.rdb.Context()
	p := bk.rdb.Pipeline()
	// save service fields, and set expire time
	p.HSet(ctx, key, fields).Err()
	p.Expire(ctx, key, bk.conf.ServiceExpire)
	// save id to name set
	nameKey := bk.conf.ServiceByNameKey + name
	p.ZAdd(ctx, nameKey, &goredislib.Z{Score: 0, Member: id})
	_, err = p.Exec(ctx)
	if err != nil {
		err = fmt.Errorf("add service failed, %v", err)
		return "", err
	}

	return id, nil
}

func (bk *registryRedisBack) Get(id string) (name string, fields store.ServiceFields, err error) {
	key := bk.conf.ServiceKey + id
	fields, err = bk.rdb.HGetAll(bk.rdb.Context(), key).Result()
	if err != nil {
		return "", nil, fmt.Errorf("get service '%s' failed, %v", id, err)
	}
	name, ok := fields[svNameField]
	if !ok {
		return "", nil, fmt.Errorf("get service '%s' failed, not found name", id)
	}
	return name, fields, nil
}

func (bk *registryRedisBack) Del(id string) error {
	ctx := bk.rdb.Context()
	key := bk.conf.ServiceKey + id
	name, err := bk.rdb.HGet(ctx, key, svNameField).Result()
	if err != nil {
		return fmt.Errorf("get name of service '%s' failed, %v", id, err)
	}

	p := bk.rdb.Pipeline()
	p.Del(ctx, key)
	nameKey := bk.conf.ServiceByNameKey + name
	p.ZRem(ctx, nameKey, name)
	_, err = p.Exec(ctx)
	if err != nil {
		return fmt.Errorf("del service '%s' failed, %v", id, err)
	}

	return nil
}

func (bk *registryRedisBack) KeepAlive(id string, score float64, expire time.Duration) error {
	ctx := bk.rdb.Context()
	key := bk.conf.ServiceKey + id
	name, err := bk.rdb.HGet(ctx, key, svNameField).Result()
	if err != nil {
		return fmt.Errorf("get name of service '%s' failed, %v", id, err)
	}

	p := bk.rdb.Pipeline()
	p.Expire(ctx, key, expire).Err()
	nameKey := bk.conf.ServiceByNameKey + name
	p.ZAdd(ctx, nameKey, &goredislib.Z{Score: score, Member: id})
	_, err = p.Exec(ctx)
	if err != nil {
		return fmt.Errorf("keep alive service '%s' failed, %v", id, err)
	}
	return nil
}

func (bk *registryRedisBack) QueryByName(name string, limit int64) (map[string]store.ServiceFields, error) {
	nameKey := bk.conf.ServiceByNameKey + name
	ctx := bk.rdb.Context()

	svs := make(map[string]store.ServiceFields, limit)

	var start, stop int64
	stop = limit - 1
	count := limit
	for {
		ids, err := bk.rdb.ZRange(ctx, nameKey, start, stop).Result()
		if err != nil {
			return nil, fmt.Errorf("query service by name '%s' failed(1), %v", name, err)
		}
		if len(ids) == 0 {
			break
		}

		p := bk.rdb.Pipeline()
		for _, id := range ids {
			p.HGetAll(ctx, bk.conf.ServiceKey+id)
		}
		res, err := p.Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("query service by name '%s' failed(3), %v", name, err)
		}

		rems := make([]string, 0, len(ids))
		for i, r := range res {
			id := ids[i]
			if r.Err() != nil {
				rems = append(rems, bk.conf.ServiceKey+id)
				continue
			}

			mr := r.(*goredislib.StringStringMapCmd)
			if mr == nil || len(mr.Val()) == 0{
				rems = append(rems, bk.conf.ServiceKey+id)
				continue
			}

			svs[id] = mr.Val()
			count--
		}

		if len(rems) > 0 {
			// remove no longer exist service
			bk.rdb.Del(ctx, rems...)
		}

		if count <= 0 {
			break
		}

		start += limit
		stop += limit
	}

	return svs, nil
}
