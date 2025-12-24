package organization

import (
	"fmt"
	"testing"

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

func TestConverterImpl_ModelOrganizationToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelOrganizationToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整组织转换", func(t *testing.T) {
		now := int64(1640995200000) // 2022-01-01 00:00:00 UTC
		orgID := uuid.New()
		parentID := uuid.New()

		model := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:                "HOSPITAL_001",
			Name:                "北京协和医院",
			ParentID:            parentID,
			FacilityType:        "综合医院",
			AccreditationStatus: "JCI认证",
			ProvinceCity:        models.StringSlice{"北京市", "东城区"},
		}

		result := converter.ModelOrganizationToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, orgID.String(), *result.ID)
		assert.Equal(t, "HOSPITAL_001", *result.Code)
		assert.Equal(t, "北京协和医院", *result.Name)
		require.NotNil(t, result.ParentID)
		assert.Equal(t, parentID.String(), *result.ParentID)
		require.NotNil(t, result.FacilityType)
		assert.Equal(t, "综合医院", *result.FacilityType)
		require.NotNil(t, result.AccreditationStatus)
		assert.Equal(t, "JCI认证", *result.AccreditationStatus)
		require.NotNil(t, result.ProvinceCity)
		assert.Equal(t, []string{"北京市", "东城区"}, result.ProvinceCity)
		assert.Equal(t, now, *result.CreatedAt)
		assert.Equal(t, now, *result.UpdatedAt)
	})

	t.Run("最小组织转换", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		model := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:         "ORG_001",
			Name:         "测试组织",
			ParentID:     uuid.Nil, // 根组织
			FacilityType: "",        // 空值
			ProvinceCity: models.StringSlice{}, // 空切片
		}

		result := converter.ModelOrganizationToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, orgID.String(), *result.ID)
		assert.Equal(t, "ORG_001", *result.Code)
		assert.Equal(t, "测试组织", *result.Name)
		assert.Nil(t, result.ParentID) // uuid.Nil应该被跳过
		assert.Nil(t, result.FacilityType) // 空字符串应该被跳过
		assert.Nil(t, result.AccreditationStatus) // 空字符串应该被跳过
		assert.Empty(t, result.ProvinceCity) // 空切片
		assert.Equal(t, now, *result.CreatedAt)
		assert.Equal(t, now, *result.UpdatedAt)
	})

	t.Run("无认证状态组织", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		model := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:         "CLINIC_001",
			Name:         "社区诊所",
			ParentID:     uuid.Nil,
			FacilityType: "社区诊所",
			// AccreditationStatus 为空
			ProvinceCity: models.StringSlice{"上海市"},
		}

		result := converter.ModelOrganizationToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, orgID.String(), *result.ID)
		assert.Equal(t, "CLINIC_001", *result.Code)
		assert.Equal(t, "社区诊所", *result.Name)
		assert.Nil(t, result.ParentID)
		require.NotNil(t, result.FacilityType)
		assert.Equal(t, "社区诊所", *result.FacilityType)
		assert.Nil(t, result.AccreditationStatus) // 空值应该被跳过
		require.NotNil(t, result.ProvinceCity)
		assert.Equal(t, []string{"上海市"}, result.ProvinceCity)
	})

	t.Run("多省市组织", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		model := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:         "CHAIN_001",
			Name:         "连锁集团",
			ParentID:     uuid.Nil,
			FacilityType: "医疗集团",
			ProvinceCity: models.StringSlice{"北京市", "上海市", "广东省", "深圳市"},
		}

		result := converter.ModelOrganizationToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, orgID.String(), *result.ID)
		assert.Equal(t, "CHAIN_001", *result.Code)
		assert.Equal(t, "连锁集团", *result.Name)
		assert.Nil(t, result.ParentID)
		require.NotNil(t, result.FacilityType)
		assert.Equal(t, "医疗集团", *result.FacilityType)
		require.NotNil(t, result.ProvinceCity)
		assert.Equal(t, []string{"北京市", "上海市", "广东省", "深圳市"}, result.ProvinceCity)
	})
}

