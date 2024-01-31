package common

import "net/http"

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
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
