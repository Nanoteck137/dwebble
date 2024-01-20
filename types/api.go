package types

import (
	"net/http"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err ApiError) Error() string {
	return err.Message
}

func NewApiError(status int, message string, data ...any) ApiError {
	var d any = nil

	if data != nil && len(data) > 0 && data[0] != nil {
		d = data[0]
	}

	return ApiError{
		Status:  status,
		Message: message,
		Data:    d,
	}
}

func ApiBadRequestError(message string, data ...any) ApiError {
	return NewApiError(http.StatusBadRequest, message, data...)
}

type ApiResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func NewApiResponse(data any) ApiResponse {
	return ApiResponse{
		Status: 200,
		Data:   data,
	}
}
