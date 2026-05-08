package iamclient

import (
	"context"
	"testing"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/stretchr/testify/assert"
)

func TestSubjectFromContext(t *testing.T) {
	c := &Client{}

	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		ctx = metainfo.WithPersistentValue(ctx, MetaUserID, "u1")
		ctx = metainfo.WithPersistentValue(ctx, MetaUserName, "alice")
		ctx = metainfo.WithPersistentValue(ctx, MetaTenantID, "org-1")
		ctx = metainfo.WithPersistentValue(ctx, MetaUserRoles, "doctor,admin")
		ctx = metainfo.WithPersistentValue(ctx, MetaRequestID, "req-xyz")
		ctx = metainfo.WithPersistentValue(ctx, MetaJTI, "jti-abc")

		subject, err := c.SubjectFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "u1", subject.UserID)
		assert.Equal(t, "alice", subject.UserName)
		assert.Equal(t, "org-1", subject.TenantID)
		assert.Equal(t, []string{"doctor", "admin"}, subject.Roles)
		assert.Equal(t, "req-xyz", subject.RequestID)
		assert.Equal(t, "jti-abc", subject.Jti)
		assert.Same(t, c, subject.client)
	})

	t.Run("missing user id", func(t *testing.T) {
		ctx := context.Background()
		ctx = metainfo.WithPersistentValue(ctx, MetaUserName, "alice")

		subject, err := c.SubjectFromContext(ctx)
		assert.Nil(t, subject)
		assert.ErrorIs(t, err, ErrMissingUserID)
	})

	t.Run("empty context", func(t *testing.T) {
		subject, err := c.SubjectFromContext(context.Background())
		assert.Nil(t, subject)
		assert.ErrorIs(t, err, ErrMissingUserID)
	})
}
