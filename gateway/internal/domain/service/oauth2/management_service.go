package oauth2

import (
	"context"
	"time"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	httpBase "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	oauth2model "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/oauth2"
	oauth2asm "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/oauth2"
	domaincommon "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
	identity_srv "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// OAuth2ManagementService OAuth2 管理服务接口
type OAuth2ManagementService interface {
	CreateClient(
		ctx context.Context,
		req *oauth2model.CreateOAuth2ClientRequestDTO,
	) (*oauth2model.CreateOAuth2ClientResponseDTO, error)
	GetClient(ctx context.Context, clientID string) (*oauth2model.OAuth2ClientResponseDTO, error)
	UpdateClient(
		ctx context.Context,
		clientID string,
		req *oauth2model.UpdateOAuth2ClientRequestDTO,
	) (*oauth2model.OAuth2ClientResponseDTO, error)
	DeleteClient(ctx context.Context, clientID string) (*httpBase.OperationStatusResponseDTO, error)
	ListClients(
		ctx context.Context,
		req *oauth2model.ListOAuth2ClientsRequestDTO,
	) (*oauth2model.ListOAuth2ClientsResponseDTO, error)
	RotateClientSecret(ctx context.Context, clientID string) (*oauth2model.RotateOAuth2ClientSecretResponseDTO, error)
	ListScopes(ctx context.Context) (*oauth2model.ListOAuth2ScopesResponseDTO, error)
	GetConfig(ctx context.Context) (*oauth2model.GetOAuth2ConfigResponseDTO, error)
	ListMyConsents(
		ctx context.Context,
		userID string,
		req *oauth2model.ListMyOAuth2ConsentsRequestDTO,
	) (*oauth2model.ListOAuth2ConsentsResponseDTO, error)
	RevokeMyConsent(ctx context.Context, userID, clientID string) (*httpBase.OperationStatusResponseDTO, error)
}

type oauth2ManagementServiceImpl struct {
	*domaincommon.BaseService
	identityClient identitycli.IdentityClient
	assembler      oauth2asm.Assembler
}

// NewOAuth2ManagementService 创建 OAuth2 管理服务实例
func NewOAuth2ManagementService(
	identityClient identitycli.IdentityClient,
	assembler oauth2asm.Assembler,
	logger *hertzZerolog.Logger,
) OAuth2ManagementService {
	return &oauth2ManagementServiceImpl{
		BaseService:    domaincommon.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

func (s *oauth2ManagementServiceImpl) CreateClient(
	ctx context.Context,
	req *oauth2model.CreateOAuth2ClientRequestDTO,
) (*oauth2model.CreateOAuth2ClientResponseDTO, error) {
	rpcReq := s.assembler.ToRPCCreateClientRequest(req)

	result, err := s.ProcessRPCCall(ctx, "创建 OAuth2 客户端",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.CreateOAuth2Client(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.CreateOAuth2ClientResponse)

	return &oauth2model.CreateOAuth2ClientResponseDTO{
		BaseResp:     s.ResponseBuilder().BuildSuccessResponse(),
		Client:       s.assembler.ToHTTPClient(rpcResp.Client),
		ClientSecret: rpcResp.ClientSecret,
	}, nil
}

func (s *oauth2ManagementServiceImpl) GetClient(
	ctx context.Context,
	clientID string,
) (*oauth2model.OAuth2ClientResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取 OAuth2 客户端",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.GetOAuth2Client(ctx, &identity_srv.GetOAuth2ClientRequest{Id: &clientID})
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.OAuth2Client)

	return &oauth2model.OAuth2ClientResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Client:   s.assembler.ToHTTPClient(rpcResp),
	}, nil
}

func (s *oauth2ManagementServiceImpl) UpdateClient(
	ctx context.Context,
	clientID string,
	req *oauth2model.UpdateOAuth2ClientRequestDTO,
) (*oauth2model.OAuth2ClientResponseDTO, error) {
	rpcReq := &identity_srv.UpdateOAuth2ClientRequest{
		Id:           &clientID,
		ClientName:   req.ClientName,
		Description:  req.Description,
		RedirectURIs: req.RedirectURIs,
		Scopes:       req.Scopes,
		LogoURI:      req.LogoURI,
		ClientURI:    req.ClientURI,
		IsActive:     req.IsActive,
	}

	if req.AccessTokenLifespan != nil {
		rpcReq.AccessTokenLifespan = req.AccessTokenLifespan
	}

	if req.RefreshTokenLifespan != nil {
		rpcReq.RefreshTokenLifespan = req.RefreshTokenLifespan
	}

	result, err := s.ProcessRPCCall(ctx, "更新 OAuth2 客户端",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.UpdateOAuth2Client(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.OAuth2Client)

	return &oauth2model.OAuth2ClientResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Client:   s.assembler.ToHTTPClient(rpcResp),
	}, nil
}

func (s *oauth2ManagementServiceImpl) DeleteClient(
	ctx context.Context,
	clientID string,
) (*httpBase.OperationStatusResponseDTO, error) {
	_, err := s.ProcessRPCCall(ctx, "删除 OAuth2 客户端",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.DeleteOAuth2Client(ctx, &identity_srv.DeleteOAuth2ClientRequest{Id: &clientID})
		},
	)
	if err != nil {
		return nil, err
	}

	return &httpBase.OperationStatusResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
	}, nil
}

