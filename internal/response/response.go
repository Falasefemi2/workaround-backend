package response

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Data any `json:"data"`
}

type errorResponse struct {
	Error  string `json:"error"`
	Fields any    `json:"fields,omitempty"` // only appears on validation errors
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(successResponse{Data: data})
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func ValidationError(w http.ResponseWriter, fields any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(errorResponse{
		Error:  "validation failed",
		Fields: fields,
	})
}
