package department

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_ModelDepartmentToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelDepartmentToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整部门转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		model := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Engineering",
			OrganizationID:     orgID,
			DepartmentType:     "Technical",
			AvailableEquipment: `["equip1", "equip2", "equip3"]`,
		}

		result := converter.ModelDepartmentToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, deptID.String(), *result.ID)
		assert.Equal(t, "Engineering", *result.Name)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Equal(t, "Technical", *result.DepartmentType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		require.Len(t, result.AvailableEquipment, 3)
		assert.Equal(t, "equip1", result.AvailableEquipment[0])
		assert.Equal(t, "equip2", result.AvailableEquipment[1])
		assert.Equal(t, "equip3", result.AvailableEquipment[2])
	})

	t.Run("最小部门转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		model := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "HR",
			OrganizationID:     orgID,
			DepartmentType:     "",
			AvailableEquipment: "",
		}

		result := converter.ModelDepartmentToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, deptID.String(), *result.ID)
		assert.Equal(t, "HR", *result.Name)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Nil(t, result.DepartmentType) // 空字符串应该转换为nil
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Empty(t, result.AvailableEquipment) // 空设备列表
	})

	t.Run("无效JSON设备列表处理", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		model := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Finance",
			OrganizationID:     orgID,
			DepartmentType:     "Administrative",
			AvailableEquipment: "invalid json",
		}

		result := converter.ModelDepartmentToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, deptID.String(), *result.ID)
		assert.Equal(t, "Finance", *result.Name)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Equal(t, "Administrative", *result.DepartmentType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Empty(t, result.AvailableEquipment) // 无效JSON应该返回空列表
	})

	t.Run("空JSON数组设备列表处理", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		model := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Marketing",
			OrganizationID:     orgID,
			DepartmentType:     "Business",
			AvailableEquipment: "[]",
		}

		result := converter.ModelDepartmentToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, deptID.String(), *result.ID)
		assert.Equal(t, "Marketing", *result.Name)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Equal(t, "Business", *result.DepartmentType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Empty(t, result.AvailableEquipment) // 空数组应该返回空列表
	})

	t.Run("单个设备列表处理", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		model := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Lab",
			OrganizationID:     orgID,
			DepartmentType:     "Research",
			AvailableEquipment: `["single-equip"]`,
		}

		result := converter.ModelDepartmentToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, deptID.String(), *result.ID)
		assert.Equal(t, "Lab", *result.Name)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Equal(t, "Research", *result.DepartmentType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		require.Len(t, result.AvailableEquipment, 1)
		assert.Equal(t, "single-equip", result.AvailableEquipment[0])
	})
}

func TestConverterImpl_ModelDepartmentsToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ModelDepartmentsToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("空切片转换", func(t *testing.T) {
		result := converter.ModelDepartmentsToThrift([]*models.Department{})
		assert.Nil(t, result)
	})

	t.Run("多个部门转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID1 := uuid.New()
		deptID2 := uuid.New()
		orgID1 := uuid.New()
		orgID2 := uuid.New()

		models := []*models.Department{
			{
				BaseModel: models.BaseModel{
					ID:        deptID1,
					CreatedAt: nowMs,
					UpdatedAt: nowMs,
				},
				Name:               "Engineering",
				OrganizationID:     orgID1,
				DepartmentType:     "Technical",
				AvailableEquipment: `["equip1"]`,
			},
			{
				BaseModel: models.BaseModel{
					ID:        deptID2,
					CreatedAt: nowMs,
					UpdatedAt: nowMs,
				},
				Name:               "HR",
				OrganizationID:     orgID2,
				DepartmentType:     "Administrative",
				AvailableEquipment: "",
			},
		}

		result := converter.ModelDepartmentsToThrift(models)

		require.NotNil(t, result)
		require.Len(t, result, 2)

		// 验证第一个部门
		dept1 := result[0]
		assert.Equal(t, deptID1.String(), *dept1.ID)
		assert.Equal(t, "Engineering", *dept1.Name)
		assert.Equal(t, orgID1.String(), *dept1.OrganizationID)
		assert.Equal(t, "Technical", *dept1.DepartmentType)
		require.Len(t, dept1.AvailableEquipment, 1)
		assert.Equal(t, "equip1", dept1.AvailableEquipment[0])

		// 验证第二个部门
		dept2 := result[1]
		assert.Equal(t, deptID2.String(), *dept2.ID)
		assert.Equal(t, "HR", *dept2.Name)
		assert.Equal(t, orgID2.String(), *dept2.OrganizationID)
		assert.Equal(t, "Administrative", *dept2.DepartmentType)
		assert.Empty(t, dept2.AvailableEquipment)
	})

	t.Run("包含nil部门的切片转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		models := []*models.Department{
			{
				BaseModel: models.BaseModel{
					ID:        deptID,
					CreatedAt: nowMs,
					UpdatedAt: nowMs,
				},
				Name:           "Engineering",
				OrganizationID: orgID,
			},
			nil,
			nil,
		}

		result := converter.ModelDepartmentsToThrift(models)

		require.NotNil(t, result)
		require.Len(t, result, 1) // 只有非nil的部门被转换

		dept := result[0]
		assert.Equal(t, deptID.String(), *dept.ID)
		assert.Equal(t, "Engineering", *dept.Name)
		assert.Equal(t, orgID.String(), *dept.OrganizationID)
	})
}

