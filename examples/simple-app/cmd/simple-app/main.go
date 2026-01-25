package main

import (
	"fmt"
	"log"

	"github.com/luaxlou/glow/starter/glowconfig"
)

func main() {
	// Load configuration
	config, err := glowconfig.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config.json: %v", err)
		log.Println("Using default values...")
		config = make(glowconfig.Config)
	}

	// Read configuration values
	logLevel := config.GetString("log_level")
	maxConnections := config.GetInt("max_connections")
	mysqlDSN := config.GetString("mysql_dsn")
	redisAddr := config.GetString("redis_addr")

	// Print configuration
	fmt.Println("=== Glow Simple App ===")
	fmt.Printf("Log Level: %s\n", logLevel)
	fmt.Printf("Max Connections: %d\n", maxConnections)
	fmt.Printf("MySQL DSN: %s\n", mysqlDSN)
	fmt.Printf("Redis Addr: %s\n", redisAddr)

	// Check if running in Glow environment
	port := fmt.Sprintf("%s", config.Get("port"))
	if port != "" {
		fmt.Printf("Port from config: %s\n", port)
	}

	fmt.Println("=== Application Started ===")
	fmt.Println("Hello, World!")
}