func TestConverterImpl_ThriftOrganizationToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ThriftOrganizationToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整组织转换", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()
		parentID := uuid.New()

		orgIDStr := orgID.String()
		parentIDStr := parentID.String()
		code := "HOSPITAL_001"
		name := "北京协和医院"
		facilityType := "综合医院"
		accreditationStatus := "JCI认证"

		dto := &identity_srv.Organization{
			ID:                  &orgIDStr,
			Code:                &code,
			Name:                &name,
			ParentID:            &parentIDStr,
			FacilityType:        &facilityType,
			AccreditationStatus: &accreditationStatus,
			ProvinceCity:        []string{"北京市", "东城区"},
			CreatedAt:           &now,
			UpdatedAt:           &now,
		}

		result := converter.ThriftOrganizationToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, orgID, result.ID)
		assert.Equal(t, "HOSPITAL_001", result.Code)
		assert.Equal(t, "北京协和医院", result.Name)
		assert.Equal(t, parentID, result.ParentID)
		assert.Equal(t, "综合医院", result.FacilityType)
		assert.Equal(t, "JCI认证", result.AccreditationStatus)
		assert.Equal(t, models.StringSlice{"北京市", "东城区"}, result.ProvinceCity)
		// 注意：ThrfitOrganizationToModel 不设置时间字段，这些字段会在数据库层设置
		assert.Equal(t, int64(0), result.CreatedAt)
		assert.Equal(t, int64(0), result.UpdatedAt)
	})

	t.Run("最小组织转换", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		orgIDStr := orgID.String()
		code := "ORG_001"
		name := "测试组织"

		dto := &identity_srv.Organization{
			ID:        &orgIDStr,
			Code:      &code,
			Name:      &name,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		result := converter.ThriftOrganizationToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, orgID, result.ID)
		assert.Equal(t, "ORG_001", result.Code)
		assert.Equal(t, "测试组织", result.Name)
		assert.Equal(t, uuid.Nil, result.ParentID) // nil ParentID
		assert.Equal(t, "", result.FacilityType) // nil FacilityType
		assert.Equal(t, "", result.AccreditationStatus) // nil AccreditationStatus
		assert.Empty(t, result.ProvinceCity) // nil ProvinceCity
		// 注意：ThrfitOrganizationToModel 不设置时间字段，这些字段会在数据库层设置
		assert.Equal(t, int64(0), result.CreatedAt)
		assert.Equal(t, int64(0), result.UpdatedAt)
	})

	t.Run("根组织转换", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		orgIDStr := orgID.String()
		code := "ROOT_001"
		name := "根组织"
		facilityType := "总部"

		dto := &identity_srv.Organization{
			ID:           &orgIDStr,
			Code:         &code,
			Name:         &name,
			ParentID:     nil, // 根组织没有父组织
			FacilityType: &facilityType,
			ProvinceCity: []string{"北京市"},
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		result := converter.ThriftOrganizationToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, orgID, result.ID)
		assert.Equal(t, "ROOT_001", result.Code)
		assert.Equal(t, "根组织", result.Name)
		assert.Equal(t, uuid.Nil, result.ParentID) // nil ParentID 应该是 uuid.Nil
		assert.Equal(t, "总部", result.FacilityType)
		assert.Equal(t, models.StringSlice{"北京市"}, result.ProvinceCity)
	})

	t.Run("空省市列表处理", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		orgIDStr := orgID.String()
		code := "EMPTY_001"
		name := "空省市组织"

		dto := &identity_srv.Organization{
			ID:           &orgIDStr,
			Code:         &code,
			Name:         &name,
			ProvinceCity: []string{}, // 空列表
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		result := converter.ThriftOrganizationToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, orgID, result.ID)
		assert.Equal(t, "EMPTY_001", result.Code)
		assert.Equal(t, "空省市组织", result.Name)
		assert.Empty(t, result.ProvinceCity) // 空列表应该保持为空
	})
}

