package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lugondev/bbs-selective-disclosure-example/internal/holder"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/issuer"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/verifier"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/did"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/vc"
)

// TestFullLifecycle tests the complete DID -> VC -> VP workflow
func TestFullLifecycle(t *testing.T) {
	// Setup services
	didRepo := did.NewInMemoryRepository()
	didService := did.NewService(didRepo)
	bbsService := bbs.NewService()
	credRepo := vc.NewInMemoryCredentialRepository()
	presRepo := vc.NewInMemoryPresentationRepository()
	vcService := vc.NewService(bbsService, credRepo, presRepo)

	// Setup use cases
	issuerUC := issuer.NewUseCase(didService, vcService, bbsService)
	holderUC := holder.NewUseCase(didService, vcService, credRepo)
	verifierUC := verifier.NewUseCase(didService, vcService, presRepo)

	t.Run("Complete Selective Disclosure Workflow", func(t *testing.T) {
		// Step 1: Setup participants
		issuerSetup, err := issuerUC.SetupIssuer("test")
		require.NoError(t, err)
		assert.NotNil(t, issuerSetup.DID)
		assert.NotNil(t, issuerSetup.BBSKeyPair)

		holderSetup, err := holderUC.SetupHolder("test")
		require.NoError(t, err)
		assert.NotNil(t, holderSetup.DID)

		_, err = verifierUC.SetupVerifier("test")
		require.NoError(t, err)

		// Step 2: Issue credential
		claims := []vc.Claim{
			{Key: "firstName", Value: "John"},
			{Key: "lastName", Value: "Doe"},
			{Key: "age", Value: 25},
			{Key: "nationality", Value: "American"},
			{Key: "email", Value: "john.doe@example.com"},
		}

		credential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
			IssuerDID:  issuerSetup.DID.String(),
			SubjectDID: holderSetup.DID.String(),
			Claims:     claims,
		})
		require.NoError(t, err)
		assert.NotNil(t, credential)
		assert.Equal(t, issuerSetup.DID.String(), credential.Issuer)
		assert.Equal(t, holderSetup.DID.String(), credential.CredentialSubject["id"])

		// Step 3: Holder stores credential
		err = holderUC.StoreCredential(credential)
		require.NoError(t, err)

		// Verify credential is stored
		storedCreds, err := holderUC.ListCredentials(holderSetup.DID.String())
		require.NoError(t, err)
		assert.Len(t, storedCreds, 1)
		assert.Equal(t, credential.ID, storedCreds[0].ID)

		// Step 4: Create selective disclosure presentation
		// Only reveal age and nationality, hide name and email
		selectiveDisclosure := []vc.SelectiveDisclosureRequest{
			{
				CredentialID:       credential.ID,
				RevealedAttributes: []string{"age", "nationality"},
			},
		}

		presentation, err := holderUC.CreatePresentation(holder.PresentationRequest{
			HolderDID:           holderSetup.DID.String(),
			CredentialIDs:       []string{credential.ID},
			SelectiveDisclosure: selectiveDisclosure,
		})
		require.NoError(t, err)
		assert.NotNil(t, presentation)
		assert.Equal(t, holderSetup.DID.String(), presentation.Holder)
		assert.Len(t, presentation.VerifiableCredential, 1)

		// Step 5: Verify presentation
		verificationResult, err := verifierUC.VerifyPresentation(verifier.VerificationRequest{
			Presentation:   presentation,
			RequiredClaims: []string{"age", "nationality"},
			TrustedIssuers: []string{issuerSetup.DID.String()},
		})
		require.NoError(t, err)
		assert.True(t, verificationResult.Valid)
		assert.Len(t, verificationResult.Errors, 0)

		// Verify only requested attributes are revealed
		assert.Equal(t, 25, verificationResult.RevealedClaims["age"])
		assert.Equal(t, "American", verificationResult.RevealedClaims["nationality"])

		// Verify hidden attributes are not present
		assert.NotContains(t, verificationResult.RevealedClaims, "firstName")
		assert.NotContains(t, verificationResult.RevealedClaims, "lastName")
		assert.NotContains(t, verificationResult.RevealedClaims, "email")
	})
}

