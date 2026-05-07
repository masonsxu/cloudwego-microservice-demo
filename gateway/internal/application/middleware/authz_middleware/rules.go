package middleware

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultPolicy 未匹配任何规则时的默认行为
type DefaultPolicy string

const (
	// DefaultAllow 已认证用户默认放行（未在 roles 中显式限制的路径）
	DefaultAllow DefaultPolicy = "allow"
	// DefaultDeny 默认拒绝，强制 YAML 显式覆盖所有路径（白名单模式）
	DefaultDeny DefaultPolicy = "deny"
)

// Rules 网关路由级 ACL 规则
type Rules struct {
	Default       DefaultPolicy
	Public        []Endpoint
	Authenticated []Endpoint
	Roles         []RolePrefix
}

// Endpoint method+path 端点声明
//
// Method 取值为大写 HTTP 方法（GET/POST/PUT/DELETE/PATCH 等），
// "*" 表示任意方法。
type Endpoint struct {
	Method string
	Path   string
}

// RolePrefix 路径前缀级角色门禁，Require 之间为 OR 关系
type RolePrefix struct {
	Prefix  string
	Require []string
}

type rawRules struct {
	Default       string          `yaml:"default"`
	Public        []string        `yaml:"public"`
	Authenticated []string        `yaml:"authenticated"`
	Roles         []rawRolePrefix `yaml:"roles"`
}

type rawRolePrefix struct {
	Prefix  string   `yaml:"prefix"`
	Require []string `yaml:"require"`
}

// LoadRulesFromFile 从 YAML 文件加载规则
func LoadRulesFromFile(path string) (*Rules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 authz_rules 文件失败: %w", err)
	}

	return ParseRules(data)
}

// ParseRules 从 YAML 字节流解析规则
func ParseRules(data []byte) (*Rules, error) {
	var raw rawRules
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("解析 authz_rules YAML 失败: %w", err)
	}

	rules := &Rules{Default: DefaultAllow}

	switch strings.ToLower(strings.TrimSpace(raw.Default)) {
	case "", string(DefaultAllow):
		rules.Default = DefaultAllow
	case string(DefaultDeny):
		rules.Default = DefaultDeny
	default:
		return nil, fmt.Errorf("非法的 default 值 %q（仅支持 allow / deny）", raw.Default)
	}

	for _, e := range raw.Public {
		ep, err := parseEndpoint(e)
		if err != nil {
			return nil, fmt.Errorf("public 规则 %q 解析失败: %w", e, err)
		}

		rules.Public = append(rules.Public, ep)
	}

	for _, e := range raw.Authenticated {
		ep, err := parseEndpoint(e)
		if err != nil {
			return nil, fmt.Errorf("authenticated 规则 %q 解析失败: %w", e, err)
		}

		rules.Authenticated = append(rules.Authenticated, ep)
	}

	for _, r := range raw.Roles {
		prefix := strings.TrimSpace(r.Prefix)
		if prefix == "" {
			return nil, fmt.Errorf("roles 规则缺少 prefix")
		}

		if !strings.HasPrefix(prefix, "/") {
			return nil, fmt.Errorf("roles prefix %q 必须以 / 开头", prefix)
		}

		if len(r.Require) == 0 {
			return nil, fmt.Errorf("roles 规则 prefix=%s 缺少 require", prefix)
		}

		rules.Roles = append(rules.Roles, RolePrefix{
			Prefix:  prefix,
			Require: r.Require,
		})
	}

	return rules, nil
}

// parseEndpoint 解析 "GET /api/v1/foo" 形式（method 空白 path），允许 method=*
func parseEndpoint(s string) (Endpoint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Endpoint{}, fmt.Errorf("端点声明为空")
	}

	parts := strings.Fields(s)
	if len(parts) != 2 {
		return Endpoint{}, fmt.Errorf("端点声明格式错误，需为 'METHOD /path'")
	}

	method := strings.ToUpper(parts[0])

	path := parts[1]
	if !strings.HasPrefix(path, "/") {
		return Endpoint{}, fmt.Errorf("path 必须以 / 开头")
	}

	return Endpoint{Method: method, Path: path}, nil
}

// MatchPublic method+path 是否命中 public 列表
func (r *Rules) MatchPublic(method, path string) bool {
	return matchEndpoints(r.Public, method, path)
}

// MatchAuthenticated 是否命中 authenticated 列表
func (r *Rules) MatchAuthenticated(method, path string) bool {
	return matchEndpoints(r.Authenticated, method, path)
}

// MatchRolePrefix 命中第一个匹配的角色规则；hit=false 表示无匹配
//
// 规则按声明顺序匹配，建议把更具体的前缀写在前面。
func (r *Rules) MatchRolePrefix(path string) (RolePrefix, bool) {
	for _, rp := range r.Roles {
		if strings.HasPrefix(path, rp.Prefix) {
			return rp, true
		}
	}

	return RolePrefix{}, false
}

// HasAnyRole 用户角色与 require 是否有交集
func (rp RolePrefix) HasAnyRole(userRoles []string) bool {
	if len(userRoles) == 0 {
		return false
	}

	needed := make(map[string]struct{}, len(rp.Require))
	for _, r := range rp.Require {
		needed[strings.TrimSpace(r)] = struct{}{}
	}

	for _, r := range userRoles {
		if _, ok := needed[strings.TrimSpace(r)]; ok {
			return true
		}
	}

	return false
}

func matchEndpoints(eps []Endpoint, method, path string) bool {
	for _, ep := range eps {
		if (ep.Method == "*" || ep.Method == method) && ep.Path == path {
			return true
		}
	}

	return false
}
