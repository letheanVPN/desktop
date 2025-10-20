package openpgp

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/stretchr/testify/assert"
)

// Mock key store to simulate filesystem for GetPublicKey and GetPrivateKey
var mockKeyStore = make(map[string]*openpgp.Entity)

// mockReadEntity simulates the readEntity function from key.go
func mockReadEntity(keyArmored string) (*openpgp.Entity, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(keyArmored))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}
	if len(entityList) == 0 {
		return nil, fmt.Errorf("no entity found in key")
	}
	return entityList[0], nil
}

// mockGetPublicKey simulates GetPublicKey using the mockKeyStore
func mockGetPublicKey(id string) (*openpgp.Entity, error) {
	entity, ok := mockKeyStore[id]
	if !ok {
		return nil, fmt.Errorf("public key for id %s not found in mock store", id)
	}
	return entity, nil
}

// mockGetPrivateKey simulates GetPrivateKey using the mockKeyStore
func mockGetPrivateKey(id, passphrase string) (*openpgp.Entity, error) {
	entity, ok := mockKeyStore[id]
	if !ok {
		return nil, fmt.Errorf("private key for id %s not found in mock store", id)
	}

	// Create a copy to avoid modifying the stored entity directly
	// This is important because Decrypt modifies the PrivateKey field
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

// testEncryptPGP mimics the EncryptPGP function but uses mock key retrieval
func testEncryptPGP(recipientID, data string, signerID, signerPassphrase *string) (string, error) {
	recipient, err := mockGetPublicKey(recipientID)
	if err != nil {
		return "", fmt.Errorf("failed to get recipient public key: %w", err)
	}

	var signer *openpgp.Entity
	if signerID != nil && signerPassphrase != nil {
		signer, err = mockGetPrivateKey(*signerID, *signerPassphrase)
		if err != nil {
			return "", fmt.Errorf("could not get private key for signing: %w", err)
		}
	}

	buf := new(bytes.Buffer)
	armoredWriter, err := armor.Encode(buf, pgpMessageHeader, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create armored writer: %w", err)
	}

	plaintextWriter, err := openpgp.Encrypt(armoredWriter, []*openpgp.Entity{recipient}, signer, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}

	if _, err := plaintextWriter.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("failed to write plaintext data: %w", err)
	}

	if err := plaintextWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close plaintext writer: %w", err)
	}
	if err := armoredWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close armored writer: %w", err)
	}

	return buf.String(), nil
}

// testDecryptPGP mimics the DecryptPGP function but uses mock key retrieval
func testDecryptPGP(recipientID, message, passphrase string, signerID *string) (string, error) {
	privateKeyEntity, err := mockGetPrivateKey(recipientID, passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	keyring := openpgp.EntityList{privateKeyEntity}
	var expectedSigner *openpgp.Entity

	if signerID != nil {
		publicKeyEntity, err := mockGetPublicKey(*signerID)
		if err != nil {
			return "", fmt.Errorf("could not get public key for verification: %w", err)
		}
		keyring = append(keyring, publicKeyEntity)
		expectedSigner = publicKeyEntity
	}

	md, err := openpgp.ReadMessage(strings.NewReader(message), keyring, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read PGP message: %w", err)
	}

	decrypted, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", fmt.Errorf("failed to read decrypted body: %w", err)
	}

	if signerID != nil {
		if md.SignatureError != nil {
			return "", fmt.Errorf("signature verification failed: %w", md.SignatureError)
		}
		if md.SignedBy == nil {
			return "", fmt.Errorf("message is not signed, but signature verification was requested")
		}
		if expectedSigner.PrimaryKey.KeyId != md.SignedBy.PublicKey.KeyId {
			return "", fmt.Errorf("signature from unexpected key id: got %X, want %X", md.SignedBy.PublicKey.KeyId, expectedSigner.PrimaryKey.KeyId)
		}
	}

	return string(decrypted), nil
}

