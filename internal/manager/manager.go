package manager

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/luaxlou/glow/internal/appcenter"
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/internal/provisioner"
	"github.com/luaxlou/glow/internal/statemanager"
	"github.com/luaxlou/glow/pkg/api"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	mu sync.RWMutex
)

func ProvisionResource(req api.ProvisionRequest) (map[string]any, error) {
	if req.ResourceType == "mysql" {
		var mysqlConfig api.MySQLConfig
		if err := configmanager.GetSystemConfigJSON("mysql_info", &mysqlConfig); err != nil {
			return nil, fmt.Errorf("mysql info not found: %w", err)
		}
		if mysqlConfig.Host == "" {
			return nil, fmt.Errorf("mysql service not configured")
		}

		p := provisioner.NewMySQL(api.ServiceSpec{
			Port:          mysqlConfig.Port,
			AdminUser:     mysqlConfig.User,
			AdminPassword: mysqlConfig.Password,
		})
		if err := p.Check(); err != nil {
			return nil, fmt.Errorf("mysql check failed: %w", err)
		}

		user, pass, err := p.Provision(req.ResourceName)
		if err != nil {
			return nil, fmt.Errorf("failed to provision mysql: %w", err)
		}

		configFragment := map[string]any{
			"mysql": map[string]interface{}{
				"dsn": fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s?parseTime=true", user, pass, mysqlConfig.Port, req.ResourceName),
			},
		}

		if err := configmanager.Set(req.AppName, configFragment, true); err != nil {
			return nil, fmt.Errorf("failed to save config: %w", err)
		}

		return configFragment, nil
	}

	return nil, fmt.Errorf("unsupported resource type: %s", req.ResourceType)
}

func StartApp(req api.StartAppRequest) error {
	mu.Lock()
	defer mu.Unlock()

	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	serverURL, _ := configmanager.GetSystemConfig("server_url")
	apiKey, _ := configmanager.GetSystemConfig("api_key")

	var existingRestartCount int
	var existingDomain string
	var existingConfigHash string
	var existingBinaryHash string

	if existingApp, err := statemanager.GetApp(req.Name); err == nil {
		if existingApp.Status == "RUNNING" {
			// Idempotent: return success if already running
			return nil
		}
		if existingApp.Status != "STOPPED" {
			existingRestartCount = existingApp.RestartCount
		}
		existingDomain = existingApp.Domain
		existingConfigHash = existingApp.ConfigHash
		existingBinaryHash = existingApp.BinaryHash
	}

	port, err := GetFreePort()
	// Optional port allocation: if not specified in config, try to allocate
	// But spec says "Make port allocation optional (read from config)".
	// If OP_APP_PORT env or config exists, use it?
	// For now, let's keep allocation but maybe skip if config says so?
	// Simplified: Always allocate for now as per "system MUST allocate" in old spec,
	// but new spec says "if not specified".
	// Let's check if config has port.

	// Check if app config has port
	if cfgVal, ok := configmanager.GetValue(req.Name, "port"); ok {
		if p, ok := cfgVal.(float64); ok { // JSON numbers are float64
			port = int(p)
			err = nil
		} else if p, ok := cfgVal.(int); ok {
			port = p
			err = nil
		}
	}

	if err != nil {
		return fmt.Errorf("failed to assign port: %w", err)
	}

	// 1. Prepare Directory: apps/<name>
	appDir := filepath.Join(dataDir, "apps", req.Name)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app dir: %w", err)
	}

	// 2. Rename Binary: glow_<name>
	srcBinary := req.Command
	dstBinaryName := "glow_" + req.Name
	dstBinaryPath := filepath.Join(appDir, dstBinaryName)

	// Attempt to copy if source exists
	if _, err := os.Stat(srcBinary); err == nil {
		if err := copyFile(srcBinary, dstBinaryPath); err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}
		if err := os.Chmod(dstBinaryPath, 0755); err != nil {
			return fmt.Errorf("failed to chmod binary: %w", err)
		}
	} else {
		// If source not found, check if destination already exists (restart scenario)
		if _, err := os.Stat(dstBinaryPath); err != nil {
			// If neither exists, fallback if it looks like a command
			if strings.Contains(req.Command, string(os.PathSeparator)) {
				return fmt.Errorf("binary not found: %s", req.Command)
			}
			dstBinaryPath = req.Command
		}
	}

	// Calculate Binary Hash
	currentBinaryHash, err := calculateFileHash(dstBinaryPath)
	if err == nil {
		existingBinaryHash = currentBinaryHash
	}

	app := api.AppInfo{
		Name:         req.Name,
		Command:      dstBinaryPath,
		Args:         req.Args,
		WorkingDir:   req.WorkingDir,
		Env:          req.Env,
		Port:         port,
		Domain:       existingDomain,
		AutoRestart:  req.AutoRestart,
		RestartCount: existingRestartCount,
		Status:       "STARTING",
		ConfigHash:   existingConfigHash,
		BinaryHash:   existingBinaryHash,
	}

	if app.WorkingDir == "" {
		app.WorkingDir = appDir
	}

	if err := GenerateNginxConfig(dataDir, NginxConfig{
		Name:   app.Name,
		Port:   app.Port,
		Domain: app.Domain,
	}); err != nil {
		fmt.Printf("Warning: Failed to generate nginx config: %v\n", err)
	}

	cmd := exec.Command(app.Command, app.Args...)
	cmd.Dir = app.WorkingDir

	cmd.Env = os.Environ()
	for k, v := range app.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("OP_APP_PORT=%d", port))
	cmd.Env = append(cmd.Env, fmt.Sprintf("OP_SERVER_URL=%s", serverURL))
	cmd.Env = append(cmd.Env, fmt.Sprintf("OP_APP_NAME=%s", app.Name))
	cmd.Env = append(cmd.Env, fmt.Sprintf("OP_API_KEY=%s", apiKey))

	// 3. Run as glow user
	if u, err := user.Lookup("glow"); err == nil {
		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
		}
	}

	// 4. Logs with rotation
	logDir := filepath.Join(appDir, "logs")
	os.MkdirAll(logDir, 0755)
	logFile := filepath.Join(logDir, app.Name+".log")

	rotator, err := NewLogRotator(logFile, 10*1024*1024, 5)
	if err == nil {
		cmd.Stdout = rotator
		cmd.Stderr = rotator
	} else {
		f, _ := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		cmd.Stdout = f
		cmd.Stderr = f
	}

	if err := cmd.Start(); err != nil {
		app.Status = "ERROR"
		statemanager.SaveApp(app)
		return fmt.Errorf("failed to start app: %w", err)
	}

	app.Status = "RUNNING"
	app.Pid = cmd.Process.Pid
	app.StartTime = time.Now().UnixMilli()
	statemanager.SaveApp(app)

	go waitForExit(app.Name, cmd)

	return nil
}

