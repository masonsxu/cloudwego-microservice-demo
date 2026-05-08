package iamclient

import (
	"errors"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// decisionCache 是 PDP 决策的本地 LRU 缓存。
//
// 设计：
//   - 单 Decision 的存活时间通过 entry.expireAt 控制；
//   - 命中且未过期才返回；过期 / 未命中时调用方需重新走 RPC；
//   - LRU 容量驱逐由 hashicorp/golang-lru 实现，goroutine-safe。
//
// Phase 4 暂不实现"策略变更通知主动失效"——策略变更 → policy_srv 内部
// reload → 业务侧 TTL 内仍可能命中旧决策。后续若引入 etcd watch，
// 可在此处暴露 invalidate(jti) 与 invalidateAll() API。
type decisionCache struct {
	lru *lru.Cache[string, cacheEntry]
	ttl time.Duration
}

type cacheEntry struct {
	decision *Decision
	expireAt time.Time
}

// newDecisionCache 创建缓存。size <= 0 表示禁用缓存（返回 nil cache + nil err）。
func newDecisionCache(size int, ttl time.Duration) (*decisionCache, error) {
	if size <= 0 {
		return nil, nil
	}

	if ttl <= 0 {
		return nil, errors.New("iamclient: cache TTL must be positive")
	}

	c, err := lru.New[string, cacheEntry](size)
	if err != nil {
		return nil, err
	}

	return &decisionCache{
		lru: c,
		ttl: ttl,
	}, nil
}

// get 命中且未过期才返回 (decision, true)；其余返回 (nil, false)。
//
// 命中过期时主动从 LRU 中移除，避免堆积无效条目。
func (c *decisionCache) get(key string) (*Decision, bool) {
	if c == nil || c.lru == nil {
		return nil, false
	}

	entry, ok := c.lru.Get(key)
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.expireAt) {
		c.lru.Remove(key)
		return nil, false
	}

	return entry.decision, true
}

// set 写入或更新缓存条目。
func (c *decisionCache) set(key string, d *Decision) {
	if c == nil || c.lru == nil || d == nil {
		return
	}

	c.lru.Add(key, cacheEntry{
		decision: d,
		expireAt: time.Now().Add(c.ttl),
	})
}

// purge 清空全部缓存条目（保留供未来 invalidate 实现使用）。
func (c *decisionCache) purge() {
	if c == nil || c.lru == nil {
		return
	}

	c.lru.Purge()
}

// len 返回当前缓存条目数（包含未过期与已过期但未驱逐的）。供测试与监控使用。
func (c *decisionCache) len() int {
	if c == nil || c.lru == nil {
		return 0
	}

	return c.lru.Len()
}
