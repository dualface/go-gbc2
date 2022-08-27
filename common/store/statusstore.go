package store

import "time"

type (
	StatusStore interface {
		// 自动生成一个新 ID，并保存状态及绑定数据，然后返回新 ID
		Add(st uint8, b []byte) (id string, err error)

		// 保存特定 ID 的状态，以及绑定数据
		Save(id string, st uint8, b []byte) error

		// 载入特定 ID 的状态，以及绑定数据
		Load(id string) (st uint8, b []byte, err error)

		// 设置特定 ID 的状态
		SetStatus(id string, st uint8) error

		// 检查特定 ID 的状态
		GetStatus(id string) (uint8, error)

		// 删除特定 ID 及其绑定数据
		Del(id string) error

		// 保持 ID 活跃，避免被自动清理
		KeepAlive(id string, expire time.Duration) error
	}
)