func TestConverterImpl_ModelOrganizationsToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ModelOrganizationsToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("空切片转换", func(t *testing.T) {
		result := converter.ModelOrganizationsToThrift([]*models.Organization{})
		assert.Nil(t, result)
	})

	t.Run("多个组织转换", func(t *testing.T) {
		now := int64(1640995200000)
		org1ID := uuid.New()
		org2ID := uuid.New()

		org1 := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        org1ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:     "ORG_001",
			Name:     "组织1",
			ParentID: uuid.Nil,
		}

		org2 := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        org2ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:     "ORG_002",
			Name:     "组织2",
			ParentID: uuid.Nil,
		}

		models := []*models.Organization{org1, org2}
		result := converter.ModelOrganizationsToThrift(models)

		require.NotNil(t, result)
		require.Len(t, result, 2)

		// 验证第一个组织
		assert.Equal(t, org1ID.String(), *result[0].ID)
		assert.Equal(t, "ORG_001", *result[0].Code)
		assert.Equal(t, "组织1", *result[0].Name)

		// 验证第二个组织
		assert.Equal(t, org2ID.String(), *result[1].ID)
		assert.Equal(t, "ORG_002", *result[1].Code)
		assert.Equal(t, "组织2", *result[1].Name)
	})

	t.Run("包含nil元素的切片", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		org := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:     "ORG_001",
			Name:     "测试组织",
			ParentID: uuid.Nil,
		}

		models := []*models.Organization{org, nil, org} // 包含nil元素
		result := converter.ModelOrganizationsToThrift(models)

		require.NotNil(t, result)
		require.Len(t, result, 2) // nil元素被过滤掉了

		// 第一个组织应该正常转换
		assert.Equal(t, orgID.String(), *result[0].ID)
		assert.Equal(t, "ORG_001", *result[0].Code)

		// 第二个组织应该正常转换
		assert.Equal(t, orgID.String(), *result[1].ID)
		assert.Equal(t, "ORG_001", *result[1].Code)
	})
}

func TestConverterImpl_CreateRequestToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.CreateRequestToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整创建请求转换", func(t *testing.T) {
		parentID := uuid.New()
		name := "新医院"
		facilityType := "专科医院"
		accreditationStatus := "CAP认证"
		provinceCity := []string{"北京市", "海淀区"}

		parentIDStr := parentID.String()
		req := &identity_srv.CreateOrganizationRequest{
			Name:                &name,
			ParentID:            &parentIDStr,
			FacilityType:        &facilityType,
			AccreditationStatus: &accreditationStatus,
			ProvinceCity:        provinceCity,
		}

		result := converter.CreateRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, "", result.Code) // Code应该在repository层生成
		assert.Equal(t, parentID, result.ParentID)
		assert.Equal(t, facilityType, result.FacilityType)
		assert.Equal(t, accreditationStatus, result.AccreditationStatus)
		assert.Equal(t, models.StringSlice(provinceCity), result.ProvinceCity)
	})

	t.Run("最小创建请求转换", func(t *testing.T) {
		name := "最小组织"

		req := &identity_srv.CreateOrganizationRequest{
			Name: &name,
		}

		result := converter.CreateRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, "", result.Code) // Code应该在repository层生成
		assert.Equal(t, uuid.Nil, result.ParentID) // nil ParentID
		assert.Equal(t, "", result.FacilityType) // nil FacilityType
		assert.Equal(t, "", result.AccreditationStatus) // nil AccreditationStatus
		assert.Empty(t, result.ProvinceCity) // nil ProvinceCity
	})

	t.Run("根组织创建请求", func(t *testing.T) {
		name := "根组织"
		facilityType := "总部"

		req := &identity_srv.CreateOrganizationRequest{
			Name:         &name,
			ParentID:     nil, // 根组织
			FacilityType: &facilityType,
		}

		result := converter.CreateRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, uuid.Nil, result.ParentID) // nil ParentID 应该是 uuid.Nil
		assert.Equal(t, facilityType, result.FacilityType)
		assert.Empty(t, result.ProvinceCity)
	})
}

