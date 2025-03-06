package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ensureGoToolchain ensures the Go toolchain is available.
// If autoGoToolchain is true and Go is not installed, it will download and install Go.
func ensureGoToolchain() error {
	// Skip if auto-obtaining is disabled
	if !autoGoToolchain {
		return nil
	}

	// Check if Go is already installed
	if _, err := exec.LookPath("go"); err == nil {
		return nil // Go is already installed
	}

	// Determine OS and architecture
	goOS := runtime.GOOS
	goArch := runtime.GOARCH

	// Determine latest stable Go version
	goVersion, err := getLatestGoVersion()
	if err != nil {
		return fmt.Errorf("failed to determine latest Go version: %v", err)
	}

	if verbose {
		fmt.Printf("Go toolchain not found. Downloading Go %s for %s/%s\n", goVersion, goOS, goArch)
	}

	// Get download URL for the platform
	downloadURL := fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.tar.gz", goVersion, goOS, goArch)
	if goOS == "windows" {
		downloadURL = fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.zip", goVersion, goOS, goArch)
	}

	// Determine installation directory
	var installDir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	if goOS == "windows" {
		installDir = filepath.Join(homeDir, "go")
	} else {
		installDir = "/usr/local/go"
		// Use home directory if /usr/local/go is not writable
		if err := os.MkdirAll(installDir, 0755); err != nil {
			installDir = filepath.Join(homeDir, "go")
		}
	}

	// Create temporary directory for download
	tempDir, err := os.MkdirTemp("", "go-download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Download Go archive
	archivePath := filepath.Join(tempDir, filepath.Base(downloadURL))
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("failed to download Go: %v", err)
	}

	// Extract archive
	if goOS == "windows" {
		// Extract zip file for Windows
		if err := extractZip(archivePath, filepath.Dir(installDir)); err != nil {
			return fmt.Errorf("failed to extract Go archive: %v", err)
		}
	} else {
		// Extract tar.gz for Unix systems
		if err := extractTarGz(archivePath, filepath.Dir(installDir)); err != nil {
			return fmt.Errorf("failed to extract Go archive: %v", err)
		}
	}

	// Add Go bin to PATH for the current process
	path := os.Getenv("PATH")
	newPath := filepath.Join(installDir, "bin") + string(os.PathListSeparator) + path
	os.Setenv("PATH", newPath)

	if verbose {
		fmt.Printf("Go %s installed to %s\n", goVersion, installDir)
		fmt.Printf("Added %s to PATH\n", filepath.Join(installDir, "bin"))
	}

	// Verify installation
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Go installation verification failed: %v", err)
	}

	if verbose {
		fmt.Printf("Go installation verified: %s\n", string(output))
	}

	return nil
}

// getLatestGoVersion fetches the latest stable Go version
func getLatestGoVersion() (string, error) {
	// For simplicity, we'll use a hardcoded recent version
	// In a real implementation, you might want to fetch this from https://go.dev/dl/?mode=json
	return "1.22.0", nil
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer to track download progress
	var writer io.Writer = out
	if verbose {
		writer = io.MultiWriter(out, newProgressWriter())
	}

	// Write the body to file
	_, err = io.Copy(writer, resp.Body)
	if verbose {
		fmt.Println() // End the progress line
	}
	return err
}

// progressWriter is a simple writer that shows download progress
type progressWriter struct {
	totalBytes int64
	lastReport int64
}

func newProgressWriter() *progressWriter {
	return &progressWriter{}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.totalBytes += int64(n)
	
	// Report progress every 1MB
	if pw.totalBytes-pw.lastReport >= 1024*1024 {
		fmt.Printf("\rDownloading... %d MB", pw.totalBytes/(1024*1024))
		pw.lastReport = pw.totalBytes
	}
	
	return n, nil
}

// extractTarGz extracts a .tar.gz file to a destination directory
func extractTarGz(tarGzFile, destDir string) error {
	// Use tar command for Unix systems
	cmd := exec.Command("tar", "-xzf", tarGzFile, "-C", destDir)
	if verbose {
		fmt.Printf("Extracting %s to %s...\n", tarGzFile, destDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// extractZip extracts a .zip file to a destination directory
func extractZip(zipFile, destDir string) error {
	if verbose {
		fmt.Printf("Extracting %s to %s...\n", zipFile, destDir)
	}
	
	// Use PowerShell for Windows
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", 
			fmt.Sprintf("Expand-Archive -Path %s -DestinationPath %s -Force", zipFile, destDir))
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		return cmd.Run()
	}
	
	// Use unzip for Unix systems that might have it
	cmd := exec.Command("unzip", "-o", zipFile, "-d", destDir)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}