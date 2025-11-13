package payment

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    code,
	})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, send a plain text error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
	}
}

// respondWithSuccess sends a success response
func respondWithSuccess(w http.ResponseWriter, data interface{}, message string) {
	respondWithJSON(w, http.StatusOK, SuccessResponse{
		Data:    data,
		Message: message,
	})
}
