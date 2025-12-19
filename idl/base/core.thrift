namespace go core

// =================================================================
//                      API Gateway 基础数据类型定义
// =================================================================
typedef string UUID // 通用唯一标识符 UUID
typedef i64 TimestampMS // Unix 毫秒时间戳 (推荐)
typedef string URL // 统一资源定位符
typedef string IPAddress // IP地址

// =================================================================
//                          请求上下文
// =================================================================

/**
   * 请求上下文 - 所有请求的核心上下文信息
   * 这是整个系统中传递的核心数据结构
   * 字段顺序与现有Context中间件保持一致
   */
struct RequestContext {
    // 追踪信息 - 核心标识符
    1: optional string requestID,              // 请求唯一标识
    2: optional string traceID,                // 链路追踪标识
    3: optional TimestampMS requestTime,       // 请求发起时间 (Unix毫秒)
    // 客户端信息
    4: optional string clientIP,               // 客户端IP地址
    5: optional string userAgent,              // 用户代理字符串
    // 从认证信息中提炼出的、对所有服务都有用的核心身份
    6: optional string userID,                 // 用户ID
    7: optional string organizationID,         // 组织ID
    8: optional string tenantID,               // 租户ID
    // 其他信息
    9: optional string source,                 // 请求来源标识 (gateway|internal|cli)
    10: optional string locale,                // 语言环境 (zh-CN|en-US)
    // 扩展字段
    11: optional map<string, string> metadata, // 自定义元数据 (键值对)
}
