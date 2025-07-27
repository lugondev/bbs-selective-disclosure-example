package bbs

import (
	"fmt"
	"strings"
	"time"
)

// Provider represents the BBS implementation provider type
type Provider string

const (
	// ProviderSimple uses the current simple implementation
	ProviderSimple Provider = "simple"
	// ProviderAries uses Hyperledger Aries Framework Go
	ProviderAries Provider = "aries"
	// ProviderProduction uses the production BLS12-381 implementation
	ProviderProduction Provider = "production"
)

// String returns the string representation of the Provider
func (p Provider) String() string {
	return string(p)
}

// ParseProvider parses a string into a Provider
func ParseProvider(s string) (Provider, error) {
	provider := Provider(strings.ToLower(s))
	switch provider {
	case ProviderSimple, ProviderAries, ProviderProduction:
		return provider, nil
	default:
		return "", fmt.Errorf("unknown provider: %s", s)
	}
}

// BBSInterface defines the core BBS+ operations interface
type BBSInterface interface {
	// Core BBS+ operations
	GenerateKeyPair() (*KeyPair, error)
	Sign(privateKey []byte, messages [][]byte) (*Signature, error)
	Verify(publicKey []byte, signature *Signature, messages [][]byte) error
	CreateProof(signature *Signature, publicKey []byte, messages [][]byte, revealedIndices []int, nonce []byte) (*Proof, error)
	VerifyProof(publicKey []byte, proof *Proof, revealedMessages [][]byte, nonce []byte) error

	// Validation and utility methods
	ValidateKeyPair(keyPair *KeyPair) error
	GetMessageCount(signature *Signature, publicKey []byte) (int, error)

	// Security features
	ConstantTimeVerify(publicKey []byte, signature *Signature, messages [][]byte) error
	SecureErase(data []byte)

	// Metadata
	GetProvider() Provider
	GetVersion() string
	IsProductionReady() bool
}

// Config holds configuration for BBS service initialization
type Config struct {
	Provider Provider `json:"provider"`
	// Performance settings
	EnableLogging    bool          `json:"enable_logging"`
	OperationTimeout time.Duration `json:"operation_timeout"`

	// Security settings
	ConstantTimeOps bool `json:"constant_time_ops"`
	SecureMemory    bool `json:"secure_memory"`

	// Aries-specific settings
	AriesConfig *AriesConfig `json:"aries_config,omitempty"`
}

// AriesConfig holds Aries Framework specific configuration
type AriesConfig struct {
	// Aries framework settings
	KMSType         string `json:"kms_type"`         // "local" or "remote"
	StorageProvider string `json:"storage_provider"` // "mem" or "leveldb"
	CryptoSuite     string `json:"crypto_suite"`     // BBS crypto suite

	// Remote KMS settings (if applicable)
	RemoteKMSURL string `json:"remote_kms_url,omitempty"`
	AuthToken    string `json:"auth_token,omitempty"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Provider:         ProviderProduction,
		EnableLogging:    true,
		OperationTimeout: 30 * time.Second,
		ConstantTimeOps:  true,
		SecureMemory:     true,
		AriesConfig: &AriesConfig{
			KMSType:         "local",
			StorageProvider: "mem",
			CryptoSuite:     "BLS12381G2",
		},
	}
}

// ServiceInfo provides metadata about the BBS service implementation
type ServiceInfo struct {
	Provider          Provider  `json:"provider"`
	Version           string    `json:"version"`
	IsProductionReady bool      `json:"is_production_ready"`
	SupportedFeatures []string  `json:"supported_features"`
	CreatedAt         time.Time `json:"created_at"`
}

// PerformanceMetrics tracks operation performance
type PerformanceMetrics struct {
	KeyGenerationTime time.Duration `json:"key_generation_time"`
	SigningTime       time.Duration `json:"signing_time"`
	VerificationTime  time.Duration `json:"verification_time"`
	ProofCreationTime time.Duration `json:"proof_creation_time"`
	ProofVerifyTime   time.Duration `json:"proof_verify_time"`
	TotalOperations   int64         `json:"total_operations"`
	SuccessRate       float64       `json:"success_rate"`
}

// BBSServiceFactory creates BBS service instances based on provider
type BBSServiceFactory interface {
	CreateService(provider Provider, config *Config) (BBSInterface, error)
	GetSupportedProviders() []Provider
	ValidateConfig(provider Provider, config *Config) error
}
