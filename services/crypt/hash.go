package crypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"github.com/letheanVPN/desktop/services/crypt/lib/lthn"
)

// Hash computes a hash of the payload using the specified algorithm.
func (s *Service) Hash(lib HashType, payload string) string {
	switch lib {
	case LTHN:
		return lthn.Hash(payload)
	case SHA512:
		hash := sha512.Sum512([]byte(payload))
		return hex.EncodeToString(hash[:])
	case SHA1:
		hash := sha1.Sum([]byte(payload))
		return hex.EncodeToString(hash[:])
	case MD5:
		hash := md5.Sum([]byte(payload))
		return hex.EncodeToString(hash[:])
	case SHA256:
		fallthrough
	default:
		hash := sha256.Sum256([]byte(payload))
		return hex.EncodeToString(hash[:])
	}
}
