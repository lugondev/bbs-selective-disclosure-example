package main

import (
	"encoding/json"
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
	fmt.Println("🔐 BBS+ Selective Disclosure Demo")
	fmt.Println("=====================================")

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
	if err := runDemo(issuerUC, holderUC, verifierUC); err != nil {
		log.Fatalf("Demo failed: %v", err)
	}

	fmt.Println("\n✅ Demo completed successfully!")
}

func runDemo(issuerUC *issuer.UseCase, holderUC *holder.UseCase, verifierUC *verifier.UseCase) error {
	// Step 1: Setup Issuer (Government ID Authority)
	fmt.Println("\n📋 Step 1: Setting up Issuer (Government ID Authority)")
	issuerSetup, err := issuerUC.SetupIssuer("example")
	if err != nil {
		return fmt.Errorf("failed to setup issuer: %w", err)
	}
	fmt.Printf("✓ Issuer DID: %s\n", issuerSetup.DID.String())

	// Step 2: Setup Holder (Citizen)
	fmt.Println("\n👤 Step 2: Setting up Holder (Citizen)")
	holderSetup, err := holderUC.SetupHolder("example")
	if err != nil {
		return fmt.Errorf("failed to setup holder: %w", err)
	}
	fmt.Printf("✓ Holder DID: %s\n", holderSetup.DID.String())

	// Step 3: Setup Verifier (Cinema)
	fmt.Println("\n🎬 Step 3: Setting up Verifier (Cinema)")
	verifierSetup, err := verifierUC.SetupVerifier("example")
	if err != nil {
		return fmt.Errorf("failed to setup verifier: %w", err)
	}
	fmt.Printf("✓ Verifier DID: %s\n", verifierSetup.DID.String())

	// Step 4: Issue Digital ID Credential
	fmt.Println("\n📄 Step 4: Issuing Digital ID Credential")

	// Create claims for a digital ID
	claims := []vc.Claim{
		{Key: "firstName", Value: "An"},
		{Key: "lastName", Value: "Nguyen Van"},
		{Key: "dateOfBirth", Value: "2000-01-20"},
		{Key: "nationality", Value: "Vietnamese"},
		{Key: "address", Value: "123 Nguyen Trai St, Ho Chi Minh City"},
		{Key: "idNumber", Value: "123456789"},
	}

	credential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
		IssuerDID:  issuerSetup.DID.String(),
		SubjectDID: holderSetup.DID.String(),
		Claims:     claims,
	})
	if err != nil {
		return fmt.Errorf("failed to issue credential: %w", err)
	}

	fmt.Printf("✓ Credential issued with ID: %s\n", credential.ID)
	fmt.Printf("  Claims: %v\n", getClaimKeys(claims))

	// Step 5: Holder stores the credential
	fmt.Println("\n💾 Step 5: Holder storing credential")
	if err := holderUC.StoreCredential(credential); err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}
	fmt.Println("✓ Credential stored successfully")

	// Step 6: Cinema requests age and nationality verification
	fmt.Println("\n🎭 Step 6: Cinema requests age and nationality verification")
	fmt.Println("  Cinema needs to verify:")
	fmt.Println("  - Age (18+): needs dateOfBirth")
	fmt.Println("  - Nationality: needs nationality")
	fmt.Println("  - Does NOT need: firstName, lastName, address, idNumber")

	// Generate verification nonce
	verificationNonce := "cinema-verification-" + fmt.Sprintf("%d", time.Now().UnixMilli())
	fmt.Printf("  Generated verification nonce: %s\n", verificationNonce)

	// Step 7: Holder creates selective disclosure presentation
	fmt.Println("\n🎪 Step 7: Creating selective disclosure presentation")

	selectiveDisclosure := []vc.SelectiveDisclosureRequest{
		{
			CredentialID:       credential.ID,
			RevealedAttributes: []string{"dateOfBirth", "nationality"}, // Only reveal these
		},
	}

	presentation, err := holderUC.CreatePresentation(holder.PresentationRequest{
		HolderDID:           holderSetup.DID.String(),
		CredentialIDs:       []string{credential.ID},
		SelectiveDisclosure: selectiveDisclosure,
		Nonce:               verificationNonce, // Use the verification nonce
	})
	if err != nil {
		return fmt.Errorf("failed to create presentation: %w", err)
	}

	fmt.Printf("✓ Presentation created with ID: %s\n", presentation.ID)
	fmt.Println("  Revealed attributes: dateOfBirth, nationality")
	fmt.Println("  Hidden attributes: firstName, lastName, address, idNumber")

	// Step 8: Verifier verifies the presentation
	fmt.Println("\n🔍 Step 8: Cinema verifying presentation")

	verificationResult, err := verifierUC.VerifyPresentation(verifier.VerificationRequest{
		Presentation:      presentation,
		RequiredClaims:    []string{"dateOfBirth", "nationality"},
		TrustedIssuers:    []string{issuerSetup.DID.String()},
		VerificationNonce: verificationNonce, // Use the same verification nonce
	})
	if err != nil {
		return fmt.Errorf("failed to verify presentation: %w", err)
	}

	fmt.Printf("✓ Verification result: %v\n", verificationResult.Valid)
	if len(verificationResult.Errors) > 0 {
		fmt.Printf("  Errors: %v\n", verificationResult.Errors)
	}

	// Step 9: Display revealed information
	fmt.Println("\n📊 Step 9: Information available to Cinema")
	fmt.Printf("  Holder DID: %s\n", verificationResult.HolderDID)
	fmt.Printf("  Issuer DIDs: %v\n", verificationResult.IssuerDIDs)
	fmt.Printf("  Revealed claims:\n")
	for key, value := range verificationResult.RevealedClaims {
		fmt.Printf("    %s: %v\n", key, value)
	}

	// Step 10: Age verification logic
	fmt.Println("\n🎂 Step 10: Age verification")
	if dateOfBirth, ok := verificationResult.RevealedClaims["dateOfBirth"].(string); ok {
		age := calculateAge(dateOfBirth)
		fmt.Printf("  Calculated age: %d years\n", age)
		if age >= 18 {
			fmt.Println("  ✅ Age verification: PASSED (18+)")
		} else {
			fmt.Println("  ❌ Age verification: FAILED (under 18)")
		}
	}

	if nationality, ok := verificationResult.RevealedClaims["nationality"].(string); ok {
		fmt.Printf("  Nationality: %s\n", nationality)
		fmt.Println("  ✅ Nationality verification: PASSED")
	}

	// Step 11: Privacy demonstration
	fmt.Println("\n🔒 Step 11: Privacy Protection Demonstration")
	fmt.Println("  Cinema CANNOT see:")
	fmt.Println("    - firstName: [HIDDEN]")
	fmt.Println("    - lastName: [HIDDEN]")
	fmt.Println("    - address: [HIDDEN]")
	fmt.Println("    - idNumber: [HIDDEN]")
	fmt.Println("  But can still verify that these attributes exist and are signed by a trusted issuer!")

	// Step 12: Display presentation structure
	fmt.Println("\n📋 Step 12: Technical Details")
	if len(presentation.VerifiableCredential) > 0 {
		credentialData, _ := json.MarshalIndent(presentation.VerifiableCredential[0], "  ", "  ")
		fmt.Printf("  Selective disclosure credential structure:\n  %s\n", credentialData)
	}

	return nil
}

func getClaimKeys(claims []vc.Claim) []string {
	keys := make([]string, len(claims))
	for i, claim := range claims {
		keys[i] = claim.Key
	}
	return keys
}

func calculateAge(dateOfBirth string) int {
	// Parse date of birth (format: YYYY-MM-DD)
	birthTime, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return 0
	}

	now := time.Now()
	age := now.Year() - birthTime.Year()

	// Adjust if birthday hasn't occurred this year
	if now.YearDay() < birthTime.YearDay() {
		age--
	}

	return age
}
