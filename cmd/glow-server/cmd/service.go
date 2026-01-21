package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/spf13/cobra"
)

const systemdTemplate = `[Unit]
Description=Glow Server
After=network.target

[Service]
Type=simple
ExecStart={{.BinaryPath}} serve --data-dir={{.DataDir}}
Restart=always
User={{.User}}
WorkingDirectory={{.DataDir}}
EnvironmentFile=-{{.EnvFile}}

[Install]
WantedBy=multi-user.target
`

const launchdTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.luaxlou.glow-server</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.BinaryPath}}</string>
        <string>serve</string>
        <string>--data-dir={{.DataDir}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.DataDir}}/logs/server.log</string>
    <key>StandardErrorPath</key>
    <string>{{.DataDir}}/logs/server.log</string>
    <key>WorkingDirectory</key>
    <string>{{.DataDir}}</string>
    <key>EnvironmentFiles</key>
    <array>
        <string>{{.EnvFile}}</string>
    </array>
</dict>
</plist>
`

type ServiceConfig struct {
	BinaryPath string
	User       string
	DataDir    string
	EnvFile    string
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage glow-server system service",
}

var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install glow-server as a system service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := installService(); err != nil {
			fmt.Printf("Failed to install service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service installed successfully!")
	},
}

var serviceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start glow-server service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startService(); err != nil {
			fmt.Printf("Failed to start service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service started successfully!")
	},
}

var serviceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop glow-server service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := stopService(); err != nil {
			fmt.Printf("Failed to stop service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service stopped successfully!")
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
}

func installService() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Use fixed data directory
	dataDir := "/var/lib/glow-server"
	envFile := "/etc/default/glow-server"

	config := ServiceConfig{
		BinaryPath: exePath,
		User:       os.Getenv("USER"),
		DataDir:    dataDir,
		EnvFile:    envFile,
	}

	if config.User == "" {
		config.User = "root"
	}

	// Create data directory and subdirectories
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	dirs := []string{
		filepath.Join(dataDir, "db"),
		filepath.Join(dataDir, "logs"),
		filepath.Join(dataDir, "apps"),
		filepath.Join(dataDir, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create environment file
	envContent := "# Glow Server Environment Configuration\n"
	envContent += "# PORT=32102\n"
	envContent += "# APP_CENTER_PORT=32101\n"

	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to create environment file: %v", err)
	}

	switch runtime.GOOS {
	case "linux":
		return installSystemd(config)
	case "darwin":
		return installLaunchd(config)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func installSystemd(config ServiceConfig) error {
	servicePath := "/etc/systemd/system/glow-server.service"
	f, err := os.OpenFile(servicePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create service file (try running with sudo): %v", err)
	}
	defer f.Close()

	tmpl, _ := template.New("systemd").Parse(systemdTemplate)
	if err := tmpl.Execute(f, config); err != nil {
		return err
	}

	fmt.Println("Reloading systemd daemon...")
	exec.Command("systemctl", "daemon-reload").Run()
	fmt.Println("Enabling glow-server service...")
	exec.Command("systemctl", "enable", "glow-server").Run()
	fmt.Println("Starting glow-server service...")
	exec.Command("systemctl", "start", "glow-server").Run()

	return nil
}

func installLaunchd(config ServiceConfig) error {
	homeDir, _ := os.UserHomeDir()
	plistPath := fmt.Sprintf("%s/Library/LaunchAgents/com.luaxlou.glow-server.plist", homeDir)

	f, err := os.OpenFile(plistPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %v", err)
	}
	defer f.Close()

	tmpl, _ := template.New("launchd").Parse(launchdTemplate)
	if err := tmpl.Execute(f, config); err != nil {
		return err
	}

	fmt.Println("Loading launchd service...")
	exec.Command("launchctl", "unload", plistPath).Run() // Try unload first if exists
	if err := exec.Command("launchctl", "load", plistPath).Run(); err != nil {
		return fmt.Errorf("failed to load service: %v", err)
	}

	return nil
}

func startService() error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("systemctl", "start", "glow-server").Run()
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		plistPath := filepath.Join(homeDir, "Library/LaunchAgents/com.luaxlou.glow-server.plist")
		return exec.Command("launchctl", "load", plistPath).Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func stopService() error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("systemctl", "stop", "glow-server").Run()
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		plistPath := filepath.Join(homeDir, "Library/LaunchAgents/com.luaxlou.glow-server.plist")
		return exec.Command("launchctl", "unload", plistPath).Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
