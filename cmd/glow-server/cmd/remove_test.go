package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/starter/glowsqlite"
)

func TestRunRemoveResource_NotConfigured(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "glow.db")
	glowsqlite.Reload()
	glowsqlite.Init(dbPath)

	var out bytes.Buffer
	if err := runRemoveResource(bytes.NewBufferString("y\n"), &out, "mysql", false); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if got := out.String(); got == "" {
		t.Fatalf("expected output")
	}
}

func TestRunRemoveResource_YesSkipsConfirmation(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "glow.db")
	glowsqlite.Reload()
	glowsqlite.Init(dbPath)

	if err := configmanager.SetSystemConfig("mysql_info", `{"host":"127.0.0.1"}`); err != nil {
		t.Fatalf("set mysql_info: %v", err)
	}
	if err := configmanager.SetSystemConfig("mysql_users", `{"db1":{"user":"u","password":"p"}}`); err != nil {
		t.Fatalf("set mysql_users: %v", err)
	}

	var out bytes.Buffer
	if err := runRemoveResource(bytes.NewBufferString("n\n"), &out, "mysql", true); err != nil {
		t.Fatalf("remove: %v", err)
	}

	v, err := configmanager.GetSystemConfig("mysql_info")
	if err == nil && v != "" {
		t.Fatalf("expected mysql_info removed, got: %q", v)
	}
}

func TestRunRemoveResource_ConfirmYes(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "glow.db")
	glowsqlite.Reload()
	glowsqlite.Init(dbPath)

	if err := configmanager.SetSystemConfig("redis_info", `{"host":"127.0.0.1"}`); err != nil {
		t.Fatalf("set redis_info: %v", err)
	}

	var out bytes.Buffer
	if err := runRemoveResource(bytes.NewBufferString("y\n"), &out, "redis", false); err != nil {
		t.Fatalf("remove: %v", err)
	}

	v, err := configmanager.GetSystemConfig("redis_info")
	if err == nil && v != "" {
		t.Fatalf("expected redis_info removed, got: %q", v)
	}
}

func TestRunRemoveResource_ConfirmNo(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "glow.db")
	glowsqlite.Reload()
	glowsqlite.Init(dbPath)

	if err := configmanager.SetSystemConfig("nginx_info", `{"binary_path":"/usr/sbin/nginx"}`); err != nil {
		t.Fatalf("set nginx_info: %v", err)
	}

	var out bytes.Buffer
	if err := runRemoveResource(bytes.NewBufferString("n\n"), &out, "nginx", false); err != nil {
		t.Fatalf("remove: %v", err)
	}

	v, err := configmanager.GetSystemConfig("nginx_info")
	if err != nil || v == "" {
		t.Fatalf("expected nginx_info still present")
	}
}
