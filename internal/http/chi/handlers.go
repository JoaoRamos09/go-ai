package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	middle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/joaoramos09/go-ai/ia"
	"github.com/joaoramos09/go-ai/internal/http/middlewares"
	"github.com/joaoramos09/go-ai/internal/user"
	"github.com/joaoramos09/go-ai/internal/user/auth"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func Handlers(aiService ia.UseCase, userService user.UseCase, authService auth.UseCase, middlewareService middlewares.UseCase) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middle.Recoverer)
	r.Use(middle.Logger)
	r.Use(middle.RequestID)
	r.Use(middle.RealIP)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Method(http.MethodPost, "/register", createUser(userService, authService))
			r.Method(http.MethodPost, "/login", loginUser(userService, authService))
		})
		r.Route("/ai", func(r chi.Router) {
			r.Use(middlewareService.Authenticate)
			r.Method(http.MethodPost, "/invoke", invokeAI(aiService))
			r.Route("/documents", func(r chi.Router) {
				r.With(middlewareService.RequireRole(user.RoleAdmin)).Method(http.MethodPost, "/", insertDocuments(aiService))
				r.Method(http.MethodGet, "/query", queryDocuments(aiService))
			})
		})
	})

	return r
}
