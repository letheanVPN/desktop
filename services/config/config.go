package config

// Config holds the resolved paths for the application.
type Config struct {
	DataDir       string // For user-specific data files.
	ConfigDir     string // For user-specific configuration files.
	CacheDir      string // For non-essential (cached) data.
	WorkspacesDir string // The root directory for all workspaces.
	RootDir       string // The top-level application directory (e.g., ~/Lethean).
	UserHomeDir   string // The user's home directory.
}
