package iamclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecisionCache_HitMiss(t *testing.T) {
	c, err := newDecisionCache(10, time.Minute)
	require.NoError(t, err)

	got, ok := c.get("k1")
	assert.False(t, ok)
	assert.Nil(t, got)

	d := &Decision{Allowed: true, Reason: "ok"}
	c.set("k1", d)

	got, ok = c.get("k1")
	assert.True(t, ok)
	assert.Same(t, d, got)
}

func TestDecisionCache_TTLExpiry(t *testing.T) {
	c, err := newDecisionCache(10, 10*time.Millisecond)
	require.NoError(t, err)

	c.set("k1", &Decision{Allowed: true})

	got, ok := c.get("k1")
	require.True(t, ok)
	assert.True(t, got.Allowed)

	time.Sleep(15 * time.Millisecond)

	got, ok = c.get("k1")
	assert.False(t, ok)
	assert.Nil(t, got)
	assert.Equal(t, 0, c.len(), "expired entry should be removed on miss")
}

func TestDecisionCache_LRUEviction(t *testing.T) {
	c, err := newDecisionCache(2, time.Minute)
	require.NoError(t, err)

	c.set("a", &Decision{Allowed: true, Reason: "a"})
	c.set("b", &Decision{Allowed: true, Reason: "b"})
	c.set("c", &Decision{Allowed: true, Reason: "c"})

	// "a" should have been evicted (LRU)
	_, ok := c.get("a")
	assert.False(t, ok)

	got, ok := c.get("b")
	assert.True(t, ok)
	assert.Equal(t, "b", got.Reason)

	got, ok = c.get("c")
	assert.True(t, ok)
	assert.Equal(t, "c", got.Reason)
}

func TestDecisionCache_DisabledWhenSizeZero(t *testing.T) {
	c, err := newDecisionCache(0, time.Minute)
	assert.NoError(t, err)
	assert.Nil(t, c)

	// Disabled cache must tolerate get/set without panic.
	_, ok := c.get("anything")
	assert.False(t, ok)
	c.set("anything", &Decision{Allowed: true})
}

func TestDecisionCache_RejectsZeroTTL(t *testing.T) {
	_, err := newDecisionCache(10, 0)
	assert.Error(t, err)
}

func TestDecisionCache_Purge(t *testing.T) {
	c, err := newDecisionCache(10, time.Minute)
	require.NoError(t, err)

	c.set("k1", &Decision{Allowed: true})
	c.set("k2", &Decision{Allowed: true})
	require.Equal(t, 2, c.len())

	c.purge()
	assert.Equal(t, 0, c.len())
}

func TestSubjectCacheKey_Stable(t *testing.T) {
	s := &Subject{
		UserID: "u1",
		Jti:    "jti-1",
	}

	k1 := s.cacheKey("read", "patient:7", map[string]string{"a": "1", "b": "2"})
	k2 := s.cacheKey("read", "patient:7", map[string]string{"b": "2", "a": "1"})
	assert.Equal(t, k1, k2, "attribute map order should not affect key")
}

func TestSubjectCacheKey_RolesOrderStable_WhenJtiMissing(t *testing.T) {
	s1 := &Subject{
		UserID: "u1",
		Roles:  []string{"doctor", "admin"},
	}
	s2 := &Subject{
		UserID: "u1",
		Roles:  []string{"admin", "doctor"},
	}

	k1 := s1.cacheKey("read", "patient", nil)
	k2 := s2.cacheKey("read", "patient", nil)
	assert.Equal(t, k1, k2, "roles slice order should not affect key when jti missing")
}

func TestSubjectCacheKey_DifferentResourcesProduceDifferentKeys(t *testing.T) {
	s := &Subject{UserID: "u1", Jti: "jti-1"}

	a := s.cacheKey("read", "patient:1", nil)
	b := s.cacheKey("read", "patient:2", nil)
	assert.NotEqual(t, a, b)
}
