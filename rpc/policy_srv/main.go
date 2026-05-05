package main

import (
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv/policyservice"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/wire"
)

func main() {
	container, cleanup, err := wire.InitializeApp()
	if err != nil {
		klog.Fatalf("初始化应用失败: %v", err)
	}
	defer cleanup()

	container.StartHealthCheck()

	serviceImpl := NewPolicyServiceImpl(container.Decision, container.Enforcer)
	svr := policyservice.NewServer(serviceImpl, container.ServerOptions.Options...)

	klog.Infof("Policy service 启动于 %s", container.ServerOptions.Addr.String())
	if err := svr.Run(); err != nil {
		klog.Fatalf("服务异常退出: %v", err)
	}
}