func TestConverterImpl_ApplyUpdateToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入处理", func(t *testing.T) {
		// 两个都是nil
		result := converter.ApplyUpdateToModel(nil, nil)
		assert.Nil(t, result)

		// existing为nil
		result = converter.ApplyUpdateToModel(nil, &identity_srv.UpdateOrganizationRequest{})
		assert.Nil(t, result)

		// req为nil
		existing := &models.Organization{}
		result = converter.ApplyUpdateToModel(existing, nil)
		assert.Equal(t, existing, result) // 应该返回原对象
	})

	t.Run("完整更新应用", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()
		newParentID := uuid.New()

		existing := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:                "OLD_CODE",
			Name:                "旧名称",
			ParentID:            uuid.Nil,
			FacilityType:        "旧类型",
			AccreditationStatus: "旧状态",
			ProvinceCity:        models.StringSlice{"旧省市"},
		}

		newName := "新名称"
		facilityType := "新类型"
		accreditationStatus := "新状态"
		provinceCity := []string{"新省市1", "新省市2"}

		newParentIDStr := newParentID.String()
		req := &identity_srv.UpdateOrganizationRequest{
			Name:                &newName,
			ParentID:            &newParentIDStr,
			FacilityType:        &facilityType,
			AccreditationStatus: &accreditationStatus,
			ProvinceCity:        provinceCity,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.NotEqual(t, existing, result) // 应该是新的对象
		assert.Equal(t, orgID, result.ID) // ID不应该改变
		assert.Equal(t, "OLD_CODE", result.Code) // Code不应该改变
		assert.Equal(t, newName, result.Name) // Name应该更新
		assert.Equal(t, newParentID, result.ParentID) // ParentID应该更新
		assert.Equal(t, facilityType, result.FacilityType) // FacilityType应该更新
		assert.Equal(t, accreditationStatus, result.AccreditationStatus) // AccreditationStatus应该更新
		assert.Equal(t, models.StringSlice(provinceCity), result.ProvinceCity) // ProvinceCity应该更新
	})

	t.Run("部分更新应用", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		existing := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:                "ORG_001",
			Name:                "原始名称",
			ParentID:            uuid.Nil,
			FacilityType:        "原始类型",
			AccreditationStatus: "原始状态",
			ProvinceCity:        models.StringSlice{"原始省市"},
		}

		newName := "更新的名称"
		req := &identity_srv.UpdateOrganizationRequest{
			Name: &newName, // 只更新名称
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, orgID, result.ID)
		assert.Equal(t, "ORG_001", result.Code) // Code不应该改变
		assert.Equal(t, newName, result.Name) // Name应该更新
		assert.Equal(t, uuid.Nil, result.ParentID) // ParentID不应该改变
		assert.Equal(t, "原始类型", result.FacilityType) // FacilityType不应该改变
		assert.Equal(t, "原始状态", result.AccreditationStatus) // AccreditationStatus不应该改变
		assert.Equal(t, models.StringSlice{"原始省市"}, result.ProvinceCity) // ProvinceCity不应该改变
	})

	t.Run("空值更新处理", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		existing := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:                "ORG_001",
			Name:                "原始名称",
			ParentID:            uuid.New(), // 有父组织
			FacilityType:        "原始类型",
			AccreditationStatus: "原始状态",
			ProvinceCity:        models.StringSlice{"原始省市"},
		}

		req := &identity_srv.UpdateOrganizationRequest{
			ParentID:            nil, // nil ParentID
			FacilityType:        func() *string { s := ""; return &s }(), // 空字符串
			AccreditationStatus: nil, // nil
			ProvinceCity:        []string{}, // 空列表
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.NotEqual(t, existing, result) // 应该是新对象，不是原对象
		assert.Equal(t, orgID, result.ID)
		assert.Equal(t, "ORG_001", result.Code)
		assert.Equal(t, "原始名称", result.Name) // Name没有更新
		assert.Equal(t, existing.ParentID, result.ParentID) // ParentID 不应该改变，因为 req.ParentID 是 nil
		assert.Equal(t, "", result.FacilityType) // 空字符串
		assert.Equal(t, "原始状态", result.AccreditationStatus) // nil 不更新
		assert.Equal(t, models.StringSlice{}, result.ProvinceCity) // 空列表
	})
}

