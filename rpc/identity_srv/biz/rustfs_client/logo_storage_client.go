package rustfsclient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// LogoStorageClient 组织Logo文件存储客户端接口
// 专门用于组织Logo的上传、下载、删除和生命周期管理
type LogoStorageClient interface {
	// UploadTemporaryLogo 上传临时Logo文件
	// 上传的Logo会被标记为 Status=temporary，7天后自动删除（除非转为永久）
	//
	// 参数:
	//   - ctx: 上下文
	//   - uploaderID: 上传者ID（用于文件路径构建）
	//   - fileName: 原始文件名
	//   - content: 文件内容
	//   - mimeType: MIME类型（用于验证）
	//
	// 返回:
	//   - fileID: 文件标识（格式: bucket/objectKey）
	//   - error: 错误信息
	UploadTemporaryLogo(
		ctx context.Context,
		uploaderID uuid.UUID,
		fileName string,
		content []byte,
		mimeType string,
	) (string, error)

	// DeleteLogo 删除Logo文件
	//
	// 参数:
	//   - ctx: 上下文
	//   - fileID: 文件标识（格式: bucket/objectKey）
	//
	// 返回:
	//   - error: 错误信息
	DeleteLogo(ctx context.Context, fileID string) error

	// GetLogoURL 生成Logo的预签名下载URL
	//
	// 参数:
	//   - ctx: 上下文
	//   - fileID: 文件标识（格式: bucket/objectKey）
	//   - expireSeconds: URL过期时间（秒），0表示使用默认值（7天）
	//
	// 返回:
	//   - url: 预签名URL
	//   - error: 错误信息
	GetLogoURL(ctx context.Context, fileID string, expireSeconds int) (string, error)

	// UpdateLogoTagToPermanent 将Logo标签从temporary更新为permanent
	// 当Logo被正式绑定到组织后调用，防止被自动删除
	//
	// 参数:
	//   - ctx: 上下文
	//   - fileID: 文件标识（格式: bucket/objectKey）
	//
	// 返回:
	//   - error: 错误信息
	UpdateLogoTagToPermanent(ctx context.Context, fileID string) error

	// ConfigureS3LifecyclePolicy 配置S3生命周期策略
	// 设置自动删除标记为temporary的文件（7天后）
	//
	// 参数:
	//   - ctx: 上下文
	//
	// 返回:
	//   - error: 错误信息
	ConfigureS3LifecyclePolicy(ctx context.Context) error

	// ValidateFileSize 验证文件大小（最大10MB）
	ValidateFileSize(size int64) error

	// ValidateFileType 验证文件MIME类型（仅支持图片）
	ValidateFileType(mimeType string) error
}

// logoStorageClientImpl Logo存储客户端实现（基于 AWS S3 SDK v2）
type logoStorageClientImpl struct {
	s3Client       *s3.Client // 内部端点客户端（用于上传、删除等操作）
	s3PublicClient *s3.Client // 公共端点客户端（用于生成预签名 URL）
	cfg            *config.LogoStorageConfig
}

// NewLogoStorageClient 创建 Logo 存储客户端
func NewLogoStorageClient(cfg *config.LogoStorageConfig) (LogoStorageClient, error) {
	if cfg == nil {
		return nil, errors.New("logo storage config is nil")
	}

	// 验证必要的配置参数
	if err := validateLogoConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 构建 AWS SDK 配置
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		// 静态凭证
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 创建内部端点 S3 客户端（用于上传、删除等操作）
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// 使用 BaseEndpoint 设置内部端点
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		// 根据配置选择 Path Style 或 Virtual Host Style
		o.UsePathStyle = cfg.UsePathStyle
	})

	// 确定公共端点（如果未配置则使用内部端点）
	publicEndpoint := cfg.S3PublicEndpoint
	if publicEndpoint == "" {
		publicEndpoint = cfg.S3Endpoint
	}

	// 创建公共端点 S3 客户端（用于生成预签名 URL）
	s3PublicClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// 使用 BaseEndpoint 设置公共端点
		o.BaseEndpoint = aws.String(publicEndpoint)
		// 根据配置选择 Path Style 或 Virtual Host Style
		o.UsePathStyle = cfg.UsePathStyle
	})

	// 创建客户端实例
	client := &logoStorageClientImpl{
		s3Client:       s3Client,
		s3PublicClient: s3PublicClient,
		cfg:            cfg,
	}

	// 确保存储桶存在
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketName := getLogoBucketName()
	if err := client.ensureBucket(ctx, bucketName); err != nil {
		return nil, fmt.Errorf("failed to ensure logo bucket exists: %w", err)
	}

	// 配置生命周期策略
	if err := client.ConfigureS3LifecyclePolicy(ctx); err != nil {
		// 生命周期策略配置失败不阻止客户端初始化，仅记录警告
		// 实际生产环境应该使用 logger 记录
		fmt.Printf("Warning: failed to configure S3 lifecycle policy: %v\n", err)
	}

	return client, nil
}

