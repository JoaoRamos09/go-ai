package json

import (
	"encoding/json"
	"net/http"
	"github.com/joaoramos09/go-ai/internal/errs"
	"log"
)

func Write(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func Read(w http.ResponseWriter, r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(v)
	if err != nil {
		log.Printf("Error reading request: %v", err)
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, statusCode int, message string) error {
	return Write(w, statusCode, &errs.DefaultError{Message: message})
}

func WriteErrorWithParams(w http.ResponseWriter, statusCode int, message string, params string) error {
	log.Printf("Error with params: %v, %v", message, params)
	return Write(w, statusCode, &errs.InvalidParamsError{Message: message, Params: params})
}




