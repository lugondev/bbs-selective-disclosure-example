package bbs

import (
	"fmt"
	"log"
)

// DemoServiceSwitching demonstrates switching between different BBS providers
func DemoServiceSwitching() error {
	fmt.Println("=== BBS Provider Switching Demo ===")

	// 1. List supported providers
	providers := GetSupportedProviders()
	fmt.Printf("Supported providers: %v\n", providers)

	// 2. Compare providers
	comparisons := CompareProviders()
	for provider, comparison := range comparisons {
		fmt.Printf("\nProvider: %s\n", provider)
		fmt.Printf("  Security Level: %s\n", comparison.SecurityLevel)
		fmt.Printf("  Performance: %s\n", comparison.Performance)
		fmt.Printf("  Production Ready: %t\n", comparison.ProductionReady)
		fmt.Printf("  Features: %v\n", comparison.Features)
		fmt.Printf("  Recommended Use: %s\n", comparison.RecommendedUse)
	}

	// 3. Create simple service
	fmt.Println("\n=== Creating Simple Service ===")
	simpleService, err := NewSimpleBBSService()
	if err != nil {
		return fmt.Errorf("failed to create simple service: %w", err)
	}

	fmt.Printf("Created service with provider: %s\n", simpleService.GetProvider())
	fmt.Printf("Version: %s\n", simpleService.GetVersion())
	fmt.Printf("Production Ready: %t\n", simpleService.IsProductionReady())

	// Test simple service
	if err := testService(simpleService, "Simple"); err != nil {
		return fmt.Errorf("simple service test failed: %w", err)
	}

	// 4. Create production service
	fmt.Println("\n=== Creating Production Service ===")
	productionService, err := NewProductionBBSService()
	if err != nil {
		return fmt.Errorf("failed to create production service: %w", err)
	}

	fmt.Printf("Created service with provider: %s\n", productionService.GetProvider())
	fmt.Printf("Version: %s\n", productionService.GetVersion())
	fmt.Printf("Production Ready: %t\n", productionService.IsProductionReady())

	// Test production service
	if err := testService(productionService, "Production"); err != nil {
		log.Printf("Production service test failed (expected for some operations): %v", err)
	}

	// 5. Switch between providers
	fmt.Println("\n=== Provider Switching Demo ===")
	config := DefaultConfig()
	config.EnableLogging = true

	newService, err := SwitchProvider(simpleService, ProviderProduction, config)
	if err != nil {
		return fmt.Errorf("failed to switch provider: %w", err)
	}

	fmt.Printf("Switched to provider: %s\n", newService.GetProvider())

	// 6. Demo Aries service (will show error since not implemented)
	fmt.Println("\n=== Aries Service Demo ===")
	ariesConfig := &AriesConfig{
		KMSType:         "local",
		StorageProvider: "mem",
		CryptoSuite:     "BLS12381G2",
	}

	ariesService, err := NewAriesBBSService(ariesConfig)
	if err != nil {
		fmt.Printf("Expected error for Aries service: %v\n", err)
		fmt.Println("Integration guide:")
		fmt.Println(AriesIntegrationGuide())
	} else {
		fmt.Printf("Aries service created: %s\n", ariesService.GetProvider())
	}

	// 7. Migration demo
	fmt.Println("\n=== Migration Demo ===")
	migrationHelper := NewMigrationHelper(ProviderSimple, ProviderProduction, config)

	if err := migrationHelper.ValidateMigration(); err != nil {
		log.Printf("Migration validation failed: %v", err)
	} else {
		fmt.Println("Migration validation passed")

		if err := migrationHelper.PerformMigration(); err != nil {
			log.Printf("Migration failed: %v", err)
		} else {
			fmt.Println("Migration completed successfully")
		}
	}

	// 8. Benchmarking demo
	fmt.Println("\n=== Benchmarking Demo ===")
	benchmarkProviders := []Provider{ProviderSimple, ProviderProduction}

	results, err := BenchmarkProviders(benchmarkProviders, 3)
	if err != nil {
		log.Printf("Benchmarking failed: %v", err)
	} else {
		for provider, metrics := range results {
			fmt.Printf("Provider %s metrics:\n", provider)
			fmt.Printf("  Total Operations: %d\n", metrics.TotalOperations)
			fmt.Printf("  Success Rate: %.2f%%\n", metrics.SuccessRate*100)
		}
	}

	fmt.Println("\n=== Demo Completed ===")
	return nil
}

