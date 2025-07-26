package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lugon/bbs-selective-disclosure-example/interfaces/http/dto"
	"github.com/lugon/bbs-selective-disclosure-example/internal/holder"
)

// HolderHandler handles holder-related HTTP requests
type HolderHandler struct {
	holderUC *holder.UseCase
}

// NewHolderHandler creates a new holder handler
func NewHolderHandler(holderUC *holder.UseCase) *HolderHandler {
	return &HolderHandler{
		holderUC: holderUC,
	}
}

// SetupHolder handles POST /api/holder/setup
func (h *HolderHandler) SetupHolder(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.SetupHolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Setup holder
	setup, err := h.holderUC.SetupHolder(req.Method)
	if err != nil {
		writeErrorResponse(w, "Failed to setup holder", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.SetupHolderResponse{
		DID:    setup.DID.String(),
		Status: "success",
	}

	writeSuccessResponse(w, response)
}

// StoreCredential handles POST /api/holder/credentials
func (h *HolderHandler) StoreCredential(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.StoreCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Store credential
	if err := h.holderUC.StoreCredential(req.Credential); err != nil {
		writeErrorResponse(w, "Failed to store credential", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.StoreCredentialResponse{
		Status: "success",
	}

	writeSuccessResponse(w, response)
}

// CreatePresentation handles POST /api/holder/presentations
func (h *HolderHandler) CreatePresentation(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.CreatePresentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Convert DTO to use case request
	ucReq := holder.PresentationRequest{
		HolderDID:           req.HolderDID,
		CredentialIDs:       req.CredentialIDs,
		SelectiveDisclosure: dto.ToVCSelectiveDisclosure(req.SelectiveDisclosure),
	}

	// Create presentation
	presentation, err := h.holderUC.CreatePresentation(ucReq)
	if err != nil {
		writeErrorResponse(w, "Failed to create presentation", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.CreatePresentationResponse{
		PresentationID: presentation.ID,
		Presentation:   presentation,
	}

	writeSuccessResponse(w, response)
}

// ListCredentials handles GET /api/holder/credentials?holderDid={did}
func (h *HolderHandler) ListCredentials(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	holderDID := r.URL.Query().Get("holderDid")
	if holderDID == "" {
		writeErrorResponse(w, "holderDid parameter is required", http.StatusBadRequest, "")
		return
	}

	// List credentials
	credentials, err := h.holderUC.ListCredentials(holderDID)
	if err != nil {
		writeErrorResponse(w, "Failed to list credentials", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ListCredentialsResponse{
		Credentials: credentials,
	}

	writeSuccessResponse(w, response)
}
