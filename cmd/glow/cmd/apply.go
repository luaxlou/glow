package cmd

import (
	"fmt"
	"os"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/luaxlou/glow/pkg/manifest"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by filename or stdin",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("filename")
		if file == "" {
			fmt.Println("Error: must specify -f")
			return
		}

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		docs, err := manifest.Parse(data)
		if err != nil {
			fmt.Printf("Error parsing manifest: %v\n", err)
			return
		}

		for _, doc := range docs {
			switch obj := doc.(type) {
			case api.Host:
				fmt.Printf("Applying Node: %s\n", obj.Metadata.Name)
				applyResource("/apply/host", obj)
			case api.App:
				fmt.Printf("Applying Deployment: %s\n", obj.Metadata.Name)
				applyResource("/apply/app", obj)
			default:
				fmt.Printf("Unknown resource type: %T\n", obj)
			}
		}
	},
}

func applyResource(endpoint string, body any) {
	var resp api.Response
	if err := request("POST", endpoint, body, &resp); err != nil {
		fmt.Printf("Failed to apply: %v\n", err)
		return
	}
	if !resp.Success {
		fmt.Printf("Error: %s\n", resp.Message)
	} else {
		fmt.Printf("Success: %s\n", resp.Message)
	}
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("filename", "f", "", "Filename to use to create the resource")
}
