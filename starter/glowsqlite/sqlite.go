package glowsqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/luaxlou/glow/starter/glowapp"
	_ "modernc.org/sqlite"
)

var (
	db          *sql.DB
	initialized bool
	mu          sync.RWMutex
	dbName      string
	localPath   string
	schemas     []string
)

// Init configures the SQLite connection.
// This must be called before DB().
func Init(path string) {
	mu.Lock()
	defer mu.Unlock()
	localPath = path
}

// ensureDefaults sets default values if they haven't been configured.
func ensureDefaults() {
	if dbName == "" {
		dbName = "glow"
	}
}

// RegisterSchema adds SQL statements to be executed when the connection is established.
// Useful for creating tables.
func RegisterSchema(sql string) {
	schemas = append(schemas, sql)
}

// DB returns a singleton SQLite connection.
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

	ensureDefaults()

	var dbFile string

	if localPath != "" {
		dbFile = localPath
	} else {
		dbFile = fmt.Sprintf("./%s.db", dbName)
	}

	// Ensure directory exists if it has a path
	dir := filepath.Dir(dbFile)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for sqlite: %w", err)
		}
	}

	if dbFile == "" {
		return nil, fmt.Errorf("sqlite file path not found")
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}

	// Apply schemas
	for _, schema := range schemas {
		if _, err := conn.Exec(schema); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to apply schema: %w", err)
		}
	}

	db = conn
	initialized = true
	log.Println("SQLite Starter initialized successfully.")

	glowapp.RegisterCleanup("SQLite Starter", func() {
		if db != nil {
			log.Println("Closing SQLite connection...")
			db.Close()
		}
	})

	return db, nil
}

// Reload forces re-initialization of the SQLite connection.
func Reload() {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		db.Close()
		db = nil
	}
	initialized = false
	log.Println("SQLite Starter reset for re-initialization.")
}
