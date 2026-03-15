package fositestore

import "github.com/ory/fosite"

// DefaultClient 是一个简单的 fosite.Client 实现，用于内存存储。
type DefaultClient struct {
	ID            string
	Secret        []byte
	RedirectURIs  []string
	GrantTypes    []string
	ResponseTypes []string
	Scopes        []string
	Public        bool
}

func (c *DefaultClient) GetID() string                      { return c.ID }
func (c *DefaultClient) GetHashedSecret() []byte            { return c.Secret }
func (c *DefaultClient) GetRedirectURIs() []string          { return c.RedirectURIs }
func (c *DefaultClient) GetGrantTypes() fosite.Arguments    { return c.GrantTypes }
func (c *DefaultClient) GetResponseTypes() fosite.Arguments { return c.ResponseTypes }
func (c *DefaultClient) GetScopes() fosite.Arguments        { return c.Scopes }
func (c *DefaultClient) IsPublic() bool                     { return c.Public }
func (c *DefaultClient) GetAudience() fosite.Arguments      { return nil }

var _ fosite.Client = (*DefaultClient)(nil)
