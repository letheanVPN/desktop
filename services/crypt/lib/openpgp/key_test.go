package openpgp

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/letheanVPN/desktop/services/filesystem"
)

// MockFileSystemMedium implements the filesystem.Medium interface for testing purposes.
type MockFileSystemMedium struct {
	Files map[string]string
	Dirs  map[string]bool
}

func NewMockFileSystemMedium() *MockFileSystemMedium {
	return &MockFileSystemMedium{
		Files: make(map[string]string),
		Dirs:  make(map[string]bool),
	}
}

func (m *MockFileSystemMedium) Read(path string) (string, error) {
	content, ok := m.Files[path]
	if !ok {
		return "", fmt.Errorf("file not found: %s", path)
	}
	return content, nil
}

func (m *MockFileSystemMedium) Write(path, content string) error {
	m.Files[path] = content
	return nil
}

func (m *MockFileSystemMedium) EnsureDir(path string) error {
	m.Dirs[path] = true
	return nil
}

func (m *MockFileSystemMedium) IsFile(path string) bool {
	_, ok := m.Files[path]
	return ok
}

// saveOriginalLocal and restoreOriginalLocal are helper functions to manage the global filesystem.Local
var originalLocal filesystem.Medium

func saveOriginalLocal() {
	originalLocal = filesystem.Local
}

func restoreOriginalLocal() {
	filesystem.Local = originalLocal
}

func TestCreateKeyPair(t *testing.T) {
	username := "testuser"

	// Call CreateKeyPair without a password
	keyPair, err := CreateKeyPair(username)
	assert.NoError(t, err)
	assert.NotNil(t, keyPair)

	assert.NotEmpty(t, keyPair.PublicKey, "Public key should not be empty")
	assert.NotEmpty(t, keyPair.PrivateKey, "Private key should not be empty")
	assert.NotEmpty(t, keyPair.RevocationCertificate, "Revocation certificate should not be empty")

	// Basic check for PGP armor headers
	assert.Contains(t, keyPair.PublicKey, "-----BEGIN PGP PUBLIC KEY BLOCK-----")
	assert.Contains(t, keyPair.PrivateKey, "-----BEGIN PGP PRIVATE KEY BLOCK-----")
	assert.Contains(t, keyPair.RevocationCertificate, "-----BEGIN PGP SIGNATURE-----")
}

func TestCreateKeyPairPassword(t *testing.T) {
	username := "testuser"
	password := "testpassword"

	keyPair, err := CreateKeyPair(username, password)
	assert.NoError(t, err)
	assert.NotNil(t, keyPair)

	assert.NotEmpty(t, keyPair.PublicKey, "Public key should not be empty")
	assert.NotEmpty(t, keyPair.PrivateKey, "Private key should not be empty")
	assert.NotEmpty(t, keyPair.RevocationCertificate, "Revocation certificate should not be empty")

	// Basic check for PGP armor headers
	assert.Contains(t, keyPair.PublicKey, "-----BEGIN PGP PUBLIC KEY BLOCK-----")
	assert.Contains(t, keyPair.PrivateKey, "-----BEGIN PGP PRIVATE KEY BLOCK-----")
	assert.Contains(t, keyPair.RevocationCertificate, "-----BEGIN PGP SIGNATURE-----")
}

