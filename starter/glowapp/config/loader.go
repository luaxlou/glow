package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/viper"
)

const (
	EnvConfig    = "OP_APP_CONFIG"
	EnvServerURL = "OP_SERVER_URL"
	EnvAppName   = "OP_APP_NAME"
)

var (
	v    *viper.Viper
	once sync.Once

	// AppIdentity holds the manually set app name if EnvAppName is missing
	AppIdentity string
)

func checkVerbose() bool {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--verbose" {
			return true
		}
	}
	return false
}

func verboseLog(format string, v ...any) {
	if checkVerbose() {
		log.Printf("[VERBOSE] "+format, v...)
	}
}

func getLocalConfigFile() string {
	appName := os.Getenv(EnvAppName)
	if appName == "" {
		appName = AppIdentity
	}
	if appName != "" {
		return fmt.Sprintf("%s_local_config.json", appName)
	}
	return "local_config.json"
}

func load() {
	once.Do(func() {
		v = viper.New()
		v.SetConfigType("json")

		localConfigFile := getLocalConfigFile()
		// Try loading from local config file
		if _, err := os.Stat(localConfigFile); err == nil {
			f, err := os.Open(localConfigFile)
			if err == nil {
				defer f.Close()
				if err := v.MergeConfig(f); err != nil {
					log.Printf("Warning: Failed to parse local config file: %v", err)
				} else {
					log.Printf("Loaded config from %s", localConfigFile)
					if checkVerbose() {
						allSettings := v.AllSettings()
						jsonSettings, _ := json.MarshalIndent(allSettings, "", "  ")
						verboseLog("Initial Local Config Data: \n%s", string(jsonSettings))
					}
				}
			}
		}
	})
}

// Start initiates the connection to AppCenter, sends AppInfo, and handles config.
// It blocks for up to 1 second waiting for the initial config.
func Start(appInfo api.AppInfo) {
	// Ensure viper is initialized (loads local config if available)
	load()

	serverAddr := os.Getenv(EnvServerURL)
	if serverAddr == "" {
		serverAddr = "127.0.0.1:32101"
	}
	serverAddr = strings.TrimPrefix(serverAddr, "http://")
	serverAddr = strings.TrimPrefix(serverAddr, "tcp://")

	verboseLog("Connecting to AppCenter at %s...", serverAddr)

	conn, err := net.DialTimeout("tcp", serverAddr, 2*time.Second)
	if err != nil {
		log.Printf("Warning: Failed to connect to AppCenter (%s): %v. Using local config.", serverAddr, err)
		return
	}

	// Prepare Start Request
	payload, _ := json.Marshal(appInfo)
	req := api.TCPRequest{
		Action:  api.ActionAppStart,
		AppName: appInfo.Name,
		Payload: payload,
	}

	verboseLog("Sending Start Request for app: %s", appInfo.Name)

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		log.Printf("Warning: Failed to send start request: %v. Using local config.", err)
		conn.Close()
		return
	}

	// Wait for initial config with timeout
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var resp api.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		// Timeout or error
		log.Printf("Warning: No config received from AppCenter (timeout 1s) or error: %v. Using local config.", err)
		// Clear deadline for future reads? Or close and rely on local?
		// Requirement: "if no config ... read local config execute startup".
		// But implies if connection is good but slow, we just proceed.
		// Should we keep listening? "received new one then update".
		// If timeout occurred, we can still keep listening.
	} else {
		// Success
		if resp.Success && resp.Data != nil {
			verboseLog("Received initial config from AppCenter")
			applyConfig(resp.Data)
		} else {
			log.Printf("Warning: AppCenter returned no config: %s", resp.Message)
		}
	}

	// Reset deadline for persistent connection
	conn.SetReadDeadline(time.Time{})

	// Start background monitor
	go monitorConfig(conn)
}

func monitorConfig(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	for {
		var resp api.Response
		if err := decoder.Decode(&resp); err != nil {
			log.Printf("Connection to AppCenter lost: %v", err)
			return
		}

		if resp.Success && resp.Data != nil {
			log.Println("Received config update from AppCenter")
			applyConfig(resp.Data)
		}
	}
}

func applyConfig(data any) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling config data: %v", err)
		return
	}

	if checkVerbose() {
		// Pretty print the config data
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, dataBytes, "", "  "); err == nil {
			verboseLog("Applying Config Data: \n%s", prettyJSON.String())
		} else {
			verboseLog("Applying Config Data: %s", string(dataBytes))
		}
	}

	// Update Viper
	if err := v.MergeConfig(bytes.NewReader(dataBytes)); err != nil {
		log.Printf("Error merging config: %v", err)
	}

	// Write to local file
	localConfigFile := getLocalConfigFile()
	if err := os.WriteFile(localConfigFile, dataBytes, 0644); err != nil {
		log.Printf("Error writing local config file: %v", err)
	}
}

// Register is deprecated, use Start instead. Kept for backward compatibility if needed,
// but currently just wraps Start (though Start blocks).
// Actually, Register was used to send AppInfo. Start does that too.
// We can remove Register if we update the caller.
func Register(appInfo api.AppInfo) error {
	// For compatibility, we can just call Start, but Start doesn't return error.
	Start(appInfo)
	return nil
}

// Get unmarshals a specific key into the target structure.
// This works for both the JSON blob and individual environment variables.
func Get(key string, target interface{}) error {
	load()
	sub := v.Sub(key)
	if sub != nil {
		return sub.Unmarshal(target)
	}

	// If it's a simple value (not a nested object), Unmarshal the main viper
	return v.UnmarshalKey(key, target)
}

// IsSet checks if a configuration key exists.
func IsSet(key string) bool {
	load()
	return v.IsSet(key)
}
