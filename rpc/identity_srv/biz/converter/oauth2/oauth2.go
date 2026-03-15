package oauth2

import (
	"strings"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter OAuth2 数据转换器接口
type Converter interface {
	// Client 转换
	ClientModelToThrift(client *models.OAuth2Client) *identity_srv.OAuth2Client
	ClientModelsToThrift(clients []*models.OAuth2Client) []*identity_srv.OAuth2Client
	CreateClientRequestToModel(req *identity_srv.CreateOAuth2ClientRequest) *models.OAuth2Client

	// Scope 转换
	ScopeModelToThrift(scope *models.OAuth2Scope) *identity_srv.OAuth2Scope
	ScopeModelsToThrift(scopes []*models.OAuth2Scope) []*identity_srv.OAuth2Scope

	// Consent 转换
	ConsentModelToThrift(consent *models.OAuth2Consent) *identity_srv.OAuth2Consent
	ConsentModelsToThrift(consents []*models.OAuth2Consent) []*identity_srv.OAuth2Consent
}

type converterImpl struct{}

// NewConverter 创建 OAuth2 转换器实例
func NewConverter() Converter {
	return &converterImpl{}
}

// --- Client ---

func (c *converterImpl) ClientModelToThrift(client *models.OAuth2Client) *identity_srv.OAuth2Client {
	if client == nil {
		return nil
	}

	id := client.ID.String()
	createdAt := client.CreatedAt
	updatedAt := client.UpdatedAt
	isActive := client.IsActive
	accessLifespan := int32(client.AccessTokenLifespan)
	refreshLifespan := int32(client.RefreshTokenLifespan)

	result := &identity_srv.OAuth2Client{
		Id:                   &id,
		ClientID:             &client.ClientID,
		ClientName:           &client.ClientName,
		Description:          &client.Description,
		ClientType:           clientTypeToThrift(client.ClientType),
		GrantTypes:           grantTypesToThrift(client.GrantTypes),
		RedirectURIs:         []string(client.RedirectURIs),
		Scopes:               []string(client.Scopes),
		LogoURI:              &client.LogoURI,
		ClientURI:            &client.ClientURI,
		AccessTokenLifespan:  &accessLifespan,
		RefreshTokenLifespan: &refreshLifespan,
		IsActive:             &isActive,
		CreatedAt:            &createdAt,
		UpdatedAt:            &updatedAt,
	}

	if client.OwnerID != nil {
		ownerStr := client.OwnerID.String()
		result.OwnerID = &ownerStr
	}

	return result
}

func (c *converterImpl) ClientModelsToThrift(clients []*models.OAuth2Client) []*identity_srv.OAuth2Client {
	result := make([]*identity_srv.OAuth2Client, 0, len(clients))
	for _, client := range clients {
		result = append(result, c.ClientModelToThrift(client))
	}

	return result
}

func (c *converterImpl) CreateClientRequestToModel(req *identity_srv.CreateOAuth2ClientRequest) *models.OAuth2Client {
	if req == nil {
		return nil
	}

	client := &models.OAuth2Client{
		ClientName:   ptrStr(req.ClientName),
		Description:  ptrStr(req.Description),
		ClientType:   clientTypeFromThrift(req.ClientType),
		GrantTypes:   grantTypesFromThrift(req.GrantTypes),
		Scopes:       models.StringSlice(req.Scopes),
		RedirectURIs: models.StringSlice(req.RedirectURIs),
		LogoURI:      ptrStr(req.LogoURI),
		ClientURI:    ptrStr(req.ClientURI),
		IsActive:     true,
	}

	if req.AccessTokenLifespan != nil {
		client.AccessTokenLifespan = int(*req.AccessTokenLifespan)
	}

	if req.RefreshTokenLifespan != nil {
		client.RefreshTokenLifespan = int(*req.RefreshTokenLifespan)
	}

	return client
}

// --- Scope ---

func (c *converterImpl) ScopeModelToThrift(scope *models.OAuth2Scope) *identity_srv.OAuth2Scope {
	if scope == nil {
		return nil
	}

	id := scope.ID.String()
	isDefault := scope.IsDefault
	isSystem := scope.IsSystem
	createdAt := scope.CreatedAt
	updatedAt := scope.UpdatedAt

	return &identity_srv.OAuth2Scope{
		Id:          &id,
		Name:        &scope.Name,
		DisplayName: &scope.DisplayName,
		Description: &scope.Description,
		IsDefault:   &isDefault,
		IsSystem:    &isSystem,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}
}

func (c *converterImpl) ScopeModelsToThrift(scopes []*models.OAuth2Scope) []*identity_srv.OAuth2Scope {
	result := make([]*identity_srv.OAuth2Scope, 0, len(scopes))
	for _, scope := range scopes {
		result = append(result, c.ScopeModelToThrift(scope))
	}

	return result
}

// --- Consent ---

func (c *converterImpl) ConsentModelToThrift(consent *models.OAuth2Consent) *identity_srv.OAuth2Consent {
	if consent == nil {
		return nil
	}

	id := consent.ID.String()
	userID := consent.UserID.String()
	grantedAt := consent.GrantedAt
	isRevoked := consent.IsRevoked
	createdAt := consent.CreatedAt
	updatedAt := consent.UpdatedAt

	scopeList := strings.Split(consent.Scopes, " ")

	return &identity_srv.OAuth2Consent{
		Id:        &id,
		UserID:    &userID,
		ClientID:  &consent.ClientID,
		Scopes:    scopeList,
		GrantedAt: &grantedAt,
		IsRevoked: &isRevoked,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

func (c *converterImpl) ConsentModelsToThrift(consents []*models.OAuth2Consent) []*identity_srv.OAuth2Consent {
	result := make([]*identity_srv.OAuth2Consent, 0, len(consents))
	for _, consent := range consents {
		result = append(result, c.ConsentModelToThrift(consent))
	}

	return result
}

// --- 辅助函数 ---

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func clientTypeToThrift(t models.OAuth2ClientType) *identity_srv.OAuth2ClientType {
	var result identity_srv.OAuth2ClientType

	switch t {
	case models.OAuth2ClientTypeConfidential:
		result = identity_srv.OAuth2ClientType_CONFIDENTIAL
	case models.OAuth2ClientTypePublic:
		result = identity_srv.OAuth2ClientType_PUBLIC
	default:
		result = identity_srv.OAuth2ClientType_CONFIDENTIAL
	}

	return &result
}

func clientTypeFromThrift(t *identity_srv.OAuth2ClientType) models.OAuth2ClientType {
	if t == nil {
		return models.OAuth2ClientTypeConfidential
	}

	switch *t {
	case identity_srv.OAuth2ClientType_PUBLIC:
		return models.OAuth2ClientTypePublic
	default:
		return models.OAuth2ClientTypeConfidential
	}
}

func grantTypesToThrift(types models.StringSlice) []identity_srv.OAuth2GrantType {
	result := make([]identity_srv.OAuth2GrantType, 0, len(types))

	for _, t := range types {
		result = append(result, stringToGrantType(t))
	}

	return result
}

func grantTypesFromThrift(types []identity_srv.OAuth2GrantType) models.StringSlice {
	result := make(models.StringSlice, 0, len(types))

	for _, t := range types {
		result = append(result, grantTypeToString(t))
	}

	return result
}

func stringToGrantType(s string) identity_srv.OAuth2GrantType {
	switch s {
	case "authorization_code":
		return identity_srv.OAuth2GrantType_AUTHORIZATION_CODE
	case "client_credentials":
		return identity_srv.OAuth2GrantType_CLIENT_CREDENTIALS
	case "refresh_token":
		return identity_srv.OAuth2GrantType_REFRESH_TOKEN
	default:
		return identity_srv.OAuth2GrantType_AUTHORIZATION_CODE
	}
}

func grantTypeToString(t identity_srv.OAuth2GrantType) string {
	switch t {
	case identity_srv.OAuth2GrantType_AUTHORIZATION_CODE:
		return "authorization_code"
	case identity_srv.OAuth2GrantType_CLIENT_CREDENTIALS:
		return "client_credentials"
	case identity_srv.OAuth2GrantType_REFRESH_TOKEN:
		return "refresh_token"
	default:
		return "authorization_code"
	}
}
