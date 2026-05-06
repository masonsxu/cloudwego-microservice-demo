package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/hertz-contrib/requestid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwtmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
)

func setupServer(t *testing.T, buf *bytes.Buffer, useRequestID bool) *server.Hertz {
	t.Helper()

	logger := zerolog.New(buf)
	h := server.Default()

	if useRequestID {
		h.Use(requestid.New())
	}

	mw := NewAccessLogMiddleware(&logger)
	h.Use(mw.MiddlewareFunc())

	h.GET("/foo", func(ctx context.Context, c *app.RequestContext) {
		c.Status(http.StatusCreated)
	})
	h.GET("/server-error", func(ctx context.Context, c *app.RequestContext) {
		c.Status(http.StatusInternalServerError)
	})

	return h
}

func parseLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	out := strings.TrimSpace(buf.String())
	require.NotEmpty(t, out, "expected access log output")

	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &entry))

	return entry
}

func TestAccessLog_LogsCoreFields(t *testing.T) {
	var buf bytes.Buffer
	h := setupServer(t, &buf, false)

	w := ut.PerformRequest(
		h.Engine,
		"GET",
		"/foo",
		nil,
		ut.Header{Key: jwtmw.HeaderUserID, Value: "user-1"},
	)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode())

	entry := parseLog(t, &buf)
	assert.Equal(t, "access", entry["message"])
	assert.Equal(t, "access_log", entry["component"])
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/foo", entry["path"])
	assert.EqualValues(t, http.StatusCreated, entry["status"])
	assert.Equal(t, "user-1", entry["user_id"])
	assert.Contains(t, entry, "duration")
}

func TestAccessLog_OmitsUserIDWhenAbsent(t *testing.T) {
	var buf bytes.Buffer
	h := setupServer(t, &buf, false)

	w := ut.PerformRequest(h.Engine, "GET", "/foo", nil)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode())

	entry := parseLog(t, &buf)
	_, hasUser := entry["user_id"]
	assert.False(t, hasUser, "user_id should be omitted when X-User-Id missing")
}

func TestAccessLog_PicksUpRequestID(t *testing.T) {
	var buf bytes.Buffer
	h := setupServer(t, &buf, true)

	w := ut.PerformRequest(h.Engine, "GET", "/foo", nil)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode())

	entry := parseLog(t, &buf)
	reqID, ok := entry["request_id"].(string)
	require.True(t, ok, "request_id should be present and a string")
	assert.NotEmpty(t, reqID)
}

func TestAccessLog_RecordsErrorStatus(t *testing.T) {
	var buf bytes.Buffer
	h := setupServer(t, &buf, false)

	w := ut.PerformRequest(h.Engine, "GET", "/server-error", nil)
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode())

	entry := parseLog(t, &buf)
	assert.EqualValues(t, http.StatusInternalServerError, entry["status"])
	assert.Equal(t, "/server-error", entry["path"])
}
