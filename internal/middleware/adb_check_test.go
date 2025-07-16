package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/basiooo/andromodem/internal/middleware"
	adb "github.com/basiooo/goadb"
	"github.com/basiooo/goadb/wire"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAdbChecker_WithNilAdb(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	mw := middleware.AdbChecker(nil, logger)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := mw(next)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	assert.False(t, called)

	body := w.Body.String()
	assert.Contains(t, body, `"success":false`)
	assert.True(t, strings.Contains(body, "ADB server is currently not running") ||
		strings.Contains(body, "ADB is not installed"))
}

func TestAdbChecker_WithValidAdb(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()

	s := &adb.MockServer{
		Status:   wire.StatusSuccess,
		Messages: []string{"output"},
	}

	fakeAdb := &adb.Adb{Server: s}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AdbChecker(fakeAdb, logger)
	handler := mw(next)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.True(t, called)
}
