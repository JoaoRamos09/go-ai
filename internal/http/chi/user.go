package chi

import (
	"log"
	"net/http"
	"github.com/joaoramos09/go-ai/internal/user"
	"github.com/joaoramos09/go-ai/internal/user/auth"
	"github.com/joaoramos09/go-ai/internal/errs"
	"github.com/joaoramos09/go-ai/internal/json"
	"github.com/go-playground/validator/v10"
	"strings"
)

type userRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type userResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func createUser(userService user.UseCase, authService auth.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req userRequest
		err := json.Read(w, r, &req)
		if err != nil {
			log.Printf("Error reading request: %v", err)
			json.WriteError(w, http.StatusBadRequest, errs.ErrSintaxError.Error())
			return
		}

		validate := Validator.Struct(req)
		if validate != nil {
			var errors []string
			for _, err := range validate.(validator.ValidationErrors) {
				errors = append(errors, err.Field())
			}
			json.WriteErrorWithParams(w, http.StatusBadRequest, errs.ErrInvalidParams.Error(), strings.Join(errors, ", "))
			return
		}

		hash, err := authService.EncodePassword(req.Password)
		if err != nil {
			json.WriteError(w, http.StatusInternalServerError, "")
			return
		}

		ur := user.User{
			Username: req.Username,
			Email:    req.Email,
			Password: hash,
		}

		u, err := userService.Create(r.Context(), &ur)
		if err != nil {
			switch err {
			case errs.ErrUserAlreadyExists:
				json.WriteError(w, http.StatusBadRequest, errs.ErrUserAlreadyExists.Error())
			default:
				json.WriteError(w, http.StatusInternalServerError, "")
			}
			return
		}

		res := userResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		}

		json.Write(w, http.StatusCreated, res)
	}
}

func loginUser(userService user.UseCase, authService auth.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		err := json.Read(w, r, &req)
		if err != nil {
			json.WriteError(w, http.StatusBadRequest, errs.ErrSintaxError.Error())
			return
		}

		validate := Validator.Struct(req)
		if validate != nil {
			var errors []string
			for _, err := range validate.(validator.ValidationErrors) {
				errors = append(errors, err.Field())
			}
			json.WriteErrorWithParams(w, http.StatusBadRequest, errs.ErrInvalidParams.Error(), strings.Join(errors, ", "))
			return
		}

		u, err := userService.GetByEmail(r.Context(), req.Email)
		if err != nil {
			switch err {
			case errs.ErrUserNotFound:
				json.WriteError(w, http.StatusNotFound, errs.ErrUserNotFound.Error())
			default:
				json.WriteError(w, http.StatusInternalServerError, "")
			}
			return
		}

		match := authService.MatchPassword(req.Password, u.Password)
		if !match {
			json.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidCredentials.Error())
			return
		}
		
		token, err := authService.Generate(u.ID)
		if err != nil {
			json.WriteError(w, http.StatusInternalServerError, "")
			return
		}

		json.Write(w, http.StatusOK, token)
		
	}
}
