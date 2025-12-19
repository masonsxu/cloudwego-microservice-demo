namespace go rpc_base

// =================================================================
//                      通用请求/响应结构
// =================================================================

/**
 * 分页请求
 */
struct PageRequest {
    // 分页参数
    1: optional i32 page = 1,               // 页码，从1开始
    2: optional i32 limit = 20,             // 每页数量
    // 搜索与过滤
    3: optional string search,              // 全局搜索关键字
    4: optional map<string, string> filter, // 字段级过滤
    // 排序
    5: optional string sort,                // 排序字段: field1,-field2
    // 字段选择（可选增强）
    6: optional list<string> fields,        // 返回指定字段
    7: optional bool include_total,         // 是否返回总数
    8: optional bool fetch_all,             // 是否获取所有数据（不分页）
}

/**
 * 分页响应
 */
struct PageResponse {
    1: optional i32 total,       // 总记录数
    2: optional i32 page,        // 当前页码
    3: optional i32 limit,       // 每页数量
    4: optional i32 total_pages, // 总页数
    5: optional bool has_next,   // 是否有下一页
    6: optional bool has_prev,   // 是否有上一页
}