package did

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDIDString(t *testing.T) {
	did := DID{
		Method:     "example",
		Identifier: "123456789",
	}

	expected := "did:example:123456789"
	assert.Equal(t, expected, did.String())
}

func TestGenerateDID(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewService(repo)

	did, keyPair, err := service.GenerateDID("test")
	require.NoError(t, err)

	assert.NotNil(t, did)
	assert.Equal(t, "test", did.Method)
	assert.NotEmpty(t, did.Identifier)

	assert.NotNil(t, keyPair)
	assert.Len(t, keyPair.PublicKey, 32)
	assert.Len(t, keyPair.PrivateKey, 64)
	assert.Contains(t, keyPair.KeyID, did.String())
}

func TestCreateDIDDocument(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewService(repo)

	did, keyPair, err := service.GenerateDID("test")
	require.NoError(t, err)

	doc, err := service.CreateDIDDocument(did, keyPair)
	require.NoError(t, err)

	assert.Equal(t, did.String(), doc.ID)
	assert.Len(t, doc.VerificationMethod, 1)
	assert.Equal(t, keyPair.KeyID, doc.VerificationMethod[0].ID)
	assert.Equal(t, "Ed25519VerificationKey2020", doc.VerificationMethod[0].Type)
	assert.Equal(t, did.String(), doc.VerificationMethod[0].Controller)

	assert.Contains(t, doc.Authentication, keyPair.KeyID)
	assert.Contains(t, doc.AssertionMethod, keyPair.KeyID)
}

func TestInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository()

	t.Run("Store and Retrieve", func(t *testing.T) {
		doc := &DIDDocument{
			ID: "did:test:123",
			VerificationMethod: []VerificationMethod{
				{
					ID:         "did:test:123#key-1",
					Type:       "Ed25519VerificationKey2020",
					Controller: "did:test:123",
				},
			},
		}

		err := repo.Create(doc)
		require.NoError(t, err)

		retrieved, err := repo.Resolve("did:test:123")
		require.NoError(t, err)
		assert.Equal(t, doc.ID, retrieved.ID)
	})

	t.Run("Resolve Non-existent", func(t *testing.T) {
		_, err := repo.Resolve("did:test:nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DID document not found")
	})

	t.Run("Update", func(t *testing.T) {
		doc := &DIDDocument{
			ID: "did:test:update",
		}

		err := repo.Create(doc)
		require.NoError(t, err)

		doc.VerificationMethod = []VerificationMethod{
			{ID: "new-key"},
		}

		err = repo.Update("did:test:update", doc)
		require.NoError(t, err)

		updated, err := repo.Resolve("did:test:update")
		require.NoError(t, err)
		assert.Len(t, updated.VerificationMethod, 1)
		assert.Equal(t, "new-key", updated.VerificationMethod[0].ID)
	})

	t.Run("Deactivate", func(t *testing.T) {
		doc := &DIDDocument{
			ID: "did:test:deactivate",
		}

		err := repo.Create(doc)
		require.NoError(t, err)

		err = repo.Deactivate("did:test:deactivate")
		require.NoError(t, err)

		_, err = repo.Resolve("did:test:deactivate")
		assert.Error(t, err)
	})
}

func TestVerifyDIDDocument(t *testing.T) {
	service := NewService(NewInMemoryRepository())

	t.Run("Valid Document", func(t *testing.T) {
		doc := &DIDDocument{
			ID: "did:test:valid",
			VerificationMethod: []VerificationMethod{
				{
					ID:         "did:test:valid#key-1",
					Type:       "Ed25519VerificationKey2020",
					Controller: "did:test:valid",
				},
			},
			Authentication:  []string{"did:test:valid#key-1"},
			AssertionMethod: []string{"did:test:valid#key-1"},
		}

		err := service.VerifyDIDDocument(doc)
		assert.NoError(t, err)
	})

	t.Run("Nil Document", func(t *testing.T) {
		err := service.VerifyDIDDocument(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DID document is nil")
	})

	t.Run("Empty ID", func(t *testing.T) {
		doc := &DIDDocument{}
		err := service.VerifyDIDDocument(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DID document ID is empty")
	})

	t.Run("No Verification Methods", func(t *testing.T) {
		doc := &DIDDocument{
			ID: "did:test:no-vm",
		}
		err := service.VerifyDIDDocument(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have at least one verification method")
	})

	t.Run("Invalid Authentication Reference", func(t *testing.T) {
		doc := &DIDDocument{
			ID: "did:test:invalid-auth",
			VerificationMethod: []VerificationMethod{
				{ID: "did:test:invalid-auth#key-1"},
			},
			Authentication: []string{"did:test:invalid-auth#nonexistent"},
		}
		err := service.VerifyDIDDocument(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication method")
		assert.Contains(t, err.Error(), "not found")
	})
}
