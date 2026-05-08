package iamclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitRoles(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want []string
	}{
		{"empty", "", nil},
		{"single", "doctor", []string{"doctor"}},
		{"multiple", "doctor,admin", []string{"doctor", "admin"}},
		{"trims spaces", " doctor , admin ", []string{"doctor", "admin"}},
		{"skips empty", "doctor,,admin", []string{"doctor", "admin"}},
		{"only commas", ",,,", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, splitRoles(tc.raw))
		})
	}
}

func TestSubjectFromHeader(t *testing.T) {
	c := &Client{}

	t.Run("happy path", func(t *testing.T) {
		h := http.Header{}
		h.Set(HeaderUserID, "u1")
		h.Set(HeaderUserName, "alice")
		h.Set(HeaderTenantID, "org-1")
		h.Set(HeaderUserRoles, "doctor,admin")
		h.Set(HeaderRequestID, "req-xyz")
		h.Set(HeaderJTI, "jti-abc")

		subject, err := c.SubjectFromHeader(h)
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
		h := http.Header{}
		h.Set(HeaderUserName, "alice")

		subject, err := c.SubjectFromHeader(h)
		assert.Nil(t, subject)
		assert.ErrorIs(t, err, ErrMissingUserID)
	})

	t.Run("nil header", func(t *testing.T) {
		subject, err := c.SubjectFromHeader(nil)
		assert.Nil(t, subject)
		assert.ErrorIs(t, err, ErrMissingUserID)
	})

	t.Run("only user id", func(t *testing.T) {
		h := http.Header{}
		h.Set(HeaderUserID, "u1")

		subject, err := c.SubjectFromHeader(h)
		assert.NoError(t, err)
		assert.Equal(t, "u1", subject.UserID)
		assert.Empty(t, subject.UserName)
		assert.Nil(t, subject.Roles)
	})
}

// fakeHeader 实现 HeaderGetter，确保接口对非 net/http.Header 类型也工作。
type fakeHeader map[string]string

func (f fakeHeader) Get(name string) string { return f[name] }

func TestSubjectFromHeader_AnyHeaderGetter(t *testing.T) {
	c := &Client{}
	h := fakeHeader{
		HeaderUserID:    "u2",
		HeaderUserRoles: "admin",
	}

	subject, err := c.SubjectFromHeader(h)
	assert.NoError(t, err)
	assert.Equal(t, "u2", subject.UserID)
	assert.Equal(t, []string{"admin"}, subject.Roles)
}