func TestConverterImpl_ModelToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("别名方法测试", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		model := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:           "Test",
			OrganizationID: orgID,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, deptID.String(), *result.ID)
		assert.Equal(t, "Test", *result.Name)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
	})
}

func TestConverterImpl_CreateRequestToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.CreateRequestToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整创建请求转换", func(t *testing.T) {
		orgID := uuid.New()
		orgIDStr := orgID.String()
		name := "Engineering"
		deptType := "Technical"

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgIDStr,
			Name:           &name,
			DepartmentType: &deptType,
		}

		result := converter.CreateRequestToModel(req)

		require.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil, result.ID) // 应该生成新的ID
		assert.Equal(t, name, result.Name)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, deptType, result.DepartmentType)
	})

	t.Run("最小创建请求转换", func(t *testing.T) {
		orgID := uuid.New()
		orgIDStr := orgID.String()
		name := "HR"

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgIDStr,
			Name:           &name,
		}

		result := converter.CreateRequestToModel(req)

		require.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil, result.ID) // 应该生成新的ID
		assert.Equal(t, name, result.Name)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, "", result.DepartmentType) // 默认空值
	})
}

func TestConverterImpl_ApplyUpdateToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入处理", func(t *testing.T) {
		// 两个都是nil
		result := converter.ApplyUpdateToModel(nil, nil)
		assert.Nil(t, result)

		// existing为nil
		result = converter.ApplyUpdateToModel(nil, &identity_srv.UpdateDepartmentRequest{})
		assert.Nil(t, result)

		// req为nil
		result = converter.ApplyUpdateToModel(&models.Department{}, nil)
		assert.NotNil(t, result) // 应该返回原始模型
	})

	t.Run("完整更新应用", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		existing := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Old Name",
			OrganizationID:     orgID,
			DepartmentType:     "Old Type",
			AvailableEquipment: "[]",
		}

		newName := "New Name"
		newType := "New Type"
		req := &identity_srv.UpdateDepartmentRequest{
			Name:           &newName,
			DepartmentType: &newType,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, deptID, result.ID) // ID应该保持不变
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, orgID, result.OrganizationID) // 组织ID应该保持不变
		assert.Equal(t, "New Type", result.DepartmentType)
		assert.Equal(t, "[]", result.AvailableEquipment) // 未更新的字段应该保持不变
	})

	t.Run("部分更新应用", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		existing := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Original Name",
			OrganizationID:     orgID,
			DepartmentType:     "Original Type",
			AvailableEquipment: `["equip1"]`,
		}

		// 只更新名称
		newName := "Updated Name"
		req := &identity_srv.UpdateDepartmentRequest{
			Name: &newName,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, deptID, result.ID)
		assert.Equal(t, "Updated Name", result.Name)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, "Original Type", result.DepartmentType)  // 未更新的字段应该保持不变
		assert.Equal(t, `["equip1"]`, result.AvailableEquipment) // 未更新的字段应该保持不变
	})

	t.Run("空更新请求", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		existing := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Test",
			OrganizationID:     orgID,
			DepartmentType:     "Test Type",
			AvailableEquipment: "",
		}

		req := &identity_srv.UpdateDepartmentRequest{} // 空更新请求

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, deptID, result.ID)
		assert.Equal(t, "Test", result.Name) // 所有字段应该保持不变
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, "Test Type", result.DepartmentType)
		assert.Equal(t, "", result.AvailableEquipment)
	})

	t.Run("更新为空值", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		deptID := uuid.New()
		orgID := uuid.New()

		existing := &models.Department{
			BaseModel: models.BaseModel{
				ID:        deptID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               "Test",
			OrganizationID:     orgID,
			DepartmentType:     "Original Type",
			AvailableEquipment: "",
		}

		// 更新为空字符串
		emptyType := ""
		req := &identity_srv.UpdateDepartmentRequest{
			DepartmentType: &emptyType,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, deptID, result.ID)
		assert.Equal(t, "Test", result.Name)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, "", result.DepartmentType) // 应该被更新为空字符串
		assert.Equal(t, "", result.AvailableEquipment)
	})
}

