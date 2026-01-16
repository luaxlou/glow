package cmd

import (
	"fmt"
	"os"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Glow Server configuration and status",
	Run:   runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) {
	fmt.Println("Glow Server Information")
	fmt.Println("-----------------------")

	// 1. Version and PID
	fmt.Printf("PID: %d\n", os.Getpid())

	// 5. Nginx
	var nginxInfo api.NginxSystemConfig
	if err := getSystemConfigJSON("nginx_info", &nginxInfo); err == nil && nginxInfo.BinaryPath != "" {
		fmt.Printf("Nginx: %s (Version: %s)\n", nginxInfo.BinaryPath, nginxInfo.Version)
	} else {
		fmt.Println("Nginx: [NOT CONFIGURED]")
	}

	// 6. Service status (simplistic check)
	// On Linux, we could check systemctl status glow-server
	// For now, just placeholder
	fmt.Println("Service: [STATUS CHECK NOT IMPLEMENTED]")
}