func waitForExit(name string, cmd *exec.Cmd) {
	err := cmd.Wait()

	mu.Lock()
	defer mu.Unlock()

	app, errGet := statemanager.GetApp(name)
	if errGet != nil {
		return
	}

	// Check if we are still the active process
	if app.Pid != cmd.Process.Pid {
		return
	}

	// If manual stop, status is already STOPPED
	if app.Status == "STOPPED" {
		app.Pid = 0
		statemanager.SaveApp(*app)
		return
	}

	if err != nil {
		app.Status = "ERROR"
	} else {
		// Normal exit is also treated as ERROR for auto-restart purposes unless manual stop
		// Spec says: "WHEN 应用进程退出 (无论 Exit Code 为何) 且状态非 STOPPED -> THEN 系统应更新状态为 ERROR"
		app.Status = "ERROR"
	}
	app.Pid = 0
	statemanager.SaveApp(*app)
}

func StopApp(name string) error {
	mu.Lock()
	defer mu.Unlock()

	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	app, err := statemanager.GetApp(name)
	if err != nil {
		return fmt.Errorf("app %s not found", name)
	}

	if app.Pid == 0 {
		return nil
	}

	process, err := os.FindProcess(app.Pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		process.Kill()
	}

	app.Status = "STOPPED"
	app.Pid = 0
	RemoveNginxConfig(dataDir, name)
	return statemanager.SaveApp(*app)
}

func DeleteApp(name string) error {
	mu.Lock()
	defer mu.Unlock()

	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	app, err := statemanager.GetApp(name)
	if err != nil {
		// If app doesn't exist in DB, check filesystem just in case
		appDir := filepath.Join(dataDir, "apps", name)
		if _, err := os.Stat(appDir); !os.IsNotExist(err) {
			os.RemoveAll(appDir)
		}
		return nil
	}

	// 1. Stop if running
	if app.Pid > 0 {
		if proc, err := os.FindProcess(app.Pid); err == nil {
			proc.Signal(syscall.SIGTERM)
			// Give it a moment, then force kill if needed?
			// For now just fire signal
		}
	}

	// 2. Remove Nginx Config
	RemoveNginxConfig(dataDir, name)

	// 3. Remove App Directory (logs, binaries)
	appDir := filepath.Join(dataDir, "apps", name)
	os.RemoveAll(appDir)

	// 4. Remove from State
	return statemanager.DeleteApp(name)
}

