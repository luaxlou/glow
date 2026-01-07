package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

var getNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "List nodes",
	Run: func(cmd *cobra.Command, args []string) {
		var resp api.Response
		if err := request("GET", "/node/status", nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Handle list or single object
		data, _ := json.Marshal(resp.Data)
		var nodes []api.Node
		
		// Try to unmarshal as list
		if err := json.Unmarshal(data, &nodes); err != nil {
			// Try as single object
			var node api.Node
			if err := json.Unmarshal(data, &node); err == nil {
				nodes = append(nodes, node)
			}
		}

		headers := []string{"NAME", "STATUS", "CPU%", "MEM%", "DISK%"}
		var rows [][]string
		for _, node := range nodes {
			rows = append(rows, []string{
				node.Name,
				"Ready", // Simplified status
				fmt.Sprintf("%.1f%%", node.Status.CPUUsage),
				fmt.Sprintf("%.1f%%", node.Status.MemUsage),
				fmt.Sprintf("%.1f%%", node.Status.DiskUsage),
			})
		}
		printTable(headers, rows)
	},
}

var describeNodeCmd = &cobra.Command{
	Use:   "node [name]",
	Short: "Describe node details",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// If name provided, filter? For single-node MVP, we ignore name or check match
		var resp api.Response
		if err := request("GET", "/node/status", nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		data, _ := json.Marshal(resp.Data)
		var nodes []api.Node
		if err := json.Unmarshal(data, &nodes); err != nil {
			var node api.Node
			json.Unmarshal(data, &node)
			nodes = append(nodes, node)
		}
		
		// If name arg present, find it
		targetName := ""
		if len(args) > 0 {
			targetName = args[0]
		}
		
		for _, node := range nodes {
			if targetName != "" && node.Name != targetName && node.Status.Hostname != targetName {
				continue
			}
			
			fmt.Printf("Name:     %s\n", node.Name)
			fmt.Printf("Hostname: %s\n", node.Status.Hostname)
			fmt.Printf("OS/Arch:  %s/%s\n", node.Status.OS, node.Status.Arch)
			fmt.Printf("Kernel:   %s\n", node.Status.Kernel)
			fmt.Println("System Info:")
			fmt.Printf("  CPU Usage:    %.2f%%\n", node.Status.CPUUsage)
			fmt.Printf("  Mem Usage:    %.2f%%\n", node.Status.MemUsage)
			fmt.Printf("  Disk Usage:   %.2f%%\n", node.Status.DiskUsage)
			
			if len(node.Status.Resources) > 0 {
				fmt.Println("Managed Resources:")
				for _, res := range node.Status.Resources {
					fmt.Printf("  - %s (%s): Port %d\n", res.Name, res.Kind, res.Port)
				}
			}
			fmt.Println("---")
		}
	},
}

func init() {
	getCmd.AddCommand(getNodeCmd)
	describeCmd.AddCommand(describeNodeCmd)
}
