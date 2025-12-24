package authentication

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_ModelUserMembershipsToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("空切片转换", func(t *testing.T) {
		result := converter.ModelUserMembershipsToThrift([]*models.UserMembership{})
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ModelUserMembershipsToThrift(nil)
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("单个成员关系转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		deptID := uuid.New()
		membershipID := uuid.New()

		membership := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        membershipID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   deptID,
			IsPrimary:      true,
		}

		result := converter.ModelUserMembershipsToThrift([]*models.UserMembership{membership})

		require.Len(t, result, 1)
		
		thriftMembership := result[0]
		assert.Equal(t, membershipID.String(), *thriftMembership.ID)
		assert.Equal(t, userID.String(), *thriftMembership.UserID)
		assert.Equal(t, orgID.String(), *thriftMembership.OrganizationID)
		assert.Equal(t, deptID.String(), *thriftMembership.DepartmentID)
		assert.Equal(t, true, *thriftMembership.IsPrimary)
		assert.Equal(t, nowMs, *thriftMembership.CreatedAt)
		assert.Equal(t, nowMs, *thriftMembership.UpdatedAt)
	})

	t.Run("多个成员关系转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		userID1 := uuid.New()
		userID2 := uuid.New()
		orgID1 := uuid.New()
		orgID2 := uuid.New()
		deptID1 := uuid.New()
		membershipID1 := uuid.New()
		membershipID2 := uuid.New()

		memberships := []*models.UserMembership{
			{
				BaseModel: models.BaseModel{
					ID:        membershipID1,
					CreatedAt: nowMs,
					UpdatedAt: nowMs,
				},
				UserID:         userID1,
				OrganizationID: orgID1,
				DepartmentID:   deptID1,
				IsPrimary:      true,
			},
			{
				BaseModel: models.BaseModel{
					ID:        membershipID2,
					CreatedAt: nowMs,
					UpdatedAt: nowMs,
				},
				UserID:         userID2,
				OrganizationID: orgID2,
				DepartmentID:   uuid.Nil, // 无部门ID
				IsPrimary:      false,
			},
		}

		result := converter.ModelUserMembershipsToThrift(memberships)

		require.Len(t, result, 2)

		// 验证第一个成员关系
		thriftMembership1 := result[0]
		assert.Equal(t, membershipID1.String(), *thriftMembership1.ID)
		assert.Equal(t, userID1.String(), *thriftMembership1.UserID)
		assert.Equal(t, orgID1.String(), *thriftMembership1.OrganizationID)
		assert.Equal(t, deptID1.String(), *thriftMembership1.DepartmentID)
		assert.Equal(t, true, *thriftMembership1.IsPrimary)

		// 验证第二个成员关系
		thriftMembership2 := result[1]
		assert.Equal(t, membershipID2.String(), *thriftMembership2.ID)
		assert.Equal(t, userID2.String(), *thriftMembership2.UserID)
		assert.Equal(t, orgID2.String(), *thriftMembership2.OrganizationID)
		assert.Nil(t, thriftMembership2.DepartmentID) // 无部门ID应为nil
		assert.Nil(t, thriftMembership2.IsPrimary)    // IsPrimary为false时不设置字段
	})

	t.Run("包含nil成员的切片", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		membershipID := uuid.New()

		membership := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        membershipID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			UserID:         userID,
			OrganizationID: orgID,
			IsPrimary:      false,
		}

		// 包含nil成员的切片
		memberships := []*models.UserMembership{
			membership,
			nil,
			nil,
		}

		result := converter.ModelUserMembershipsToThrift(memberships)

		require.Len(t, result, 3)
		
		// 第一个元素应该正确转换
		thriftMembership := result[0]
		assert.NotNil(t, thriftMembership)
		assert.Equal(t, membershipID.String(), *thriftMembership.ID)
		assert.Equal(t, userID.String(), *thriftMembership.UserID)
		assert.Equal(t, orgID.String(), *thriftMembership.OrganizationID)
		assert.Nil(t, thriftMembership.IsPrimary) // IsPrimary为false时不设置字段

		// nil元素应该被跳过，结果为nil
		assert.Nil(t, result[1])
		assert.Nil(t, result[2])
	})

	t.Run("所有字段为nil的成员关系", func(t *testing.T) {
		membership := &models.UserMembership{
			UserID:         uuid.Nil,
			OrganizationID: uuid.Nil,
			DepartmentID:   uuid.Nil,
			IsPrimary:      false,
		}

		result := converter.ModelUserMembershipsToThrift([]*models.UserMembership{membership})

		require.Len(t, result, 1)
		
		thriftMembership := result[0]
		assert.Equal(t, uuid.Nil.String(), *thriftMembership.ID)
		assert.Equal(t, uuid.Nil.String(), *thriftMembership.UserID)
		assert.Equal(t, uuid.Nil.String(), *thriftMembership.OrganizationID)
		assert.Nil(t, thriftMembership.DepartmentID) // uuid.Nil应该被跳过
		assert.Nil(t, thriftMembership.IsPrimary)    // IsPrimary为false时不设置字段
	})

	t.Run("零值DepartmentID处理", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		membershipID := uuid.New()

		membership := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        membershipID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   uuid.Nil, // 零值
			IsPrimary:      false,
		}

		result := converter.ModelUserMembershipsToThrift([]*models.UserMembership{membership})

		require.Len(t, result, 1)
		
		thriftMembership := result[0]
		assert.Equal(t, membershipID.String(), *thriftMembership.ID)
		assert.Equal(t, userID.String(), *thriftMembership.UserID)
		assert.Equal(t, orgID.String(), *thriftMembership.OrganizationID)
		assert.Nil(t, thriftMembership.DepartmentID) // 零值DepartmentID应该为nil
		assert.Nil(t, thriftMembership.IsPrimary)    // IsPrimary为false时不设置字段
	})

	t.Run("IsPrimary为false时的处理", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		membershipID := uuid.New()

		membership := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        membershipID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			UserID:         userID,
			OrganizationID: orgID,
			IsPrimary:      false, // 明确设置为false
		}

		result := converter.ModelUserMembershipsToThrift([]*models.UserMembership{membership})

		require.Len(t, result, 1)
		
		thriftMembership := result[0]
		assert.Nil(t, thriftMembership.IsPrimary) // IsPrimary为false时不设置字段
	})
}

