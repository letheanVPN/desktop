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
)

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
client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/xmrig/xmrig/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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
	defer os.Remove(tmpfile.Name())

	// Download the release
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(tmpfile, resp.Body); err != nil {
		return err
	}

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	m.Path = filepath.Join(homeDir, "xmrig")

	// Create the installation directory if it doesn't exist
	if _, err := os.Stat(m.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(m.Path, 0755); err != nil {
			return err
		}
	}

	// Extract the release
	if strings.HasSuffix(url, ".zip") {
		return m.unzip(tmpfile.Name(), m.Path)
	}
	return m.untar(tmpfile.Name(), m.Path)
}


// Start the miner
func (m *XMRigMiner) Start(config *Config) error {
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

	cmd := exec.Command(filepath.Join(m.Path, executableName), "-c", m.ConfigPath)
	if err := cmd.Start(); err != nil {
		return err
	}

	m.Pid = cmd.Process.Pid
	m.Running = true

	go func() {
		cmd.Wait()
		m.Running = false
	}()

	return nil
}

// Stop the miner
func (m *XMRigMiner) Stop() error {
	if !m.Running {
		return errors.New("miner is not running")
	}

	p, err := os.FindProcess(m.Pid)
	if err != nil {
		return err
	}

	// Try to gracefully stop the process
	if err := p.Signal(os.Interrupt); err != nil {
		// Fallback to killing the process if sending Interrupt fails
		if err := p.Kill(); err != nil {
			return err
		}
	}

	// Wait for the process to exit
	_, err = p.Wait()
	if err != nil {
		// If waiting fails, try to kill the process
		if err := p.Kill(); err != nil {
			return err
		}
	}

	m.Running = false
	return nil
}

// GetStats returns the stats for the miner
func (m *XMRigMiner) GetStats() (*PerformanceMetrics, error) {
	if !m.Running {
		return nil, errors.New("miner is not running")
	}

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/2/summary", m.API.ListenHost, m.API.ListenPort))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	m.ConfigPath = filepath.Join(homeDir, ".xmrig.json")

	// Create the config
	c := map[string]interface{}{
		"api": map[string]interface{}{
			"enabled":    m.API.Enabled,
			"listen":     fmt.Sprintf("%s:%d", m.API.ListenHost, m.API.ListenPort),
			"access-token": nil,
			"restricted": true,
		},
		"pools": []map[string]interface{}{
			{
				"url":   config.Pool,
				"user":  config.Wallet,
				"pass":  "x",
				"keepalive": true,
				"tls":     true,
			},
		},
		"cpu": map[string]interface{}{
			"enabled":    true,
			"threads":    config.Threads,
			"huge-pages": true,
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
		// Strip the top-level directory
		parts := strings.Split(f.Name, "/")
		var newName string
		if len(parts) > 1 {
			newName = strings.Join(parts[1:], "/")
		} else {
			newName = parts[0]
		}
		if newName == "" {
			continue
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, newName)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
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
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
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

		// Strip the top-level directory
		parts := strings.Split(header.Name, "/")
		var newName string
		if len(parts) > 1 {
			newName = strings.Join(parts[1:], "/")
		} else {
			newName = parts[0]
		}
		if newName == "" {
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dest, newName)

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
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file to wait until all operations have completed.
			f.Close()
		}
	}
}
