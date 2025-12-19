namespace go http_base

include "../../base/core.thrift"
include "../../base/enums.thrift"

// =================================================================
//                       JWT 相关数据结构
// =================================================================

/**
 * JWT 令牌声明（Claims）信息
 * 这些信息将直接嵌入到 JWT Token 中
 */
struct JWTClaimsDTO {
    1: optional string userProfileID (go.tag = "json:\"user_profile_id,omitempty\" form:\"user_profile_id\" query:\"user_profile_id\""),  // 用户ID
    2: optional string username (go.tag = "json:\"username,omitempty\" form:\"username\" query:\"username\""),                            // 用户名
    3: optional enums.UserStatus status (go.tag = "json:\"status,omitempty\" form:\"status\" query:\"status\""),                          // 用户状态
    4: optional string roleID (go.tag = "json:\"role_id,omitempty\" form:\"role_id\" query:\"role_id\""),                                 // 角色ID
    5: optional string organizationID (go.tag = "json:\"organization_id,omitempty\" form:\"organization_id\" query:\"organization_id\""), // 所属组织ID
    6: optional string departmentID (go.tag = "json:\"department_id,omitempty\" form:\"department_id\" query:\"department_id\""),         // 所属部门ID
    7: optional string permission (go.tag = "json:\"permission,omitempty\" form:\"permission\" query:\"permission\""),                    // 权限
    8: optional i64 exp (go.tag = "json:\"exp,omitempty\" form:\"exp\" query:\"exp\""),                                                   // 过期时间（Unix时间戳）
    9: optional i64 iat (go.tag = "json:\"iat,omitempty\" form:\"iat\" query:\"iat\""),                                                   // 签发时间（Unix时间戳）
}

/**
 * 令牌信息DTO
 */
struct TokenInfoDTO {
    1: optional string accessToken (go.tag = "json:\"access_token,omitempty\" form:\"access_token\" query:\"access_token\""), // 访问令牌
    2: optional i64 expiresIn (go.tag = "json:\"expires_in,omitempty\" form:\"expires_in\" query:\"expires_in\""),            // 访问令牌过期时间（秒）
    3: optional string tokenType (go.tag = "json:\"token_type,omitempty\" form:\"token_type\" query:\"token_type\""),         // 令牌类型（通常是"Bearer"）
}

// =================================================================
//                      通用请求/响应结构
// =================================================================

/**
 * API Gateway 统一响应基类
 * 针对 REST API 优化的响应结构
 */
struct BaseResponseDTO {
    1: required i32 code = 0 (go.tag = "json:\"code\""),           // 错误码，0表示成功
    2: required string message = "success" (go.tag = "json:\"message\""), // 人类可读的提示信息
}

/**
 * 通用操作状态响应DTO
 * 适用于删除、更新等只需返回操作状态的接口
 * 只保留基础响应信息以确保API响应格式统一
 */
struct OperationStatusResponseDTO {
    1: required BaseResponseDTO base_resp (go.tag = "json:\"base_resp\""), // 基础响应信息
}

/**
 * 分页请求
 * 针对 REST API 优化的分页参数
 */
struct PageRequestDTO {
    // 分页参数
    1: optional i32 page = 1 (api.query = "page", api.vd = "!isset(FetchAll) || *(FetchAll)$ || $>0", go.tag = "json:\"page,omitempty\" form:\"page\" query:\"page\""),                   // 页码，从1开始
    2: optional i32 limit = 20 (api.query = "limit", api.vd = "!isset(FetchAll) || *(FetchAll)$ || ($>0 && $<=200)", go.tag = "json:\"limit,omitempty\" form:\"limit\" query:\"limit\""), // 每页数量
    // 搜索与过滤
    3: optional string search (api.query = "search", go.tag = "json:\"search,omitempty\" form:\"search\" query:\"search\""),                                                              // 全局搜索关键字
    4: optional map<string, string> filter (api.query = "filter", go.tag = "json:\"filter,omitempty\" form:\"filter\" query:\"filter\""),                                                 // 字段级过滤
    // 排序
    5: optional string sort (api.query = "sort", go.tag = "json:\"sort,omitempty\" form:\"sort\" query:\"sort\""),                                                                        // 排序字段: field1,-field2
    // 字段选择（可选增强）
    6: optional list<string> fields (api.query = "fields", go.tag = "json:\"fields,omitempty\" form:\"fields\" query:\"fields\""),                                                        // 返回指定字段
    7: optional bool include_total (api.query = "include_total", go.tag = "json:\"include_total,omitempty\" form:\"include_total\" query:\"include_total\""),                             // 是否返回总数
    8: optional bool fetch_all (api.query = "fetch_all", go.tag = "json:\"fetch_all,omitempty\" form:\"fetch_all\" query:\"fetch_all\""),                                                 // 是否获取所有数据（不分页）
}

/**
 * 分页响应
 * 针对 REST API 优化的分页信息
 */
struct PageResponseDTO {
    1: optional i32 total (go.tag = "json:\"total,omitempty\""),             // 总记录数
    2: optional i32 page (go.tag = "json:\"page,omitempty\""),               // 当前页码
    3: optional i32 limit (go.tag = "json:\"limit,omitempty\""),             // 每页数量
    4: optional i32 total_pages (go.tag = "json:\"total_pages,omitempty\""), // 总页数
    5: optional bool has_next (go.tag = "json:\"has_next,omitempty\""),      // 是否有下一页
    6: optional bool has_prev (go.tag = "json:\"has_prev,omitempty\""),      // 是否有上一页
}
