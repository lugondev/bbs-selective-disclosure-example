package bbs

import (
	"fmt"

	bls12381 "github.com/kilic/bls12-381"
)

// SimpleService adapter wraps the simple BBS implementation
type SimpleService struct {
	config   *Config
	provider Provider
	version  string
}

// newSimpleService creates a new simple BBS service
func newSimpleService(config *Config) BBSInterface {
	return &SimpleService{
		config:   config,
		provider: ProviderSimple,
		version:  "1.0.0-simple",
	}
}

// GenerateKeyPair generates a simple key pair
func (s *SimpleService) GenerateKeyPair() (*KeyPair, error) {
	// This is a simplified implementation for demo purposes
	// In production, this should use secure random generation
	privateKey := make([]byte, 32)
	publicKey := make([]byte, 32)

	// Simple demo key generation (NOT secure)
	for i := range privateKey {
		privateKey[i] = byte(i + 1)
	}
	for i := range publicKey {
		publicKey[i] = byte((i + 1) * 2)
	}

	return &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// Sign creates a simple signature
func (s *SimpleService) Sign(privateKey []byte, messages [][]byte) (*Signature, error) {
	if len(privateKey) == 0 {
		return nil, fmt.Errorf("private key cannot be empty")
	}

	// Simple signature for demo (NOT secure)
	signature := &Signature{
		A: make([]byte, 32),
		E: make([]byte, 32),
		S: make([]byte, 32),
	}

	// Fill with demo data
	for i := range signature.A {
		signature.A[i] = byte(i)
		signature.E[i] = byte(i + 1)
		signature.S[i] = byte(i + 2)
	}

	return signature, nil
}

// Verify verifies a simple signature
func (s *SimpleService) Verify(publicKey []byte, signature *Signature, messages [][]byte) error {
	if len(publicKey) == 0 {
		return fmt.Errorf("public key cannot be empty")
	}

	if signature == nil {
		return fmt.Errorf("signature cannot be nil")
	}

	// Simple verification (always passes for demo)
	return nil
}

// CreateProof creates a simple proof
func (s *SimpleService) CreateProof(signature *Signature, publicKey []byte, messages [][]byte, revealedIndices []int, nonce []byte) (*Proof, error) {
	if signature == nil {
		return nil, fmt.Errorf("signature cannot be nil")
	}

	// Simple proof for demo
	proof := &Proof{
		A_prime:            make([]byte, 32),
		A_bar:              make([]byte, 32),
		C:                  make([]byte, 32),
		R2:                 make([]byte, 32),
		R3:                 make([]byte, 32),
		HiddenResponses:    [][]byte{},
		RevealedAttributes: revealedIndices,
		Nonce:              nonce,
	}

	// Fill with demo data
	for i := range proof.A_prime {
		proof.A_prime[i] = byte(i)
		proof.A_bar[i] = byte(i + 1)
		proof.C[i] = byte(i + 2)
		proof.R2[i] = byte(i + 3)
		proof.R3[i] = byte(i + 4)
	}

	return proof, nil
}

// VerifyProof verifies a simple proof
func (s *SimpleService) VerifyProof(publicKey []byte, proof *Proof, revealedMessages [][]byte, nonce []byte) error {
	if proof == nil {
		return fmt.Errorf("proof cannot be nil")
	}

	// Simple verification (always passes for demo)
	return nil
}

// ValidateKeyPair validates a key pair
func (s *SimpleService) ValidateKeyPair(keyPair *KeyPair) error {
	if keyPair == nil {
		return fmt.Errorf("key pair cannot be nil")
	}

	if len(keyPair.PrivateKey) == 0 {
		return fmt.Errorf("private key cannot be empty")
	}

	if len(keyPair.PublicKey) == 0 {
		return fmt.Errorf("public key cannot be empty")
	}

	return nil
}

// GetMessageCount returns message count
func (s *SimpleService) GetMessageCount(signature *Signature, publicKey []byte) (int, error) {
	return 0, fmt.Errorf("message count not available in simple implementation")
}

// ConstantTimeVerify performs verification
func (s *SimpleService) ConstantTimeVerify(publicKey []byte, signature *Signature, messages [][]byte) error {
	// Simple implementation doesn't have constant time guarantees
	return s.Verify(publicKey, signature, messages)
}

// SecureErase clears data
func (s *SimpleService) SecureErase(data []byte) {
	// Simple secure erase
	for i := range data {
		data[i] = 0
	}
}

// GetProvider returns provider type
func (s *SimpleService) GetProvider() Provider {
	return s.provider
}

// GetVersion returns version
func (s *SimpleService) GetVersion() string {
	return s.version
}

// IsProductionReady returns production readiness
func (s *SimpleService) IsProductionReady() bool {
	return false // Simple implementation is not production ready
}

// ProductionServiceAdapter adapts the existing ProductionService to the new interface
type ProductionServiceAdapter struct {
	service *ProductionService
	config  *Config
	version string
}

// newProductionService creates a new production BBS service adapter
func newProductionService(config *Config) BBSInterface {
	return &ProductionServiceAdapter{
		service: &ProductionService{
			g1:     bls12381.NewG1(),
			g2:     bls12381.NewG2(),
			gt:     bls12381.NewGT(),
			engine: bls12381.NewEngine(),
		},
		config:  config,
		version: "1.0.0-production",
	}
}

// GenerateKeyPair generates a production key pair
func (a *ProductionServiceAdapter) GenerateKeyPair() (*KeyPair, error) {
	return a.service.GenerateKeyPair()
}

// Sign creates a production signature
func (a *ProductionServiceAdapter) Sign(privateKey []byte, messages [][]byte) (*Signature, error) {
	return a.service.Sign(privateKey, messages)
}

// Verify verifies a production signature
func (a *ProductionServiceAdapter) Verify(publicKey []byte, signature *Signature, messages [][]byte) error {
	return a.service.Verify(publicKey, signature, messages)
}

// CreateProof creates a production proof
func (a *ProductionServiceAdapter) CreateProof(signature *Signature, publicKey []byte, messages [][]byte, revealedIndices []int, nonce []byte) (*Proof, error) {
	return a.service.CreateProof(signature, publicKey, messages, revealedIndices, nonce)
}

// VerifyProof verifies a production proof
func (a *ProductionServiceAdapter) VerifyProof(publicKey []byte, proof *Proof, revealedMessages [][]byte, nonce []byte) error {
	return a.service.VerifyProof(publicKey, proof, revealedMessages, nonce)
}

// ValidateKeyPair validates a key pair
func (a *ProductionServiceAdapter) ValidateKeyPair(keyPair *KeyPair) error {
	return a.service.ValidateKeyPair(keyPair)
}

// GetMessageCount returns message count
func (a *ProductionServiceAdapter) GetMessageCount(signature *Signature, publicKey []byte) (int, error) {
	return a.service.GetMessageCount(signature, publicKey)
}

// ConstantTimeVerify performs constant time verification
func (a *ProductionServiceAdapter) ConstantTimeVerify(publicKey []byte, signature *Signature, messages [][]byte) error {
	return a.service.ConstantTimeVerify(publicKey, signature, messages)
}

// SecureErase securely erases data
func (a *ProductionServiceAdapter) SecureErase(data []byte) {
	a.service.SecureErase(data)
}

// GetProvider returns provider type
func (a *ProductionServiceAdapter) GetProvider() Provider {
	return ProviderProduction
}

// GetVersion returns version
func (a *ProductionServiceAdapter) GetVersion() string {
	return a.version
}

// IsProductionReady returns production readiness
func (a *ProductionServiceAdapter) IsProductionReady() bool {
	return true
}
