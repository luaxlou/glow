package manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNginxManager(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := NginxConfig{
		Name:   "test-ingress",
		Port:   8080,
		Domain: "test.example.com",
	}

	// 1. Generate Config
	// We pass dataDir so it writes to tmpDir/nginx
	if err := GenerateNginxConfig(tmpDir, cfg); err != nil {
		t.Logf("GenerateNginxConfig returned error (nginx may be unavailable in test env): %v", err)
	}

	// Verify file exists
	confFile := filepath.Join(tmpDir, "nginx", "test-ingress.conf")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		t.Fatalf("Nginx config file not created")
	}

	// 2. Get Ingress
	got, err := GetIngress(tmpDir, "test-ingress")
	if err != nil {
		t.Fatalf("GetIngress failed: %v", err)
	}
	if got.Name != cfg.Name || got.Port != cfg.Port || got.Domain != cfg.Domain {
		t.Errorf("GetIngress returned mismatch: %+v, expected %+v", got, cfg)
	}

	// 3. List Ingress
	list, err := ListIngress(tmpDir)
	if err != nil {
		t.Fatalf("ListIngress failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("ListIngress length mismatch: got %d, expected 1", len(list))
	} else if list[0].Name != cfg.Name {
		t.Errorf("ListIngress content mismatch: got %s, expected %s", list[0].Name, cfg.Name)
	}

	// 4. Remove Ingress
	if err := RemoveNginxConfig(tmpDir, "test-ingress"); err != nil {
		t.Fatalf("RemoveNginxConfig failed: %v", err)
	}

	// Verify file removed
	if _, err := os.Stat(confFile); !os.IsNotExist(err) {
		t.Errorf("Nginx config file not removed")
	}
}
