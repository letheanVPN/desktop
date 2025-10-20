package filesystem

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var letheanRoot string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not get user home directory")
	}
	err = SetRoot(filepath.Join(home, "Lethean"))
	if err != nil {
		panic("could not set root directory")
	}
}

// SetRoot sets the root directory for file operations.
// It also ensures the directory exists.
func SetRoot(rootPath string) error {
	letheanRoot = rootPath
	return os.MkdirAll(letheanRoot, os.ModePerm)
}

// GetRoot returns the current root directory.
func GetRoot() string {
	return letheanRoot
}

// Path returns a full, safe path within the lethean root.
func Path(pathname string) (string, error) {
	if strings.Contains(pathname, "..") {
		return "", fmt.Errorf("path traversal attempt detected")
	}
	return filepath.Join(letheanRoot, pathname), nil
}

// Read reads a file from the storage.
func Read(path string) (string, error) {
	safePath, err := Path(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsDir checks if a path is a directory.
func IsDir(path string) bool {
	safePath, err := Path(path)
	if err != nil {
		return false
	}
	info, err := os.Stat(safePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile checks if a path is a file.
func IsFile(path string) bool {
	safePath, err := Path(path)
	if err != nil {
		return false
	}
	info, err := os.Stat(safePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// List lists all files in a directory (non-recursive), ignoring hidden files.
func List(path string) ([]string, error) {
	safePath, err := Path(path)
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(safePath)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			result = append(result, file.Name())
		}
	}
	return result, nil
}

// DetailedList lists all files in a directory with details.
func DetailedList(path string) ([]os.FileInfo, error) {
	safePath, err := Path(path)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, err
	}
	infos := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// Write writes data to a file in the storage.
func Write(path string, data string) error {
	safePath, err := Path(path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(safePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(safePath, []byte(data), 0644)
}

// EnsureDir ensures a directory exists.
func EnsureDir(path string) error {
	safePath, err := Path(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(safePath, os.ModePerm)
}

// Delete deletes a file or directory.
func Delete(path string, recursive bool) error {
	safePath, err := Path(path)
	if err != nil {
		return err
	}

	if safePath == letheanRoot {
		return fmt.Errorf("cannot delete the root storage directory")
	}

	if recursive {
		return os.RemoveAll(safePath)
	}
	return os.Remove(safePath)
}

// --- Download Functionality ---

// DownloadDestination specifies where to save a downloaded file.
type DownloadDestination struct {
	Dir  string
	File string
}

// DownloadedFile holds information about the completed download.
type DownloadedFile struct {
	File     string
	Dir      string
	FullPath string
	Size     int64
}

// ProgressWriter is used to track download progress.
type ProgressWriter struct {
	Total      int64
	Downloaded int64
	File       string
	Dir        string
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += int64(n)
	// Placeholder for sending progress updates (e.g., via ZeroMQ or another event system)
	log.Printf("Download progress for %s: %d / %d bytes", pw.File, pw.Downloaded, pw.Total)
	return n, nil
}

// DownloadContents downloads and optionally extracts a file to the given destination.
func DownloadContents(downloadURL, destDir string, unpack bool) (*DownloadedFile, error) {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	filename := filepath.Base(u.Path)
	safeDestDir, err := Path(destDir)
	if err != nil {
		return nil, fmt.Errorf("invalid destination path: %w", err)
	}

	if err := EnsureDir(safeDestDir); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	log.Printf("Attempting to download %s", downloadURL)
	downloadedFile, err := Download(u, &DownloadDestination{Dir: safeDestDir, File: filename})
	if err != nil {
		return nil, err
	}

	if unpack {
		log.Printf("Extracting %s to: %s", downloadedFile.File, downloadedFile.Dir)
		if err := unarchive(downloadedFile.FullPath, downloadedFile.Dir); err != nil {
			return nil, fmt.Errorf("failed to extract archive: %w", err)
		}
		// Clean up the downloaded archive after extraction
		if err := os.Remove(downloadedFile.FullPath); err != nil {
			log.Printf("Warning: failed to remove archive %s: %v", downloadedFile.FullPath, err)
		}
	}

	return downloadedFile, nil
}

// Download fetches a file from a URL and saves it.
func Download(url *url.URL, destination *DownloadDestination) (*DownloadedFile, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var dir, file string
	if destination == nil || destination.Dir == "" {
		dir, err = os.MkdirTemp("", "leth-dl-")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp dir: %w", err)
		}
	} else {
		dir = destination.Dir
	}

	if destination == nil || destination.File == "" {
		file = filepath.Base(resp.Request.URL.Path)
	} else {
		file = destination.File
	}

	fullPath := filepath.Join(dir, file)
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	progressWriter := &ProgressWriter{
		Total: resp.ContentLength,
		File:  file,
		Dir:   dir,
	}

	reader := io.TeeReader(resp.Body, progressWriter)
	_, err = io.Copy(out, reader)
	if err != nil {
		return nil, err
	}

	return &DownloadedFile{
		File:     file,
		Dir:      dir,
		FullPath: fullPath,
		Size:     progressWriter.Downloaded,
	}, nil
}

// unarchive determines the archive type and extracts it.
func unarchive(src, dest string) error {
	switch {
	case strings.HasSuffix(src, ".zip"):
		return unzip(src, dest)
	case strings.HasSuffix(src, ".tar.gz"):
		return untarGz(src, dest)
	case strings.HasSuffix(src, ".tar.bz2"):
		return untarBz2(src, dest)
	case strings.HasSuffix(src, ".tar"):
		return untar(src, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", src)
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		// Prevent ZipSlip path traversal vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func untar(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	return processTar(file, dest)
}

func untarGz(src, dest string) error {
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
	return processTar(gzr, dest)
}

func untarBz2(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	bz2r := bzip2.NewReader(file)
	return processTar(bz2r, dest)
}

func processTar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dest, header.Name)
		// Prevent path traversal
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", target)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
}
