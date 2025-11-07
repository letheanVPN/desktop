package openpgp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// readRecipientEntity reads an armored PGP public key from the given path.
func readRecipientEntity(path string) (entity *openpgp.Entity, err error) {
	recipientFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to open recipient public key file at %s: %w", path, err)
	}
	defer func() {
		if closeErr := recipientFile.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("openpgp: failed to close recipient key file: %w", closeErr)
		}
	}()

	block, err := armor.Decode(recipientFile)
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to decode armored key from %s: %w", path, err)
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, fmt.Errorf("openpgp: invalid key type in %s: expected public key, got %s", path, block.Type)
	}

	entity, err = openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to read entity from public key: %w", err)
	}
	return entity, nil
}

// readSignerEntity reads and decrypts an armored PGP private key.
func readSignerEntity(path, passphrase string) (entity *openpgp.Entity, err error) {
	signerFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to open signer private key file at %s: %w", path, err)
	}
	defer func() {
		if closeErr := signerFile.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("openpgp: failed to close signer key file: %w", closeErr)
		}
	}()

	block, err := armor.Decode(signerFile)
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to decode armored key from %s: %w", path, err)
	}

	if block.Type != openpgp.PrivateKeyType {
		return nil, fmt.Errorf("openpgp: invalid key type in %s: expected private key, got %s", path, block.Type)
	}

	entity, err = openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to read entity from private key: %w", err)
	}

	// Decrypt the primary private key.
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		if err := entity.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
			return nil, fmt.Errorf("openpgp: failed to decrypt private key: %w", err)
		}
	}

	// Decrypt all subkeys.
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			if err := subkey.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
				return nil, fmt.Errorf("openpgp: failed to decrypt subkey: %w", err)
			}
		}
	}

	return entity, nil
}

// readRecipientKeyRing reads an armored PGP key ring from the given path.
func readRecipientKeyRing(path string) (entityList openpgp.EntityList, err error) {
	recipientFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to open recipient key file at %s: %w", path, err)
	}
	defer func() {
		if closeErr := recipientFile.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("openpgp: failed to close recipient key file: %w", closeErr)
		}
	}()

	entityList, err = openpgp.ReadArmoredKeyRing(recipientFile)
	if err != nil {
		return nil, fmt.Errorf("openpgp: failed to read armored key ring from %s: %w", path, err)
	}
	if len(entityList) == 0 {
		return nil, fmt.Errorf("openpgp: no keys found in recipient key file %s", path)
	}

	return entityList, nil
}

// EncryptPGP encrypts a string using PGP, writing the armored, encrypted
// result to the provided io.Writer.
func EncryptPGP(writer io.Writer, recipientPath, data string, signerPath, signerPassphrase *string) error {
	// 1. Read the recipient's public key
	recipientEntity, err := readRecipientEntity(recipientPath)
	if err != nil {
		return err
	}

	// 2. Set up the list of recipients
	to := openpgp.EntityList{recipientEntity}

	// 3. Handle optional signing
	var signer *openpgp.Entity
	if signerPath != nil {
		var passphrase string
		if signerPassphrase != nil {
			passphrase = *signerPassphrase
		}
		signer, err = readSignerEntity(*signerPath, passphrase)
		if err != nil {
			return fmt.Errorf("openpgp: failed to prepare signer: %w", err)
		}
	}

	// 4. Create an armored writer and encrypt the message
	armoredWriter, err := armor.Encode(writer, "PGP MESSAGE", nil)
	if err != nil {
		return fmt.Errorf("openpgp: failed to create armored writer: %w", err)
	}

	plaintext, err := openpgp.Encrypt(armoredWriter, to, signer, nil, nil)
	if err != nil {
		_ = armoredWriter.Close() // Attempt to close, but prioritize the encryption error.
		return fmt.Errorf("openpgp: failed to begin encryption: %w", err)
	}

	_, err = plaintext.Write([]byte(data))
	if err != nil {
		_ = plaintext.Close()
		_ = armoredWriter.Close()
		return fmt.Errorf("openpgp: failed to write data to encryption stream: %w", err)
	}

	// 5. Explicitly close the writers to finalize the message.
	if err := plaintext.Close(); err != nil {
		return fmt.Errorf("openpgp: failed to finalize plaintext writer: %w", err)
	}
	if err := armoredWriter.Close(); err != nil {
		return fmt.Errorf("openpgp: failed to finalize armored writer: %w", err)
	}

	return nil
}

// DecryptPGP decrypts an armored PGP message.
func DecryptPGP(recipientPath, message, passphrase string, signerPath *string) (string, error) {
	// 1. Read the recipient's private key
	entityList, err := readRecipientKeyRing(recipientPath)
	if err != nil {
		return "", err
	}

	// 2. Decode the armored message
	block, err := armor.Decode(strings.NewReader(message))
	if err != nil {
		return "", fmt.Errorf("openpgp: failed to decode armored message: %w", err)
	}
	if block.Type != "PGP MESSAGE" {
		return "", fmt.Errorf("openpgp: invalid message type: got %s, want PGP MESSAGE", block.Type)
	}

	// 3. If signature verification is required, add signer's public key to keyring
	var signerEntity *openpgp.Entity
	keyring := entityList
	if signerPath != nil {
		signerEntity, err = readRecipientEntity(*signerPath)
		if err != nil {
			return "", fmt.Errorf("openpgp: failed to read signer public key: %w", err)
		}
		keyring = append(keyring, signerEntity)
	}

	// 4. Decrypt the message body
	md, err := openpgp.ReadMessage(block.Body, keyring, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return []byte(passphrase), nil
	}, nil)
	if err != nil {
		return "", fmt.Errorf("openpgp: failed to read PGP message: %w", err)
	}

	// Buffer the unverified body. Do not return or act on it until signature checks pass.
	plaintextBuffer := new(bytes.Buffer)
	if _, err := io.Copy(plaintextBuffer, md.UnverifiedBody); err != nil {
		return "", fmt.Errorf("openpgp: failed to buffer plaintext message body: %w", err)
	}

	// 5. Handle optional signature verification
	if signerPath != nil {
		// First, ensure a signature actually exists when one is expected.
		if md.SignedByKeyId == 0 {
			return "", fmt.Errorf("openpgp: signature verification failed: message is not signed")
		}

		if md.SignatureError != nil {
			return "", fmt.Errorf("openpgp: signature verification failed: %w", md.SignatureError)
		}
		if signerEntity != nil && md.SignedByKeyId != signerEntity.PrimaryKey.KeyId {
			match := false
			for _, subkey := range signerEntity.Subkeys {
				if subkey.PublicKey != nil && subkey.PublicKey.KeyId == md.SignedByKeyId {
					match = true
					break
				}
			}
			if !match {
				return "", fmt.Errorf("openpgp: signature from unexpected key id: got %d, want one of signer key IDs", md.SignedByKeyId)
			}
		}
	}

	return plaintextBuffer.String(), nil
}
