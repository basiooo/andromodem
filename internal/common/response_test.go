package common_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDeviceNotFoundResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	common.DeviceNotFoundResponse(recorder)

	assert.Equal(t, http.StatusNotFound, recorder.Code, "Expected status code 404 Not Found")

	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

	var actualResponse model.BaseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err, "Failed to unmarshal response body")

	expectedResponse := model.BaseResponse{
		Success: false,
		Message: "Device not found",
	}
	assert.Equal(t, expectedResponse, actualResponse, "Expected response body to match")
}

func TestErrorResponse(t *testing.T) {
	t.Run("Error response with custom status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		message := "Custom error message"
		statusCode := http.StatusBadRequest

		common.ErrorResponse(recorder, message, statusCode)

		assert.Equal(t, statusCode, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response model.BaseResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, message, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("Error response with internal server error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		message := "Internal server error"
		statusCode := http.StatusInternalServerError

		common.ErrorResponse(recorder, message, statusCode)

		assert.Equal(t, statusCode, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response model.BaseResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, message, response.Message)
		assert.Nil(t, response.Data)
	})
}

func TestSuccessResponse(t *testing.T) {
	t.Run("Success response with data", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		testData := map[string]string{"key": "value"}
		message := "Success message"
		statusCode := http.StatusOK

		common.SuccessResponse(recorder, message, testData, statusCode)

		assert.Equal(t, statusCode, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response model.BaseResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, message, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("Success response without data", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		message := "Success without data"
		statusCode := http.StatusOK

		common.SuccessResponse(recorder, message, nil, statusCode)

		assert.Equal(t, statusCode, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response model.BaseResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, message, response.Message)
		assert.Nil(t, response.Data)
	})
}
