package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of the error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SendError sends a JSON error response with the given status code and error message
func SendError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := ErrorResponse{Error: errorMessage}

	// Encode the error response to JSON and write to the response body
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
