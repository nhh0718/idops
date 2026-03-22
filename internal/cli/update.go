package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nhh0718/idops/internal/ui"
	"github.com/spf13/cobra"
)

const githubRepo = "nhh0718/idops"

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update idops to the latest version",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for updates...")

	// Fetch latest release info
	release, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check updates: %w", err)
	}

	latest := release.TagName
	current := version
	// Normalize: ensure both have v prefix for comparison
	if !strings.HasPrefix(current, "v") && current != "dev" {
		current = "v" + current
	}

	fmt.Printf("  Current: %s\n", current)
	fmt.Printf("  Latest:  %s\n", latest)

	if current == latest {
		fmt.Println(ui.RenderSuccess("Already up to date!"))
		return nil
	}

	if current == "dev" {
		fmt.Println(ui.RenderWarning("Running dev build, updating to latest release..."))
	}

	// Find matching asset
	assetName := buildAssetName(latest)
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s (looking for %s)", runtime.GOOS, runtime.GOARCH, assetName)
	}

	fmt.Printf("  Downloading %s...\n", assetName)

	// Download to temp file with correct extension
	ext := ".tar.gz"
	if strings.HasSuffix(assetName, ".zip") {
		ext = ".zip"
	}
	tmpPath, err := downloadFile(downloadURL, ext)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpPath)

	// Replace current binary
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary: %w", err)
	}
	execPath, _ = filepath.EvalSymlinks(execPath)

	if err := replaceBinary(tmpPath, execPath); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Println(ui.RenderSuccess(fmt.Sprintf("Updated to %s!", latest)))
	return nil
}

func fetchLatestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func buildAssetName(tag string) string {
	ver := strings.TrimPrefix(tag, "v")
	osName := runtime.GOOS
	arch := runtime.GOARCH
	ext := "tar.gz"
	if osName == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("idops_%s_%s_%s.%s", ver, osName, arch, ext)
}

func downloadFile(url string, ext string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "idops-update-*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func replaceBinary(archivePath, currentPath string) error {
	// Extract full archive (binary + dashboard) to temp dir
	tmpDir, err := os.MkdirTemp("", "idops-extract-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := extractAll(archivePath, tmpDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Find the binary in extracted files
	binaryName := "idops"
	if runtime.GOOS == "windows" {
		binaryName = "idops.exe"
	}
	extractedBinary := filepath.Join(tmpDir, binaryName)
	if _, err := os.Stat(extractedBinary); err != nil {
		return fmt.Errorf("binary not found in archive")
	}

	// Replace binary
	// Windows locks the running exe, so we copy new binary next to it
	// and use a helper script to swap after this process exits.
	if runtime.GOOS == "windows" {
		newPath := currentPath + ".new"
		os.Remove(newPath)
		if err := copyFile(extractedBinary, newPath); err != nil {
			return fmt.Errorf("cannot write new binary: %w", err)
		}
		// Create a batch script that waits, swaps, and cleans up
		batPath := filepath.Join(filepath.Dir(currentPath), "idops-update.bat")
		batContent := fmt.Sprintf("@echo off\r\n"+
			"timeout /t 2 /nobreak >nul\r\n"+
			"del \"%s\"\r\n"+
			"move \"%s\" \"%s\"\r\n"+
			"del \"%s\"\r\n",
			currentPath, newPath, currentPath, batPath)
		if err := os.WriteFile(batPath, []byte(batContent), 0755); err != nil {
			return fmt.Errorf("cannot create update script: %w", err)
		}
		// Launch the batch script detached
		cmd := exec.Command("cmd", "/C", "start", "/b", batPath)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("cannot launch update script: %w", err)
		}
		fmt.Println("  ⏳ Update will complete in a few seconds after exit...")
	} else {
		// Unix: rename is atomic and works on running binaries
		oldPath := currentPath + ".old"
		os.Remove(oldPath)
		if err := os.Rename(currentPath, oldPath); err != nil {
			return fmt.Errorf("cannot rename current binary: %w", err)
		}
		if err := copyFile(extractedBinary, currentPath); err != nil {
			os.Rename(oldPath, currentPath) // rollback
			return err
		}
		os.Chmod(currentPath, 0755)
		os.Remove(oldPath)
	}

	// Copy dashboard if present in archive → ~/.idops/dashboard
	extractedDashboard := filepath.Join(tmpDir, "dashboard")
	if _, err := os.Stat(extractedDashboard); err == nil {
		homeDir, _ := os.UserHomeDir()
		destDashboard := filepath.Join(homeDir, ".idops", "dashboard")

		os.MkdirAll(filepath.Join(homeDir, ".idops"), 0755)
		os.RemoveAll(destDashboard)
		if err := copyDir(extractedDashboard, destDashboard); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: dashboard copy failed: %v\n", err)
		} else {
			fmt.Println("  📦 Dashboard installed to ~/.idops/dashboard")
			// Run npm install if node_modules missing
			nmPath := filepath.Join(destDashboard, "node_modules")
			if _, err := os.Stat(nmPath); err != nil {
				fmt.Println("  📥 Installing dashboard dependencies...")
				npmCmd := exec.Command("npm", "install", "--production")
				npmCmd.Dir = destDashboard
				npmCmd.Stdout = os.Stdout
				npmCmd.Stderr = os.Stderr
				if err := npmCmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: npm install failed: %v\n", err)
				}
			}
		}
	}

	return nil
}

