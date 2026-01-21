package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/spf13/cobra"
	"github.com/luaxlou/glow/starter/glowsqlite"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate or retrieve the API Key",
	Run:   runKeygen,
}

func init() {
	rootCmd.AddCommand(keygenCmd)
}

func initDB() error {
	// Determine data directory
	dataDir := os.Getenv("GLOW_DATA_DIR")
	if dataDir == "" {
		dataDir = "/var/lib/glow-server"
	}

	// Create data directory if it doesn't exist
	dbDir := filepath.Join(dataDir, "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	// Initialize DB
	dbPath := filepath.Join(dbDir, "glow.db")
	glowsqlite.Init(dbPath)
	return nil
}

func runKeygen(cmd *cobra.Command, args []string) {
	// Initialize DB
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Check if key exists
	key, err := configmanager.GetSystemConfig("api_key")
	if err == nil && key != "" {
		fmt.Printf("Existing API Key: %s\n", key)
		return
	}

	// Generate new key
	newKey, err := generateRandomKey(32)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	// Save key
	err = configmanager.SetSystemConfig("api_key", newKey)
	if err != nil {
		log.Fatalf("Failed to save key: %v", err)
	}

	fmt.Printf("Generated New API Key: %s\n", newKey)
}

func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
