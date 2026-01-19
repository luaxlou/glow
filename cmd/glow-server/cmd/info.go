package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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

	fmt.Printf("PID: %d\n", os.Getpid())

	fmt.Println()
	fmt.Println("Managed Resources")
	fmt.Println("-----------------")

	var mysqlInfo api.MySQLConfig
	if err := getSystemConfigJSON("mysql_info", &mysqlInfo); err == nil && mysqlInfo.Host != "" {
		fmt.Println("MySQL:")
		fmt.Printf("  Host: %s\n", mysqlInfo.Host)
		fmt.Printf("  Port: %d\n", mysqlInfo.Port)
		fmt.Printf("  Root User: %s\n", mysqlInfo.User)
		fmt.Printf("  Root Password: %s\n", mysqlInfo.Password)

		if len(mysqlInfo.Databases) > 0 {
			fmt.Println("  Databases:")
			for _, db := range mysqlInfo.Databases {
				fmt.Printf("    - %s (charset=%s)\n", db.Name, db.Charset)
			}
		}

		b, _ := json.MarshalIndent(mysqlInfo, "", "  ")
		fmt.Println("  Raw Config:")
		fmt.Println(string(b))
	} else {
		fmt.Println("MySQL: [NOT CONFIGURED]")
	}

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
