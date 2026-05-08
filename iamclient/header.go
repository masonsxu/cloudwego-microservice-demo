package iamclient

// HTTP header 名称常量（业务侧契约，提案 §5.2）。
//
// 这些 header 由网关 jwt_middleware 验签通过后注入下游请求；
// 业务侧统一从 header 读身份，禁止再次解析 JWT。
const (
	HeaderUserID    = "X-User-Id"
	HeaderUserName  = "X-User-Name"
	HeaderTenantID  = "X-Tenant-Id"
	HeaderUserRoles = "X-User-Roles"
	HeaderRequestID = "X-Request-Id"
	HeaderJTI       = "X-Auth-Token-Jti"
)

// Kitex metainfo 持久键（透传到下游 RPC 服务）。
//
// 网关 RPC 客户端在调用下游前，把 HTTP header 中的身份字段写入 metainfo
// 持久值；下游 Kitex 服务通过 metainfo.GetPersistentValue 读取。
//
// 命名规则：与 HTTP header 镜像，但全小写（metainfo key 大小写敏感）。
const (
	MetaUserID    = "x-user-id"
	MetaUserName  = "x-user-name"
	MetaTenantID  = "x-tenant-id"
	MetaUserRoles = "x-user-roles"
	MetaRequestID = "request_id" // 与 gateway/internal/infrastructure/errors/trace.go 保持一致
	MetaJTI       = "x-auth-token-jti"
)
