package handlers

import (
	"net/http"

	"github.com/lugon/bbs-selective-disclosure-example/interfaces/http/dto"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	response := dto.HealthResponse{
		Status:  "healthy",
		Service: "BBS+ Selective Disclosure API",
		Version: "1.0.0",
	}

	writeSuccessResponse(w, response)
}
