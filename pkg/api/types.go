package api

import "encoding/json"

type AppInfo struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	Config      map[string]any    `json:"config"` // New: Application Config
	Port        int               `json:"port"`   // New: Assigned Port
	Domain       string            `json:"domain"` // New: Domain for Nginx
	AutoRestart  bool              `json:"auto_restart"`
	RestartCount int               `json:"restart_count"`
	Status       string            `json:"status"` // RUNNING, STOPPED, ERROR
	Pid         int               `json:"pid"`
	Stats       AppStats          `json:"stats"`
}

type AppStats struct {
	CPUPercent   float64 `json:"cpu_percent"`
	MemoryUsage  uint64  `json:"memory_usage"` // bytes
	IOReadBytes  uint64  `json:"io_read_bytes"`
	IOWriteBytes uint64  `json:"io_write_bytes"`
}

type StartAppRequest struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	Config      map[string]any    `json:"config"` // New: Application Config
	AutoRestart bool              `json:"auto_restart"`
}

type StopAppRequest struct {
	Name string `json:"name"`
}

type IngressUpdateRequest struct {
	AppName string `json:"app_name"`
	Domain  string `json:"domain"`
	Port    int    `json:"port"` // Optional: If 0, try to find running app port
}

type IngressDeleteRequest struct {
	AppName string `json:"app_name"`
}

type ProvisionRequest struct {
	AppName      string `json:"app_name"`
	ResourceType string `json:"resource_type"` // e.g., "mysql", "redis"
	ResourceName string `json:"resource_name"` // e.g., "billing_db"
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type TCPAction string

const (
	ActionGetConfig TCPAction = "get_config"
	ActionProvision TCPAction = "provision"
	ActionRegister  TCPAction = "register"
	ActionAppStart  TCPAction = "app_start"
)

type TCPRequest struct {
	Action  TCPAction       `json:"action"`
	AppName string          `json:"app_name"`
	APIKey  string          `json:"api_key"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
