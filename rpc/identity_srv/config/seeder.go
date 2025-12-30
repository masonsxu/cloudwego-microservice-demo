package config

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/password"
)

// SeedDatabase 初始化数据库种子数据
// 该函数是幂等的，可以安全地重复执行
func SeedDatabase(db *gorm.DB, logger *zerolog.Logger, cfg *DatabaseConfig) error {
	logger.Info().Msg("开始数据库种子数据初始化...")

	// 1. 创建默认组织
	orgID, err := seedDefaultOrganization(db, logger)
	if err != nil {
		logger.Error().Err(err).Msg("创建默认组织失败")
		return fmt.Errorf("创建默认组织失败: %w", err)
	}

	// 2. 创建超级管理员用户
	userID, err := seedSuperAdminUser(db, logger)
	if err != nil {
		logger.Error().Err(err).Msg("创建超级管理员用户失败")
		return fmt.Errorf("创建超级管理员用户失败: %w", err)
	}

	// 1. 创建系统角色定义
	if err := seedSystemRoles(db); err != nil {
		log.Printf("创建系统角色定义失败: %v", err)
		return fmt.Errorf("创建系统角色定义失败: %w", err)
	}

	// 2. 分配超级管理员角色
	if err := seedSuperAdminRoleAssignment(db, cfg); err != nil {
		log.Printf("分配超级管理员角色失败: %v", err)
		return fmt.Errorf("分配超级管理员角色失败: %w", err)
	}

	logger.Info().
		Str("default_org_id", orgID.String()).
		Str("superadmin_user_id", userID.String()).
		Msg("数据库种子数据初始化完成")

	return nil
}

// seedDefaultOrganization 创建默认组织
// 使用 code="DEFAULT" 作为唯一标识，实现幂等性
func seedDefaultOrganization(db *gorm.DB, logger *zerolog.Logger) (uuid.UUID, error) {
	logger.Info().Msg("正在创建或验证默认组织...")

	org := &models.Organization{}

	// 使用 FirstOrCreate 实现幂等性
	// Attrs: 仅在创建时设置的属性
	result := db.Where("code = ?", "DEFAULT").
		Attrs(&models.Organization{
			Name:                "默认组织",
			ParentID:            uuid.Nil, // 根组织
			FacilityType:        "Default Facility",
			AccreditationStatus: "N/A",
			ProvinceCity:        models.StringSlice{},
		}).
		FirstOrCreate(org, &models.Organization{
			Code: "DEFAULT",
		})

	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("创建默认组织失败: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		logger.Info().
			Str("org_id", org.ID.String()).
			Str("code", org.Code).
			Msg("✅ 默认组织创建成功")
	} else {
		logger.Info().
			Str("org_id", org.ID.String()).
			Str("code", org.Code).
			Msg("ℹ️  默认组织已存在，跳过创建")
	}

	return org.ID, nil
}

// handleExistingSuperAdminUser 处理已存在的超级管理员用户
// 包括恢复软删除的用户和更新系统用户标记
func handleExistingSuperAdminUser(db *gorm.DB, logger *zerolog.Logger, user *models.UserProfile) (uuid.UUID, error) {
	// 如果是软删除的用户，恢复它
	if user.DeletedAt.Valid {
		if err := db.Unscoped().Model(user).Updates(map[string]interface{}{
			"deleted_at":     nil,
			"is_system_user": true, // 确保标记为系统用户
		}).Error; err != nil {
			return uuid.Nil, fmt.Errorf("恢复软删除的超级管理员用户失败: %w", err)
		}

		logger.Info().
			Str("user_id", user.ID.String()).
			Str("username", user.Username).
			Bool("is_system_user", true).
			Msg("✅ 超级管理员用户已恢复（从软删除状态）")

		return user.ID, nil
	}

	// 确保现有用户标记为系统用户
	if !user.IsSystemUser {
		if err := db.Model(user).Update("is_system_user", true).Error; err != nil {
			logger.Warn().Err(err).Msg("更新系统用户标记失败")
		} else {
			logger.Info().
				Str("user_id", user.ID.String()).
				Str("username", user.Username).
				Msg("✅ 已更新 superadmin 用户的系统用户标记")
		}
	}

	logger.Info().
		Str("user_id", user.ID.String()).
		Str("username", user.Username).
		Bool("is_system_user", user.IsSystemUser).
		Msg("ℹ️  超级管理员用户已存在，跳过创建")

	return user.ID, nil
}

