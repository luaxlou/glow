package glowmysql

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/luaxlou/glow/starter/glowapp"
	"github.com/luaxlou/glow/starter/glowapp/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	gdb         *gorm.DB
	db          *sql.DB
	initialized bool
	mu          sync.RWMutex
	dbName      string
)

func Init(name string) {
	dbName = name
}

func Gorm() (*gorm.DB, error) {
	mu.RLock()
	if initialized && gdb != nil {
		defer mu.RUnlock()
		return gdb, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if initialized && gdb != nil {
		return gdb, nil
	}

	appName := config.AppIdentity
	if appName == "" {
		return nil, fmt.Errorf("app identity not set. call app.Init() first")
	}

	if dbName == "" {
		return nil, fmt.Errorf("mysql db name not configured. use mysql.Init(dbName)")
	}

	log.Printf("Lazy initializing MySQL Starter for %s (db: %s)...", appName, dbName)

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

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection via gorm: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get mysql sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	gdb = conn
	db = sqlDB
	initialized = true
	log.Println("MySQL Starter initialized successfully.")

	glowapp.RegisterCleanup("MySQL Starter", func() {
		if db != nil {
			log.Println("Closing MySQL connection...")
			db.Close()
		}
	})

	return gdb, nil
}

func DB() (*sql.DB, error) {
	mu.RLock()
	if initialized && db != nil {
		defer mu.RUnlock()
		return db, nil
	}
	mu.RUnlock()

	conn, err := Gorm()
	if err != nil {
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get mysql sql.DB: %w", err)
	}
	return sqlDB, nil
}

func Reload() {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		db.Close()
		db = nil
	}
	gdb = nil
	initialized = false
	log.Println("MySQL Starter reset for re-initialization.")
}