// TestMultipleCredentialsPresentation tests presenting multiple credentials
func TestMultipleCredentialsPresentation(t *testing.T) {
	// Setup
	didRepo := did.NewInMemoryRepository()
	didService := did.NewService(didRepo)
	bbsService := bbs.NewService()
	credRepo := vc.NewInMemoryCredentialRepository()
	presRepo := vc.NewInMemoryPresentationRepository()
	vcService := vc.NewService(bbsService, credRepo, presRepo)

	issuerUC := issuer.NewUseCase(didService, vcService, bbsService)
	holderUC := holder.NewUseCase(didService, vcService, credRepo)
	verifierUC := verifier.NewUseCase(didService, vcService, presRepo)

	// Setup participants
	issuerSetup, err := issuerUC.SetupIssuer("test")
	require.NoError(t, err)

	holderSetup, err := holderUC.SetupHolder("test")
	require.NoError(t, err)

	_, err = verifierUC.SetupVerifier("test")
	require.NoError(t, err)

	// Issue ID credential
	idClaims := []vc.Claim{
		{Key: "fullName", Value: "Jane Smith"},
		{Key: "dateOfBirth", Value: "1995-03-15"},
		{Key: "nationality", Value: "Canadian"},
	}

	idCredential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
		IssuerDID:  issuerSetup.DID.String(),
		SubjectDID: holderSetup.DID.String(),
		Claims:     idClaims,
	})
	require.NoError(t, err)

	// Issue degree credential
	degreeClaims := []vc.Claim{
		{Key: "degree", Value: "Bachelor of Science"},
		{Key: "major", Value: "Computer Science"},
		{Key: "graduationYear", Value: 2017},
		{Key: "university", Value: "University of Toronto"},
	}

	degreeCredential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
		IssuerDID:  issuerSetup.DID.String(),
		SubjectDID: holderSetup.DID.String(),
		Claims:     degreeClaims,
	})
	require.NoError(t, err)

	// Store credentials
	err = holderUC.StoreCredential(idCredential)
	require.NoError(t, err)

	err = holderUC.StoreCredential(degreeCredential)
	require.NoError(t, err)

	// Create presentation with selective disclosure from both credentials
	selectiveDisclosure := []vc.SelectiveDisclosureRequest{
		{
			CredentialID:       idCredential.ID,
			RevealedAttributes: []string{"nationality"}, // Only nationality from ID
		},
		{
			CredentialID:       degreeCredential.ID,
			RevealedAttributes: []string{"degree", "major"}, // Only degree and major
		},
	}

	presentation, err := holderUC.CreatePresentation(holder.PresentationRequest{
		HolderDID:           holderSetup.DID.String(),
		CredentialIDs:       []string{idCredential.ID, degreeCredential.ID},
		SelectiveDisclosure: selectiveDisclosure,
	})
	require.NoError(t, err)
	assert.Len(t, presentation.VerifiableCredential, 2)

	// Verify presentation
	verificationResult, err := verifierUC.VerifyPresentation(verifier.VerificationRequest{
		Presentation:   presentation,
		RequiredClaims: []string{"nationality", "degree", "major"},
		TrustedIssuers: []string{issuerSetup.DID.String()},
	})
	require.NoError(t, err)
	assert.True(t, verificationResult.Valid)

	// Verify revealed claims
	assert.Equal(t, "Canadian", verificationResult.RevealedClaims["nationality"])
	assert.Equal(t, "Bachelor of Science", verificationResult.RevealedClaims["degree"])
	assert.Equal(t, "Computer Science", verificationResult.RevealedClaims["major"])

	// Verify hidden claims
	assert.NotContains(t, verificationResult.RevealedClaims, "fullName")
	assert.NotContains(t, verificationResult.RevealedClaims, "dateOfBirth")
	assert.NotContains(t, verificationResult.RevealedClaims, "graduationYear")
	assert.NotContains(t, verificationResult.RevealedClaims, "university")
}

