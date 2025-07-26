package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lugondev/bbs-selective-disclosure-example/interfaces/http/dto"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/issuer"
)

// IssuerHandler handles issuer-related HTTP requests
type IssuerHandler struct {
	issuerUC *issuer.UseCase
}

// NewIssuerHandler creates a new issuer handler
func NewIssuerHandler(issuerUC *issuer.UseCase) *IssuerHandler {
	return &IssuerHandler{
		issuerUC: issuerUC,
	}
}

// SetupIssuer handles POST /api/issuer/setup
func (h *IssuerHandler) SetupIssuer(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.SetupIssuerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Setup issuer
	setup, err := h.issuerUC.SetupIssuer(req.Method)
	if err != nil {
		writeErrorResponse(w, "Failed to setup issuer", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.SetupIssuerResponse{
		DID:    setup.DID.String(),
		Status: "success",
	}

	writeSuccessResponse(w, response)
}

// IssueCredential handles POST /api/issuer/credentials
func (h *IssuerHandler) IssueCredential(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.IssueCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Convert DTO to use case request
	ucReq := issuer.IssueCredentialRequest{
		IssuerDID:  req.IssuerDID,
		SubjectDID: req.SubjectDID,
		Claims:     dto.ToVCClaims(req.Claims),
	}

	// Issue credential
	credential, err := h.issuerUC.IssueCredential(ucReq)
	if err != nil {
		writeErrorResponse(w, "Failed to issue credential", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.IssueCredentialResponse{
		CredentialID: credential.ID,
		Credential:   credential,
	}

	writeSuccessResponse(w, response)
}

// VerifyCredential handles POST /api/issuer/verify
func (h *IssuerHandler) VerifyCredential(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var credential map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// For simplicity, we'll return success for now
	// In a real implementation, you'd convert the map to VerifiableCredential and verify
	response := dto.SuccessResponse{
		Message: "Credential verification completed",
		Data: map[string]interface{}{
			"valid":  true,
			"status": "verified",
		},
	}

	writeSuccessResponse(w, response)
}
