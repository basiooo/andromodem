package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/basiooo/andromodem/internal/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRecoverer_PanicHandled(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Meledak boom.!!!!")
	})

	handler := middleware.Recoverer(logger)(panicHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Contains(t, w.Body.String(), `"success":false`)
	assert.Contains(t, w.Body.String(), `"Internal Server Error"`)
}

func TestRecoverer_NoPanic(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()

	called := false
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"success":true}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	handler := middleware.Recoverer(logger)(okHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, called)
	assert.True(t, strings.Contains(w.Body.String(), `"success":true`))
}
