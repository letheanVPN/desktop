package openpgp

import (
	"fmt"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
)

// Mock key store to simulate key retrieval for testing purposes.
var mockSignKeyStore = make(map[string]*openpgp.Entity)

// mockSignReadEntity simulates the readEntity function from key.go
func mockSignReadEntity(keyArmored string) (*openpgp.Entity, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(keyArmored))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}
	if len(entityList) == 0 {
		return nil, fmt.Errorf("no entity found in key")
	}
	return entityList[0], nil
}

// mockSignGetPublicKey simulates GetPublicKey using the mockSignKeyStore
func mockSignGetPublicKey(id string) (*openpgp.Entity, error) {
	entity, ok := mockSignKeyStore[id]
	if !ok {
		return nil, fmt.Errorf("public key for id %s not found in mock store", id)
	}
	return entity, nil
}

// mockSignGetPrivateKey simulates GetPrivateKey using the mockSignKeyStore
func mockSignGetPrivateKey(id, passphrase string) (*openpgp.Entity, error) {
	entity, ok := mockSignKeyStore[id]
	if !ok {
		return nil, fmt.Errorf("private key for id %s not found in mock store", id)
	}

	// Create a copy to avoid modifying the stored entity directly
	copiedEntity := *entity
	copiedEntity.PrivateKey = entity.PrivateKey // Shallow copy of pointer is fine here, as Decrypt operates on the struct it points to.

	if copiedEntity.PrivateKey == nil {
		return nil, fmt.Errorf("no private key found for id %s", id)
	}

	if copiedEntity.PrivateKey.Encrypted {
		if err := copiedEntity.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
			return nil, fmt.Errorf("failed to decrypt private key for id %s: %w", id, err)
		}
	}
	return &copiedEntity, nil
}

//func TestSignAndVerify(t *testing.T) {
//	// Clear mock key store before each test
//	mockSignKeyStore = make(map[string]*openpgp.Entity)
//
//	signerID := "signer@example.com"
//	signerPassphrase := "signer_passphrase"
//	data := "This is the data to be signed."
//
//	// 1. Generate signer key pair and store in mock store
//	signerKeyPair, err := CreateKeyPair(signerID, signerPassphrase)
//	assert.NoError(t, err)
//	signerEntity, err := mockSignReadEntity(signerKeyPair.PrivateKey) // Use private key to get full entity
//	assert.NoError(t, err)
//	mockSignKeyStore[signerID] = signerEntity
//
//	// Replace actual GetPrivateKey and GetPublicKey with mocks for this test
//	originalGetPrivateKey := GetPrivateKey
//	originalGetPublicKey := GetPublicKey
//	GetPrivateKey = mockSignGetPrivateKey
//	GetPublicKey = mockSignGetPublicKey
//	defer func() {
//		GetPrivateKey = originalGetPrivateKey
//		GetPublicKey = originalGetPublicKey
//	}()
//
//	// Test successful signing and verification
//	t.Run("successful signing and verification", func(t *testing.T) {
//		signature, err := Sign(data, signerID, signerPassphrase)
//		assert.NoError(t, err)
//		assert.NotEmpty(t, signature)
//
//		verified, err := Verify(data, signature, signerID)
//		assert.NoError(t, err)
//		assert.True(t, verified, "Signature should be valid")
//	})
//
//	// Test signing with non-existent private key ID
//	t.Run("sign with non-existent private key ID", func(t *testing.T) {
//		_, err := Sign(data, "nonexistent_signer", signerPassphrase)
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "private key for id nonexistent_signer not found in mock store")
//	})
//
//	// Test signing with wrong passphrase
//	t.Run("sign with wrong passphrase", func(t *testing.T) {
//		_, err := Sign(data, signerID, "wrong_passphrase")
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "failed to decrypt private key")
//	})
//
//	// Test verification with non-existent public key ID
//	t.Run("verify with non-existent public key ID", func(t *testing.T) {
//		signature, err := Sign(data, signerID, signerPassphrase)
//		assert.NoError(t, err)
//
//		_, err = Verify(data, signature, "nonexistent_verifier")
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "public key for id nonexistent_verifier not found in mock store")
//	})
//
//	// Test verification with tampered signature
//	t.Run("verify with tampered signature", func(t *testing.T) {
//		signature, err := Sign(data, signerID, signerPassphrase)
//		assert.NoError(t, err)
//
//		tamperedSignature := signature + "tamper"
//		verified, err := Verify(data, tamperedSignature, signerID)
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "signature verification failed")
//		assert.False(t, verified)
//	})
//
//	// Test verification with tampered data
//	t.Run("verify with tampered data", func(t *testing.T) {
//		signature, err := Sign(data, signerID, signerPassphrase)
//		assert.NoError(t, err)
//
//		tamperedData := data + "tamper"
//		verified, err := Verify(tamperedData, signature, signerID)
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "signature verification failed")
//		assert.False(t, verified)
//	})
//
//	// Test verification with a different public key (not the signer's)
//	t.Run("verify with different public key", func(t *testing.T) {
//		// Create another key pair
//		anotherSignerID := "another_signer@example.com"
//		anotherSignerPassphrase := "another_signer_passphrase"
//		anotherSignerKeyPair, err := CreateKeyPair(anotherSignerID, anotherSignerPassphrase)
//		assert.NoError(t, err)
//		anotherSignerEntity, err := mockSignReadEntity(anotherSignerKeyPair.PrivateKey)
//		assert.NoError(t, err)
//		mockSignKeyStore[anotherSignerID] = anotherSignerEntity
//
//		signature, err := Sign(data, signerID, signerPassphrase)
//		assert.NoError(t, err)
//
//		// Try to verify with the public key of 'anotherSignerID'
//		verified, err := Verify(data, signature, anotherSignerID)
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "signature verification failed")
//		assert.False(t, verified)
//	})
//}
