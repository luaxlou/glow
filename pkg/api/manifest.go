package api

type TypeMeta struct {
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
}

type ObjectMeta struct {
	Name string `yaml:"name" json:"name"`
}

// Host Manifest
type Host struct {
	TypeMeta   `yaml:",inline"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       HostSpec   `yaml:"spec" json:"spec"`
}

type HostSpec struct {
	PublicIP string                  `yaml:"publicIP" json:"publicIP"`
	Services map[string]ServiceSpec  `yaml:"services" json:"services"` // key: mysql, redis
}

type ServiceSpec struct {
	Port          int    `yaml:"port" json:"port"`
	AdminUser     string `yaml:"adminUser" json:"adminUser"`
	AdminPassword string `yaml:"adminPassword" json:"adminPassword"` // For simplicity in this demo
}

// App Manifest
type App struct {
	TypeMeta   `yaml:",inline"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       AppSpec    `yaml:"spec" json:"spec"`
}

type AppSpec struct {
	Binary       string                 `yaml:"binary" json:"binary"` // Path or name of binary
	Command      string                 `yaml:"command" json:"command"`
	Args         []string               `yaml:"args" json:"args"`
	WorkingDir   string                 `yaml:"workingDir" json:"workingDir"`
	Domain       string                 `yaml:"domain" json:"domain"`
	Dependencies map[string]DepSpec     `yaml:"dependencies" json:"dependencies"` // key: mysql, redis
}

type DepSpec struct {
	DBName string `yaml:"dbName,omitempty" json:"dbName,omitempty"` // for mysql
	DB     int    `yaml:"db,omitempty" json:"db,omitempty"`         // for redis
}
