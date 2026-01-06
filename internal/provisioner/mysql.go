package provisioner

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luaxlou/glow/pkg/api"
)

type MySQL struct {
	AdminUser     string
	AdminPassword string
	Port          int
}

func NewMySQL(spec api.ServiceSpec) *MySQL {
	user := spec.AdminUser
	if user == "" {
		user = "root"
	}
	return &MySQL{
		AdminUser:     user,
		AdminPassword: spec.AdminPassword,
		Port:          spec.Port,
	}
}

func (m *MySQL) Check() error {
	// Try to connect
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/", m.AdminUser, m.AdminPassword, m.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Second * 5)
	return db.Ping()
}

func (m *MySQL) Provision(dbName string) (string, string, error) {
	// Returns username, password, error
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/", m.AdminUser, m.AdminPassword, m.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	// Create DB
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		return "", "", fmt.Errorf("failed to create db: %w", err)
	}

	// Create User
	appUser := dbName + "_user"
	appPass, _ := generateRandomKey(16)

	// Note: In real world, check if user exists or update password
	query := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", appUser, appPass)
	if _, err := db.Exec(query); err != nil {
		// Try altering if exists? For now, just error or ignore
		// return "", "", fmt.Errorf("failed to create user: %w", err)
	}

	// Grant
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%'", dbName, appUser))
	if err != nil {
		return "", "", fmt.Errorf("failed to grant privileges: %w", err)
	}

	return appUser, appPass, nil
}

func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
