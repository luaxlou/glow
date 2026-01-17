package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/luaxlou/glow/internal/appcenter"
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/internal/manager"
	"github.com/luaxlou/glow/internal/statemanager"
	"github.com/luaxlou/glow/pkg/api"
	"github.com/shirou/gopsutil/v3/process"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", s.handleHealth)

	// --- Config Management ---
	r.Any("/config/*appName", s.handleConfig)

	// --- App Management ---
	r.POST("/apps/upload", s.handleUploadApp)
	r.POST("/apps/start", s.handleStartApp)
	r.POST("/apps/stop", s.handleStopApp)
	r.POST("/apps/restart", s.handleRestartApp)
	r.POST("/apps/delete", s.handleDeleteApp)
	r.GET("/apps/list", s.handleListApps)
	r.GET("/apps/logs", s.handleAppLogs)

	// --- Node Management ---
	r.GET("/node/status", s.handleNodeStatus) // New

	// --- Ingress Management ---
	r.POST("/ingress/update", s.handleUpdateIngress)
	r.POST("/ingress/delete", s.handleDeleteIngress)
	r.GET("/ingress/list", s.handleListIngress)
}

func (s *Server) handleHealth(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (s *Server) handleUpdateIngress(c *gin.Context) {
	var req api.IngressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}

	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	port := req.Port
	if port == 0 {
		// Try to find app port
		apps := manager.ListApps()
		for _, app := range apps {
			if app.Name == req.AppName {
				port = app.Port
				break
			}
		}
	}

	if port == 0 {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "App not found or port not specified"})
		return
	}

	if err := manager.GenerateNginxConfig(dataDir, manager.NginxConfig{
		Name:   req.AppName,
		Port:   port,
		Domain: req.Domain,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}

	// Update AppInfo in StateManager
	if app, err := statemanager.GetApp(req.AppName); err == nil {
		app.Domain = req.Domain
		statemanager.SaveApp(*app)
	}

	// Update AppInfo in StateManager
	if app, err := statemanager.GetApp(req.AppName); err == nil {
		app.Domain = req.Domain
		statemanager.SaveApp(*app)
	}

	c.JSON(http.StatusOK, api.Response{Success: true, Message: "Ingress updated"})
}

func (s *Server) handleDeleteIngress(c *gin.Context) {
	var req api.IngressDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}

	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	if err := manager.RemoveNginxConfig(dataDir, req.AppName); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}

	// Update AppInfo in StateManager
	if app, err := statemanager.GetApp(req.AppName); err == nil {
		app.Domain = ""
		statemanager.SaveApp(*app)
	}

	c.JSON(http.StatusOK, api.Response{Success: true, Message: "Ingress deleted"})
}

func (s *Server) handleListIngress(c *gin.Context) {
	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	configs, err := manager.ListIngress(dataDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true, Data: configs})
}

func (s *Server) handleConfig(c *gin.Context) {
	appName := c.Param("appName")
	if len(appName) > 1 {
		appName = appName[1:]
	} else {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "App name required"})
		return
	}

	if c.Request.Method == http.MethodGet {
		config, err := configmanager.Get(appName)
		if err != nil {
			c.JSON(http.StatusNotFound, api.Response{Success: false, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, api.Response{Success: true, Data: config})
		return
	}

	if c.Request.Method == http.MethodPut {
		var newConfig map[string]any
		if err := c.ShouldBindJSON(&newConfig); err != nil {
			c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid JSON"})
			return
		}

		merge := c.Query("merge") != "false"
		if err := configmanager.Set(appName, newConfig, merge); err != nil {
			c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Failed to update config"})
			return
		}

		// Notify AppCenter (Hot Reload)
		if err := appcenter.SendConfigUpdate(appName, newConfig); err != nil {
			log.Printf("Warning: failed to push config update to app %s: %v\n", appName, err)
		}

		c.JSON(http.StatusOK, api.Response{Success: true, Message: "Config updated"})
		return
	}

	c.JSON(http.StatusMethodNotAllowed, api.Response{Success: false, Message: "Method not allowed"})
}

func (s *Server) handleUploadApp(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "File required"})
		return
	}

	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}
	if absDir, err := filepath.Abs(dataDir); err == nil {
		dataDir = absDir
	}
	tempDir := filepath.Join(dataDir, "tmp", "uploads")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Failed to create temp dir"})
		return
	}

	// Use original filename but safe
	dst := filepath.Join(tempDir, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true, Data: dst})
}

