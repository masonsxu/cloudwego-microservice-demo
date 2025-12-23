package logo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/logo"
	rustfsclient "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/rustfs_client"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl Logo业务逻辑实现
type LogicImpl struct {
	repo              logo.LogoRepository
	converter         converter.Converter
	logoStorageClient rustfsclient.LogoStorageClient
}

// NewLogic 创建Logo业务逻辑实例
func NewLogic(
	repo logo.LogoRepository,
	converter converter.Converter,
	logoStorageClient rustfsclient.LogoStorageClient,
) LogoLogic {
	return &LogicImpl{
		repo:              repo,
		converter:         converter,
		logoStorageClient: logoStorageClient,
	}
}

// UploadTemporaryLogo 上传临时Logo（7天过期）
func (l *LogicImpl) UploadTemporaryLogo(
	ctx context.Context,
	req *identity_srv.UploadTemporaryLogoRequest,
) (*identity_srv.OrganizationLogo, error) {
	// 1. 参数验证
	if req.FileName == nil || *req.FileName == "" {
		return nil, errno.ErrInvalidParams.WithMessage("文件名不能为空")
	}

	if len(req.FileContent) == 0 {
		return nil, errno.ErrInvalidParams.WithMessage("文件内容不能为空")
	}

	if req.MimeType == nil || *req.MimeType == "" {
		return nil, errno.ErrInvalidParams.WithMessage("文件MIME类型不能为空")
	}

	if req.UploadedBy == nil || *req.UploadedBy == "" {
		return nil, errno.ErrInvalidParams.WithMessage("上传者ID不能为空")
	}

	// 2. 验证上传者ID格式
	uploaderID, err := uuid.Parse(*req.UploadedBy)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的上传者ID格式")
	}

	// 3. 验证文件大小
	fileSize := int64(len(req.FileContent))
	if err := l.logoStorageClient.ValidateFileSize(fileSize); err != nil {
		return nil, errno.ErrFileSizeExceeded.WithMessage(err.Error())
	}

	// 4. 验证文件MIME类型
	if err := l.logoStorageClient.ValidateFileType(*req.MimeType); err != nil {
		return nil, errno.ErrInvalidFileType.WithMessage(err.Error())
	}

	// 5. 生成唯一LogoID
	logoID := uuid.New()

	// 6. 上传文件到S3（自动添加Status=temporary标签）
	fileID, err := l.logoStorageClient.UploadTemporaryLogo(
		ctx,
		uploaderID,
		*req.FileName,
		req.FileContent,
		*req.MimeType,
	)
	if err != nil {
		return nil, errno.ErrFileUploadFailed.WithMessage(fmt.Sprintf("文件上传失败: %v", err))
	}

	// 7. 转换请求为Model并设置S3信息
	logoModel := l.converter.Logo().UploadRequestToModel(req)
	logoModel.ID = logoID
	logoModel.FileID = fileID
	logoModel.UploadedBy = uploaderID
	logoModel.FileSize = fileSize

	// 8. 保存元数据到数据库
	if err := l.repo.Create(ctx, logoModel); err != nil {
		// 如果数据库保存失败，尝试删除已上传的文件
		_ = l.logoStorageClient.DeleteLogo(ctx, fileID)
		return nil, errno.WrapDatabaseError(err, "保存Logo元数据失败")
	}

	// 9. 转换为Thrift响应
	response := l.converter.Logo().ModelToThrift(logoModel)

	// 10. 生成下载URL（临时Logo使用长期预签名URL，7天过期）
	downloadURL, err := l.generateDownloadURL(ctx, fileID, 7*24*3600) // 7天
	if err == nil {
		response.DownloadUrl = &downloadURL
	}

	return response, nil
}

// GetOrganizationLogo 获取Logo详情
func (l *LogicImpl) GetOrganizationLogo(
	ctx context.Context,
	req *identity_srv.GetOrganizationLogoRequest,
) (*identity_srv.OrganizationLogo, error) {
	// 1. 参数验证
	if req.LogoID == nil || *req.LogoID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("LogoID不能为空")
	}

	// 2. 解析LogoID
	logoID, err := uuid.Parse(*req.LogoID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的LogoID格式")
	}

	// 3. 从数据库查询Logo
	logoModel, err := l.repo.GetByID(ctx, logoID.String())
	if err != nil {
		return nil, errno.ErrLogoNotFound.WithMessage("Logo不存在")
	}

	// 4. 检查Logo是否已过期（仅对临时Logo）
	if logoModel.IsTemporary() && logoModel.IsExpired() {
		return nil, errno.ErrLogoExpired.WithMessage("Logo已过期")
	}

	// 5. 转换为Thrift响应
	response := l.converter.Logo().ModelToThrift(logoModel)

	// 6. 生成下载URL
	expireSeconds := 3600 // 默认1小时
	if logoModel.IsBound() {
		expireSeconds = 7 * 24 * 3600 // 已绑定的Logo使用7天过期
	}

	downloadURL, err := l.generateDownloadURL(ctx, logoModel.FileID, expireSeconds)
	if err == nil {
		response.DownloadUrl = &downloadURL
	}

	return response, nil
}

