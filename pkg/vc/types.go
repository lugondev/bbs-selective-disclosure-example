package vc

import (
	"time"

	"github.com/lugon/bbs-selective-disclosure-example/pkg/bbs"
)

// VerifiableCredential represents a W3C Verifiable Credential
type VerifiableCredential struct {
	Context           []string               `json:"@context"`
	ID                string                 `json:"id"`
	Type              []string               `json:"type"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      time.Time              `json:"issuanceDate"`
	ExpirationDate    *time.Time             `json:"expirationDate,omitempty"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Proof             *Proof                 `json:"proof,omitempty"`
}

// VerifiablePresentation represents a W3C Verifiable Presentation
type VerifiablePresentation struct {
	Context              []string      `json:"@context"`
	ID                   string        `json:"id"`
	Type                 []string      `json:"type"`
	Holder               string        `json:"holder"`
	VerifiableCredential []interface{} `json:"verifiableCredential"`
	Proof                *Proof        `json:"proof,omitempty"`
}

// Proof represents a cryptographic proof
type Proof struct {
	Type               string    `json:"type"`
	Created            time.Time `json:"created"`
	VerificationMethod string    `json:"verificationMethod"`
	ProofPurpose       string    `json:"proofPurpose"`
	ProofValue         string    `json:"proofValue,omitempty"`
	// BBS+ specific fields
	Nonce              string `json:"nonce,omitempty"`
	RevealedAttributes []int  `json:"revealedAttributes,omitempty"`
}

// Claim represents a single claim in a credential
type Claim struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SelectiveDisclosureRequest represents what attributes to reveal
type SelectiveDisclosureRequest struct {
	CredentialID       string   `json:"credentialId"`
	RevealedAttributes []string `json:"revealedAttributes"`
	Nonce              string   `json:"nonce,omitempty"`
}

// CredentialService interface for credential operations
type CredentialService interface {
	SetIssuerKeyPair(issuerDID string, keyPair *bbs.KeyPair)
	IssueCredential(issuerDID string, subjectDID string, claims []Claim) (*VerifiableCredential, error)
	VerifyCredential(vc *VerifiableCredential) error
	CreatePresentation(holderDID string, credentials []*VerifiableCredential, disclosureRequests []SelectiveDisclosureRequest) (*VerifiablePresentation, error)
	VerifyPresentation(vp *VerifiablePresentation) error
}

// CredentialRepository interface for credential storage
type CredentialRepository interface {
	Store(vc *VerifiableCredential) error
	Retrieve(id string) (*VerifiableCredential, error)
	List(holderDID string) ([]*VerifiableCredential, error)
}

// PresentationRepository interface for presentation storage
type PresentationRepository interface {
	Store(vp *VerifiablePresentation) error
	Retrieve(id string) (*VerifiablePresentation, error)
	List(holderDID string) ([]*VerifiablePresentation, error)
}