// TestVerificationFailures tests various verification failure scenarios
func TestVerificationFailures(t *testing.T) {
	// Setup
	didRepo := did.NewInMemoryRepository()
	didService := did.NewService(didRepo)
	bbsService := bbs.NewService()
	credRepo := vc.NewInMemoryCredentialRepository()
	presRepo := vc.NewInMemoryPresentationRepository()
	vcService := vc.NewService(bbsService, credRepo, presRepo)

	issuerUC := issuer.NewUseCase(didService, vcService, bbsService)
	holderUC := holder.NewUseCase(didService, vcService, credRepo)
	verifierUC := verifier.NewUseCase(didService, vcService, presRepo)

	// Setup participants
	issuerSetup, err := issuerUC.SetupIssuer("test")
	require.NoError(t, err)

	untrustedIssuerSetup, err := issuerUC.SetupIssuer("test")
	require.NoError(t, err)

	holderSetup, err := holderUC.SetupHolder("test")
	require.NoError(t, err)

	_, err = verifierUC.SetupVerifier("test")
	require.NoError(t, err)

	t.Run("Untrusted Issuer", func(t *testing.T) {
		// Issue credential from untrusted issuer
		claims := []vc.Claim{
			{Key: "name", Value: "Test User"},
			{Key: "role", Value: "Admin"},
		}

		credential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
			IssuerDID:  untrustedIssuerSetup.DID.String(),
			SubjectDID: holderSetup.DID.String(),
			Claims:     claims,
		})
		require.NoError(t, err)

		err = holderUC.StoreCredential(credential)
		require.NoError(t, err)

		presentation, err := holderUC.CreatePresentation(holder.PresentationRequest{
			HolderDID:     holderSetup.DID.String(),
			CredentialIDs: []string{credential.ID},
			SelectiveDisclosure: []vc.SelectiveDisclosureRequest{
				{
					CredentialID:       credential.ID,
					RevealedAttributes: []string{"name", "role"},
				},
			},
		})
		require.NoError(t, err)

		// Verify with trusted issuers list (untrusted issuer not included)
		verificationResult, err := verifierUC.VerifyPresentation(verifier.VerificationRequest{
			Presentation:   presentation,
			RequiredClaims: []string{"name", "role"},
			TrustedIssuers: []string{issuerSetup.DID.String()}, // Only trusted issuer
		})
		require.NoError(t, err)
		assert.False(t, verificationResult.Valid)
		assert.Contains(t, verificationResult.Errors[0], "is not trusted")
	})

	t.Run("Missing Required Claims", func(t *testing.T) {
		claims := []vc.Claim{
			{Key: "name", Value: "Test User"},
		}

		credential, err := issuerUC.IssueCredential(issuer.IssueCredentialRequest{
			IssuerDID:  issuerSetup.DID.String(),
			SubjectDID: holderSetup.DID.String(),
			Claims:     claims,
		})
		require.NoError(t, err)

		err = holderUC.StoreCredential(credential)
		require.NoError(t, err)

		presentation, err := holderUC.CreatePresentation(holder.PresentationRequest{
			HolderDID:     holderSetup.DID.String(),
			CredentialIDs: []string{credential.ID},
			SelectiveDisclosure: []vc.SelectiveDisclosureRequest{
				{
					CredentialID:       credential.ID,
					RevealedAttributes: []string{"name"}, // Only reveal name
				},
			},
		})
		require.NoError(t, err)

		// Verify requiring both name and age (age not revealed)
		verificationResult, err := verifierUC.VerifyPresentation(verifier.VerificationRequest{
			Presentation:   presentation,
			RequiredClaims: []string{"name", "age"}, // age is required but not revealed
			TrustedIssuers: []string{issuerSetup.DID.String()},
		})
		require.NoError(t, err)
		assert.False(t, verificationResult.Valid)
		assert.Contains(t, verificationResult.Errors[0], "required claim 'age' is missing")
	})
}

