package dto

import "github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"

// SetupHolderRequest represents the request to setup a holder
type SetupHolderRequest struct {
	Method      string `json:"method" validate:"required"`
	BBSProvider string `json:"bbsProvider,omitempty"`
}

// SetupHolderResponse represents the response from setting up a holder
type SetupHolderResponse struct {
	DID    string `json:"did"`
	Status string `json:"status"`
}

// StoreCredentialRequest represents the request to store a credential
type StoreCredentialRequest struct {
	Credential *vc.VerifiableCredential `json:"credential" validate:"required"`
}

// StoreCredentialResponse represents the response from storing a credential
type StoreCredentialResponse struct {
	Status string `json:"status"`
}

// CreatePresentationRequest represents the request to create a presentation
type CreatePresentationRequest struct {
	HolderDID           string                          `json:"holderDid" validate:"required"`
	CredentialIDs       []string                        `json:"credentialIds" validate:"required,min=1"`
	SelectiveDisclosure []SelectiveDisclosureRequestDTO `json:"selectiveDisclosure" validate:"required,min=1"`
	Nonce               string                          `json:"nonce,omitempty"`
	BBSProvider         string                          `json:"bbsProvider,omitempty"`
}

// SelectiveDisclosureRequestDTO represents a selective disclosure request
type SelectiveDisclosureRequestDTO struct {
	CredentialID       string   `json:"credentialId" validate:"required"`
	RevealedAttributes []string `json:"revealedAttributes" validate:"required,min=1"`
	Nonce              string   `json:"nonce,omitempty"`
}

// CreatePresentationResponse represents the response from creating a presentation
type CreatePresentationResponse struct {
	PresentationID string                     `json:"presentationId"`
	Presentation   *vc.VerifiablePresentation `json:"presentation"`
}

// ListCredentialsResponse represents the response from listing credentials
type ListCredentialsResponse struct {
	Credentials []*vc.VerifiableCredential `json:"credentials"`
}

// ToVCSelectiveDisclosure converts DTO to vc.SelectiveDisclosureRequest slice
func ToVCSelectiveDisclosure(dtos []SelectiveDisclosureRequestDTO) []vc.SelectiveDisclosureRequest {
	vcReqs := make([]vc.SelectiveDisclosureRequest, len(dtos))
	for i, dto := range dtos {
		vcReqs[i] = vc.SelectiveDisclosureRequest{
			CredentialID:       dto.CredentialID,
			RevealedAttributes: dto.RevealedAttributes,
			Nonce:              dto.Nonce,
		}
	}
	return vcReqs
}
