package api

import "encoding/json"

// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

// ObjectMeta is metadata that all persisted resources must have.
type ObjectMeta struct {
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// --- Deployment (App) ---

// Deployment represents an application deployment.
// It maps to the concept of an "App" in Glow.
type Deployment struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       AppSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     AppStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// AppSpec defines the desired state of an application.
// Replaces old AppSpec/StartAppRequest mix.
type AppSpec struct {
	Command     string            `json:"command" yaml:"command"`
	Args        []string          `json:"args,omitempty" yaml:"args,omitempty"`
	WorkingDir  string            `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	Env         map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Port        int               `json:"port,omitempty" yaml:"port,omitempty"`
	Domain      string            `json:"domain,omitempty" yaml:"domain,omitempty"`
	Replicas    int               `json:"replicas,omitempty" yaml:"replicas,omitempty"` // For future scaling
	AutoRestart bool              `json:"autoRestart,omitempty" yaml:"autoRestart,omitempty"`
	Config      map[string]any    `json:"config,omitempty" yaml:"config,omitempty"`
}

// AppStatus defines the observed state of an application.
type AppStatus struct {
	Phase        string   `json:"phase,omitempty" yaml:"phase,omitempty"` // RUNNING, STOPPED, ERROR
	Pid          int      `json:"pid,omitempty" yaml:"pid,omitempty"`
	RestartCount int      `json:"restartCount,omitempty" yaml:"restartCount,omitempty"`
	Stats        AppStats `json:"stats,omitempty" yaml:"stats,omitempty"`
}

// --- Node (Host) ---

// Node represents a host machine.
type Node struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Status     NodeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type NodeStatus struct {
	Hostname   string       `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	OS         string       `json:"os,omitempty" yaml:"os,omitempty"`
	Arch       string       `json:"arch,omitempty" yaml:"arch,omitempty"`
	Kernel     string       `json:"kernel,omitempty" yaml:"kernel,omitempty"`
	CPUUsage   float64      `json:"cpuUsage,omitempty" yaml:"cpuUsage,omitempty"`
	MemUsage   float64      `json:"memUsage,omitempty" yaml:"memUsage,omitempty"` // Percent
	DiskUsage  float64      `json:"diskUsage,omitempty" yaml:"diskUsage,omitempty"` // Percent
	Resources  []ResourceRef `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type ResourceRef struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	Port int    `json:"port,omitempty"`
}

// --- Ingress ---

type Ingress struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       IngressSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type IngressSpec struct {
	Domain  string `json:"domain" yaml:"domain"`
	Service string `json:"service" yaml:"service"` // App Name
	Port    int    `json:"port" yaml:"port"`
}

// --- Legacy & Shared Types ---

type AppStats struct {
	CPUPercent   float64 `json:"cpu_percent"`
	MemoryUsage  uint64  `json:"memory_usage"` // bytes
	IOReadBytes  uint64  `json:"io_read_bytes"`
	IOWriteBytes uint64  `json:"io_write_bytes"`
}

// Helper types for API communication (Requests/Responses)
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// Deprecated: StartAppRequest (Migrate to Deployment)
type StartAppRequest struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	Config      map[string]any    `json:"config"`
	Domain      string            `json:"domain"`
	AutoRestart bool              `json:"auto_restart"`
}

// Deprecated: StopAppRequest
type StopAppRequest struct {
	Name string `json:"name"`
}

// Deprecated: AppInfo (Use Deployment.Status instead)
type AppInfo struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	Config      map[string]any    `json:"config"`
	Port        int               `json:"port"`
	Domain       string            `json:"domain"`
	AutoRestart  bool              `json:"auto_restart"`
	RestartCount int               `json:"restart_count"`
	StartTime    int64             `json:"start_time"` // Unix timestamp
	Status       string            `json:"status"`
	Pid          int               `json:"pid"`
	Stats        AppStats          `json:"stats"`
}

type ProvisionRequest struct {
	AppName      string `json:"app_name"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
}

type IngressUpdateRequest struct {
	AppName string `json:"app_name"`
	Domain  string `json:"domain"`
	Port    int    `json:"port"`
}

type IngressDeleteRequest struct {
	AppName string `json:"app_name"`
}

// Host Manifest (Old) - Keep for backward compatibility or refactor
type Host struct {
	TypeMeta   `yaml:",inline"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       HostSpec   `yaml:"spec" json:"spec"`
}

type HostSpec struct {
	PublicIP string                  `yaml:"publicIP" json:"publicIP"`
	Services map[string]ServiceSpec  `yaml:"services" json:"services"`
}

type ServiceSpec struct {
	Port          int    `yaml:"port" json:"port"`
	AdminUser     string `yaml:"adminUser" json:"adminUser"`
	AdminPassword string `yaml:"adminPassword" json:"adminPassword"`
}

// App Manifest (Old)
type App struct {
	TypeMeta   `yaml:",inline"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       AppSpecOld `yaml:"spec" json:"spec"`
}

type AppSpecOld struct {
	Binary       string                 `yaml:"binary" json:"binary"`
	Command      string                 `yaml:"command" json:"command"`
	Args         []string               `yaml:"args" json:"args"`
	WorkingDir   string                 `yaml:"workingDir" json:"workingDir"`
	Domain       string                 `yaml:"domain" json:"domain"`
	Dependencies map[string]DepSpec     `yaml:"dependencies" json:"dependencies"`
}

type DepSpec struct {
	DBName string `yaml:"dbName,omitempty" json:"dbName,omitempty"`
	DB     int    `yaml:"db,omitempty" json:"db,omitempty"`
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
