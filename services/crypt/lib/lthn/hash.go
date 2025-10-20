package lthn

import (
	"crypto/sha256"
	"encoding/hex"
)

// SetKeyMap sets the key map for the notarisation process.
func SetKeyMap(newKeyMap map[rune]rune) {
	keyMap = newKeyMap
}

// GetKeyMap gets the current key map.
func GetKeyMap() map[rune]rune {
	return keyMap
}

// Hash creates a reproducible hash from a string.
func Hash(input string) string {
	salt := createSalt(input)
	hash := sha256.Sum256([]byte(input + salt))
	return hex.EncodeToString(hash[:])
}

// createSalt creates a quasi-salt from a string by reversing it and swapping characters.
func createSalt(input string) string {
	if input == "" {
		return ""
	}
	runes := []rune(input)
	salt := make([]rune, len(runes))
	for i := 0; i < len(runes); i++ {
		char := runes[len(runes)-1-i]
		if replacement, ok := keyMap[char]; ok {
			salt[i] = replacement
		} else {
			salt[i] = char
		}
	}
	return string(salt)
}

// Verify checks if an input string matches a given hash.
func Verifyf(input string, hash string) bool {
	return Hash(input) == hash
}
