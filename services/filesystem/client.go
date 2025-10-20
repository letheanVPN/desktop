package filesystem

import (
	"github.com/letheanVPN/desktop/services/filesystem/sftp"
	"github.com/letheanVPN/desktop/services/filesystem/webdav"
)

// NewSFTPMedium creates and returns a new SFTP medium.
func NewSFTPMedium(cfg sftp.ConnectionConfig) (Medium, error) {
	return sftp.New(cfg)
}

// NewWebDAVMedium creates and returns a new WebDAV medium.
func NewWebDAVMedium(cfg webdav.ConnectionConfig) (Medium, error) {
	return webdav.New(cfg)
}

// Read retrieves the content of a file from the given medium.
func Read(m Medium, path string) (string, error) {
	return m.Read(path)
}

// Write saves content to a file on the given medium.
func Write(m Medium, path, content string) error {
	return m.Write(path, content)
}

// EnsureDir ensures a directory exists on the given medium.
func EnsureDir(m Medium, path string) error {
	return m.EnsureDir(path)
}

// IsFile checks if a path is a file on the given medium.
func IsFile(m Medium, path string) bool {
	return m.IsFile(path)
}

// Copy copies a file from a source medium to a destination medium.
func Copy(sourceMedium Medium, sourcePath string, destMedium Medium, destPath string) error {
	content, err := sourceMedium.Read(sourcePath)
	if err != nil {
		return err
	}
	return destMedium.Write(destPath, content)
}
