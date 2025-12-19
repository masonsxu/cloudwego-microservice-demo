package models

import (
	"errors"
	"time"
)

// models.go - 模型包的统一导出文件
// 提供所有数据模型的统一访问入口

// 导出所有模型结构体，便于其他包使用
// 这些模型与 IDL 定义保持一致，支持完整的身份管理系统

// 核心实体模型
// - UserProfile: 用户档案（自然人信息）
// - UserMembership: 用户组织成员关系（多组织支持）
// - Organization: 组织结构（层级组织架构）
// - Department: 部门
// - PostureResource: 姿态资源（用户姿态资源文件管理）

// 枚举类型
// - AccountType: 账户类型（超级管理员/管理员/用户）
// - UserStatus: 用户状态（活跃/未激活/暂停/锁定）
// - MembershipStatus: 成员关系状态（活跃/待处理/已结束）

// 工具类型
// - PageOptions: 分页查询参数
// - PageResult: 分页结果元信息

// 特点：
// 1. 无外键约束设计 - 通过应用层维护数据一致性
// 2. 灵活的架构设计 - 支持多种业务场景和扩展需求
// 3. 多组织支持 - 用户可归属多个组织，支持复杂角色关系
// 4. 完整审计跟踪 - 包含创建、更新、删除等审计字段
// 5. 版本控制支持 - UserProfile 支持乐观锁并发控制
// 6. 状态管理完善 - 支持用户状态、成员关系状态管理

// 注意事项：
// - 所有 ID 字段使用 ULID 格式（26位字符串）
// - 时间戳使用毫秒级 Unix 时间戳
// - JSON 字段用于存储数组和复杂数据结构
// - 索引优化用于提升查询性能
// - GORM 钩子用于数据验证和完整性检查

// ============================================================================
// 错误定义
// ============================================================================

// OrganizationLogo 相关错误
var (
	ErrLogoEmptyFileID         = errors.New("logo file ID cannot be empty")
	ErrLogoEmptyFileName       = errors.New("logo file name cannot be empty")
	ErrLogoInvalidFileSize     = errors.New("logo file size must be greater than 0")
	ErrLogoEmptyMimeType       = errors.New("logo MIME type cannot be empty")
	ErrLogoInvalidMimeType     = errors.New("logo MIME type must be a valid image type")
	ErrLogoInvalidUploader     = errors.New("logo uploader ID is invalid")
	ErrLogoInvalidStatus       = errors.New("logo status is invalid")
	ErrLogoMissingExpiry       = errors.New("temporary logo must have an expiry time")
	ErrLogoInvalidExpiry       = errors.New("logo expiry time must be in the future")
	ErrLogoMissingOrganization = errors.New("bound logo must have an organization ID")
)

// ============================================================================
// 工具函数
// ============================================================================

// GetCurrentTimestamp 获取当前毫秒级时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
