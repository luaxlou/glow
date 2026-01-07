package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	var cfg ClientConfig
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Config file doesn't exist, proceed with empty config
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	updated := false
	reader := bufio.NewReader(os.Stdin)

	if cfg.ServerURL == "" {
		fmt.Print("Enter Server URL (default: http://localhost:32102): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			cfg.ServerURL = "http://localhost:32102"
		} else {
			cfg.ServerURL = input
		}
		updated = true
	}

	if cfg.APIKey == "" {
		fmt.Print("Enter API Key: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		cfg.APIKey = strings.TrimSpace(input)
		if cfg.APIKey != "" {
			updated = true
		}
	}

	if updated {
		if err := SaveClientConfig(&cfg); err != nil {
			return nil, fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Configuration saved to %s\n", path)
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required")
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
