package e

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE_Good(t *testing.T) {
	err := E("test.op", "test message", assert.AnError)
	assert.Error(t, err)
	assert.Equal(t, "test.op: test message: assert.AnError general error for testing", err.Error())

	err = E("test.op", "test message", nil)
	assert.Error(t, err)
	assert.Equal(t, "test.op: test message", err.Error())
}

func TestE_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := E("test.op", "test message", originalErr)

	assert.True(t, errors.Is(err, originalErr))

	var eErr *Error
	assert.True(t, errors.As(err, &eErr))
	assert.Equal(t, "test.op", eErr.Op)
}
