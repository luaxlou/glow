package manager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// NginxInfo stores detected nginx information
type NginxInfo struct {
	BinaryPath string
	ConfPath   string
	PrefixPath string
	ErrorLog   string
	AccessLog  string
	Version    string
}

const nginxTemplate = `
upstream {{.Name}} {
    server 127.0.0.1:{{.Port}};
}

server {
    listen 80;
    server_name {{.Domain}};
    client_max_body_size 100M;

    location / {
        proxy_pass http://{{.Name}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }
}
`

type NginxConfig struct {
	Name   string
	Port   int
	Domain string
}

// DetectNginx attempts to find Nginx and its configuration
func DetectNginx() (*NginxInfo, error) {
	path, err := exec.LookPath("nginx")
	if err != nil {
		return nil, fmt.Errorf("nginx not found in PATH")
	}

	info := &NginxInfo{
		BinaryPath: path,
	}

	// Parse 'nginx -V' output
	cmd := exec.Command(path, "-V")
	// nginx -V writes to stderr
	out, err := cmd.CombinedOutput()
	if err != nil {
		return info, fmt.Errorf("failed to run nginx -V: %v", err)
	}

	output := string(out)

	// Parse Version
	if matches := regexp.MustCompile(`nginx version: nginx/([\d.]+)`).FindStringSubmatch(output); len(matches) > 1 {
		info.Version = matches[1]
	}

	// Parse Paths
	args := strings.Fields(output)
	for _, arg := range args {
		if strings.HasPrefix(arg, "--conf-path=") {
			info.ConfPath = strings.TrimPrefix(arg, "--conf-path=")
		} else if strings.HasPrefix(arg, "--error-log-path=") {
			info.ErrorLog = strings.TrimPrefix(arg, "--error-log-path=")
		} else if strings.HasPrefix(arg, "--http-log-path=") {
			info.AccessLog = strings.TrimPrefix(arg, "--http-log-path=")
		} else if strings.HasPrefix(arg, "--prefix=") {
			info.PrefixPath = strings.TrimPrefix(arg, "--prefix=")
		}
	}

	return info, nil
}

// GenerateNginxConfig creates the config file and reloads nginx
func GenerateNginxConfig(dataDir string, cfg NginxConfig) error {
	if cfg.Domain == "" {
		return nil // No domain specified, skip nginx
	}

	// 1. Detect Nginx
	info, err := DetectNginx()
	if err != nil {
		log.Printf("Nginx detection failed: %v. Only generating local config file.", err)
	}

	// 2. Generate Config File in local dataDir (as backup/reference)
	nginxDir := filepath.Join(dataDir, "nginx")
	if err := os.MkdirAll(nginxDir, 0755); err != nil {
		return err
	}

	tmpl, err := template.New("nginx").Parse(nginxTemplate)
	if err != nil {
		return err
	}

	fileName := filepath.Join(nginxDir, cfg.Name+".conf")
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(f, cfg); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// 3. Link or Copy to Nginx Config Directory (if detected)
	if err != nil || info == nil || info.ConfPath == "" {
		if err != nil {
			return fmt.Errorf("nginx not available: %w", err)
		}
		return fmt.Errorf("nginx not available: conf path not detected")
	}

	confDir := filepath.Dir(info.ConfPath)

	// Common include directories
	candidates := []string{
		filepath.Join(confDir, "servers"),
		filepath.Join(confDir, "conf.d"),
		filepath.Join(confDir, "sites-enabled"),
	}

	var targetDir string
	for _, dir := range candidates {
		if s, statErr := os.Stat(dir); statErr == nil && s.IsDir() {
			targetDir = dir
			break
		}
	}

	if targetDir == "" {
		return fmt.Errorf("no suitable nginx include directory found under %s", confDir)
	}

	targetFile := filepath.Join(targetDir, cfg.Name+".conf")
	log.Printf("Installing Nginx config to: %s", targetFile)

	input, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return readErr
	}
	if writeErr := os.WriteFile(targetFile, input, 0644); writeErr != nil {
		return fmt.Errorf("failed to install nginx config to %s: %w", targetFile, writeErr)
	}
	if reloadErr := ReloadNginx(info.BinaryPath); reloadErr != nil {
		return fmt.Errorf("failed to reload nginx: %w", reloadErr)
	}

	return nil
}

func ReloadNginx(binaryPath string) error {
	cmd := exec.Command(binaryPath, "-s", "reload")
	return cmd.Run()
}

func ListIngress(dataDir string) ([]NginxConfig, error) {
	nginxDir := filepath.Join(dataDir, "nginx")
	files, err := os.ReadDir(nginxDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []NginxConfig{}, nil
		}
		return nil, err
	}

	var configs []NginxConfig
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".conf") {
			name := strings.TrimSuffix(f.Name(), ".conf")
			cfg, err := GetIngress(dataDir, name)
			if err == nil {
				configs = append(configs, *cfg)
			}
		}
	}
	return configs, nil
}

func GetIngress(dataDir string, name string) (*NginxConfig, error) {
	fileName := filepath.Join(dataDir, "nginx", name+".conf")
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cfg := &NginxConfig{Name: name}

	// Simple parsing from file content (upstream and server_name)
	// Upstream: server 127.0.0.1:PORT;
	rePort := regexp.MustCompile(`server 127.0.0.1:(\d+);`)
	if matches := rePort.FindStringSubmatch(string(content)); len(matches) > 1 {
		cfg.Port, _ = strconv.Atoi(matches[1])
	}

	// Server Name: server_name DOMAIN;
	reDomain := regexp.MustCompile(`server_name ([\w.-]+);`)
	if matches := reDomain.FindStringSubmatch(string(content)); len(matches) > 1 {
		cfg.Domain = matches[1]
	}

	return cfg, nil
}

func RemoveNginxConfig(dataDir string, name string) error {
	// Remove local
	fileName := filepath.Join(dataDir, "nginx", name+".conf")
	os.Remove(fileName)

	// Remove system (if possible)
	info, err := DetectNginx()
	if err == nil && info.ConfPath != "" {
		confDir := filepath.Dir(info.ConfPath)
		candidates := []string{
			filepath.Join(confDir, "servers"),
			filepath.Join(confDir, "conf.d"),
			filepath.Join(confDir, "sites-enabled"),
		}
		for _, dir := range candidates {
			targetFile := filepath.Join(dir, name+".conf")
			if _, err := os.Stat(targetFile); err == nil {
				if err := os.Remove(targetFile); err == nil {
					ReloadNginx(info.BinaryPath)
				}
			}
		}
	}
	return nil
}