func TestEncryptDecryptPGP(t *testing.T) {
	// Clear mock key store before each test
	mockKeyStore = make(map[string]*openpgp.Entity)

	recipientID := "recipient"
	recipientPassphrase := "recipient_passphrase"
	data := "This is a secret message."

	// 1. Generate recipient key pair and store in mock store
	recipientKeyPair, err := CreateKeyPair(recipientID, recipientPassphrase)
	assert.NoError(t, err)
	recipientEntity, err := mockReadEntity(recipientKeyPair.PrivateKey) // Use private key to get full entity
	assert.NoError(t, err)
	mockKeyStore[recipientID] = recipientEntity

	// Test basic encryption and decryption
	t.Run("basic encryption and decryption", func(t *testing.T) {
		encrypted, err := testEncryptPGP(recipientID, data, nil, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		decrypted, err := testDecryptPGP(recipientID, encrypted, recipientPassphrase, nil)
		assert.NoError(t, err)
		assert.Equal(t, data, decrypted)
	})

	// Test decryption with wrong passphrase
	t.Run("decryption with wrong passphrase", func(t *testing.T) {
		encrypted, err := testEncryptPGP(recipientID, data, nil, nil)
		assert.NoError(t, err)

		_, err = testDecryptPGP(recipientID, encrypted, "wrong_passphrase", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decrypt private key")
	})

	// Test signed encryption and decryption
	t.Run("signed encryption and decryption", func(t *testing.T) {
		signerID := "signer@example.com"
		signerPassphrase := "signer_passphrase"
		signerKeyPair, err := CreateKeyPair(signerID, signerPassphrase)
		assert.NoError(t, err)
		signerEntity, err := mockReadEntity(signerKeyPair.PrivateKey)
		assert.NoError(t, err)
		mockKeyStore[signerID] = signerEntity

		encryptedSigned, err := testEncryptPGP(recipientID, data, &signerID, &signerPassphrase)
		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedSigned)

		decryptedSigned, err := testDecryptPGP(recipientID, encryptedSigned, recipientPassphrase, &signerID)
		assert.NoError(t, err)
		assert.Equal(t, data, decryptedSigned)
	})

	// Test signed encryption and decryption with wrong signer ID for verification
	t.Run("signed encryption and decryption with wrong signer verification", func(t *testing.T) {
		signerID := "signer@example.com"
		signerPassphrase := "signer_passphrase"
		signerKeyPair, err := CreateKeyPair(signerID, signerPassphrase)
		assert.NoError(t, err)
		signerEntity, err := mockReadEntity(signerKeyPair.PrivateKey)
		assert.NoError(t, err)
		mockKeyStore[signerID] = signerEntity

		// Create a different signer
		wrongSignerID := "wrong_signer@example.com"
		wrongSignerPassphrase := "wrong_signer_passphrase"
		wrongSignerKeyPair, err := CreateKeyPair(wrongSignerID, wrongSignerPassphrase)
		assert.NoError(t, err)
		wrongSignerEntity, err := mockReadEntity(wrongSignerKeyPair.PrivateKey)
		assert.NoError(t, err)
		mockKeyStore[wrongSignerID] = wrongSignerEntity

		encryptedSigned, err := testEncryptPGP(recipientID, data, &signerID, &signerPassphrase)
		assert.NoError(t, err)

		// Attempt to decrypt and verify with the wrong signer ID
		_, err = testDecryptPGP(recipientID, encryptedSigned, recipientPassphrase, &wrongSignerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature from unexpected key id")
	})

	// Test decryption of unsigned message when signature verification is requested
	t.Run("decryption of unsigned message with signature verification requested", func(t *testing.T) {
		encrypted, err := testEncryptPGP(recipientID, data, nil, nil) // Encrypt without signing
		assert.NoError(t, err)

		signerID := "signer@example.com" // Request verification with a signer
		_, err = testDecryptPGP(recipientID, encrypted, recipientPassphrase, &signerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message is not signed, but signature verification was requested")
	})
}
