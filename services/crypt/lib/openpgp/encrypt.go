package openpgp

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

const (
	pgpMessageHeader = "PGP MESSAGE"
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

	return buf.String(), nil
}

// DecryptPGP decrypts a PGP message, optionally verifying the signature.
func DecryptPGP(recipientID, message, passphrase string, signerID *string) (string, error) {
	privateKeyEntity, err := GetPrivateKey(recipientID, passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	var verificationKeyRing openpgp.EntityList
	if signerID != nil {
		publicKeyEntity, err := GetPublicKey(*signerID)
		if err != nil {
			return "", fmt.Errorf("could not get public key for verification: %w", err)
		}
		verificationKeyRing = openpgp.EntityList{publicKeyEntity}
	}

	md, err := openpgp.ReadMessage(strings.NewReader(message), openpgp.EntityList{privateKeyEntity}, verificationKeyRing, nil)
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
		// Check if the signer is who we expect
		if verificationKeyRing[0].PrimaryKey.KeyId != md.SignedBy.PublicKey.KeyId {
			return "", fmt.Errorf("signature from unexpected key id: got %X, want %X", md.SignedBy.PublicKey.KeyId, verificationKeyRing[0].PrimaryKey.KeyId)
		}
	}

	return string(decrypted), nil
}
