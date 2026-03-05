package converter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type ConverterTestSuite struct {
	suite.Suite
	conv Converter
}

func (s *ConverterTestSuite) SetupTest() {
	s.conv = NewConverter()
}

func (s *ConverterTestSuite) TestNewConverter_Success() {
	conv := NewConverter()

	require.NotNil(s.T(), conv)
	assert.NotNil(s.T(), conv.UserProfile())
	assert.NotNil(s.T(), conv.Membership())
	assert.NotNil(s.T(), conv.Authentication())
	assert.NotNil(s.T(), conv.Organization())
	assert.NotNil(s.T(), conv.Department())
	assert.NotNil(s.T(), conv.Logo())
	assert.NotNil(s.T(), conv.Menu())
	assert.NotNil(s.T(), conv.RoleDefinition())
	assert.NotNil(s.T(), conv.UserRoleAssignment())
	assert.NotNil(s.T(), conv.Enum())
	assert.NotNil(s.T(), conv.Base())
}

func (s *ConverterTestSuite) TestUserProfile_Accessor() {
	converter := s.conv.UserProfile()

	require.NotNil(s.T(), converter)
	// 验证返回的是同一个实例（单例模式）
	assert.Same(s.T(), converter, s.conv.UserProfile())
}

func (s *ConverterTestSuite) TestMembership_Accessor() {
	converter := s.conv.Membership()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Membership())
}

func (s *ConverterTestSuite) TestAuthentication_Accessor() {
	converter := s.conv.Authentication()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Authentication())
}

func (s *ConverterTestSuite) TestOrganization_Accessor() {
	converter := s.conv.Organization()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Organization())
}

func (s *ConverterTestSuite) TestDepartment_Accessor() {
	converter := s.conv.Department()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Department())
}

func (s *ConverterTestSuite) TestLogo_Accessor() {
	converter := s.conv.Logo()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Logo())
}

func (s *ConverterTestSuite) TestMenu_Accessor() {
	converter := s.conv.Menu()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Menu())
}

func (s *ConverterTestSuite) TestRoleDefinition_Accessor() {
	converter := s.conv.RoleDefinition()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.RoleDefinition())
}

func (s *ConverterTestSuite) TestUserRoleAssignment_Accessor() {
	converter := s.conv.UserRoleAssignment()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.UserRoleAssignment())
}

func (s *ConverterTestSuite) TestEnum_Accessor() {
	converter := s.conv.Enum()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Enum())
}

func (s *ConverterTestSuite) TestBase_Accessor() {
	converter := s.conv.Base()

	require.NotNil(s.T(), converter)
	assert.Same(s.T(), converter, s.conv.Base())
}

func (s *ConverterTestSuite) TestBuildLoginResponse_CompleteUser() {
	// 创建用户档案
	userID := uuid.New()
	userProfile := &models.UserProfile{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		Username: "testuser",
		RealName: "测试用户",
		Email:    "test@example.com",
		Status:   models.UserStatusActive,
	}

	// 创建成员关系
	orgID := uuid.New()
	deptID := uuid.New()
	memberships := []*models.UserMembership{
		{
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   deptID,
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		},
	}

	// 构建登录响应
	response := s.conv.BuildLoginResponse(userProfile, memberships)

	// 验证响应
	require.NotNil(s.T(), response)
	assert.NotNil(s.T(), response.UserProfile)
	assert.NotNil(s.T(), response.UserProfile.Username)
	assert.Equal(s.T(), "testuser", *response.UserProfile.Username)
	assert.NotNil(s.T(), response.UserProfile.RealName)
	assert.Equal(s.T(), "测试用户", *response.UserProfile.RealName)
	assert.NotNil(s.T(), response.Memberships)
	assert.Len(s.T(), response.Memberships, 1)
}

func (s *ConverterTestSuite) TestBuildLoginResponse_UserWithoutMemberships() {
	userID := uuid.New()
	userProfile := &models.UserProfile{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		Username: "testuser",
		Status:   models.UserStatusActive,
	}

	response := s.conv.BuildLoginResponse(userProfile, nil)

	require.NotNil(s.T(), response)
	assert.NotNil(s.T(), response.UserProfile)
	assert.NotNil(s.T(), response.Memberships)
	assert.Len(s.T(), response.Memberships, 0)
}

func (s *ConverterTestSuite) TestBuildLoginResponse_UserWithMultipleMemberships() {
	userID := uuid.New()
	org1ID := uuid.New()
	org2ID := uuid.New()

	userProfile := &models.UserProfile{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		Username: "multiorg",
		Status:   models.UserStatusActive,
	}

	memberships := []*models.UserMembership{
		{
			UserID:         userID,
			OrganizationID: org1ID,
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		},
		{
			UserID:         userID,
			OrganizationID: org2ID,
			IsPrimary:      false,
			Status:         models.MembershipStatusActive,
		},
	}

	response := s.conv.BuildLoginResponse(userProfile, memberships)

	require.NotNil(s.T(), response)
	assert.NotNil(s.T(), response.UserProfile)
	assert.Len(s.T(), response.Memberships, 2)
}

func (s *ConverterTestSuite) TestImpl_AllConvertersInitialized() {
	impl, ok := s.conv.(*Impl)
	require.True(s.T(), ok, "Converter 应该是 *Impl 类型")

	// 验证所有子转换器都已初始化
	assert.NotNil(s.T(), impl.authenticationConverter)
	assert.NotNil(s.T(), impl.userProfileConverter)
	assert.NotNil(s.T(), impl.membershipConverter)
	assert.NotNil(s.T(), impl.organizationConverter)
	assert.NotNil(s.T(), impl.departmentConverter)
	assert.NotNil(s.T(), impl.logoConverter)
	assert.NotNil(s.T(), impl.menuConverter)
	assert.NotNil(s.T(), impl.roleDefinitionConverter)
	assert.NotNil(s.T(), impl.userRoleAssignmentConverter)
	assert.NotNil(s.T(), impl.enumConverter)
	assert.NotNil(s.T(), impl.baseConverter)
}

func TestConverterSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}
