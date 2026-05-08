package iamclient

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefaults(t *testing.T) {
	cfg := applyDefaults(Config{
		EtcdEndpoints: []string{"etcd:2379"},
	})
	assert.Equal(t, defaultPolicyService, cfg.PolicyService)
	assert.Equal(t, defaultRPCTimeout, cfg.RPCTimeout)
	assert.Equal(t, defaultCacheSize, cfg.CacheSize)
	assert.Equal(t, defaultCacheTTL, cfg.CacheTTL)
}

func TestApplyDefaults_RespectsExplicitValues(t *testing.T) {
	custom := Config{
		EtcdEndpoints: []string{"etcd:2379"},
		PolicyService: "custom-pdp",
		RPCTimeout:    2 * time.Second,
		CacheSize:     50,
		CacheTTL:      time.Minute,
	}
	got := applyDefaults(custom)
	assert.Equal(t, custom, got)
}

func TestNew_RejectsEmptyEndpoints(t *testing.T) {
	_, err := New(Config{})
	assert.Error(t, err)
}

func TestPermissionDeniedError_Error(t *testing.T) {
	t.Run("with reason", func(t *testing.T) {
		e := &PermissionDeniedError{Action: "read", Resource: "p:1", Reason: "no role"}
		assert.Contains(t, e.Error(), "read")
		assert.Contains(t, e.Error(), "p:1")
		assert.Contains(t, e.Error(), "no role")
		assert.True(t, errors.Is(e, ErrPermissionDenied))
	})

	t.Run("without reason", func(t *testing.T) {
		e := &PermissionDeniedError{Action: "read", Resource: "p:1"}
		assert.Contains(t, e.Error(), "read")
		assert.Contains(t, e.Error(), "p:1")
	})
}

func TestClient_Close(t *testing.T) {
	c := &Client{}
	assert.NoError(t, c.Close())
}
