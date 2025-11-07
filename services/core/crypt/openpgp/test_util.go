package openpgp

import (
	"os"
	"path/filepath"
	"testing"
)

// encryptedPrivateKey is a pre-generated, armored PGP private key, encrypted with the passphrase "test-passphrase".
// This key is used in tests where programmatic key generation and encryption is not feasible due to library limitations.
const encryptedPrivateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
hw==
=u7WH
-----END PGP PRIVATE KEY BLOCK-----`

// createEncryptedKeyFile creates a temporary file containing a pre-generated, encrypted private key.
// It returns the path to the temporary file and a cleanup function to remove the temporary directory.
func createEncryptedKeyFile(t *testing.T) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "pgp-test-key-*")
	if err != nil {
		t.Fatalf("test setup: failed to create temp dir for encrypted key: %v", err)
	}

	privKeyPath := filepath.Join(tempDir, "encrypted-key.asc")
	err = os.WriteFile(privKeyPath, []byte(encryptedPrivateKey), 0600)
	if err != nil {
		t.Fatalf("test setup: failed to write encrypted key to file: %v", err)
	}

	cleanup := func() { os.RemoveAll(tempDir) }
	return privKeyPath, cleanup
}