func (s *oauth2ManagementServiceImpl) ListClients(
	ctx context.Context,
	req *oauth2model.ListOAuth2ClientsRequestDTO,
) (*oauth2model.ListOAuth2ClientsResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "列出 OAuth2 客户端",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.ListOAuth2Clients(ctx, &identity_srv.ListOAuth2ClientsRequest{
				IsActive: req.IsActive,
			})
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.ListOAuth2ClientsResponse)

	return &oauth2model.ListOAuth2ClientsResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Clients:  s.assembler.ToHTTPClients(rpcResp.Clients),
	}, nil
}

func (s *oauth2ManagementServiceImpl) RotateClientSecret(
	ctx context.Context,
	clientID string,
) (*oauth2model.RotateOAuth2ClientSecretResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "轮换 OAuth2 客户端密钥",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.RotateOAuth2ClientSecret(
				ctx,
				&identity_srv.RotateOAuth2ClientSecretRequest{Id: &clientID},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.RotateOAuth2ClientSecretResponse)

	return &oauth2model.RotateOAuth2ClientSecretResponseDTO{
		BaseResp:     s.ResponseBuilder().BuildSuccessResponse(),
		ClientSecret: rpcResp.ClientSecret,
	}, nil
}

func (s *oauth2ManagementServiceImpl) ListScopes(
	ctx context.Context,
) (*oauth2model.ListOAuth2ScopesResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "列出 OAuth2 作用域",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.ListOAuth2Scopes(ctx, &identity_srv.ListOAuth2ScopesRequest{})
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.ListOAuth2ScopesResponse)

	return &oauth2model.ListOAuth2ScopesResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Scopes:   s.assembler.ToHTTPScopes(rpcResp.Scopes),
	}, nil
}

func (s *oauth2ManagementServiceImpl) GetConfig(
	_ context.Context,
) (*oauth2model.GetOAuth2ConfigResponseDTO, error) {
	cfg := getEffectiveOAuth2Config()

	enabled := cfg.Enabled
	issuer := cfg.Issuer
	accessTokenLifespan := int64(cfg.AccessTokenLifespan / time.Second)
	refreshTokenLifespan := int64(cfg.RefreshTokenLifespan / time.Second)
	authCodeLifespan := int64(cfg.AuthCodeLifespan / time.Second)
	enforcePKCE := cfg.EnforcePKCE
	consentPageURL := cfg.ConsentPageURL

	return &oauth2model.GetOAuth2ConfigResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Config: &oauth2model.OAuth2ConfigDTO{
			Enabled:              &enabled,
			Issuer:               &issuer,
			AccessTokenLifespan:  &accessTokenLifespan,
			RefreshTokenLifespan: &refreshTokenLifespan,
			AuthCodeLifespan:     &authCodeLifespan,
			EnforcePKCE:          &enforcePKCE,
			ConsentPageURL:       &consentPageURL,
		},
	}, nil
}

func getEffectiveOAuth2Config() config.OAuth2Config {
	if config.Config != nil {
		return config.Config.OAuth2
	}

	return config.OAuth2Config{
		Enabled:              true,
		Issuer:               "http://localhost:8080",
		AccessTokenLifespan:  time.Hour,
		RefreshTokenLifespan: 30 * 24 * time.Hour,
		AuthCodeLifespan:     10 * time.Minute,
		EnforcePKCE:          true,
		ConsentPageURL:       "/oauth2/consent",
	}
}

func (s *oauth2ManagementServiceImpl) ListMyConsents(
	ctx context.Context,
	userID string,
	_ *oauth2model.ListMyOAuth2ConsentsRequestDTO,
) (*oauth2model.ListOAuth2ConsentsResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "列出用户授权同意",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.ListOAuth2Consents(ctx, &identity_srv.ListOAuth2ConsentsRequest{UserID: &userID})
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.ListOAuth2ConsentsResponse)

	return &oauth2model.ListOAuth2ConsentsResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Consents: s.assembler.ToHTTPConsents(rpcResp.Consents),
	}, nil
}

func (s *oauth2ManagementServiceImpl) RevokeMyConsent(
	ctx context.Context,
	userID, clientID string,
) (*httpBase.OperationStatusResponseDTO, error) {
	_, rpcErr := s.identityClient.RevokeOAuth2Consent(ctx, &identity_srv.RevokeOAuth2ConsentRequest{
		UserID:   &userID,
		ClientID: &clientID,
	})
	if err := errors.ProcessRPCError(rpcErr, "撤销授权同意失败"); err != nil {
		return nil, err
	}

	return &httpBase.OperationStatusResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
	}, nil
}
