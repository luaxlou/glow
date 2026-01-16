package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	serverURL string
	apiKey    string
)

type Config struct {
	CurrentContext string    `json:"current-context"`
	Contexts       []Context `json:"contexts"`

	// Legacy fields for migration
	LegacyServerURL string `json:"server_url,omitempty"`
	LegacyAPIKey    string `json:"api_key,omitempty"`
}

type Context struct {
	Name      string `json:"name"`
	ServerURL string `json:"server_url"`
	APIKey    string `json:"api_key"`
}

var rootCmd = &cobra.Command{
	Use:   "glow",
	Short: "Glow CLI - Kubernetes-style application governance",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config check for auth/context/help commands
		if cmd.Name() == "help" {
			return nil
		}
		if cmd.Parent() != nil && (cmd.Parent().Name() == "auth" || cmd.Parent().Name() == "context") {
			return nil
		}
		// Also skip for context/auth root commands themselves
		if cmd.Name() == "auth" || cmd.Name() == "context" {
			return nil
		}
		
		return ensureConfig()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	home, _ := os.UserHomeDir()
	cfgFile = filepath.Join(home, ".glow.json")
}

func ensureConfig() error {
	cfg, err := loadConfig()
	// Migration logic handled in loadConfig, but if still no contexts:
	if err != nil || len(cfg.Contexts) == 0 {
		fmt.Println("No context found, entering interactive setup...")
		reader := bufio.NewReader(os.Stdin)

		defaultURL := "http://localhost:32102"
		fmt.Printf("Enter Server URL [%s]: ", defaultURL)
		url, _ := reader.ReadString('\n')
		url = strings.TrimSpace(url)
		if url == "" {
			url = defaultURL
		}

		fmt.Print("Enter API Key: ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		// Create default context
		newCtx := Context{
			Name:      "default",
			ServerURL: url,
			APIKey:    key,
		}
		cfg = &Config{
			CurrentContext: "default",
			Contexts:       []Context{newCtx},
		}

		if err := saveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("Configuration saved.")
	}

	// Set globals from current context
	for _, ctx := range cfg.Contexts {
		if ctx.Name == cfg.CurrentContext {
			serverURL = ctx.ServerURL
			apiKey = ctx.APIKey
			break
		}
	}
	
	if serverURL == "" {
		return fmt.Errorf("current context '%s' not found or invalid", cfg.CurrentContext)
	}

	return nil
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Migration Logic
	if len(cfg.Contexts) == 0 && cfg.LegacyServerURL != "" {
		// Migrate legacy to default context
		cfg.Contexts = []Context{{
			Name:      "default",
			ServerURL: cfg.LegacyServerURL,
			APIKey:    cfg.LegacyAPIKey,
		}}
		cfg.CurrentContext = "default"
		cfg.LegacyServerURL = ""
		cfg.LegacyAPIKey = ""
		saveConfig(&cfg) // Persist migration
	}

	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	os.MkdirAll(filepath.Dir(cfgFile), 0755)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(cfgFile, data, 0644)
}