// seedSuperAdminUser 创建超级管理员用户
// 使用 username="superadmin" 作为唯一标识，实现幂等性
// 支持软删除记录的恢复
func seedSuperAdminUser(db *gorm.DB, logger *zerolog.Logger) (uuid.UUID, error) {
	logger.Info().Msg("正在创建或验证超级管理员用户...")

	// 生成默认密码哈希
	defaultPassword := "password123"

	passwordHash, err := password.HashPassword(defaultPassword)
	if err != nil {
		return uuid.Nil, fmt.Errorf("生成密码哈希失败: %w", err)
	}

	user := &models.UserProfile{}

	// 1. 先检查是否存在（包括软删除的记录）
	err = db.Unscoped().Where("username = ?", "superadmin").First(user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// 查询出错
			return uuid.Nil, fmt.Errorf("查询超级管理员用户失败: %w", err)
		}
		// 用户不存在，继续创建流程
	} else {
		// 用户已存在，处理并返回
		return handleExistingSuperAdminUser(db, logger, user)
	}

	// 2. 用户不存在，创建新用户
	user = &models.UserProfile{
		Username:           "superadmin",
		PasswordHash:       passwordHash,
		Email:              "superadmin@masonsxu.local",
		RealName:           "超级管理员",
		Status:             models.UserStatusActive,
		MustChangePassword: false,
		IsSystemUser:       true, // 标记为系统用户
	}

	if err := db.Create(user).Error; err != nil {
		return uuid.Nil, fmt.Errorf("创建超级管理员用户失败: %w", err)
	}

	logger.Info().
		Str("user_id", user.ID.String()).
		Str("username", user.Username).
		Bool("is_system_user", user.IsSystemUser).
		Str("default_password", defaultPassword).
		Msg("✅ 超级管理员用户创建成功")
	logger.Warn().
		Str("username", "superadmin").
		Str("password", defaultPassword).
		Msg("⚠️  请及时修改默认密码！")

	return user.ID, nil
}

// systemRoleDefinition 系统角色定义结构
type systemRoleDefinition struct {
	Name        string
	Description string
}

// seedSystemRoles 创建系统角色定义
// 使用 name 作为唯一标识，实现幂等性
func seedSystemRoles(db *gorm.DB) error {
	log.Println("正在创建或验证系统角色定义...")

	// 定义 18 个系统角色
	roles := []systemRoleDefinition{
		// 系统管理角色
		{"superadmin", "超级管理员 - 拥有系统全部权限，负责系统维护和最高级别管理"},
	}

	createdCount := 0
	skippedCount := 0

	for _, roleDef := range roles {
		role := &models.RoleDefinition{}

		// 使用 FirstOrCreate 实现幂等性
		// Attrs: 仅在创建时设置的属性（避免 JSON 类型比较）
		result := db.Where("name = ?", roleDef.Name).
			Attrs(&models.RoleDefinition{
				Description:  roleDef.Description,
				Status:       models.RoleStatusActive,
				Permissions:  models.Permissions{},
				IsSystemRole: true,
			}).
			FirstOrCreate(role, &models.RoleDefinition{
				Name: roleDef.Name,
			})

		if result.Error != nil {
			return fmt.Errorf("创建角色 %s 失败: %w", roleDef.Name, result.Error)
		}

		if result.RowsAffected > 0 {
			createdCount++

			log.Printf("✅ 角色创建成功: %s (ID: %s)", roleDef.Name, role.ID)
		} else {
			skippedCount++
		}
	}

	log.Printf("系统角色定义初始化完成 - 创建: %d, 跳过: %d", createdCount, skippedCount)

	return nil
}

