package did

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

// ServiceImpl implements DIDService interface
type ServiceImpl struct {
	repository DIDRepository
}

// NewService creates a new DID service
func NewService(repo DIDRepository) DIDService {
	return &ServiceImpl{
		repository: repo,
	}
}

// GenerateDID generates a new DID with key pair
func (s *ServiceImpl) GenerateDID(method string) (*DID, *KeyPair, error) {
	// Generate Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Create identifier from public key
	identifier := base58.Encode(publicKey)

	did := &DID{
		Method:     method,
		Identifier: identifier,
	}

	keyPair := &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		KeyID:      did.String() + "#key-1",
	}

	return did, keyPair, nil
}

// CreateDIDDocument creates a DID document for the given DID and key pair
func (s *ServiceImpl) CreateDIDDocument(did *DID, keyPair *KeyPair) (*DIDDocument, error) {
	now := time.Now()

	verificationMethod := VerificationMethod{
		ID:                 keyPair.KeyID,
		Type:               "Ed25519VerificationKey2020",
		Controller:         did.String(),
		PublicKeyMultibase: "z" + base58.Encode(keyPair.PublicKey),
	}

	doc := &DIDDocument{
		Context: []string{
			"https://www.w3.org/ns/did/v1",
			"https://w3id.org/security/suites/ed25519-2020/v1",
		},
		ID:                 did.String(),
		VerificationMethod: []VerificationMethod{verificationMethod},
		Authentication:     []string{keyPair.KeyID},
		AssertionMethod:    []string{keyPair.KeyID},
		Created:            now,
		Updated:            now,
	}

	return doc, nil
}

// ResolveDID resolves a DID to its DID Document
func (s *ServiceImpl) ResolveDID(didString string) (*DIDDocument, error) {
	return s.repository.Resolve(didString)
}

// VerifyDIDDocument verifies the integrity of a DID Document
func (s *ServiceImpl) VerifyDIDDocument(doc *DIDDocument) error {
	if doc == nil {
		return fmt.Errorf("DID document is nil")
	}

	if doc.ID == "" {
		return fmt.Errorf("DID document ID is empty")
	}

	if len(doc.VerificationMethod) == 0 {
		return fmt.Errorf("DID document must have at least one verification method")
	}

	// Verify that authentication methods reference valid verification methods
	for _, authMethod := range doc.Authentication {
		found := false
		for _, vm := range doc.VerificationMethod {
			if vm.ID == authMethod {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("authentication method %s not found in verification methods", authMethod)
		}
	}

	return nil
}

// InMemoryRepository implements DIDRepository interface for testing
type InMemoryRepository struct {
	documents map[string]*DIDDocument
}

// NewInMemoryRepository creates a new in-memory DID repository
func NewInMemoryRepository() DIDRepository {
	return &InMemoryRepository{
		documents: make(map[string]*DIDDocument),
	}
}

// Create stores a DID document
func (r *InMemoryRepository) Create(doc *DIDDocument) error {
	if doc == nil {
		return fmt.Errorf("DID document is nil")
	}
	r.documents[doc.ID] = doc
	return nil
}

// Resolve retrieves a DID document by DID
func (r *InMemoryRepository) Resolve(did string) (*DIDDocument, error) {
	doc, exists := r.documents[did]
	if !exists {
		return nil, fmt.Errorf("DID document not found: %s", did)
	}
	return doc, nil
}

// Update updates an existing DID document
func (r *InMemoryRepository) Update(did string, doc *DIDDocument) error {
	if _, exists := r.documents[did]; !exists {
		return fmt.Errorf("DID document not found: %s", did)
	}
	doc.Updated = time.Now()
	r.documents[did] = doc
	return nil
}

// Deactivate removes a DID document
func (r *InMemoryRepository) Deactivate(did string) error {
	delete(r.documents, did)
	return nil
}
