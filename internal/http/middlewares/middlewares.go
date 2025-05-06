package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joaoramos09/go-ai/internal/errs"
	"github.com/joaoramos09/go-ai/internal/json"
	"github.com/joaoramos09/go-ai/internal/user"
	"github.com/joaoramos09/go-ai/internal/user/auth"
)

type UseCase interface {
	Authenticate(next http.Handler) http.Handler
	RequireRole(requiredRole user.Role) func(http.Handler) http.Handler
}

type Service struct {
	authService auth.UseCase
	userService user.UseCase
}

func NewService(authService auth.UseCase, userService user.UseCase) UseCase {
	return &Service{
		authService: authService,
		userService: userService,
	}
}

func (s *Service) RequireRole(requiredRole user.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := s.userService.GetFromContext(r.Context())
			if err != nil {
				log.Printf("[MIDDLEWARE] Error getting user: %v", err)
				json.WriteError(w, http.StatusInternalServerError, "")
				return
			}

			if user.Role.Value() < requiredRole.Value() {
				log.Printf("[MIDDLEWARE] User does not have %v role", requiredRole)
				json.WriteError(w, http.StatusForbidden, errs.ErrUnauthorized.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (s *Service) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" {
			log.Printf("[MIDDLEWARE] Token is required")
			json.WriteError(w, http.StatusBadRequest, errs.ErrTokenRequired.Error())
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[MIDDLEWARE] Invalid token structure: %v", parts)
			json.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidTokenStructure.Error())
			return
		}

		token = parts[1]

		jwtToken, err := s.authService.Validate(token)
		if err != nil {
			log.Printf("[MIDDLEWARE] Invalid token: %v", err)
			json.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidToken.Error())
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !jwtToken.Valid {
			log.Printf("[MIDDLEWARE] Invalid token claims: %v", err)
			json.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidTokenClaims.Error())
			return
		}

		if claims["exp"] == nil {
			log.Printf("[MIDDLEWARE] Invalid token claims: %v", err)
			json.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidTokenClaims.Error())
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			log.Printf("[MIDDLEWARE] Invalid token claims: %v", err)
			json.WriteError(w, http.StatusInternalServerError, "")
			return
		}

		if time.Now().Unix() > int64(exp) {
			log.Printf("[MIDDLEWARE] Token expired")
			json.WriteError(w, http.StatusUnauthorized, errs.ErrTokenExpired.Error())
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			log.Printf("[MIDDLEWARE] Invalid token claims: %v", err)
			json.WriteError(w, http.StatusInternalServerError, "")
			return
		}

		cxt := context.WithValue(r.Context(), user.UserIDKey, int(userID))
		next.ServeHTTP(w, r.WithContext(cxt))
	})
}
