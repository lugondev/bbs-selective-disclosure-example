package did

import (
	"crypto/ed25519"
	"time"
)

// DID represents a Decentralized Identifier
type DID struct {
	Method     string `json:"method"`
	Identifier string `json:"identifier"`
}

// String returns the full DID string
func (d DID) String() string {
	return "did:" + d.Method + ":" + d.Identifier
}

// DIDDocument represents a DID Document structure
type DIDDocument struct {
	Context            []string             `json:"@context"`
	ID                 string               `json:"id"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
	Authentication     []string             `json:"authentication"`
	AssertionMethod    []string             `json:"assertionMethod"`
	KeyAgreement       []string             `json:"keyAgreement,omitempty"`
	Service            []Service            `json:"service,omitempty"`
	Created            time.Time            `json:"created"`
	Updated            time.Time            `json:"updated"`
}

// VerificationMethod represents a verification method in DID Document
type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

// Service represents a service endpoint in DID Document
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// KeyPair represents a cryptographic key pair
type KeyPair struct {
	PublicKey  ed25519.PublicKey  `json:"publicKey"`
	PrivateKey ed25519.PrivateKey `json:"privateKey"`
	KeyID      string             `json:"keyId"`
}

// DIDRepository interface for DID operations
type DIDRepository interface {
	Create(doc *DIDDocument) error
	Resolve(did string) (*DIDDocument, error)
	Update(did string, doc *DIDDocument) error
	Deactivate(did string) error
}

// DIDService interface for DID business logic
type DIDService interface {
	GenerateDID(method string) (*DID, *KeyPair, error)
	CreateDIDDocument(did *DID, keyPair *KeyPair) (*DIDDocument, error)
	ResolveDID(didString string) (*DIDDocument, error)
	VerifyDIDDocument(doc *DIDDocument) error
}
