package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type errorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func ErrorJSON(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, errorResponse{Message: message})
}

func ValidationErrorJSON(w http.ResponseWriter, err error) {
	fields := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		fields[field] = e.Tag()
	}

	WriteJSON(w, http.StatusBadRequest, errorResponse{
		Message: "validation failed",
		Errors:  fields,
	})
}
