package openpgp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDecryptWithWrongPassphrase checks that DecryptPGP returns an error when the wrong passphrase is used.
func TestDecryptWithWrongPassphrase(t *testing.T) {
	recipientPub, _, cleanup := generateTestKeys(t, "recipient", "") // Unencrypted key for encryption
	defer cleanup()

	// Use the pre-generated encrypted key for decryption test
	encryptedPrivKeyPath, cleanup2 := createEncryptedKeyFile(t)
	defer cleanup2()

	originalMessage := "This message should fail to decrypt."

	var encryptedBuf bytes.Buffer
	err := EncryptPGP(&encryptedBuf, recipientPub, originalMessage, nil, nil)
	assert.NoError(t, err, "Encryption failed unexpectedly")
	encryptedMessage := encryptedBuf.String()

	_, err = DecryptPGP(encryptedPrivKeyPath, encryptedMessage, "wrong-passphrase", nil)
	assert.Error(t, err, "Decryption was expected to fail with wrong passphrase, but it succeeded.")
	assert.Contains(t, err.Error(), "failed to read PGP message", "Expected error message about failing to read PGP message")
}

// TestDecryptMalformedMessage checks that DecryptPGP handles non-PGP or malformed input gracefully.
func TestDecryptMalformedMessage(t *testing.T) {
	// Generate an unencrypted key for this test, as we expect failure before key usage.
	_, recipientPriv, cleanup := generateTestKeys(t, "recipient", "")
	defer cleanup()

	malformedMessage := "This is not a PGP message."

	// The passphrase here is irrelevant as the key is not encrypted, but we pass one
	// to satisfy the function signature.
	_, err := DecryptPGP(recipientPriv, malformedMessage, "any-pass", nil)
	assert.Error(t, err, "Decryption should fail for a malformed message, but it did not.")
	assert.Contains(t, err.Error(), "failed to decode armored message", "Expected error about decoding armored message")
}

// TestEncryptWithNonexistentRecipient checks that EncryptPGP fails when the recipient's public key file does not exist.
func TestEncryptWithNonexistentRecipient(t *testing.T) {
	var encryptedBuf bytes.Buffer
	err := EncryptPGP(&encryptedBuf, "/path/to/nonexistent/key.pub", "message", nil, nil)
	assert.Error(t, err, "Encryption should fail if recipient key does not exist, but it succeeded.")
	assert.Contains(t, err.Error(), "failed to open recipient public key file", "Expected file open error for recipient key")
}

// TestEncryptAndSignWithWrongPassphrase checks that signing during encryption fails with an incorrect passphrase.
func TestEncryptAndSignWithWrongPassphrase(t *testing.T) {
	recipientPub, _, rCleanup := generateTestKeys(t, "recipient", "")
	defer rCleanup()

	// Use the pre-generated encrypted key for the signer
	signerPriv, sCleanup := createEncryptedKeyFile(t)
	defer sCleanup()

	originalMessage := "This message should fail to sign."
	wrongPassphrase := "wrong-signer-pass"

	var encryptedBuf bytes.Buffer
	err := EncryptPGP(&encryptedBuf, recipientPub, originalMessage, &signerPriv, &wrongPassphrase)

	assert.Error(t, err, "Encryption with signing was expected to fail with a wrong passphrase, but it succeeded.")
	assert.Contains(t, err.Error(), "failed to decrypt private key", "Expected error about private key decryption failure")
}
