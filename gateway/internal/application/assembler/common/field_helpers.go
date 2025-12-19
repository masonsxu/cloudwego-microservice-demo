package common

// field_helpers.go 提供网关层对象转换的通用工具函数
// 用于 HTTP DTO 与 RPC Request/Response 之间的字段安全转换

// ========================================================================
// 字段安全复制工具（HTTP -> RPC / RPC -> HTTP）
// ========================================================================

// CopyStringPtr 安全复制字符串指针字段
// 用于 RPC -> HTTP 转换，将 RPC 的 optional string 复制到 HTTP DTO
func CopyStringPtr(src *string) *string {
	if src == nil {
		return nil
	}

	return src
}

// CopyInt32Ptr 安全复制 int32 指针字段
// 用于 RPC -> HTTP 转换
func CopyInt32Ptr(src *int32) *int32 {
	if src == nil {
		return nil
	}

	return src
}

// CopyInt64Ptr 安全复制 int64 指针字段（通常用于时间戳）
// 用于 RPC -> HTTP 转换
func CopyInt64Ptr(src *int64) *int64 {
	if src == nil {
		return nil
	}

	return src
}

// CopyBoolPtr 安全复制 bool 指针字段
// 用于 RPC -> HTTP 转换
func CopyBoolPtr(src *bool) *bool {
	if src == nil {
		return nil
	}

	return src
}

// CopyStringSlice 安全复制字符串切片
// 如果源切片为空，返回 nil
func CopyStringSlice(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	// 创建新切片避免共享底层数组
	dest := make([]string, len(src))
	copy(dest, src)

	return dest
}

// ========================================================================
// 条件字段设置工具（HTTP -> RPC）
// ========================================================================

// SetIfNotNil 仅在源指针非 nil 时调用 setter 函数
// 用于将 HTTP DTO 的 optional 字段设置到 RPC Request
//
// 示例:
//
//	var req identity_srv.UpdateUserRequest
//	SetIfNotNil(dto.Email, func(v *string) { req.Email = v })
func SetIfNotNil[T any](src *T, setter func(*T)) {
	if src != nil {
		setter(src)
	}
}

// SetIfNotEmpty 仅在字符串指针非 nil 且非空时调用 setter
// 用于过滤空字符串
func SetIfNotEmpty(src *string, setter func(*string)) {
	if src != nil && *src != "" {
		setter(src)
	}
}

// SetIfSliceNotEmpty 仅在切片非空时调用 setter
// 用于避免设置空切片
func SetIfSliceNotEmpty[T any](src []T, setter func([]T)) {
	if len(src) > 0 {
		setter(src)
	}
}

// ========================================================================
// Thrift IsSet 检查辅助（HTTP -> RPC）
// ========================================================================

// ApplyIfSet 根据 Thrift 生成的 IsSetXxx() 方法决定是否应用字段
// 这是处理 Thrift optional 字段的标准模式
//
// 示例:
//
//	ApplyIfSet(dto.IsSetEmail, dto.Email, func(v *string) { req.Email = v })
func ApplyIfSet[T any](isSet func() bool, value *T, setter func(*T)) {
	if isSet() {
		setter(value)
	}
}

// ApplyIfSetSlice 用于切片类型的 IsSet 检查
// 示例:
//
//	ApplyIfSetSlice(dto.IsSetSpecialties, dto.Specialties, func(v []string) { req.Specialties = v })
func ApplyIfSetSlice[T any](isSet func() bool, value []T, setter func([]T)) {
	if isSet() {
		setter(value)
	}
}

// ========================================================================
// 复合字段处理工具
// ========================================================================

// Int32Value 安全获取 int32 指针的值，nil 返回默认值 0
func Int32Value(ptr *int32) int32 {
	if ptr == nil {
		return 0
	}

	return *ptr
}

// Int32Ptr 创建 int32 指针，0 值返回 nil（避免传递无意义的 0）
func Int32Ptr(val int32) *int32 {
	if val == 0 {
		return nil
	}

	return &val
}

// BoolPtr 创建 bool 指针（不过滤 false 值）
func BoolPtr(val bool) *bool {
	return &val
}

// StringPtr 创建字符串指针，空字符串返回 nil
func StringPtr(val string) *string {
	if val == "" {
		return nil
	}

	return &val
}

// StringValue 安全获取字符串指针的值，nil 返回空字符串
func StringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}
