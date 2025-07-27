package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lugondev/bbs-selective-disclosure-example/internal/holder"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/issuer"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/verifier"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/did"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"
)

func main() {
	fmt.Println("🔐 BBS+ Age Verification Demo (18+ without revealing exact age/DOB)")
	fmt.Println("==================================================================")

	// Initialize services
	didRepo := did.NewInMemoryRepository()
	didService := did.NewService(didRepo)
	bbsService := bbs.NewService()
	credRepo := vc.NewInMemoryCredentialRepository()
	presRepo := vc.NewInMemoryPresentationRepository()
	vcService := vc.NewService(bbsService, credRepo, presRepo)

	// Initialize use cases
	issuerUC := issuer.NewUseCase(didService, vcService, bbsService)
	holderUC := holder.NewUseCase(didService, vcService, credRepo)
	verifierUC := verifier.NewUseCase(didService, vcService, presRepo)

	// Demo scenario
	if err := runAgeVerificationDemo(issuerUC, holderUC, verifierUC); err != nil {
		log.Fatalf("Age verification demo failed: %v", err)
	}

	fmt.Println("\n✅ Age verification demo completed successfully!")
}

func runAgeVerificationDemo(issuerUC *issuer.UseCase, holderUC *holder.UseCase, verifierUC *verifier.UseCase) error {
	// Step 1: Setup Government Authority (Issuer)
	fmt.Println("\n🏛️  Step 1: Setting up Government Authority (Digital ID Issuer)")
	issuerSetup, err := issuerUC.SetupIssuer("example")
	if err != nil {
		return fmt.Errorf("failed to setup issuer: %w", err)
	}
	fmt.Printf("✓ Government Authority DID: %s\n", issuerSetup.DID.String())

	// Step 2: Setup Citizen (Holder)
	fmt.Println("\n👤 Step 2: Setting up Citizen (Credential Holder)")
	holderSetup, err := holderUC.SetupHolder("example")
	if err != nil {
		return fmt.Errorf("failed to setup holder: %w", err)
	}
	fmt.Printf("✓ Citizen DID: %s\n", holderSetup.DID.String())

	// Step 3: Setup Age-Restricted Service (Verifier)
	fmt.Println("\n🎮 Step 3: Setting up Age-Restricted Service (Online Gaming Platform)")
	verifierSetup, err := verifierUC.SetupVerifier("example")
	if err != nil {
		return fmt.Errorf("failed to setup verifier: %w", err)
	}
	fmt.Printf("✓ Gaming Platform DID: %s\n", verifierSetup.DID.String())

	// Step 4: Government issues enhanced credential with age proofs
	fmt.Println("\n📄 Step 4: Government issuing enhanced digital ID with age verification claims")

	// Create enhanced claims including age verification attributes
	// Instead of just storing dateOfBirth, we create multiple derived claims
	birthYear := 1995 // Example: person born in 1995 (28 years old)

	claims := []vc.Claim{
		// Personal information
		{Key: "firstName", Value: "Minh"},
		{Key: "lastName", Value: "Tran Duc"},
		{Key: "fullName", Value: "Tran Duc Minh"},
		{Key: "dateOfBirth", Value: "1995-03-15"}, // Actual birth date
		{Key: "placeOfBirth", Value: "Ha Noi, Vietnam"},
		{Key: "nationality", Value: "Vietnamese"},
		{Key: "idNumber", Value: "987654321"},
		{Key: "address", Value: "456 Le Loi St, District 1, Ho Chi Minh City"},

		// Age verification claims (derived from dateOfBirth)
		{Key: "ageOver13", Value: true},
		{Key: "ageOver16", Value: true},
		{Key: "ageOver18", Value: true},
		{Key: "ageOver21", Value: true},
		{Key: "ageOver25", Value: true},
		{Key: "birthYear", Value: birthYear},
		{Key: "ageCategory", Value: "adult"}, // child, teen, adult, senior

		// Additional verification claims
		{Key: "issuedAt", Value: time.Now().Format("2006-01-02")},
		{Key: "documentType", Value: "national_id"},
		{Key: "validUntil", Value: "2030-03-15"},
	}

	credential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
		IssuerDID:  issuerSetup.DID.String(),
		SubjectDID: holderSetup.DID.String(),
		Claims:     claims,
	})
	if err != nil {
		return fmt.Errorf("failed to issue credential: %w", err)
	}

	fmt.Printf("✓ Enhanced credential issued with ID: %s\n", credential.ID)
	fmt.Printf("  Total claims: %d\n", len(claims))
	fmt.Println("  Age verification claims: ageOver13, ageOver16, ageOver18, ageOver21, ageOver25")

	// Step 5: Citizen stores the credential
	fmt.Println("\n💾 Step 5: Citizen storing enhanced credential")
	if err := holderUC.StoreCredential(credential); err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}
	fmt.Println("✓ Enhanced credential stored successfully")

	// Step 6: Gaming platform requests age verification (18+)
	fmt.Println("\n🎮 Step 6: Gaming platform requesting age verification")
	fmt.Println("  Gaming platform requirements:")
	fmt.Println("  - User must be 18+ years old")
	fmt.Println("  - Verify nationality for regional content")
	fmt.Println("  - Does NOT need: exact age, birth date, name, address, ID number")

	verificationNonce := fmt.Sprintf("gaming-age-verification-%d", time.Now().UnixMilli())
	fmt.Printf("  Generated verification nonce: %s\n", verificationNonce)

	// Step 7: Citizen creates selective disclosure presentation (Privacy-Preserving)
	fmt.Println("\n🔒 Step 7: Creating privacy-preserving age verification presentation")

	// Only reveal the minimum necessary information
	selectiveDisclosure := []vc.SelectiveDisclosureRequest{
		{
			CredentialID: credential.ID,
			RevealedAttributes: []string{
				"ageOver18",    // Boolean claim - proves 18+ without revealing exact age
				"nationality",  // Required for regional content
				"documentType", // Proves this is from official government ID
			},
		},
	}

	presentation, err := holderUC.CreatePresentation(holder.PresentationRequest{
		HolderDID:           holderSetup.DID.String(),
		CredentialIDs:       []string{credential.ID},
		SelectiveDisclosure: selectiveDisclosure,
		Nonce:               verificationNonce,
	})
	if err != nil {
		return fmt.Errorf("failed to create presentation: %w", err)
	}

	fmt.Printf("✓ Privacy-preserving presentation created with ID: %s\n", presentation.ID)
	fmt.Println("  ✅ REVEALED attributes:")
	fmt.Println("    - ageOver18: proves user is 18+ (boolean)")
	fmt.Println("    - nationality: required for content")
	fmt.Println("    - documentType: proves government-issued")
	fmt.Println("  🔒 HIDDEN attributes (remain private):")
	fmt.Println("    - firstName, lastName, fullName")
	fmt.Println("    - dateOfBirth (exact birth date)")
	fmt.Println("    - exact age (28 years old)")
	fmt.Println("    - birthYear (1995)")
	fmt.Println("    - address, idNumber")
	fmt.Println("    - placeOfBirth")

	// Step 8: Gaming platform verifies the presentation
	fmt.Println("\n🔍 Step 8: Gaming platform verifying age presentation")

	verificationResult, err := verifierUC.VerifyPresentation(verifier.VerificationRequest{
		Presentation:      presentation,
		RequiredClaims:    []string{"ageOver18", "nationality"},
		TrustedIssuers:    []string{issuerSetup.DID.String()},
		VerificationNonce: verificationNonce,
	})
	if err != nil {
		return fmt.Errorf("failed to verify presentation: %w", err)
	}

	fmt.Printf("✓ Verification result: %v\n", verificationResult.Valid)
	if len(verificationResult.Errors) > 0 {
		fmt.Printf("  Errors: %v\n", verificationResult.Errors)
	}

	// Step 9: Display verification results
	fmt.Println("\n📊 Step 9: Gaming platform verification results")
	fmt.Printf("  Holder DID: %s\n", verificationResult.HolderDID)
	fmt.Printf("  Trusted Issuer: %v\n", verificationResult.IssuerDIDs)
	fmt.Printf("  Revealed claims:\n")
	for key, value := range verificationResult.RevealedClaims {
		fmt.Printf("    %s: %v\n", key, value)
	}

	// Step 10: Business logic for age verification
	fmt.Println("\n✅ Step 10: Age verification business logic")
	if ageOver18, ok := verificationResult.RevealedClaims["ageOver18"].(bool); ok {
		if ageOver18 {
			fmt.Println("  🎉 ACCESS GRANTED: User is verified to be 18+ years old")
			fmt.Println("  🎮 User can access age-restricted gaming content")
		} else {
			fmt.Println("  ❌ ACCESS DENIED: User is under 18 years old")
		}
	}

	if nationality, ok := verificationResult.RevealedClaims["nationality"].(string); ok {
		fmt.Printf("  🌍 Regional content: Available for %s users\n", nationality)
	}

	if docType, ok := verificationResult.RevealedClaims["documentType"].(string); ok {
		fmt.Printf("  📄 Document verification: %s (government-issued)\n", docType)
	}

	// Step 11: Privacy protection demonstration
	fmt.Println("\n🛡️  Step 11: Privacy Protection Achieved")
	fmt.Println("  Gaming platform knows:")
	fmt.Println("    ✓ User is 18+ (meets legal requirement)")
	fmt.Println("    ✓ User is Vietnamese (for regional content)")
	fmt.Println("    ✓ Verified by government authority (trustworthy)")
	fmt.Println()
	fmt.Println("  Gaming platform CANNOT determine:")
	fmt.Println("    🔒 User's exact age (could be 18, 25, 35, 50...)")
	fmt.Println("    🔒 User's birth date (month/day/year unknown)")
	fmt.Println("    🔒 User's real name")
	fmt.Println("    🔒 User's address")
	fmt.Println("    🔒 User's ID number")
	fmt.Println("    🔒 User's birth year")

	// Step 12: Demonstrate multiple use cases
	fmt.Println("\n🎯 Step 12: Multiple Age Verification Scenarios")

	// Different services requiring different age thresholds
	scenarios := []struct {
		service     string
		requirement string
		claim       string
	}{
		{"Social Media", "13+", "ageOver13"},
		{"Movie Theater (R-rated)", "17+", "ageOver16"}, // Close enough to 17
		{"Alcohol Purchase", "21+", "ageOver21"},
		{"Senior Discount", "65+", "ageOver25"}, // We only have up to 25 in this example
	}

	for _, scenario := range scenarios {
		if claimValue, exists := claims[findClaimIndex(claims, scenario.claim)].Value.(bool); exists && claimValue {
			fmt.Printf("  ✅ %s (%s): ELIGIBLE\n", scenario.service, scenario.requirement)
		} else {
			fmt.Printf("  ❌ %s (%s): NOT ELIGIBLE\n", scenario.service, scenario.requirement)
		}
	}

	// Step 13: Technical details
	fmt.Println("\n🔧 Step 13: Technical Implementation Details")
	fmt.Println("  BBS+ Features utilized:")
	fmt.Println("    • Selective disclosure of attributes")
	fmt.Println("    • Zero-knowledge proof of age threshold")
	fmt.Println("    • Unlinkable presentations")
	fmt.Println("    • Tamper-evident credentials")
	fmt.Println()
	fmt.Println("  Privacy-by-Design principles:")
	fmt.Println("    • Minimal data disclosure")
	fmt.Println("    • Purpose limitation")
	fmt.Println("    • Data minimization")
	fmt.Println("    • User control over personal data")

	return nil
}

func findClaimIndex(claims []vc.Claim, key string) int {
	for i, claim := range claims {
		if claim.Key == key {
			return i
		}
	}
	return -1
}