// TestDIDOperations tests DID creation and resolution
func TestDIDOperations(t *testing.T) {
	didRepo := did.NewInMemoryRepository()
	didService := did.NewService(didRepo)

	t.Run("DID Creation and Resolution", func(t *testing.T) {
		// Generate DID
		generatedDID, keyPair, err := didService.GenerateDID("test")
		require.NoError(t, err)
		assert.NotNil(t, generatedDID)
		assert.NotNil(t, keyPair)
		assert.Equal(t, "test", generatedDID.Method)

		// Create DID document
		didDoc, err := didService.CreateDIDDocument(generatedDID, keyPair)
		require.NoError(t, err)
		assert.NotNil(t, didDoc)
		assert.Equal(t, generatedDID.String(), didDoc.ID)
		assert.Len(t, didDoc.VerificationMethod, 1)

		// Store DID document
		err = didRepo.Create(didDoc)
		require.NoError(t, err)

		// Resolve DID
		resolvedDoc, err := didService.ResolveDID(generatedDID.String())
		require.NoError(t, err)
		assert.Equal(t, didDoc.ID, resolvedDoc.ID)
		assert.Equal(t, didDoc.VerificationMethod[0].ID, resolvedDoc.VerificationMethod[0].ID)

		// Verify DID document
		err = didService.VerifyDIDDocument(resolvedDoc)
		require.NoError(t, err)
	})

	t.Run("DID Document Validation", func(t *testing.T) {
		// Test invalid DID document
		invalidDoc := &did.DIDDocument{}
		err := didService.VerifyDIDDocument(invalidDoc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DID document ID is empty")

		// Test DID document with no verification methods but with authentication
		invalidDoc2 := &did.DIDDocument{
			ID:             "did:test:123",
			Authentication: []string{"invalid-key-id"},
		}
		err = didService.VerifyDIDDocument(invalidDoc2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DID document must have at least one verification method")
	})
}

// TestBBSOperations tests BBS+ signing and proof operations
func TestBBSOperations(t *testing.T) {
	bbsService := bbs.NewService()

	t.Run("Key Generation and Signing", func(t *testing.T) {
		// Generate key pair
		keyPair, err := bbsService.GenerateKeyPair()
		require.NoError(t, err)
		assert.NotNil(t, keyPair.PublicKey)
		assert.NotNil(t, keyPair.PrivateKey)
		assert.Len(t, keyPair.PrivateKey, 32)
		assert.Len(t, keyPair.PublicKey, 32)

		// Sign messages
		messages := [][]byte{
			[]byte("message1"),
			[]byte("message2"),
			[]byte("message3"),
		}

		signature, err := bbsService.Sign(keyPair.PrivateKey, messages)
		require.NoError(t, err)
		assert.NotNil(t, signature)
		assert.Len(t, signature.Value, 32)

		// Verify signature
		err = bbsService.Verify(keyPair.PublicKey, signature, messages)
		require.NoError(t, err)
	})

	t.Run("Selective Disclosure Proof", func(t *testing.T) {
		keyPair, err := bbsService.GenerateKeyPair()
		require.NoError(t, err)

		messages := [][]byte{
			[]byte("secret1"),
			[]byte("secret2"),
			[]byte("public1"),
			[]byte("public2"),
		}

		signature, err := bbsService.Sign(keyPair.PrivateKey, messages)
		require.NoError(t, err)

		// Create proof revealing only indices 2 and 3 (public1, public2)
		revealedIndices := []int{2, 3}
		nonce := []byte("verification-nonce")

		proof, err := bbsService.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
		require.NoError(t, err)
		assert.NotNil(t, proof)
		assert.Equal(t, revealedIndices, proof.RevealedAttributes)
		assert.Equal(t, nonce, proof.Nonce)

		// Verify proof with revealed messages
		revealedMessages := [][]byte{
			messages[2], // public1
			messages[3], // public2
		}

		err = bbsService.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce)
		require.NoError(t, err)
	})
}
