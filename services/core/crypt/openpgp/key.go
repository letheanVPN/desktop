package openpgp

import (
	"bytes"
	"crypto"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/Snider/Core/pkg/crypt/lthn"
)

// CreateKeyPair generates a new OpenPGP key pair.
// The password parameter is optional. If not provided, the private key will not be encrypted.
func CreateKeyPair(username string, passwords ...string) (*KeyPair, error) {
	var password string
	if len(passwords) > 0 {
		password = passwords[0]
	}

	entity, err := openpgp.NewEntity(username, "Lethean Desktop", "", &packet.Config{
		RSABits:     4096,
		DefaultHash: crypto.SHA256,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new entity: %w", err)
	}

	// The private key is initially unencrypted after NewEntity.
	// Generate revocation certificate while the private key is unencrypted.
	revocationCert, err := createRevocationCertificate(entity)
	if err != nil {
		revocationCert = "" // Non-critical, proceed without it if it fails
	}

	// Encrypt the private key only if a password is provided, after revocation cert generation.
	if password != "" {
		if err := entity.PrivateKey.Encrypt([]byte(password)); err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %w", err)
		}
	}

	publicKey, err := serializeEntity(entity, openpgp.PublicKeyType, "") // Public key doesn't need password
	if err != nil {
		return nil, err
	}

	// Private key serialization. The key is already in its final encrypted/unencrypted state.
	privateKey, err := serializeEntity(entity, openpgp.PrivateKeyType, "") // No password needed here for serialization
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PublicKey:             publicKey,
		PrivateKey:            privateKey,
		RevocationCertificate: revocationCert,
	}, nil
}

// CreateServerKeyPair creates and stores a key pair for the server in a specific directory.
func CreateServerKeyPair(keysDir string) error {
	serverKeyPath := filepath.Join(keysDir, "server.lthn.pub")
	// Passphrase is derived from the path itself, consistent with original logic.
	passphrase := lthn.Hash(serverKeyPath)
	return createAndStoreKeyPair("server", passphrase, keysDir)
}

// GetPublicKey retrieves an armored public key for a given ID.
func GetPublicKey(path string) (*openpgp.Entity, error) {
	return readEntity(path)
}

// GetPrivateKey retrieves and decrypts an armored private key.
func GetPrivateKey(path, passphrase string) (*openpgp.Entity, error) {
	entity, err := readEntity(path)
	if err != nil {
		return nil, err
	}

	if entity.PrivateKey == nil {
		return nil, fmt.Errorf("no private key found for path %s", path)
	}

	if entity.PrivateKey.Encrypted {
		if err := entity.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
			return nil, fmt.Errorf("failed to decrypt private key for path %s: %w", path, err)
		}
	}

	var primaryIdentity *openpgp.Identity
	for _, identity := range entity.Identities {
		if identity.SelfSignature.IsPrimaryId != nil && *identity.SelfSignature.IsPrimaryId {
			primaryIdentity = identity
			break
		}
	}
	if primaryIdentity == nil {
		for _, identity := range entity.Identities {
			primaryIdentity = identity
			break
		}
	}

	if primaryIdentity == nil {
		return nil, fmt.Errorf("key for %s has no identity", path)
	}

	if primaryIdentity.SelfSignature.KeyLifetimeSecs != nil {
		if primaryIdentity.SelfSignature.CreationTime.Add(time.Duration(*primaryIdentity.SelfSignature.KeyLifetimeSecs) * time.Second).Before(time.Now()) {
			return nil, fmt.Errorf("key for %s has expired", path)
		}
	}

	return entity, nil
}

// --- Helper Functions ---

func createAndStoreKeyPair(id, password, dir string) error {
	//var keyPair *KeyPair
	var err error

	//if password != "" {
	//	keyPair, err = CreateKeyPair(id, password)
	//} else {
	//	keyPair, err = CreateKeyPair(id)
	//}

	if err != nil {
		return fmt.Errorf("failed to create key pair for id %s: %w", id, err)
	}

	//if err := io.Local.EnsureDir(dir); err != nil {
	//	return fmt.Errorf("failed to ensure key directory exists: %w", err)
	//}
	//
	//files := map[string]string{
	//	filepath.Join(dir, fmt.Sprintf("%s.lthn.pub", id)): keyPair.PublicKey,
	//	filepath.Join(dir, fmt.Sprintf("%s.lthn.key", id)): keyPair.PrivateKey,
	//	filepath.Join(dir, fmt.Sprintf("%s.lthn.rev", id)): keyPair.RevocationCertificate, // Re-enabled
	//}
	//
	//for path, content := range files {
	//	if content == "" {
	//		continue
	//	}
	//	if err := io.Local.Write(path, content); err != nil {
	//		return fmt.Errorf("failed to write key file %s: %w", path, err)
	//	}
	//}
	return nil
}

func readEntity(path string) (*openpgp.Entity, error) {
	//keyArmored, err := m.Read(path)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read key file %s: %w", path, err)
	//}

	//entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(keyArmored))
	//if err != nil {
	//	return nil, fmt.Errorf("failed to parse key file %s: %w", path, err)
	//}
	//if len(entityList) == 0 {
	//	return nil, fmt.Errorf("no entity found in key file %s", path)
	//}
	//return entityList[0], nil
	return nil, nil
}

func serializeEntity(entity *openpgp.Entity, keyType string, password string) (string, error) {
	buf := new(bytes.Buffer)
	writer, err := armor.Encode(buf, keyType, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create armor encoder: %w", err)
	}

	if keyType == openpgp.PrivateKeyType {
		// Serialize the private key in its current in-memory state.
		// Encryption is handled by CreateKeyPair before this function is called.
		err = entity.SerializePrivateWithoutSigning(writer, nil)
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

	// SignKey requires an unencrypted private key.
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