// DeleteOrganizationLogo 删除Logo（软删除+S3文件删除）
func (l *LogicImpl) DeleteOrganizationLogo(
	ctx context.Context,
	req *identity_srv.DeleteOrganizationLogoRequest,
) error {
	// 1. 参数验证
	if req.LogoID == nil || *req.LogoID == "" {
		return errno.ErrInvalidParams.WithMessage("LogoID不能为空")
	}

	// 2. 解析LogoID
	logoID, err := uuid.Parse(*req.LogoID)
	if err != nil {
		return errno.ErrInvalidParams.WithMessage("无效的LogoID格式")
	}

	// 3. 查询Logo元数据
	logoModel, err := l.repo.GetByID(ctx, logoID.String())
	if err != nil {
		return errno.ErrLogoNotFound.WithMessage("Logo不存在")
	}

	// 4. 删除S3文件
	if err := l.logoStorageClient.DeleteLogo(ctx, logoModel.FileID); err != nil {
		return errno.ErrFileDeleteFailed.WithMessage(fmt.Sprintf("删除文件失败: %v", err))
	}

	// 5. 软删除数据库记录
	if err := l.repo.Delete(ctx, logoID.String()); err != nil {
		return errno.WrapDatabaseError(err, "删除Logo元数据失败")
	}

	return nil
}

// BindLogoToOrganization 绑定Logo到组织（临时→永久）
func (l *LogicImpl) BindLogoToOrganization(
	ctx context.Context,
	req *identity_srv.BindLogoToOrganizationRequest,
) (*identity_srv.OrganizationLogo, error) {
	// 1. 参数验证
	if req.LogoID == nil || *req.LogoID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("LogoID不能为空")
	}

	if req.OrganizationID == nil || *req.OrganizationID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("OrganizationID不能为空")
	}

	// 2. 解析ID
	logoID, err := uuid.Parse(*req.LogoID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的LogoID格式")
	}

	orgID, err := uuid.Parse(*req.OrganizationID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的OrganizationID格式")
	}

	// 3. 查询Logo元数据
	logoModel, err := l.repo.GetByID(ctx, logoID.String())
	if err != nil {
		return nil, errno.ErrLogoNotFound.WithMessage("Logo不存在")
	}

	// 4. 验证Logo状态
	if !logoModel.IsTemporary() {
		return nil, errno.ErrLogoAlreadyBound.WithMessage("Logo已被绑定")
	}

	// 5. 验证Logo是否过期
	if logoModel.IsExpired() {
		return nil, errno.ErrLogoExpired.WithMessage("Logo已过期，无法绑定")
	}

	// 6. 检查并清理组织的旧Logo（如果存在）
	oldLogo, err := l.repo.GetByOrganizationID(ctx, orgID)
	if err == nil && oldLogo != nil {
		// 删除旧Logo的S3文件
		if deleteErr := l.logoStorageClient.DeleteLogo(ctx, oldLogo.FileID); deleteErr != nil {
			// S3删除失败记录日志，但不阻止绑定流程
			_ = deleteErr
		}

		// 软删除旧Logo的数据库记录
		if deleteErr := l.repo.Delete(ctx, oldLogo.ID.String()); deleteErr != nil {
			// 数据库删除失败记录日志，但不阻止绑定流程
			_ = deleteErr
		}
	}

	// 7. 绑定Logo到组织
	if err := l.repo.BindToOrganization(ctx, logoID, orgID); err != nil {
		return nil, errno.ErrLogoBindingFailed.WithMessage(fmt.Sprintf("绑定Logo失败: %v", err))
	}

	// 8. 更新S3对象标签（从temporary改为permanent）
	if err := l.logoStorageClient.UpdateLogoTagToPermanent(ctx, logoModel.FileID); err != nil {
		// S3标签更新失败不影响业务流程，仅记录错误
		// 实际生产环境应该使用 logger 记录
		_ = err
	}

	// 9. 查询更新后的Logo
	updatedLogo, err := l.repo.GetByID(ctx, logoID.String())
	if err != nil {
		return nil, errno.WrapDatabaseError(err, "查询更新后的Logo失败")
	}

	// 10. 转换为Thrift响应
	response := l.converter.Logo().ModelToThrift(updatedLogo)

	// 11. 生成下载URL（绑定后的Logo使用长期URL，7天过期）
	downloadURL, err := l.generateDownloadURL(ctx, updatedLogo.FileID, 7*24*3600)
	if err == nil {
		response.DownloadUrl = &downloadURL
	}

	return response, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// generateDownloadURL 生成下载URL
func (l *LogicImpl) generateDownloadURL(
	ctx context.Context,
	fileID string,
	expireSeconds int,
) (string, error) {
	// 使用Logo存储客户端生成预签名URL
	downloadURL, err := l.logoStorageClient.GetLogoURL(ctx, fileID, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("生成下载URL失败: %v", err)
	}

	return downloadURL, nil
}

// CleanupExpiredLogos 定期清理过期的临时Logo（供定时任务调用）
func (l *LogicImpl) CleanupExpiredLogos(ctx context.Context) (int64, error) {
	// 1. 获取所有过期的临时Logo
	expiredLogos, err := l.repo.ListExpiredTemporaryLogos(ctx)
	if err != nil {
		return 0, errno.WrapDatabaseError(err, "查询过期Logo失败")
	}

	if len(expiredLogos) == 0 {
		return 0, nil
	}

	// 2. 删除S3文件并软删除数据库记录
	var deletedCount int64

	for _, logo := range expiredLogos {
		// 删除S3文件
		if err := l.logoStorageClient.DeleteLogo(ctx, logo.FileID); err != nil {
			// 记录错误但继续处理下一个
			continue
		}

		// 软删除数据库记录
		if err := l.repo.Delete(ctx, logo.ID.String()); err != nil {
			continue
		}

		deletedCount++
	}

	return deletedCount, nil
}
