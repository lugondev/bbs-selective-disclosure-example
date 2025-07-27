package bbs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair(t *testing.T) {
	service := NewService()

	keyPair, err := service.GenerateKeyPair()
	require.NoError(t, err)

	assert.NotNil(t, keyPair)
	assert.Len(t, keyPair.PublicKey, 192) // G2 point is 192 bytes
	assert.Len(t, keyPair.PrivateKey, 32)

	// Ensure different key pairs are generated
	keyPair2, err := service.GenerateKeyPair()
	require.NoError(t, err)
	assert.NotEqual(t, keyPair.PrivateKey, keyPair2.PrivateKey)
	assert.NotEqual(t, keyPair.PublicKey, keyPair2.PublicKey)

	// Test ValidateKeyPair
	err = service.ValidateKeyPair(keyPair)
	assert.NoError(t, err)
}

func TestSignAndVerify(t *testing.T) {
	service := NewService()

	keyPair, err := service.GenerateKeyPair()
	require.NoError(t, err)

	messages := [][]byte{
		[]byte("message1"),
		[]byte("message2"),
		[]byte("message3"),
	}

	t.Run("Valid Signature", func(t *testing.T) {
		signature, err := service.Sign(keyPair.PrivateKey, messages)
		require.NoError(t, err)
		assert.NotNil(t, signature)
		assert.Len(t, signature.A, 96) // G1 point is 96 bytes
		assert.Len(t, signature.E, 32) // Scalar is 32 bytes
		assert.Len(t, signature.S, 32) // Scalar is 32 bytes

		err = service.Verify(keyPair.PublicKey, signature, messages)
		assert.NoError(t, err)
	})

	t.Run("Invalid Private Key Length", func(t *testing.T) {
		invalidKey := []byte("invalid")
		_, err := service.Sign(invalidKey, messages)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid private key length")
	})

	t.Run("Invalid Public Key Length", func(t *testing.T) {
		signature, err := service.Sign(keyPair.PrivateKey, messages)
		require.NoError(t, err)

		invalidPubKey := []byte("invalid")
		err = service.Verify(invalidPubKey, signature, messages)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid public key length")
	})
}

func TestCreateAndVerifyProof(t *testing.T) {
	service := NewService()

	keyPair, err := service.GenerateKeyPair()
	require.NoError(t, err)

	messages := [][]byte{
		[]byte("secret1"), // index 0 - hidden
		[]byte("secret2"), // index 1 - hidden
		[]byte("public1"), // index 2 - revealed
		[]byte("public2"), // index 3 - revealed
	}

	signature, err := service.Sign(keyPair.PrivateKey, messages)
	require.NoError(t, err)

	t.Run("Valid Proof", func(t *testing.T) {
		revealedIndices := []int{2, 3}
		nonce := []byte("test-nonce")

		proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
		require.NoError(t, err)

		assert.NotNil(t, proof)
		assert.Equal(t, revealedIndices, proof.RevealedAttributes)
		assert.Equal(t, nonce, proof.Nonce)
		assert.Len(t, proof.A_prime, 96) // G1 point
		assert.Len(t, proof.A_bar, 96)   // G1 point
		assert.Len(t, proof.C, 32)       // Challenge hash

		// Verify proof with revealed messages
		revealedMessages := [][]byte{
			messages[2], // public1
			messages[3], // public2
		}

		err = service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce)
		assert.NoError(t, err)
	})

	t.Run("Empty Nonce", func(t *testing.T) {
		revealedIndices := []int{2}
		emptyNonce := []byte{}

		_, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, emptyNonce)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nonce is required")
	})

	t.Run("Invalid Revealed Index", func(t *testing.T) {
		revealedIndices := []int{10} // out of range
		nonce := []byte("test-nonce")

		_, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid revealed indices")
	})

	t.Run("Invalid Public Key for Proof Verification", func(t *testing.T) {
		revealedIndices := []int{2}
		nonce := []byte("test-nonce")

		proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
		require.NoError(t, err)

		invalidPubKey := []byte("invalid")
		revealedMessages := [][]byte{messages[2]}

		err = service.VerifyProof(invalidPubKey, proof, revealedMessages, nonce)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid public key length")
	})

	t.Run("Mismatched Revealed Messages and Indices", func(t *testing.T) {
		revealedIndices := []int{2, 3}
		nonce := []byte("test-nonce")

		proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
		require.NoError(t, err)

		// Wrong number of revealed messages
		wrongRevealedMessages := [][]byte{messages[2]} // should be 2 messages

		err = service.VerifyProof(keyPair.PublicKey, proof, wrongRevealedMessages, nonce)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mismatch between revealed messages and indices")
	})
}

func TestEncodeDecodeProof(t *testing.T) {
	service := NewService()

	keyPair, err := service.GenerateKeyPair()
	require.NoError(t, err)

	messages := [][]byte{
		[]byte("message1"),
		[]byte("message2"),
	}

	signature, err := service.Sign(keyPair.PrivateKey, messages)
	require.NoError(t, err)

	revealedIndices := []int{0, 1}
	nonce := []byte("test-nonce-for-encoding")

	proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
	require.NoError(t, err)

	t.Run("Encode and Decode", func(t *testing.T) {
		encoded := EncodeProof(proof)
		assert.NotEmpty(t, encoded)

		decoded, err := DecodeProof(encoded)
		require.NoError(t, err)

		assert.Equal(t, proof.A_prime, decoded.A_prime)
		assert.Equal(t, proof.A_bar, decoded.A_bar)
		assert.Equal(t, proof.RevealedAttributes, decoded.RevealedAttributes)
		assert.Equal(t, proof.Nonce, decoded.Nonce)
	})

	t.Run("Invalid Base64", func(t *testing.T) {
		invalidEncoded := "invalid-base64!!!"
		_, err := DecodeProof(invalidEncoded)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode proof")
	})

	t.Run("Invalid Data Length", func(t *testing.T) {
		// Create a short base64 string that decodes to less than 300 bytes
		shortData := "dGVzdA==" // "test" in base64
		_, err := DecodeProof(shortData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proof data length")
	})
}

func TestMultipleMessages(t *testing.T) {
	service := NewService()

	keyPair, err := service.GenerateKeyPair()
	require.NoError(t, err)

	// Test with many messages
	messages := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		messages[i] = []byte(fmt.Sprintf("message%d", i))
	}

	signature, err := service.Sign(keyPair.PrivateKey, messages)
	require.NoError(t, err)

	// Reveal only some messages
	revealedIndices := []int{2, 5, 8}
	nonce := []byte("multi-message-nonce")

	proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
	require.NoError(t, err)

	revealedMessages := [][]byte{
		messages[2],
		messages[5],
		messages[8],
	}

	err = service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce)
	assert.NoError(t, err)
}
