package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	"github.com/segoSambel/learn-golang-restapi/internal/auth"
	"github.com/segoSambel/learn-golang-restapi/internal/httputil"
	"github.com/segoSambel/learn-golang-restapi/internal/user"
)

func New(logger *slog.Logger, authHandler *auth.Handler, userHandler *user.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		Level:         slog.LevelInfo,
		Schema:        httplog.SchemaECS.Concise(true),
		RecoverPanics: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = httputil.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())
		r.Mount("/users", userHandler.Routes(authHandler.RequireAuth))
	})

	return r
}
