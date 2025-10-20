package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for testing
	testRoot, err := os.MkdirTemp("", "local_test_root")
	assert.NoError(t, err)
	defer os.RemoveAll(testRoot) // Clean up after the test

	// Test successful creation
	medium, err := New(testRoot)
	assert.NoError(t, err)
	assert.NotNil(t, medium)
	assert.Equal(t, testRoot, medium.root)

	// Verify the root directory exists
	info, err := os.Stat(testRoot)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test creating a new instance with an existing directory (should not error)
	medium2, err := New(testRoot)
	assert.NoError(t, err)
	assert.NotNil(t, medium2)
}

func TestPath(t *testing.T) {
	testRoot := "/tmp/test_root"
	medium := &Medium{root: testRoot}

	// Valid path
	validPath, err := medium.path("file.txt")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testRoot, "file.txt"), validPath)

	// Subdirectory path
	subDirPath, err := medium.path("dir/sub/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testRoot, "dir", "sub", "file.txt"), subDirPath)

	// Path traversal attempt
	_, err = medium.path("../secret.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal attempt detected")

	_, err = medium.path("dir/../../secret.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal attempt detected")
}

func TestReadWrite(t *testing.T) {
	testRoot, err := os.MkdirTemp("", "local_read_write_test")
	assert.NoError(t, err)
	defer os.RemoveAll(testRoot)

	medium, err := New(testRoot)
	assert.NoError(t, err)

	fileName := "testfile.txt"
	filePath := filepath.Join("subdir", fileName)
	content := "Hello, Gopher!\nThis is a test file."

	// Test Write
	err = medium.Write(filePath, content)
	assert.NoError(t, err)

	// Verify file content by reading directly from OS
	readContent, err := os.ReadFile(filepath.Join(testRoot, filePath))
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))

	// Test Read
	readByMedium, err := medium.Read(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, readByMedium)

	// Test Read non-existent file
	_, err = medium.Read("nonexistent.txt")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))

	// Test Write to a path with traversal attempt
	writeErr := medium.Write("../badfile.txt", "malicious content")
	assert.Error(t, writeErr)
	assert.Contains(t, writeErr.Error(), "path traversal attempt detected")
}

func TestEnsureDir(t *testing.T) {
	testRoot, err := os.MkdirTemp("", "local_ensure_dir_test")
	assert.NoError(t, err)
	defer os.RemoveAll(testRoot)

	medium, err := New(testRoot)
	assert.NoError(t, err)

	dirName := "newdir/subdir"
	dirPath := filepath.Join(testRoot, dirName)

	// Test creating a new directory
	err = medium.EnsureDir(dirName)
	assert.NoError(t, err)
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test ensuring an existing directory (should not error)
	err = medium.EnsureDir(dirName)
	assert.NoError(t, err)

	// Test ensuring a directory with path traversal attempt
	err = medium.EnsureDir("../bad_dir")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal attempt detected")
}

func TestIsFile(t *testing.T) {
	testRoot, err := os.MkdirTemp("", "local_is_file_test")
	assert.NoError(t, err)
	defer os.RemoveAll(testRoot)

	medium, err := New(testRoot)
	assert.NoError(t, err)

	// Create a test file
	fileName := "existing_file.txt"
	filePath := filepath.Join(testRoot, fileName)
	err = os.WriteFile(filePath, []byte("content"), 0644)
	assert.NoError(t, err)

	// Create a test directory
	dirName := "existing_dir"
	dirPath := filepath.Join(testRoot, dirName)
	err = os.Mkdir(dirPath, 0755)
	assert.NoError(t, err)

	// Test with an existing file
	assert.True(t, medium.IsFile(fileName))

	// Test with a non-existent file
	assert.False(t, medium.IsFile("nonexistent_file.txt"))

	// Test with a directory
	assert.False(t, medium.IsFile(dirName))

	// Test with path traversal attempt
	assert.False(t, medium.IsFile("../bad_file.txt"))
}
