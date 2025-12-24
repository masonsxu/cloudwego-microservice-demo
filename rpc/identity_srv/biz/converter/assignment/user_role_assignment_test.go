package assignment

import (
	"testing"

	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	if converter == nil {
		t.Fatal("NewConverter() should not return nil")
	}
}

func TestModelToThrift_NilInput(t *testing.T) {
	converter := NewConverter()

	result := converter.ModelToThrift(nil)

	if result != nil {
		t.Errorf("ModelToThrift(nil) = %v, want nil", result)
	}
}

func TestModelToThrift_FullModel(t *testing.T) {
	converter := NewConverter()

	// 准备测试数据
	id := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()
	createdBy := uuid.New()
	updatedBy := uuid.New()
	var createdAt int64 = 1703472000000 // 2023-12-25 00:00:00 UTC
	var updatedAt int64 = 1703558400000 // 2023-12-26 00:00:00 UTC

	model := &models.UserRoleAssignment{
		BaseModel: models.BaseModel{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		UserID:    userID,
		RoleID:    roleID,
		CreatedBy: &createdBy,
		UpdatedBy: &updatedBy,
	}

	result := converter.ModelToThrift(model)

	// 验证结果不为 nil
	if result == nil {
		t.Fatal("ModelToThrift() should not return nil for valid input")
	}

	// 验证 ID 转换
	if result.Id == nil || *result.Id != id.String() {
		t.Errorf("Id = %v, want %v", safeDeref(result.Id), id.String())
	}

	// 验证 UserID 转换
	if result.UserID == nil || *result.UserID != userID.String() {
		t.Errorf("UserID = %v, want %v", safeDeref(result.UserID), userID.String())
	}

	// 验证 RoleID 转换
	if result.RoleID == nil || *result.RoleID != roleID.String() {
		t.Errorf("RoleID = %v, want %v", safeDeref(result.RoleID), roleID.String())
	}

	// 验证 CreatedBy 转换
	if result.CreatedBy == nil || *result.CreatedBy != createdBy.String() {
		t.Errorf("CreatedBy = %v, want %v", safeDeref(result.CreatedBy), createdBy.String())
	}

	// 验证 UpdatedBy 转换
	if result.UpdatedBy == nil || *result.UpdatedBy != updatedBy.String() {
		t.Errorf("UpdatedBy = %v, want %v", safeDeref(result.UpdatedBy), updatedBy.String())
	}

	// 验证 CreatedAt 转换
	if result.CreatedAt == nil || *result.CreatedAt != createdAt {
		t.Errorf("CreatedAt = %v, want %v", safeDerefInt64(result.CreatedAt), createdAt)
	}

	// 验证 UpdatedAt 转换
	if result.UpdatedAt == nil || *result.UpdatedAt != updatedAt {
		t.Errorf("UpdatedAt = %v, want %v", safeDerefInt64(result.UpdatedAt), updatedAt)
	}
}

func TestModelToThrift_OptionalFieldsNil(t *testing.T) {
	converter := NewConverter()

	// 准备测试数据 - 可选字段为 nil
	id := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()
	var createdAt int64 = 1703472000000
	var updatedAt int64 = 1703558400000

	model := &models.UserRoleAssignment{
		BaseModel: models.BaseModel{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		UserID:    userID,
		RoleID:    roleID,
		CreatedBy: nil, // 可选字段为 nil
		UpdatedBy: nil, // 可选字段为 nil
	}

	result := converter.ModelToThrift(model)

	// 验证结果不为 nil
	if result == nil {
		t.Fatal("ModelToThrift() should not return nil for valid input")
	}

	// 验证必填字段仍然正确转换
	if result.Id == nil || *result.Id != id.String() {
		t.Errorf("Id = %v, want %v", safeDeref(result.Id), id.String())
	}

	if result.UserID == nil || *result.UserID != userID.String() {
		t.Errorf("UserID = %v, want %v", safeDeref(result.UserID), userID.String())
	}

	if result.RoleID == nil || *result.RoleID != roleID.String() {
		t.Errorf("RoleID = %v, want %v", safeDeref(result.RoleID), roleID.String())
	}

	// 验证可选字段为 nil
	if result.CreatedBy != nil {
		t.Errorf("CreatedBy = %v, want nil", *result.CreatedBy)
	}

	if result.UpdatedBy != nil {
		t.Errorf("UpdatedBy = %v, want nil", *result.UpdatedBy)
	}
}

func TestModelToThrift_ZeroTimestamps(t *testing.T) {
	converter := NewConverter()

	// 准备测试数据 - 时间戳为 0
	model := &models.UserRoleAssignment{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: 0,
			UpdatedAt: 0,
		},
		UserID: uuid.New(),
		RoleID: uuid.New(),
	}

	result := converter.ModelToThrift(model)

	if result == nil {
		t.Fatal("ModelToThrift() should not return nil for valid input")
	}

	// 验证零值时间戳正确转换
	if result.CreatedAt == nil || *result.CreatedAt != 0 {
		t.Errorf("CreatedAt = %v, want 0", safeDerefInt64(result.CreatedAt))
	}

	if result.UpdatedAt == nil || *result.UpdatedAt != 0 {
		t.Errorf("UpdatedAt = %v, want 0", safeDerefInt64(result.UpdatedAt))
	}
}

func TestModelToThrift_UUIDFormat(t *testing.T) {
	converter := NewConverter()

	// 使用固定 UUID 验证格式
	id, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	userID, _ := uuid.Parse("660e8400-e29b-41d4-a716-446655440001")
	roleID, _ := uuid.Parse("770e8400-e29b-41d4-a716-446655440002")

	model := &models.UserRoleAssignment{
		BaseModel: models.BaseModel{
			ID: id,
		},
		UserID: userID,
		RoleID: roleID,
	}

	result := converter.ModelToThrift(model)

	if result == nil {
		t.Fatal("ModelToThrift() should not return nil")
	}

	// 验证 UUID 字符串格式正确
	expectedID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUserID := "660e8400-e29b-41d4-a716-446655440001"
	expectedRoleID := "770e8400-e29b-41d4-a716-446655440002"

	if result.Id == nil || *result.Id != expectedID {
		t.Errorf("Id = %v, want %v", safeDeref(result.Id), expectedID)
	}

	if result.UserID == nil || *result.UserID != expectedUserID {
		t.Errorf("UserID = %v, want %v", safeDeref(result.UserID), expectedUserID)
	}

	if result.RoleID == nil || *result.RoleID != expectedRoleID {
		t.Errorf("RoleID = %v, want %v", safeDeref(result.RoleID), expectedRoleID)
	}
}

// 辅助函数：安全解引用字符串指针
func safeDeref(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

// 辅助函数：安全解引用 int64 指针
func safeDerefInt64(i *int64) int64 {
	if i == nil {
		return -1
	}
	return *i
}
