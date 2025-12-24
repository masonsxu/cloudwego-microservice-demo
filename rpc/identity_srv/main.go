package main

import (
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv/identityservice"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/wire"
)

func main() {
	// 统一初始化所有依赖（只初始化一次）
	// Wire 自动管理依赖图和生命周期
	container, cleanup, err := wire.InitializeApp()
	if err != nil {
		klog.Fatalf("failed to init app: %v", err)
	}
	defer cleanup()

	// 启动健康检查服务器（独立的 goroutine）
	container.StartHealthCheck()

	// 创建 Handler 实例
	serviceImpl := NewIdentityServiceImpl(container)

	// 创建并启动 Kitex Server
	svr := identityservice.NewServer(serviceImpl, container.ServerOptions.Options...)

	klog.Infof("Identity service starting on %s", container.ServerOptions.Addr.String())

	if err := svr.Run(); err != nil {
		klog.Fatalf("server stopped with error: %v", err)
	}
}
