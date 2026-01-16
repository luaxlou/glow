package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
}

func runDeploy(cmd *cobra.Command, args []string) {
	binaryPath := args[0]
	name, _ := cmd.Flags().GetString("name")
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
		if existingApp.BinaryHash == localHash {
			fmt.Printf("App '%s' binary is unchanged (hash: %s). No update needed.\n", name, localHash)
			return
		}
		fmt.Printf("App '%s' binary changed (remote: %s, local: %s). Updating...\n", name, existingApp.BinaryHash, localHash)
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
		Name:    name,
		Command: uploadedPath,
	}
	
	// Preserve existing config/env if updating
	if existingApp != nil {
		req.Args = existingApp.Args
		req.Env = existingApp.Env
		req.WorkingDir = existingApp.WorkingDir
		req.AutoRestart = existingApp.AutoRestart
		req.Config = existingApp.Config
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
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.Close()

	client := &http.Client{Timeout: 60 * time.Second} // Longer timeout for upload
	url := strings.TrimSuffix(serverURL, "/") + "/apps/upload"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiResp api.Response
		json.NewDecoder(resp.Body).Decode(&apiResp)
		return "", fmt.Errorf("upload failed: %s (%s)", resp.Status, apiResp.Message)
	}

	var apiResp api.Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	if !apiResp.Success {
		return "", fmt.Errorf("upload failed: %s", apiResp.Message)
	}

	// Data should be the path (string)
	if pathStr, ok := apiResp.Data.(string); ok {
		return pathStr, nil
	}
	return "", fmt.Errorf("invalid response data from upload")
}