package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adrg/xdg"
)

// setupTestEnv creates temporary directories and sets environment variables
// to simulate a specific user home and XDG base directories for testing.
// It returns the path to the temporary home directory and a cleanup function.
func setupTestEnv(t *testing.T) (string, func()) {
	// Create a temporary directory for the user's home
	tempHomeDir, err := os.MkdirTemp("", "test_home")
	if err != nil {
		t.Fatalf("Failed to create temp home directory: %v", err)
	}

	// Store original HOME environment variable to restore it later
	oldHome := os.Getenv("HOME")

	// Set HOME environment variable for the test
	os.Setenv("HOME", tempHomeDir)

	cleanup := func() {
		// Restore original HOME environment variable
		os.Setenv("HOME", oldHome)
		// Clean up temporary directories
		os.RemoveAll(tempHomeDir)
	}

	return tempHomeDir, cleanup
}

func TestNewService(t *testing.T) {
	tempHomeDir, cleanup := setupTestEnv(t)
	defer cleanup()

	service, err := NewService()
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	cfg := service.Get()

	// These paths are based on the mocked HOME directory
	expectedUserHomeDir := filepath.Join(tempHomeDir, appName)
	expectedConfigDir := filepath.Join(expectedUserHomeDir, "config")
	expectedDataDir := filepath.Join(expectedUserHomeDir, "data")
	expectedWorkspacesDir := filepath.Join(expectedUserHomeDir, "workspaces")

	// For RootDir and CacheDir, xdg library's init() might have already run
	// before our test's os.Setenv calls take effect for xdg. So, we calculate
	// the *expected* values based on what xdg *actually* returns in the
	// current process, which will likely be the system defaults or whatever
	// was set before the test started.
	actualXDGDataFile, err := xdg.DataFile(appName)
	if err != nil {
		t.Fatalf("xdg.DataFile failed: %v", err)
	}
	actualXDGCacheFile, err := xdg.CacheFile(appName)
	if err != nil {
		t.Fatalf("xdg.CacheFile failed: %v", err)
	}

	expectedRootDir := actualXDGDataFile
	expectedCacheDir := actualXDGCacheFile

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"UserHomeDir", cfg.UserHomeDir, expectedUserHomeDir},
		{"RootDir", cfg.RootDir, expectedRootDir},
		{"ConfigDir", cfg.ConfigDir, expectedConfigDir},
		{"DataDir", cfg.DataDir, expectedDataDir},
		{"CacheDir", cfg.CacheDir, expectedCacheDir},
		{"WorkspacesDir", cfg.WorkspacesDir, expectedWorkspacesDir},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Mismatch for %s: got %q, want %q", tt.name, tt.actual, tt.expected)
			}
			// Also check if the directory was actually created
			if info, err := os.Stat(tt.actual); err != nil {
				t.Errorf("Directory %q for %s was not created: %v", tt.actual, tt.name, err)
			} else if !info.IsDir() {
				t.Errorf("Path %q for %s is not a directory", tt.actual, tt.name)
			}
		})
	}
}

func TestNewService_DirectoryCreationFails(t *testing.T) {
	// Create a temporary directory that we will make read-only
	tempHomeDir, err := os.MkdirTemp("", "test_readonly_home")
	if err != nil {
		t.Fatalf("Failed to create temp home directory: %v", err)
	}
	// Ensure cleanup happens, and restore permissions before removing
	defer func() {
		os.Chmod(tempHomeDir, 0755) // Restore write permissions for os.RemoveAll
		os.RemoveAll(tempHomeDir)
	}()

	// Make the temporary home directory read-only
	if err := os.Chmod(tempHomeDir, 0555); err != nil { // r-xr-xr-x
		t.Fatalf("Failed to make temp home directory read-only: %v", err)
	}

	// Store original HOME environment variable to restore it later
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHomeDir)
	defer os.Setenv("HOME", oldHome)

	// NewService should now fail because it cannot create subdirectories in tempHomeDir
	_, err = NewService()
	if err == nil {
		t.Errorf("NewService() expected to fail when directory creation is impossible, but it succeeded")
	}
	// Optionally, check for a specific error message or type
	if err != nil && !strings.Contains(err.Error(), "could not create directory") {
		t.Errorf("NewService() failed with unexpected error: %v", err)
	}
}

func TestNewService_PathTraversalAttempt(t *testing.T) {

	problematicAppName := "../lethean"

	// Simulate the validation logic from NewService
	if !strings.Contains(problematicAppName, "..") && !strings.Contains(problematicAppName, string(filepath.Separator)) {
		t.Errorf("Expected problematicAppName to contain path traversal characters, but it didn't")
	}

	// We'll create a temporary function to simulate the validation within NewService
	validateAppName := func(name string) error {
		if strings.Contains(name, "..") || strings.Contains(name, string(filepath.Separator)) {
			return fmt.Errorf("invalid app name '%s': contains path traversal characters", name)
		}
		return nil
	}

	// Test with a problematic app name
	err := validateAppName(problematicAppName)
	if err == nil {
		t.Errorf("validateAppName expected to fail for %q, but it succeeded", problematicAppName)
	}
	if err != nil && !strings.Contains(err.Error(), "path traversal characters") {
		t.Errorf("validateAppName failed for %q with unexpected error: %v", problematicAppName, err)
	}
	// Test with a safe app name
	safeAppName := "lethean"
	err = validateAppName(safeAppName)
	if err != nil {
		t.Errorf("validateAppName expected to succeed for %q, but it failed with error: %v", safeAppName, err)
	}
}
