package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	checkOnly bool
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update glow-server to the latest version",
	Run:   runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for updates without installing")
}

func runUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Checking for updates...")

	latestVersion, downloadURL, err := getLatestRelease()
	if err != nil {
		fmt.Printf("Failed to check for updates: %v\n", err)
		os.Exit(1)
	}

	currentVersion := version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	fmt.Printf("Current version: %s\n", currentVersion)
	fmt.Printf("Latest version: %s\n", latestVersion)

	if currentVersion == latestVersion || currentVersion == "dev" {
		fmt.Println("Already up to date!")
		return
	}

	if checkOnly {
		fmt.Println("Update available!")
		return
	}

	fmt.Printf("\nUpdating to %s...\n", latestVersion)

	if err := performUpdate(downloadURL, latestVersion); err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nUpdate successful!")
	fmt.Println("Please restart glow-server to use the new version.")

	// Attempt to restart service
	restartService()
}

func getLatestRelease() (string, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/luaxlou/glow/releases/latest")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	// Construct asset name based on platform
	arch := runtime.GOARCH
	if arch == "arm64" {
		arch = "arm64"
	} else {
		arch = "amd64"
	}

	assetName := fmt.Sprintf("glow-server-%s-%s", runtime.GOOS, arch)
	var downloadURL string

	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, assetName) && !strings.HasSuffix(asset.Name, ".sha256") {
			downloadURL = asset.URL
			break
		}
	}

	if downloadURL == "" {
		return "", "", fmt.Errorf("no binary found for %s-%s", runtime.GOOS, arch)
	}

	return strings.TrimPrefix(release.TagName, "v"), downloadURL, nil
}

func performUpdate(downloadURL, newVersion string) error {
	// Download binary
	tmpDir := os.TempDir()
	tmpBinary := filepath.Join(tmpDir, fmt.Sprintf("glow-server-%s", newVersion))

	// Get download URL (redirect to actual asset)
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Follow redirect manually
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		location := resp.Header.Get("Location")
		resp, err = client.Get(location)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Verify checksum if available
	checksumURL := downloadURL + ".sha256"
	if checksum, err := downloadChecksum(checksumURL); err == nil {
		fmt.Println("Verifying checksum...")
		hasher := sha256.New()
		tmpFile, err := os.Create(tmpBinary)
		if err != nil {
			return err
		}

		multi := io.MultiWriter(hasher, tmpFile)
		if _, err := io.Copy(multi, resp.Body); err != nil {
			tmpFile.Close()
			os.Remove(tmpBinary)
			return err
		}
		tmpFile.Close()

		calculatedHash := hex.EncodeToString(hasher.Sum(nil))
		if calculatedHash != checksum {
			os.Remove(tmpBinary)
			return fmt.Errorf("checksum mismatch: expected %s, got %s", checksum, calculatedHash)
		}
		fmt.Println("Checksum verified successfully")
	} else {
		// No checksum available, download directly
		f, err := os.Create(tmpBinary)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, resp.Body); err != nil {
			return err
		}
	}

	// Make executable
	if err := os.Chmod(tmpBinary, 0755); err != nil {
		return err
	}

	// Create backup of current binary
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	backupPath := exePath + ".backup"
	if err := copyFile(exePath, backupPath); err != nil {
		return err
	}

	// Atomic replacement by renaming
	if err := atomicReplace(tmpBinary, exePath); err != nil {
		// Rollback on failure
		copyFile(backupPath, exePath)
		os.Remove(backupPath)
		return fmt.Errorf("update failed, rolled back: %w", err)
	}

	// Clean up
	os.Remove(backupPath)

	return nil
}

func atomicReplace(src, dst string) error {
	// On Unix systems, rename is atomic
	return os.Rename(src, dst)
}

func downloadChecksum(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		location := resp.Header.Get("Location")
		resp, err = client.Get(location)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checksum download failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func restartService() {
	switch runtime.GOOS {
	case "linux":
		// Check if running as service
		if _, err := os.Stat("/etc/systemd/system/glow-server.service"); err == nil {
			fmt.Println("\nRestarting glow-server service...")
			exec.Command("systemctl", "restart", "glow-server").Run()
		}
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		plistPath := filepath.Join(homeDir, "Library/LaunchAgents/com.luaxlou.glow-server.plist")
		if _, err := os.Stat(plistPath); err == nil {
			fmt.Println("\nRestarting glow-server service...")
			exec.Command("launchctl", "unload", plistPath).Run()
			exec.Command("launchctl", "load", plistPath).Run()
		}
	}
}
