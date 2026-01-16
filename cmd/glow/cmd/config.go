package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configurations",
}

var viewConfigCmd = &cobra.Command{
	Use:   "view <app>",
	Short: "View application configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		var resp api.Response
		if err := request("GET", "/config/"+appName, nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if !resp.Success {
			fmt.Printf("Error: %s\n", resp.Message)
			return
		}
		data, _ := json.MarshalIndent(resp.Data, "", "  ")
		fmt.Println(string(data))
	},
}

var editConfigCmd = &cobra.Command{
	Use:   "edit <app>",
	Short: "Edit configuration interactively",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		// 1. Get
		var resp api.Response
		request("GET", "/config/"+appName, nil, &resp)

		currentData := resp.Data
		if currentData == nil {
			currentData = make(map[string]any)
		}

		// 2. Temp File
		f, err := os.CreateTemp("", "glow-config-*.json")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer os.Remove(f.Name())

		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		encoder.Encode(currentData)
		f.Close()

		// 3. Editor
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		c := exec.Command(editor, f.Name())
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("Error running editor: %v\n", err)
			return
		}

		// 4. Read back
		updatedContent, _ := os.ReadFile(f.Name())
		var newData map[string]any
		if err := json.Unmarshal(updatedContent, &newData); err != nil {
			fmt.Printf("Invalid JSON: %v\n", err)
			return
		}

		// 5. Update
		var upResp api.Response
		if err := request("PUT", "/config/"+appName+"?merge=false", newData, &upResp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if upResp.Success {
			fmt.Printf("Config updated\n")
		} else {
			fmt.Printf("Error: %s\n", upResp.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(viewConfigCmd)
	configCmd.AddCommand(editConfigCmd)
}