package convutil

import (
	"encoding/json"
	"strings"
)

// StringSliceToJSON 将字符串切片转换为JSON字符串
// 用于处理 UserProfile.Specialties 字段
func StringSliceToJSON(slice []string) string {
	if len(slice) == 0 {
		return ""
	}

	data, err := json.Marshal(slice)
	if err != nil {
		return ""
	}

	return string(data)
}

// JSONToStringSlice 将JSON字符串转换为字符串切片
// 用于处理 UserProfile.Specialties 字段
func JSONToStringSlice(jsonStr string) []string {
	if jsonStr == "" {
		return nil
	}

	var slice []string
	if err := json.Unmarshal([]byte(jsonStr), &slice); err != nil {
		return nil
	}

	return slice
}

// ULIDSliceToJSON 将ULID切片转换为JSON字符串
// 用于处理 Department.AvailableEquipment 字段
func ULIDSliceToJSON(slice []string) string {
	if len(slice) == 0 {
		return ""
	}

	data, err := json.Marshal(slice)
	if err != nil {
		return ""
	}

	return string(data)
}

// JSONToULIDSlice 将JSON字符串转换为ULID切片
// 用于处理 Department.AvailableEquipment 字段
func JSONToULIDSlice(jsonStr string) []string {
	if jsonStr == "" {
		return nil
	}

	var slice []string
	if err := json.Unmarshal([]byte(jsonStr), &slice); err != nil {
		return nil
	}

	return slice
}

// StringPtr 返回字符串指针，用于处理可选字段
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

// StringValue 安全获取字符串指针的值
func StringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// Int32Ptr 返回int32指针，用于处理可选字段
func Int32Ptr(i int32) *int32 {
	if i == 0 {
		return nil
	}

	return &i
}

// Int32Value 安全获取int32指针的值
func Int32Value(i *int32) int32 {
	if i == nil {
		return 0
	}

	return *i
}

// Int64Ptr 返回int64指针，用于处理可选时间戳字段
func Int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}

	return &i
}

// Int64Value 安全获取int64指针的值
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}

// BoolPtr 返回bool指针，用于处理可选布尔字段
func BoolPtr(b bool) *bool {
	return &b
}

// BoolValue 安全获取bool指针的值
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}

// TrimSpace 安全处理字符串去空格
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmptyString 检查字符串是否为空（包括只有空格的情况）
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
