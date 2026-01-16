package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/luaxlou/glow/pkg/manifest"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by filename or stdin",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("filename")
		if path == "" {
			fmt.Println("Error: must specify -f")
			return
		}

		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Error accessing path: %v\n", err)
			return
		}

		if info.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Printf("Error reading directory: %v\n", err)
				return
			}
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				ext := strings.ToLower(filepath.Ext(f.Name()))
				if ext == ".yaml" || ext == ".yml" || ext == ".json" {
					processFile(filepath.Join(path, f.Name()))
				}
			}
		} else {
			processFile(path)
		}
	},
}

func processFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	docs, err := manifest.Parse(data)
	if err != nil {
		fmt.Printf("Error parsing manifest %s: %v\n", filename, err)
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
		case api.Config:
			fmt.Printf("Applying Config: %s\n", obj.Name)
			// Config uses PUT /config/<appName>
			var resp api.Response
			if err := request("PUT", "/config/"+obj.Name, obj.Data, &resp); err != nil {
				fmt.Printf("Failed to apply config: %v\n", err)
			} else if !resp.Success {
				fmt.Printf("Error applying config: %s\n", resp.Message)
			} else {
				fmt.Printf("Success: %s\n", "Config updated")
			}
		case api.Ingress:
			fmt.Printf("Applying Ingress: %s\n", obj.Name)
			// Ingress uses POST /ingress/update
			appName := obj.Spec.Service
			if appName == "" {
				appName = obj.Name
			}
			req := api.IngressUpdateRequest{
				AppName: appName,
				Domain:  obj.Spec.Domain,
				Port:    obj.Spec.Port,
			}
			applyResource("/ingress/update", req)
		default:
			fmt.Printf("Unknown resource type in %s: %T\n", filename, obj)
		}
	}
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
