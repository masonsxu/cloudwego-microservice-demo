package oauth2

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// OAuth2Logic OAuth2 业务逻辑接口
type OAuth2Logic interface {
	// Client 管理
	CreateOAuth2Client(
		ctx context.Context,
		req *identity_srv.CreateOAuth2ClientRequest,
	) (*identity_srv.CreateOAuth2ClientResponse, error)
	GetOAuth2Client(ctx context.Context, req *identity_srv.GetOAuth2ClientRequest) (*identity_srv.OAuth2Client, error)
	UpdateOAuth2Client(
		ctx context.Context,
		req *identity_srv.UpdateOAuth2ClientRequest,
	) (*identity_srv.OAuth2Client, error)
	DeleteOAuth2Client(ctx context.Context, req *identity_srv.DeleteOAuth2ClientRequest) error
	ListOAuth2Clients(
		ctx context.Context,
		req *identity_srv.ListOAuth2ClientsRequest,
	) (*identity_srv.ListOAuth2ClientsResponse, error)
	RotateOAuth2ClientSecret(
		ctx context.Context,
		req *identity_srv.RotateOAuth2ClientSecretRequest,
	) (*identity_srv.RotateOAuth2ClientSecretResponse, error)
	GetOAuth2ClientForAuth(ctx context.Context, clientID string) (*identity_srv.GetOAuth2ClientForAuthResponse, error)

	// Scope 管理
	ListOAuth2Scopes(
		ctx context.Context,
		req *identity_srv.ListOAuth2ScopesRequest,
	) (*identity_srv.ListOAuth2ScopesResponse, error)

	// Consent 管理
	SaveOAuth2Consent(ctx context.Context, req *identity_srv.SaveOAuth2ConsentRequest) error
	GetOAuth2Consent(
		ctx context.Context,
		req *identity_srv.GetOAuth2ConsentRequest,
	) (*identity_srv.GetOAuth2ConsentResponse, error)
	ListOAuth2Consents(
		ctx context.Context,
		req *identity_srv.ListOAuth2ConsentsRequest,
	) (*identity_srv.ListOAuth2ConsentsResponse, error)
	RevokeOAuth2Consent(ctx context.Context, req *identity_srv.RevokeOAuth2ConsentRequest) error
}

type logicImpl struct {
	dal  dal.DAL
	conv converter.Converter
}

// NewLogic 创建 OAuth2 业务逻辑实例
func NewLogic(dal dal.DAL, conv converter.Converter) OAuth2Logic {
	return &logicImpl{dal: dal, conv: conv}
}

// --- Client 管理 ---

func (l *logicImpl) CreateOAuth2Client(
	ctx context.Context,
	req *identity_srv.CreateOAuth2ClientRequest,
) (*identity_srv.CreateOAuth2ClientResponse, error) {
	client := l.conv.OAuth2().CreateClientRequestToModel(req)

	if req.OwnerID != nil {
		ownerUUID, err := uuid.Parse(*req.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("无效的 ownerID: %w", err)
		}

		client.OwnerID = &ownerUUID
	}

	plainSecret, err := client.GenerateSecret()
	if err != nil {
		return nil, fmt.Errorf("生成客户端密钥失败: %w", err)
	}

	if err := l.dal.OAuth2Client().Create(ctx, client); err != nil {
		return nil, fmt.Errorf("创建 OAuth2 客户端失败: %w", err)
	}

	return &identity_srv.CreateOAuth2ClientResponse{
		Client:       l.conv.OAuth2().ClientModelToThrift(client),
		ClientSecret: &plainSecret,
	}, nil
}

func (l *logicImpl) GetOAuth2Client(
	ctx context.Context,
	req *identity_srv.GetOAuth2ClientRequest,
) (*identity_srv.OAuth2Client, error) {
	if req.Id == nil {
		return nil, fmt.Errorf("客户端 ID 不能为空")
	}

	client, err := l.dal.OAuth2Client().GetByID(ctx, *req.Id)
	if err != nil {
		return nil, err
	}

	return l.conv.OAuth2().ClientModelToThrift(client), nil
}

