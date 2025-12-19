package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
)

// ============================================================================
// 聚合组装器实现 - 统一访问入口
// ============================================================================

// identityAssembler 身份管理聚合组装器实现
type identityAssembler struct {
	authAssembler       IAuthAssembler
	userAssembler       IUserAssembler
	orgAssembler        IOrgAssembler
	departmentAssembler IDepartmentAssembler
	membershipAssembler IMembershipAssembler
	logoAssembler       ILogoAssembler
}

// NewIdentityAggregateAssembler 创建身份管理聚合组装器
func NewIdentityAggregateAssembler(
	authAssembler IAuthAssembler,
	userAssembler IUserAssembler,
	orgAssembler IOrgAssembler,
	departmentAssembler IDepartmentAssembler,
	membershipAssembler IMembershipAssembler,
	logoAssembler ILogoAssembler,
) Assembler {
	return &identityAssembler{
		authAssembler:       authAssembler,
		userAssembler:       userAssembler,
		orgAssembler:        orgAssembler,
		departmentAssembler: departmentAssembler,
		membershipAssembler: membershipAssembler,
		logoAssembler:       logoAssembler,
	}
}

// 获取各个业务领域的组装器
func (a *identityAssembler) Auth() IAuthAssembler             { return a.authAssembler }
func (a *identityAssembler) User() IUserAssembler             { return a.userAssembler }
func (a *identityAssembler) Organization() IOrgAssembler      { return a.orgAssembler }
func (a *identityAssembler) Department() IDepartmentAssembler { return a.departmentAssembler }
func (a *identityAssembler) Membership() IMembershipAssembler { return a.membershipAssembler }
func (a *identityAssembler) Logo() ILogoAssembler             { return a.logoAssembler }

// 通用转换方法
func (a *identityAssembler) ToHTTPPageResponse(
	rpc *rpc_base.PageResponse,
) *http_base.PageResponseDTO {
	return ToHTTPPageResponse(rpc)
}

func (a *identityAssembler) ToRPCPageRequest(http *http_base.PageRequestDTO) *rpc_base.PageRequest {
	return ToRPCPageRequest(http)
}
