package openpgp

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

// EncryptPGP encrypts data for a recipient, optionally signing it.
func EncryptPGP(recipientID, data string, signerID, signerPassphrase *string) (string, error) {
	recipient, err := GetPublicKey(recipientID)
	if err != nil {
		return "", fmt.Errorf("failed to get recipient public key: %w", err)
	}

	var signer *openpgp.Entity
	if signerID != nil && signerPassphrase != nil {
		signer, err = GetPrivateKey(*signerID, *signerPassphrase)
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

	// Debug print the encrypted message
	fmt.Printf("Encrypted Message:\n%s\n", buf.String())

	return buf.String(), nil
}

// DecryptPGP decrypts a PGP message, optionally verifying the signature.
func DecryptPGP(recipientID, message, passphrase string, signerID *string) (string, error) {
	privateKeyEntity, err := GetPrivateKey(recipientID, passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	// For this API version, the keyring must contain all keys for decryption and verification.
	keyring := openpgp.EntityList{privateKeyEntity}
	var expectedSigner *openpgp.Entity

	if signerID != nil {
		publicKeyEntity, err := GetPublicKey(*signerID)
		if err != nil {
			return "", fmt.Errorf("could not get public key for verification: %w", err)
		}
		keyring = append(keyring, publicKeyEntity)
		expectedSigner = publicKeyEntity
	}

	// Debug print the message before decryption
	fmt.Printf("Message to Decrypt:\n%s\n", message)

	// We pass the combined keyring, and nil for the prompt function because the private key is already decrypted.
	md, err := openpgp.ReadMessage(strings.NewReader(message), keyring, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read PGP message: %w", err)
	}

	decrypted, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", fmt.Errorf("failed to read decrypted body: %w", err)
	}

	// The signature is checked automatically if the public key is in the keyring.
	// We still need to check for errors and that the signer was who we expected.
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