// UploadTemporaryLogo 上传临时Logo文件
func (c *logoStorageClientImpl) UploadTemporaryLogo(
	ctx context.Context,
	uploaderID uuid.UUID,
	fileName string,
	content []byte,
	mimeType string,
) (string, error) {
	// 验证文件大小
	if err := c.ValidateFileSize(int64(len(content))); err != nil {
		return "", err
	}

	// 验证文件类型
	if err := c.ValidateFileType(mimeType); err != nil {
		return "", err
	}

	// 构建对象路径: {uploaderID}/{timestamp}_{fileName}
	timestamp := time.Now().UnixMilli()
	objectKey := fmt.Sprintf("%s/%d_%s",
		uploaderID.String(),
		timestamp,
		sanitizeFileName(fileName),
	)

	bucket := getLogoBucketName()

	// 上传文件，并添加 temporary 标签
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		Body:        strings.NewReader(string(content)),
		ContentType: aws.String(mimeType),
		// 添加对象标签：Status=temporary（临时Logo，7天后自动删除）
		Tagging: aws.String("Status=temporary"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload logo: %w", err)
	}

	// 构建 fileID（格式: bucket/objectKey）
	fileID := fmt.Sprintf("%s/%s", bucket, objectKey)

	return fileID, nil
}

// DeleteLogo 删除Logo文件
func (c *logoStorageClientImpl) DeleteLogo(ctx context.Context, fileID string) error {
	if fileID == "" {
		return errors.New("file ID is empty")
	}

	// 解析 fileID（格式: bucket/objectKey）
	bucket, objectKey := parseFileID(fileID)
	if bucket == "" || objectKey == "" {
		return errors.New("invalid file ID format")
	}

	// 删除对象
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete logo: %w", err)
	}

	return nil
}

// GetLogoURL 生成Logo的预签名下载URL
// 使用公共端点客户端生成 URL，确保外部客户端（如浏览器）可以访问
func (c *logoStorageClientImpl) GetLogoURL(
	ctx context.Context,
	fileID string,
	expireSeconds int,
) (string, error) {
	if fileID == "" {
		return "", errors.New("file ID is empty")
	}

	// 解析 fileID（格式: bucket/objectKey）
	bucket, objectKey := parseFileID(fileID)
	if bucket == "" || objectKey == "" {
		return "", errors.New("invalid file ID format")
	}

	// 如果未指定过期时间，使用默认值（7天）
	if expireSeconds <= 0 {
		expireSeconds = 7 * 24 * 3600 // 7天
	}

	// 使用公共端点客户端创建 Presign 客户端
	// 这样生成的 URL 包含公共端点地址，外部客户端可以访问
	presignClient := s3.NewPresignClient(c.s3PublicClient)

	// 生成预签名 URL
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(time.Duration(expireSeconds)*time.Second))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedReq.URL, nil
}

