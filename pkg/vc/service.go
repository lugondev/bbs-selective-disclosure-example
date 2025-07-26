package vc

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
)

// ServiceImpl implements CredentialService interface
type ServiceImpl struct {
	bbsService bbs.BBSService
	credRepo   CredentialRepository
	presRepo   PresentationRepository
	keyStore   map[string]*bbs.KeyPair // DID -> KeyPair mapping
}

// NewService creates a new credential service
func NewService(bbsService bbs.BBSService, credRepo CredentialRepository, presRepo PresentationRepository) CredentialService {
	return &ServiceImpl{
		bbsService: bbsService,
		credRepo:   credRepo,
		presRepo:   presRepo,
		keyStore:   make(map[string]*bbs.KeyPair),
	}
}

// SetIssuerKeyPair sets the BBS+ key pair for an issuer DID
func (s *ServiceImpl) SetIssuerKeyPair(issuerDID string, keyPair *bbs.KeyPair) {
	s.keyStore[issuerDID] = keyPair
}

// IssueCredential creates and signs a new verifiable credential
func (s *ServiceImpl) IssueCredential(issuerDID string, subjectDID string, claims []Claim) (*VerifiableCredential, error) {
	keyPair, exists := s.keyStore[issuerDID]
	if !exists {
		return nil, fmt.Errorf("no key pair found for issuer DID: %s", issuerDID)
	}

	// Create credential subject
	credentialSubject := make(map[string]interface{})
	credentialSubject["id"] = subjectDID

	// Convert claims to messages for BBS+ signing
	var messages [][]byte
	var claimKeys []string

	for _, claim := range claims {
		credentialSubject[claim.Key] = claim.Value
		claimKeys = append(claimKeys, claim.Key)

		// Convert claim value to bytes
		valueBytes, err := json.Marshal(claim.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal claim value: %w", err)
		}
		messages = append(messages, valueBytes)
	}

	// Create the credential
	now := time.Now()
	credential := &VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://w3id.org/security/bbs/v1",
		},
		ID:                uuid.New().String(),
		Type:              []string{"VerifiableCredential"},
		Issuer:            issuerDID,
		IssuanceDate:      now,
		CredentialSubject: credentialSubject,
	}

	// Sign with BBS+
	signature, err := s.bbsService.Sign(keyPair.PrivateKey, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to sign credential: %w", err)
	}

	// Create proof
	credential.Proof = &Proof{
		Type:               "BbsBlsSignature2020",
		Created:            now,
		VerificationMethod: issuerDID + "#bbs-key-1",
		ProofPurpose:       "assertionMethod",
		ProofValue:         bbs.EncodeProof(&bbs.Proof{ProofValue: signature.Value}),
	}

	// Store metadata for later proof creation
	credential.Proof.RevealedAttributes = make([]int, len(claims))
	for i := range claims {
		credential.Proof.RevealedAttributes[i] = i
	}

	return credential, nil
}

// VerifyCredential verifies a verifiable credential
func (s *ServiceImpl) VerifyCredential(vc *VerifiableCredential) error {
	if vc == nil {
		return fmt.Errorf("credential is nil")
	}

	if vc.Proof == nil {
		return fmt.Errorf("credential has no proof")
	}

	// For demonstration, we'll skip actual BBS+ verification
	// In production, you would:
	// 1. Resolve issuer DID to get public key
	// 2. Reconstruct messages from credential subject
	// 3. Verify BBS+ signature

	return nil
}

// CreatePresentation creates a verifiable presentation with selective disclosure
func (s *ServiceImpl) CreatePresentation(holderDID string, credentials []*VerifiableCredential, disclosureRequests []SelectiveDisclosureRequest) (*VerifiablePresentation, error) {
	if len(credentials) != len(disclosureRequests) {
		return nil, fmt.Errorf("mismatch between credentials and disclosure requests")
	}

	var presentedCredentials []interface{}

	for i, credential := range credentials {
		request := disclosureRequests[i]

		// Create selective disclosure proof
		derivedCredential, err := s.createSelectiveDisclosureCredential(credential, request)
		if err != nil {
			return nil, fmt.Errorf("failed to create selective disclosure: %w", err)
		}

		presentedCredentials = append(presentedCredentials, derivedCredential)
	}

	// Create presentation
	presentation := &VerifiablePresentation{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://w3id.org/security/bbs/v1",
		},
		ID:                   uuid.New().String(),
		Type:                 []string{"VerifiablePresentation"},
		Holder:               holderDID,
		VerifiableCredential: presentedCredentials,
	}

	// Add presentation proof (simplified)
	now := time.Now()
	presentation.Proof = &Proof{
		Type:               "BbsBlsSignatureProof2020",
		Created:            now,
		VerificationMethod: holderDID + "#key-1",
		ProofPurpose:       "authentication",
	}

	return presentation, nil
}

