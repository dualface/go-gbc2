package store

type (
	// 存储接口
	Store interface {
		// 产生一个 int64 ID，并用 base58 进行编码为字符串
		GenID() string

		// 产生一个 int64 ID
		GenIDInt() int64

		// 返回 RegistryStore 接口
		Registry() RegistryStore

		// 返回 StatusStore 接口
		Status() StatusStore
	}
)
