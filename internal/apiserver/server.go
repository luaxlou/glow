package apiserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/internal/manager"
	"github.com/luaxlou/glow/internal/statemanager"
	"github.com/luaxlou/glow/pkg/api"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", s.handleHealth)

	// Apply Auth Middleware
	r.Use(authMiddleware())

	// --- Config Management ---
	r.Any("/config/*appName", s.handleConfig)

	// --- App Management ---
	r.POST("/apps/start", s.handleStartApp)
	r.POST("/apps/stop", s.handleStopApp)
	r.POST("/apps/restart", s.handleRestartApp)
	r.POST("/apps/delete", s.handleDeleteApp)
	r.GET("/apps/list", s.handleListApps)
	r.GET("/apps/logs", s.handleAppLogs)

	// --- Resource Provisioning ---
	r.POST("/resources/provision", s.handleProvisionResource)
	r.GET("/resources/list", s.handleListResources) // New

	// --- Node Management ---
	r.GET("/node/status", s.handleNodeStatus) // New

	// --- Ingress Management ---
	r.POST("/ingress/update", s.handleUpdateIngress)
	r.POST("/ingress/delete", s.handleDeleteIngress)
	r.GET("/ingress/list", s.handleListIngress)

	// --- Manifest Application ---
	r.POST("/apply/host", s.handleApplyHost)
	r.POST("/apply/app", s.handleApplyApp)
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
		if err := configmanager.Set(appName, newConfig, true); err != nil {
			c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Failed to update config"})
			return
		}
		c.JSON(http.StatusOK, api.Response{Success: true, Message: "Config updated"})
		return
	}

	c.JSON(http.StatusMethodNotAllowed, api.Response{Success: false, Message: "Method not allowed"})
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
	if err := manager.StopApp(req.Name); err != nil {
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

func (s *Server) handleProvisionResource(c *gin.Context) {
	var req api.ProvisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid request body"})
		return
	}
	config, err := manager.ProvisionResource(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Message: "Resource provisioned", Data: config})
}

func (s *Server) handleListResources(c *gin.Context) {
	resources, err := manager.ListResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Data: resources})
}

func (s *Server) handleNodeStatus(c *gin.Context) {
	node, err := manager.GetNodeStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Data: node})
}

func (s *Server) handleApplyHost(c *gin.Context) {
	var hostReq api.Host
	if err := c.ShouldBindJSON(&hostReq); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid manifest"})
		return
	}
	if err := configmanager.SaveHostConfig(hostReq); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Failed to save host config"})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Message: "Host config applied"})
}

func (s *Server) handleApplyApp(c *gin.Context) {
	var appReq api.App
	if err := c.ShouldBindJSON(&appReq); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Success: false, Message: "Invalid manifest"})
		return
	}
	if err := manager.ApplyApp(appReq); err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, api.Response{Success: true, Message: "App applied and started"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, err := configmanager.GetSystemConfig("api_key")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Internal Server Error: Failed to retrieve API Key"})
			return
		}
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Response{Success: false, Message: "Internal Server Error: API Key not configured"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{Success: false, Message: "Unauthorized: Missing Authorization header"})
			return
		}

		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{Success: false, Message: "Unauthorized: Invalid Authorization header format"})
			return
		}

		token := authHeader[len(prefix):]
		if token != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{Success: false, Message: "Unauthorized: Invalid API Key"})
			return
		}

		c.Next()
	}
}