// 基准测试
func BenchmarkConverterImpl_ModelUserMembershipsToThrift(b *testing.B) {
	converter := &ConverterImpl{}
	
	now := time.Now()
	nowMs := now.UnixMilli()
	// 准备测试数据
	memberships := make([]*models.UserMembership, 1000)
	for i := 0; i < 1000; i++ {
		memberships[i] = &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			UserID:         uuid.New(),
			OrganizationID: uuid.New(),
			DepartmentID:   uuid.New(),
			IsPrimary:      i%2 == 0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.ModelUserMembershipsToThrift(memberships)
	}
}

// 表格驱动测试
func TestConverterImpl_ModelUserMembershipsToThrift_TableDriven(t *testing.T) {
	converter := &ConverterImpl{}
	
	tests := []struct {
		name        string
		input       []*models.UserMembership
		expectedLen int
		description string
	}{
		{
			name:        "空切片",
			input:       []*models.UserMembership{},
			expectedLen: 0,
			description: "空输入应该返回空切片",
		},
		{
			name:        "nil输入",
			input:       nil,
			expectedLen: 0,
			description: "nil输入应该返回空切片",
		},
		{
			name: "单个有效成员",
			input: []*models.UserMembership{
				{
					BaseModel: models.BaseModel{
						ID:        uuid.New(),
						CreatedAt: time.Now().UnixMilli(),
						UpdatedAt: time.Now().UnixMilli(),
					},
					UserID:         uuid.New(),
					OrganizationID: uuid.New(),
					IsPrimary:      true,
				},
			},
			expectedLen: 1,
			description: "单个有效成员应该返回长度为1的切片",
		},
		{
			name: "包含nil成员",
			input: []*models.UserMembership{
				{
					BaseModel: models.BaseModel{
						ID:        uuid.New(),
						CreatedAt: time.Now().UnixMilli(),
						UpdatedAt: time.Now().UnixMilli(),
					},
					UserID:         uuid.New(),
					OrganizationID: uuid.New(),
					IsPrimary:      true,
				},
				nil,
			},
			expectedLen: 2,
			description: "包含nil成员的切片应该保持原长度",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ModelUserMembershipsToThrift(tt.input)
			assert.Equal(t, tt.expectedLen, len(result), tt.description)
		})
	}
}