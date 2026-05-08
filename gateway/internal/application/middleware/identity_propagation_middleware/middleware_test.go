package middleware

import (
	"context"
	"testing"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/iamclient"
)

func TestIdentityPropagation_HeadersMirroredToMetainfo(t *testing.T) {
	ctx := context.Background()
	c := app.NewContext(0)

	c.Request.Header.Set(iamclient.HeaderUserID, "u1")
	c.Request.Header.Set(iamclient.HeaderUserName, "alice")
	c.Request.Header.Set(iamclient.HeaderTenantID, "org-1")
	c.Request.Header.Set(iamclient.HeaderUserRoles, "doctor,admin")
	c.Request.Header.Set(iamclient.HeaderJTI, "jti-1")

	mw := NewIdentityPropagation().MiddlewareFunc()

	var captured context.Context

	c.SetHandlers([]app.HandlerFunc{
		mw,
		func(ctx context.Context, _ *app.RequestContext) { captured = ctx },
	})
	c.Next(ctx)

	require.NotNil(t, captured)

	got, ok := metainfo.GetPersistentValue(captured, iamclient.MetaUserID)
	assert.True(t, ok)
	assert.Equal(t, "u1", got)

	got, ok = metainfo.GetPersistentValue(captured, iamclient.MetaUserName)
	assert.True(t, ok)
	assert.Equal(t, "alice", got)

	got, ok = metainfo.GetPersistentValue(captured, iamclient.MetaTenantID)
	assert.True(t, ok)
	assert.Equal(t, "org-1", got)

	got, ok = metainfo.GetPersistentValue(captured, iamclient.MetaUserRoles)
	assert.True(t, ok)
	assert.Equal(t, "doctor,admin", got)

	got, ok = metainfo.GetPersistentValue(captured, iamclient.MetaJTI)
	assert.True(t, ok)
	assert.Equal(t, "jti-1", got)
}

func TestIdentityPropagation_NoHeadersNoMetainfo(t *testing.T) {
	ctx := context.Background()
	c := app.NewContext(0)

	mw := NewIdentityPropagation().MiddlewareFunc()

	var captured context.Context

	c.SetHandlers([]app.HandlerFunc{
		mw,
		func(ctx context.Context, _ *app.RequestContext) { captured = ctx },
	})
	c.Next(ctx)

	require.NotNil(t, captured)

	_, ok := metainfo.GetPersistentValue(captured, iamclient.MetaUserID)
	assert.False(t, ok, "no header should not pollute metainfo")
}

func TestIdentityPropagation_PartialHeaders(t *testing.T) {
	ctx := context.Background()
	c := app.NewContext(0)
	c.Request.Header.Set(iamclient.HeaderUserID, "u1")

	mw := NewIdentityPropagation().MiddlewareFunc()

	var captured context.Context

	c.SetHandlers([]app.HandlerFunc{
		mw,
		func(ctx context.Context, _ *app.RequestContext) { captured = ctx },
	})
	c.Next(ctx)

	got, ok := metainfo.GetPersistentValue(captured, iamclient.MetaUserID)
	assert.True(t, ok)
	assert.Equal(t, "u1", got)

	_, ok = metainfo.GetPersistentValue(captured, iamclient.MetaUserName)
	assert.False(t, ok)
}
