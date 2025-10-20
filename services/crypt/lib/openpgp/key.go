package openpgp

import (
	"bytes"
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/letheanVPN/desktop/services/crypt/lib/lthn"
	"github.com/letheanVPN/desktop/services/filesystem"
)

const (
	userDirectory = "users"
)

// KeyPair holds the generated armored keys and revocation certificate.
type KeyPair struct {
	PublicKey             string
	PrivateKey            string
	RevocationCertificate string
}

// CreateKeyPair generates a new OpenPGP key pair.
func CreateKeyPair(username, password string) (*KeyPair, error) {
	entity, err := openpgp.NewEntity(username, "Lethean Desktop", "", &packet.Config{
		RSABits:     4096,
		DefaultHash: crypto.SHA256,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new entity: %w", err)
	}

	// Encrypt the private key with the given passphrase.
	if err := entity.PrivateKey.Encrypt([]byte(password)); err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Serialize public key
	publicKey, err := serializeEntity(entity, openpgp.PublicKeyType)
	if err != nil {
		return nil, err
	}

	// Serialize private key
	privateKey, err := serializeEntity(entity, openpgp.PrivateKeyType)
	if err != nil {
		return nil, err
	}

	// Create and serialize revocation certificate
	revocationCert, err := createRevocationCertificate(entity)
	if err != nil {
		// This is not a fatal error, but we should log it.
		revocationCert = ""
	}

	return &KeyPair{
		PublicKey:             publicKey,
		PrivateKey:            privateKey,
		RevocationCertificate: revocationCert,
	}, nil
}

// GetPublicKey retrieves an armored public key for a given ID.
func GetPublicKey(id string) (*openpgp.Entity, error) {
	path := fmt.Sprintf("%s/%s.lthn.pub", userDirectory, id)
	return readEntity(path)
}

// GetPrivateKey retrieves and decrypts an armored private key.
func GetPrivateKey(id, passphrase string) (*openpgp.Entity, error) {
	path := fmt.Sprintf("%s/%s.lthn.key", userDirectory, id)
	entity, err := readEntity(path)
	if err != nil {
		return nil, err
	}

	if entity.PrivateKey == nil {
		return nil, fmt.Errorf("no private key found for id %s", id)
	}

	if entity.PrivateKey.Encrypted {
		if err := entity.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
			return nil, fmt.Errorf("failed to decrypt private key for id %s: %w", id, err)
		}
	}

	// Check for key expiry
	var primaryIdentity *openpgp.Identity
	for _, identity := range entity.Identities {
		if identity.SelfSignature.IsPrimaryId != nil && *identity.SelfSignature.IsPrimaryId {
			primaryIdentity = identity
			break
		}
	}
	if primaryIdentity == nil {
		for _, identity := range entity.Identities {
			primaryIdentity = identity // fallback to first identity
			break
		}
	}

	if primaryIdentity == nil {
		return nil, fmt.Errorf("key for %s has no identity", id)
	}

	if primaryIdentity.SelfSignature.KeyLifetimeSecs != nil {
		if primaryIdentity.SelfSignature.CreationTime.Add(time.Duration(*primaryIdentity.SelfSignature.KeyLifetimeSecs) * time.Second).Before(time.Now()) {
			return nil, fmt.Errorf("key for %s has expired", id)
		}
	}

	return entity, nil
}

// CreateUserKeyPair creates and stores a key pair for a user.
func CreateUserKeyPair(username, password string) error {
	usernameHash := lthn.Hash(username)
	return createAndStoreKeyPair(usernameHash, password, userDirectory)
}

// CreateServerKeyPair creates and stores a key pair for the server.
func CreateServerKeyPair() error {
	pubKeyPath := fmt.Sprintf("%s/server.lthn.pub", userDirectory)
	fullPath, err := filesystem.Path(pubKeyPath)
	if err != nil {
		return fmt.Errorf("could not resolve server public key path: %w", err)
	}
	passphrase := lthn.Hash(fullPath)
	return createAndStoreKeyPair("server", passphrase, userDirectory)
}

// --- Helper Functions ---

// createAndStoreKeyPair generates keys and writes them to the filesystem.
func createAndStoreKeyPair(id, password, dir string) error {
	keyPair, err := CreateKeyPair(id, password)
	if err != nil {
		return fmt.Errorf("failed to create key pair for id %s: %w", id, err)
	}

	if err := filesystem.EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to ensure user directory exists: %w", err)
	}

	files := map[string]string{
		fmt.Sprintf("%s/%s.lthn.pub", dir, id): keyPair.PublicKey,
		fmt.Sprintf("%s/%s.lthn.key", dir, id): keyPair.PrivateKey,
		fmt.Sprintf("%s/%s.lthn.rev", dir, id): keyPair.RevocationCertificate,
	}

	for path, content := range files {
		if content == "" {
			continue // Don't write empty revocation certs
		}
		if err := filesystem.Write(path, content); err != nil {
			return fmt.Errorf("failed to write key file %s: %w", path, err)
		}
	}
	return nil
}

// readEntity reads an armored PGP entity (key) from the filesystem.
func readEntity(path string) (*openpgp.Entity, error) {
	keyArmored, err := filesystem.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file %s: %w", path, err)
	}

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(keyArmored))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key file %s: %w", path, err)
	}
	if len(entityList) == 0 {
		return nil, fmt.Errorf("no entity found in key file %s", path)
	}
	return entityList[0], nil
}

// serializeEntity armors and serializes a PGP entity.
func serializeEntity(entity *openpgp.Entity, keyType string) (string, error) {
	buf := new(bytes.Buffer)
	writer, err := armor.Encode(buf, keyType, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create armor encoder: %w", err)
	}

	if keyType == openpgp.PrivateKeyType {
		err = entity.SerializePrivate(writer, nil)
	} else {
		err = entity.Serialize(writer)
	}

	if err != nil {
		return "", fmt.Errorf("failed to serialize entity: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close armor writer: %w", err)
	}
	return buf.String(), nil
}

// createRevocationCertificate generates an armored revocation certificate for an entity.
func createRevocationCertificate(entity *openpgp.Entity) (string, error) {
	buf := new(bytes.Buffer)
	writer, err := armor.Encode(buf, openpgp.SignatureType, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create armor encoder for revocation: %w", err)
	}

	sig := &packet.Signature{
		SigType:      packet.SigTypeKeyRevocation,
		PubKeyAlgo:   entity.PrimaryKey.PubKeyAlgo,
		Hash:         crypto.SHA256,
		CreationTime: time.Now(),
		IssuerKeyId:  &entity.PrimaryKey.KeyId,
	}

	if err := sig.SignKey(entity.PrimaryKey, entity.PrivateKey, nil); err != nil {
		return "", fmt.Errorf("failed to sign revocation: %w", err)
	}
	if err := sig.Serialize(writer); err != nil {
		return "", fmt.Errorf("failed to serialize revocation signature: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close revocation writer: %w", err)
	}
	return buf.String(), nil
}
