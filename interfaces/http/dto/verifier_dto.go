package dto

import "github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"

// SetupVerifierRequest represents the request to setup a verifier
type SetupVerifierRequest struct {
	Method string `json:"method" validate:"required"`
}

// SetupVerifierResponse represents the response from setting up a verifier
type SetupVerifierResponse struct {
	DID    string `json:"did"`
	Status string `json:"status"`
}

// VerifyPresentationRequest represents the request to verify a presentation
type VerifyPresentationRequest struct {
	Presentation      *vc.VerifiablePresentation `json:"presentation" validate:"required"`
	RequiredClaims    []string                   `json:"requiredClaims"`
	TrustedIssuers    []string                   `json:"trustedIssuers"`
	VerificationNonce string                     `json:"verificationNonce"`
}

// VerifyPresentationResponse represents the response from verifying a presentation
type VerifyPresentationResponse struct {
	Valid           bool                   `json:"valid"`
	Errors          []string               `json:"errors,omitempty"`
	RevealedClaims  map[string]interface{} `json:"revealedClaims,omitempty"`
	HolderDID       string                 `json:"holderDid"`
	IssuerDIDs      []string               `json:"issuerDids"`
	CredentialTypes []string               `json:"credentialTypes"`
}

// CreateVerificationRequestRequest represents the request to create a verification request
type CreateVerificationRequestRequest struct {
	RequiredClaims    []string `json:"requiredClaims" validate:"required,min=1"`
	TrustedIssuers    []string `json:"trustedIssuers"`
	VerificationNonce string   `json:"verificationNonce"`
}

// CreateVerificationRequestResponse represents the response from creating a verification request
type CreateVerificationRequestResponse struct {
	RequiredClaims    []string `json:"requiredClaims"`
	TrustedIssuers    []string `json:"trustedIssuers"`
	VerificationNonce string   `json:"verificationNonce"`
}

// ListPresentationsResponse represents the response from listing presentations
type ListPresentationsResponse struct {
	Presentations []*vc.VerifiablePresentation `json:"presentations"`
}
