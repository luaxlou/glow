package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ClientConfigName = "client.yaml"
)

type ClientConfig struct {
	ServerURL string `yaml:"server_url"`
	APIKey    string `yaml:"api_key"`
}

func GetClientConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDir, ClientConfigName), nil
}

func LoadClientConfig() (*ClientConfig, error) {
	path, err := GetClientConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found at %s. Please run 'init' or configure manually", path)
		}
		return nil, err
	}

	var cfg ClientConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveClientConfig(cfg *ClientConfig) error {
	path, err := GetClientConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