func (l *logicImpl) UpdateOAuth2Client(
	ctx context.Context,
	req *identity_srv.UpdateOAuth2ClientRequest,
) (*identity_srv.OAuth2Client, error) {
	if req.Id == nil {
		return nil, fmt.Errorf("客户端 ID 不能为空")
	}

	client, err := l.dal.OAuth2Client().GetByID(ctx, *req.Id)
	if err != nil {
		return nil, err
	}

	if req.ClientName != nil {
		client.ClientName = *req.ClientName
	}

	if req.Description != nil {
		client.Description = *req.Description
	}

	if req.RedirectURIs != nil {
		client.RedirectURIs = models.StringSlice(req.RedirectURIs)
	}

	if req.Scopes != nil {
		client.Scopes = models.StringSlice(req.Scopes)
	}

	if req.LogoURI != nil {
		client.LogoURI = *req.LogoURI
	}

	if req.ClientURI != nil {
		client.ClientURI = *req.ClientURI
	}

	if req.AccessTokenLifespan != nil {
		client.AccessTokenLifespan = int(*req.AccessTokenLifespan)
	}

	if req.RefreshTokenLifespan != nil {
		client.RefreshTokenLifespan = int(*req.RefreshTokenLifespan)
	}

	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	if err := l.dal.OAuth2Client().Update(ctx, client); err != nil {
		return nil, fmt.Errorf("更新 OAuth2 客户端失败: %w", err)
	}

	return l.conv.OAuth2().ClientModelToThrift(client), nil
}

func (l *logicImpl) DeleteOAuth2Client(ctx context.Context, req *identity_srv.DeleteOAuth2ClientRequest) error {
	if req.Id == nil {
		return fmt.Errorf("客户端 ID 不能为空")
	}

	return l.dal.OAuth2Client().Delete(ctx, *req.Id)
}

func (l *logicImpl) ListOAuth2Clients(
	ctx context.Context,
	req *identity_srv.ListOAuth2ClientsRequest,
) (*identity_srv.ListOAuth2ClientsResponse, error) {
	page, limit := int32(1), int32(20)
	if req.Page != nil {
		page = req.Page.GetPage()
		limit = req.Page.GetLimit()
	}

	var ownerID *string
	if req.OwnerID != nil {
		ownerID = req.OwnerID
	}

	clients, total, err := l.dal.OAuth2Client().List(ctx, ownerID, req.IsActive, int(page), int(limit))
	if err != nil {
		return nil, err
	}

	totalI32 := int32(total)
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &identity_srv.ListOAuth2ClientsResponse{
		Clients: l.conv.OAuth2().ClientModelsToThrift(clients),
		Page: &rpc_base.PageResponse{
			Total:      &totalI32,
			Page:       &page,
			Limit:      &limit,
			TotalPages: &totalPages,
		},
	}, nil
}

func (l *logicImpl) RotateOAuth2ClientSecret(
	ctx context.Context,
	req *identity_srv.RotateOAuth2ClientSecretRequest,
) (*identity_srv.RotateOAuth2ClientSecretResponse, error) {
	if req.Id == nil {
		return nil, fmt.Errorf("客户端 ID 不能为空")
	}

	client, err := l.dal.OAuth2Client().GetByID(ctx, *req.Id)
	if err != nil {
		return nil, err
	}

	plainSecret, err := client.GenerateSecret()
	if err != nil {
		return nil, fmt.Errorf("生成新密钥失败: %w", err)
	}

	if err := l.dal.OAuth2Client().Update(ctx, client); err != nil {
		return nil, fmt.Errorf("更新客户端密钥失败: %w", err)
	}

	return &identity_srv.RotateOAuth2ClientSecretResponse{
		ClientSecret: &plainSecret,
	}, nil
}