// 表格驱动测试
func TestConverterImpl_ModelDepartmentToThrift_TableDriven(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name        string
		input       *models.Department
		expectNil   bool
		description string
	}{
		{
			name:        "nil输入",
			input:       nil,
			expectNil:   true,
			description: "nil输入应该返回nil",
		},
		{
			name:        "空部门模型",
			input:       &models.Department{},
			expectNil:   false,
			description: "空部门模型应该返回非nil结果",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ModelDepartmentToThrift(tt.input)

			if tt.expectNil {
				assert.Nil(t, result, tt.description)
			} else {
				assert.NotNil(t, result, tt.description)
			}
		})
	}
}

// 基准测试
func BenchmarkConverterImpl_ModelDepartmentToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	now := time.Now()
	nowMs := now.UnixMilli()
	deptID := uuid.New()
	orgID := uuid.New()

	model := &models.Department{
		BaseModel: models.BaseModel{
			ID:        deptID,
			CreatedAt: nowMs,
			UpdatedAt: nowMs,
		},
		Name:               "Benchmark Department",
		OrganizationID:     orgID,
		DepartmentType:     "Technical",
		AvailableEquipment: `["equip1", "equip2", "equip3"]`,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelDepartmentToThrift(model)
	}
}

func BenchmarkConverterImpl_ModelDepartmentsToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	now := time.Now()
	nowMs := now.UnixMilli()

	departments := make([]*models.Department, 100)
	for i := 0; i < 100; i++ {
		departments[i] = &models.Department{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:               fmt.Sprintf("Department %d", i),
			OrganizationID:     uuid.New(),
			DepartmentType:     "Technical",
			AvailableEquipment: `["equip1"]`,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelDepartmentsToThrift(departments)
	}
}

func BenchmarkConverterImpl_CreateRequestToModel(b *testing.B) {
	converter := &ConverterImpl{}

	orgID := uuid.New()
	orgIDStr := orgID.String()
	name := "Benchmark Department"
	deptType := "Technical"

	req := &identity_srv.CreateDepartmentRequest{
		OrganizationID: &orgIDStr,
		Name:           &name,
		DepartmentType: &deptType,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.CreateRequestToModel(req)
	}
}

func BenchmarkConverterImpl_ApplyUpdateToModel(b *testing.B) {
	converter := &ConverterImpl{}

	now := time.Now()
	nowMs := now.UnixMilli()
	deptID := uuid.New()
	orgID := uuid.New()

	existing := &models.Department{
		BaseModel: models.BaseModel{
			ID:        deptID,
			CreatedAt: nowMs,
			UpdatedAt: nowMs,
		},
		Name:               "Original Name",
		OrganizationID:     orgID,
		DepartmentType:     "Original Type",
		AvailableEquipment: "",
	}

	newName := "Updated Name"
	req := &identity_srv.UpdateDepartmentRequest{
		Name: &newName,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ApplyUpdateToModel(existing, req)
	}
}
