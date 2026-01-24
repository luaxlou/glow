package glowconfig

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 1. Create a dummy config.json
	configContent := map[string]interface{}{
		"foo": "bar",
		"port": 8080,
	}
	data, _ := json.Marshal(configContent)
	if err := os.WriteFile("config.json", data, 0644); err != nil {
		t.Fatalf("Failed to create config.json: %v", err)
	}
	defer os.Remove("config.json")

	// 2. Load it
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// 3. Verify values
	if cfg.GetString("foo") != "bar" {
		t.Errorf("Expected foo=bar, got %v", cfg.Get("foo"))
	}
	if cfg.GetInt("port") != 8080 {
		t.Errorf("Expected port=8080, got %v", cfg.Get("port"))
	}
}
