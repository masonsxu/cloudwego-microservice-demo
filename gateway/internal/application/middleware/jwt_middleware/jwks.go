package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
)

// jwksKey JWKS 格式的单键描述
type jwksKey struct {
	KTY string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	KID string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// jwksResponse JWKS 响应体
type jwksResponse struct {
	Keys []jwksKey `json:"keys"`
}

// JWKSHandler 返回 JWKS 端点 handler（从公钥文件读取并组装 JWKS）
func JWKSHandler(pubKeyPath string) app.HandlerFunc {
	jwksBytes, err := buildJWKSJSON(pubKeyPath)
	if err != nil {
		panic("failed to build JWKS from public key: " + err.Error())
	}

	return func(_ context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, json.RawMessage(jwksBytes))
	}
}

// buildJWKSJSON 从 RSA 公钥 PEM 文件构建 JWKS JSON 字节
func buildJWKSJSON(pubKeyPath string) ([]byte, error) {
	data, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取公钥文件: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无效的 PEM 格式")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥不是 RSA 类型")
	}

	// kid: 公钥 modulus 的 SHA-256 前 8 字节的 base64url
	kidHash := sha256.Sum256(rsaPub.N.Bytes())
	kid := base64.RawURLEncoding.EncodeToString(kidHash[:8])

	eBytes := bigEndianBytes(rsaPub.E)

	jwks := jwksResponse{
		Keys: []jwksKey{{
			KTY: "RSA",
			Use: "sig",
			Alg: "RS256",
			KID: kid,
			N:   base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(eBytes),
		}},
	}

	return json.Marshal(jwks)
}

// bigEndianBytes 将 int 编码为大端字节（去除前导零）
func bigEndianBytes(v int) []byte {
	if v == 0 {
		return []byte{0}
	}

	var buf []byte
	for v > 0 {
		buf = append([]byte{byte(v & 0xff)}, buf...)
		v >>= 8
	}

	// 去除前导零（JWK 规范要求）
	for len(buf) > 1 && buf[0] == 0 {
		buf = buf[1:]
	}

	return buf
}
