package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

// --- Get App ---
var getAppCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"apps", "application"},
	Short:   "List all applications",
	Run: func(cmd *cobra.Command, args []string) {
		var resp api.Response
		if err := request("GET", "/apps/list", nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		data, _ := json.Marshal(resp.Data)
		var apps []api.AppInfo
		json.Unmarshal(data, &apps)

		headers := []string{"NAME", "STATUS", "RESTARTS", "AGE", "CPU", "MEM", "PID", "PORT", "DOMAIN"}
		var rows [][]string
		for _, app := range apps {
			age := "-"
			if app.StartTime > 0 {
				age = time.Since(time.UnixMilli(app.StartTime)).Round(time.Second).String()
			}
			mem := fmt.Sprintf("%.1f MB", float64(app.Stats.MemoryUsage)/1024/1024)

			rows = append(rows, []string{
				app.Name,

				app.Status,
				strconv.Itoa(app.RestartCount),
				age,
				fmt.Sprintf("%.1f%%", app.Stats.CPUPercent),
				mem,
				strconv.Itoa(app.Pid),
				strconv.Itoa(app.Port),
				app.Domain,
			})
		}
		printTable(headers, rows)
	},
}

// --- Describe App ---
var describeAppCmd = &cobra.Command{
	Use:     "app [name]",
	Aliases: []string{"apps", "application"},
	Short:   "Show details of an application",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		var resp api.Response
		request("GET", "/apps/list", nil, &resp)
		data, _ := json.Marshal(resp.Data)
		var apps []api.AppInfo
		json.Unmarshal(data, &apps)

		var target *api.AppInfo
		for _, app := range apps {
			if app.Name == name {
				target = &app
				break
			}
		}

		if target == nil {
			fmt.Printf("App %s not found\n", name)
			return
		}

		fmt.Printf("Name:       %s\n", target.Name)
		fmt.Printf("Status:     %s\n", target.Status)
		fmt.Printf("PID:        %d\n", target.Pid)
		fmt.Printf("Port:       %d\n", target.Port)
		fmt.Printf("Domain:     %s\n", target.Domain)
		fmt.Printf("Restarts:   %d\n", target.RestartCount)
		fmt.Printf("Command:    %s\n", target.Command)
		fmt.Printf("Args:       %v\n", target.Args)
		fmt.Printf("WorkDir:    %s\n", target.WorkingDir)

		age := "-"
		if target.StartTime > 0 {
			age = time.Since(time.UnixMilli(target.StartTime)).Round(time.Second).String()
		}
		fmt.Printf("Age:        %s\n", age)

		fmt.Println("Resources:")
		fmt.Printf("  CPU:      %.2f%%\n", target.Stats.CPUPercent)
		fmt.Printf("  Memory:   %d bytes\n", target.Stats.MemoryUsage)

		fmt.Println("Config:")
		if target.Config != nil {
			configBytes, _ := json.MarshalIndent(target.Config, "  ", "  ")
			fmt.Println(string(configBytes))
		} else {
			fmt.Println("  (none)")
		}
	},
}

// --- Delete App ---
var deleteAppCmd = &cobra.Command{
	Use:     "app [name]",
	Aliases: []string{"apps", "application"},
	Short:   "Delete an application",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		req := map[string]string{"name": name}
		var resp api.Response
		if err := request("POST", "/apps/delete", req, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if resp.Success {
			fmt.Printf("app.apps/%s deleted\n", name)
		} else {
			fmt.Printf("Error: %s\n", resp.Message)
		}
	},
}

// --- Lifecycle Commands ---
func runStartApp(cmd *cobra.Command, args []string) {
	name := args[0]
	req := api.StartAppRequest{Name: name}
	var resp api.Response
	if err := request("POST", "/apps/start", req, &resp); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if resp.Success {
		fmt.Printf("app.apps/%s started\n", name)
	} else {
		fmt.Printf("Error: %s\n", resp.Message)
	}
}

func runStopApp(cmd *cobra.Command, args []string) {
	name := args[0]
	req := map[string]string{"name": name}
	var resp api.Response
	if err := request("POST", "/apps/stop", req, &resp); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if resp.Success {
		fmt.Printf("app.apps/%s stopped\n", name)
	} else {
		fmt.Printf("Error: %s\n", resp.Message)
	}
}

func runRestartApp(cmd *cobra.Command, args []string) {
	name := args[0]
	req := map[string]string{"name": name}
	var resp api.Response
	if err := request("POST", "/apps/restart", req, &resp); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if resp.Success {
		fmt.Printf("app.apps/%s restarted\n", name)
	} else {
		fmt.Printf("Error: %s\n", resp.Message)
	}
}

func init() {
	// glow get app
	getCmd.AddCommand(getAppCmd)
	// glow describe app
	describeCmd.AddCommand(describeAppCmd)
	// glow delete app
	deleteCmd.AddCommand(deleteAppCmd)

	// glow start
	startCmd.Use = "start [name]"
	startCmd.Short = "Start an application"
	startCmd.Args = cobra.ExactArgs(1)
	startCmd.Run = runStartApp

	// glow stop
	stopCmd.Use = "stop [name]"
	stopCmd.Short = "Stop an application"
	stopCmd.Args = cobra.ExactArgs(1)
	stopCmd.Run = runStopApp

	// glow restart
	restartCmd.Use = "restart [name]"
	restartCmd.Short = "Restart an application"
	restartCmd.Args = cobra.ExactArgs(1)
	restartCmd.Run = runRestartApp

	// glow logs <name> (Flat command, implement directly)
	logsCmd.Run = func(cmd *cobra.Command, args []string) {
		name := args[0]
		var resp api.Response
		if err := request("GET", "/apps/logs?name="+name, nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if resp.Success {
			fmt.Println(resp.Data)
		} else {
			fmt.Printf("Error: %s\n", resp.Message)
		}
	}
}
