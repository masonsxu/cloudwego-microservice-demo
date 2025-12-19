package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	rustfsclient "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/rustfs_client"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl 组织管理业务逻辑实现
type LogicImpl struct {
	dal               dal.DAL
	converter         converter.Converter
	logoStorageClient rustfsclient.LogoStorageClient
}

// NewLogic 创建组织管理业务逻辑实例
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
	logoStorageClient rustfsclient.LogoStorageClient,
) OrganizationLogic {
	return &LogicImpl{
		dal:               dal,
		converter:         converter,
		logoStorageClient: logoStorageClient,
	}
}

// ============================================================================
// 组织基础操作
// ============================================================================

// CreateOrganization 创建组织
func (l *LogicImpl) CreateOrganization(
	ctx context.Context,
	req *identity_srv.CreateOrganizationRequest,
) (*identity_srv.Organization, error) {
	// 参数验证
	if err := l.validateCreateOrganizationRequest(req); err != nil {
		return nil, err
	}

	// 转换请求为模型
	org := l.converter.Organization().CreateRequestToModel(req)

	// 在事务中创建组织
	var result *models.Organization

	err := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		err := txDAL.Organization().Create(ctx, org)
		if err != nil {
			return err
		}

		result = org

		// 如果提供了logoID，绑定Logo到组织
		if req.LogoID != nil && *req.LogoID != "" {
			logoID, parseErr := uuid.Parse(*req.LogoID)
			if parseErr != nil {
				return errno.ErrInvalidParams.WithMessage("无效的LogoID格式")
			}

			// 验证Logo是否存在且为临时状态
			logo, getErr := txDAL.Logo().GetByID(ctx, logoID.String())
			if getErr != nil {
				return errno.ErrLogoNotFound.WithMessage("Logo不存在或已被删除")
			}

			// 验证Logo状态
			if logo.Status != models.LogoStatusTemporary {
				return errno.ErrLogoAlreadyBound.WithMessage("Logo已被绑定或已删除")
			}

			// 验证Logo未过期
			if logo.IsExpired() {
				return errno.ErrLogoExpired.WithMessage("Logo已过期，无法绑定")
			}

			// 绑定Logo到组织
			if bindErr := txDAL.Logo().BindToOrganization(ctx, logoID, org.ID); bindErr != nil {
				return errno.ErrLogoBindingFailed.WithMessage(fmt.Sprintf("绑定Logo失败: %v", bindErr))
			}
		}

		return nil
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, errno.ErrParentOrganizationNotFound
		}

		return nil, err // 已经是errno类型，直接返回
	}

	// 转换为Thrift对象
	thriftOrg := l.converter.Organization().ModelToThrift(result)

	// 查询组织绑定的Logo并生成URL
	logo, err := l.getOrganizationLogo(ctx, result.ID)
	if err != nil {
		// 记录警告但不中断流程
		// TODO: 添加日志记录
	} else if logo != nil {
		// 生成Logo下载URL
		logoURL, urlErr := l.generateLogoURL(ctx, logo.FileID)
		if urlErr == nil && logoURL != "" {
			thriftOrg.Logo = &logoURL
		}

		// 设置Logo ID
		logoIDStr := logo.ID.String()
		thriftOrg.LogoID = &logoIDStr
	}

	return thriftOrg, nil
}

// GetOrganization 根据ID获取组织信息
func (l *LogicImpl) GetOrganization(
	ctx context.Context,
	req *identity_srv.GetOrganizationRequest,
) (*identity_srv.Organization, error) {
	if req.OrganizationID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("组织ID不能为空")
	}

	org, err := l.dal.Organization().GetByID(ctx, *req.OrganizationID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrOrganizationNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取组织信息失败: " + err.Error())
	}

	// 转换为Thrift对象
	thriftOrg := l.converter.Organization().ModelToThrift(org)

	// 查询组织绑定的Logo并生成URL
	logo, err := l.getOrganizationLogo(ctx, org.ID)
	if err != nil {
		// 记录警告但不中断流程
		// TODO: 添加日志记录
	} else if logo != nil {
		// 生成Logo下载URL
		logoURL, urlErr := l.generateLogoURL(ctx, logo.FileID)
		if urlErr == nil && logoURL != "" {
			thriftOrg.Logo = &logoURL
		}

		// 设置Logo ID
		logoIDStr := logo.ID.String()
		thriftOrg.LogoID = &logoIDStr
	}

	return thriftOrg, nil
}

