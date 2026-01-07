package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

var getResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "List managed resources",
	Run: func(cmd *cobra.Command, args []string) {
		var resp api.Response
		if err := request("GET", "/resources/list", nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		data, _ := json.Marshal(resp.Data)
		var resources []api.ResourceRef
		json.Unmarshal(data, &resources)

		headers := []string{"KIND", "NAME", "PORT"}
		var rows [][]string
		for _, res := range resources {
			rows = append(rows, []string{
				res.Kind,
				res.Name,
				strconv.Itoa(res.Port),
			})
		}
		printTable(headers, rows)
	},
}

// Generic Describe (Hooked into verbs.go via init or similar)
// For now, standalone command or try to hook describeCmd?
// In node.go I added `describeNodeCmd`.
// Here I add `describeCmd` fallback logic?
// Cobra executes checking subcommands. If none match, it runs `describeCmd.Run`.
// So if I set `describeCmd.Run` in `verbs.go` or here?
// `verbs.go` initialized `describeCmd`.
// I can set `describeCmd.Run` here.

func init() {
	getCmd.AddCommand(getResourcesCmd)
	
	// Generic Describe Logic
	describeCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		name := args[0]
		
		// 1. Try Resources
		var resp api.Response
		request("GET", "/resources/list", nil, &resp)
		data, _ := json.Marshal(resp.Data)
		var resources []api.ResourceRef
		json.Unmarshal(data, &resources)
		
		for _, res := range resources {
			if res.Name == name {
				fmt.Printf("Resource: %s\nKind:     %s\nPort:     %d\n", res.Name, res.Kind, res.Port)
				return
			}
		}
		
		// 2. Try Deployments
		request("GET", "/apps/list", nil, &resp)
		data, _ = json.Marshal(resp.Data)
		var apps []api.AppInfo
		json.Unmarshal(data, &apps)
		for _, app := range apps {
			if app.Name == name {
				fmt.Printf("Deployment: %s\nStatus:     %s\nPID:        %d\n", app.Name, app.Status, app.Pid)
				return
			}
		}
		
		fmt.Printf("Resource %s not found\n", name)
	}
}
