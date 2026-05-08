// Package iamclient 是业务系统接入 IAM 的统一 SDK。
//
// 提供两大能力：
//  1. Subject 还原：从网关注入的 HTTP header 或 Kitex metadata 还原调用主体身份
//     （UserID / Tenant / Roles / Jti / RequestID）。
//  2. 权限决策：通过 PDP 服务（rpc/policy_srv）做单点鉴权，内置 LRU 缓存。
//
// 业务侧禁止再次解析 JWT、禁止持有 Casbin Enforcer、禁止直连 identity_srv 数据库。
//
// 典型用法：
//
//	cli, err := iamclient.New(iamclient.Config{
//	    EtcdEndpoints: []string{"etcd:2379"},
//	    PolicyService: "policy-service",
//	    CallerService: "my-business-service",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer cli.Close()
//
//	// HTTP 业务（Hertz）
//	subject, err := cli.SubjectFromHeader(c.Request.Header)
//	if err != nil {
//	    c.AbortWithStatus(http.StatusUnauthorized)
//	    return
//	}
//	if err := subject.MustCheck(ctx, "read", "patient:"+id); err != nil {
//	    c.AbortWithError(http.StatusForbidden, err)
//	    return
//	}
//
//	// Kitex 业务
//	subject, err := cli.SubjectFromContext(ctx)
//	// 同上
//
// 设计参考：docs/04-权限管理/重构提案-网关边界与权限模型.md §5.6。
package iamclient
