package manager

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/pkg/api"
)

func TestAppManager_StartApp(t *testing.T) {
	// Setup temp dir
	tmpDir := t.TempDir()

	// Initialize Config
	if err := configmanager.EnsureInitialized(); err != nil {
		t.Fatalf("Failed to init config manager: %v", err)
	}
	configmanager.SetSystemConfig("data_dir", tmpDir)
	configmanager.SetSystemConfig("api_key", "test-key")
	configmanager.SetSystemConfig("server_url", "127.0.0.1:8080")

	// Create dummy binary (shell script)
	dummySrc := filepath.Join(tmpDir, "dummy_app")
	scriptContent := `#!/bin/sh
echo "Starting dummy app"
while true; do
  echo "running"
  sleep 1
done
`
	if err := os.WriteFile(dummySrc, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create dummy app: %v", err)
	}

	// Start App
	req := api.StartAppRequest{
		Name:        "test-app",
		Command:     dummySrc,
		AutoRestart: false,
	}

	if err := StartApp(req); err != nil {
		t.Fatalf("StartApp failed: %v", err)
	}

	// Verify Dir
	appDir := filepath.Join(tmpDir, "apps", "test-app")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		t.Errorf("App dir not created")
	}

	// Verify Binary Renamed
	renamedBin := filepath.Join(appDir, "glow_test-app")
	if _, err := os.Stat(renamedBin); os.IsNotExist(err) {
		t.Errorf("Binary not renamed/copied")
	}

	// Verify Logs Created
	logDir := filepath.Join(appDir, "logs")
	logFile := filepath.Join(logDir, "test-app.log")
	// Allow some time for log creation
	time.Sleep(500 * time.Millisecond)
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file not created")
	}

	// Verify Process Running
	time.Sleep(1 * time.Second)
	apps := ListApps()
	if len(apps) == 0 {
		t.Errorf("No apps listed")
	} else {
		found := false
		for _, app := range apps {
			if app.Name == "test-app" {
				found = true
				if app.Status != "RUNNING" {
					t.Errorf("App status is %s, expected RUNNING", app.Status)
				}
				if app.Pid == 0 {
					t.Errorf("App PID is 0")
				}
			}
		}
		if !found {
			t.Errorf("test-app not found in list")
		}
	}

	// Stop App
	if err := StopApp("test-app"); err != nil {
		t.Fatalf("StopApp failed: %v", err)
	}

	// Verify Stopped
	apps = ListApps()
	for _, app := range apps {
		if app.Name == "test-app" {
			if app.Status != "STOPPED" {
				t.Errorf("App status is %s, expected STOPPED", app.Status)
			}
		}
	}

	// Clean up DB
	os.Remove("glow.db")
}
