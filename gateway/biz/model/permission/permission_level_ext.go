// permission_level_ext.go
// 扩展 thriftgo 生成的 PermissionLevel 枚举，添加自定义 JSON 序列化
// 确保输出前端期望的小写字符串格式: "none" / "read" / "write" / "full"

package permission

import (
	"encoding/json"
	"fmt"
)

// permissionLevelToString 将枚举值转换为前端期望的小写字符串
var permissionLevelToString = map[PermissionLevel]string{
	PermissionLevel_NONE:  "none",
	PermissionLevel_READ:  "read",
	PermissionLevel_WRITE: "write",
	PermissionLevel_FULL:  "full",
}

// stringToPermissionLevel 将前端字符串转换为枚举值
var stringToPermissionLevel = map[string]PermissionLevel{
	"none":  PermissionLevel_NONE,
	"read":  PermissionLevel_READ,
	"write": PermissionLevel_WRITE,
	"full":  PermissionLevel_FULL,
	// 兼容大写格式
	"NONE":  PermissionLevel_NONE,
	"READ":  PermissionLevel_READ,
	"WRITE": PermissionLevel_WRITE,
	"FULL":  PermissionLevel_FULL,
}

// MarshalJSON 自定义 JSON 序列化，输出小写字符串
func (p PermissionLevel) MarshalJSON() ([]byte, error) {
	if str, ok := permissionLevelToString[p]; ok {
		return json.Marshal(str)
	}

	return json.Marshal("none")
}

// UnmarshalJSON 自定义 JSON 反序列化，支持小写和大写字符串
func (p *PermissionLevel) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		// 尝试作为数字解析（向后兼容）
		var num int64
		if numErr := json.Unmarshal(data, &num); numErr != nil {
			return fmt.Errorf("PermissionLevel must be a string or number: %w", err)
		}

		*p = PermissionLevel(num)

		return nil
	}

	if level, ok := stringToPermissionLevel[str]; ok {
		*p = level
		return nil
	}

	return fmt.Errorf("invalid PermissionLevel value: %s", str)
}

// ToLowerString 返回前端期望的小写字符串表示
func (p PermissionLevel) ToLowerString() string {
	if str, ok := permissionLevelToString[p]; ok {
		return str
	}

	return "none"
}
