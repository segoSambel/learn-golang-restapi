package httputil

import (
	"net/http"

	"github.com/go-chi/httplog/v3"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	_ = JSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	_ = httplog.SetError(r.Context(), err)
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}
