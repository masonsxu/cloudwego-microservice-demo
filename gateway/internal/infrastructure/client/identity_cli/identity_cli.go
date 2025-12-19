package identitycli

import "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv/identityservice"

// IdentityClient 聚合所有用户相关服务的统一接口
type IdentityClient interface {
	identityservice.Client
}
