package dto

import "github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"

// SetupIssuerRequest represents the request to setup an issuer
type SetupIssuerRequest struct {
	Method      string `json:"method" validate:"required"`
	BBSProvider string `json:"bbsProvider,omitempty"`
}

// SetupIssuerResponse represents the response from setting up an issuer
type SetupIssuerResponse struct {
	DID    string `json:"did"`
	Status string `json:"status"`
}

// IssueCredentialRequest represents the request to issue a credential
type IssueCredentialRequest struct {
	IssuerDID   string     `json:"issuerDid" validate:"required"`
	SubjectDID  string     `json:"subjectDid" validate:"required"`
	Claims      []ClaimDTO `json:"claims" validate:"required,min=1"`
	BBSProvider string     `json:"bbsProvider,omitempty"`
}

// ClaimDTO represents a claim in the credential
type ClaimDTO struct {
	Key   string      `json:"key" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

// IssueCredentialResponse represents the response from issuing a credential
type IssueCredentialResponse struct {
	CredentialID string                   `json:"credentialId"`
	Credential   *vc.VerifiableCredential `json:"credential"`
}

// ToVCClaims converts ClaimDTO slice to vc.Claim slice
func ToVCClaims(claims []ClaimDTO) []vc.Claim {
	vcClaims := make([]vc.Claim, len(claims))
	for i, claim := range claims {
		vcClaims[i] = vc.Claim{
			Key:   claim.Key,
			Value: claim.Value,
		}
	}
	return vcClaims
}
