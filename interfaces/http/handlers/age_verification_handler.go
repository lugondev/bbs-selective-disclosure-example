package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lugondev/bbs-selective-disclosure-example/internal/holder"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/issuer"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/verifier"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"
)

type AgeVerificationHandler struct {
	issuerUC   *issuer.UseCase
	holderUC   *holder.UseCase
	verifierUC *verifier.UseCase
}

func NewAgeVerificationHandler(
	issuerUC *issuer.UseCase,
	holderUC *holder.UseCase,
	verifierUC *verifier.UseCase,
) *AgeVerificationHandler {
	return &AgeVerificationHandler{
		issuerUC:   issuerUC,
		holderUC:   holderUC,
		verifierUC: verifierUC,
	}
}

// AgeCredentialRequest represents the request to issue an age verification credential
type AgeCredentialRequest struct {
	IssuerDID   string `json:"issuerDid"`
	SubjectDID  string `json:"subjectDid"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
	Address     string `json:"address"`
	IDNumber    string `json:"idNumber"`
}

// AgeVerificationRequest represents a request for age verification
type AgeVerificationRequest struct {
	HolderDID      string   `json:"holderDid"`
	CredentialID   string   `json:"credentialId"`
	MinAge         int      `json:"minAge"`
	RequiredClaims []string `json:"requiredClaims"`
	ServiceType    string   `json:"serviceType"` // gaming, cinema, alcohol, etc.
}

// AgeVerificationResponse represents the response from age verification
type AgeVerificationResponse struct {
	Success          bool                   `json:"success"`
	AccessGranted    bool                   `json:"accessGranted"`
	ServiceType      string                 `json:"serviceType"`
	MinAgeRequired   int                    `json:"minAgeRequired"`
	AgeVerified      bool                   `json:"ageVerified"`
	RevealedClaims   map[string]interface{} `json:"revealedClaims"`
	HiddenAttributes []string               `json:"hiddenAttributes"`
	PrivacyProtected bool                   `json:"privacyProtected"`
	Message          string                 `json:"message"`
	Error            string                 `json:"error,omitempty"`
}

// POST /api/age-verification/credential - Issue enhanced age verification credential
func (h *AgeVerificationHandler) IssueAgeCredential(w http.ResponseWriter, r *http.Request) {
	var req AgeCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Parse birth date to calculate age-related claims
	birthTime, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		writeErrorResponse(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest, err.Error())
		return
	}

	currentAge := calculateAge(birthTime)
	birthYear := birthTime.Year()

	// Create enhanced claims with age verification
	claims := []vc.Claim{
		// Personal information (will be hidden in age verification)
		{Key: "firstName", Value: req.FirstName},
		{Key: "lastName", Value: req.LastName},
		{Key: "fullName", Value: fmt.Sprintf("%s %s", req.FirstName, req.LastName)},
		{Key: "dateOfBirth", Value: req.DateOfBirth},
		{Key: "nationality", Value: req.Nationality},
		{Key: "address", Value: req.Address},
		{Key: "idNumber", Value: req.IDNumber},

		// Age verification claims (boolean - privacy-preserving)
		{Key: "ageOver13", Value: currentAge >= 13},
		{Key: "ageOver16", Value: currentAge >= 16},
		{Key: "ageOver18", Value: currentAge >= 18},
		{Key: "ageOver21", Value: currentAge >= 21},
		{Key: "ageOver25", Value: currentAge >= 25},
		{Key: "ageOver65", Value: currentAge >= 65},

		// Additional metadata
		{Key: "birthYear", Value: birthYear},
		{Key: "ageCategory", Value: getAgeCategory(currentAge)},
		{Key: "documentType", Value: "national_id"},
		{Key: "issuedAt", Value: time.Now().Format("2006-01-02")},
		{Key: "validUntil", Value: time.Now().AddDate(10, 0, 0).Format("2006-01-02")},
	}

	credential, err := h.issuerUC.IssueCredential(issuer.IssueCredentialRequest{
		IssuerDID:  req.IssuerDID,
		SubjectDID: req.SubjectDID,
		Claims:     claims,
	})
	if err != nil {
		writeErrorResponse(w, "Failed to issue credential", http.StatusInternalServerError, err.Error())
		return
	}

	// Automatically store the credential for the holder
	if err := h.holderUC.StoreCredential(credential); err != nil {
		writeErrorResponse(w, "Failed to store credential", http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"success":         true,
		"credential":      credential,
		"currentAge":      currentAge,
		"ageVerification": map[string]bool{
			"ageOver13": currentAge >= 13,
			"ageOver16": currentAge >= 16,
			"ageOver18": currentAge >= 18,
			"ageOver21": currentAge >= 21,
			"ageOver25": currentAge >= 25,
			"ageOver65": currentAge >= 65,
		},
		"message": fmt.Sprintf("Enhanced age verification credential issued and stored for %d-year-old citizen", currentAge),
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

// POST /api/age-verification/verify - Verify age with privacy preservation
func (h *AgeVerificationHandler) VerifyAge(w http.ResponseWriter, r *http.Request) {
	var req AgeVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Validate input
	if req.HolderDID == "" || req.CredentialID == "" {
		writeErrorResponse(w, "Missing required fields: holderDid and credentialId", http.StatusBadRequest, "")
		return
	}

	// Get the age claim key based on minimum age requirement
	ageClaimKey := getAgeClaimKey(req.MinAge)
	if ageClaimKey == "" {
		writeErrorResponse(w, fmt.Sprintf("Unsupported minimum age: %d", req.MinAge), http.StatusBadRequest, "Supported ages: 13, 16, 18, 21, 25, 65")
		return
	}

	// Create selective disclosure - only reveal age verification and required claims
	revealedAttributes := []string{ageClaimKey, "nationality", "documentType"}
	if req.RequiredClaims != nil {
		revealedAttributes = append(revealedAttributes, req.RequiredClaims...)
	}

	// Remove duplicates
	revealedAttributes = removeDuplicates(revealedAttributes)

	selectiveDisclosure := []vc.SelectiveDisclosureRequest{
		{
			CredentialID:       req.CredentialID,
			RevealedAttributes: revealedAttributes,
		},
	}

	// Generate verification nonce
	verificationNonce := fmt.Sprintf("%s-age-verification-%d", req.ServiceType, time.Now().UnixMilli())

	// Create presentation with error handling
	presentation, err := h.holderUC.CreatePresentation(holder.PresentationRequest{
		HolderDID:           req.HolderDID,
		CredentialIDs:       []string{req.CredentialID},
		SelectiveDisclosure: selectiveDisclosure,
		Nonce:               verificationNonce,
	})
	if err != nil {
		writeErrorResponse(w, "Failed to create presentation", http.StatusInternalServerError, 
			fmt.Sprintf("Could not create presentation for credential %s. Error: %v", req.CredentialID, err))
		return
	}

	// For demo purposes, we'll simulate the verifier part here
	// In a real scenario, this would be done by the verifier service
	// We need to get the issuer DID from the credential for trusted issuers list
	var trustedIssuers []string
	if len(presentation.VerifiableCredential) > 0 {
		if credMap, ok := presentation.VerifiableCredential[0].(map[string]interface{}); ok {
			if issuer, exists := credMap["issuer"]; exists {
				if issuerStr, ok := issuer.(string); ok {
					trustedIssuers = []string{issuerStr}
				}
			}
		}
	}

	verificationResult, err := h.verifierUC.VerifyPresentation(verifier.VerificationRequest{
		Presentation:      presentation,
		RequiredClaims:    []string{ageClaimKey},
		TrustedIssuers:    trustedIssuers,
		VerificationNonce: verificationNonce,
	})
	if err != nil {
		writeErrorResponse(w, "Failed to verify presentation", http.StatusInternalServerError, err.Error())
		return
	}

	// Check age verification result
	var ageVerified bool
	var accessGranted bool
	if ageValue, ok := verificationResult.RevealedClaims[ageClaimKey].(bool); ok {
		ageVerified = ageValue
		accessGranted = ageValue
	}

	// Identify hidden attributes (privacy-protected information)
	hiddenAttributes := []string{
		"firstName", "lastName", "fullName", "dateOfBirth", "address", "idNumber", "birthYear",
	}

	// Generate appropriate message
	message := generateAgeVerificationMessage(req.ServiceType, req.MinAge, accessGranted)

	response := AgeVerificationResponse{
		Success:          verificationResult.Valid,
		AccessGranted:    accessGranted,
		ServiceType:      req.ServiceType,
		MinAgeRequired:   req.MinAge,
		AgeVerified:      ageVerified,
		RevealedClaims:   verificationResult.RevealedClaims,
		HiddenAttributes: hiddenAttributes,
		PrivacyProtected: true,
		Message:          message,
	}

	if !verificationResult.Valid {
		response.Error = "Presentation verification failed"
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// GET /api/age-verification/scenarios - Get supported age verification scenarios
func (h *AgeVerificationHandler) GetAgeScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios := []map[string]interface{}{
		{
			"service":       "Social Media",
			"minAge":        13,
			"claimKey":      "ageOver13",
			"description":   "Access to social media platforms",
			"privacyLevel":  "High - Only verifies 13+ status",
		},
		{
			"service":       "Movie Theater (PG-13)",
			"minAge":        13,
			"claimKey":      "ageOver13",
			"description":   "Access to PG-13 rated movies",
			"privacyLevel":  "High - Only verifies 13+ status",
		},
		{
			"service":       "Movie Theater (R-rated)",
			"minAge":        16,
			"claimKey":      "ageOver16",
			"description":   "Access to R-rated movies (approximated with 16+)",
			"privacyLevel":  "High - Only verifies 16+ status",
		},
		{
			"service":       "Online Gaming",
			"minAge":        18,
			"claimKey":      "ageOver18",
			"description":   "Access to adult gaming content",
			"privacyLevel":  "High - Only verifies 18+ status",
		},
		{
			"service":       "Alcohol Purchase",
			"minAge":        21,
			"claimKey":      "ageOver21",
			"description":   "Purchase alcoholic beverages",
			"privacyLevel":  "High - Only verifies 21+ status",
		},
		{
			"service":       "Senior Discount",
			"minAge":        65,
			"claimKey":      "ageOver65",
			"description":   "Eligibility for senior citizen discounts",
			"privacyLevel":  "High - Only verifies 65+ status",
		},
	}

	response := map[string]interface{}{
		"scenarios": scenarios,
		"privacy_benefits": []string{
			"Exact age never revealed",
			"Birth date remains private",
			"Personal information hidden",
			"Only boolean age verification disclosed",
			"Unlinkable presentations",
			"Zero-knowledge age proof",
		},
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// POST /api/age-verification/demo - Run complete age verification demo
func (h *AgeVerificationHandler) RunAgeDemo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServiceType string `json:"serviceType"`
		MinAge      int    `json:"minAge"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.ServiceType = "gaming"
		req.MinAge = 18
	}

	// This would run the complete demo flow
	// For now, return a demo execution plan
	response := map[string]interface{}{
		"demo_plan": []map[string]string{
			{
				"step":        "1",
				"title":       "Setup Government Authority",
				"description": "Initialize digital ID issuer with BBS+ keys",
				"status":      "ready",
			},
			{
				"step":        "2",
				"title":       "Setup Citizen",
				"description": "Create citizen DID for credential holder",
				"status":      "ready",
			},
			{
				"step":        "3",
				"title":       fmt.Sprintf("Setup %s Service", req.ServiceType),
				"description": fmt.Sprintf("Initialize %s platform for age verification", req.ServiceType),
				"status":      "ready",
			},
			{
				"step":        "4",
				"title":       "Issue Enhanced ID",
				"description": "Government issues digital ID with age verification claims",
				"status":      "ready",
			},
			{
				"step":        "5",
				"title":       "Age Verification Request",
				"description": fmt.Sprintf("Service requests %d+ age verification", req.MinAge),
				"status":      "ready",
			},
			{
				"step":        "6",
				"title":       "Privacy-Preserving Presentation",
				"description": "Citizen creates selective disclosure presentation",
				"status":      "ready",
			},
			{
				"step":        "7",
				"title":       "Verification",
				"description": "Service verifies age without seeing personal details",
				"status":      "ready",
			},
		},
		"service_type": req.ServiceType,
		"min_age":      req.MinAge,
		"privacy_protection": map[string]interface{}{
			"revealed": []string{
				fmt.Sprintf("ageOver%d (boolean)", req.MinAge),
				"nationality",
				"documentType",
			},
			"hidden": []string{
				"firstName", "lastName", "dateOfBirth", "exactAge",
				"address", "idNumber", "birthYear",
			},
		},
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// Helper functions

func calculateAge(birthTime time.Time) int {
	now := time.Now()
	age := now.Year() - birthTime.Year()
	if now.YearDay() < birthTime.YearDay() {
		age--
	}
	return age
}

func getAgeCategory(age int) string {
	switch {
	case age < 13:
		return "child"
	case age < 18:
		return "teen"
	case age < 65:
		return "adult"
	default:
		return "senior"
	}
}

func getAgeClaimKey(minAge int) string {
	switch minAge {
	case 13:
		return "ageOver13"
	case 16:
		return "ageOver16"
	case 18:
		return "ageOver18"
	case 21:
		return "ageOver21"
	case 25:
		return "ageOver25"
	case 65:
		return "ageOver65"
	default:
		return ""
	}
}

func generateAgeVerificationMessage(serviceType string, minAge int, accessGranted bool) string {
	if accessGranted {
		return fmt.Sprintf("🎉 ACCESS GRANTED: User verified to be %d+ years old for %s service. Privacy protected - exact age and personal details remain hidden.", minAge, serviceType)
	}
	return fmt.Sprintf("❌ ACCESS DENIED: User is under %d years old for %s service.", minAge, serviceType)
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}