func ListApps() []api.AppInfo {
	dbApps, err := statemanager.ListApps()
	if err != nil {
		return []api.AppInfo{}
	}

	activeApps := appcenter.GetActiveApps()
	var result []api.AppInfo

	for _, app := range dbApps {
		if active, ok := activeApps[app.Name]; ok {
			// Use active info (which might have fresher data if updated, though mostly PID)
			// But DB has config/env potentially updated?
			// Let's assume DB is source of config, AppCenter is source of Liveness/PID
			app.Status = "RUNNING"
			app.Pid = active.Pid
			app.Port = active.Port // Ensure port is also updated from active info if available
			// We could merge stats here if appcenter had them
			// Actually, scanAndMonitor updates stats in DB, so dbApps should have them if queried recently.
			// But scanAndMonitor runs every 5s.
			// If we want real-time stats on ListApps, we should fetch them now.
			if proc, err := process.NewProcess(int32(app.Pid)); err == nil {
				cpu, _ := proc.CPUPercent()
				mem, _ := proc.MemoryInfo()
				ioStat, _ := proc.IOCounters()
				app.Stats.CPUPercent = cpu
				if mem != nil {
					app.Stats.MemoryUsage = mem.RSS
				}
				if ioStat != nil {
					app.Stats.IOReadBytes = ioStat.ReadBytes
					app.Stats.IOWriteBytes = ioStat.WriteBytes
				}
			}
		} else {
			// Not active
			if app.Status == "RUNNING" {
				// It was running but now disconnected?
				// ListApps is just a view, scanAndMonitor handles the state transition.
				// But for display, we might want to show it as "UNKNOWN" or keep DB state
				// Let's keep DB state, monitor will fix it shortly.
			}
		}
		result = append(result, app)
	}
	return result
}

func StartMonitor() {
	go runMonitor()
}

func runMonitor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scanAndMonitor()
		}
	}
}

func scanAndMonitor() {
	apps, err := statemanager.ListApps()
	if err != nil {
		return
	}

	activeApps := appcenter.GetActiveApps()

	for _, app := range apps {
		if activeInfo, ok := activeApps[app.Name]; ok {
			// App is connected to AppCenter -> It is RUNNING
			app.Status = "RUNNING"
			app.Pid = activeInfo.Pid

			// Update Stats using PID if process exists locally
			if app.Pid > 0 {
				if proc, err := process.NewProcess(int32(app.Pid)); err == nil {
					cpu, _ := proc.CPUPercent()
					mem, _ := proc.MemoryInfo()
					ioStat, _ := proc.IOCounters()
					createTime, _ := proc.CreateTime()

					app.Stats.CPUPercent = cpu
					if mem != nil {
						app.Stats.MemoryUsage = mem.RSS
					}
					if ioStat != nil {
						app.Stats.IOReadBytes = ioStat.ReadBytes
						app.Stats.IOWriteBytes = ioStat.WriteBytes
					}
					app.StartTime = createTime
				}
			}

			statemanager.SaveApp(app)
		} else {
			// App is NOT connected to AppCenter
			if app.Status == "RUNNING" || app.Status == "STARTING" {
				// It should be running but isn't connected -> ERROR (Crash/Disconnect)
				// Note: STARTING gives it a grace period?
				// For now, if it's not connected, it's not running.
				// But StartApp sets status to STARTING then RUNNING.
				// If we mark ERROR immediately, it might race with startup connection.
				// Ideally, we should check if process exists by PID from DB if we want to be sure it's not just a network issue?
				// But user said "use appcenter status". AppCenter status = Connected.
				// Let's assume if it's not in AppCenter, it's dead.

				// Grace period logic could be added here if needed.
				// For now, strictly follow "Not in AppCenter = Dead"
				app.Status = "ERROR"
				app.Pid = 0
				statemanager.SaveApp(app)
			}
		}

		// Restart Logic
		if app.AutoRestart && app.Status != "RUNNING" && app.Status != "STOPPED" {
			if app.RestartCount > 5 {
				continue
			}

			fmt.Printf("Restarting app %s (Attempt %d)\n", app.Name, app.RestartCount+1)

			app.RestartCount++
			statemanager.SaveApp(app)

			req := api.StartAppRequest{
				Name:       app.Name,
				Command:    app.Command,
				Args:       app.Args,
				WorkingDir: app.WorkingDir,
				Env:         app.Env,
				AutoRestart: app.AutoRestart,
			}

			go func(r api.StartAppRequest) {
				if err := StartApp(r); err != nil {
					log.Printf("Failed to restart app %s: %v", r.Name, err)
				}
			}(req)
		}
	}
}

func copyFile(src, dst string) error {
	if src == dst {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func calculateFileHash(path string) (string, error) {
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
