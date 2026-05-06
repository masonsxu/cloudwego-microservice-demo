package middleware

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
)

func ptrString(s string) *string { return &s }

func TestInjectIdentityHeaders_AllFields(t *testing.T) {
	ctx := &app.RequestContext{}
	ctx.Request.Header = protocol.RequestHeader{}

	claims := &http_base.JWTClaimsDTO{
		UserProfileID:  ptrString("user-123"),
		Username:       ptrString("alice"),
		OrganizationID: ptrString("org-456"),
		RoleIDs:        []string{"admin", "doctor"},
	}

	injectIdentityHeaders(ctx, claims)

	assert.Equal(t, "user-123", string(ctx.Request.Header.Get(HeaderUserID)))
	assert.Equal(t, "alice", string(ctx.Request.Header.Get(HeaderUserName)))
	assert.Equal(t, "org-456", string(ctx.Request.Header.Get(HeaderTenantID)))
	assert.Equal(t, "admin,doctor", string(ctx.Request.Header.Get(HeaderUserRoles)))
}

func TestInjectIdentityHeaders_NilClaims(t *testing.T) {
	ctx := &app.RequestContext{}
	ctx.Request.Header = protocol.RequestHeader{}

	injectIdentityHeaders(ctx, nil)

	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserID)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserName)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderTenantID)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserRoles)))
}

func TestInjectIdentityHeaders_PartialFields(t *testing.T) {
	ctx := &app.RequestContext{}
	ctx.Request.Header = protocol.RequestHeader{}

	claims := &http_base.JWTClaimsDTO{
		UserProfileID: ptrString("user-only"),
		// Username / OrganizationID / RoleIDs 缺省
	}

	injectIdentityHeaders(ctx, claims)

	assert.Equal(t, "user-only", string(ctx.Request.Header.Get(HeaderUserID)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserName)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderTenantID)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserRoles)))
}

func TestInjectIdentityHeaders_EmptyStringNotInjected(t *testing.T) {
	ctx := &app.RequestContext{}
	ctx.Request.Header = protocol.RequestHeader{}

	empty := ""
	claims := &http_base.JWTClaimsDTO{
		UserProfileID:  &empty,
		Username:       &empty,
		OrganizationID: &empty,
		RoleIDs:        []string{},
	}

	injectIdentityHeaders(ctx, claims)

	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserID)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserName)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderTenantID)))
	assert.Empty(t, string(ctx.Request.Header.Get(HeaderUserRoles)))
}

func TestInjectIdentityHeaders_SingleRole(t *testing.T) {
	ctx := &app.RequestContext{}
	ctx.Request.Header = protocol.RequestHeader{}

	claims := &http_base.JWTClaimsDTO{
		RoleIDs: []string{"admin"},
	}

	injectIdentityHeaders(ctx, claims)

	assert.Equal(t, "admin", string(ctx.Request.Header.Get(HeaderUserRoles)))
}
