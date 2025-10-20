package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalInitialization(t *testing.T) {
	// The init() function in filesystem.go is automatically called when the package is imported.
	// This test simply verifies that the 'Local' variable is not nil after initialization,
	// implying that the setup of the local filesystem medium was successful and did not
	// result in a fatal error.
	assert.NotNil(t, Local, "filesystem.Local should be initialized after package init()")
}
