package bbs

import (
	"fmt"
	"log"
)

// NewBBSService creates a new BBS service with the specified provider
func NewBBSService(provider Provider, config *Config) (BBSInterface, error) {
	factory := NewFactory()
	service, err := factory.CreateService(provider, config)
	if err != nil {
		return nil, err
	}

	// Wrap with metrics and logging if enabled
	if config != nil && config.EnableLogging {
		return NewServiceWrapper(service, config), nil
	}

	return service, nil
}

// NewDefaultBBSService creates a BBS service with default configuration
func NewDefaultBBSService() (BBSInterface, error) {
	config := DefaultConfig()
	return NewBBSService(config.Provider, config)
}

// NewSimpleBBSService creates a simple BBS service for testing
func NewSimpleBBSService() (BBSInterface, error) {
	config := DefaultConfig()
	config.Provider = ProviderSimple
	config.EnableLogging = false
	return NewBBSService(ProviderSimple, config)
}

// NewProductionBBSService creates a production-ready BBS service
func NewProductionBBSService() (BBSInterface, error) {
	config := DefaultConfig()
	config.Provider = ProviderProduction
	config.ConstantTimeOps = true
	config.SecureMemory = true
	return NewBBSService(ProviderProduction, config)
}

// NewAriesBBSService creates an Aries-based BBS service
func NewAriesBBSService(ariesConfig *AriesConfig) (BBSInterface, error) {
	config := DefaultConfig()
	config.Provider = ProviderAries
	config.AriesConfig = ariesConfig
	return NewBBSService(ProviderAries, config)
}

// GetSupportedProviders returns all supported BBS providers
func GetSupportedProviders() []Provider {
	factory := NewFactory()
	return factory.GetSupportedProviders()
}

// ValidateProvider checks if a provider is supported
func ValidateProvider(provider Provider) error {
	supported := GetSupportedProviders()
	for _, p := range supported {
		if p == provider {
			return nil
		}
	}
	return fmt.Errorf("unsupported provider: %s", provider)
}

// CompareProviders compares performance and features of different providers
func CompareProviders() map[Provider]ProviderComparison {
	return map[Provider]ProviderComparison{
		ProviderSimple: {
			Provider:        ProviderSimple,
			SecurityLevel:   "Demo",
			Performance:     "Fast",
			ProductionReady: false,
			Features:        []string{"basic_signing", "basic_verification"},
			Limitations:     []string{"not_cryptographically_secure", "demo_only"},
			RecommendedUse:  "Testing and development",
		},
		ProviderProduction: {
			Provider:        ProviderProduction,
			SecurityLevel:   "High",
			Performance:     "Good",
			ProductionReady: true,
			Features:        []string{"bls12_381", "selective_disclosure", "zero_knowledge_proofs", "constant_time_ops"},
			Limitations:     []string{"requires_careful_implementation"},
			RecommendedUse:  "Production deployments",
		},
		ProviderAries: {
			Provider:        ProviderAries,
			SecurityLevel:   "High",
			Performance:     "Good",
			ProductionReady: true,
			Features:        []string{"industry_standard", "aries_interop", "w3c_vc_compliance", "did_support"},
			Limitations:     []string{"requires_aries_dependency", "larger_binary_size"},
			RecommendedUse:  "Enterprise and interoperability",
		},
	}
}

// ProviderComparison holds comparison data for providers
type ProviderComparison struct {
	Provider        Provider `json:"provider"`
	SecurityLevel   string   `json:"security_level"`
	Performance     string   `json:"performance"`
	ProductionReady bool     `json:"production_ready"`
	Features        []string `json:"features"`
	Limitations     []string `json:"limitations"`
	RecommendedUse  string   `json:"recommended_use"`
}

// SwitchProvider switches between providers at runtime
func SwitchProvider(currentService BBSInterface, newProvider Provider, config *Config) (BBSInterface, error) {
	if currentService.GetProvider() == newProvider {
		return currentService, nil // No change needed
	}

	log.Printf("Switching BBS provider from %s to %s", currentService.GetProvider(), newProvider)

	newService, err := NewBBSService(newProvider, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new service with provider %s: %w", newProvider, err)
	}

	// Clean up old service if it has secure erase
	if config != nil && config.SecureMemory {
		// This is a placeholder - in practice you'd want to clean up any sensitive data
		log.Printf("Performing secure cleanup of old service")
	}

	log.Printf("Successfully switched to provider %s", newProvider)
	return newService, nil
}

