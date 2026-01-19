package configmanager

import (
	"encoding/json"

	"github.com/luaxlou/glow/starter/glowsqlite"
)

const systemConfigSchema = `
CREATE TABLE IF NOT EXISTS system_config (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

func init() {
	glowsqlite.RegisterSchema(systemConfigSchema)
}

func GetSystemConfig(key string) (string, error) {
	db, err := glowsqlite.DB()
	if err != nil {
		return "", err
	}
	var value string
	err = db.QueryRow("SELECT value FROM system_config WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func GetSystemConfigJSON(key string, v interface{}) error {
	val, err := GetSystemConfig(key)
	if err != nil {
		return err
	}
	if val == "" {
		return nil
	}
	return json.Unmarshal([]byte(val), v)
}

func SetSystemConfig(key, value string) error {
	db, err := glowsqlite.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO system_config (key, value, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET 
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP;
	`, key, value)
	return err
}

func SetSystemConfigJSON(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return SetSystemConfig(key, string(b))
}

func DeleteSystemConfig(key string) error {
	db, err := glowsqlite.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM system_config WHERE key = ?", key)
	return err
}
