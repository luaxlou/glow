package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/luaxlou/glow/internal/manager"
	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Interactive installer for Glow Server",
	Run:   runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Glow Server Installer!")
	fmt.Println("-------------------------------------")

	// 1. Service Installation
	if err := stepService(reader); err != nil {
		fmt.Printf("Service installation failed: %v\n", err)
		return
	}

	// 2. Ingress / Reverse Proxy
	if err := stepIngress(reader); err != nil {
		fmt.Printf("Ingress setup failed: %v\n", err)
		return
	}

	fmt.Println("\n-------------------------------------")
	fmt.Println("Glow Server installation/configuration complete!")
	fmt.Println("You can now run 'glow-server serve' to start the server.")
}

func stepService(reader *bufio.Reader) error {
	if isServiceInstalled() {
		fmt.Printf("\nGlow Server service is already installed. Reinstall? [y/N]: ")
	} else {
		fmt.Printf("\nWould you like to install Glow Server as a system service? [y/N]: ")
	}

	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" {
		return nil
	}

	fmt.Println("Installing system service...")
	return installService()
}

func stepIngress(reader *bufio.Reader) error {
	var nginxInfo api.NginxSystemConfig
	err := getSystemConfigJSON("nginx_info", &nginxInfo)
	if err != nil || nginxInfo.BinaryPath == "" {
		// Nginx not configured, skip
		return nil
	}

	fmt.Printf("\nNginx detected. Would you like to configure a reverse proxy for Glow Server? [y/N]: ")
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" {
		return nil
	}

	fmt.Print("Enter domain name (e.g., glow.example.com): ")
	domain, _ := reader.ReadString('\n')
	domain = strings.TrimSpace(domain)
	if domain == "" {
		fmt.Println("No domain provided, skipping.")
		return nil
	}

	fmt.Printf("Configuring Nginx for %s...\n", domain)

	cfg := manager.NginxConfig{
		Name:   "glow-server",
		Port:   32102, // Default glow-server port
		Domain: domain,
	}

	return manager.GenerateNginxConfig(".", cfg)
}
