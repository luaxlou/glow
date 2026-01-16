package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

var getIngressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "List ingress rules",
	Run: func(cmd *cobra.Command, args []string) {
		var resp api.Response
		if err := request("GET", "/ingress/list", nil, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		data, _ := json.Marshal(resp.Data)
		
		type IngressItem struct {
			Name   string `json:"name"`
			Port   int    `json:"port"`
			Domain string `json:"domain"`
		}
		var items []IngressItem
		json.Unmarshal(data, &items)

		headers := []string{"NAME", "DOMAIN", "PORT"}
		var rows [][]string
		for _, item := range items {
			rows = append(rows, []string{
				item.Name,
				item.Domain,
				strconv.Itoa(item.Port),
			})
		}
		printTable(headers, rows)
	},
}

var deleteIngressCmd = &cobra.Command{
	Use:   "ingress <name>",
	Short: "Delete an ingress rule",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		req := api.IngressDeleteRequest{AppName: name}
		var resp api.Response
		if err := request("POST", "/ingress/delete", req, &resp); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if resp.Success {
			fmt.Printf("Ingress %s deleted\n", name)
		} else {
			fmt.Printf("Error: %s\n", resp.Message)
		}
	},
}

func init() {
	getCmd.AddCommand(getIngressCmd)
	deleteCmd.AddCommand(deleteIngressCmd)
}