// UpdateOrganization 更新组织信息
func (l *LogicImpl) UpdateOrganization(
	ctx context.Context,
	req *identity_srv.UpdateOrganizationRequest,
) (*identity_srv.Organization, error) {
	if req.OrganizationID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("组织ID不能为空")
	}

	// 验证更新请求
	if err := l.validateUpdateOrganizationRequest(req); err != nil {
		return nil, err
	}

	// 获取现有组织
	existingOrg, err := l.dal.Organization().GetByID(ctx, *req.OrganizationID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrOrganizationNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取组织信息失败: " + err.Error())
	}

	// 应用更新
	updatedOrg := l.converter.Organization().ApplyUpdateToModel(existingOrg, req)

	// 在事务中更新
	var result *models.Organization

	txErr := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		// 处理Logo更新逻辑
		if req.LogoID != nil && *req.LogoID != "" {
			newLogoID, parseErr := uuid.Parse(*req.LogoID)
			if parseErr != nil {
				return errno.ErrInvalidParams.WithMessage("无效的LogoID格式")
			}

			// 1. 查询组织是否已有绑定的Logo
			oldLogo, getOldErr := txDAL.Logo().GetByOrganizationID(ctx, existingOrg.ID)
			if getOldErr != nil && !errno.IsRecordNotFound(getOldErr) {
				return errno.WrapDatabaseError(getOldErr, "查询旧Logo失败")
			}

			// 2. 如果存在旧Logo，删除旧Logo（S3文件+数据库记录）
			if oldLogo != nil {
				// 删除S3文件
				if l.logoStorageClient != nil {
					if deleteErr := l.logoStorageClient.DeleteLogo(ctx, oldLogo.FileID); deleteErr != nil {
						// 记录警告但继续（S3文件删除失败不应阻止Logo更新）
						// TODO: 添加日志记录
					}
				}

				// 软删除数据库记录
				if deleteErr := txDAL.Logo().Delete(ctx, oldLogo.ID.String()); deleteErr != nil {
					return errno.WrapDatabaseError(deleteErr, "删除旧Logo记录失败")
				}
			}

			// 3. 验证新Logo是否存在且为临时状态
			newLogo, getNewErr := txDAL.Logo().GetByID(ctx, newLogoID.String())
			if getNewErr != nil {
				return errno.ErrLogoNotFound.WithMessage("新Logo不存在或已被删除")
			}

			// 验证Logo状态
			if newLogo.Status != models.LogoStatusTemporary {
				return errno.ErrLogoAlreadyBound.WithMessage("新Logo已被绑定或已删除")
			}

			// 验证Logo未过期
			if newLogo.IsExpired() {
				return errno.ErrLogoExpired.WithMessage("新Logo已过期，无法绑定")
			}

			// 4. 绑定新Logo到组织
			if bindErr := txDAL.Logo().BindToOrganization(ctx, newLogoID, existingOrg.ID); bindErr != nil {
				return errno.ErrLogoBindingFailed.WithMessage(fmt.Sprintf("绑定新Logo失败: %v", bindErr))
			}
		}

		// 更新组织信息
		if err := txDAL.Organization().Update(ctx, updatedOrg); err != nil {
			return err
		}

		result = updatedOrg

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	// 转换为Thrift对象
	thriftOrg := l.converter.Organization().ModelToThrift(result)

	// 查询组织绑定的Logo并生成URL
	logo, err := l.getOrganizationLogo(ctx, result.ID)
	if err != nil {
		// 记录警告但不中断流程
		// TODO: 添加日志记录
	} else if logo != nil {
		// 生成Logo下载URL
		logoURL, urlErr := l.generateLogoURL(ctx, logo.FileID)
		if urlErr == nil && logoURL != "" {
			thriftOrg.Logo = &logoURL
		}

		// 设置Logo ID
		logoIDStr := logo.ID.String()
		thriftOrg.LogoID = &logoIDStr
	}

	return thriftOrg, nil
}

// DeleteOrganization 删除组织（软删除）
func (l *LogicImpl) DeleteOrganization(
	ctx context.Context,
	organizationID string,
) error {
	if organizationID == "" {
		return errno.ErrInvalidParams.WithMessage("组织ID不能为空")
	}

	// 软删除组织
	err := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.Organization().SoftDelete(ctx, organizationID)
	})
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("删除组织失败: " + err.Error())
	}

	return nil
}

