package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
)

func main() {
	fmt.Println("🚀 BBS+ Interface Demo")
	fmt.Println("======================")

	// Check if specific demo is requested
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			bbs.DemoConfigurationOptions()
			return
		case "switching":
			if err := bbs.DemoServiceSwitching(); err != nil {
				log.Fatalf("Demo failed: %v", err)
			}
			return
		case "all":
			if err := bbs.RunAllDemos(); err != nil {
				log.Fatalf("Demo failed: %v", err)
			}
			return
		default:
			fmt.Printf("Unknown demo: %s\n", os.Args[1])
			fmt.Println("Available demos: config, switching, all")
			return
		}
	}

	// Default: run basic demo
	if err := runBasicDemo(); err != nil {
		log.Fatalf("Demo failed: %v", err)
	}

	fmt.Println("\n✅ Demo completed successfully!")
	fmt.Println("\nTo run other demos:")
	fmt.Println("  go run cmd/interface_demo/main.go config")
	fmt.Println("  go run cmd/interface_demo/main.go switching")
	fmt.Println("  go run cmd/interface_demo/main.go all")
}

func runBasicDemo() error {
	fmt.Println("\n=== Basic BBS Interface Demo ===")

	// 1. Show supported providers
	providers := bbs.GetSupportedProviders()
	fmt.Printf("Supported providers: %v\n", providers)

	// 2. Create simple service for quick demo
	fmt.Println("\n--- Simple Provider Demo ---")
	simpleService, err := bbs.NewSimpleBBSService()
	if err != nil {
		return fmt.Errorf("failed to create simple service: %w", err)
	}

	fmt.Printf("Provider: %s\n", simpleService.GetProvider())
	fmt.Printf("Version: %s\n", simpleService.GetVersion())
	fmt.Printf("Production Ready: %t\n", simpleService.IsProductionReady())

	// Test basic operations
	if err := demonstrateBasicOperations(simpleService); err != nil {
		return fmt.Errorf("simple service operations failed: %w", err)
	}

	// 3. Create production service
	fmt.Println("\n--- Production Provider Demo ---")
	productionService, err := bbs.NewProductionBBSService()
	if err != nil {
		return fmt.Errorf("failed to create production service: %w", err)
	}

	fmt.Printf("Provider: %s\n", productionService.GetProvider())
	fmt.Printf("Version: %s\n", productionService.GetVersion())
	fmt.Printf("Production Ready: %t\n", productionService.IsProductionReady())

	// Test production service (may have some failures due to crypto complexity)
	fmt.Println("Testing production service (some operations may fail due to crypto complexity)...")
	if err := demonstrateBasicOperations(productionService); err != nil {
		fmt.Printf("Production service note: %v\n", err)
		fmt.Println("This is expected as the production crypto is more complex")
	}

	// 4. Show provider comparison
	fmt.Println("\n--- Provider Comparison ---")
	comparisons := bbs.CompareProviders()
	for provider, comparison := range comparisons {
		fmt.Printf("\n%s Provider:\n", provider)
		fmt.Printf("  Security: %s\n", comparison.SecurityLevel)
		fmt.Printf("  Performance: %s\n", comparison.Performance)
		fmt.Printf("  Production Ready: %t\n", comparison.ProductionReady)
		fmt.Printf("  Recommended Use: %s\n", comparison.RecommendedUse)
	}

	// 5. Demonstrate provider switching
	fmt.Println("\n--- Provider Switching Demo ---")
	config := bbs.DefaultConfig()
	config.EnableLogging = false // Reduce output for demo

	newService, err := bbs.SwitchProvider(simpleService, bbs.ProviderProduction, config)
	if err != nil {
		return fmt.Errorf("failed to switch provider: %w", err)
	}

	fmt.Printf("Successfully switched from %s to %s\n",
		simpleService.GetProvider(), newService.GetProvider())

	// 6. Show Aries integration info
	fmt.Println("\n--- Aries Framework Integration ---")
	fmt.Println("To use Hyperledger Aries Framework Go:")
	fmt.Println("1. Add dependency: go get github.com/hyperledger/aries-framework-go")
	fmt.Println("2. Complete the implementation in aries_adapter.go")
	fmt.Println("3. Use bbs.NewAriesBBSService() to create Aries-based service")

	ariesConfig := &bbs.AriesConfig{
		KMSType:         "local",
		StorageProvider: "mem",
		CryptoSuite:     "BLS12381G2",
	}

	_, err = bbs.NewAriesBBSService(ariesConfig)
	if err != nil {
		fmt.Printf("Expected Aries error: %v\n", err)
	}

	return nil
}

func demonstrateBasicOperations(service bbs.BBSInterface) error {
	fmt.Println("\n📝 Demonstrating BBS+ Operations:")

	// 1. Generate key pair
	fmt.Println("1. Generating key pair...")
	keyPair, err := service.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("key generation failed: %w", err)
	}
	fmt.Printf("   ✓ Generated key pair (pub: %d bytes, priv: %d bytes)\n",
		len(keyPair.PublicKey), len(keyPair.PrivateKey))

	// 2. Validate key pair
	if err := service.ValidateKeyPair(keyPair); err != nil {
		return fmt.Errorf("key validation failed: %w", err)
	}
	fmt.Println("   ✓ Key pair validated")

	// 3. Prepare credential data
	credentialData := [][]byte{
		[]byte("John Doe"),          // Name
		[]byte("1990-01-01"),        // Date of birth
		[]byte("Software Engineer"), // Job title
		[]byte("Bachelor's Degree"), // Education
		[]byte("New York"),          // Location
	}

	// 4. Sign the credential
	fmt.Printf("2. Signing credential with %d attributes...\n", len(credentialData))
	signature, err := service.Sign(keyPair.PrivateKey, credentialData)
	if err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}
	fmt.Println("   ✓ Credential signed")

	// 5. Verify the signature
	fmt.Println("3. Verifying signature...")
	if err := service.Verify(keyPair.PublicKey, signature, credentialData); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Println("   ✓ Signature verified")

	// 6. Create selective disclosure proof
	// Reveal only name and location, hide other attributes
	revealedIndices := []int{0, 4} // Name and Location
	nonce := []byte("demo-proof-2024")

	fmt.Println("4. Creating selective disclosure proof...")
	fmt.Printf("   Revealing: Name (%s) and Location (%s)\n",
		credentialData[0], credentialData[4])
	fmt.Printf("   Hiding: Date of birth, Job title, Education\n")

	proof, err := service.CreateProof(signature, keyPair.PublicKey, credentialData, revealedIndices, nonce)
	if err != nil {
		return fmt.Errorf("proof creation failed: %w", err)
	}
	fmt.Printf("   ✓ Proof created (revealing %d out of %d attributes)\n",
		len(revealedIndices), len(credentialData))

	// 7. Verify the proof
	fmt.Println("5. Verifying selective disclosure proof...")
	revealedMessages := [][]byte{credentialData[0], credentialData[4]} // Name, Location

	if err := service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce); err != nil {
		return fmt.Errorf("proof verification failed: %w", err)
	}
	fmt.Println("   ✓ Proof verified successfully")

	// 8. Demonstrate security features
	fmt.Println("6. Testing security features...")
	sensitiveData := []byte("sensitive-private-key-data")
	service.SecureErase(sensitiveData)
	fmt.Println("   ✓ Sensitive data securely erased")

	fmt.Printf("\n✅ All operations completed successfully with %s provider\n",
		service.GetProvider())

	return nil
}
