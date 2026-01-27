package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/luaxlou/glow/internal/apiserver"
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/internal/manager"
	"github.com/luaxlou/glow/internal/statemanager"
	"github.com/luaxlou/glow/pkg/api"
	"github.com/luaxlou/glow/starter/glowhttp"
	"github.com/luaxlou/glow/starter/glowsqlite"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

var (
	serverPort int
	dataDir    string
	maxAgeDays int
	maxTotalMB int
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 32102, "HTTP Port")
	serveCmd.Flags().StringVar(&dataDir, "data-dir", "", "Data directory path (default: /var/lib/glow-server)")
	serveCmd.Flags().IntVar(&maxAgeDays, "log-max-age-days", 30, "Maximum age of log files in days")
	serveCmd.Flags().IntVar(&maxTotalMB, "log-max-total-mb", 500, "Maximum total size of log files in MB")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Glow Server (HTTP API)",
	Run: func(cmd *cobra.Command, args []string) {
		// Set up data directory
		if err := setupDataDir(); err != nil {
			log.Fatalf("Failed to setup data directory: %v", err)
		}

		// Start log cleaner
		logDir := filepath.Join(dataDir, "logs")
		manager.StartPeriodicCleanup(logDir, maxAgeDays, maxTotalMB, 24*time.Hour)
		log.Printf("Started log cleaner (max age: %d days, max total: %d MB)", maxAgeDays, maxTotalMB)

		// Initialize Config Manager
		if err := configmanager.EnsureInitialized(); err != nil {
			log.Fatalf("Failed to init config manager: %v", err)
		}

		// Self-register glow-server ASAP
		serverAppInfo := api.AppInfo{
			Name:      "glow-server",
			Status:    "RUNNING",
			Pid:       os.Getpid(),
			Port:      serverPort,
			StartTime: time.Now().UnixMilli(),
		}

		fmt.Printf("Registering glow-server in DB (PID: %d)...\n", os.Getpid())
		if err := statemanager.SaveApp(serverAppInfo); err != nil {
			fmt.Printf("ERROR: Failed to register glow-server in DB: %v\n", err)
		} else {
			fmt.Printf("Successfully registered glow-server in DB.\n")
		}

		// Get API Key (Checking existence)
		apiKey, err := configmanager.GetSystemConfig("api_key")
		if err != nil || apiKey == "" {
			log.Fatal("API Key not found. Please run 'glow-server keygen' first.")
		}

		glowhttp.Init(serverPort)

		// Setup API Server
		server := apiserver.New()
		server.RegisterRoutes(glowhttp.Router())

		glowhttp.Run()

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
	},
}

func setupDataDir() error {
	// Determine data directory
	if dataDir == "" {
		dataDir = "/var/lib/glow-server"
	}
	if absDir, err := filepath.Abs(dataDir); err == nil {
		dataDir = absDir
	}

	// Create data directory and subdirectories
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{
		filepath.Join(dataDir, "db"),
		filepath.Join(dataDir, "logs"),
		filepath.Join(dataDir, "apps"),
		filepath.Join(dataDir, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Update DB path to use data directory
	dbPath := filepath.Join(dataDir, "db", "glow.db")
	glowsqlite.Init(dbPath)

	if err := configmanager.SetSystemConfig("data_dir", dataDir); err != nil {
		return err
	}

	return nil
}
