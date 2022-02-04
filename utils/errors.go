package utils

import "net/http"

type HTTPError struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func InternalServerErrorHandler() HTTPError {
	return HTTPError{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
	}
}
