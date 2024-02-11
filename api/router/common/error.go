package common

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

var (
	QueryNoResultErr                      = errors.New("no record")
	FailUpdateError                       = errors.New("update failure")
	FailCreateError                       = errors.New("create failure")
	FailDeleteError                       = errors.New("delete failure")
	DuplicateValueError                   = errors.New("duplicate value")
	UUIDBadRequest                        = errors.New("invalid uuid")
	CPassAndPassNotMatch                  = errors.New("password and confirmed password must be the same")
	TokenExpired                          = errors.New("expired token")
	OldAccessTokenAndRefreshTokenNotMatch = errors.New("expired access token and refresh token to request new access token not match")
)

func BadRequestResponse(err error) ErrorResponse {
	msgErr := "bad request"

	if err != nil {
		msgErr = err.Error()
	}

	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Error:   "Bad Request Error",
		Message: msgErr,
	}
}

func InternalErrorResponse(err error) ErrorResponse {
	msgErr := "Internal error"

	if err != nil {
		msgErr = err.Error()
	}

	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Error:   "Internal Server Error",
		Message: msgErr,
	}
}

func NotFoundErrorResponse(err error) ErrorResponse {
	msgErr := "api not found"

	if err != nil {
		msgErr = err.Error()
	}

	return ErrorResponse{
		Status:  http.StatusNotFound,
		Error:   "Not Found",
		Message: msgErr,
	}
}

func UnauthorizedErrorResponse(err error) ErrorResponse {
	msgErr := "unauthorized request"

	if err != nil {
		msgErr = err.Error()
	}

	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Error:   "Unauthorized Request",
		Message: msgErr,
	}
}
