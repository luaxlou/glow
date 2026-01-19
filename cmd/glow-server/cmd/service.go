package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"
)

const systemdTemplate = `[Unit]
Description=Glow Server
After=network.target

[Service]
Type=simple
ExecStart={{.BinaryPath}} serve
Restart=always
User={{.User}}
WorkingDirectory={{.WorkDir}}

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
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.WorkDir}}/server.log</string>
    <key>StandardErrorPath</key>
    <string>{{.WorkDir}}/server.log</string>
    <key>WorkingDirectory</key>
    <string>{{.WorkDir}}</string>
</dict>
</plist>
`

type ServiceConfig struct {
	BinaryPath string
	User       string
	WorkDir    string
}

func installService() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	config := ServiceConfig{
		BinaryPath: exePath,
		User:       os.Getenv("USER"),
		WorkDir:    workDir,
	}

	if config.User == "" {
		config.User = "root"
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

func isServiceInstalled() bool {
	switch runtime.GOOS {
	case "linux":
		_, err := os.Stat("/etc/systemd/system/glow-server.service")
		return err == nil
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		plistPath := fmt.Sprintf("%s/Library/LaunchAgents/com.luaxlou.glow-server.plist", homeDir)
		_, err := os.Stat(plistPath)
		return err == nil
	}
	return false
}
