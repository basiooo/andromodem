package common

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/go-playground/validator/v10"
)

func DeviceNotFoundResponse(writer http.ResponseWriter) {
	WriteToResponseBody(writer, model.BaseResponse{
		Success: false,
		Message: "Device not found",
	}, http.StatusNotFound)
}

func ErrorResponse(writter http.ResponseWriter, message string, statusCode int) {
	WriteToResponseBody(writter, model.BaseResponse{
		Success: false,
		Message: message,
	}, statusCode)
}

func ValidationErrorResponse(writer http.ResponseWriter, message string, statusCode int, err error) {
	errs := make(map[string]string)

	for _, fieldError := range err.(validator.ValidationErrors) {
		field := fieldError.Field()
		tag := fieldError.Tag()
		param := fieldError.Param()
		switch tag {
		case "required":
			errs[field] = field + " is required"
		case "email":
			errs[field] = "must be a valid email address"
		case "url":
			errs[field] = "must be a valid URL"
		case "uuid":
			errs[field] = "must be a valid UUID"
		case "min":
			errs[field] = field + " must be at least " + param + " characters"
		case "max":
			errs[field] = field + " must be at most " + param + " characters"
		case "len":
			errs[field] = field + " must be exactly " + param + " characters"
		case "eq":
			errs[field] = field + " must be equal to " + param
		case "ne":
			errs[field] = field + " must not be equal to " + param
		case "gt":
			errs[field] = field + " must be greater than " + param
		case "gte":
			errs[field] = field + " must be greater than or equal to " + param
		case "lt":
			errs[field] = field + " must be less than " + param
		case "lte":
			errs[field] = field + " must be less than or equal to " + param
		case "oneof":
			errs[field] = field + " must be one of: " + param
		case "notin":
			errs[field] = "field " + field + " cannot be one of: " + param
		default:
			errs[field] = field + " is invalid"
		}
	}

	WriteToResponseBody(writer, model.BaseResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	}, statusCode)
}

func SuccessResponse(writter http.ResponseWriter, message string, data any, statusCode int) {
	WriteToResponseBody(writter, model.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	}, statusCode)
}

func SSESetResponseHeader(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	writer.Header().Set("Content-Type", "text/event-stream")
}