// testService performs basic operations on a BBS service
func testService(service BBSInterface, name string) error {
	fmt.Printf("\n--- Testing %s Service ---\n", name)

	// Generate key pair
	fmt.Println("1. Generating key pair...")
	keyPair, err := service.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("key generation failed: %w", err)
	}
	fmt.Printf("   ✓ Key pair generated (pub: %d bytes, priv: %d bytes)\n",
		len(keyPair.PublicKey), len(keyPair.PrivateKey))

	// Validate key pair
	fmt.Println("2. Validating key pair...")
	if err := service.ValidateKeyPair(keyPair); err != nil {
		return fmt.Errorf("key validation failed: %w", err)
	}
	fmt.Println("   ✓ Key pair is valid")

	// Prepare messages
	messages := [][]byte{
		[]byte("Alice"),
		[]byte("25"),
		[]byte("Engineer"),
		[]byte("New York"),
	}

	// Sign messages
	fmt.Printf("3. Signing %d messages...\n", len(messages))
	signature, err := service.Sign(keyPair.PrivateKey, messages)
	if err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}
	fmt.Println("   ✓ Messages signed")

	// Verify signature
	fmt.Println("4. Verifying signature...")
	if err := service.Verify(keyPair.PublicKey, signature, messages); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Println("   ✓ Signature verified")

	// Create proof (reveal only name and city, hide age and job)
	fmt.Println("5. Creating selective disclosure proof...")
	revealedIndices := []int{0, 3} // Reveal Alice and New York
	nonce := []byte("test-nonce-123")

	proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
	if err != nil {
		return fmt.Errorf("proof creation failed: %w", err)
	}
	fmt.Printf("   ✓ Proof created (revealing %d out of %d attributes)\n",
		len(revealedIndices), len(messages))

	// Verify proof
	fmt.Println("6. Verifying proof...")
	revealedMessages := [][]byte{messages[0], messages[3]} // Alice, New York

	if err := service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce); err != nil {
		return fmt.Errorf("proof verification failed: %w", err)
	}
	fmt.Println("   ✓ Proof verified")

	// Security operations
	fmt.Println("7. Testing security operations...")
	testData := []byte("sensitive-data-to-erase")
	service.SecureErase(testData)
	fmt.Println("   ✓ Secure erase completed")

	fmt.Printf("--- %s Service Test Completed Successfully ---\n", name)
	return nil
}

// DemoConfigurationOptions shows different configuration options
func DemoConfigurationOptions() {
	fmt.Println("=== Configuration Options Demo ===")

	// Default config
	defaultConfig := DefaultConfig()
	fmt.Printf("Default Config:\n")
	fmt.Printf("  Provider: %s\n", defaultConfig.Provider)
	fmt.Printf("  Enable Logging: %t\n", defaultConfig.EnableLogging)
	fmt.Printf("  Constant Time Ops: %t\n", defaultConfig.ConstantTimeOps)
	fmt.Printf("  Secure Memory: %t\n", defaultConfig.SecureMemory)
	fmt.Printf("  Operation Timeout: %v\n", defaultConfig.OperationTimeout)

	// Custom config for production
	prodConfig := &Config{
		Provider:         ProviderProduction,
		EnableLogging:    true,
		OperationTimeout: defaultConfig.OperationTimeout,
		ConstantTimeOps:  true,
		SecureMemory:     true,
	}

	fmt.Printf("\nProduction Config:\n")
	fmt.Printf("  Provider: %s\n", prodConfig.Provider)
	fmt.Printf("  Constant Time Ops: %t\n", prodConfig.ConstantTimeOps)
	fmt.Printf("  Secure Memory: %t\n", prodConfig.SecureMemory)

	// Aries config
	ariesConfig := &Config{
		Provider:         ProviderAries,
		EnableLogging:    true,
		OperationTimeout: defaultConfig.OperationTimeout,
		ConstantTimeOps:  true,
		SecureMemory:     true,
		AriesConfig: &AriesConfig{
			KMSType:         "remote",
			StorageProvider: "leveldb",
			CryptoSuite:     "BLS12381G2",
			RemoteKMSURL:    "https://kms.example.com",
			AuthToken:       "your-auth-token",
		},
	}

	fmt.Printf("\nAries Config:\n")
	fmt.Printf("  Provider: %s\n", ariesConfig.Provider)
	fmt.Printf("  KMS Type: %s\n", ariesConfig.AriesConfig.KMSType)
	fmt.Printf("  Storage Provider: %s\n", ariesConfig.AriesConfig.StorageProvider)
	fmt.Printf("  Crypto Suite: %s\n", ariesConfig.AriesConfig.CryptoSuite)
	fmt.Printf("  Remote KMS URL: %s\n", ariesConfig.AriesConfig.RemoteKMSURL)
}

// RunAllDemos runs all available demos
func RunAllDemos() error {
	fmt.Println("🚀 Starting BBS Interface Demo")
	fmt.Println("=====================================")

	// Configuration demo
	DemoConfigurationOptions()

	fmt.Println()

	// Main service switching demo
	if err := DemoServiceSwitching(); err != nil {
		return fmt.Errorf("service switching demo failed: %w", err)
	}

	fmt.Println("\n🎉 All demos completed successfully!")
	return nil
}
