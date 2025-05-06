package chi

import (
	"net/http"
	"strings"
	"github.com/go-playground/validator/v10"
	"github.com/joaoramos09/go-ai/ia"
	"github.com/joaoramos09/go-ai/internal/errs"
	"github.com/joaoramos09/go-ai/internal/json"
)

type invokeRequest struct {
	Model string `json:"model"`
	Input string `json:"input" validate:"required"`
}

type invokeResponse struct {
	Response string `json:"response"`
}

type documentRequest struct {
	Text string `json:"text" validate:"required"`
	Category string `json:"category" validate:"required"`
}

type queryRequest struct {
	Input string `json:"input" validate:"required"`
	Model string `json:"model"`
}

func invokeAI(ai ia.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request invokeRequest
		err := json.Read(w, r, &request)
		if err != nil {
			json.WriteError(w, http.StatusBadRequest, errs.ErrSintaxError.Error())
			return
		}

		if err := Validator.Struct(request); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field())
			}
			json.WriteErrorWithParams(w, http.StatusBadRequest, errs.ErrInvalidParams.Error(), strings.Join(errors, ", "))
			return
		}

		response, err := ai.Invoke(r.Context(), request.Input, request.Model)
		if err != nil {
			json.WriteError(w, http.StatusInternalServerError, errs.ErrAIInvoke.Error())
			return
		}

		var result invokeResponse
		result.Response = response

		json.Write(w, http.StatusOK, result)
	}
}

func insertDocuments(ai ia.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requests []documentRequest
		err := json.Read(w, r, &requests)
		if err != nil {
			json.WriteError(w, http.StatusBadRequest, errs.ErrSintaxError.Error())
			return
		}

		for _, request := range requests {
			if err := Validator.Struct(request); err != nil {
				var errors []string
				for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field())
			}
				json.WriteErrorWithParams(w, http.StatusBadRequest, errs.ErrInvalidParams.Error(), strings.Join(errors, ", "))
				return
			}
		}

		documents := make([]ia.Document, len(requests))

		for i, request := range requests {
			documents[i] = ia.Document{
				Text: request.Text,
				Category: request.Category,
			}
		}
		if err := ai.Insert(r.Context(), documents); err != nil {
			json.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		json.Write(w, http.StatusOK, nil)
	}
}

func queryDocuments(ai ia.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request queryRequest
		err := json.Read(w, r, &request)
		if err != nil {
			json.WriteError(w, http.StatusBadRequest, errs.ErrSintaxError.Error())
			return
		}

		if err := Validator.Struct(request); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field())
			}
			json.WriteErrorWithParams(w, http.StatusBadRequest, errs.ErrInvalidParams.Error(), strings.Join(errors, ", "))
			return
		}

		documents, err := ai.Query(r.Context(), request.Input)
		if err != nil {
			json.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response, err := ai.InvokeWithSystemPrompt(r.Context(), request.Input, documents, request.Model)
		if err != nil {
			json.WriteError(w, http.StatusInternalServerError, errs.ErrAIInvoke.Error())
			return
		}

		var result invokeResponse
		result.Response = response
		
		json.Write(w, http.StatusOK, result)
	}
}
