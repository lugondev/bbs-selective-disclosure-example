package bbs

import (
	"fmt"
	"log"
	"time"
)

// DefaultFactory implements BBSServiceFactory
type DefaultFactory struct {
	supportedProviders []Provider
}

// NewFactory creates a new BBS service factory
func NewFactory() BBSServiceFactory {
	return &DefaultFactory{
		supportedProviders: []Provider{
			ProviderSimple,
			ProviderProduction,
			ProviderAries,
		},
	}
}

// CreateService creates a BBS service instance based on provider
func (f *DefaultFactory) CreateService(provider Provider, config *Config) (BBSInterface, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Validate configuration
	if err := f.ValidateConfig(provider, config); err != nil {
		return nil, fmt.Errorf("invalid config for provider %s: %w", provider, err)
	}

	switch provider {
	case ProviderSimple:
		return newSimpleService(config), nil
	case ProviderProduction:
		return newProductionService(config), nil
	case ProviderAries:
		return newAriesService(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// GetSupportedProviders returns list of supported providers
func (f *DefaultFactory) GetSupportedProviders() []Provider {
	return f.supportedProviders
}

// ValidateConfig validates configuration for a specific provider
func (f *DefaultFactory) ValidateConfig(provider Provider, config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Common validation
	if config.OperationTimeout <= 0 {
		return fmt.Errorf("operation timeout must be positive")
	}

	// Provider-specific validation
	switch provider {
	case ProviderSimple:
		// Simple provider doesn't need special validation
		return nil
	case ProviderProduction:
		// Production provider validation
		if !config.ConstantTimeOps {
			log.Printf("Warning: constant time operations disabled for production provider")
		}
		return nil
	case ProviderAries:
		return f.validateAriesConfig(config.AriesConfig)
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
}

// validateAriesConfig validates Aries-specific configuration
func (f *DefaultFactory) validateAriesConfig(ariesConfig *AriesConfig) error {
	if ariesConfig == nil {
		return fmt.Errorf("aries config is required for aries provider")
	}

	validKMSTypes := map[string]bool{
		"local":  true,
		"remote": true,
	}
	if !validKMSTypes[ariesConfig.KMSType] {
		return fmt.Errorf("invalid KMS type: %s", ariesConfig.KMSType)
	}

	validStorageProviders := map[string]bool{
		"mem":     true,
		"leveldb": true,
	}
	if !validStorageProviders[ariesConfig.StorageProvider] {
		return fmt.Errorf("invalid storage provider: %s", ariesConfig.StorageProvider)
	}

	if ariesConfig.KMSType == "remote" {
		if ariesConfig.RemoteKMSURL == "" {
			return fmt.Errorf("remote KMS URL is required for remote KMS type")
		}
	}

	return nil
}

// ServiceWrapper wraps a BBS service with common functionality
type ServiceWrapper struct {
	service BBSInterface
	config  *Config
	metrics *PerformanceMetrics
	info    *ServiceInfo
}

// NewServiceWrapper creates a new service wrapper
func NewServiceWrapper(service BBSInterface, config *Config) *ServiceWrapper {
	return &ServiceWrapper{
		service: service,
		config:  config,
		metrics: &PerformanceMetrics{
			TotalOperations: 0,
			SuccessRate:     1.0,
		},
		info: &ServiceInfo{
			Provider:          service.GetProvider(),
			Version:           service.GetVersion(),
			IsProductionReady: service.IsProductionReady(),
			CreatedAt:         time.Now(),
			SupportedFeatures: []string{
				"key_generation",
				"signing",
				"verification",
				"proof_creation",
				"proof_verification",
				"selective_disclosure",
			},
		},
	}
}

// GenerateKeyPair generates a key pair with metrics tracking
func (w *ServiceWrapper) GenerateKeyPair() (*KeyPair, error) {
	start := time.Now()
	w.metrics.TotalOperations++

	result, err := w.service.GenerateKeyPair()
	w.metrics.KeyGenerationTime = time.Since(start)

	if err != nil {
		w.updateSuccessRate(false)
		if w.config.EnableLogging {
			log.Printf("Key generation failed: %v", err)
		}
		return nil, err
	}

	w.updateSuccessRate(true)
	if w.config.EnableLogging {
		log.Printf("Key generation completed in %v", w.metrics.KeyGenerationTime)
	}

	return result, nil
}

// Sign creates a signature with metrics tracking
func (w *ServiceWrapper) Sign(privateKey []byte, messages [][]byte) (*Signature, error) {
	start := time.Now()
	w.metrics.TotalOperations++

	result, err := w.service.Sign(privateKey, messages)
	w.metrics.SigningTime = time.Since(start)

	if err != nil {
		w.updateSuccessRate(false)
		if w.config.EnableLogging {
			log.Printf("Signing failed: %v", err)
		}
		return nil, err
	}

	w.updateSuccessRate(true)
	if w.config.EnableLogging {
		log.Printf("Signing completed in %v for %d messages", w.metrics.SigningTime, len(messages))
	}

	return result, nil
}

// Verify verifies a signature with metrics tracking
func (w *ServiceWrapper) Verify(publicKey []byte, signature *Signature, messages [][]byte) error {
	start := time.Now()
	w.metrics.TotalOperations++

	var err error
	if w.config.ConstantTimeOps {
		err = w.service.ConstantTimeVerify(publicKey, signature, messages)
	} else {
		err = w.service.Verify(publicKey, signature, messages)
	}

	w.metrics.VerificationTime = time.Since(start)

	if err != nil {
		w.updateSuccessRate(false)
		if w.config.EnableLogging {
			log.Printf("Verification failed: %v", err)
		}
		return err
	}

	w.updateSuccessRate(true)
	if w.config.EnableLogging {
		log.Printf("Verification completed in %v", w.metrics.VerificationTime)
	}

	return nil
}

// CreateProof creates a proof with metrics tracking
func (w *ServiceWrapper) CreateProof(signature *Signature, publicKey []byte, messages [][]byte, revealedIndices []int, nonce []byte) (*Proof, error) {
	start := time.Now()
	w.metrics.TotalOperations++

	result, err := w.service.CreateProof(signature, publicKey, messages, revealedIndices, nonce)
	w.metrics.ProofCreationTime = time.Since(start)

	if err != nil {
		w.updateSuccessRate(false)
		if w.config.EnableLogging {
			log.Printf("Proof creation failed: %v", err)
		}
		return nil, err
	}

	w.updateSuccessRate(true)
	if w.config.EnableLogging {
		log.Printf("Proof creation completed in %v", w.metrics.ProofCreationTime)
	}

	return result, nil
}

// VerifyProof verifies a proof with metrics tracking
func (w *ServiceWrapper) VerifyProof(publicKey []byte, proof *Proof, revealedMessages [][]byte, nonce []byte) error {
	start := time.Now()
	w.metrics.TotalOperations++

	err := w.service.VerifyProof(publicKey, proof, revealedMessages, nonce)
	w.metrics.ProofVerifyTime = time.Since(start)

	if err != nil {
		w.updateSuccessRate(false)
		if w.config.EnableLogging {
			log.Printf("Proof verification failed: %v", err)
		}
		return err
	}

	w.updateSuccessRate(true)
	if w.config.EnableLogging {
		log.Printf("Proof verification completed in %v", w.metrics.ProofVerifyTime)
	}

	return nil
}

// ValidateKeyPair validates a key pair
func (w *ServiceWrapper) ValidateKeyPair(keyPair *KeyPair) error {
	return w.service.ValidateKeyPair(keyPair)
}

// GetMessageCount gets message count
func (w *ServiceWrapper) GetMessageCount(signature *Signature, publicKey []byte) (int, error) {
	return w.service.GetMessageCount(signature, publicKey)
}

// ConstantTimeVerify performs constant time verification
func (w *ServiceWrapper) ConstantTimeVerify(publicKey []byte, signature *Signature, messages [][]byte) error {
	return w.service.ConstantTimeVerify(publicKey, signature, messages)
}

// SecureErase securely erases data
func (w *ServiceWrapper) SecureErase(data []byte) {
	if w.config.SecureMemory {
		w.service.SecureErase(data)
	}
}

// GetProvider returns the provider type
func (w *ServiceWrapper) GetProvider() Provider {
	return w.service.GetProvider()
}

// GetVersion returns the version
func (w *ServiceWrapper) GetVersion() string {
	return w.service.GetVersion()
}

// IsProductionReady returns production readiness
func (w *ServiceWrapper) IsProductionReady() bool {
	return w.service.IsProductionReady()
}

// GetMetrics returns performance metrics
func (w *ServiceWrapper) GetMetrics() *PerformanceMetrics {
	return w.metrics
}

// GetInfo returns service information
func (w *ServiceWrapper) GetInfo() *ServiceInfo {
	return w.info
}

// updateSuccessRate updates the success rate metric
func (w *ServiceWrapper) updateSuccessRate(success bool) {
	if w.metrics.TotalOperations == 1 {
		if success {
			w.metrics.SuccessRate = 1.0
		} else {
			w.metrics.SuccessRate = 0.0
		}
		return
	}

	// Calculate running average
	current := w.metrics.SuccessRate * float64(w.metrics.TotalOperations-1)
	if success {
		current += 1.0
	}
	w.metrics.SuccessRate = current / float64(w.metrics.TotalOperations)
}
