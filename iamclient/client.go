package iamclient

import (
	"errors"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv/policyservice"
)

// Config 是创建 Client 所需的配置。
type Config struct {
	// EtcdEndpoints policy_srv 服务发现地址（必填）。
	EtcdEndpoints []string

	// PolicyService policy_srv 注册名（默认 "policy-service"）。
	PolicyService string

	// CallerService 当前业务服务名（用于 RPC 调用方标识）。
	CallerService string

	// RPCTimeout 调 policy_srv 的超时（默认 1s）。
	RPCTimeout time.Duration

	// CacheSize LRU 缓存条目数（默认 10000，传 -1 关闭缓存）。
	CacheSize int

	// CacheTTL 单条缓存有效期（默认 30s）。
	CacheTTL time.Duration
}

const (
	defaultPolicyService = "policy-service"
	defaultRPCTimeout    = time.Second
	defaultCacheSize     = 10000
	defaultCacheTTL      = 30 * time.Second
)

// Client 是业务系统接入 IAM 的唯一入口。
//
// 持有 policy_srv kitex 客户端与本地决策缓存，goroutine-safe，应用启动时
// 创建一次并复用。
type Client struct {
	cfg    Config
	policy policyservice.Client
	cache  *decisionCache
}

// New 创建并初始化 IAM 客户端。
//
// 调用方负责在应用关闭时调用 Client.Close 释放资源。
func New(cfg Config) (*Client, error) {
	if len(cfg.EtcdEndpoints) == 0 {
		return nil, errors.New("iamclient: EtcdEndpoints is required")
	}

	cfg = applyDefaults(cfg)

	resolver, err := etcd.NewEtcdResolver(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("iamclient: create etcd resolver: %w", err)
	}

	opts := []client.Option{
		client.WithResolver(resolver),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithRPCTimeout(cfg.RPCTimeout),
	}

	if cfg.CallerService != "" {
		opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: cfg.CallerService,
		}))
	}

	policyCli, err := policyservice.NewClient(cfg.PolicyService, opts...)
	if err != nil {
		return nil, fmt.Errorf("iamclient: create policy client: %w", err)
	}

	cache, err := newDecisionCache(cfg.CacheSize, cfg.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("iamclient: init cache: %w", err)
	}

	return &Client{
		cfg:    cfg,
		policy: policyCli,
		cache:  cache,
	}, nil
}

// Close 释放底层资源（当前实现为空，保留语义）。
func (c *Client) Close() error {
	return nil
}

func applyDefaults(cfg Config) Config {
	if cfg.PolicyService == "" {
		cfg.PolicyService = defaultPolicyService
	}

	if cfg.RPCTimeout <= 0 {
		cfg.RPCTimeout = defaultRPCTimeout
	}

	if cfg.CacheSize == 0 {
		cfg.CacheSize = defaultCacheSize
	}

	if cfg.CacheTTL <= 0 {
		cfg.CacheTTL = defaultCacheTTL
	}

	return cfg
}