// createSelectiveDisclosureCredential creates a derived credential with only revealed attributes
func (s *ServiceImpl) createSelectiveDisclosureCredential(credential *VerifiableCredential, request SelectiveDisclosureRequest) (map[string]interface{}, error) {
	// Create derived credential with only revealed attributes
	derivedCredential := map[string]interface{}{
		"@context":          credential.Context,
		"id":                credential.ID,
		"type":              credential.Type,
		"issuer":            credential.Issuer,
		"issuanceDate":      credential.IssuanceDate,
		"credentialSubject": make(map[string]interface{}),
	}

	// Include subject ID
	if subjectID, ok := credential.CredentialSubject["id"]; ok {
		derivedCredential["credentialSubject"].(map[string]interface{})["id"] = subjectID
	}

	// Include only revealed attributes
	for _, attr := range request.RevealedAttributes {
		if value, exists := credential.CredentialSubject[attr]; exists {
			derivedCredential["credentialSubject"].(map[string]interface{})[attr] = value
		}
	}

	// Use provided nonce or generate one if not provided
	var nonceStr string
	if request.Nonce != "" {
		nonceStr = request.Nonce
	} else {
		// Generate nonce for proof
		nonce := make([]byte, 32)
		if _, err := rand.Read(nonce); err != nil {
			return nil, fmt.Errorf("failed to generate nonce: %w", err)
		}
		nonceStr = fmt.Sprintf("%x", nonce)
	}

	// Create selective disclosure proof
	// In a real implementation, this would use the original BBS+ signature
	// to create a proof for only the revealed attributes
	derivedCredential["proof"] = map[string]interface{}{
		"type":               "BbsBlsSignatureProof2020",
		"created":            time.Now(),
		"verificationMethod": credential.Proof.VerificationMethod,
		"proofPurpose":       "assertionMethod",
		"proofValue":         "derived-proof-placeholder",
		"nonce":              nonceStr,
		"revealedAttributes": request.RevealedAttributes,
	}

	return derivedCredential, nil
}

// VerifyPresentation verifies a verifiable presentation
func (s *ServiceImpl) VerifyPresentation(vp *VerifiablePresentation) error {
	if vp == nil {
		return fmt.Errorf("presentation is nil")
	}

	if vp.Proof == nil {
		return fmt.Errorf("presentation has no proof")
	}

	// Verify each credential in the presentation
	for _, credInterface := range vp.VerifiableCredential {
		// In a real implementation, you would:
		// 1. Parse the derived credential
		// 2. Verify the selective disclosure proof
		// 3. Ensure only requested attributes are revealed
		_ = credInterface
	}

	return nil
}

// InMemoryCredentialRepository implements CredentialRepository interface
type InMemoryCredentialRepository struct {
	credentials map[string]*VerifiableCredential
}

// NewInMemoryCredentialRepository creates a new in-memory credential repository
func NewInMemoryCredentialRepository() CredentialRepository {
	return &InMemoryCredentialRepository{
		credentials: make(map[string]*VerifiableCredential),
	}
}

// Store stores a verifiable credential
func (r *InMemoryCredentialRepository) Store(vc *VerifiableCredential) error {
	if vc == nil {
		return fmt.Errorf("credential is nil")
	}
	r.credentials[vc.ID] = vc
	return nil
}

// Retrieve retrieves a verifiable credential by ID
func (r *InMemoryCredentialRepository) Retrieve(id string) (*VerifiableCredential, error) {
	vc, exists := r.credentials[id]
	if !exists {
		return nil, fmt.Errorf("credential not found: %s", id)
	}
	return vc, nil
}

// List lists all credentials for a holder DID
func (r *InMemoryCredentialRepository) List(holderDID string) ([]*VerifiableCredential, error) {
	var credentials []*VerifiableCredential
	for _, vc := range r.credentials {
		if subjectID, ok := vc.CredentialSubject["id"].(string); ok && subjectID == holderDID {
			credentials = append(credentials, vc)
		}
	}
	return credentials, nil
}

// InMemoryPresentationRepository implements PresentationRepository interface
type InMemoryPresentationRepository struct {
	presentations map[string]*VerifiablePresentation
}

// NewInMemoryPresentationRepository creates a new in-memory presentation repository
func NewInMemoryPresentationRepository() PresentationRepository {
	return &InMemoryPresentationRepository{
		presentations: make(map[string]*VerifiablePresentation),
	}
}

// Store stores a verifiable presentation
func (r *InMemoryPresentationRepository) Store(vp *VerifiablePresentation) error {
	if vp == nil {
		return fmt.Errorf("presentation is nil")
	}
	r.presentations[vp.ID] = vp
	return nil
}

// Retrieve retrieves a verifiable presentation by ID
func (r *InMemoryPresentationRepository) Retrieve(id string) (*VerifiablePresentation, error) {
	vp, exists := r.presentations[id]
	if !exists {
		return nil, fmt.Errorf("presentation not found: %s", id)
	}
	return vp, nil
}

// List lists all presentations for a holder DID
func (r *InMemoryPresentationRepository) List(holderDID string) ([]*VerifiablePresentation, error) {
	var presentations []*VerifiablePresentation
	for _, vp := range r.presentations {
		if vp.Holder == holderDID {
			presentations = append(presentations, vp)
		}
	}
	return presentations, nil
}
