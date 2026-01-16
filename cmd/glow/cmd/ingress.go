package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/luaxlou/glow/pkg/api"
	"github.com/spf13/cobra"
)

// --- Shared Logic ---

func runListIngress(cmd *cobra.Command, args []string) {
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
}

func runCreateIngress(cmd *cobra.Command, args []string) {
	name := args[0]
	domain, _ := cmd.Flags().GetString("domain")
	port, _ := cmd.Flags().GetInt("port")
	service, _ := cmd.Flags().GetString("service")

	if domain == "" {
		fmt.Println("Error: --domain is required")
		return
	}

	req := api.IngressUpdateRequest{
		Domain: domain,
		Port:   port,
	}
	if service != "" {
		req.AppName = service
	} else {
		req.AppName = name
	}

	var resp api.Response
	if err := request("POST", "/ingress/update", req, &resp); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if resp.Success {
		fmt.Printf("Ingress %s created\n", name)
	} else {
		fmt.Printf("Error: %s\n", resp.Message)
	}
}

func runDeleteIngress(cmd *cobra.Command, args []string) {
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
}

// --- Verb-Noun Commands (glow get ingress, etc.) ---

var getIngressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "List ingress rules",
	Run:   runListIngress,
}

var createIngressCmd = &cobra.Command{
	Use:   "ingress <name>",
	Short: "Create an ingress rule",
	Args:  cobra.ExactArgs(1),
	Run:   runCreateIngress,
}

var deleteIngressCmd = &cobra.Command{
	Use:   "ingress <name>",
	Short: "Delete an ingress rule",
	Args:  cobra.ExactArgs(1),
	Run:   runDeleteIngress,
}

// --- Noun-Verb Commands (glow ingress list, etc.) ---

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Manage ingress",
}

var ingressListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ingress rules",
	Run:   runListIngress,
}

var ingressApplyCmd = &cobra.Command{
	Use:   "apply <name>",
	Short: "Create or update an ingress rule",
	Args:  cobra.ExactArgs(1),
	Run:   runCreateIngress,
}

var ingressDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an ingress rule",
	Args:  cobra.ExactArgs(1),
	Run:   runDeleteIngress,
}

func init() {
	// Verb-Noun registration
	getCmd.AddCommand(getIngressCmd)
	createCmd.AddCommand(createIngressCmd)
	deleteCmd.AddCommand(deleteIngressCmd)

	// Noun-Verb registration
	rootCmd.AddCommand(ingressCmd)
	ingressCmd.AddCommand(ingressListCmd)
	ingressCmd.AddCommand(ingressApplyCmd)
	ingressCmd.AddCommand(ingressDeleteCmd)

	// Flags for create/apply
	for _, cmd := range []*cobra.Command{createIngressCmd, ingressApplyCmd} {
		cmd.Flags().String("domain", "", "Domain name")
		cmd.Flags().Int("port", 0, "Target port (optional if app running)")
		cmd.Flags().String("service", "", "Target service/app name (optional if same as ingress name)")
	}
}
