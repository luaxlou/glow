package cmd

import (
	"database/sql"

	"github.com/luaxlou/glow/internal/configmanager"
)

// Wrapper functions to maintain backward compatibility with existing commands
// while delegating to the new centralized configmanager.

func getSystemConfig(key string) (string, error) {
	val, err := configmanager.GetSystemConfig(key)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func getSystemConfigJSON(key string, v interface{}) error {
	return configmanager.GetSystemConfigJSON(key, v)
}

func setSystemConfig(key, value string) error {
	return configmanager.SetSystemConfig(key, value)
}

func setSystemConfigJSON(key string, v interface{}) error {
	return configmanager.SetSystemConfigJSON(key, v)
}
