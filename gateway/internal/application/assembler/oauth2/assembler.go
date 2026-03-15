package oauth2

import (
	identity_srv "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/oauth2"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	rpc "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Assembler OAuth2 数据装配器接口
type Assembler interface {
	// Client 转换
	ToRPCCreateClientRequest(req *identity_srv.CreateOAuth2ClientRequestDTO) *rpc.CreateOAuth2ClientRequest
	ToHTTPClient(rpcClient *rpc.OAuth2Client) *identity_srv.OAuth2ClientDTO
	ToHTTPClients(rpcClients []*rpc.OAuth2Client) []*identity_srv.OAuth2ClientDTO

	// Scope 转换
	ToHTTPScope(rpcScope *rpc.OAuth2Scope) *identity_srv.OAuth2ScopeDTO
	ToHTTPScopes(rpcScopes []*rpc.OAuth2Scope) []*identity_srv.OAuth2ScopeDTO

	// Consent 转换
	ToHTTPConsent(rpcConsent *rpc.OAuth2Consent) *identity_srv.OAuth2ConsentDTO
	ToHTTPConsents(rpcConsents []*rpc.OAuth2Consent) []*identity_srv.OAuth2ConsentDTO
}

type assemblerImpl struct{}

func NewAssembler() Assembler {
	return &assemblerImpl{}
}

// --- Client ---

func (a *assemblerImpl) ToRPCCreateClientRequest(
	req *identity_srv.CreateOAuth2ClientRequestDTO,
) *rpc.CreateOAuth2ClientRequest {
	if req == nil {
		return nil
	}

	rpcReq := &rpc.CreateOAuth2ClientRequest{
		ClientName:   req.ClientName,
		Description:  req.Description,
		ClientType:   req.ClientType,
		GrantTypes:   req.GrantTypes,
		RedirectURIs: req.RedirectURIs,
		Scopes:       req.Scopes,
		LogoURI:      req.LogoURI,
		ClientURI:    req.ClientURI,
	}

	if req.AccessTokenLifespan != nil {
		rpcReq.AccessTokenLifespan = req.AccessTokenLifespan
	}

	if req.RefreshTokenLifespan != nil {
		rpcReq.RefreshTokenLifespan = req.RefreshTokenLifespan
	}

	return rpcReq
}

func (a *assemblerImpl) ToHTTPClient(rpcClient *rpc.OAuth2Client) *identity_srv.OAuth2ClientDTO {
	if rpcClient == nil {
		return nil
	}

	return &identity_srv.OAuth2ClientDTO{
		ID:                   common.CopyStringPtr(rpcClient.Id),
		ClientID:             common.CopyStringPtr(rpcClient.ClientID),
		ClientName:           common.CopyStringPtr(rpcClient.ClientName),
		Description:          common.CopyStringPtr(rpcClient.Description),
		ClientType:           common.CopyStringPtr(rpcClient.ClientType),
		GrantTypes:           rpcClient.GrantTypes,
		RedirectURIs:         rpcClient.RedirectURIs,
		Scopes:               rpcClient.Scopes,
		LogoURI:              common.CopyStringPtr(rpcClient.LogoURI),
		ClientURI:            common.CopyStringPtr(rpcClient.ClientURI),
		AccessTokenLifespan:  common.CopyInt32Ptr(rpcClient.AccessTokenLifespan),
		RefreshTokenLifespan: common.CopyInt32Ptr(rpcClient.RefreshTokenLifespan),
		IsActive:             common.CopyBoolPtr(rpcClient.IsActive),
		OwnerID:              common.CopyStringPtr(rpcClient.OwnerID),
		CreatedAt:            common.CopyInt64Ptr(rpcClient.CreatedAt),
		UpdatedAt:            common.CopyInt64Ptr(rpcClient.UpdatedAt),
	}
}

func (a *assemblerImpl) ToHTTPClients(rpcClients []*rpc.OAuth2Client) []*identity_srv.OAuth2ClientDTO {
	result := make([]*identity_srv.OAuth2ClientDTO, 0, len(rpcClients))
	for _, c := range rpcClients {
		result = append(result, a.ToHTTPClient(c))
	}

	return result
}

// --- Scope ---

func (a *assemblerImpl) ToHTTPScope(rpcScope *rpc.OAuth2Scope) *identity_srv.OAuth2ScopeDTO {
	if rpcScope == nil {
		return nil
	}

	return &identity_srv.OAuth2ScopeDTO{
		ID:          common.CopyStringPtr(rpcScope.Id),
		Name:        common.CopyStringPtr(rpcScope.Name),
		DisplayName: common.CopyStringPtr(rpcScope.DisplayName),
		Description: common.CopyStringPtr(rpcScope.Description),
		IsDefault:   common.CopyBoolPtr(rpcScope.IsDefault),
		IsSystem:    common.CopyBoolPtr(rpcScope.IsSystem),
	}
}

func (a *assemblerImpl) ToHTTPScopes(rpcScopes []*rpc.OAuth2Scope) []*identity_srv.OAuth2ScopeDTO {
	result := make([]*identity_srv.OAuth2ScopeDTO, 0, len(rpcScopes))
	for _, s := range rpcScopes {
		result = append(result, a.ToHTTPScope(s))
	}

	return result
}

// --- Consent ---

func (a *assemblerImpl) ToHTTPConsent(rpcConsent *rpc.OAuth2Consent) *identity_srv.OAuth2ConsentDTO {
	if rpcConsent == nil {
		return nil
	}

	return &identity_srv.OAuth2ConsentDTO{
		ID:         common.CopyStringPtr(rpcConsent.Id),
		ClientID:   common.CopyStringPtr(rpcConsent.ClientID),
		ClientName: common.CopyStringPtr(rpcConsent.ClientName),
		Scopes:     rpcConsent.Scopes,
		GrantedAt:  common.CopyInt64Ptr(rpcConsent.GrantedAt),
		IsRevoked:  common.CopyBoolPtr(rpcConsent.IsRevoked),
	}
}

func (a *assemblerImpl) ToHTTPConsents(rpcConsents []*rpc.OAuth2Consent) []*identity_srv.OAuth2ConsentDTO {
	result := make([]*identity_srv.OAuth2ConsentDTO, 0, len(rpcConsents))
	for _, c := range rpcConsents {
		result = append(result, a.ToHTTPConsent(c))
	}

	return result
}
