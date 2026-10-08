package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/httplog/v3"
	"github.com/go-playground/validator/v10"
	"github.com/segoSambel/learn-golang-restapi/internal/httputil"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", validationMessage(err))
		return
	}

	response, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			httputil.Error(w, http.StatusConflict, "EMAIL_TAKEN", "email is already registered")
			return
		}

		_ = httplog.SetError(r.Context(), err)
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	_ = httputil.JSON(w, http.StatusCreated, response)
}

func validationMessage(err error) string {
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) || len(fieldErrors) == 0 {
		return "request validation failed"
	}

	first := fieldErrors[0]
	return fmt.Sprintf("%s is invalid (%s)", first.Field(), first.Tag())
}