// UpdateLogoTagToPermanent 将Logo标签从temporary更新为permanent
func (c *logoStorageClientImpl) UpdateLogoTagToPermanent(ctx context.Context, fileID string) error {
	if fileID == "" {
		return errors.New("file ID is empty")
	}

	// 解析 fileID（格式: bucket/objectKey）
	bucket, objectKey := parseFileID(fileID)
	if bucket == "" || objectKey == "" {
		return errors.New("invalid file ID format")
	}

	// 更新对象标签为 permanent（永久Logo，不会被自动删除）
	_, err := c.s3Client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Tagging: &types.Tagging{
			TagSet: []types.Tag{
				{
					Key:   aws.String("Status"),
					Value: aws.String("permanent"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update logo tag to permanent: %w", err)
	}

	return nil
}

// ConfigureS3LifecyclePolicy 配置S3生命周期策略
func (c *logoStorageClientImpl) ConfigureS3LifecyclePolicy(ctx context.Context) error {
	bucket := getLogoBucketName()

	// 定义生命周期策略：7天后删除标记为 temporary 的对象
	lifecycleRules := []types.LifecycleRule{
		{
			// 规则ID
			ID: aws.String("delete-temporary-logos-after-7-days"),
			// 规则状态：启用
			Status: types.ExpirationStatusEnabled,
			// 过滤条件：仅应用于标记为 Status=temporary 的对象
			Filter: &types.LifecycleRuleFilter{
				Tag: &types.Tag{
					Key:   aws.String("Status"),
					Value: aws.String("temporary"),
				},
			},
			// 过期策略：7天后删除
			Expiration: &types.LifecycleExpiration{
				Days: aws.Int32(7),
			},
		},
	}

	// 应用生命周期策略到存储桶
	_, err := c.s3Client.PutBucketLifecycleConfiguration(
		ctx,
		&s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: lifecycleRules,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to configure S3 lifecycle policy: %w", err)
	}

	return nil
}

// ValidateFileSize 验证文件大小（最大10MB）
func (c *logoStorageClientImpl) ValidateFileSize(size int64) error {
	if size <= 0 {
		return errors.New("file size must be greater than 0")
	}

	// Logo最大10MB
	maxSize := int64(10 * 1024 * 1024)
	if c.cfg.MaxFileSize > 0 {
		maxSize = c.cfg.MaxFileSize
	}

	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", size, maxSize)
	}

	return nil
}

// ValidateFileType 验证文件MIME类型（仅支持常见图片格式）
func (c *logoStorageClientImpl) ValidateFileType(mimeType string) error {
	// 标准化 MIME 类型：去除空格并转为小写
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	if mimeType == "" {
		return errors.New("MIME类型不能为空")
	}

	// Logo允许的图片类型
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}

	// MIME 类型别名映射（处理不同客户端的兼容性）
	mimeTypeAliases := map[string][]string{
		"image/jpeg": {"image/jpg"},
		"image/jpg":  {"image/jpeg"},
	}

	// 检查是否在允许列表中（包括别名）
	for _, allowedType := range allowedTypes {
		// 精确匹配
		if allowedType == mimeType {
			return nil
		}

		// 检查别名匹配
		if aliases, exists := mimeTypeAliases[allowedType]; exists {
			for _, alias := range aliases {
				if alias == mimeType {
					return nil
				}
			}
		}
	}

	return fmt.Errorf(
		"不支持的文件类型 '%s'。仅支持图片格式: %s",
		mimeType,
		strings.Join(allowedTypes, ", "),
	)
}

// ============================================================================
// 辅助函数
// ============================================================================

// getLogoBucketName 获取Logo存储桶名称（固定）
func getLogoBucketName() string {
	return "organization-logos"
}

// ensureBucket 确保 bucket 存在，不存在则创建
func (c *logoStorageClientImpl) ensureBucket(ctx context.Context, bucket string) error {
	// 检查 bucket 是否存在
	_, err := c.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		// Bucket 已存在
		return nil
	}

	// Bucket 不存在，尝试创建
	_, err = c.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		// 检查是否是因为 bucket 已存在导致的错误（并发创建）
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "BucketAlreadyOwnedByYou" ||
				apiErr.ErrorCode() == "BucketAlreadyExists" {
				// Bucket 已存在，忽略错误
				return nil
			}
		}

		return fmt.Errorf("create bucket failed: %w", err)
	}

	return nil
}

// validateLogoConfig 验证Logo存储配置参数
func validateLogoConfig(cfg *config.LogoStorageConfig) error {
	if cfg.S3Endpoint == "" {
		return errors.New("S3 endpoint is required")
	}

	if cfg.S3Region == "" {
		return errors.New("S3 region is required")
	}

	if cfg.AccessKey == "" {
		return errors.New("access key is required")
	}

	if cfg.SecretKey == "" {
		return errors.New("secret key is required")
	}

	// 验证端点格式
	if !strings.HasPrefix(cfg.S3Endpoint, "http://") &&
		!strings.HasPrefix(cfg.S3Endpoint, "https://") {
		return errors.New("S3 endpoint must start with http:// or https://")
	}

	return nil
}

// parseFileID 解析fileID为bucket和objectKey
// fileID格式: "bucket/objectKey"
func parseFileID(fileID string) (bucket, objectKey string) {
	parts := strings.SplitN(fileID, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

// sanitizeFileName 清理文件名，确保安全
func sanitizeFileName(fileName string) string {
	// 只保留文件名部分，移除路径
	fileName = strings.TrimSpace(fileName)
	// 移除路径分隔符
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	fileName = strings.ReplaceAll(fileName, "..", "_")

	if fileName == "" || fileName == "." {
		return "unnamed_file"
	}

	return fileName
}
