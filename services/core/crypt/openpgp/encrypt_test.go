package openpgp

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// generateTestKeys creates a new PGP entity and saves the public and private keys to temporary files.
func generateTestKeys(t *testing.T, name, passphrase string) (string, string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "pgp-keys-*")
	if err != nil {
		t.Fatalf("test setup: failed to create temp dir for keys: %v", err)
	}

	config := &packet.Config{
		RSABits: 2048, // Use a reasonable key size for tests
	}

	entity, err := openpgp.NewEntity(name, "", name, config)
	if err != nil {
		t.Fatalf("test setup: failed to create new PGP entity: %v", err)
	}

	// --- Save Public Key ---
	pubKeyPath := filepath.Join(tempDir, name+".pub")
	pubKeyFile, err := os.Create(pubKeyPath)
	if err != nil {
		t.Fatalf("test setup: failed to create public key file: %v", err)
	}
	pubKeyWriter, err := armor.Encode(pubKeyFile, openpgp.PublicKeyType, nil)
	if err != nil {
		t.Fatalf("test setup: failed to create armored writer for public key: %v", err)
	}
	if err := entity.Serialize(pubKeyWriter); err != nil {
		t.Fatalf("test setup: failed to serialize public key: %v", err)
	}
	if err := pubKeyWriter.Close(); err != nil {
		t.Fatalf("test setup: failed to close public key writer: %v", err)
	}
	if err := pubKeyFile.Close(); err != nil {
		t.Fatalf("test setup: failed to close public key file: %v", err)
	}

	// --- Save Private Key (unencrypted for test setup) ---
	privKeyPath := filepath.Join(tempDir, name+".asc")
	privKeyFile, err := os.Create(privKeyPath)
	if err != nil {
		t.Fatalf("test setup: failed to create private key file: %v", err)
	}
	privKeyWriter, err := armor.Encode(privKeyFile, openpgp.PrivateKeyType, nil)
	if err != nil {
		t.Fatalf("test setup: failed to create armored writer for private key: %v", err)
	}

	// Serialize the whole entity with an unencrypted private key.
	if err := entity.SerializePrivate(privKeyWriter, nil); err != nil {
		t.Fatalf("test setup: failed to serialize private key: %v", err)
	}
	if err := privKeyWriter.Close(); err != nil {
		t.Fatalf("test setup: failed to close private key writer: %v", err)
	}
	if err := privKeyFile.Close(); err != nil {
		t.Fatalf("test setup: failed to close private key file: %v", err)
	}

	cleanup := func() { os.RemoveAll(tempDir) }
	return pubKeyPath, privKeyPath, cleanup
}

func TestEncryptDecryptPGP(t *testing.T) {
	recipientPub, recipientPriv, cleanup := generateTestKeys(t, "recipient", "recipient-pass")
	defer cleanup()

	originalMessage := "This is a secret message."

	// --- Test Encryption ---
	var encryptedBuf bytes.Buffer
	err := EncryptPGP(&encryptedBuf, recipientPub, originalMessage, nil, nil)
	if err != nil {
		t.Fatalf("EncryptPGP() failed unexpectedly: %v", err)
	}
	encryptedMessage := encryptedBuf.String()

	if !strings.Contains(encryptedMessage, "-----BEGIN PGP MESSAGE-----") {
		t.Errorf("Encrypted message does not appear to be PGP armored")
	}

	// --- Test Decryption ---
	decryptedMessage, err := DecryptPGP(recipientPriv, encryptedMessage, "recipient-pass", nil)
	if err != nil {
		t.Fatalf("DecryptPGP() failed unexpectedly: %v", err)
	}

	if decryptedMessage != originalMessage {
		t.Errorf("Decrypted message mismatch: got=%q, want=%q", decryptedMessage, originalMessage)
	}
}

func TestSignAndVerifyPGP(t *testing.T) {
	recipientPub, recipientPriv, rCleanup := generateTestKeys(t, "recipient", "recipient-pass")
	defer rCleanup()

	signerPub, signerPriv, sCleanup := generateTestKeys(t, "signer", "signer-pass")
	defer sCleanup()

	originalMessage := "This is a signed and verified message."

	// --- Encrypt and Sign ---
	var encryptedBuf bytes.Buffer
	signerPass := "signer-pass"
	err := EncryptPGP(&encryptedBuf, recipientPub, originalMessage, &signerPriv, &signerPass)
	if err != nil {
		t.Fatalf("EncryptPGP() with signing failed unexpectedly: %v", err)
	}
	encryptedMessage := encryptedBuf.String()

	// --- Decrypt and Verify ---
	decryptedMessage, err := DecryptPGP(recipientPriv, encryptedMessage, "recipient-pass", &signerPub)
	if err != nil {
		t.Fatalf("DecryptPGP() with verification failed unexpectedly: %v", err)
	}

	if decryptedMessage != originalMessage {
		t.Errorf("Decrypted message mismatch after signing: got=%q, want=%q", decryptedMessage, originalMessage)
	}
}

func TestVerificationFailure(t *testing.T) {
	recipientPub, recipientPriv, rCleanup := generateTestKeys(t, "recipient", "recipient-pass")
	defer rCleanup()

	_, signerPriv, sCleanup := generateTestKeys(t, "signer", "signer-pass")
	defer sCleanup()

	// Generate a third, unexpected key to test verification failure
	unexpectedSignerPub, _, uCleanup := generateTestKeys(t, "unexpected", "unexpected-pass")
	defer uCleanup()

	originalMessage := "This message should fail verification."

	// --- Encrypt and Sign with the actual signer key ---
	var encryptedBuf bytes.Buffer
	signerPass := "signer-pass"
	err := EncryptPGP(&encryptedBuf, recipientPub, originalMessage, &signerPriv, &signerPass)
	if err != nil {
		t.Fatalf("EncryptPGP() with signing failed unexpectedly: %v", err)
	}
	encryptedMessage := encryptedBuf.String()

	// --- Attempt to Decrypt and Verify with the WRONG public key ---
	_, err = DecryptPGP(recipientPriv, encryptedMessage, "recipient-pass", &unexpectedSignerPub)
	if err == nil {
		t.Fatal("DecryptPGP() did not fail, but verification with an incorrect key was expected to fail.")
	}

	if !strings.Contains(err.Error(), "signature from unexpected key") {
		t.Errorf("Expected error to contain 'signature from unexpected key', but got: %v", err)
	}
}
