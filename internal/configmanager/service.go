package configmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/luaxlou/glow/starter/glowsqlite"
)

var (
	once  sync.Once
	mu    sync.Mutex
	cache map[string]any
)

const schema = `
	CREATE TABLE IF NOT EXISTS app_configs (
		app_name TEXT PRIMARY KEY,
		config_json TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

func init() {
	glowsqlite.RegisterSchema(schema)
}

func Init() {
	once.Do(func() {
		cache = make(map[string]any)
	})
}

// EnsureInitialized ensures that the service is initialized and cache is loaded.
func EnsureInitialized() error {
	Init()

	mu.Lock()
	defer mu.Unlock()

	if len(cache) > 0 {
		return nil
	}

	return loadCache()
}

func loadCache() error {
	db, err := glowsqlite.DB()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT app_name, config_json FROM app_configs")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var appName, configJSON string
		if err := rows.Scan(&appName, &configJSON); err != nil {
			return err
		}

		var config map[string]any
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			log.Printf("Error unmarshaling config for app %s: %v", appName, err)
			continue
		}
		cache[appName] = config
	}
	return nil
}

// Get returns the configuration for a given app.
func Get(appName string) (map[string]any, error) {
	if err := EnsureInitialized(); err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	if config, ok := cache[appName]; ok {
		return config.(map[string]any), nil
	}
	return nil, fmt.Errorf("config not found for app: %s", appName)
}

// GetValue retrieves a specific key from an app's configuration.
func GetValue(appName, key string) (any, bool) {
	cfg, err := Get(appName)
	if err != nil {
		return nil, false
	}
	val, ok := cfg[key]
	return val, ok
}

// Set updates the configuration for a given app.
// It merges with existing config if merge is true, otherwise overwrites.
func Set(appName string, newConfig map[string]any, merge bool) error {
	if err := EnsureInitialized(); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	var finalConfig map[string]any

	if merge {
		finalConfig = make(map[string]any)
		// Start with existing
		if existing, ok := cache[appName]; ok {
			if existingMap, ok := existing.(map[string]any); ok {
				for k, v := range existingMap {
					finalConfig[k] = v
				}
			}
		}
		// Merge new
		for k, v := range newConfig {
			finalConfig[k] = v
		}
	} else {
		finalConfig = newConfig
	}

	configJSON, err := json.Marshal(finalConfig)
	if err != nil {
		return err
	}

	db, err := glowsqlite.DB()
	if err != nil {
		return err
	}

	// Upsert into DB
	query := `
	INSERT INTO app_configs (app_name, config_json, updated_at) 
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(app_name) DO UPDATE SET 
		config_json = excluded.config_json,
		updated_at = CURRENT_TIMESTAMP;
	`
	if _, err := db.Exec(query, appName, string(configJSON)); err != nil {
		return err
	}

	// Update Cache
	cache[appName] = finalConfig
	return nil
}
