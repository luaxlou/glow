package cmd

import (
	"fmt"
	"os"

	"github.com/luaxlou/glow/internal/configmanager"
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

	// 2. API Key
	apiKey, _ := configmanager.GetSystemConfig("api_key")
	if apiKey != "" {
		fmt.Printf("API Key: [CONFIGURED]\n")
	} else {
		fmt.Printf("API Key: [NOT FOUND]\n")
	}

	// 3. MySQL
	var mysqlInfo MySQLConfig
	if err := getSystemConfigJSON("mysql_info", &mysqlInfo); err == nil && mysqlInfo.Host != "" {
		fmt.Printf("MySQL: %s:%d (User: %s, Databases: %d)\n", mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.User, len(mysqlInfo.Databases))
	} else {
		fmt.Println("MySQL: [NOT CONFIGURED]")
	}

	// 4. Redis
	var redisInfo RedisConfig
	if err := getSystemConfigJSON("redis_info", &redisInfo); err == nil && redisInfo.Host != "" {
		fmt.Printf("Redis: %s:%d\n", redisInfo.Host, redisInfo.Port)
	} else {
		fmt.Println("Redis: [NOT CONFIGURED]")
	}

	// 5. Nginx
	var nginxInfo NginxSystemConfig
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
