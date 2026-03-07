package mock

import (
	"context"

	"go.uber.org/mock/gomock"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
)

// TestMocks 聚合所有 Mock 对象，方便测试代码引用
type TestMocks struct {
	Ctrl      *gomock.Controller
	DAL       *MockDAL
	Converter converter.Converter // 使用真实 Converter（纯函数，无副作用）

	// 子仓储 Mock
	UserRepo       *MockUserProfileRepository
	OrgRepo        *MockOrganizationRepository
	DeptRepo       *MockDepartmentRepository
	MembershipRepo *MockUserMembershipRepository
	DefinitionRepo *MockRoleDefinitionRepository
	AssignmentRepo *MockUserRoleAssignmentRepository
	MenuRepo       *MockMenuRepository
	RoleMenuRepo   *MockRoleMenuPermissionRepository
	LogoRepo       *MockLogoRepository
}

// NewTestMocks 创建完整的测试 Mock 环境
//
// 自动完成以下设置：
//   - 创建 MockDAL 和所有子仓储 Mock
//   - 配置 MockDAL 的子仓储访问方法返回对应的 Mock
//   - 配置 WithTransaction 直接执行回调函数（使用同一个 MockDAL）
//   - 创建真实的 Converter 实例
func NewTestMocks(ctrl *gomock.Controller) *TestMocks {
	m := &TestMocks{
		Ctrl:      ctrl,
		Converter: converter.NewConverter(),

		DAL:            NewMockDAL(ctrl),
		UserRepo:       NewMockUserProfileRepository(ctrl),
		OrgRepo:        NewMockOrganizationRepository(ctrl),
		DeptRepo:       NewMockDepartmentRepository(ctrl),
		MembershipRepo: NewMockUserMembershipRepository(ctrl),
		DefinitionRepo: NewMockRoleDefinitionRepository(ctrl),
		AssignmentRepo: NewMockUserRoleAssignmentRepository(ctrl),
		MenuRepo:       NewMockMenuRepository(ctrl),
		RoleMenuRepo:   NewMockRoleMenuPermissionRepository(ctrl),
		LogoRepo:       NewMockLogoRepository(ctrl),
	}

	// 配置 DAL 子仓储访问方法（AnyTimes 避免测试中每次都需要 EXPECT）
	m.DAL.EXPECT().UserProfile().Return(m.UserRepo).AnyTimes()
	m.DAL.EXPECT().Organization().Return(m.OrgRepo).AnyTimes()
	m.DAL.EXPECT().Department().Return(m.DeptRepo).AnyTimes()
	m.DAL.EXPECT().UserMembership().Return(m.MembershipRepo).AnyTimes()
	m.DAL.EXPECT().RoleDefinition().Return(m.DefinitionRepo).AnyTimes()
	m.DAL.EXPECT().UserRoleAssignment().Return(m.AssignmentRepo).AnyTimes()
	m.DAL.EXPECT().Menu().Return(m.MenuRepo).AnyTimes()
	m.DAL.EXPECT().RoleMenuPermission().Return(m.RoleMenuRepo).AnyTimes()
	m.DAL.EXPECT().Logo().Return(m.LogoRepo).AnyTimes()

	// 配置 WithTransaction：直接执行回调函数，使用同一个 MockDAL
	m.DAL.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context, dal.DAL) error) error {
			return fn(ctx, m.DAL)
		}).AnyTimes()

	return m
}
