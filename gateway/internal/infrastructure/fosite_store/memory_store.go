// Package fositestore 实现了 fosite 所需的各种存储接口。
// 当前使用内存存储（开发阶段），生产环境应通过 RPC 调用 identity_srv 持久化。
package fositestore

import (
	"context"
	"sync"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/pkce"
)

// MemoryStore 内存存储实现，用于开发和测试。
// 实现了 fosite 需要的所有存储接口：
//   - fosite.ClientManager
//   - oauth2.AuthorizeCodeStorage
//   - oauth2.AccessTokenStorage
//   - oauth2.RefreshTokenStorage
//   - oauth2.TokenRevocationStorage
//   - pkce.PKCERequestStorage
type MemoryStore struct {
	mu sync.RWMutex

	clients        map[string]fosite.Client
	authCodes      map[string]storedRequest
	accessTokens   map[string]storedRequest
	refreshTokens  map[string]storedRequest
	pkceSessions   map[string]storedRequest
	accessTokenIdx map[string]string // requestID -> signature
	refreshIdx     map[string]string // requestID -> signature
}

type storedRequest struct {
	request   fosite.Requester
	createdAt time.Time
}

// NewMemoryStore 创建内存存储实例。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		clients:        make(map[string]fosite.Client),
		authCodes:      make(map[string]storedRequest),
		accessTokens:   make(map[string]storedRequest),
		refreshTokens:  make(map[string]storedRequest),
		pkceSessions:   make(map[string]storedRequest),
		accessTokenIdx: make(map[string]string),
		refreshIdx:     make(map[string]string),
	}
}

// RegisterClient 注册一个 OAuth2 客户端到存储中。
func (s *MemoryStore) RegisterClient(client fosite.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client.GetID()] = client
}

// -- fosite.ClientManager --

func (s *MemoryStore) GetClient(_ context.Context, id string) (fosite.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, ok := s.clients[id]
	if !ok {
		return nil, fosite.ErrNotFound
	}

	return client, nil
}

func (s *MemoryStore) ClientAssertionJWTValid(_ context.Context, _ string) error {
	return nil
}

func (s *MemoryStore) SetClientAssertionJWT(_ context.Context, _ string, _ time.Time) error {
	return nil
}

// -- oauth2.AuthorizeCodeStorage --

func (s *MemoryStore) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.authCodes[code] = storedRequest{request: req, createdAt: time.Now()}

	return nil
}

func (s *MemoryStore) GetAuthorizeCodeSession(
	_ context.Context,
	code string,
	session fosite.Session,
) (fosite.Requester, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored, ok := s.authCodes[code]
	if !ok {
		return nil, fosite.ErrNotFound
	}

	return stored.request, nil
}

func (s *MemoryStore) InvalidateAuthorizeCodeSession(_ context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.authCodes, code)

	return nil
}

// -- oauth2.AccessTokenStorage --

func (s *MemoryStore) CreateAccessTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.accessTokens[signature] = storedRequest{request: req, createdAt: time.Now()}
	s.accessTokenIdx[req.GetID()] = signature

	return nil
}

func (s *MemoryStore) GetAccessTokenSession(
	_ context.Context,
	signature string,
	session fosite.Session,
) (fosite.Requester, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored, ok := s.accessTokens[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}

	return stored.request, nil
}

func (s *MemoryStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.accessTokens, signature)

	return nil
}

// -- oauth2.RefreshTokenStorage --

func (s *MemoryStore) CreateRefreshTokenSession(
	_ context.Context,
	signature string,
	_ string,
	req fosite.Requester,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.refreshTokens[signature] = storedRequest{request: req, createdAt: time.Now()}
	s.refreshIdx[req.GetID()] = signature

	return nil
}

func (s *MemoryStore) GetRefreshTokenSession(
	_ context.Context,
	signature string,
	session fosite.Session,
) (fosite.Requester, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored, ok := s.refreshTokens[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}

	return stored.request, nil
}

func (s *MemoryStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.refreshTokens, signature)

	return nil
}

// -- oauth2.TokenRevocationStorage --

func (s *MemoryStore) RevokeRefreshToken(_ context.Context, requestID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sig, ok := s.refreshIdx[requestID]; ok {
		delete(s.refreshTokens, sig)
		delete(s.refreshIdx, requestID)
	}

	return nil
}

func (s *MemoryStore) RevokeRefreshTokenMaybeGracePeriod(_ context.Context, requestID string, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sig, ok := s.refreshIdx[requestID]; ok {
		delete(s.refreshTokens, sig)
		delete(s.refreshIdx, requestID)
	}

	return nil
}

func (s *MemoryStore) RevokeAccessToken(_ context.Context, requestID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sig, ok := s.accessTokenIdx[requestID]; ok {
		delete(s.accessTokens, sig)
		delete(s.accessTokenIdx, requestID)
	}

	return nil
}

func (s *MemoryStore) RotateRefreshToken(_ context.Context, requestID string, refreshTokenSignature string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 删除旧的 refresh token
	if oldSig, ok := s.refreshIdx[requestID]; ok {
		delete(s.refreshTokens, oldSig)
	}

	// 更新索引指向新签名
	s.refreshIdx[requestID] = refreshTokenSignature

	return nil
}

// -- pkce.PKCERequestStorage --

func (s *MemoryStore) CreatePKCERequestSession(_ context.Context, signature string, req fosite.Requester) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pkceSessions[signature] = storedRequest{request: req, createdAt: time.Now()}

	return nil
}

func (s *MemoryStore) GetPKCERequestSession(
	_ context.Context,
	signature string,
	session fosite.Session,
) (fosite.Requester, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored, ok := s.pkceSessions[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}

	return stored.request, nil
}

func (s *MemoryStore) DeletePKCERequestSession(_ context.Context, signature string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pkceSessions, signature)

	return nil
}

// 编译时接口验证
var (
	_ fosite.ClientManager          = (*MemoryStore)(nil)
	_ oauth2.AuthorizeCodeStorage   = (*MemoryStore)(nil)
	_ oauth2.AccessTokenStorage     = (*MemoryStore)(nil)
	_ oauth2.RefreshTokenStorage    = (*MemoryStore)(nil)
	_ oauth2.TokenRevocationStorage = (*MemoryStore)(nil)
	_ pkce.PKCERequestStorage       = (*MemoryStore)(nil)
)
