package glowmysql

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luaxlou/glow/starter/glowapp"
	"github.com/luaxlou/glow/starter/glowapp/config"
)

var (
	db          *sql.DB
	initialized bool
	mu          sync.RWMutex
	dbName      string
)

func Init(name string) {
	dbName = name
}

// DB returns a singleton MySQL connection.
// It initializes the connection on the first call using the configuration
// associated with the current AppIdentity.
func DB() (*sql.DB, error) {
	mu.RLock()
	if initialized && db != nil {
		defer mu.RUnlock()
		return db, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	// Double check
	if initialized && db != nil {
		return db, nil
	}

	appName := config.AppIdentity
	if appName == "" {
		return nil, fmt.Errorf("app identity not set. call app.Init() first")
	}

	if dbName == "" {
		return nil, fmt.Errorf("mysql db name not configured. use mysql.Init(dbName)")
	}

	log.Printf("Lazy initializing MySQL Starter for %s (db: %s)...", appName, dbName)

	// This calls into sdk/config, which uses the registered AppIdentity
	cfg, err := config.ProvisionResource("mysql", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to provision mysql: %w", err)
	}

	var dsn string
	if mysqlCfg, ok := cfg["mysql"].(map[string]any); ok {
		if d, ok := mysqlCfg["dsn"].(string); ok {
			dsn = d
		}
	}

	if dsn == "" {
		return nil, fmt.Errorf("dsn not found in mysql config")
	}

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	db = conn
	initialized = true
	log.Println("MySQL Starter initialized successfully.")

	glowapp.RegisterCleanup("MySQL Starter", func() {
		if db != nil {
			log.Println("Closing MySQL connection...")
			db.Close()
		}
	})

	return db, nil
}

// Reload forces re-initialization of the MySQL connection.
// This can be called when configuration updates are received.
func Reload() {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		db.Close()
		db = nil
	}
	initialized = false
	log.Println("MySQL Starter reset for re-initialization.")
}
