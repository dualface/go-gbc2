package store

import "time"

type (
	// 定义一个服务
	Service struct {
		ID         string `json:"id"`                    // 服务唯一 ID
		Name       string `json:"name"`                  // 同一个名字可以有多个服务实例
		URL        string `json:"url,omitempty"`         // 对外提供服务的 URL
		CmdSubject string `json:"cmd_subject,omitempty"` // 服务订阅的命令通道
		ReqSubject string `json:"req_subject,omitempty"` // 服务订阅的请求通道
	}

	// 服务注册表接口
	RegistryStore interface {
		// 注册服务
		Add(sv *Service) error

		// 取得特定 ID 的服务
		Get(id string) (sv *Service, err error)

		// 删除特定服务
		Del(sv *Service) error

		// 保持服务活跃，并更新服务的排序
		KeepAlive(sv *Service, score int64, expire time.Duration) error

		// 按名字查询服务，按 score 从小到大排序
		QueryByName(name string, limit int) ([]*Service, error)
	}
)
