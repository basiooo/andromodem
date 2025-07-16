package common_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFromRequestBody(t *testing.T) {
	t.Parallel()
	// Test case 1: Successful decoding
	t.Run("Successful decoding", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		input := TestStruct{Name: "John Doe", Age: 30}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		var result TestStruct
		err := common.ReadFromRequestBody(req, &result)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, result.Name, "Expected name to match")
		assert.Equal(t, input.Age, result.Age, "Expected age to match")
	})

	// Test case 2: Invalid JSON
	t.Run("Invalid JSON should panic", func(t *testing.T) {
		invalidBody := []byte(`{"name": "John", "age": "thirty"}`) // Age is string, should be int
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(invalidBody))

		var result struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		err := common.ReadFromRequestBody(req, &result)
		assert.Error(t, err)
	})

	// Test case 3: Empty body
	t.Run("Empty body should panic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte{}))

		var result struct {
			Name string `json:"name"`
		}

		err := common.ReadFromRequestBody(req, &result)
		assert.Error(t, err)
	})
}

func TestWriteToResponseBody(t *testing.T) {
	t.Parallel()
	// Test case 1: Successful encoding and status code
	t.Run("Successful encoding and status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		response := map[string]string{"message": "Success"}
		statusCode := http.StatusOK

		common.WriteToResponseBody(recorder, response, statusCode)

		assert.Equal(t, statusCode, recorder.Code, "Expected status code to match")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		var result map[string]string
		err := json.Unmarshal(recorder.Body.Bytes(), &result)
		require.NoError(t, err, "Failed to unmarshal response body")
		assert.Equal(t, "Success", result["message"], "Expected message to be 'Success'")
	})

	// Test case 2: Different status code
	t.Run("Different status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		response := map[string]string{"error": "Not Found"}
		statusCode := http.StatusNotFound

		common.WriteToResponseBody(recorder, response, statusCode)

		assert.Equal(t, statusCode, recorder.Code, "Expected status code to match")
	})

	// Test case 3: Nil response body
	t.Run("Nil response body", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		statusCode := http.StatusNoContent // 204 No Content

		common.WriteToResponseBody(recorder, nil, statusCode)

		assert.Equal(t, statusCode, recorder.Code, "Expected status code to match")
		assert.Empty(t, recorder.Body.Bytes(), "Expected empty body for nil response")
	})

	// Test case 4: Unencodable response, should panic
	t.Run("Unencodable response should panic", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		// Channel types cannot be marshaled to JSON, causing an error
		unencodableResponse := make(chan int)
		statusCode := http.StatusOK

		assert.Panics(t, func() {
			common.WriteToResponseBody(recorder, unencodableResponse, statusCode)
		}, "The code should panic on unencodable response")
	})
}

// MockResponseWriter is a mock implementation of http.ResponseWriter for testing panic cases
// This is still useful if you need to simulate specific Write/WriteHeader behaviors that cause panics
// beyond what httptest.NewRecorder can easily do, though for these specific tests, testify's Panics
// assertion handles the direct panic from the function under test.
type MockResponseWriter struct {
	HeaderMap http.Header
	Body      *bytes.Buffer
	Status    int
	Panics    bool
}

func NewMockResponseWriter(panics bool) *MockResponseWriter {
	return &MockResponseWriter{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Panics:    panics,
	}
}

func (m *MockResponseWriter) Header() http.Header {
	return m.HeaderMap
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	if m.Panics {
		panic("mock write panic")
	}
	return m.Body.Write(b)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Status = statusCode
}

// Ensure that MockResponseWriter implements http.ResponseWriter
var _ http.ResponseWriter = (*MockResponseWriter)(nil)
