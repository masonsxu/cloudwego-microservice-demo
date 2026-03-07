package auditlog

import (
	"math"

	"github.com/google/uuid"

	auditlogDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/auditlog"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter 审计日志转换器接口
type Converter interface {
	// ModelToThrift 将审计日志 Model 转换为 Thrift DTO
	ModelToThrift(log *models.AuditLog) *identity_srv.AuditLog

	// ModelsToThrift 批量将审计日志 Model 转换为 Thrift DTO
	ModelsToThrift(logs []*models.AuditLog) []*identity_srv.AuditLog

	// CreateRequestToModel 将创建请求转换为 Model
	CreateRequestToModel(req *identity_srv.CreateAuditLogRequest) *models.AuditLog

	// StatsToThrift 将 DAL 统计结果转换为 Thrift DTO
	StatsToThrift(stats *auditlogDAL.AuditLogStats) *identity_srv.AuditLogStats
}

// ConverterImpl 审计日志转换器实现
type ConverterImpl struct{}

// NewConverter 创建审计日志转换器
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ModelToThrift 将审计日志 Model 转换为 Thrift DTO
func (c *ConverterImpl) ModelToThrift(log *models.AuditLog) *identity_srv.AuditLog {
	if log == nil {
		return nil
	}

	id := log.ID.String()
	action := identity_srv.AuditAction(log.Action)

	dto := &identity_srv.AuditLog{
		Id:          &id,
		RequestID:   &log.RequestID,
		TraceID:     &log.TraceID,
		Username:    &log.Username,
		Action:      &action,
		Resource:    &log.Resource,
		ResourceID:  &log.ResourceID,
		StatusCode:  &log.StatusCode,
		Success:     &log.Success,
		ClientIP:    &log.ClientIP,
		UserAgent:   &log.UserAgent,
		RequestBody: &log.RequestBody,
		DurationMs:  &log.DurationMs,
		CreatedAt:   &log.CreatedAt,
	}

	if log.UserID != nil {
		userID := log.UserID.String()
		dto.UserID = &userID
	}

	if log.OrganizationID != nil {
		orgID := log.OrganizationID.String()
		dto.OrganizationID = &orgID
	}

	return dto
}

// ModelsToThrift 批量将审计日志 Model 转换为 Thrift DTO
func (c *ConverterImpl) ModelsToThrift(logs []*models.AuditLog) []*identity_srv.AuditLog {
	if len(logs) == 0 {
		return nil
	}

	dtos := make([]*identity_srv.AuditLog, 0, len(logs))

	for _, log := range logs {
		if dto := c.ModelToThrift(log); dto != nil {
			dtos = append(dtos, dto)
		}
	}

	return dtos
}

// CreateRequestToModel 将创建请求转换为 Model
func (c *ConverterImpl) CreateRequestToModel(req *identity_srv.CreateAuditLogRequest) *models.AuditLog {
	if req == nil {
		return nil
	}

	log := &models.AuditLog{
		ID: uuid.New(),
	}

	if req.RequestID != nil {
		log.RequestID = *req.RequestID
	}

	if req.TraceID != nil {
		log.TraceID = *req.TraceID
	}

	if req.UserID != nil && *req.UserID != "" {
		userID, err := uuid.Parse(*req.UserID)
		if err == nil {
			log.UserID = &userID
		}
	}

	if req.Username != nil {
		log.Username = *req.Username
	}

	if req.OrganizationID != nil && *req.OrganizationID != "" {
		orgID, err := uuid.Parse(*req.OrganizationID)
		if err == nil {
			log.OrganizationID = &orgID
		}
	}

	if req.Action != nil {
		log.Action = models.AuditAction(*req.Action)
	}

	if req.Resource != nil {
		log.Resource = *req.Resource
	}

	if req.ResourceID != nil {
		log.ResourceID = *req.ResourceID
	}

	if req.StatusCode != nil {
		log.StatusCode = *req.StatusCode
	}

	if req.Success != nil {
		log.Success = *req.Success
	}

	if req.ClientIP != nil {
		log.ClientIP = *req.ClientIP
	}

	if req.UserAgent != nil {
		log.UserAgent = *req.UserAgent
	}

	if req.RequestBody != nil {
		log.RequestBody = *req.RequestBody
	}

	if req.DurationMs != nil {
		log.DurationMs = *req.DurationMs
	}

	return log
}

// StatsToThrift 将 DAL 统计结果转换为 Thrift DTO
func (c *ConverterImpl) StatsToThrift(stats *auditlogDAL.AuditLogStats) *identity_srv.AuditLogStats {
	if stats == nil {
		return nil
	}

	avgDurationMs := int32(math.Round(stats.AvgDurationMs))

	return &identity_srv.AuditLogStats{
		TotalCount:    &stats.TotalCount,
		SuccessCount:  &stats.SuccessCount,
		AvgDurationMs: &avgDurationMs,
	}
}
