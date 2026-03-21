package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	// Download to temp file
	tmpPath, err := downloadFile(downloadURL)
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

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "idops-update-*")
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
	// Extract binary from archive to temp dir
	tmpDir, err := os.MkdirTemp("", "idops-extract-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	extractedPath, err := extractBinary(archivePath, tmpDir)
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// On Windows, can't replace running binary directly.
	// Rename current -> .old, then move extracted -> current
	oldPath := currentPath + ".old"
	os.Remove(oldPath)

	if err := os.Rename(currentPath, oldPath); err != nil {
		return fmt.Errorf("cannot rename current binary: %w", err)
	}

	src, err := os.Open(extractedPath)
	if err != nil {
		os.Rename(oldPath, currentPath) // rollback
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(currentPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		os.Rename(oldPath, currentPath)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(currentPath)
		os.Rename(oldPath, currentPath)
		return err
	}

	os.Remove(oldPath)
	return nil
}
