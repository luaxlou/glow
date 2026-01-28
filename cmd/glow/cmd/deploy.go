package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

func init() {
	// Configure the existing deployCmd defined in verbs.go
	deployCmd.Use = "deploy [binary]"
	deployCmd.Short = "Deploy or update an application"
	deployCmd.Args = cobra.ExactArgs(1)
	deployCmd.Run = runDeploy

	deployCmd.Flags().String("name", "", "Name of the application (defaults to binary name)")
	deployCmd.Flags().BoolP("force", "f", false, "Force update even if binary is unchanged")
}

func runDeploy(cmd *cobra.Command, args []string) {
	binaryPath := args[0]
	name, _ := cmd.Flags().GetString("name")
	force, _ := cmd.Flags().GetBool("force")
	if name == "" {
		name = filepath.Base(binaryPath)
	}

	// 1. Calculate local hash
	localHash, err := calculateLocalFileHash(binaryPath)
	if err != nil {
		fmt.Printf("Error calculating local hash: %v\n", err)
		return
	}

	// 2. Check remote hash
	var resp api.Response
	if err := request("GET", "/apps/list", nil, &resp); err != nil {
		fmt.Printf("Error getting app list: %v\n", err)
		return
	}

	// Parse list
	data, _ := json.Marshal(resp.Data)
	var apps []api.AppInfo
	json.Unmarshal(data, &apps)

	var existingApp *api.AppInfo
	for _, app := range apps {
		if app.Name == name {
			existingApp = &app
			break
		}
	}

	if existingApp != nil {
		if existingApp.BinaryHash == localHash && !force {
			fmt.Printf("App '%s' binary is unchanged (hash: %s). No update needed.\n", name, localHash)
			return
		}
		if force {
			fmt.Printf("Forcing update for app '%s'...\n", name)
		} else {
			fmt.Printf("App '%s' binary changed (remote: %s, local: %s). Updating...\n", name, existingApp.BinaryHash, localHash)
		}
	} else {
		fmt.Printf("Deploying new app '%s'...\n", name)
	}

	// 3. Upload binary
	fmt.Printf("Uploading %s...\n", binaryPath)
	uploadedPath, err := uploadFile(binaryPath)
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		return
	}

	// 4. Start/Update App
	req := api.StartAppRequest{
		Name:        name,
		Command:     uploadedPath,
		SkipIngress: true,
	}

	// Preserve existing config/env if updating
	if existingApp != nil {
		req.Args = existingApp.Args
		req.Env = existingApp.Env
		req.WorkingDir = existingApp.WorkingDir
		req.AutoRestart = existingApp.AutoRestart
		req.Config = existingApp.Config

		// Stop the app first to ensure restart with new binary
		if existingApp.Status == "RUNNING" {
			fmt.Printf("Stopping app '%s' for update...\n", name)
			var stopResp api.Response
			stopReq := api.StopAppRequest{Name: name, KeepIngress: true}
			if err := request("POST", "/apps/stop", stopReq, &stopResp); err != nil {
				fmt.Printf("Warning: Failed to stop app: %v\n", err)
			}
		}
	}

	var startResp api.Response
	if err := request("POST", "/apps/start", req, &startResp); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		return
	}

	if startResp.Success {
		fmt.Printf("App '%s' deployed successfully.\n", name)
	} else {
		fmt.Printf("Error deploying app: %s\n", startResp.Message)
	}
}

func calculateLocalFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func uploadFile(path string) (string, error) {
	uploadURL := strings.TrimSuffix(sanitizeServerURL(serverURL), "/") + "/apps/upload"
	key := strings.TrimSpace(apiKey)
	if key == "" {
		return "", fmt.Errorf("missing api key")
	}

	        const maxAttempts = 3
	        timeout := 60 * time.Minute
	
	        var lastErr error
	        for attempt := 1; attempt <= maxAttempts; attempt++ {
		uploadedPath, err := uploadFileOnce(uploadURL, key, path, timeout)
		if err == nil {
			return uploadedPath, nil
		}
		lastErr = err
		if !isRetryableUploadError(err) || attempt == maxAttempts {
			break
		}
		time.Sleep(time.Duration(500*(1<<(attempt-1))) * time.Millisecond)
	}
	return "", lastErr
}

func uploadFileOnce(uploadURL, apiKey, filePath string, timeout time.Duration) (string, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	writeErrCh := make(chan error, 1)

	req, err := http.NewRequest("POST", uploadURL, pr)
	if err != nil {
		_ = pw.CloseWithError(err)
		return "", err
	}

	go func() {
		file, err := os.Open(filePath)
		if err != nil {
			_ = pw.CloseWithError(err)
			writeErrCh <- err
			return
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			_ = pw.CloseWithError(err)
			writeErrCh <- err
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			_ = pw.CloseWithError(err)
			writeErrCh <- err
			return
		}

		if err := writer.Close(); err != nil {
			_ = pw.CloseWithError(err)
			writeErrCh <- err
			return
		}

		_ = pw.Close()
		writeErrCh <- nil
	}()

	client := &http.Client{Timeout: timeout}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		_ = pw.CloseWithError(err)
		_ = <-writeErrCh
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiResp api.Response
		_ = json.NewDecoder(resp.Body).Decode(&apiResp)
		<-writeErrCh
		return "", fmt.Errorf("upload failed: %s (%s)", resp.Status, apiResp.Message)
	}

	var apiResp api.Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		<-writeErrCh
		return "", err
	}

	if err := <-writeErrCh; err != nil {
		return "", err
	}

	if !apiResp.Success {
		return "", fmt.Errorf("upload failed: %s", apiResp.Message)
	}

	if pathStr, ok := apiResp.Data.(string); ok {
		return pathStr, nil
	}
	return "", fmt.Errorf("invalid response data from upload")
}

func sanitizeServerURL(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.Trim(s, "`")
	return s
}

func isRetryableUploadError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "Client.Timeout exceeded") ||
		strings.Contains(msg, "context deadline exceeded") ||
		strings.Contains(msg, "timeout")
}
