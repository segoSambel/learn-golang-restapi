package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/segoSambel/learn-golang-restapi/internal/httputil"
)

type userIDKey struct{}

func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}

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

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)

	return r
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !httputil.Bind(w, r, h.validate, &req) {
		return
	}

	response, err := h.service.Login(r.Context(), req)
	if errors.Is(err, ErrInvalidCredentials) {
		httputil.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		return
	}
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	_ = httputil.JSON(w, http.StatusOK, response)
}

func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, isBearer := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")

		userID, err := h.service.Verify(accessToken)
		if !isBearer || err != nil {
			httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing or invalid access token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
