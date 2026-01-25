package glowconfig

import (
	"encoding/json"
	"os"
)

// Config represents the generic configuration map.
type Config map[string]interface{}

// Load reads config.json from the current working directory.
func Load() (Config, error) {
	return LoadFromFile("config.json")
}

// LoadFromFile reads configuration from a specific file path.
func LoadFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// Get returns a value from the config map using dot notation for nested keys.
func (c Config) Get(key string) interface{} {
	return c.getNested(key)
}

func (c Config) getNested(key string) interface{} {
	keys := splitKey(key)
	var current interface{} = c

	for _, k := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			// Try Config type alias
			if cm, ok := current.(Config); ok {
				m = map[string]interface{}(cm)
			} else {
				return nil
			}
		}
		val, ok := m[k]
		if !ok {
			return nil
		}
		current = val
	}
	return current
}

func splitKey(key string) []string {
	// Simple split by dot, but we can't import strings due to minimal imports?
	// We can import strings.
	var parts []string
	start := 0
	for i := 0; i < len(key); i++ {
		if key[i] == '.' {
			parts = append(parts, key[start:i])
			start = i + 1
		}
	}
	parts = append(parts, key[start:])
	return parts
}

// GetString returns a string value or empty string.
func (c Config) GetString(key string) string {
	val := c.Get(key)
	if v, ok := val.(string); ok {
		return v
	}
	return ""
}

// GetInt returns an int value or 0.
func (c Config) GetInt(key string) int {
	val := c.Get(key)
	if v, ok := val.(float64); ok {
		return int(v)
	}
	if v, ok := val.(int); ok {
		return v
	}
	return 0
}
