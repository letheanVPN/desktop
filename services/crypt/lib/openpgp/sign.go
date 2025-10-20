package openpgp

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
)

// Sign creates a detached signature for the data.
func Sign(data, privateKeyID, passphrase string) (string, error) {
	signer, err := GetPrivateKey(privateKeyID, passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to get private key for signing: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := openpgp.ArmoredDetachSign(buf, signer, strings.NewReader(data), nil); err != nil {
		return "", fmt.Errorf("failed to create detached signature: %w", err)
	}

	return buf.String(), nil
}

// Verify checks a detached signature.
func Verify(data, signature, publicKeyID string) (bool, error) {
	keyring, err := GetPublicKey(publicKeyID)
	if err != nil {
		return false, fmt.Errorf("failed to get public key for verification: %w", err)
	}

	_, err = openpgp.CheckArmoredDetachedSignature(openpgp.EntityList{keyring}, strings.NewReader(data), strings.NewReader(signature), nil)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %w", err)
	}
	return true, nil
}
