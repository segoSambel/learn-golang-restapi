package user

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/segoSambel/learn-golang-restapi/internal/auth"
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

func (h *Handler) Routes(requireAuth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.Create)
	r.With(requireAuth).Get("/me", h.Me)

	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if !httputil.Bind(w, r, h.validate, &req) {
		return
	}

	response, err := h.service.Create(r.Context(), req)
	if errors.Is(err, ErrEmailTaken) {
		httputil.Error(w, http.StatusConflict, "EMAIL_TAKEN", "email is already registered")
		return
	}
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	_ = httputil.JSON(w, http.StatusCreated, response)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.Get(r.Context(), auth.UserID(r.Context()))
	if errors.Is(err, ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		return
	}
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	_ = httputil.JSON(w, http.StatusOK, response)
}
