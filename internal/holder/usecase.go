package holder

import (
	"fmt"

	"github.com/lugondev/bbs-selective-disclosure-example/pkg/did"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"
)

// UseCase represents the holder use case
type UseCase struct {
	didService did.DIDService
	vcService  vc.CredentialService
	credRepo   vc.CredentialRepository
}

// NewUseCase creates a new holder use case
func NewUseCase(didService did.DIDService, vcService vc.CredentialService, credRepo vc.CredentialRepository) *UseCase {
	return &UseCase{
		didService: didService,
		vcService:  vcService,
		credRepo:   credRepo,
	}
}

// HolderSetup represents the setup process for a holder
type HolderSetup struct {
	DID     *did.DID
	DIDDoc  *did.DIDDocument
	KeyPair *did.KeyPair
}

// SetupHolder sets up a new holder with DID
func (uc *UseCase) SetupHolder(method string) (*HolderSetup, error) {
	// Generate DID and key pair
	holderDID, keyPair, err := uc.didService.GenerateDID(method)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DID: %w", err)
	}

	// Create DID document
	didDoc, err := uc.didService.CreateDIDDocument(holderDID, keyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to create DID document: %w", err)
	}

	return &HolderSetup{
		DID:     holderDID,
		DIDDoc:  didDoc,
		KeyPair: keyPair,
	}, nil
}

// StoreCredential stores a received credential
func (uc *UseCase) StoreCredential(credential *vc.VerifiableCredential) error {
	if credential == nil {
		return fmt.Errorf("credential is nil")
	}

	// Verify credential before storing
	if err := uc.vcService.VerifyCredential(credential); err != nil {
		return fmt.Errorf("credential verification failed: %w", err)
	}

	// Store credential
	if err := uc.credRepo.Store(credential); err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}

	return nil
}

// ListCredentials lists all credentials for a holder
func (uc *UseCase) ListCredentials(holderDID string) ([]*vc.VerifiableCredential, error) {
	credentials, err := uc.credRepo.List(holderDID)
	if err != nil {
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}

	return credentials, nil
}

// PresentationRequest represents a presentation request
type PresentationRequest struct {
	HolderDID           string
	CredentialIDs       []string
	SelectiveDisclosure []vc.SelectiveDisclosureRequest
	Nonce               string
}

// CreatePresentation creates a verifiable presentation with selective disclosure
func (uc *UseCase) CreatePresentation(req PresentationRequest) (*vc.VerifiablePresentation, error) {
	if req.HolderDID == "" {
		return nil, fmt.Errorf("holder DID is required")
	}

	if len(req.CredentialIDs) == 0 {
		return nil, fmt.Errorf("at least one credential ID is required")
	}

	if len(req.CredentialIDs) != len(req.SelectiveDisclosure) {
		return nil, fmt.Errorf("mismatch between credential IDs and selective disclosure requests")
	}

	// Retrieve credentials
	var credentials []*vc.VerifiableCredential
	for _, credID := range req.CredentialIDs {
		credential, err := uc.credRepo.Retrieve(credID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve credential %s: %w", credID, err)
		}

		// Verify holder owns the credential
		if subjectID, ok := credential.CredentialSubject["id"].(string); !ok || subjectID != req.HolderDID {
			return nil, fmt.Errorf("credential %s does not belong to holder %s", credID, req.HolderDID)
		}

		credentials = append(credentials, credential)
	}

	// Set nonce for each selective disclosure request if provided
	disclosureRequests := make([]vc.SelectiveDisclosureRequest, len(req.SelectiveDisclosure))
	for i, sd := range req.SelectiveDisclosure {
		disclosureRequests[i] = sd
		if req.Nonce != "" {
			disclosureRequests[i].Nonce = req.Nonce
		}
	}

	// Create presentation
	presentation, err := uc.vcService.CreatePresentation(req.HolderDID, credentials, disclosureRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to create presentation: %w", err)
	}

	return presentation, nil
}

// GetCredential retrieves a specific credential
func (uc *UseCase) GetCredential(credentialID string) (*vc.VerifiableCredential, error) {
	credential, err := uc.credRepo.Retrieve(credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credential: %w", err)
	}

	return credential, nil
}
