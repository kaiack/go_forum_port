package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// handleValidationError processes validation errors and sends a JSON response with error details.
func HandleValidationError(err error, w http.ResponseWriter) {
	// If validation fails, return a JSON response with the error messages
	var validationErrors []string
	for _, e := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", e.Field(), e.Tag()))
	}

	// Set response header and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	// Send the validation error response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Validation failed",
		"errors":  validationErrors,
	})
}