func (s *Server) handleStartApp(c *gin.Context) {
	var req api.StartAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}
	if err := manager.StartApp(req); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Message: "App started"})
}

func (s *Server) handleStopApp(c *gin.Context) {
	var req api.StopAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}
	if err := manager.StopAppWithOptions(req.Name, manager.StopAppOptions{KeepIngress: req.KeepIngress}); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Message: "App stopped"})
}

func (s *Server) handleDeleteApp(c *gin.Context) {
	var req api.StopAppRequest // Re-use StopAppRequest as it just needs Name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}
	if err := manager.DeleteApp(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Message: "App deleted"})
}

func (s *Server) handleRestartApp(c *gin.Context) {
	var req api.StopAppRequest // Use StopAppRequest struct as it just needs Name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}

	// 1. Find existing app info to preserve config/env/args
	apps := manager.ListApps()
	var targetApp *api.AppInfo
	for _, app := range apps {
		if app.Name == req.Name {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		c.JSON(http.StatusNotFound, api.Response{Success: false, Message: "App not found"})
		return
	}

	// Legacy support: apps may have registered without Command/WorkingDir. Try to infer from PID.
	if targetApp.Command == "" && targetApp.Pid != 0 {
		if p, err := process.NewProcess(int32(targetApp.Pid)); err == nil {
			if exe, err := p.Exe(); err == nil && exe != "" {
				targetApp.Command = exe
			}
			if targetApp.WorkingDir == "" {
				if cwd, err := p.Cwd(); err == nil && cwd != "" {
					targetApp.WorkingDir = cwd
				}
			}
		}
	}

	// If Command is still empty, check for deployed binary
	if targetApp.Command == "" {
		dataDir, _ := configmanager.GetSystemConfig("data_dir")
		if dataDir == "" {
			dataDir = "."
		}
		if absDir, err := filepath.Abs(dataDir); err == nil {
			dataDir = absDir
		}
		appDir := filepath.Join(dataDir, "apps", targetApp.Name)
		dstBinaryPath := filepath.Join(appDir, "glow_"+targetApp.Name)
		
		// Check if deployed binary exists
		if _, err := os.Stat(dstBinaryPath); err == nil {
			targetApp.Command = dstBinaryPath
		} else {
			c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: fmt.Sprintf("Failed to restart: app '%s' not found or no command specified", req.Name)})
			return
		}
	}

	// 2. Stop first (ignore error if already stopped)
	manager.StopApp(req.Name)

	// 3. Start again with same parameters
	startReq := api.StartAppRequest{
		Name:        targetApp.Name,
		Command:     targetApp.Command,
		Args:        targetApp.Args,
		WorkingDir:  targetApp.WorkingDir,
		Env:         targetApp.Env,
		AutoRestart: targetApp.AutoRestart,
		Config:      targetApp.Config,
	}

	if err := manager.StartApp(startReq); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: fmt.Sprintf("Failed to restart: %v", err)})
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true, Message: "App restarted"})
}

func (s *Server) handleListApps(c *gin.Context) {
	apps := manager.ListApps()
	c.JSON(http.StatusOK, api.Response{Success: true, Data: apps})
}

func (s *Server) handleAppLogs(c *gin.Context) {
	appName := c.Query("name")
	if appName == "" {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "name required"})
		return
	}

	// Get dataDir from system config
	dataDir, _ := configmanager.GetSystemConfig("data_dir")
	if dataDir == "" {
		dataDir = "."
	}

	// Simple log reading
	logFile := filepath.Join(dataDir, "apps", appName, "logs", appName+".log")
	if _, err := os.Stat(logFile); err != nil {
		c.JSON(http.StatusNotFound, api.Response{Success: false, Message: "Log file not found"})
		return
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Failed to read logs"})
		return
	}

	// Limit log size? For now return all.
	c.JSON(http.StatusOK, api.Response{Success: true, Data: string(content)})
}

func (s *Server) handleNodeStatus(c *gin.Context) {
	node, err := manager.GetNodeStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Data: node})
}
