package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lugon/bbs-selective-disclosure-example/interfaces/http/dto"
	"github.com/lugon/bbs-selective-disclosure-example/internal/verifier"
)

// VerifierHandler handles verifier-related HTTP requests
type VerifierHandler struct {
	verifierUC *verifier.UseCase
}

// NewVerifierHandler creates a new verifier handler
func NewVerifierHandler(verifierUC *verifier.UseCase) *VerifierHandler {
	return &VerifierHandler{
		verifierUC: verifierUC,
	}
}

// SetupVerifier handles POST /api/verifier/setup
func (h *VerifierHandler) SetupVerifier(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.SetupVerifierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Setup verifier
	setup, err := h.verifierUC.SetupVerifier(req.Method)
	if err != nil {
		writeErrorResponse(w, "Failed to setup verifier", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.SetupVerifierResponse{
		DID:    setup.DID.String(),
		Status: "success",
	}

	writeSuccessResponse(w, response)
}

// VerifyPresentation handles POST /api/verifier/verify
func (h *VerifierHandler) VerifyPresentation(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.VerifyPresentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Convert DTO to use case request
	ucReq := verifier.VerificationRequest{
		Presentation:      req.Presentation,
		RequiredClaims:    req.RequiredClaims,
		TrustedIssuers:    req.TrustedIssuers,
		VerificationNonce: req.VerificationNonce,
	}

	// Verify presentation
	result, err := h.verifierUC.VerifyPresentation(ucReq)
	if err != nil {
		writeErrorResponse(w, "Failed to verify presentation", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.VerifyPresentationResponse{
		Valid:           result.Valid,
		Errors:          result.Errors,
		RevealedClaims:  result.RevealedClaims,
		HolderDID:       result.HolderDID,
		IssuerDIDs:      result.IssuerDIDs,
		CredentialTypes: result.CredentialTypes,
	}

	writeSuccessResponse(w, response)
}

// CreateVerificationRequest handles POST /api/verifier/verification-request
func (h *VerifierHandler) CreateVerificationRequest(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.CreateVerificationRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Convert DTO to use case request
	params := verifier.CreateVerificationRequestParams{
		RequiredClaims:    req.RequiredClaims,
		TrustedIssuers:    req.TrustedIssuers,
		VerificationNonce: req.VerificationNonce,
	}

	// Create verification request
	result, err := h.verifierUC.CreateVerificationRequest(params)
	if err != nil {
		writeErrorResponse(w, "Failed to create verification request", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.CreateVerificationRequestResponse{
		RequiredClaims:    result.RequiredClaims,
		TrustedIssuers:    result.TrustedIssuers,
		VerificationNonce: result.VerificationNonce,
	}

	writeSuccessResponse(w, response)
}

// ListPresentations handles GET /api/verifier/presentations?verifierDid={did}
func (h *VerifierHandler) ListPresentations(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	verifierDID := r.URL.Query().Get("verifierDid")
	if verifierDID == "" {
		writeErrorResponse(w, "verifierDid parameter is required", http.StatusBadRequest, "")
		return
	}

	// List presentations
	presentations, err := h.verifierUC.ListVerifiedPresentations(verifierDID)
	if err != nil {
		writeErrorResponse(w, "Failed to list presentations", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ListPresentationsResponse{
		Presentations: presentations,
	}

	writeSuccessResponse(w, response)
}
