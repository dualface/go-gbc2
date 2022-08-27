package store

import "time"

type (
	ServiceFields = map[string]string

	RegistryStore interface {
		// register service with auto generated id
		Add(name string, fields ServiceFields) (id string, err error)

		// get service by id
		Get(id string) (name string, fields ServiceFields, err error)

		// delete service
		Del(id string) error

		// keep alive service, update score and expire time
		KeepAlive(id string, score float64, expire time.Duration) error

		// query services by name
		QueryByName(name string, limit int64) (map[string]ServiceFields, error)
	}
)
