package appcenter

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/internal/statemanager"
	"github.com/luaxlou/glow/pkg/api"
)

var (
// activeConns sync.Map // map[string]net.Conn (AppName -> Conn) -- Replaced by store.go
)

func Start(port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Printf("App Center TCP Server listening on :%d", port)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept error: %v", err)
				continue
			}
			go handleConnection(conn)
		}
	}()
	return nil

}

func Close() {

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var connectedAppName string
	defer func() {
		if connectedAppName != "" {
			UnregisterActiveApp(connectedAppName)
		}
	}()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	// Keep connection open for multiple requests?
	// Or one request per connection?
	// MQ style usually keeps connection open.
	for {
		var req api.TCPRequest
		if err := decoder.Decode(&req); err != nil {
			// EOF is expected when client closes
			return
		}

		// Security Check: Deny any request attempting to act as or query "glow-server"
		if req.AppName == "glow-server" {
			log.Printf("Security Alert: Attempt to access reserved app name 'glow-server' via TCP from %s", conn.RemoteAddr())
			resp := api.Response{Success: false, Message: "Access denied: 'glow-server' is a reserved system application."}
			encoder.Encode(resp)
			return // Close connection immediately
		}

		var resp api.Response
		switch req.Action {
		case api.ActionGetConfig:
			// Deprecated active get, but keeping for compatibility or direct config fetch
			config, err := configmanager.Get(req.AppName)
			if err != nil {
				// Not found is not necessarily an error to log loudly, but we return false
				resp = api.Response{Success: false, Message: err.Error()}
			} else {
				// No need to sanitize here if we strictly enforce AppName != glow-server
				// and assume app configs don't contain sensitive system keys.
				// But sanitizing is still a good defense-in-depth if we had other shared keys.
				// For now, based on user instruction "not a filter issue", we return as is.
				resp = api.Response{Success: true, Data: config}
			}

		case api.ActionRegister:
			var appInfo api.AppInfo
			if err := json.Unmarshal(req.Payload, &appInfo); err != nil {
				resp = api.Response{Success: false, Message: "Invalid payload"}
			} else {
				// Additional check inside payload
				if appInfo.Name == "glow-server" {
					resp = api.Response{Success: false, Message: "Registration denied: 'glow-server' is reserved."}
				} else {
					// We use store directly to save/update app info
					// This allows external apps to register themselves
					if err := statemanager.SaveApp(appInfo); err != nil {
						resp = api.Response{Success: false, Message: err.Error()}
					} else {
						resp = api.Response{Success: true, Message: "Registered"}
					}
				}
			}

		case api.ActionAppStart:
			var appInfo api.AppInfo
			if err := json.Unmarshal(req.Payload, &appInfo); err != nil {
				resp = api.Response{Success: false, Message: "Invalid payload"}
			} else {
				// Additional check inside payload
				if appInfo.Name == "glow-server" {
					resp = api.Response{Success: false, Message: "Start denied: 'glow-server' is reserved."}
				} else {
					connectedAppName = appInfo.Name
					RegisterActiveApp(appInfo, conn)

					// 1. Record state (memory & sqlite)
					if err := statemanager.SaveApp(appInfo); err != nil {
						log.Printf("Error saving app state for %s: %v", appInfo.Name, err)
						// We continue even if save fails, to try to return config
					}

					// 2. Get Config from SQLite (via configmanager)
					config, err := configmanager.Get(appInfo.Name)
					if err != nil {
						// Config might not exist yet, which is fine
						log.Printf("No config found for app %s: %v", appInfo.Name, err)
						// We send empty success so client knows we handled the start
						resp = api.Response{Success: true, Message: "App started, no config found", Data: nil}
					} else {
						// 3. Pass config to application
						// Direct return, no sanitization needed as per new requirement
						resp = api.Response{Success: true, Data: config}
					}
				}
			}

		case api.ActionProvision:
			resp = HandleProvision(req)

		default:
			resp = api.Response{Success: false, Message: "Unknown action"}
		}

		if err := encoder.Encode(resp); err != nil {
			log.Printf("Encode error: %v", err)
			return
		}
	}
}
