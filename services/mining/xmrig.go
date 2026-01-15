package mining

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/adrg/xdg"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// NewXMRigMiner creates a new XMRig miner
func NewXMRigMiner() *XMRigMiner {
	return &XMRigMiner{
		Name:    "XMRig",
		Version: "latest",
		URL:     "https://github.com/xmrig/xmrig/releases",
		API: &API{
			Enabled:    true,
			ListenHost: "127.0.0.1",
			ListenPort: 9000,
		},
	}
}

// GetName returns the name of the miner
func (m *XMRigMiner) GetName() string {
	return m.Name
}

// GetLatestVersion returns the latest version of XMRig
func (m *XMRigMiner) GetLatestVersion() (string, error) {
	resp, err := httpClient.Get("https://api.github.com/repos/xmrig/xmrig/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get latest release: unexpected status code %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

// Download and install the latest version of XMRig
func (m *XMRigMiner) Install() error {
	version, err := m.GetLatestVersion()
	if err != nil {
		return err
	}
	m.Version = version

	// Construct the download URL
	var url string
	switch runtime.GOOS {
	case "windows":
		url = fmt.Sprintf("https://github.com/xmrig/xmrig/releases/download/%s/xmrig-%s-msvc-win64.zip", version, strings.TrimPrefix(version, "v"))
	case "linux":
		url = fmt.Sprintf("https://github.com/xmrig/xmrig/releases/download/%s/xmrig-%s-linux-x64.tar.gz", version, strings.TrimPrefix(version, "v"))
	case "darwin":
		url = fmt.Sprintf("https://github.com/xmrig/xmrig/releases/download/%s/xmrig-%s-macos-x64.tar.gz", version, strings.TrimPrefix(version, "v"))
	default:
		return errors.New("unsupported operating system")
	}

	// Create a temporary file to download the release to
	tmpfile, err := os.CreateTemp("", "xmrig-")
	if err != nil {
		return err
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	// Download the release
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download release: unexpected status code %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmpfile, resp.Body); err != nil {
		return err
	}

	// Get the application-specific data path
	dataPath, err := xdg.DataFile("lethean-desktop/miners/xmrig")
	if err != nil {
		return err
	}
	m.Path = dataPath

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(m.Path, 0755); err != nil {
		return err
	}

	// Extract the release
	if strings.HasSuffix(url, ".zip") {
		return m.unzip(tmpfile.Name(), m.Path)
	}
	return m.untar(tmpfile.Name(), m.Path)
}

// Start the miner
func (m *XMRigMiner) Start(config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Running {
		return errors.New("miner is already running")
	}

	// Create the config file
	if err := m.createConfig(config); err != nil {
		return err
	}

	var executableName string
	if runtime.GOOS == "windows" {
		executableName = "xmrig.exe"
	} else {
		executableName = "xmrig"
	}

	executablePath := filepath.Join(m.Path, fmt.Sprintf("xmrig-%s", strings.TrimPrefix(m.Version, "v")), executableName)
	if _, err := os.Stat(executablePath); os.IsNotExist(err) {
		return fmt.Errorf("xmrig executable not found at %s", executablePath)
	}

	m.cmd = exec.Command(executablePath, "-c", m.ConfigPath)
	if err := m.cmd.Start(); err != nil {
		return err
	}

	m.Running = true

	go func() {
		_ = m.cmd.Wait() // Error intentionally ignored; process state tracked via m.Running
		m.mu.Lock()
		m.Running = false
		m.cmd = nil
		m.mu.Unlock()
	}()

	return nil
}

// Stop the miner
func (m *XMRigMiner) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.Running || m.cmd == nil {
		return errors.New("miner is not running")
	}

	// Kill the process. The goroutine in Start() will handle Wait() and state change.
	return m.cmd.Process.Kill()
}

// GetStats returns the stats for the miner
func (m *XMRigMiner) GetStats() (*PerformanceMetrics, error) {
	m.mu.Lock()
	running := m.Running
	m.mu.Unlock()

	if !running {
		return nil, errors.New("miner is not running")
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/2/summary", m.API.ListenHost, m.API.ListenPort))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get stats: unexpected status code %d", resp.StatusCode)
	}

	var summary XMRigSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, err
	}

	var hashrate int
	if len(summary.Hashrate.Total) > 0 {
		hashrate = int(summary.Hashrate.Total[0])
	}

	return &PerformanceMetrics{
		Hashrate:  hashrate,
		Shares:    int(summary.Results.SharesGood),
		Rejected:  int(summary.Results.SharesTotal - summary.Results.SharesGood),
		Uptime:    int(summary.Uptime),
		Algorithm: summary.Algorithm,
	}, nil
}

func (m *XMRigMiner) createConfig(config *Config) error {
	configPath, err := xdg.ConfigFile("lethean-desktop/xmrig.json")
	if err != nil {
		// Fallback to home directory if XDG is not available
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(homeDir, ".config", "lethean-desktop", "xmrig.json")
	}
	m.ConfigPath = configPath

	if err := os.MkdirAll(filepath.Dir(m.ConfigPath), 0755); err != nil {
		return err
	}

	// Create the config
	c := map[string]interface{}{
		"api": map[string]interface{}{
			"enabled":      m.API.Enabled,
			"listen":       fmt.Sprintf("%s:%d", m.API.ListenHost, m.API.ListenPort),
			"access-token": nil,
			"restricted":   true,
		},
		"pools": []map[string]interface{}{
			{
				"url":       config.Pool,
				"user":      config.Wallet,
				"pass":      "x",
				"keepalive": true,
				"tls":       config.TLS,
			},
		},
		"cpu": map[string]interface{}{
			"enabled":    true,
			"threads":    config.Threads,
			"huge-pages": config.HugePages,
		},
	}

	// Write the config to the file
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.ConfigPath, data, 0644)
}

func (m *XMRigMiner) unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		closeErr := outFile.Close()
		_ = rc.Close() // Intentionally ignoring rc.Close error as outFile error is more important
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

func (m *XMRigMiner) untar(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// Sanitize the header name to prevent path traversal
		cleanedName := filepath.Clean(header.Name)
		if strings.HasPrefix(cleanedName, "..") || strings.HasPrefix(cleanedName, "/") || cleanedName == "." {
			continue
		}

		target := filepath.Join(dest, cleanedName)
		rel, err := filepath.Rel(dest, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}

			// manually close here after each file operation; defering would cause each file to wait until all operations have completed.
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
}