// seedSuperAdminRoleAssignment 分配超级管理员角色
// 从 identity_srv 数据库获取超级管理员用户 ID
func seedSuperAdminRoleAssignment(db *gorm.DB, cfg *DatabaseConfig) error {
	log.Println("正在分配超级管理员角色...")

	// 1. 获取超级管理员角色 ID
	var superadminRole models.RoleDefinition
	if err := db.Where("name = ?", "superadmin").First(&superadminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("超级管理员角色不存在，请先确保角色初始化完成")
		}

		return fmt.Errorf("查询超级管理员角色失败: %w", err)
	}

	log.Printf("✅ 超级管理员角色已找到 (ID: %s)", superadminRole.ID)

	// 2. 从 identity_srv 数据库获取超级管理员用户 ID
	// 创建到 identity_srv 数据库的临时连接
	superadminUserID, err := querySuperAdminUserID(cfg)
	if err != nil {
		log.Printf("⚠️  查询超级管理员用户失败: %v", err)
		log.Println("提示：请确保 identity_srv 服务已正确初始化超级管理员用户")

		return nil // 不返回错误，允许服务启动
	}

	if superadminUserID == uuid.Nil {
		log.Println("⚠️  超级管理员用户不存在，跳过角色分配")
		return nil
	}

	log.Printf("✅ 从 identity_srv 成功获取超级管理员用户 ID: %s", superadminUserID)

	// 3. 检查角色分配是否已存在
	var existingAssignment models.UserRoleAssignment

	err = db.Where("user_id = ? AND role_id = ?", superadminUserID, superadminRole.ID).
		First(&existingAssignment).Error
	if err == nil {
		log.Printf("ℹ️  超级管理员角色分配已存在，跳过创建 (ID: %s)", existingAssignment.ID)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("检查角色分配失败: %w", err)
	}

	// 4. 创建角色分配
	assignment := &models.UserRoleAssignment{
		UserID: superadminUserID,
		RoleID: superadminRole.ID,
	}

	if err := db.Create(assignment).Error; err != nil {
		return fmt.Errorf("创建角色分配失败: %w", err)
	}

	log.Printf("✅ 成功为超级管理员分配角色")
	log.Printf("   用户 ID: %s", superadminUserID)
	log.Printf("   角色 ID: %s", superadminRole.ID)

	return nil
}

// querySuperAdminUserID 从 identity_srv 数据库查询超级管理员用户 ID
// 建立临时数据库连接来跨数据库查询
func querySuperAdminUserID(cfg *DatabaseConfig) (uuid.UUID, error) {
	// 构建 identity_srv 数据库连接 DSN
	identityDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.DBName, // 连接到 identity_srv 数据库
		cfg.Port,
		cfg.SSLMode,
		cfg.Timezone,
	)

	// 建立临时连接
	identityDB, err := gorm.Open(postgres.Open(identityDSN), &gorm.Config{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("连接 identity_srv 数据库失败: %w", err)
	}

	// 确保连接关闭
	sqlDB, err := identityDB.DB()
	if err != nil {
		return uuid.Nil, fmt.Errorf("获取 identity_srv 数据库连接失败: %w", err)
	}
	defer sqlDB.Close()

	// 查询超级管理员用户 ID
	var result struct {
		ID uuid.UUID `gorm:"column:id"`
	}

	err = identityDB.Table("user_profiles").
		Select("id").
		Where("username = ? AND deleted_at IS NULL", "superadmin").
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, nil // 用户不存在，返回 Nil UUID
		}

		return uuid.Nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return result.ID, nil
}