func (l *logicImpl) GetOAuth2ClientForAuth(
	ctx context.Context,
	clientID string,
) (*identity_srv.GetOAuth2ClientForAuthResponse, error) {
	client, err := l.dal.OAuth2Client().GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	id := client.ID.String()
	accessLifespan := int32(client.AccessTokenLifespan)
	refreshLifespan := int32(client.RefreshTokenLifespan)

	grantTypes := make([]string, len(client.GrantTypes))
	copy(grantTypes, client.GrantTypes)

	return &identity_srv.GetOAuth2ClientForAuthResponse{
		Id:                   &id,
		ClientID:             &client.ClientID,
		ClientSecretHash:     &client.ClientSecret,
		ClientName:           &client.ClientName,
		ClientType:           l.conv.OAuth2().ClientModelToThrift(client).ClientType,
		GrantTypes:           grantTypes,
		RedirectURIs:         []string(client.RedirectURIs),
		Scopes:               []string(client.Scopes),
		AccessTokenLifespan:  &accessLifespan,
		RefreshTokenLifespan: &refreshLifespan,
		IsActive:             &client.IsActive,
	}, nil
}

// --- Scope 管理 ---

func (l *logicImpl) ListOAuth2Scopes(
	ctx context.Context,
	req *identity_srv.ListOAuth2ScopesRequest,
) (*identity_srv.ListOAuth2ScopesResponse, error) {
	var scopes []*models.OAuth2Scope

	var err error

	if req.DefaultOnly != nil && *req.DefaultOnly {
		scopes, err = l.dal.OAuth2Scope().ListDefaults(ctx)
	} else {
		scopes, err = l.dal.OAuth2Scope().ListAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &identity_srv.ListOAuth2ScopesResponse{
		Scopes: l.conv.OAuth2().ScopeModelsToThrift(scopes),
	}, nil
}

// --- Consent 管理 ---

func (l *logicImpl) SaveOAuth2Consent(ctx context.Context, req *identity_srv.SaveOAuth2ConsentRequest) error {
	if req.UserID == nil || req.ClientID == nil {
		return fmt.Errorf("userID 和 clientID 不能为空")
	}

	userUUID, err := uuid.Parse(*req.UserID)
	if err != nil {
		return fmt.Errorf("无效的 userID: %w", err)
	}

	consent := &models.OAuth2Consent{
		UserID:   userUUID,
		ClientID: *req.ClientID,
		Scopes:   strings.Join(req.Scopes, " "),
	}

	return l.dal.OAuth2Consent().Save(ctx, consent)
}

func (l *logicImpl) GetOAuth2Consent(
	ctx context.Context,
	req *identity_srv.GetOAuth2ConsentRequest,
) (*identity_srv.GetOAuth2ConsentResponse, error) {
	if req.UserID == nil || req.ClientID == nil {
		return nil, fmt.Errorf("userID 和 clientID 不能为空")
	}

	consent, err := l.dal.OAuth2Consent().GetByUserAndClient(ctx, *req.UserID, *req.ClientID)
	if err != nil {
		found := false
		return &identity_srv.GetOAuth2ConsentResponse{Found: &found}, nil
	}

	found := true

	return &identity_srv.GetOAuth2ConsentResponse{
		Consent: l.conv.OAuth2().ConsentModelToThrift(consent),
		Found:   &found,
	}, nil
}

func (l *logicImpl) ListOAuth2Consents(
	ctx context.Context,
	req *identity_srv.ListOAuth2ConsentsRequest,
) (*identity_srv.ListOAuth2ConsentsResponse, error) {
	if req.UserID == nil {
		return nil, fmt.Errorf("userID 不能为空")
	}

	page, limit := 1, 20
	if req.Page != nil {
		page = int(req.Page.GetPage())
		limit = int(req.Page.GetLimit())
	}

	consents, _, err := l.dal.OAuth2Consent().ListByUserID(ctx, *req.UserID, page, limit)
	if err != nil {
		return nil, err
	}

	return &identity_srv.ListOAuth2ConsentsResponse{
		Consents: l.conv.OAuth2().ConsentModelsToThrift(consents),
	}, nil
}

func (l *logicImpl) RevokeOAuth2Consent(ctx context.Context, req *identity_srv.RevokeOAuth2ConsentRequest) error {
	if req.UserID == nil || req.ClientID == nil {
		return fmt.Errorf("userID 和 clientID 不能为空")
	}

	return l.dal.OAuth2Consent().Revoke(ctx, *req.UserID, *req.ClientID)
}
