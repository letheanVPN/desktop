package crypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"github.com/letheanVPN/desktop/services/crypt/lib/lthn"
)

func Hash(lib string, payload string) string {
	switch lib {
	case "lthn":
		return lthn.Hash(payload)
	case "sha512":
		hash := sha512.Sum512([]byte(payload))
		return hex.EncodeToString(hash[:])
	case "sha1":
		hash := sha1.Sum([]byte(payload))
		return hex.EncodeToString(hash[:])
	case "md5":
		hash := md5.Sum([]byte(payload))
		return hex.EncodeToString(hash[:])
	case "sha256":
		fallthrough
	default:
		hash := sha256.Sum256([]byte(payload))
		return hex.EncodeToString(hash[:])
	}
}
