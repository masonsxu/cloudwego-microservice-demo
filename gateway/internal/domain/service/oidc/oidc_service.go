package oidc

import (
	"net/http"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/oidcstore"
)

// Service OIDC 领域服务接口
type Service interface {
	http.Handler
	Storage() *oidcstore.Storage
	Issuer() string
}
