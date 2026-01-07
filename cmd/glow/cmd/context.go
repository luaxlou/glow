package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage multiple Glow environments",
}

var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contexts",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}
		
		fmt.Println("CURRENT  NAME       URL")
		for _, ctx := range cfg.Contexts {
			current := " "
			if ctx.Name == cfg.CurrentContext {
				current = "*"
			}
			fmt.Printf("%s        %-10s %s\n", current, ctx.Name, ctx.ServerURL)
		}
	},
}

var contextUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to a different context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg, err := loadConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		found := false
		for _, ctx := range cfg.Contexts {
			if ctx.Name == name {
				found = true
				break
			}
		}
		
		if !found {
			fmt.Printf("Context '%s' not found\n", name)
			return
		}
		
		cfg.CurrentContext = name
		saveConfig(cfg)
		fmt.Printf("Switched to context '%s'\n", name)
	},
}

var contextAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		url, _ := cmd.Flags().GetString("url")
		key, _ := cmd.Flags().GetString("key")
		
		if url == "" || key == "" {
			fmt.Println("Error: --url and --key are required")
			return
		}
		
		cfg, err := loadConfig()
		if err != nil {
			cfg = &Config{}
		}
		
		for _, ctx := range cfg.Contexts {
			if ctx.Name == name {
				fmt.Printf("Context '%s' already exists\n", name)
				return
			}
		}
		
		cfg.Contexts = append(cfg.Contexts, Context{
			Name: name,
			ServerURL: url,
			APIKey: key,
		})
		
		if cfg.CurrentContext == "" {
			cfg.CurrentContext = name
		}
		
saveConfig(cfg)
		fmt.Printf("Context '%s' added\n", name)
	},
}

var contextDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg, err := loadConfig()
		if err != nil {
			return
		}
		
		newContexts := []Context{}
		for _, ctx := range cfg.Contexts {
			if ctx.Name != name {
				newContexts = append(newContexts, ctx)
			}
		}
		
		if len(newContexts) == len(cfg.Contexts) {
			fmt.Printf("Context '%s' not found\n", name)
			return
		}
		
		cfg.Contexts = newContexts
		if cfg.CurrentContext == name {
			cfg.CurrentContext = ""
			if len(cfg.Contexts) > 0 {
				cfg.CurrentContext = cfg.Contexts[0].Name
				fmt.Printf("Active context deleted. Switched to '%s'\n", cfg.CurrentContext)
			} else {
				fmt.Println("Active context deleted. No contexts remaining.")
			}
		}
	saveConfig(cfg)
		fmt.Printf("Context '%s' deleted\n", name)
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextListCmd)
	contextCmd.AddCommand(contextUseCmd)
	contextCmd.AddCommand(contextAddCmd)
	contextCmd.AddCommand(contextDeleteCmd)
	
	contextAddCmd.Flags().String("url", "", "Server URL")
	contextAddCmd.Flags().String("key", "", "API Key")
}
