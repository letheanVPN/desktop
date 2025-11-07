package crypt

import (
	"github.com/Snider/Core/pkg/core"
	"github.com/Snider/Core/pkg/crypt/internal"
)

// Options holds configuration for the crypt service.
type Options = internal.Options

// Service provides cryptographic functions to the application.
type Service = internal.Service

// HashType defines the supported hashing algorithms.
type HashType = internal.HashType

const (
	LTHN   = internal.LTHN
	SHA512 = internal.SHA512
	SHA256 = internal.SHA256
	SHA1   = internal.SHA1
	MD5    = internal.MD5
)

// New is the constructor for static dependency injection.
func New() (*Service, error) {
	return internal.New()
}

// Register is the constructor for dynamic dependency injection.
func Register(c *core.Core) (any, error) {
	return internal.Register(c)
}
