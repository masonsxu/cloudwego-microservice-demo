package converter

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/enum"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/user"
)

// 包级别工厂函数 - 用于独立使用特定转换器
//
// 使用场景：
//   - 工具脚本：数据迁移、批量转换等场景
//   - 独立测试：只需测试某个转换器，无需完整的 Converter 实例
//   - 轻量级函数：避免创建整个 Converter 聚合
//
// 注意事项：
//   - 如果业务逻辑层已有 Converter 实例，应优先使用 conv.XXX() 而非这些工厂函数
//   - 工厂函数返回的是独立实例，不会共享依赖（如 enumMapper）
//   - 适合一次性使用，不建议缓存返回值（缓存应使用 Converter 聚合实例）

// NewStandaloneUserProfileConverter 创建独立的用户档案转换器
//
// 注意：UserProfileConverter 依赖 EnumMapper，此工厂函数会创建新的 EnumMapper 实例
//
// 使用示例：
//
//	// 在测试中使用
//	func TestUserProfileConversion(t *testing.T) {
//	    conv := converter.NewStandaloneUserProfileConverter()
//	    result := conv.ModelUserProfileToThrift(userModel)
//	    assert.NotNil(t, result)
//	}
//
// 返回值：
//   - user.UserProfileConverter: 独立的用户档案转换器实例
func NewStandaloneUserProfileConverter() user.Converter {
	enumConverter := enum.NewConverter()
	return user.NewConverter(enumConverter)
}

// NewStandaloneEnumMapper 创建独立的枚举映射器
//
// 使用示例：
//
//	// 在独立函数中使用
//	func ConvertUserStatus(status models.UserStatus) identity_srv.UserStatus {
//	    converter := converter.NewStandaloneEnumConverter()
//	    return converter.ModelUserStatusToThrift(status)
//	}
//
// 返回值：
//   - enum.EnumMapper: 独立的枚举映射器实例
func NewStandaloneEnumMapper() enum.Converter {
	return enum.NewConverter()
}

// NewStandaloneBaseConverter 创建独立的基础转换器
//
// 使用示例：
//
//	// 在工具函数中使用
//	func FormatTimestamp(t time.Time) int64 {
//	    conv := converter.NewStandaloneBaseConverter()
//	    return conv.TimeToTimestampMS(t)
//	}
//
// 返回值：
//   - base.BaseConverter: 独立的基础转换器实例
func NewStandaloneBaseConverter() base.Converter {
	return base.NewConverter()
}
