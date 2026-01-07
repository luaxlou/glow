package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luaxlou/glow/internal/manager"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

// DatabaseInfo stores information about a database
type DatabaseInfo struct {
	Name    string `json:"name"`
	Charset string `json:"charset"`
}

// MySQLConfig stores the collected MySQL configuration
type MySQLConfig struct {
	Host      string         `json:"host"`
	Port      int            `json:"port"`
	User      string         `json:"user"`
	Password  string         `json:"password"`
	Databases []DatabaseInfo `json:"databases"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// RedisConfig stores the collected Redis configuration
type RedisConfig struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NginxSystemConfig stores the collected Nginx configuration
type NginxSystemConfig struct {
	BinaryPath string    `json:"binary_path"`
	ConfPath   string    `json:"conf_path"`
	Version    string    `json:"version"`
	UpdatedAt  time.Time `json:"updated_at"`
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to the system",
}

var addMysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Add a MySQL database",
	Run:   runAddMysql,
}

var addRedisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Add a Redis server",
	Run:   runAddRedis,
}

var addNginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "Add/Discover system Nginx",
	Run:   runAddNginx,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addMysqlCmd)
	addCmd.AddCommand(addRedisCmd)
	addCmd.AddCommand(addNginxCmd)
}

func runAddNginx(cmd *cobra.Command, args []string) {
	fmt.Println("Probing for system Nginx...")

	info, err := manager.DetectNginx()
	if err != nil {
		fmt.Printf("Error: Nginx detection failed: %v\n", err)
		return
	}

	fmt.Printf("Found Nginx %s\n", info.Version)
	fmt.Printf("  Binary: %s\n", info.BinaryPath)
	fmt.Printf("  Config: %s\n", info.ConfPath)
	fmt.Printf("  Prefix: %s\n", info.PrefixPath)
	fmt.Printf("  Error Log: %s\n", info.ErrorLog)
	fmt.Printf("  Access Log: %s\n", info.AccessLog)

	// Save to system config
	config := NginxSystemConfig{
		BinaryPath: info.BinaryPath,
		ConfPath:   info.ConfPath,
		Version:    info.Version,
		UpdatedAt:  time.Now(),
	}

	err = setSystemConfigJSON("nginx_info", config)
	if err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Println("Successfully saved Nginx info.")
}

func runAddRedis(cmd *cobra.Command, args []string) {
	fmt.Println("Probing for local Redis on port 6379...")

	var password string
	var rdb *redis.Client
	var err error

	// 0. Check for existing config
	var existingConfig RedisConfig
	if err := getSystemConfigJSON("redis_info", &existingConfig); err == nil && existingConfig.Host != "" {
		fmt.Println("Found existing configuration. Trying to connect with stored credentials...")
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", existingConfig.Host, existingConfig.Port),
			Password: existingConfig.Password,
			DB:       0,
		})

		if err := rdb.Ping(context.Background()).Err(); err == nil {
			fmt.Println("Successfully connected using stored credentials.")
			password = existingConfig.Password
		} else {
			fmt.Printf("Failed to connect with stored credentials: %v\n", err)
			rdb.Close()
			rdb = nil
		}
	}

	if rdb == nil {
		// 1. Detect port 6379
		conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 1*time.Second)
		if err != nil {
			fmt.Printf("Error: Could not connect to local port 6379: %v\n", err)
			return
		}
		conn.Close()
		fmt.Println("Port 6379 is open.")

		// 2. Try to connect with empty password
		rdb = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})

		err = rdb.Ping(context.Background()).Err()
		if err != nil {
			fmt.Println("Connection with empty password failed.")
			rdb.Close()

			// 3. Interactive password prompt
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Please enter Redis password: ")
			passInput, _ := reader.ReadString('\n')
			password = strings.TrimSpace(passInput)

			rdb = redis.NewClient(&redis.Options{
				Addr:     "127.0.0.1:6379",
				Password: password,
				DB:       0,
			})

			if err := rdb.Ping(context.Background()).Err(); err != nil {
				fmt.Printf("Failed to connect with provided password: %v\n", err)
				rdb.Close()
				return
			}
		}
	}
	defer rdb.Close()

	fmt.Println("Successfully connected to Redis.")

	// 4. Save to system_configs
	config := RedisConfig{
		Host:      "127.0.0.1",
		Port:      6379,
		Password:  password,
		UpdatedAt: time.Now(),
	}

	err = setSystemConfigJSON("redis_info", config)
	if err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Println("Successfully saved Redis info.")
}

func runAddMysql(cmd *cobra.Command, args []string) {
	fmt.Println("Probing for local MySQL on port 3306...")

	var password string
	var db *sql.DB
	var err error

	// 0. Check for existing config
	var existingConfig MySQLConfig
	if err := getSystemConfigJSON("mysql_info", &existingConfig); err == nil && existingConfig.Host != "" {
		fmt.Println("Found existing configuration. Trying to connect with stored credentials...")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", existingConfig.User, existingConfig.Password, existingConfig.Host, existingConfig.Port)
		tempDB, err := sql.Open("mysql", dsn)
		if err == nil {
			if err = tempDB.Ping(); err == nil {
				fmt.Println("Successfully connected using stored credentials.")
				password = existingConfig.Password
				db = tempDB
			} else {
				fmt.Printf("Failed to connect with stored credentials: %v\n", err)
				tempDB.Close()
			}
		}
	}

	if db == nil {
		// 1. Detect port 3306
		conn, err := net.DialTimeout("tcp", "127.0.0.1:3306", 1*time.Second)
		if err != nil {
			fmt.Printf("Error: Could not connect to local port 3306: %v\n", err)
			return
		}
		conn.Close()
		fmt.Println("Port 3306 is open.")

		// 2. Try to connect with root and empty password
		dsn := "root:@tcp(127.0.0.1:3306)/"
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("Error initializing database client: %v\n", err)
			return
		}

		// Try ping
		err = db.Ping()
		if err != nil {
			fmt.Println("Connection with empty password failed.")
			db.Close() // Close previous connection attempt

			// 3. Interactive password prompt
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Please enter MySQL root password: ")
			passInput, _ := reader.ReadString('\n')
			password = strings.TrimSpace(passInput)

			dsn = fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/", password)
			db, err = sql.Open("mysql", dsn)
			if err != nil {
				fmt.Printf("Error initializing database client: %v\n", err)
				return
			}

			err = db.Ping()
			if err != nil {
				fmt.Printf("Failed to connect with provided password: %v\n", err)
				db.Close()
				return
			}
		}
	}

	defer db.Close()

	fmt.Println("Successfully connected to MySQL.")

	// 4. Get info
	rows, err := db.Query("SELECT SCHEMA_NAME, DEFAULT_CHARACTER_SET_NAME FROM information_schema.SCHEMATA")
	if err != nil {
		fmt.Printf("Failed to query databases: %v\n", err)
		return
	}
	defer rows.Close()

	var dbs []DatabaseInfo
	for rows.Next() {
		var name, charset string
		if err := rows.Scan(&name, &charset); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		dbs = append(dbs, DatabaseInfo{Name: name, Charset: charset})
	}

	config := MySQLConfig{
		Host:      "127.0.0.1",
		Port:      3306,
		User:      "root",
		Password:  password,
		Databases: dbs,
		UpdatedAt: time.Now(),
	}

	// 5. Save to system_configs
	err = setSystemConfigJSON("mysql_info", config)
	if err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Printf("Successfully saved MySQL info for %d databases.\n", len(dbs))
}