func TestConverterImpl_ModelToThrift_Alias(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("别名方法测试", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		model := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code: "ALIAS_TEST",
			Name: "别名测试",
		}

		// 使用别名方法
		result1 := converter.ModelToThrift(model)
		// 使用原方法
		result2 := converter.ModelOrganizationToThrift(model)

		// 结果应该相同
		assert.Equal(t, result1, result2)
		require.NotNil(t, result1)
		assert.Equal(t, orgID.String(), *result1.ID)
		assert.Equal(t, "ALIAS_TEST", *result1.Code)
		assert.Equal(t, "别名测试", *result1.Name)
	})
}

// 完整往返转换测试
func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("Model -> Thrift -> Model 往返转换", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()
		parentID := uuid.New()

		originalModel := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:                "ROUND_TRIP",
			Name:                "往返测试",
			ParentID:            parentID,
			FacilityType:        "测试类型",
			AccreditationStatus: "测试状态",
			ProvinceCity:        models.StringSlice{"测试省", "测试市"},
		}

		// Model -> Thrift
		thriftOrg := converter.ModelOrganizationToThrift(originalModel)
		require.NotNil(t, thriftOrg)

		// Thrift -> Model
		resultModel := converter.ThriftOrganizationToModel(thriftOrg)
		require.NotNil(t, resultModel)

		// 验证关键字段
		assert.Equal(t, originalModel.ID, resultModel.ID)
		assert.Equal(t, originalModel.Code, resultModel.Code)
		assert.Equal(t, originalModel.Name, resultModel.Name)
		assert.Equal(t, originalModel.ParentID, resultModel.ParentID)
		assert.Equal(t, originalModel.FacilityType, resultModel.FacilityType)
		assert.Equal(t, originalModel.AccreditationStatus, resultModel.AccreditationStatus)
		assert.Equal(t, originalModel.ProvinceCity, resultModel.ProvinceCity)
		// 注意：ThrfitOrganizationToModel 不设置时间字段，这些字段会在数据库层设置
		assert.Equal(t, int64(0), resultModel.CreatedAt)
		assert.Equal(t, int64(0), resultModel.UpdatedAt)
	})

	t.Run("Thrift -> Model -> Thrift 往返转换", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()
		parentID := uuid.New()

		orgIDStr := orgID.String()
		parentIDStr := parentID.String()
		code := "ROUND_TRIP2"
		name := "往返测试2"
		facilityType := "测试类型2"
		accreditationStatus := "测试状态2"

		originalThrift := &identity_srv.Organization{
			ID:                  &orgIDStr,
			Code:                &code,
			Name:                &name,
			ParentID:            &parentIDStr,
			FacilityType:        &facilityType,
			AccreditationStatus: &accreditationStatus,
			ProvinceCity:        []string{"测试省2", "测试市2"},
			CreatedAt:           &now,
			UpdatedAt:           &now,
		}

		// Thrift -> Model
		modelOrg := converter.ThriftOrganizationToModel(originalThrift)
		require.NotNil(t, modelOrg)

		// Model -> Thrift
		resultThrift := converter.ModelOrganizationToThrift(modelOrg)
		require.NotNil(t, resultThrift)

		// 验证关键字段
		assert.Equal(t, originalThrift.ID, resultThrift.ID)
		assert.Equal(t, originalThrift.Code, resultThrift.Code)
		assert.Equal(t, originalThrift.Name, resultThrift.Name)
		assert.Equal(t, originalThrift.ParentID, resultThrift.ParentID)
		assert.Equal(t, originalThrift.FacilityType, resultThrift.FacilityType)
		assert.Equal(t, originalThrift.AccreditationStatus, resultThrift.AccreditationStatus)
		assert.Equal(t, originalThrift.ProvinceCity, resultThrift.ProvinceCity)
		// 注意：ThrfitOrganizationToModel 不设置时间字段，所以往返转换后时间字段不会保持
	// Model -> Thrift -> Model -> Thrift 的往返转换会使用 Model 的时间字段（即默认值 0）
		assert.Equal(t, int64(0), *resultThrift.CreatedAt)
		assert.Equal(t, int64(0), *resultThrift.UpdatedAt)
	})
}

