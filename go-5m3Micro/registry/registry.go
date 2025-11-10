package registry

import "context"

type Registrar interface {
	// Register 注册服务
	Register(ctx context.Context, sr *ServiceInstance) error
	// Deregister 注销服务
	Deregister(ctx context.Context, sr *ServiceInstance) error
}

type Discovery interface {
	// GetService 获取服务实例
	GetService(ctx context.Context, name string) ([]*ServiceInstance, error)
	// Watch 创建服务监听器
	Watch(ctx context.Context, name string) (Watcher, error)
}

type Watcher interface {
	// Next
	//		 1.第一次监听时，若服务实例列表不为空，则添加所有服务实例
	//		 2.相关服务实例配置发生变化时，再次添加所有服务实例，达到更新信息效果
	Next() ([]*ServiceInstance, error)
	// Stop 停止监听
	Stop() error
}
type ServiceInstance struct {
	//服务Id
	Id string `json:"id"`
	//服务名称
	Name string `json:"name"`
	//服务版本
	Version string `json:"version"`
	//服务地址
	Endpoints []string `json:"endpoints"`
	//服务元数据
	Metadata map[string]string `json:"metaData"`
}
