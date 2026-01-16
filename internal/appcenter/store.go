package appcenter

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/luaxlou/glow/pkg/api"
)

type RuntimeApp struct {
	Info          api.AppInfo
	Conn          net.Conn
	ConnectedAt   time.Time
	LastHeartbeat time.Time
}

var (
	// store holds the in-memory state of active applications.
	store = struct {
		sync.RWMutex
		apps map[string]*RuntimeApp
	}{
		apps: make(map[string]*RuntimeApp),
	}
)

// RegisterActiveApp adds or updates an active application in the memory store.
func RegisterActiveApp(info api.AppInfo, conn net.Conn) {
	store.Lock()
	defer store.Unlock()

	store.apps[info.Name] = &RuntimeApp{
		Info:          info,
		Conn:          conn,
		ConnectedAt:   time.Now(),
		LastHeartbeat: time.Now(),
	}
}

// UnregisterActiveApp removes an application from the memory store.
func UnregisterActiveApp(name string) {
	store.Lock()
	defer store.Unlock()
	delete(store.apps, name)
}

// GetActiveApps returns a snapshot of all currently connected applications.
func GetActiveApps() map[string]api.AppInfo {
	store.RLock()
	defer store.RUnlock()

	result := make(map[string]api.AppInfo)
	for name, runtime := range store.apps {
		result[name] = runtime.Info
	}
	return result
}

// GetActiveApp returns the info for a specific connected application.
func GetActiveApp(name string) (api.AppInfo, bool) {
	store.RLock()
	defer store.RUnlock()

	if runtime, ok := store.apps[name]; ok {
		return runtime.Info, true
	}
	return api.AppInfo{}, false
}

// SendConfigUpdate sends the new configuration to the connected application.
func SendConfigUpdate(appName string, config map[string]any) error {
	store.RLock()
	runtime, ok := store.apps[appName]
	store.RUnlock()

	if !ok {
		return fmt.Errorf("app %s not connected", appName)
	}

	resp := api.Response{
		Success: true,
		Data:    config,
	}

	// Create a new encoder to send the message
	// Note: We assume single-writer safety here because the client
	// on this connection is only listening (monitorConfig), so
	// handleConnection on server is blocked on Read and won't Write.
	return json.NewEncoder(runtime.Conn).Encode(resp)
}
