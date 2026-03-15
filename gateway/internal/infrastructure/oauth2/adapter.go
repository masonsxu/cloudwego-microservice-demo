package oauth2

import (
	"net/http"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// hertzToHTTPRequest 将 Hertz RequestContext 转换为标准 http.Request。
// fosite 需要标准 http.Request 来解析 OAuth2 请求参数。
func hertzToHTTPRequest(ctx *app.RequestContext) *http.Request {
	u, _ := url.Parse(string(ctx.Request.URI().FullURI()))

	r := &http.Request{
		Method: string(ctx.Method()),
		URL:    u,
		Header: make(http.Header),
		Form:   make(url.Values),
	}

	// 复制请求头
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		r.Header.Set(string(key), string(value))
	})

	// 解析表单数据（POST body）
	ctx.Request.PostArgs().VisitAll(func(key, value []byte) {
		r.Form.Set(string(key), string(value))
	})

	// 也复制 URL query 参数到 Form
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		r.Form.Set(string(key), string(value))
	})

	return r
}

// hertzResponseWriter 适配 Hertz RequestContext 为 http.ResponseWriter。
type hertzResponseWriter struct {
	ctx        *app.RequestContext
	statusCode int
	header     http.Header
	body       []byte
}

func newHertzResponseWriter(ctx *app.RequestContext) *hertzResponseWriter {
	return &hertzResponseWriter{
		ctx:        ctx,
		statusCode: consts.StatusOK,
		header:     make(http.Header),
	}
}

func (w *hertzResponseWriter) Header() http.Header {
	return w.header
}

func (w *hertzResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *hertzResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// flushResponse 将 hertzResponseWriter 中缓冲的响应写回 Hertz RequestContext。
func flushResponse(ctx *app.RequestContext, rw *hertzResponseWriter) {
	for key, values := range rw.header {
		for _, v := range values {
			ctx.Response.Header.Set(key, v)
		}
	}

	ctx.Response.SetStatusCode(rw.statusCode)

	if len(rw.body) > 0 {
		ctx.Response.SetBody(rw.body)
	}
}
