package appcenter

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/luaxlou/glow/pkg/api"
)

func TestSendConfigUpdate(t *testing.T) {
	// 1. Setup pipe to simulate connection
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	appName := "test-app-config-update"
	
	// 2. Register the app
	RegisterActiveApp(api.AppInfo{Name: appName}, serverConn)
	defer UnregisterActiveApp(appName)

	// 3. Prepare config to send
	config := map[string]any{
		"key1": "value1",
		"key2": 123.0,
	}

	// 4. Send update in a goroutine (because writing might block if pipe is full, though here it's small)
	errCh := make(chan error, 1)
	go func() {
		errCh <- SendConfigUpdate(appName, config)
	}()

	// 5. Read from clientConn
	decoder := json.NewDecoder(clientConn)
	var resp api.Response
	
	// Set read deadline to avoid hanging
	clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// 6. Verify result
	if err := <-errCh; err != nil {
		t.Fatalf("SendConfigUpdate failed: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success true, got false")
	}

	// Verify Data content
	// JSON unmarshal to map[string]any makes numbers float64
	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be map[string]any, got %T", resp.Data)
	}

	if data["key1"] != "value1" {
		t.Errorf("Expected key1=value1, got %v", data["key1"])
	}
	if data["key2"] != 123.0 {
		t.Errorf("Expected key2=123.0, got %v", data["key2"])
	}
}

func TestSendConfigUpdate_AppNotConnected(t *testing.T) {
	err := SendConfigUpdate("non-existent-app", map[string]any{})
	if err == nil {
		t.Error("Expected error for non-existent app, got nil")
	}
}
