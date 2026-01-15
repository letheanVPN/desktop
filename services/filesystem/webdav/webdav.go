package webdav

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

// New creates a new, connected instance of the WebDAV storage medium.
func New(cfg ConnectionConfig) (*Medium, error) {
	transport := &authTransport{
		Username: cfg.User,
		Password: cfg.Password,
		Wrapped:  http.DefaultTransport,
	}

	httpClient := &http.Client{Transport: transport}

	// Ping the server to ensure the connection and credentials are valid.
	// We do a PROPFIND on the root, which is a standard WebDAV operation.
	req, err := http.NewRequest("PROPFIND", cfg.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("webdav: failed to create ping request: %w", err)
	}
	req.Header.Set("Depth", "0")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webdav: connection test failed: %w", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusMultiStatus && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("webdav: connection test failed with status %s", resp.Status)
	}

	return &Medium{
		client:  httpClient,
		baseURL: cfg.URL,
	}, nil
}

// Read retrieves the content of a file from the WebDAV server.
func (m *Medium) Read(p string) (string, error) {
	url := m.resolveURL(p)
	resp, err := m.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("webdav: GET request for %s failed: %w", p, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("webdav: failed to read %s, status: %s", p, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("webdav: failed to read response body for %s: %w", p, err)
	}

	return string(data), nil
}

// Write saves the given content to a file on the WebDAV server.
func (m *Medium) Write(p, content string) error {
	// Ensure the parent directory exists first.
	dir := path.Dir(p)
	if dir != "." && dir != "/" {
		if err := m.EnsureDir(dir); err != nil {
			return err // This will be a detailed error from EnsureDir
		}
	}

	url := m.resolveURL(p)
	req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(content)))
	if err != nil {
		return fmt.Errorf("webdav: failed to create PUT request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("webdav: PUT request for %s failed: %w", p, err)
	}
	defer resp.Body.Close()

	// StatusCreated (201) or StatusNoContent (204) are success codes for PUT.
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webdav: failed to write %s, status: %s", p, resp.Status)
	}

	return nil
}

// EnsureDir makes sure a directory exists on the WebDAV server, creating parent dirs as needed.
func (m *Medium) EnsureDir(p string) error {
	// To mimic MkdirAll, we create each part of the path sequentially.
	parts := strings.Split(p, "/")
	currentPath := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		currentPath = path.Join(currentPath, part)
		url := m.resolveURL(currentPath) + "/" // MKCOL needs a trailing slash

		req, err := http.NewRequest("MKCOL", url, nil)
		if err != nil {
			return fmt.Errorf("webdav: failed to create MKCOL request for %s: %w", currentPath, err)
		}

		resp, err := m.client.Do(req)
		if err != nil {
			return fmt.Errorf("webdav: MKCOL request for %s failed: %w", currentPath, err)
		}
		_ = resp.Body.Close()

		// 405 Method Not Allowed means it already exists, which is fine for us.
		// 201 Created is a success.
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMethodNotAllowed {
			return fmt.Errorf("webdav: failed to create directory %s, status: %s", currentPath, resp.Status)
		}
	}
	return nil
}

// IsFile checks if a path exists and is a regular file on the WebDAV server.
func (m *Medium) IsFile(p string) bool {
	url := m.resolveURL(p)
	req, err := http.NewRequest("PROPFIND", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Depth", "0")

	resp, err := m.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// If we get anything other than a Multi-Status, it's probably not a file.
	if resp.StatusCode != http.StatusMultiStatus {
		return false
	}

	// A simple check: if the response body contains the string for a collection, it's a directory.
	// A more robust implementation would parse the XML response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return !strings.Contains(string(body), "<D:collection/>")
}

// resolveURL joins the base URL with a path segment, ensuring correct slashes.
func (m *Medium) resolveURL(p string) string {
	return strings.TrimSuffix(m.baseURL, "/") + "/" + strings.TrimPrefix(p, "/")
}

// authTransport is a custom http.RoundTripper to inject Basic Auth.
type authTransport struct {
	Username string
	Password string
	Wrapped  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Wrapped.RoundTrip(req)
}