// 边界情况和错误处理测试
func TestConverterImpl_EdgeCases(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("空字符串字段处理", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		model := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        orgID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:                "EMPTY_TEST",
			Name:                "空字段测试",
			ParentID:            uuid.Nil,
			FacilityType:        "", // 空字符串
			AccreditationStatus: "", // 空字符串
			ProvinceCity:        models.StringSlice{}, // 空切片
		}

		result := converter.ModelOrganizationToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "EMPTY_TEST", *result.Code)
		assert.Equal(t, "空字段测试", *result.Name)
		assert.Nil(t, result.ParentID)
		assert.Nil(t, result.FacilityType) // 空字符串应该被跳过
		assert.Nil(t, result.AccreditationStatus) // 空字符串应该被跳过
		assert.Empty(t, result.ProvinceCity) // 空切片
	})

	t.Run("nil字段处理", func(t *testing.T) {
		now := int64(1640995200000)
		orgID := uuid.New()

		orgIDStr := orgID.String()
		code := "NIL_TEST"
		name := "nil字段测试"

		dto := &identity_srv.Organization{
			ID:        &orgIDStr,
			Code:      &code,
			Name:      &name,
			ParentID:  nil, // nil
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		result := converter.ThriftOrganizationToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, "NIL_TEST", result.Code)
		assert.Equal(t, "nil字段测试", result.Name)
		assert.Equal(t, uuid.Nil, result.ParentID) // nil ParentID 应该是 uuid.Nil
		assert.Equal(t, "", result.FacilityType) // nil 应该是空字符串
		assert.Equal(t, "", result.AccreditationStatus) // nil 应该是空字符串
		assert.Empty(t, result.ProvinceCity) // nil 应该是空切片
	})
}

// 基准测试
func BenchmarkConverterImpl_ModelOrganizationToThrift(b *testing.B) {
	converter := &ConverterImpl{}
	
	now := int64(1640995200000)
	orgID := uuid.New()
	parentID := uuid.New()

	model := &models.Organization{
		BaseModel: models.BaseModel{
			ID:        orgID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Code:                "BENCHMARK_TEST",
		Name:                "基准测试",
		ParentID:            parentID,
		FacilityType:        "测试类型",
		AccreditationStatus: "测试状态",
		ProvinceCity:        models.StringSlice{"测试省", "测试市"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.ModelOrganizationToThrift(model)
	}
}

func BenchmarkConverterImpl_ThriftOrganizationToModel(b *testing.B) {
	converter := &ConverterImpl{}
	
	now := int64(1640995200000)
	orgID := uuid.New()
	parentID := uuid.New()

	orgIDStr := orgID.String()
			parentIDStr := parentID.String()
			code := "BENCHMARK_TEST"
			name := "基准测试"
			facilityType := "测试类型"
			accreditationStatus := "测试状态"
	
			dto := &identity_srv.Organization{
				ID:                  &orgIDStr,
				Code:                &code,
				Name:                &name,
				ParentID:            &parentIDStr,
				FacilityType:        &facilityType,
				AccreditationStatus: &accreditationStatus,
				ProvinceCity:        []string{"测试省", "测试市"},
				CreatedAt:           &now,
				UpdatedAt:           &now,
			}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.ThriftOrganizationToModel(dto)
	}
}

func BenchmarkConverterImpl_ModelOrganizationsToThrift(b *testing.B) {
	converter := &ConverterImpl{}
	
	now := int64(1640995200000)

	var orgs []*models.Organization
	for i := 0; i < 10; i++ {
		org := &models.Organization{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:     fmt.Sprintf("ORG_%d", i),
			Name:     fmt.Sprintf("组织 %d", i),
			ParentID: uuid.Nil,
		}
		orgs = append(orgs, org)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.ModelOrganizationsToThrift(orgs)
	}
}