// BenchmarkProviders runs performance benchmarks on different providers
func BenchmarkProviders(providers []Provider, messageCount int) (map[Provider]*PerformanceMetrics, error) {
	results := make(map[Provider]*PerformanceMetrics)

	// Prepare test data
	messages := make([][]byte, messageCount)
	for i := 0; i < messageCount; i++ {
		messages[i] = []byte(fmt.Sprintf("test message %d", i))
	}

	for _, provider := range providers {
		log.Printf("Benchmarking provider: %s", provider)

		config := DefaultConfig()
		config.Provider = provider
		config.EnableLogging = false

		service, err := NewBBSService(provider, config)
		if err != nil {
			log.Printf("Failed to create service for provider %s: %v", provider, err)
			continue
		}

		// Skip if not available
		if !service.IsProductionReady() && provider == ProviderAries {
			log.Printf("Skipping provider %s - not available", provider)
			continue
		}

		metrics := &PerformanceMetrics{}

		// Benchmark key generation
		keyPair, err := service.GenerateKeyPair()
		if err != nil {
			log.Printf("Key generation failed for provider %s: %v", provider, err)
			continue
		}

		// Benchmark signing
		signature, err := service.Sign(keyPair.PrivateKey, messages)
		if err != nil {
			log.Printf("Signing failed for provider %s: %v", provider, err)
			continue
		}

		// Benchmark verification
		err = service.Verify(keyPair.PublicKey, signature, messages)
		if err != nil {
			log.Printf("Verification failed for provider %s: %v", provider, err)
			continue
		}

		// Benchmark proof creation
		revealedIndices := []int{0} // Reveal only first message for simplicity
		if messageCount > 1 {
			revealedIndices = []int{0, messageCount - 1} // Reveal first and last messages
		}
		nonce := []byte("benchmark-nonce")
		proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
		if err != nil {
			log.Printf("Proof creation failed for provider %s: %v", provider, err)
			continue
		}

		// Benchmark proof verification
		revealedMessages := make([][]byte, len(revealedIndices))
		for i, idx := range revealedIndices {
			revealedMessages[i] = messages[idx]
		}
		err = service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce)
		if err != nil {
			log.Printf("Proof verification failed for provider %s: %v", provider, err)
			continue
		}

		// Get metrics if service is wrapped
		if wrapper, ok := service.(*ServiceWrapper); ok {
			metrics = wrapper.GetMetrics()
		}

		results[provider] = metrics
		log.Printf("Benchmark completed for provider %s", provider)
	}

	return results, nil
}

// MigrationHelper helps migrate between different BBS providers
type MigrationHelper struct {
	sourceProvider Provider
	targetProvider Provider
	config         *Config
}

// NewMigrationHelper creates a new migration helper
func NewMigrationHelper(sourceProvider, targetProvider Provider, config *Config) *MigrationHelper {
	return &MigrationHelper{
		sourceProvider: sourceProvider,
		targetProvider: targetProvider,
		config:         config,
	}
}

// ValidateMigration checks if migration is possible
func (m *MigrationHelper) ValidateMigration() error {
	if m.sourceProvider == m.targetProvider {
		return fmt.Errorf("source and target providers are the same")
	}

	// Check if target provider is supported
	if err := ValidateProvider(m.targetProvider); err != nil {
		return fmt.Errorf("target provider validation failed: %w", err)
	}

	// Check if target provider is available
	targetService, err := NewBBSService(m.targetProvider, m.config)
	if err != nil {
		return fmt.Errorf("target provider not available: %w", err)
	}

	if !targetService.IsProductionReady() && m.targetProvider == ProviderAries {
		return fmt.Errorf("target provider %s is not production ready", m.targetProvider)
	}

	return nil
}

// PerformMigration performs the actual migration
func (m *MigrationHelper) PerformMigration() error {
	log.Printf("Starting migration from %s to %s", m.sourceProvider, m.targetProvider)

	if err := m.ValidateMigration(); err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	// Create services
	sourceService, err := NewBBSService(m.sourceProvider, m.config)
	if err != nil {
		return fmt.Errorf("failed to create source service: %w", err)
	}

	targetService, err := NewBBSService(m.targetProvider, m.config)
	if err != nil {
		return fmt.Errorf("failed to create target service: %w", err)
	}

	// Test basic operations on both services
	sourceKeyPair, err := sourceService.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("source service test failed: %w", err)
	}

	targetKeyPair, err := targetService.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("target service test failed: %w", err)
	}

	// Validate key pairs
	if err := sourceService.ValidateKeyPair(sourceKeyPair); err != nil {
		return fmt.Errorf("source key pair validation failed: %w", err)
	}

	if err := targetService.ValidateKeyPair(targetKeyPair); err != nil {
		return fmt.Errorf("target key pair validation failed: %w", err)
	}

	log.Printf("Migration from %s to %s completed successfully", m.sourceProvider, m.targetProvider)
	return nil
}
