package statemanager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/luaxlou/glow/starter/glowsqlite"
)

var (
	once sync.Once
)

const schema = `
	CREATE TABLE IF NOT EXISTS apps (
		name TEXT PRIMARY KEY,
		info_json TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

func init() {
	glowsqlite.RegisterSchema(schema)
}

// SaveApp saves or updates the application information.
func SaveApp(app api.AppInfo) error {
	infoJSON, err := json.Marshal(app)
	if err != nil {
		return err
	}

	db, err := glowsqlite.DB()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO apps (name, info_json, updated_at) 
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(name) DO UPDATE SET 
		info_json = excluded.info_json,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err = db.Exec(query, app.Name, string(infoJSON))
	return err
}

// GetApp retrieves application information by name.
func GetApp(name string) (*api.AppInfo, error) {
	db, err := glowsqlite.DB()
	if err != nil {
		return nil, err
	}

	var infoJSON string
	err = db.QueryRow("SELECT info_json FROM apps WHERE name = ?", name).Scan(&infoJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("app not found: %s", name)
		}
		return nil, err
	}

	var app api.AppInfo
	if err := json.Unmarshal([]byte(infoJSON), &app); err != nil {
		return nil, err
	}
	return &app, nil
}

// DeleteApp removes an application from the database.
func DeleteApp(name string) error {
	db, err := glowsqlite.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM apps WHERE name = ?", name)
	return err
}

// ListApps retrieves all applications.
func ListApps() ([]api.AppInfo, error) {
	db, err := glowsqlite.DB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT info_json FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []api.AppInfo
	for rows.Next() {
		var infoJSON string
		if err := rows.Scan(&infoJSON); err != nil {
			return nil, err
		}
		var app api.AppInfo
		if err := json.Unmarshal([]byte(infoJSON), &app); err != nil {
			continue
		}
		apps = append(apps, app)
	}
	return apps, nil
}