// ListOrganizations 分页查询组织列表
func (l *LogicImpl) ListOrganizations(
	ctx context.Context,
	req *identity_srv.ListOrganizationsRequest,
) (*identity_srv.ListOrganizationsResponse, error) {
	// 转换分页参数（包含 page, limit, sort, search, filter, fetchAll 等所有参数）
	opts := l.converter.Base().PageRequestToQueryOptions(req.Page)

	// 添加业务特定的 ParentID 过滤条件
	if req.ParentID != nil {
		opts.WithFilter("parent_id", req.ParentID)
	}

	// 执行查询
	organizations, pageResult, err := l.dal.Organization().FindAll(ctx, opts)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询组织列表失败: " + err.Error())
	}

	// 转换结果并为每个组织生成logo URL
	orgDTOs := make([]*identity_srv.Organization, len(organizations))
	for i, org := range organizations {
		thriftOrg := l.converter.Organization().ModelToThrift(org)

		// 查询组织绑定的Logo并生成URL
		logo, err := l.getOrganizationLogo(ctx, org.ID)
		if err != nil {
			// 记录警告但不中断流程
			// TODO: 添加日志记录
		} else if logo != nil {
			// 生成Logo下载URL
			logoURL, urlErr := l.generateLogoURL(ctx, logo.FileID)
			if urlErr == nil && logoURL != "" {
				thriftOrg.Logo = &logoURL
			}

			// 设置Logo ID
			logoIDStr := logo.ID.String()
			thriftOrg.LogoID = &logoIDStr
		}

		orgDTOs[i] = thriftOrg
	}

	return &identity_srv.ListOrganizationsResponse{
		Organizations: orgDTOs,
		Page:          l.converter.Base().PageResponseToThrift(pageResult),
	}, nil
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// validateCreateOrganizationRequest 验证创建组织请求
func (l *LogicImpl) validateCreateOrganizationRequest(
	req *identity_srv.CreateOrganizationRequest,
) error {
	if req.Name == nil {
		return errno.ErrInvalidParams.WithMessage("组织名称不能为空")
	}

	// 验证省市信息格式（如果提供了）
	if len(req.ProvinceCity) > 0 {
		// 验证省市信息列表中每个项目的长度
		for i, city := range req.ProvinceCity {
			if len(city) > 100 {
				return errno.ErrInvalidParams.WithMessage(fmt.Sprintf("第%d个省市信息长度不能超过100字符", i+1))
			}

			if city == "" {
				return errno.ErrInvalidParams.WithMessage("省市信息不能为空")
			}
		}
		// 限制省市信息列表的项目数量
		if len(req.ProvinceCity) > 10 {
			return errno.ErrInvalidParams.WithMessage("省市信息列表项目数量不能超过10个")
		}
	}

	return nil
}

// validateUpdateOrganizationRequest 验证更新组织请求
func (l *LogicImpl) validateUpdateOrganizationRequest(
	req *identity_srv.UpdateOrganizationRequest,
) error {
	// 验证省市信息格式（如果提供了）
	if len(req.ProvinceCity) > 0 {
		// 验证省市信息列表中每个项目的长度
		for i, city := range req.ProvinceCity {
			if len(city) > 100 {
				return errno.ErrInvalidParams.WithMessage(fmt.Sprintf("第%d个省市信息长度不能超过100字符", i+1))
			}

			if city == "" {
				return errno.ErrInvalidParams.WithMessage("省市信息不能为空")
			}
		}
		// 限制省市信息列表的项目数量
		if len(req.ProvinceCity) > 10 {
			return errno.ErrInvalidParams.WithMessage("省市信息列表项目数量不能超过10个")
		}
	}

	return nil
}

// ============================================================================
// Logo URL 生成辅助方法
// ============================================================================

// getOrganizationLogo 根据组织ID查询绑定的Logo信息
// 返回Logo对象（包含ID和FileID），如果不存在或查询失败返回nil
func (l *LogicImpl) getOrganizationLogo(
	ctx context.Context,
	organizationID uuid.UUID,
) (*models.OrganizationLogo, error) {
	// 查询组织绑定的Logo
	logo, err := l.dal.Logo().GetByOrganizationID(ctx, organizationID)
	if err != nil {
		// Logo不存在是正常情况，不应报错
		if errno.IsRecordNotFound(err) {
			return nil, nil
		}
		// 其他错误记录日志但不中断业务
		return nil, nil
	}

	return logo, nil
}

// generateLogoURL 根据fileID生成logo的访问URL
// 返回7天有效期的预签名URL
func (l *LogicImpl) generateLogoURL(
	ctx context.Context,
	fileID string,
) (string, error) {
	if fileID == "" {
		return "", nil
	}

	// 如果 logoStorageClient 未初始化，返回原始 fileID
	if l.logoStorageClient == nil {
		return fileID, nil
	}

	// 使用LogoStorageClient生成7天长期URL
	logoURL, err := l.logoStorageClient.GetLogoURL(ctx, fileID, 7*24*3600)
	if err != nil {
		return "", fmt.Errorf("failed to generate logo URL: %w", err)
	}

	return logoURL, nil
}