func TestGetPublicKey(t *testing.T) {
	saveOriginalLocal()
	defer restoreOriginalLocal()

	mockFS := NewMockFileSystemMedium()
	filesystem.Local = mockFS // Inject mock filesystem

	username := "testuser_pub"
	password := "testpassword_pub"
	id := username // GetPublicKey uses the ID directly

	// 1. Create a key pair
	keyPair, err := CreateKeyPair(username, password)
	assert.NoError(t, err)

	// 2. Simulate storing the public key in the mock filesystem
	pubPath := filepath.Join("users", fmt.Sprintf("%s.lthn.pub", id))
	mockFS.Write(pubPath, keyPair.PublicKey)

	// 3. Retrieve the public key using GetPublicKey
	entity, err := GetPublicKey(id)
	assert.NoError(t, err)
	assert.NotNil(t, entity)
	// Debug print to inspect the entity's identities
	// fmt.Printf("TestGetPublicKey: Entity Identities for %s: %+v\n", id, entity.Identities)
	assert.Equal(t, username, entity.PrimaryIdentity().UserId.Name) // Changed from .Email

	// Test non-existent public key
	_, err = GetPublicKey("nonexistent_user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestGetPrivateKey(t *testing.T) {
	saveOriginalLocal()
	defer restoreOriginalLocal()

	mockFS := NewMockFileSystemMedium()
	filesystem.Local = mockFS // Inject mock filesystem

	username := "testuser_priv"
	password := "testpassword_priv"
	id := username // GetPrivateKey uses the ID directly

	// 1. Create a key pair
	keyPair, err := CreateKeyPair(username, password)
	assert.NoError(t, err)

	// 2. Simulate storing the private key in the mock filesystem
	privPath := filepath.Join("users", fmt.Sprintf("%s.lthn.key", id))
	mockFS.Write(privPath, keyPair.PrivateKey)

	// 3. Retrieve and decrypt the private key using GetPrivateKey
	entity, err := GetPrivateKey(id, password)
	assert.NoError(t, err)
	assert.NotNil(t, entity)
	// Debug print to inspect the entity's identities
	// fmt.Printf("TestGetPrivateKey: Entity Identities for %s: %+v\n", id, entity.Identities)
	assert.Equal(t, username, entity.PrimaryIdentity().UserId.Name) // Changed from .Email
	assert.False(t, entity.PrivateKey.Encrypted, "Private key should be decrypted")

	// Test non-existent private key
	_, err = GetPrivateKey("nonexistent_user", "any_passphrase")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")

	// Test with wrong passphrase
	_, err = GetPrivateKey(id, "wrong_passphrase")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt private key")
}

func TestCreateAndStoreKeyPair(t *testing.T) {
	saveOriginalLocal()
	defer restoreOriginalLocal()

	mockFS := NewMockFileSystemMedium()
	filesystem.Local = mockFS // Inject mock filesystem

	id := "test_store_user"
	password := "test_store_password"
	keysDir := "test_keys_dir"

	err := createAndStoreKeyPair(id, password, keysDir)
	assert.NoError(t, err)

	// Verify files were written to the mock filesystem
	pubPath := filepath.Join(keysDir, fmt.Sprintf("%s.lthn.pub", id))
	privPath := filepath.Join(keysDir, fmt.Sprintf("%s.lthn.key", id))
	revPath := filepath.Join(keysDir, fmt.Sprintf("%s.lthn.rev", id))

	assert.True(t, mockFS.IsFile(pubPath), "Public key file should exist")
	assert.True(t, mockFS.IsFile(privPath), "Private key file should exist")
	assert.True(t, mockFS.IsFile(revPath), "Revocation certificate file should exist")

	// Verify directory was ensured
	assert.True(t, mockFS.Dirs[keysDir], "Keys directory should have been ensured")

	// Try to read them back to ensure content is valid
	retrievedPubKey, err := mockFS.Read(pubPath)
	assert.NoError(t, err)
	assert.Contains(t, retrievedPubKey, "-----BEGIN PGP PUBLIC KEY BLOCK-----")

	retrievedPrivKey, err := mockFS.Read(privPath)
	assert.NoError(t, err)
	assert.Contains(t, retrievedPrivKey, "-----BEGIN PGP PRIVATE KEY BLOCK-----")
}

func TestCreateServerKeyPair(t *testing.T) {
	saveOriginalLocal()
	defer restoreOriginalLocal()

	mockFS := NewMockFileSystemMedium()
	filesystem.Local = mockFS // Inject mock filesystem

	keysDir := "server_keys"

	err := CreateServerKeyPair(keysDir)
	assert.NoError(t, err)

	// Verify server key files were created
	serverPubPath := filepath.Join(keysDir, "server.lthn.pub")
	serverPrivPath := filepath.Join(keysDir, "server.lthn.key")
	serverRevPath := filepath.Join(keysDir, "server.lthn.rev")

	assert.True(t, mockFS.IsFile(serverPubPath), "Server public key file should exist")
	assert.True(t, mockFS.IsFile(serverPrivPath), "Server private key file should exist")
	assert.True(t, mockFS.IsFile(serverRevPath), "Server revocation certificate file should exist")

	// Verify directory was ensured
	assert.True(t, mockFS.Dirs[keysDir], "Server keys directory should have been ensured")
}
