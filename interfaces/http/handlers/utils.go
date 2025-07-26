package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lugondev/bbs-selective-disclosure-example/interfaces/http/dto"
)

// writeErrorResponse writes an error response to the HTTP response writer
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := dto.ErrorResponse{
		Error:   message,
		Code:    statusCode,
		Details: details,
	}

	json.NewEncoder(w).Encode(errorResp)
}

// writeSuccessResponse writes a success response to the HTTP response writer
func writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}

// writeJSONResponse writes a JSON response with custom status code
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(data)
}

// enableCORS enables CORS for the response
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
