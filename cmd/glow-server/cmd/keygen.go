package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(keygenCmd)
}

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate or retrieve the API Key",
	Run:   runKeygen,
}

func runKeygen(cmd *cobra.Command, args []string) {
	// Initialize Config Manager
	if err := configmanager.EnsureInitialized(); err != nil {
		log.Fatalf("Failed to init config manager: %v", err)
	}

	// Try to get existing key
	apiKey, err := getSystemConfig("api_key")
	if err != nil {
		log.Fatalf("Failed to retrieve API Key: %v", err)
	}
	if apiKey != "" {
		fmt.Println(apiKey)
		return
	}

	// Generate new key
	newKey, err := generateRandomKey(32)
	if err != nil {
		log.Fatalf("Failed to generate API Key: %v", err)
	}

	if err := setSystemConfig("api_key", newKey); err != nil {
		log.Fatalf("Failed to save API Key to system_config: %v", err)
	}

	fmt.Println(newKey)
}

func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
