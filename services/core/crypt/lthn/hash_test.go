package lthn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	input := "test_string"
	expectedHash := "45d4027179b17265c38732fb1e7089a0b1adfe1d3ba4105fce66f7d46ba42f7d"

	hashed := Hash(input)
	fmt.Printf("Hash for \"%s\": %s\n", input, hashed)

	assert.Equal(t, expectedHash, hashed, "The hash should match the expected value")
}

func TestCreateSalt(t *testing.T) {
	// Test with default keyMap
	SetKeyMap(map[rune]rune{})
	assert.Equal(t, "gnirts_tset", createSalt("test_string"))
	assert.Equal(t, "", createSalt(""))
	assert.Equal(t, "A", createSalt("A"))

	// Test with a custom keyMap
	customKeyMap := map[rune]rune{
		'a': 'x',
		'b': 'y',
		'c': 'z',
	}
	SetKeyMap(customKeyMap)
	assert.Equal(t, "zyx", createSalt("abc"))
	assert.Equal(t, "gnirts_tset", createSalt("test_string")) // 'test_string' doesn't have 'a', 'b', 'c'

	// Reset keyMap to default for other tests
	SetKeyMap(map[rune]rune{})
}

func TestVerify(t *testing.T) {
	input := "another_test_string"
	hashed := Hash(input)

	assert.True(t, Verifyf(input, hashed), "Verifyf should return true for a matching hash")
	assert.False(t, Verifyf(input, "wrong_hash"), "Verifyf should return false for a non-matching hash")
	assert.False(t, Verifyf("different_input", hashed), "Verifyf should return false for different input")
}
