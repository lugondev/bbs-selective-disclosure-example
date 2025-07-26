package verifier

import (
	"fmt"

	"github.com/lugon/bbs-selective-disclosure-example/pkg/did"
	"github.com/lugon/bbs-selective-disclosure-example/pkg/vc"
)

// UseCase represents the verifier use case
type UseCase struct {
	didService did.DIDService
	vcService  vc.CredentialService
	presRepo   vc.PresentationRepository
}

// NewUseCase creates a new verifier use case
func NewUseCase(didService did.DIDService, vcService vc.CredentialService, presRepo vc.PresentationRepository) *UseCase {
	return &UseCase{
		didService: didService,
		vcService:  vcService,
		presRepo:   presRepo,
	}
}

// VerifierSetup represents the setup process for a verifier
type VerifierSetup struct {
	DID     *did.DID
	DIDDoc  *did.DIDDocument
	KeyPair *did.KeyPair
}

// SetupVerifier sets up a new verifier with DID
func (uc *UseCase) SetupVerifier(method string) (*VerifierSetup, error) {
	// Generate DID and key pair
	verifierDID, keyPair, err := uc.didService.GenerateDID(method)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DID: %w", err)
	}

	// Create DID document
	didDoc, err := uc.didService.CreateDIDDocument(verifierDID, keyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to create DID document: %w", err)
	}

	return &VerifierSetup{
		DID:     verifierDID,
		DIDDoc:  didDoc,
		KeyPair: keyPair,
	}, nil
}

// VerificationRequest represents a verification request
type VerificationRequest struct {
	Presentation      *vc.VerifiablePresentation
	RequiredClaims    []string
	TrustedIssuers    []string
	VerificationNonce string
}

// VerificationResult represents the result of verification
type VerificationResult struct {
	Valid           bool                   `json:"valid"`
	Errors          []string               `json:"errors,omitempty"`
	RevealedClaims  map[string]interface{} `json:"revealedClaims,omitempty"`
	HolderDID       string                 `json:"holderDid"`
	IssuerDIDs      []string               `json:"issuerDids"`
	CredentialTypes []string               `json:"credentialTypes"`
}

// VerifyPresentation verifies a verifiable presentation
func (uc *UseCase) VerifyPresentation(req VerificationRequest) (*VerificationResult, error) {
	result := &VerificationResult{
		Valid:           true,
		Errors:          []string{},
		RevealedClaims:  make(map[string]interface{}),
		HolderDID:       req.Presentation.Holder,
		IssuerDIDs:      []string{},
		CredentialTypes: []string{},
	}

	// Verify presentation structure
	if err := uc.vcService.VerifyPresentation(req.Presentation); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("presentation verification failed: %v", err))
		return result, nil
	}

	// Verify each credential in the presentation
	for i, credInterface := range req.Presentation.VerifiableCredential {
		credMap, ok := credInterface.(map[string]interface{})
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("credential %d: invalid format", i))
			continue
		}

		// Extract issuer
		issuer, ok := credMap["issuer"].(string)
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("credential %d: missing or invalid issuer", i))
			continue
		}

		result.IssuerDIDs = append(result.IssuerDIDs, issuer)

		// Check if issuer is trusted
		if len(req.TrustedIssuers) > 0 {
			trusted := false
			for _, trustedIssuer := range req.TrustedIssuers {
				if issuer == trustedIssuer {
					trusted = true
					break
				}
			}
			if !trusted {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("credential %d: issuer %s is not trusted", i, issuer))
				continue
			}
		}

		// Extract credential types
		if types, ok := credMap["type"].([]interface{}); ok {
			for _, t := range types {
				if typeStr, ok := t.(string); ok {
					result.CredentialTypes = append(result.CredentialTypes, typeStr)
				}
			}
		}

		// Extract revealed claims from credential subject
		if credentialSubject, ok := credMap["credentialSubject"].(map[string]interface{}); ok {
			for key, value := range credentialSubject {
				if key != "id" { // Skip subject ID
					result.RevealedClaims[key] = value
				}
			}
		}

		// Verify selective disclosure proof
		if err := uc.verifySelectiveDisclosureProof(credMap, req.VerificationNonce); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("credential %d: selective disclosure verification failed: %v", i, err))
		}
	}

	// Check if all required claims are present
	for _, requiredClaim := range req.RequiredClaims {
		if _, exists := result.RevealedClaims[requiredClaim]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("required claim '%s' is missing", requiredClaim))
		}
	}

	// Store verification result
	if result.Valid {
		if err := uc.presRepo.Store(req.Presentation); err != nil {
			// Log error but don't fail verification
			result.Errors = append(result.Errors, fmt.Sprintf("failed to store presentation: %v", err))
		}
	}

	return result, nil
}

// verifySelectiveDisclosureProof verifies the selective disclosure proof
func (uc *UseCase) verifySelectiveDisclosureProof(credMap map[string]interface{}, nonce string) error {
	proof, ok := credMap["proof"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing or invalid proof")
	}

	proofType, ok := proof["type"].(string)
	if !ok || proofType != "BbsBlsSignatureProof2020" {
		return fmt.Errorf("invalid proof type: expected BbsBlsSignatureProof2020, got %v", proofType)
	}

	// Verify nonce if provided
	if nonce != "" {
		proofNonce, ok := proof["nonce"].(string)
		if !ok || proofNonce != nonce {
			return fmt.Errorf("nonce mismatch: expected %s, got %v", nonce, proofNonce)
		}
	}

	// In a real implementation, you would:
	// 1. Resolve the issuer DID to get the public key
	// 2. Verify the BBS+ proof using the public key
	// 3. Ensure only the claimed attributes are revealed

	return nil
}

// CreateVerificationRequest creates a verification request for specific claims
type CreateVerificationRequestParams struct {
	RequiredClaims    []string
	TrustedIssuers    []string
	VerificationNonce string
}

// CreateVerificationRequest creates a verification request
func (uc *UseCase) CreateVerificationRequest(params CreateVerificationRequestParams) (*CreateVerificationRequestParams, error) {
	// Generate a nonce if not provided
	if params.VerificationNonce == "" {
		// In a real implementation, generate a cryptographically secure nonce
		params.VerificationNonce = "verification-nonce-" + fmt.Sprintf("%d", len(params.RequiredClaims))
	}

	return &params, nil
}

// ListVerifiedPresentations lists all verified presentations
func (uc *UseCase) ListVerifiedPresentations(verifierDID string) ([]*vc.VerifiablePresentation, error) {
	// In this simplified implementation, we'll return all presentations
	// In a real implementation, you might filter by verifier
	presentations, err := uc.presRepo.List("")
	if err != nil {
		return nil, fmt.Errorf("failed to list presentations: %w", err)
	}

	return presentations, nil
}
