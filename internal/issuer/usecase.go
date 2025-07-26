package issuer

import (
	"fmt"

	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/did"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"
)

// UseCase represents the issuer use case
type UseCase struct {
	didService did.DIDService
	vcService  vc.CredentialService
	bbsService bbs.BBSService
}

// NewUseCase creates a new issuer use case
func NewUseCase(didService did.DIDService, vcService vc.CredentialService, bbsService bbs.BBSService) *UseCase {
	return &UseCase{
		didService: didService,
		vcService:  vcService,
		bbsService: bbsService,
	}
}

// IssuerSetup represents the setup process for an issuer
type IssuerSetup struct {
	DID        *did.DID
	DIDDoc     *did.DIDDocument
	KeyPair    *did.KeyPair
	BBSKeyPair *bbs.KeyPair
}

// SetupIssuer sets up a new issuer with DID and keys
func (uc *UseCase) SetupIssuer(method string) (*IssuerSetup, error) {
	// Generate DID and key pair
	issuerDID, keyPair, err := uc.didService.GenerateDID(method)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DID: %w", err)
	}

	// Create DID document
	didDoc, err := uc.didService.CreateDIDDocument(issuerDID, keyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to create DID document: %w", err)
	}

	// Generate BBS+ key pair for signing credentials
	bbsKeyPair, err := uc.bbsService.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate BBS+ key pair: %w", err)
	}

	// Set up the issuer in the VC service
	uc.vcService.SetIssuerKeyPair(issuerDID.String(), bbsKeyPair)

	return &IssuerSetup{
		DID:        issuerDID,
		DIDDoc:     didDoc,
		KeyPair:    keyPair,
		BBSKeyPair: bbsKeyPair,
	}, nil
}

// IssueCredentialRequest represents a credential issuance request
type IssueCredentialRequest struct {
	IssuerDID  string
	SubjectDID string
	Claims     []vc.Claim
}

// IssueCredential issues a new verifiable credential
func (uc *UseCase) IssueCredential(req IssueCredentialRequest) (*vc.VerifiableCredential, error) {
	if req.IssuerDID == "" {
		return nil, fmt.Errorf("issuer DID is required")
	}

	if req.SubjectDID == "" {
		return nil, fmt.Errorf("subject DID is required")
	}

	if len(req.Claims) == 0 {
		return nil, fmt.Errorf("at least one claim is required")
	}

	// Issue the credential
	credential, err := uc.vcService.IssueCredential(req.IssuerDID, req.SubjectDID, req.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to issue credential: %w", err)
	}

	return credential, nil
}

// VerifyCredential verifies a verifiable credential
func (uc *UseCase) VerifyCredential(credential *vc.VerifiableCredential) error {
	return uc.vcService.VerifyCredential(credential)
}
