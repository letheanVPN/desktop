package webdav

import "net/http"

// Medium implements the filesystem.Medium interface for the WebDAV protocol.
type Medium struct {
	client  *http.Client
	baseURL string // e.g., https://dav.example.com/remote.php/dav/files/username/
}

// ConnectionConfig holds the necessary details to connect to a WebDAV server.
type ConnectionConfig struct {
	URL      string // The full base URL of the WebDAV share.
	User     string
	Password string
}
