package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/luaxlou/glow/internal/apiserver"
	"github.com/luaxlou/glow/internal/appcenter"
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/internal/statemanager"
	"github.com/luaxlou/glow/pkg/api"
	"github.com/luaxlou/glow/starter/glowapp"
	"github.com/luaxlou/glow/starter/glowhttp"
	"github.com/spf13/cobra"
)

var (
	serverPort    int
	appCenterPort int
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 32102, "HTTP Port")
	serveCmd.Flags().IntVarP(&appCenterPort, "app-center-port", "a", 32101, "App Center Port")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Glow Server (HTTP API)",
	Run: func(cmd *cobra.Command, args []string) {
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
		appcenter.RegisterActiveApp(serverAppInfo, nil)
		
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

		log.Printf("Starting App Center on port %d...", appCenterPort)
		if err := appcenter.Start(appCenterPort); err != nil {
			log.Fatalf("Failed to start App Center: %v", err)
		}

		glowapp.Init("glow-server", glowapp.WithNoRegistration())
		glowhttp.Init(serverPort)

		// Setup API Server
		server := apiserver.New()
		server.RegisterRoutes(glowhttp.Router())

		glowhttp.Run()
		glowapp.WaitForShutdown()
	},
}
