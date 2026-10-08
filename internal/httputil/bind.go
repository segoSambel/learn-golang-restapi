package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const maxBodyBytes = 1 << 20

func Bind(w http.ResponseWriter, r *http.Request, validate *validator.Validate, dst any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return false
	}

	if err := validate.Struct(dst); err != nil {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", validationMessage(err))
		return false
	}

	return true
}

func validationMessage(err error) string {
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) || len(fieldErrors) == 0 {
		return "request validation failed"
	}

	first := fieldErrors[0]
	return fmt.Sprintf("%s is invalid (%s)", first.Field(), first.Tag())
}
